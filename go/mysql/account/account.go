//
// accounts.go
//
package account

import (
    "database/sql"
    "errors"
    "fmt"
    "strings"
    "time"
    "unicode/utf8"

    "github.com/go-sql-driver/mysql"

    "github.com/k4k3ru-hub/db/go/mysql/helper"
)


const (
    DefaultAccountTableName = "accounts"
)


var (
    accountIDCounter = &helper.IdCounter{}
)


//
// Account.
//
// Version:
//   - 2026-05-09: Added.
//
type Account struct {
    ID        uint64    `json:"id,string"`
    Status    Status    `json:"status"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}


type AccountStore struct {
    tableName string
}

type AccountInsertParams struct {
    ID        uint64    `json:"id,string"`
    Status    Status    `json:"status"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
    Ignore    bool      `json:"ignore"`
}

type AccountSelectOption struct {
    ID          *uint64 `json:"id,string,omitempty"`
    Status      *Status `json:"status,omitempty"`
    NameLike    *string `json:"nameLike,omitempty"`
    IDGTE       *uint64 `json:"idGte,omitempty"`
    IDLTE       *uint64 `json:"idLte,omitempty"`
    OrderBy     string  `json:"orderBy"`
    OrderByDesc bool    `json:"orderByDesc"`
    Limit       int     `json:"limit"`
    Offset      int     `json:"offset"`
}


