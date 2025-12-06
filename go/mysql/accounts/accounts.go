//
// accounts.go
//
package account

import (
    "database/sql"
    "fmt"
    "regexp"
    "strconv"
    "strings"
    "sync"
    "time"
    "unicode/utf8"
    _ "github.com/go-sql-driver/mysql"

    myHelper "github.com/k4k3ru-hub/db/go/mysql/helper"
)


const (
    ColCountry          = "country"
    ColCreatedAt        = "created_at"
    ColEmail            = "email"
    ColId               = "id"
    ColLastLoggedIn     = "last_logged_in"
    ColPassword         = "password"
    ColPublishableToken = "publishable_token"
    ColRole             = "role"
    ColSecretToken      = "secret_token"
    ColStatus           = "status"
    ColUpdatedAt        = "updated_at"
    ColUsername         = "username"

    TableName = "accounts"

    ValRoleNone   = "none"
    ValRoleAdmin  = "admin"
    ValRoleEditor = "editor"
    ValRoleViewer = "viewer"

    ValStatusActive    = "active"
    ValStatusInactive  = "inactive"
    ValStatusPending   = "pending"
    ValStatusSuspended = "suspended"
    ValStatusDeleted   = "deleted"
)


var (
    emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
    idCounter = &myHelper.IdCounter{}
    insertStmt *sql.Stmt
)


type Account struct {
    Id               uint64
    Status           string
    Role             string
    Username         sql.NullString
    Email            sql.NullString
    Country          sql.NullString
    Password         sql.NullString
    PublishableToken sql.NullString
    SecretToken      sql.NullString
    LastLoggedIn     sql.NullTime
    CreatedAt        time.Time
    UpdatedAt        time.Time
}
type Client struct {
    Db        *sql.DB
    TableName string
}
type IdCounter struct {
    mu sync.Mutex
    id uint64
}
type InsertOption struct {
    Id               uint64
    Status           string
    Role             string
    Username         *string
    Email            *string
    Country          *string
    Password         *string
    PublishableToken *string
    SecretToken      *string
    LastLoggedIn     *time.Time
}
type SelectOption struct {
    Status          *string
    Role            *string
    UsernameLike    *string
    Email           *string
    EmailLike       *string
    Country         *string
    LastIdOrLater   *uint64
    LastIdOrEarlier *uint64
    OrderBy         string
    OrderByDesc     bool
    Limit           int
    Offset          int
}
type UpdateOption struct {
    Status           *string
    Role             *string
    Username         *string
    Email            *string
    Country          *string
    Password         *string
    PublishableToken *string
    SecretToken      *string
    LastLoggedIn     *time.Time
}


