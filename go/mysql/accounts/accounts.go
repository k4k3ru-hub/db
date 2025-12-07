//
// accounts.go
//
package account

import (
    "database/sql"
    "database/sql/driver"
    "fmt"
    "regexp"
    "strconv"
    "strings"
    "time"
    "unicode/utf8"
    _ "github.com/go-sql-driver/mysql"

    myHelper "github.com/k4k3ru-hub/db/go/mysql/helper"
)


const (
    ColCreatedAt    = "created_at"
    ColEmail        = "email"
    ColId           = "id"
    ColLastLoggedIn = "last_logged_in"
    ColPassword     = "password"
    ColPublicToken  = "public_token"
    ColRole         = "role"
    ColSecretToken  = "secret_token"
    ColStatus       = "status"
    ColUpdatedAt    = "updated_at"
    ColName         = "name"

    TableName = "accounts"
)


var (
    emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
    idCounter = &myHelper.IdCounter{}
)


type Status uint8
const (
    StatusNew Status = iota
    StatusActived
    StatusInactived
    StatusPending
    StatusSuspended
    StatusDeleted
)
func (s Status) IsValid() bool {
    return s >= StatusNew && s <= StatusDeleted
}
func (s Status) Value() (driver.Value, error) {
    if !s.IsValid() {
        return nil, fmt.Errorf("invalid status: value=%d", s)
    }
    return int64(s), nil
}
func (s *Status) Scan(src any) error {
    switch v := src.(type) {
    case int64:
        *s = Status(v)
    case []byte:
        n, err := strconv.ParseUint(string(v), 10, 8)
        if err != nil {
            return err
        }
        *s = Status(n)
    default:
        return fmt.Errorf("failed to scan: type=%T.", src)
    }
    if !s.IsValid() {
        return fmt.Errorf("invalid status: value=%d", *s)
    }
    return nil
}


type Role uint8
const (
    RoleNone Role = iota
    RoleViewer
    RoleEditor
    RoleAdmin
)
func (r Role) IsValid() bool {
    return r >= RoleNone && r <= RoleAdmin
}
func (r Role) Value() (driver.Value, error) {
    if !r.IsValid() {
        return nil, fmt.Errorf("invalid role: value=%d", r)
    }
    return int64(r), nil
}
func (r *Role) Scan(src any) error {
    switch v := src.(type) {
    case int64:
        *r = Role(v)
    case []byte:
        n, err := strconv.ParseUint(string(v), 10, 8)
        if err != nil {
            return err
        }
        *r = Role(n)
    default:
        return fmt.Errorf("failed to scan: type=%T.", src)
    }
    if !r.IsValid() {
        return fmt.Errorf("invalid role: value=%d", *r)
    }
    return nil
} 


