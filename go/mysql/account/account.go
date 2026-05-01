//
// accounts.go
//
package account

import (
    "database/sql"
    "fmt"
    "strings"
    "time"
    "unicode/utf8"
    _ "github.com/go-sql-driver/mysql"

    myHelper "github.com/k4k3ru-hub/db/go/mysql/helper"
)


const (
    DefaultAccountTableName = "accounts"
)


var (
    accountIDCounter = &myHelper.IdCounter{}
)


type Account struct {
    ID           uint64     `json:"id,string"`
    Status       Status     `json:"status"`
    Role         Role       `json:"role"`
    Name         string     `json:"name"`
    LastLoggedIn *time.Time `json:"lastLoggedIn,omitempty"`
    CreatedAt    time.Time  `json:"createdAt,omitempty"`
    UpdatedAt    time.Time  `json:"updatedAt,omitempty"`
}


type AccountStore struct {
    db        *sql.DB
    tableName string
}


type AccountSelectOption struct {
    Status      *Status `json:"status,omitempty"`
    Role        *Role   `json:"role,omitempty"`
    NameLike    *string `json:"nameLike,omitempty"`
    IDOrLater   *uint64 `json:"idOrLater,omitempty"`
    IDOrEarlier *uint64 `json:"idOrEarlier,omitempty"`
    OrderBy     string  `json:"orderBy"`
    OrderByDesc bool    `json:"orderByDesc"`
    Limit       int     `json:"limit"`
    Offset      int     `json:"offset"`
}


type AccountUpdateOption struct {
    ID           uint64     `json:"id,string"`
    Status       *Status    `json:"status"`
    Role         *Role      `json:"role"`
    Name         *string    `json:"name,omitempty"`
    LastLoggedIn *time.Time `json:"lastLoggedIn,omitempty"`
}


type AccountDeleteOption struct {
    ID uint64 `json:"id,string"`
}


//
// Validate account name.
//
func (a *Account) ValidateName() error {
    // Guard.
    if a == nil {
        return fmt.Errorf("invalid parameter: account=null")
    }

    if a.Name == "" {
        return fmt.Errorf("invalid parameter: password=empty")
    }
    if utf8.RuneCountInString(a.Name) > 64 {
        return fmt.Errorf("invalid parameter: name=%q", a.Name)
    }

    return nil
}


//
// Generate account ID.
//
func GenerateAccountID() uint64 {
    return accountIDCounter.GenerateID()
}


//
// Create new account store.
//
// Version:
//   - 2026-04-30: Added.
//
func NewAccountStore(db *sql.DB, tableName string) (*AccountStore, error) {
    // Guard.
    if db == nil {
        return nil, fmt.Errorf("failed to create account store: missing required parameter: db=null")
    }
    if tableName == "" {
        return nil, fmt.Errorf("failed to create account store: missing required parameter: table_name=empty")
    }

    return &AccountStore{
        db:        db,
        tableName: tableName,
    }, nil
}


