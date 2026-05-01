//
// password.go
//
package credential

import (
    "database/sql"
    "fmt"
    "regexp"
    "strings"
    "time"
    "unicode/utf8"
    _ "github.com/go-sql-driver/mysql"

    myHelper "github.com/k4k3ru-hub/db/go/mysql/helper"
)


const (
    DefaultPasswordCredentialTableName = "account_password_credentials"
)


var (
    emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)


type PasswordCredential struct {
    AccountID  uint64     `json:"accountId,string"`
    Status     Status     `json:"status"`
    Email      string     `json:"email"`
    Password   string     `json:"password"`
    LastUsedAt *time.Time `json:"lastUsedAt,omitempty"`
    CreatedAt  time.Time  `json:"createdAt,omitempty"`
    UpdatedAt  time.Time  `json:"updatedAt,omitempty"`
}


type PasswordCredentialStore struct {
    db               *sql.DB
    tableName        string
    accountTableName string
}


type PasswordCredentialSelectOption struct {
    AccountID   *uint64 `json:"accountId,string,omitempty"`
    Status      *Status `json:"status,omitempty"`
    Email       *string `json:"email,omitempty"`
    EmailLike   *string `json:"emailLike,omitempty"`
    OrderBy     string  `json:"orderBy"`
    OrderByDesc bool    `json:"orderByDesc"`
    Limit       int     `json:"limit"`
    Offset      int     `json:"offset"`
}


type PasswordCredentialUpdateOption struct {
    AccountID  uint64     `json:"accountId,string"`
    Status     *Status    `json:"status,omitempty"`
    Email      *string    `json:"email,omitempty"`
    Password   *string    `json:"password,omitempty"`
    LastUsedAt *time.Time `json:"lastUsedAt,omitempty"`
}


type PasswordCredentialDeleteOption struct {
    AccountID uint64 `json:"accountId,string"`
}


//
// Validate account ID.
//
func (c *PasswordCredential) ValidateAccountID() error {
    // Guard.
    if c == nil {
        return fmt.Errorf("invalid parameter: password_credential=null")
    }

    if c.AccountID == 0 {
        return fmt.Errorf("invalid parameter: account_id=0")
    }

    return nil
}


//
// Validate status.
//
func (c *PasswordCredential) ValidateStatus() error {
    // Guard.
    if c == nil {
        return fmt.Errorf("invalid parameter: password_credential=null")
    }

    if !c.Status.IsValid() {
        return fmt.Errorf("invalid parameter: status=%d", c.Status)
    }

    return nil
}


//
// Validate email.
//
func (c *PasswordCredential) ValidateEmail() error {
    // Guard.
    if c == nil {
        return fmt.Errorf("invalid parameter: password_credential=null")
    }

    if c.Email == "" {
        return fmt.Errorf("invalid parameter: email=empty")
    }
    if utf8.RuneCountInString(c.Email) > 254 {
        return fmt.Errorf("invalid parameter: email=%q", myHelper.TruncateRunes(c.Email, 254))
    }
    if !emailRegex.MatchString(c.Email) {
        return fmt.Errorf("invalid parameter: email=%q", c.Email)
    }

    return nil
}


//
// Validate password.
//
func (c *PasswordCredential) ValidatePassword() error {
    // Guard.
    if c == nil {
        return fmt.Errorf("invalid parameter: password_credential=null")
    }

    if c.Password == "" {
        return fmt.Errorf("invalid parameter: password=empty")
    }
    if utf8.RuneCountInString(c.Password) > 254 {
        return fmt.Errorf("invalid parameter: password=invalid")
    }

    return nil
}


//
// Create new account password credential store.
//
// Version:
//   - 2026-04-30: Added.
//
func NewPasswordCredentialStore(db *sql.DB, tableName, accountTableName string) (*PasswordCredentialStore, error) {
    // Guard.
    if db == nil {
        return nil, fmt.Errorf("failed to create account password credential store: missing required parameter: db=null")
    }
    if tableName == "" {
        return nil, fmt.Errorf("failed to create account password credential store: missing required parameter: table_name=empty")
    }
    if accountTableName == "" {
        return nil, fmt.Errorf("failed to create account password credential store: missing required parameter: account_table_name=empty")
    }

    return &PasswordCredentialStore{
        db:               db,
        tableName:        tableName,
        accountTableName: accountTableName,
    }, nil
}