type Account struct {
    Id           uint64
    Status       Status
    Role         Role
    Name         string
    Email        *string
    Password     *string
    PublicToken  *string
    SecretToken  *string
    LastLoggedIn *time.Time
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
type Client struct {
    db         *sql.DB
    tableName  string
    insertStmt *sql.Stmt
}
type InsertOption struct {
    Id           uint64     `json:"id"`
    Status       Status     `json:"status"`
    Role         Role       `json:"role"`
    Name         *string    `json:"name,omitempty"`
    Email        *string    `json:"email,omitempty"`
    Password     *string    `json:"password,omitempty"`
    PublicToken  *string    `json:"publicToken,omitempty"`
    SecretToken  *string    `json:"secretToken,omitempty"`
    LastLoggedIn *time.Time `json:"lastLoggedIn,omitempty"`
    CreatedAt    *time.Time `json:"createdAt,omitempty"`
    UpdatedAt    *time.Time `json:"updatedAt,omitempty"`
}
type SelectOption struct {
    Status          *Status `json:"status,omitempty"`
    Role            *Role   `json:"role,omitempty"`
    NameLike        *string `json:"nameLike,omitempty"`
    Email           *string `json:"email,omitempty"`
    EmailLike       *string `json:"emailLike,omitempty"`
    LastIdOrLater   *uint64 `json:"lastIdOrLater,omitempty"`
    LastIdOrEarlier *uint64 `json:"lastIdOrEarlier,omitempty"`
    OrderBy         string  `json:"orderBy"`
    OrderByDesc     bool    `json:"orderByDesc"`
    Limit           int     `json:"limit"`
    Offset          int     `json:"offset"`
}
type UpdateOption struct {
    Status       *myHelper.UpdateField[Status]    `json:"status,omitempty"`
    Role         *myHelper.UpdateField[Role]      `json:"role,omitempty"`
    Name         *myHelper.UpdateField[string]    `json:"name,omitempty"`
    Email        *myHelper.UpdateField[string]    `json:"email,omitempty"`
    Password     *myHelper.UpdateField[string]    `json:"password,omitempty"`
    PublicToken  *myHelper.UpdateField[string]    `json:",omitempty"`
    SecretToken  *myHelper.UpdateField[string]    `json:",omitempty"`
    LastLoggedIn *myHelper.UpdateField[time.Time] `json:"lastLoggedIn,omitempty"`
    Filter       *UpdateFilter                    `json:"filter,omitempty"`
}
type UpdateFilter struct {
    Primary *UpdateFilterPrimary `json:"primary,omitempty"`
    Unique  *UpdateFilterUnique  `json:"unique,omitempty"`
}
type UpdateFilterPrimary struct {
    Id uint64 `json:"id,string"`
}
type UpdateFilterUnique struct {
    Name string `json:"name"`
}
type DeleteOption struct {
    Id uint64 `json:"id,string"`
}


//
// New Client
//
func NewClient(db *sql.DB, tableName string) *Client {
    return &Client{
        db: db,
        tableName: tableName,
    }
}


//
// New InsertOption
//
func NewInsertOption() *InsertOption {
    return &InsertOption{
        Id: GenerateId(),
    }
}


//
// New SelectOption
//
func NewSelectOption() *SelectOption {
    return &SelectOption{}
}


//
// New UpdateOption
//
func NewUpdateOption() *UpdateOption {
    return &UpdateOption{}
}


//
// Validate email
//
func ValidateEmail(value string) bool {
    return utf8.RuneCountInString(value) <= 64 && emailRegex.MatchString(value)
}


//
// Validate password
//
func ValidatePassword(value string) bool {
    return utf8.RuneCountInString(value) <= 32
}


//
// Validate name
//
func ValidateName(value string) bool {
    return utf8.RuneCountInString(value) <= 32
}


//
// Count
//
func (c *Client) Count(option *SelectOption) (int64, error) {
    // Check if DB is connected.
    if c.db == nil {
        return 0, fmt.Errorf("missing database connection: table=%s", c.tableName)
    }

    // Generate a SELECT query.
    var query strings.Builder
    var conditions []string
    args := make([]interface{}, 0)
    query.WriteString("SELECT COUNT(*) FROM " + c.tableName)

    if option != nil {
        if option.Status != nil {
            conditions = append(conditions, fmt.Sprintf("%s = ?", ColStatus))
            args = append(args, *option.Status)
        }
        if option.Role != nil {
            conditions = append(conditions, fmt.Sprintf("%s = ?", ColRole))
            args = append(args, *option.Role)
        }
        if option.NameLike != nil {
            conditions = append(conditions, fmt.Sprintf("%s LIKE %?%", ColName))
            args = append(args, *option.NameLike)
        }
        if option.Email != nil {
            conditions = append(conditions, fmt.Sprintf("%s = ?", ColEmail))
            args = append(args, *option.Email)
        }
        if option.EmailLike != nil {
            conditions = append(conditions, fmt.Sprintf("%s LIKE %?%", ColEmail))
            args = append(args, *option.EmailLike)
        }
        if option.LastIdOrLater != nil {
            conditions = append(conditions, fmt.Sprintf("%s >= ?", ColId))
            args = append(args, *option.LastIdOrLater)
        }
        if option.LastIdOrEarlier != nil {
            conditions = append(conditions, fmt.Sprintf("%s <= ?", ColId))
            args = append(args, *option.LastIdOrEarlier)
        }

        if len(conditions) > 0 {
            query.WriteString(" WHERE ")
            query.WriteString(strings.Join(conditions, " AND "))
        }
    }

    // Execute.
    var result int64
    err := c.db.QueryRow(query.String(), args...).Scan(&result)
    if err != nil {
        return 0, err
    }

    return result, nil
}


//
// Create table
//
func (c *Client) CreateTable() error {
    // Check if DB is connected.
    if c.db == nil {
        return fmt.Errorf("missing database connection: table=%s", c.tableName)
    }

    // Generate a CREATE TABLE query.
    query := fmt.Sprintf(
        `CREATE TABLE IF NOT EXISTS %s (
            %s BIGINT UNSIGNED NOT NULL COMMENT 'ID',
            %s TINYINT UNSIGNED NOT NULL COMMENT 'Status',
            %s TINYINT UNSIGNED NOT NULL COMMENT 'Role',
            %s VARCHAR(32) NOT NULL COMMENT 'Name',
            %s VARCHAR(64) COMMENT 'Email',
            %s VARCHAR(128) COMMENT 'Password',
            %s VARCHAR(128) COMMENT 'Public token',
            %s VARCHAR(128) COMMENT 'Secret token',
            %s DATETIME COMMENT 'Last logged at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
            PRIMARY KEY (%s),
            UNIQUE KEY (%s),
            INDEX (%s),
            INDEX (%s),
            INDEX (%s));`,
        c.tableName,
        ColId,
        ColStatus,
        ColRole,
        ColName,
        ColEmail,
        ColPassword,
        ColPublicToken,
        ColSecretToken,
        ColLastLoggedIn,
        ColCreatedAt,
        ColUpdatedAt,
        ColId,
        ColName,
        ColEmail,
        ColStatus,
        ColRole,
    )

    // Execute the query.
    if _, err := c.db.Exec(query); err != nil {
        return err
    }

    return nil
}


//
// Delete.
//
func (c *Client) Delete(options []*DeleteOption) error {
    // Check if DB is connected.
    if c.db == nil {
        return fmt.Errorf("missing database connection: table=%s", c.tableName)
    }

    // Generate delete query.
    query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", c.tableName, ColId)

    // Prepare query.
    stmt, err := c.db.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    // Execute query.
    for _, option := range options {
        if option == nil { continue }
        if _, err := stmt.Exec(option.Id); err != nil {
            return err
        }
    }

    return nil
}


//
// Delete by primary key
//
func (c *Client) DeleteByPrimaryKey(id uint64) error {
    // Check if DB is connected.
    if c.db == nil {
        return fmt.Errorf("missing database connection: table=%s", c.tableName)
    }

    // Generate a query.
    query := "DELETE FROM " + c.tableName + " WHERE " + ColId + " = ?"

    // Execute.
    if _, err := c.db.Exec(query, id); err != nil {
        return err
    }

    return nil
}


//
// Insert
//
func (c *Client) Insert(option *InsertOption) error {
    // Check if DB is connected.
    if c.db == nil {
        return fmt.Errorf("missing database connection: table=%s", c.tableName)
    }

    // Set SQl query statement.
    if c.insertStmt == nil {
        var err error
        c.insertStmt, err = c.db.Prepare(
            fmt.Sprintf(
                `INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
                c.tableName,
                ColId,
                ColStatus,
                ColRole,
                ColName,
                ColEmail,
                ColPassword,
                ColPublicToken,
                ColSecretToken,
                ColLastLoggedIn,
                ColCreatedAt,
                ColUpdatedAt,
            ),
        )
        if err != nil {
            return err
        }
    }

    // Generate an ID.
    if option.Id == 0 {
        option.Id = GenerateId()
    }

    // Set default values.
    now := time.Now()
    if option.CreatedAt == nil {
        option.CreatedAt = &now
    }
    if option.UpdatedAt == nil {
        option.UpdatedAt = &now
    }

    // Execute.
    _, err := c.insertStmt.Exec(
        option.Id,
        option.Status,
        option.Role,
        option.Name,
        option.Email,
        option.Password,
        option.PublicToken,
        option.SecretToken,
        option.LastLoggedIn,
        option.CreatedAt,
        option.UpdatedAt,
    )

    return err
}


//
// Select
//
func (c *Client) Select(option *SelectOption) ([]*Account, error) {
    // Check if DB is connected.
    if c.db == nil {
        return nil, fmt.Errorf("missing database connection: table=%s", c.tableName)
    }

    // Generate a SELECT query.
    var query strings.Builder
    var conditions []string
    args := make([]interface{}, 0)
    query.WriteString("SELECT * FROM " + c.tableName)

    if option != nil {
        if option.Status != nil {
            conditions = append(conditions, fmt.Sprintf("%s = ?", ColStatus))
            args = append(args, *option.Status)
        }
        if option.Role != nil {
            conditions = append(conditions, fmt.Sprintf("%s = ?", ColRole))
            args = append(args, *option.Role)
        }
        if option.NameLike != nil {
            conditions = append(conditions, fmt.Sprintf("%s LIKE %?%", ColName))
            args = append(args, *option.NameLike)
        }
        if option.Email != nil {
            conditions = append(conditions, fmt.Sprintf("%s = ?", ColEmail))
            args = append(args, *option.Email)
        }
        if option.EmailLike != nil {
            conditions = append(conditions, fmt.Sprintf("%s LIKE %?%", ColEmail))
            args = append(args, *option.EmailLike)
        }
        if option.LastIdOrLater != nil {
            conditions = append(conditions, fmt.Sprintf("%s >= ?", ColId))
            args = append(args, *option.LastIdOrLater)
        }
        if option.LastIdOrEarlier != nil {
            conditions = append(conditions, fmt.Sprintf("%s <= ?", ColId))
            args = append(args, *option.LastIdOrEarlier)
        }

        if len(conditions) > 0 {
            query.WriteString(" WHERE ")
            query.WriteString(strings.Join(conditions, " AND "))
        }

        if option.OrderBy != "" {
            query.WriteString(" ORDER BY " + option.OrderBy)
            if option.OrderByDesc {
                query.WriteString(" DESC")
            }
        }
        if option.Limit > 0 {
            query.WriteString(" LIMIT " + strconv.Itoa(option.Limit))
        }
        if option.Offset > 0 {
            query.WriteString(" , " + strconv.Itoa(option.Offset))
        }
    }

    // Execute.
    rows, err := c.db.Query(query.String(), args...)
    if err != nil {
        return nil, err
    }

    defer rows.Close()

    var result []*Account
    for rows.Next() {
        row := &Account{}
        err := rows.Scan(
            &row.Id,
            &row.Status,
            &row.Role,
            &row.Name,
            &row.Email,
            &row.Password,
            &row.PublicToken,
            &row.SecretToken,
            &row.LastLoggedIn,
            &row.CreatedAt,
            &row.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        result = append(result, row)
    }

    return result, nil
}


//
// Select by primary key
//
func (c *Client) SelectByPrimaryKey(id uint64) (*Account, error) {
    // Check if DB is connected.
    if c.db == nil {
        return nil, fmt.Errorf("missing database connection: table=%s", c.tableName)
    }

    // Generate a select query.
    query := "SELECT * FROM " + c.tableName + " WHERE " + ColId + " = ? LIMIT 1"

    // Execute.
    row := c.db.QueryRow(query, id)

    // Scan.
    result := &Account{}
    err := row.Scan(
        &result.Id,
        &result.Status,
        &result.Role,
        &result.Name,
        &result.Email,
        &result.Password,
        &result.PublicToken,
        &result.SecretToken,
        &result.LastLoggedIn,
        &result.CreatedAt,
        &result.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }

    return result, nil
}


//
// Select by unique key.
//
func (c *Client) SelectByUniqueKey(name string) (*Account, error) {
    // Check if DB is connected.
    if c.db == nil {
        return nil, fmt.Errorf("missing database connection: table=%s", c.tableName)
    }

    // Generate a select query.
    query := "SELECT * FROM " + c.tableName + " WHERE " + ColName + " = ? LIMIT 1"

    // Execute.
    row := c.db.QueryRow(query, name)

    // Scan.
    result := &Account{}
    err := row.Scan(
        &result.Id,
        &result.Status,
        &result.Role,
        &result.Name,
        &result.Email,
        &result.Password,
        &result.PublicToken,
        &result.SecretToken,
        &result.LastLoggedIn,
        &result.CreatedAt,
        &result.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }

    return result, nil
}


//
// Update
//
func (c *Client) Update(id uint64, option *UpdateOption) error {
    // Check if DB is connected.
    if c.db == nil {
        return fmt.Errorf("missing database connection: table=%s", c.tableName)
    }

    // Check options.
    if option == nil || option.Filter == nil {
        return fmt.Errorf("missing update options")
    }
    if option.Filter.Primary == nil {
        return fmt.Errorf("missing required only one of filter.primary")
    }

    // Generate a update query.
    query := "UPDATE " + c.tableName + " SET "
    var assignmentList []string
    setArgs := make([]interface{}, 0)

    myHelper.AppendUpdateAssignment(&assignmentList, &setArgs, ColStatus, option.Status)
    myHelper.AppendUpdateAssignment(&assignmentList, &setArgs, ColRole, option.Role)
    myHelper.AppendUpdateAssignment(&assignmentList, &setArgs, ColName, option.Name)
    myHelper.AppendUpdateAssignment(&assignmentList, &setArgs, ColEmail, option.Email)
    myHelper.AppendUpdateAssignment(&assignmentList, &setArgs, ColPassword, option.Password)
    myHelper.AppendUpdateAssignment(&assignmentList, &setArgs, ColPublicToken, option.PublicToken)
    myHelper.AppendUpdateAssignment(&assignmentList, &setArgs, ColSecretToken, option.SecretToken)
    myHelper.AppendUpdateAssignment(&assignmentList, &setArgs, ColLastLoggedIn, option.LastLoggedIn)

    // Check assignmentList.
    if len(assignmentList) == 0 {
        return fmt.Errorf("missing columns to update")
    }

    // Generate query for conditions.
    var conditions []string
    var conditionArgs []interface{}
    if option.Filter.Primary != nil && option.Filter.Primary.Id != 0 {
        conditions = append(conditions, fmt.Sprintf("%s = ?", ColId))
        conditionArgs = append(conditionArgs, option.Filter.Primary.Id)
    }
    if option.Filter.Unique != nil && option.Filter.Unique.Name != "" {
        conditions = append(conditions, fmt.Sprintf("%s = ?", ColName))
        conditionArgs = append(conditionArgs, option.Filter.Unique.Name)
    }
    if len(conditions) == 0 {
        return fmt.Errorf("missing conditions to update.")
    }

    // Execute.
    args := append(setArgs, conditionArgs...)
    query += strings.Join(assignmentList, ", ") + " WHERE " + strings.Join(conditions, " AND ")
    _, err := c.db.Exec(query, args...)

    return err
}


//
// Has duplicate for update.
//
func (c *Client) HasDuplicateForUpdate(option *UpdateOption) (bool, error) {
    // Check required values.
    if option == nil || option.Filter == nil {
        return false, fmt.Errorf("missing options to update.")
    }
    if (option.Filter.Primary == nil) == (option.Filter.Unique == nil) {
        return false, fmt.Errorf("missing required only one of filter.primary or filter.unique.")
    }

    // Retrieve account by primary key or unique key.
    var account *Account
    var err error
    if option.Filter.Primary != nil && option.Filter.Primary.Id != 0 {
        account, err = c.SelectByPrimaryKey(option.Filter.Primary.Id)
    } else if option.Filter.Unique != nil && option.Filter.Unique.Name != "" {
        account, err = c.SelectByUniqueKey(option.Filter.Unique.Name)
    } else {
        return false, fmt.Errorf("missing conditions to update.")
    }
    if err != nil {
        return false, err
    }
    if account == nil {
        return false, nil
    }

    // Check if unique key is able to update or not.
    if (option.Name == nil || account.Name == option.Name.Value) {
        return false, nil
    }

    // Retrieve account by specified unique key.
    specifiedAccount, err := c.SelectByUniqueKey(option.Name.Value)
    if err != nil {
        return false, err
    }
    if specifiedAccount != nil && specifiedAccount.Id != account.Id {
        return true, nil
    }

    return false, nil
}


//
// Update last logged In
//
func (c *Client) UpdateLastLoggedIn(id uint64) error {
    option := NewUpdateOption()
    lastLoggedIn := time.Now()
    option.LastLoggedIn = &myHelper.UpdateField[time.Time]{Value:lastLoggedIn}
    return c.Update(id, option)
}


//
// Generate an ID.
//
func GenerateId() uint64 {
    return myHelper.GenerateId(idCounter)
}