//
// Count.
//
// Version:
//   - 2026-04-29: Added.
//
func (s *AccountStore) Count(option *AccountSelectOption) (int64, error) {
    // Guard.
    if s == nil {
        return 0, fmt.Errorf("failed to count accounts: missing required parameter: account_store=null")
    }
    if s.db == nil {
        return 0, fmt.Errorf("failed to count accounts: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return 0, fmt.Errorf("failed to count accounts: missing required parameter: table_name=empty")
    }
    if option == nil {
        return 0, fmt.Errorf("failed to count accounts: missing required parameter: select_option=null")
    }

    // Build conditions.
    conditions, args := buildAccountSelectCondition(option)

    // Generate a SELECT query.
    var query strings.Builder
    query.WriteString("SELECT COUNT(*) FROM ")
    query.WriteString(s.tableName)

    if len(conditions) > 0 {
        query.WriteString(" WHERE ")
        query.WriteString(strings.Join(conditions, " AND "))
    }

    // Execute.
    var result int64
    err := s.db.QueryRow(query.String(), args...).Scan(&result)
    if err != nil {
        return 0, fmt.Errorf("failed to count accounts: %w", err)
    }

    return result, nil
}


//
// Create accounts table.
//
// Version:
//   - 2026-04-29: Added.
//
func (s *AccountStore) CreateTable() error {
    // Guard.
    if s == nil {
        return fmt.Errorf("failed to create accounts table: missing required parameter: account_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to create accounts table: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to create accounts table: missing required parameter: table_name=empty")
    }

    // Generate a CREATE TABLE query.
    query := fmt.Sprintf(
        `CREATE TABLE IF NOT EXISTS %s (
            %s BIGINT UNSIGNED NOT NULL COMMENT 'ID',
            %s TINYINT UNSIGNED NOT NULL COMMENT 'Status',
            %s TINYINT UNSIGNED NOT NULL COMMENT 'Role',
            %s VARCHAR(64) NOT NULL COMMENT 'Name',
            %s DATETIME COMMENT 'Last logged at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
            PRIMARY KEY (%s),
            UNIQUE KEY uk_accounts_name (%s),
            KEY idx_accounts_status (%s),
            KEY idx_accounts_role (%s));
        `,
        s.tableName,
        ColID,
        ColStatus,
        ColRole,
        ColName,
        ColLastLoggedIn,
        ColCreatedAt,
        ColUpdatedAt,
        ColID,
        ColName,
        ColStatus,
        ColRole,
    )

    // Execute the query.
    if _, err := s.db.Exec(query); err != nil {
        return fmt.Errorf("failed to create accounts table: %w", err)
    }

    return nil
}


//
// Delete by ID.
//
// Version:
//   - 2026-04-29: Added.
//
func (s *AccountStore) DeleteByID(id uint64) error {
    // Guard.
    if s == nil {
        return fmt.Errorf("failed to delete account by id: missing required parameter: account_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to delete account by id: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to delete account by id: missing required parameter: table_name=empty")
    }
    if id == 0 {
        return fmt.Errorf("failed to delete account by id: invalid parameter: id=0")
    }

    // Generate a DELETE query.
    query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", s.tableName, ColID)

    // Execute.
    if _, err := s.db.Exec(query, id); err != nil {
        return fmt.Errorf("failed to delete account by id: %w", err)
    }

    return nil
}


//
// Insert an account.
//
// Version:
//   - 2026-04-29: Added.
//
func (s *AccountStore) Insert(row *Account) error {
    // Guard.
    if s == nil {
        return fmt.Errorf("failed to insert account: missing required parameter: account_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to insert account: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to insert account: missing required parameter: table_name=empty")
    }
    if row == nil {
        return fmt.Errorf("failed to insert account: missing required parameter: account=null table=%q", s.tableName)
    }
    if !row.Status.IsValid() {
        return fmt.Errorf("failed to insert account: invalid parameter: status=%d table=%q", row.Status, s.tableName)
    }
    if !row.Role.IsValid() {
        return fmt.Errorf("failed to insert account: invalid parameter: role=%d table=%q", row.Role, s.tableName)
    }
    if err := row.ValidateName(); err != nil {
        return fmt.Errorf("failed to insert account: %w table=%q", err, s.tableName)
    }

    // Generate an INSERT query.
    query := fmt.Sprintf(
        "INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?, ?, ?);",
        s.tableName,
        ColID,
        ColStatus,
        ColRole,
        ColName,
        ColLastLoggedIn,
        ColCreatedAt,
        ColUpdatedAt,
    )

    if row.ID == 0 {
        row.ID = GenerateAccountID()
    }

    now := time.Now()
    if row.CreatedAt.IsZero() {
        row.CreatedAt = now
    }
    if row.UpdatedAt.IsZero() {
        row.UpdatedAt = now
    }

    // Execute.
    if _, err := s.db.Exec(
        query,
        row.ID,
        row.Status,
        row.Role,
        row.Name,
        row.LastLoggedIn,
        row.CreatedAt,
        row.UpdatedAt,
    ); err != nil {
        return fmt.Errorf("failed to insert account: %w", err)
    }

    return nil
}


//
// Select accounts.
//
func (s *AccountStore) Select(option *AccountSelectOption) ([]*Account, error) {
    // Guard.
    if s == nil {
        return nil, fmt.Errorf("failed to select accounts: missing required parameter: account_store=null")
    }
    if s.db == nil {
        return nil, fmt.Errorf("failed to select accounts: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select accounts: missing required parameter: table_name=empty")
    }
    if option == nil {
        return nil, fmt.Errorf("failed to select accounts: missing required parameter: select_option=null")
    }

    // Build conditions.
    conditions, args := buildAccountSelectCondition(option)

    // Generate a SELECT query.
    var query strings.Builder
    query.WriteString("SELECT * FROM ")
    query.WriteString(s.tableName)

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
        if option.Offset < 0 {
            return nil, fmt.Errorf("failed to select accounts: invalid parameter: offset=%d", option.Offset)
        }
        query.WriteString(" LIMIT ? OFFSET ?")
        args = append(args, option.Limit, option.Offset)
    }

    // Execute.
    rows, err := s.db.Query(query.String(), args...)
    if err != nil {
        return nil, fmt.Errorf("failed to select accounts: %w", err)
    }

    defer rows.Close()

    // Scan.
    var result []*Account
    for rows.Next() {
        row := &Account{}
        err := rows.Scan(
            &row.ID,
            &row.Status,
            &row.Role,
            &row.Name,
            &row.LastLoggedIn,
            &row.CreatedAt,
            &row.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to select accounts: %w", err)
        }
        result = append(result, row)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to select accounts: %w", err)
    }

    return result, nil
}


//
// Select account by ID.
//
// Version:
//   - 2026-04-29: Added.
//
func (s *AccountStore) SelectByID(id uint64) (*Account, error) {
    // Guard.
    if s == nil {
        return nil, fmt.Errorf("failed to select account by id: missing required parameter: account_store=null")
    }
    if s.db == nil {
        return nil, fmt.Errorf("failed to select account by id: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account by id: missing required parameter: table_name=empty")
    }
    if id == 0 {
        return nil, fmt.Errorf("failed to select account by id: invalid parameter: id=empty")
    }

    // Generate a SELECT query.
    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColID)

    // Execute.
    row := s.db.QueryRow(query, id)

    // Scan.
    result := &Account{}
    err := row.Scan(
        &result.ID,
        &result.Status,
        &result.Role,
        &result.Name,
        &result.LastLoggedIn,
        &result.CreatedAt,
        &result.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to select account by id: %w", err)
    }

    return result, nil
}


