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


//
// Account.
//
// Version:
//   - 2026-05-09: Added.
//
type Account struct {
    ID           uint64     `json:"id,string"`
    Status       Status     `json:"status"`
    Role         Role       `json:"role"`
    Name         string     `json:"name"`
    LastLoggedIn *time.Time `json:"lastLoggedIn,omitempty"`
    CreatedAt    time.Time  `json:"createdAt"`
    UpdatedAt    time.Time  `json:"updatedAt"`
}


type AccountStore struct {
    executor  myHelper.Executor
    tableName string
}


type AccountSelectOption struct {
    ID          *uint64 `json:"id,string,omitempty"`
    Status      *Status `json:"status,omitempty"`
    Role        *Role   `json:"role,omitempty"`
    NameLike    *string `json:"nameLike,omitempty"`
    IDGTE       *uint64 `json:"idGte,omitempty"`
    IDLTE       *uint64 `json:"idLte,omitempty"`
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
// Validate account ID.
//
// Version:
//   - 2026-05-08: Added.
//
func ValidateAccountID(id uint64) error {
    if id == 0 {
        return fmt.Errorf("invalid parameter: id=0")
    }
    return nil
}


//
// Validate account ID.
//
// Version:
//   - 2026-05-08: Added.
//
func (a *Account) ValidateID() error {
    if a == nil {
        return fmt.Errorf("missing required parameter: account=null")
    }
    return ValidateAccountID(a.ID)
}


//
// Validate account status.
//
// Version:
//   - 2026-05-08: Added.
//
func ValidateAccountStatus(status Status) error {
    if !status.IsValid() {
        return fmt.Errorf("invalid parameter: status=%d", status)
    }
    return nil
}


//
// Validate account status.
//
// Version:
//   - 2026-05-08: Added.
//
func (a *Account) ValidateStatus() error {
    if a == nil {
        return fmt.Errorf("missing required parameter: account=null")
    }
    return ValidateAccountStatus(a.Status)
}


//
// Validate account role.
//
// Version:
//   - 2026-05-08: Added.
//
func ValidateAccountRole(role Role) error {
    if !role.IsValid() {
        return fmt.Errorf("invalid parameter: role=%d", role)
    }
    return nil
}


//
// Validate account role.
//
// Version:
//   - 2026-05-08: Added.
//
func (a *Account) ValidateRole() error {
    if a == nil {
        return fmt.Errorf("missing required parameter: account=null")
    }
    return ValidateAccountRole(a.Role)
}


//
// Validate account name.
//
func ValidateAccountName(name string) error {
    if name == "" {
        return fmt.Errorf("invalid parameter: name=%q", "empty")
    }
    if utf8.RuneCountInString(name) > 64 {
        return fmt.Errorf("invalid parameter: max_length=64 name=%q", "too long")
    }

    return nil
}


//
// Validate account name.
//
func (a *Account) ValidateName() error {
    if a == nil {
        return fmt.Errorf("missing required parameter: account=null")
    }
    return ValidateAccountName(a.Name)
}


//
// Validate account last logged in.
//
func ValidateAccountLastLoggedIn(lastLoggedIn *time.Time) error {
    if lastLoggedIn == nil {
        return nil
    }
    if lastLoggedIn.IsZero() {
        return fmt.Errorf("invalid parameter: last_logged_in=%q", "empty")
    }
    return nil
}


//
// Validate account last logged in.
//
func (a *Account) ValidateLastLoggedIn() error {
    if a == nil {
        return fmt.Errorf("missing required parameter: account=null")
    }
    return ValidateAccountLastLoggedIn(a.LastLoggedIn)
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
func NewAccountStore(executor myHelper.Executor, tableName string) (*AccountStore, error) {
    // Guard.
    if executor == nil {
        return nil, fmt.Errorf("failed to create account store: missing required parameter: executor=null")
    }
    if tableName == "" {
        return nil, fmt.Errorf("failed to create account store: missing required parameter: table_name=%q", "empty")
    }

    return &AccountStore{
        executor:  executor,
        tableName: tableName,
    }, nil
}


//
// Count accounts.
//
// Version:
//   - 2026-04-29: Added.
//
func (s *AccountStore) Count(option *AccountSelectOption) (int64, error) {
    // Guard.
    if s == nil {
        return 0, fmt.Errorf("failed to count accounts: missing required parameter: account_store=null")
    }
    if s.executor == nil {
        return 0, fmt.Errorf("failed to count accounts: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return 0, fmt.Errorf("failed to count accounts: missing required parameter: table_name=%q", "empty")
    }
    if err := option.Validate(); err != nil {
        return 0, fmt.Errorf("failed to count accounts: %w", err)
    }

    query, args := option.BuildQuery("SELECT COUNT(*) FROM " + s.tableName)

    // Execute.
    var result int64
    err := s.executor.QueryRow(query, args...).Scan(&result)
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
    if s.executor == nil {
        return fmt.Errorf("failed to create accounts table: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to create accounts table: missing required parameter: table_name=%q", "empty")
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
        ColStatus,
        ColRole,
    )

    // Execute the query.
    if _, err := s.executor.Exec(query); err != nil {
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
    if s.executor == nil {
        return fmt.Errorf("failed to delete account by id: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to delete account by id: missing required parameter: table_name=%q", "empty")
    }
    if id == 0 {
        return fmt.Errorf("failed to delete account by id: invalid parameter: id=0")
    }

    // Generate a DELETE query.
    query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", s.tableName, ColID)

    // Execute.
    if _, err := s.executor.Exec(query, id); err != nil {
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
    if s.executor == nil {
        return fmt.Errorf("failed to insert account: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to insert account: missing required parameter: table_name=%q", "empty")
    }
    if row == nil {
        return fmt.Errorf("failed to insert account: missing required parameter: account=null table=%q", s.tableName)
    }
    if err := row.ValidateStatus(); err != nil {
        return fmt.Errorf("failed to insert account: %w", err)
    }
    if err := row.ValidateRole(); err != nil {
        return fmt.Errorf("failed to insert account: %w", err)
    }
    if err := row.ValidateName(); err != nil {
        return fmt.Errorf("failed to insert account: %w", err)
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
    if _, err := s.executor.Exec(
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
    if s.executor == nil {
        return nil, fmt.Errorf("failed to select accounts: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select accounts: missing required parameter: table_name=%q", "empty")
    }
    if err := option.Validate(); err != nil {
        return nil, fmt.Errorf("failed to select accounts: %w", err)
    }

    query, args := option.BuildQuery("SELECT * FROM " + s.tableName)

    // Execute.
    rows, err := s.executor.Query(query, args...)
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
    if s.executor == nil {
        return nil, fmt.Errorf("failed to select account by id: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account by id: missing required parameter: table_name=%q", "empty")
    }
    if id == 0 {
        return nil, fmt.Errorf("failed to select account by id: invalid parameter: id=empty")
    }

    // Generate a SELECT query.
    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColID)

    // Execute.
    row := s.executor.QueryRow(query, id)

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
    if s.executor == nil {
        return fmt.Errorf("failed to update account: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update account: missing required parameter: table_name=%q", "empty")
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
        if err := ValidateAccountName(*option.Name); err != nil {
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
        return fmt.Errorf("failed to update account: invalid parameter: assignments=%q", "empty")
    }

    args = append(args, option.ID)

    // Generate a UPDATE query.
    query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?;", s.tableName, strings.Join(assignments, ", "), ColID)

    // Execute.
    if _, err := s.executor.Exec(query, args...); err != nil {
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
    if s.executor == nil {
        return fmt.Errorf("failed to update account last logged in: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update account last logged in: missing required parameter: table_name=%q", "empty")
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
        if option.IDGTE != nil {
            conditions = append(conditions, ColID + " >= ?")
            args = append(args, *option.IDGTE)
        }
        if option.IDLTE != nil {
            conditions = append(conditions, ColID + " <= ?")
            args = append(args, *option.IDLTE)
        }
    }

    return conditions, args
}


//
// Build query.
//
// Version:
//   - 2025-05-09: Added.
//
func (o *AccountSelectOption) BuildQuery(selectFromClause string) (string, []any) {
    // Guard.
    if o == nil {
        return selectFromClause, nil
    }

    var query strings.Builder
    query.WriteString(selectFromClause)

    conditions := make([]string, 0, 6)
    args := make([]any, 0, 8)

    if o.ID != nil {
        conditions = append(conditions, ColID + " = ?")
        args = append(args, *o.ID)
    }
    if o.Status != nil {
        conditions = append(conditions, ColStatus + " = ?")
        args = append(args, *o.Status)
    }
    if o.Role != nil {
        conditions = append(conditions, ColRole + " = ?")
        args = append(args, *o.Role)
    }
    if o.NameLike != nil {
        conditions = append(conditions, ColName + " LIKE ?")
        args = append(args, "%" + *o.NameLike + "%")
    }
    if o.IDGTE != nil {
        conditions = append(conditions, ColID + " >= ?")
        args = append(args, *o.IDGTE)
    }
    if o.IDLTE != nil {
        conditions = append(conditions, ColID + " <= ?")
        args = append(args, *o.IDLTE)
    }

    if len(conditions) > 0 {
        query.WriteString(" WHERE ")
        query.WriteString(strings.Join(conditions, " AND "))
    }

    if o.OrderBy != "" {
        query.WriteString(" ORDER BY ")
        query.WriteString(o.OrderBy)
        if o.OrderByDesc {
            query.WriteString(" DESC")
        }
    }

    if o.Limit > 0 {
        query.WriteString(" LIMIT ? OFFSET ?")
        args = append(args, o.Limit, o.Offset)
    }

    return query.String(), args
}


//
// Validate account email credential select option.
//
// Version:
//   - 2025-05-09: Added.
//
func (o *AccountSelectOption) Validate() error {
    // Guard.
    if o == nil {
        return fmt.Errorf("missing required parameter: account_select_option=null")
    }

    if o.ID != nil {
        if err := ValidateAccountID(*o.ID); err != nil {
            return err
        }
    }
    if o.Status != nil {
        if err := ValidateAccountStatus(*o.Status); err != nil {
            return err
        }
    }
    if o.Role != nil {
        if err := ValidateAccountRole(*o.Role); err != nil {
            return err
        }
    }
    if o.NameLike != nil {
        if err := ValidateAccountName(*o.NameLike); err != nil {
            return err
        }
    }
    if o.IDGTE != nil && o.IDLTE != nil && *o.IDGTE > *o.IDLTE {
        return fmt.Errorf("invalid parameter: id_gte=%d id_lte=%d", *o.IDGTE, *o.IDLTE)
    }

    if o.OrderBy != "" {
        switch o.OrderBy {
        case ColID,
            ColStatus,
            ColRole,
            ColName,
            ColLastLoggedIn,
            ColCreatedAt,
            ColUpdatedAt:
        default:
            return fmt.Errorf("invalid parameter: order_by=%q", o.OrderBy)
        }
    }

    if o.Limit < 0 {
        return fmt.Errorf("invalid parameter: limit=%d", o.Limit)
    }
    if o.Offset < 0 {
        return fmt.Errorf("invalid parameter: offset=%d", o.Offset)
    }

    return nil
}