type AccountUpdateOption struct {
    Status *Status `json:"status,omitempty"`
    Name   *string `json:"name,omitempty"`
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
// Validate account name.
//
func ValidateAccountName(name string) error {
    if name == "" {
        return fmt.Errorf("invalid parameter: name=%q", "empty")
    }
    if utf8.RuneCountInString(name) > 64 {
        return fmt.Errorf("invalid parameter: name=%q max_length=64", "too long")
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
// Generate account ID.
//
// Version:
//   - 2026-04-30: Added.
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
func NewAccountStore(tableName string) (*AccountStore, error) {
    // Guard.
    tableName = strings.TrimSpace(tableName)
    if tableName == "" {
        return nil, fmt.Errorf("failed to create new account store: missing required parameter: table_name=%q", "empty")
    }

    return &AccountStore{
        tableName: tableName,
    }, nil
}


//
// Count accounts.
//
// Version:
//   - 2026-04-29: Added.
//
func (s *AccountStore) Count(executor helper.Executor, option *AccountSelectOption) (int64, error) {
    // Guard.
    if s == nil {
        return 0, fmt.Errorf("failed to count accounts: missing required parameter: account_store=null")
    }
    if s.tableName == "" {
        return 0, fmt.Errorf("failed to count accounts: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return 0, fmt.Errorf("failed to count accounts: missing required parameter: executor=null")
    }
    if err := option.Validate(); err != nil {
        return 0, fmt.Errorf("failed to count accounts: %w", err)
    }

    query, args := option.BuildQuery("SELECT COUNT(*) FROM " + s.tableName)

    // Execute query.
    var result int64
    err := executor.QueryRow(query, args...).Scan(&result)
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
func (s *AccountStore) CreateTable(executor helper.Executor) error {
    // Guard.
    if s == nil {
        return fmt.Errorf("failed to create accounts table: missing required parameter: account_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to create accounts table: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to create accounts table: missing required parameter: executor=null")
    }

    // Generate CREATE TABLE query.
    query := fmt.Sprintf(
        `CREATE TABLE IF NOT EXISTS %s (
            %s BIGINT UNSIGNED NOT NULL COMMENT 'ID',
            %s TINYINT UNSIGNED NOT NULL COMMENT 'Status',
            %s VARCHAR(64) NOT NULL COMMENT 'Name',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
            PRIMARY KEY (%s),
            KEY idx_accounts_status (%s)
        );`,
        s.tableName,
        ColID,
        ColStatus,
        ColName,
        ColCreatedAt,
        ColUpdatedAt,
        ColID,
        ColStatus,
    )

    // Execute query.
    if _, err := executor.Exec(query); err != nil {
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
func (s *AccountStore) DeleteByID(executor helper.Executor, id uint64) error {
    // Guard.
    if s == nil {
        return fmt.Errorf("failed to delete account by id: missing required parameter: account_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to delete account by id: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to delete account by id: missing required parameter: executor=null")
    }
    if id == 0 {
        return fmt.Errorf("failed to delete account by id: invalid parameter: id=0")
    }

    // Generate DELETE query.
    query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", s.tableName, ColID)

    // Execute query.
    if _, err := executor.Exec(query, id); err != nil {
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
func (s *AccountStore) Insert(executor helper.Executor, params *AccountInsertParams) error {
    // Guard.
    if s == nil {
        return fmt.Errorf("failed to insert account: missing required parameter: account_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to insert account: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to insert account: missing required parameter: executor=null")
    }
    if params == nil {
        return fmt.Errorf("failed to insert account: missing required parameter: account_insert_params=null")
    }
    if err := ValidateAccountStatus(params.Status); err != nil {
        return fmt.Errorf("failed to insert account: %w", err)
    }
    if err := ValidateAccountName(params.Name); err != nil {
        return fmt.Errorf("failed to insert account: %w", err)
    }

    // Generate INSERT query.
    query := fmt.Sprintf(
        "INSERT INTO %s (%s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?);",
        s.tableName,
        ColID,
        ColStatus,
        ColName,
        ColCreatedAt,
        ColUpdatedAt,
    )

    if params.ID == 0 {
        params.ID = GenerateAccountID()
    }

    now := time.Now()
    if params.CreatedAt.IsZero() {
        params.CreatedAt = now
    }
    if params.UpdatedAt.IsZero() {
        params.UpdatedAt = now
    }

    // Execute query.
    if _, err := executor.Exec(
        query,
        params.ID,
        params.Status,
        params.Name,
        params.CreatedAt,
        params.UpdatedAt,
    ); err != nil {
        var mysqlErr *mysql.MySQLError
        if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
            return fmt.Errorf("failed to insert account: %w", helper.ErrDuplicateKey)
        }
        return fmt.Errorf("failed to insert account: %w", err)
    }

    return nil
}


//
// Select accounts.
//
func (s *AccountStore) Select(executor helper.Executor, option *AccountSelectOption) ([]*Account, error) {
    // Guard.
    if s == nil {
        return nil, fmt.Errorf("failed to select accounts: missing required parameter: account_store=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select accounts: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return nil, fmt.Errorf("failed to select accounts: missing required parameter: executor=null")
    }
    if err := option.Validate(); err != nil {
        return nil, fmt.Errorf("failed to select accounts: %w", err)
    }

    query, args := option.BuildQuery("SELECT * FROM " + s.tableName)

    // Execute query.
    rows, err := executor.Query(query, args...)
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
            &row.Name,
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
func (s *AccountStore) SelectByID(executor helper.Executor, id uint64) (*Account, error) {
    // Guard.
    if s == nil {
        return nil, fmt.Errorf("failed to select account by id: missing required parameter: account_store=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account by id: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return nil, fmt.Errorf("failed to select account by id: missing required parameter: executor=null")
    }
    if id == 0 {
        return nil, fmt.Errorf("failed to select account by id: invalid parameter: id=empty")
    }

    // Generate a SELECT query.
    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColID)

    // Execute.
    row := executor.QueryRow(query, id)

    // Scan.
    result := &Account{}
    err := row.Scan(
        &result.ID,
        &result.Status,
        &result.Name,
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
func (s *AccountStore) Update(executor helper.Executor, id uint64, option *AccountUpdateOption) error {
    // Guard.
    if s == nil {
        return fmt.Errorf("failed to update account: missing required parameter: account_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update account: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to update account: missing required parameter: executor=null")
    }
    if option == nil {
        return fmt.Errorf("failed to update account: missing required parameter: option=null")
    }
    if id == 0 {
        return fmt.Errorf("failed to update account: invalid parameter: id=0")
    }

    assignments := make([]string, 0, 2)
    args := make([]any, 0, 3)

    if option.Status != nil {
        if err := ValidateAccountStatus(*option.Status); err != nil {
            return fmt.Errorf("failed to update account: %w", err)
        }
        assignments = append(assignments, ColStatus + " = ?")
        args = append(args, *option.Status)
    }
    if option.Name != nil {
        if err := ValidateAccountName(*option.Name); err != nil {
            return fmt.Errorf("failed to update account: %w", err)
        }
        assignments = append(assignments, ColName + " = ?")
        args = append(args, *option.Name)
    }

    if len(assignments) == 0 {
        return fmt.Errorf("failed to update account: invalid parameter: assignments=%q", "empty")
    }

    args = append(args, id)

    // Generate UPDATE query.
    query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?;", s.tableName, strings.Join(assignments, ", "), ColID)

    // Execute query.
    if _, err := executor.Exec(query, args...); err != nil {
        return fmt.Errorf("failed to update account: %w", err)
    }

    return nil
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

    conditions := make([]string, 0, 5)
    args := make([]any, 0, 7)

    if o.ID != nil {
        conditions = append(conditions, ColID + " = ?")
        args = append(args, *o.ID)
    }
    if o.Status != nil {
        conditions = append(conditions, ColStatus + " = ?")
        args = append(args, *o.Status)
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
            ColName,
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