//
// Select by name.
//
// Version:
//   - 2026-04-29: Added.
//
func (s *AccountStore) SelectByName(name string) (*Account, error) {
    // Guard.
    if s == nil {
        return nil, fmt.Errorf("failed to select account by name: missing required parameter: account_store=null")
    }
    if s.db == nil {
        return nil, fmt.Errorf("failed to select account by name: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account by name: missing required parameter: table_name=empty")
    }

    a := Account{
        Name: name,
    }
    if err := a.ValidateName(); err != nil {
        return nil, fmt.Errorf("failed to select account by name: %w", err)
    }

    // Generate a SELECT query.
    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColName)

    // Execute.
    row := s.db.QueryRow(query, name)

    // Scan.
    result := &Account{}
    err := row.Scan(
        &result.ID,
        &result.Status,
        &result.Role,
        &result.Name,
        &result.LastLoggedIn,
        &result.CreatedAt,
        &result.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to select account by name: %w", err)
    }

    return result, nil
}


//
// Update account.
//
// Version:
//   - 2026-04-29: Added.
//
func (s *AccountStore) Update(option *AccountUpdateOption) error {
    // Guard.
    if s == nil {
        return fmt.Errorf("failed to update account: missing required parameter: account_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to update account: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update account: missing required parameter: table_name=empty")
    }
    if option == nil {
        return fmt.Errorf("failed to update account: missing required parameter: option=null")
    }
    if option.ID == 0 {
        return fmt.Errorf("failed to update account: invalid parameter: id=0")
    }

    assignments := make([]string, 0, 4)
    args := make([]interface{}, 0, 4)

    if option.Status != nil {
        if !option.Status.IsValid() {
            return fmt.Errorf("failed to update account: invalid parameter: status=%d", *option.Status)
        }
        assignments = append(assignments, ColStatus + " = ?")
        args = append(args, *option.Status)
    }

    if option.Role != nil {
        if !option.Role.IsValid() {
            return fmt.Errorf("failed to update account: invalid parameter: role=%d", *option.Role)
        }
        assignments = append(assignments, ColRole + " = ?")
        args = append(args, *option.Role)
    }

    if option.Name != nil {
        a := Account{
            Name: *option.Name,
        }
        if err := a.ValidateName(); err != nil {
            return fmt.Errorf("failed to update account: %w", err)
        }
        assignments = append(assignments, ColName + " = ?")
        args = append(args, *option.Name)
    }

    if option.LastLoggedIn != nil {
        assignments = append(assignments, ColLastLoggedIn + " = ?")
        args = append(args, *option.LastLoggedIn)
    }

    if len(assignments) == 0 {
        return fmt.Errorf("failed to update account: invalid parameter: assignments=empty")
    }

    args = append(args, option.ID)

    // Generate a UPDATE query.
    query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?;", s.tableName, strings.Join(assignments, ", "), ColID)

    // Execute.
    if _, err := s.db.Exec(query, args...); err != nil {
        return fmt.Errorf("failed to update account: %w", err)
    }

    return nil
}


//
// Update last logged in.
//
// Version:
//   - 2026-04-29: Added.
//
func (s *AccountStore) UpdateLastLoggedIn(id uint64) error {
    // Guard.
    if s == nil {
        return fmt.Errorf("failed to update account last logged in: missing required parameter: account_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to update account last logged in: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update account last logged in: missing required parameter: table_name=empty")
    }
	if id == 0 {
		return fmt.Errorf("failed to update account last logged in: invalid parameter: id=0")
	}

	now := time.Now()

	return s.Update(&AccountUpdateOption{
		ID:           id,
		LastLoggedIn: &now,
	})
}


//
// Build account select conditions.
//
func buildAccountSelectCondition(option *AccountSelectOption) ([]string, []any) {
    conditions := make([]string, 0, 5)
    args := make([]any, 0, 5)

    if option != nil {
        if option.Status != nil {
            conditions = append(conditions, ColStatus + " = ?")
            args = append(args, *option.Status)
        }
        if option.Role != nil {
            conditions = append(conditions, ColRole + " = ?")
            args = append(args, *option.Role)
        }
        if option.NameLike != nil {
            conditions = append(conditions, ColName + " LIKE ?")
            args = append(args, "%" + *option.NameLike + "%")
        }
        if option.IDOrLater != nil {
            conditions = append(conditions, ColID + " >= ?")
            args = append(args, *option.IDOrLater)
        }
        if option.IDOrEarlier != nil {
            conditions = append(conditions, ColID + " <= ?")
            args = append(args, *option.IDOrEarlier)
        }
    }

    return conditions, args
}