//
// Count account password credentials.
// 
// Version:
//   - 2026-04-30: Added.
//   
func (s *PasswordCredentialStore) Count(option *PasswordCredentialSelectOption) (int64, error) {
    // Guard.
    if s == nil {
        return 0, fmt.Errorf("failed to count account password credentials: missing required parameter: password_credential_store=null")
    }   
    if s.db == nil {
        return 0, fmt.Errorf("failed to count account password credentials: missing required parameter: db=null")
    }   
    if s.tableName == "" {
        return 0, fmt.Errorf("failed to count account password credentials: missing required parameter: table_name=empty")
    }   
    if option == nil {
        return 0, fmt.Errorf("failed to count account password credentials: missing required parameter: select_option=null")
    }
    
    // Build conditions. 
    conditions, args := buildPasswordCredentialSelectCondition(option)
            
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
        return 0, fmt.Errorf("failed to count account password credentials: %w", err)
    }

    return result, nil
}


//
// Create account password credentials table.
//
// Version:
//   - 2026-04-30: Added.
//
func (s *PasswordCredentialStore) CreateTable() error {
    // Guard.
    if s == nil {
        return fmt.Errorf("failed to create account password credentials table: missing required parameter: password_credential_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to create account password credentials table: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to create account password credentials table: missing required parameter: table_name=empty")
    }
    if s.accountTableName == "" {
        return fmt.Errorf("failed to create account password credentials table: missing required parameter: account_table_name=empty")
    }

    // Generate a CREATE TABLE query.
    query := fmt.Sprintf(
        `CREATE TABLE IF NOT EXISTS %s (
            %s BIGINT UNSIGNED NOT NULL COMMENT 'Account ID',
            %s TINYINT UNSIGNED NOT NULL COMMENT 'Status',
            %s VARCHAR(254) NOT NULL COMMENT 'Email',
            %s VARCHAR(254) NOT NULL COMMENT 'Password hash',
            %s DATETIME COMMENT 'Last used at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
            PRIMARY KEY (%s),
            UNIQUE KEY uk_account_password_credentials_email (%s),
            KEY idx_account_password_credentials_status (%s),
            CONSTRAINT fk_account_password_credentials_account_id FOREIGN KEY (%s) REFERENCES %s (%s) ON DELETE CASCADE ON UPDATE CASCADE);
        `,
        s.tableName,
        ColAccountID,
        ColStatus,
        ColEmail,
        ColPassword,
        ColLastUsedAt,
        ColCreatedAt,
        ColUpdatedAt,
        ColAccountID,
        ColEmail,
        ColStatus,
        ColAccountID, s.accountTableName, ColID,
    )

    // Execute the query.
    if _, err := s.db.Exec(query); err != nil {
        return fmt.Errorf("failed to create account password credentials table: %w", err)
    }

    return nil
}


//
// Delete account password credential by account ID.
//
// Version:
//   - 2026-04-30: Added.
//
func (s *PasswordCredentialStore) DeleteByAccountID(accountID uint64) error {
    // Guard.
    if s == nil {
        return fmt.Errorf("failed to delete account password credential by account id: missing required parameter: password_credential_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to delete account password credential by account id: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to delete account password credential by account id: missing required parameter: table_name=empty")
    }
    if accountID == 0 {
        return fmt.Errorf("failed to delete account password credential by account id: invalid parameter: account_id=0")
    }
    
    // Generate a DELETE query.
    query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", s.tableName, ColAccountID)
    
    // Execute.
    if _, err := s.db.Exec(query, accountID); err != nil {
        return fmt.Errorf("failed to delete account password credential by account id: %w", err)
    }   

    return nil
}