//
// New Client
//
func NewClient(db *sql.DB, tableName string) *Client {
    return &Client{
        Db: db,
        TableName: tableName,
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
// Validate country
//
func ValidateCountry(value string) bool {
    return utf8.RuneCountInString(value) == 2
}


//
// Validate password
//
func ValidatePassword(value string) bool {
    return utf8.RuneCountInString(value) <= 32
}


//
// Validate role
//
func ValidateRole(value string) bool {
    validStatuses := map[string]struct{}{
        ValRoleNone: {},
        ValRoleAdmin: {},
        ValRoleEditor: {},
        ValRoleViewer: {},
    }
    _, exists := validStatuses[value]
    return exists
}


//
// Validate status
//
func ValidateStatus(value string) bool {
    validStatuses := map[string]struct{}{
        ValStatusActive: {},
        ValStatusInactive: {},
        ValStatusPending: {},
        ValStatusSuspended: {},
        ValStatusDeleted: {},
    }
    _, exists := validStatuses[value]
    return exists
}


//
// Validate username
//
func ValidateUsername(value string) bool {
    return utf8.RuneCountInString(value) <= 32
}


//
// Count
//
func (c *Client) Count(option *SelectOption) (int64, error) {
    // Check if DB is connected.
    if c.Db == nil {
        return 0, fmt.Errorf("Failed to execute select query because DB was disconnected. (table: %s)\n", c.TableName)
    }

    // Generate a SELECT query.
    var query strings.Builder
    args := make([]interface{}, 0)
    query.WriteString("SELECT COUNT(*) FROM " + c.TableName)

    if option != nil {
        isWhere := false
        if option.Status != nil {
            query.WriteString(" WHERE " + ColStatus + " = ?")
            args = append(args, option.Status)
            isWhere = true
        }
        if option.Role != nil {
            if isWhere {
                query.WriteString(" AND " + ColRole + " = ?")
            } else {
                query.WriteString(" WHERE " + ColRole + " = ?")
                isWhere = true
            }
            args = append(args, option.Role)
        }
        if option.UsernameLike != nil {
            if isWhere {
                query.WriteString(" AND " + ColUsername + " LIKE ?")
            } else {
                query.WriteString(" WHERE " + ColUsername + " LIKE ?")
                isWhere = true
            }
            args = append(args, "%" + *option.UsernameLike + "%")
        }
        if option.Email != nil {
            if isWhere {
                query.WriteString(" AND " + ColEmail + " = ?")
            } else {
                query.WriteString(" WHERE " + ColEmail + " = ?")
                isWhere = true
            }
            args = append(args, option.Email)
        }
        if option.EmailLike != nil {
            if isWhere {
                query.WriteString(" AND " + ColEmail + " LIKE ?")
            } else {
                query.WriteString(" WHERE " + ColEmail + " LIKE ?")
                isWhere = true
            }
            args = append(args, "%" + *option.EmailLike + "%")
        }
        if option.Country != nil {
            if isWhere {
                query.WriteString(" AND " + ColCountry + " = ?")
            } else {
                query.WriteString(" WHERE " + ColCountry + " = ?")
                isWhere = true
            }
            args = append(args, option.Country)
        }
        if option.LastIdOrLater != nil {
            if isWhere {
                query.WriteString(" AND " + ColId + " >= ?")
            } else {
                query.WriteString(" WHERE " + ColId + " >= ?")
                isWhere = true
            }
            args = append(args, option.LastIdOrLater)
        }
        if option.LastIdOrEarlier != nil {
            if isWhere {
                query.WriteString(" AND " + ColId + " <= ?")
            } else {
                query.WriteString(" WHERE " + ColId + " <= ?")
                isWhere = true
            }
            args = append(args, option.LastIdOrEarlier)
        }
    }

    // Execute.
    var result int64
    err := c.Db.QueryRow(query.String(), args...).Scan(&result)
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
    if c.Db == nil {
        return fmt.Errorf("Failed to execute create table query because DB was disconnected. (table: %s)\n", c.TableName)
    }

    // Generate a CREATE TABLE query.
    query := fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s (
        %s BIGINT UNSIGNED NOT NULL COMMENT 'ID',
        %s ENUM('%s', '%s', '%s', '%s', '%s') NOT NULL DEFAULT '%s' COMMENT 'Status',
        %s ENUM('%s', '%s', '%s', '%s') NOT NULL DEFAULT '%s' COMMENT 'Role',
        %s VARCHAR(32) COMMENT 'Username',
        %s VARCHAR(64) COMMENT 'Email',
        %s VARCHAR(2) COMMENT 'Country code',
        %s VARCHAR(128) COMMENT 'Password',
        %s VARCHAR(128) COMMENT 'Publishable token',
        %s VARCHAR(128) COMMENT 'Secret token',
        %s DATETIME COMMENT 'Last logged at',
        %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
        %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
        PRIMARY KEY (%s),
        UNIQUE KEY (%s),
        INDEX (%s),
        INDEX (%s),
        INDEX (%s));`,
        c.TableName,
        ColId,
        ColStatus, ValStatusActive, ValStatusInactive, ValStatusPending, ValStatusSuspended, ValStatusDeleted, ValStatusActive,
        ColRole, ValRoleNone, ValRoleAdmin, ValRoleEditor, ValRoleViewer, ValRoleNone,
        ColUsername,
        ColEmail,
        ColCountry,
        ColPassword,
        ColPublishableToken,
        ColSecretToken,
        ColLastLoggedIn,
        ColCreatedAt,
        ColUpdatedAt,
        ColId,
        ColEmail,
        ColStatus,
        ColRole,
        ColUsername)

    // Execute the query.
    if _, err := c.Db.Exec(query); err != nil {
        return err
    }

    return nil
}


//
// Delete by primary key
//
func (c *Client) DeleteByPrimaryKey(id uint64) error {
    // Check if DB is connected.
    if c.Db == nil {
        return fmt.Errorf("Failed to execute delete query because DB was disconnected. (table: %s)\n", c.TableName)
    }

    // Generate a select query.
    query := "DELETE FROM " + c.TableName + " WHERE " + ColId + " = ?"

    // Execute.
    if _, err := c.Db.Exec(query, id); err != nil {
        return err
    }

    return nil
}


//
// Insert
//
func (c *Client) Insert(option *InsertOption) error {
    // Check if DB is connected.
    if c.Db == nil {
        return fmt.Errorf("Failed to execute insert query because DB was disconnected. (table: %s)\n", c.TableName)
    }

    // Set SQl query statement.
    if insertStmt == nil {
        var err error
        insertStmt, err = c.Db.Prepare(
            fmt.Sprintf(
                `INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
                c.TableName,
                ColId,
                ColStatus,
                ColRole,
                ColUsername,
                ColEmail,
                ColCountry,
                ColPassword,
                ColPublishableToken,
                ColSecretToken,
                ColLastLoggedIn))
        if err != nil {
            return err
        }
    }

    // Generate an ID.
    if option.Id == 0 {
        option.Id = GenerateId()
    }

    // Execute.
    _, err := insertStmt.Exec(
        option.Id,
        option.Status,
        option.Role,
        option.Username,
        option.Email,
        option.Country,
        option.Password,
        option.PublishableToken,
        option.SecretToken,
        option.LastLoggedIn,
    )

    return err
}