//
// Insert an account password credential.
//
// Version:
//   - 2026-04-30: Added.
//
func (s *PasswordCredentialStore) Insert(row *PasswordCredential) error {
    // Guard.
    if s == nil {
        return fmt.Errorf("failed to insert account password credential: missing required parameter: password_credential_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to insert account password credential: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to insert account password credential: missing required parameter: table_name=empty")
    }
    if row == nil {
        return fmt.Errorf("failed to insert account password credential: missing required parameter: account_password_credential=null table=%q", s.tableName)
    }
    if err := row.ValidateAccountID(); err != nil {
        return fmt.Errorf("failed to insert account password credential: %w table=%q", err, s.tableName)
    }
    if err := row.ValidateStatus(); err != nil {
        return fmt.Errorf("failed to insert account password credential: %w table=%q", err, s.tableName)
    }
    if err := row.ValidateEmail(); err != nil {
        return fmt.Errorf("failed to insert account password credential: %w table=%q", err, s.tableName)
    }
    if err := row.ValidatePassword(); err != nil {
        return fmt.Errorf("failed to insert account password credential: %w table=%q", err, s.tableName)
    }

    // Generate an INSERT query.
    query := fmt.Sprintf(
        "INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?, ?, ?);",
        s.tableName,
        ColAccountID,
        ColStatus,
        ColEmail,
        ColPassword,
        ColLastUsedAt,
        ColCreatedAt,
        ColUpdatedAt,
    )

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
        row.AccountID,
        row.Status,
        row.Email,
        row.Password,
        row.LastUsedAt,
        row.CreatedAt,
        row.UpdatedAt,
    ); err != nil {
        return fmt.Errorf("failed to insert account password credential: %w", err)
    }

    return nil
}


//
// Select account password credentials.
//
// Version:
//   - 2026-04-30: Added.
//
func (s *PasswordCredentialStore) Select(option *PasswordCredentialSelectOption) ([]*PasswordCredential, error) {
    // Guard.
    if s == nil {
        return nil, fmt.Errorf("failed to select account password credentials: missing required parameter: password_credential_store=null")
    }
    if s.db == nil {
        return nil, fmt.Errorf("failed to select account password credentials: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account password credentials: missing required parameter: table_name=empty")
    }
    if option == nil {
        return nil, fmt.Errorf("failed to select account password credentials: missing required parameter: select_option=null")
    }

    // Build conditions.
    conditions, args := buildPasswordCredentialSelectCondition(option)

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
            return nil, fmt.Errorf("failed to select account password credentials: invalid parameter: offset=%d", option.Offset)
        }
        query.WriteString(" LIMIT ? OFFSET ?")
        args = append(args, option.Limit, option.Offset)
    }

    // Execute.
    rows, err := s.db.Query(query.String(), args...)
    if err != nil {
        return nil, fmt.Errorf("failed to select account password credentials: %w", err)
    }

    defer rows.Close()

    // Scan.
    var result []*PasswordCredential
    for rows.Next() {
        row := &PasswordCredential{}
        err := rows.Scan(
            &row.AccountID,
            &row.Status,
            &row.Email,
            &row.Password,
            &row.LastUsedAt,
            &row.CreatedAt,
            &row.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to select account password credentials: %w", err)
        }
        result = append(result, row)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to select account password credentials: %w", err)
    }

    return result, nil
}


//
// Select account password credential by account ID.
//
// Version:
//   - 2026-04-30: Added.
//
func (s *PasswordCredentialStore) SelectByAccountID(accountID uint64) (*PasswordCredential, error) {
    // Guard.
    if s == nil {
        return nil, fmt.Errorf("failed to select account password credential by account id: missing required parameter: password_credential_store=null")
    }
    if s.db == nil {
        return nil, fmt.Errorf("failed to select account password credential by account id: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account password credential by account id: missing required parameter: table_name=empty")
    }
    if accountID == 0 {
        return nil, fmt.Errorf("failed to select account password credential by account id: invalid parameter: account_id=0")
    }

    // Generate a SELECT query.
    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColAccountID)

    // Execute.
    row := s.db.QueryRow(query, accountID)

    // Scan.
    result := &PasswordCredential{}
    err := row.Scan(
        &result.AccountID,
        &result.Status,
        &result.Email,
        &result.Password,
        &result.LastUsedAt,
        &result.CreatedAt,
        &result.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to select account password credential by account id: %w", err)
    }

    return result, nil
}


//
// Select account password credential by email.
//
// Version:
//   - 2026-04-30: Added.
//
func (s *PasswordCredentialStore) SelectByEmail(email string) (*PasswordCredential, error) {
    // Guard.
    if s == nil {
        return nil, fmt.Errorf("failed to select account password credential by email: missing required parameter: password_credential_store=null")
    }
    if s.db == nil {
        return nil, fmt.Errorf("failed to select account password credential by email: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account password credential by email: missing required parameter: table_name=empty")
    }

    p := &PasswordCredential{
        Email: email,
    }
    if err := p.ValidateEmail(); err != nil {
        return nil, fmt.Errorf("failed to select account password credential by email: invalid parameter: email=%q", email)
    }

    // Generate a SELECT query.
    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColEmail)

    // Execute.
    row := s.db.QueryRow(query, email)

    // Scan.
    result := &PasswordCredential{}
    err := row.Scan(
        &result.AccountID,
        &result.Status,
        &result.Email,
        &result.Password,
        &result.LastUsedAt,
        &result.CreatedAt,
        &result.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to select account password credential by email: %w", err)
    }

    return result, nil
}


//
// Update account password credential.
//
// Version:
//   - 2026-04-30: Added.
//
func (s *PasswordCredentialStore) Update(option *PasswordCredentialUpdateOption) error {
    // Guard.
    if s == nil {
        return fmt.Errorf("failed to update account password credential: missing required parameter: password_credential_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to update account password credential: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update account password credential: missing required parameter: table_name=empty")
    }
    if option == nil {
        return fmt.Errorf("failed to update account password credential: missing required parameter: option=null")
    }
    if option.AccountID == 0 {
        return fmt.Errorf("failed to update account password credential: invalid parameter: account_id=0")
    }

    assignments := make([]string, 0, 4)
    args := make([]interface{}, 0, 4)

    if option.Status != nil {
        p := PasswordCredential{
            Status: *option.Status,
        }
        if err := p.ValidateStatus(); err != nil {
            return fmt.Errorf("failed to update account password credential: %w", err)
        }
        assignments = append(assignments, ColStatus + " = ?")
        args = append(args, *option.Status)
    }

    if option.Email != nil {
        p := PasswordCredential{
            Email: *option.Email,
        }
        if err := p.ValidateEmail(); err != nil {
            return fmt.Errorf("failed to update account password credential: %w", err)
        }
        assignments = append(assignments, ColEmail + " = ?")
        args = append(args, *option.Email)
    }

    if option.Password != nil {
        p := PasswordCredential{
            Password: *option.Password,
        }
        if err := p.ValidatePassword(); err != nil {
            return fmt.Errorf("failed to update account password credential: %w", err)
        }
        assignments = append(assignments, ColPassword + " = ?")
        args = append(args, *option.Password)
    }

    if option.LastUsedAt != nil {
        assignments = append(assignments, ColLastUsedAt + " = ?")
        args = append(args, *option.LastUsedAt)
    }

    if len(assignments) == 0 {
        return fmt.Errorf("failed to update account password credential: invalid parameter: assignments=empty")
    }

    args = append(args, option.AccountID)

    // Generate a UPDATE query.
    query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?;", s.tableName, strings.Join(assignments, ", "), ColAccountID)

    // Execute.
    if _, err := s.db.Exec(query, args...); err != nil {
        return fmt.Errorf("failed to update account password credential: %w", err)
    }

    return nil
}


//
// Update last used at.
//
// Version:
//   - 2026-04-30: Added.
//
func (s *PasswordCredentialStore) UpdateLastUsedAt(accountID uint64) error {
    // Guard.
    if s == nil {
        return fmt.Errorf("failed to update account password credential last used at: missing required parameter: password_credential_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to update account password credential last used at: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update account password credential last used at: missing required parameter: table_name=empty")
    }
    if accountID == 0 {
        return fmt.Errorf("failed to update account password credential last used at: invalid parameter: account_id=0")
    }

    now := time.Now()

    return s.Update(&PasswordCredentialUpdateOption{
        AccountID:  accountID,
        LastUsedAt: &now,
    })
}


//
// Build account password credential select conditions.
//
// Version:
//   - 2026-04-30: Added.
//
func buildPasswordCredentialSelectCondition(option *PasswordCredentialSelectOption) ([]string, []any) {
    conditions := make([]string, 0, 4)
    args := make([]any, 0, 4)

    if option != nil {
        if option.AccountID != nil {
            conditions = append(conditions, ColAccountID + " = ?")
            args = append(args, *option.AccountID)
        }
        if option.Status != nil {
            conditions = append(conditions, ColStatus + " = ?")
            args = append(args, *option.Status)
        }
        if option.Email != nil {
            conditions = append(conditions, ColEmail + " = ?")
            args = append(args, *option.Email)
        }
        if option.EmailLike != nil {
            conditions = append(conditions, ColEmail + " LIKE ?")
            args = append(args, "%" + *option.EmailLike + "%")
        }
    }

    return conditions, args
}