//
// Select
//
func (c *Client) Select(option *SelectOption) ([]*Account, error) {
    // Check if DB is connected.
    if c.Db == nil {
        return nil, fmt.Errorf("Failed to execute select query because DB was disconnected. (table: %s)\n", c.TableName)
    }

    // Generate a SELECT query.
    var query strings.Builder
    args := make([]interface{}, 0)
    query.WriteString("SELECT * FROM " + c.TableName)

    if option != nil {
        isWhere := false
        if option.Status != nil {
            query.WriteString(" WHERE " + ColStatus + " = ?")
            args = append(args, option.Status)
            isWhere = true
        }
        if option.Role != nil {
            if isWhere {
                query.WriteString(" AND " + ColRole + " = ?")
            } else {
                query.WriteString(" WHERE " + ColRole + " = ?")
                isWhere = true
            }
            args = append(args, option.Role)
        }
        if option.UsernameLike != nil {
            if isWhere {
                query.WriteString(" AND " + ColUsername + " LIKE ?")
            } else {
                query.WriteString(" WHERE " + ColUsername + " LIKE ?")
                isWhere = true
            }
            args = append(args, "%" + *option.UsernameLike + "%")
        }
        if option.Email != nil {
            if isWhere {
                query.WriteString(" AND " + ColEmail + " = ?")
            } else {
                query.WriteString(" WHERE " + ColEmail + " = ?")
                isWhere = true
            }
            args = append(args, option.Email)
        }
        if option.EmailLike != nil {
            if isWhere {
                query.WriteString(" AND " + ColEmail + " LIKE ?")
            } else {
                query.WriteString(" WHERE " + ColEmail + " LIKE ?")
                isWhere = true
            }
            args = append(args, "%" + *option.EmailLike + "%")
        }
        if option.Country != nil {
            if isWhere {
                query.WriteString(" AND " + ColCountry + " = ?")
            } else {
                query.WriteString(" WHERE " + ColCountry + " = ?")
                isWhere = true
            }
            args = append(args, option.Country)
        }
        if option.LastIdOrLater != nil {
            if isWhere {
                query.WriteString(" AND " + ColId + " >= ?")
            } else {
                query.WriteString(" WHERE " + ColId + " >= ?")
                isWhere = true
            }
            args = append(args, option.LastIdOrLater)
        }
        if option.LastIdOrEarlier != nil {
            if isWhere {
                query.WriteString(" AND " + ColId + " <= ?")
            } else {
                query.WriteString(" WHERE " + ColId + " <= ?")
                isWhere = true
            }
            args = append(args, option.LastIdOrEarlier)
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
    rows, err := c.Db.Query(query.String(), args...)
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
            &row.Username,
            &row.Email,
            &row.Country,
            &row.Password,
            &row.PublishableToken,
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
    if c.Db == nil {
        return nil, fmt.Errorf("Failed to execute select query because DB was disconnected. (table: %s)\n", c.TableName)
    }

    // Generate a select query.
    query := "SELECT * FROM " + c.TableName + " WHERE " + ColId + " = ? LIMIT 1"

    // Execute.
    result := &Account{}
    if err := c.Db.QueryRow(query, id).Scan(
        &result.Id,
        &result.Status,
        &result.Role,
        &result.Username,
        &result.Email,
        &result.Country,
        &result.Password,
        &result.PublishableToken,
        &result.SecretToken,
        &result.LastLoggedIn,
        &result.CreatedAt,
        &result.UpdatedAt); err != nil {
            if err == sql.ErrNoRows {
                return nil, nil
            }
            return nil, err
        }

    return result, nil
}


//
// Select by email
//
func (c *Client) SelectByEmail(email string) (*Account, error) {
    // Check if DB is connected.
    if c.Db == nil {
        return nil, fmt.Errorf("Failed to execute select query because DB was disconnected. (table: %s)\n", c.TableName)
    }

    // Generate a select query.
    query := "SELECT * FROM " + c.TableName + " WHERE " + ColEmail + " = ? LIMIT 1"

    // Execute.
    result := &Account{}
    if err := c.Db.QueryRow(query, email).Scan(
        &result.Id,
        &result.Status,
        &result.Role,
        &result.Username,
        &result.Email,
        &result.Country,
        &result.Password,
        &result.PublishableToken,
        &result.SecretToken,
        &result.LastLoggedIn,
        &result.CreatedAt,
        &result.UpdatedAt); err != nil {
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
    if c.Db == nil {
        return fmt.Errorf("Failed to execute insert query because DB was disconnected. (table: %s)\n", c.TableName)
    }

    // Generate a update query.
    query := "UPDATE " + c.TableName + " SET "
    var assignmentList []string
    args := make([]interface{}, 0)

    if option.Status != nil {
        assignmentList = append(assignmentList, ColStatus + " = ?")
        args = append(args, option.Status)
    }
    if option.Role != nil {
        assignmentList = append(assignmentList, ColRole + " = ?")
        args = append(args, option.Role)
    }
    if option.Username != nil {
        assignmentList = append(assignmentList, ColUsername + " = ?")
        args = append(args, option.Username)
    }
    if option.Email != nil {
        assignmentList = append(assignmentList, ColEmail + " = ?")
        args = append(args, option.Email)
    }
    if option.Country != nil {
        assignmentList = append(assignmentList, ColCountry + " = ?")
        args = append(args, option.Country)
    }
    if option.Password != nil {
        assignmentList = append(assignmentList, ColPassword + " = ?")
        args = append(args, option.Password)
    }
    if option.PublishableToken != nil {
        assignmentList = append(assignmentList, ColPublishableToken + " = ?")
        args = append(args, option.PublishableToken)
    }
    if option.SecretToken != nil {
        assignmentList = append(assignmentList, ColSecretToken + " = ?")
        args = append(args, option.SecretToken)
    }
    if option.LastLoggedIn != nil {
        assignmentList = append(assignmentList, ColLastLoggedIn + " = ?")
        args = append(args, (*option.LastLoggedIn).Format("2006-01-02 15:04:05"))
    }

    // Execute.
    query += strings.Join(assignmentList, ", ") + " WHERE " + ColId + " = ?"
    args = append(args, id)
    _, err := c.Db.Exec(query, args...)

    return err
}


//
// Update last logged In
//
func (c *Client) UpdateLastLoggedIn(id uint64) error {
    option := NewUpdateOption()
    lastLoggedIn := time.Now()
    option.LastLoggedIn = &lastLoggedIn
    return c.Update(id, option)
}


//
// Generate an ID.
//
func GenerateId() uint64 {
    return myHelper.GenerateId(idCounter)
}
