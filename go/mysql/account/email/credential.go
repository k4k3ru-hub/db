//
// credential.go
//
package email

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "strings"
    "time"
    "unicode/utf8"

    _ "github.com/go-sql-driver/mysql"

    "github.com/k4k3ru-hub/db/go/mysql/account"
    "github.com/k4k3ru-hub/db/go/mysql/helper"
)

const (
    DefaultCredentialTableName = "account_email_credentials"
)

var (
    credentialIDCounter = &helper.IdCounter{}
)


//
// Credential.
//
// Version:
//   - 2026-05-03: Added.
//
type Credential struct {
    ID         uint64           `json:"id,string"`
    AccountID  uint64           `json:"accountId,string"`
    Email      string           `json:"email"`
    Status     CredentialStatus `json:"status"`
    LastUsedAt *time.Time       `json:"lastUsedAt,omitempty"`
    MetaData   *string          `json:"metaData,omitempty"`
    CreatedAt  time.Time        `json:"createdAt"`
    UpdatedAt  time.Time        `json:"updatedAt"`
}


//
// CredentialStore.
//
// Version:
//   - 2026-05-03: Added.
//
type CredentialStore struct {
    executor         helper.Executor
    tableName        string
    accountTableName string
}


//
// CredentialSelectOption.
//
// Version:
//   - 2026-05-03: Added.
//
type CredentialSelectOption struct {
    ID            *uint64           `json:"id,string,omitempty"`
    AccountID     *uint64           `json:"accountId,string,omitempty"`
    Email         *string           `json:"email,omitempty"`
    EmailLike     *string           `json:"emailLike,omitempty"`
    Status        *CredentialStatus `json:"status,omitempty"`
    LastUsedAtGTE *time.Time        `json:"lastUsedAtGte,omitempty"`
    LastUsedAtLTE *time.Time        `json:"lastUsedAtLte,omitempty"`
    OrderBy       string            `json:"orderBy"`
    OrderByDesc   bool              `json:"orderByDesc"`
    Limit         int               `json:"limit"`
    Offset        int               `json:"offset"`
}


//
// CredentialUpdateOption.
//
// Version:
//   - 2026-05-03: Added.
//
type CredentialUpdateOption struct {
    ID         uint64            `json:"id,string"`
    AccountID  *uint64           `json:"accountId,string,omitempty"`
    Email      *string           `json:"email,omitempty"`
    Status     *CredentialStatus `json:"status,omitempty"`
    LastUsedAt *time.Time        `json:"lastUsedAt,omitempty"`
    MetaData   *string           `json:"metaData,omitempty"`
}


//
// Generate credential ID.
// 
// Version:
//   - 2026-05-03: Added.
//
func GenerateCredentialID() uint64 {
    return credentialIDCounter.GenerateID()
}


//
// Create new credential store.
//
// Version:
//   - 2026-05-03: Added.
//
func NewCredentialStore(executor helper.Executor, tableName, accountTableName string) (*CredentialStore, error) {
    if executor == nil {
        return nil, fmt.Errorf("failed to create account email credential store: missing required parameter: executor=null")
    }
    if tableName == "" {
        return nil, fmt.Errorf("failed to create account email credential store: missing required parameter: table_name=%q", "empty")
    }
    if accountTableName == "" {
        return nil, fmt.Errorf("failed to create account email credential store: missing required parameter: account_table_name=empty")
    }

    return &CredentialStore{
        executor:               executor,
        tableName:        tableName,
        accountTableName: accountTableName,
    }, nil
}


//
// Validate account email credential ID.
//
// Version:
//   - 2026-05-09: Added.
//
func ValidateCredentialID(id uint64) error {
    if id == 0 {
        return fmt.Errorf("invalid parameter: id=0")
    }
    return nil
}


//
// Validate account email credential ID.
//
// Version:
//   - 2026-05-09: Added.
//
func (c *Credential) ValidateID() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: account_email_credential=null")
    }
    return ValidateCredentialID(c.ID)
}


//
// Validate account email credential account ID.
//
// Version:
//   - 2026-05-09: Added.
//
func ValidateCredentialAccountID(accountID uint64) error {
    if accountID == 0 {
        return fmt.Errorf("invalid parameter: account_id=0")
    }
    return nil
}


//
// Validate account email credential account ID.
//
// Version:
//   - 2026-05-09: Added.
//
func (c *Credential) ValidateAccountID() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: account_email_credential=null")
    }
    return ValidateCredentialID(c.AccountID)
}


//
// Validate account email credential email.
//
// Version:
//   - 2026-05-03: Added.
//
func ValidateCredentialEmail(email string) error {
    if email == "" {
        return fmt.Errorf("email=%q", "empty")
    }
    if utf8.RuneCountInString(email) > 255 {
        return fmt.Errorf("invalid parameter: max_length=255 email=%q", "too long")
    }
    return nil
}


//
// Validate account email credential email.
//
// Version:
//   - 2026-05-03: Added.
//
func (c *Credential) ValidateEmail(email string) error {
    if c == nil {
        return fmt.Errorf("missing required parameter: account_email_credential=null")
    }
    return ValidateCredentialEmail(c.Email)
}


//
// Validate account email credential status.
//
// Version:
//   - 2026-05-09: Added.
//
func ValidateCredentialStatus(status CredentialStatus) error {
    if !status.IsValid() {
        return fmt.Errorf("invalid parameter: status=%d", status)
    }
    return nil
}


//
// Validate account email credential status.
//
// Version:
//   - 2026-05-09: Added.
//
func (c *Credential) ValidateStatus() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: account_email_credential=null")
    }
    return ValidateCredentialStatus(c.Status)
}


//
// Validate account email credential last used at.
//
// Version:
//   - 2026-05-09: Added.
//
func ValidateCredentialLastUsedAt(lastUsedAt *time.Time) error {
    if lastUsedAt == nil {
        return nil
    }
    if lastUsedAt.IsZero() {
        return fmt.Errorf("invalid parameter: last_used_at=%q", "empty")
    }
    if lastUsedAt.After(time.Now().UTC()) {
        return fmt.Errorf("invalid parameter: last_used_at=%q", lastUsedAt.Format(time.RFC3339))
    }
    return nil
}


//
// Validate account email credential last used at.
//
// Version:
//   - 2026-05-09: Added.
//
func (c *Credential) ValidateLastUsedAt() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: account_email_credential=null")
    }
    return ValidateCredentialLastUsedAt(c.LastUsedAt)
}


//
// Validate meta data.
//
// Version:
//   - 2026-05-02: Added.
//
func ValidateCredentialMetaData(meta *string) error {
    if meta == nil {
        return nil
    }
    if !json.Valid([]byte(*meta)) {
        return fmt.Errorf("invalid parameter: meta_data=%q", helper.TruncateRunes(*meta, 1024))
    }
    return nil
}


//
// Validate meta data.
//
// Version:
//   - 2026-05-02: Added.
//
func (c *Credential) ValidateMetaData() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: account_email_credential=null")
    }
    return ValidateCredentialMetaData(c.MetaData)
}


//
// Create account email credentials table.
//
// Version:
//   - 2026-05-02: Added.
//
func (s *CredentialStore) CreateTable() error {
    if s == nil {
        return fmt.Errorf("failed to create account email credentials table: missing required parameter: credential_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to create account email credentials table: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to create account email credentials table: missing required parameter: table_name=%q", "empty")
    }
    if s.accountTableName == "" {
        return fmt.Errorf("failed to create account email credential table: missing required parameter: account_table_name=%q", "empty")
    }

    query := fmt.Sprintf(
        `CREATE TABLE IF NOT EXISTS %s (
            %s BIGINT UNSIGNED NOT NULL COMMENT 'ID',
            %s BIGINT UNSIGNED NOT NULL COMMENT 'Account ID',
            %s VARCHAR(255) NOT NULL COMMENT 'Email',
            %s TINYINT UNSIGNED NOT NULL COMMENT 'Status',
            %s DATETIME NULL COMMENT 'Last used at',
            %s JSON NULL COMMENT 'Meta data',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
            PRIMARY KEY (%s),
            UNIQUE KEY uk_account_email_credential_account_id (%s),
            UNIQUE KEY uk_account_email_credential_email (%s),
            KEY idx_account_app_email_credential_status (%s),
            CONSTRAINT fk_account_email_credential_account_id FOREIGN KEY (%s) REFERENCES %s (%s) ON DELETE CASCADE ON UPDATE CASCADE
        );`,
        s.tableName,
        ColID,
        ColAccountID,
        ColEmail,
        ColStatus,
        ColLastUsedAt,
        ColMetaData,
        ColCreatedAt,
        ColUpdatedAt,
        ColID,
        ColAccountID,
        ColEmail,
        ColStatus,
        ColAccountID, s.accountTableName, account.ColID,
    )

    if _, err := s.executor.Exec(query); err != nil {
        return fmt.Errorf("failed to create account email credentials table: %w", err)
    }

    return nil
}


//
// Insert account email credential.
//
// Version:
//   - 2026-05-04: Added.
//
func (s *CredentialStore) Insert(row *Credential) error {
    if s == nil {
        return fmt.Errorf("failed to insert account email credential: missing required parameter: credential_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to insert account email credential: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to insert account email credential: missing required parameter: table_name=%q", "empty")
    }
    if row == nil {
        return fmt.Errorf("failed to insert account email credential: missing required parameter: credential=null")
    }
    if err := ValidateCredentialAccountID(row.AccountID); err != nil {
        return fmt.Errorf("failed to insert account email credential: %w", err)
    }
    if err := ValidateCredentialEmail(row.Email); err != nil {
        return fmt.Errorf("failed to insert account email credential: %w", err)
    }
    if err := row.ValidateStatus(); err != nil {
        return fmt.Errorf("failed to insert account email credential: %w", err)
    }
    if err := row.ValidateLastUsedAt(); err != nil {
        return fmt.Errorf("failed to insert account email credential: %w", err)
    }
    if err := row.ValidateMetaData(); err != nil {
        return fmt.Errorf("failed to insert account email credential: %w", err)
    }

    if row.ID == 0 {
        row.ID = GenerateCredentialID()
    }

    now := time.Now()
    if row.CreatedAt.IsZero() {
        row.CreatedAt = now
    }
    if row.UpdatedAt.IsZero() {
        row.UpdatedAt = now
    }

    query := fmt.Sprintf(
        "INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?, ?, ?, ?);",
        s.tableName,
        ColID,
        ColAccountID,
        ColEmail,
        ColStatus,
        ColLastUsedAt,
        ColMetaData,
        ColCreatedAt,
        ColUpdatedAt,
    )

    if _, err := s.executor.Exec(
        query,
        row.ID,
        row.AccountID,
        row.Email,
        row.Status,
        row.LastUsedAt,
        row.MetaData,
        row.CreatedAt,
        row.UpdatedAt,
    ); err != nil {
        return fmt.Errorf("failed to insert account email credential: %w", err)
    }

    return nil
}


//
// Select account email credential by ID.
//
// Version:
//   - 2026-05-04: Added.
//
func (s *CredentialStore) SelectByID(id uint64) (*Credential, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select account email credential by id: missing required parameter: credential_store=null")
    }
    if s.executor == nil {
        return nil, fmt.Errorf("failed to select account email credential by id: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account email credential by id: missing required parameter: table_name=%q", "empty")
    }
    if id == 0 {
        return nil, fmt.Errorf("failed to select account email credential by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColID)

    result := &Credential{}
    err := s.executor.QueryRow(query, id).Scan(
        &result.ID,
        &result.AccountID,
        &result.Email,
        &result.Status,
        &result.LastUsedAt,
        &result.MetaData,
        &result.CreatedAt,
        &result.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to select account email credential by id: %w", err)
    }

    return result, nil
}


//
// Select account email credential by email.
//
// Version:
//   - 2026-05-05: Added.
//
func (s *CredentialStore) SelectByEmail(email string) (*Credential, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select account email credential by email: missing required parameter: credential_store=null")
    }
    if s.executor == nil {
        return nil, fmt.Errorf("failed to select account email credential by email: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account email credential by email: missing required parameter: table_name=%q", "empty")
    }
    if err := ValidateCredentialEmail(email); err != nil {
        return nil, fmt.Errorf("failed to select account email credential by email: %w", err)
    }

    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColEmail)

    result := &Credential{}
    err := s.executor.QueryRow(query, email).Scan(
        &result.ID,
        &result.AccountID,
        &result.Email,
        &result.Status,
        &result.LastUsedAt,
        &result.MetaData,
        &result.CreatedAt,
        &result.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to select account email credential by email: %w", err)
    }

    return result, nil
}


//
// Select account email credentials.
//
// Version:
//   - 2026-05-04: Added.
//
func (s *CredentialStore) Select(option *CredentialSelectOption) ([]*Credential, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select account email credentials: missing required parameter: credential_store=null")
    }
    if s.executor == nil {
        return nil, fmt.Errorf("failed to select account email credentials: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account email credentials: missing required parameter: table_name=%q", "empty")
    }
    if err := option.Validate(); err != nil {
        return nil, fmt.Errorf("failed to select account email credentials: %w", err)
    }

    query, args := option.BuildQuery("SELECT * FROM " + s.tableName)

    rows, err := s.executor.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to select account email credentials: %w", err)
    }
    defer rows.Close()

    var result []*Credential
    for rows.Next() {
        row := &Credential{}
        if err := rows.Scan(
            &row.ID,
            &row.AccountID,
            &row.Email,
            &row.Status,
            &row.LastUsedAt,
            &row.MetaData,
            &row.CreatedAt,
            &row.UpdatedAt,
        ); err != nil {
            return nil, fmt.Errorf("failed to select account email credentials: %w", err)
        }

        result = append(result, row)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to select account email credentials: %w", err)
    }

    return result, nil
}


//
// Count account email credentials.
//
// Version:
//   - 2026-05-04: Added.
//
func (s *CredentialStore) Count(option *CredentialSelectOption) (int64, error) {
    if s == nil {
        return 0, fmt.Errorf("failed to count account email credentials: missing required parameter: credential_store=null")
    }
    if s.executor == nil {
        return 0, fmt.Errorf("failed to count account email credentials: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return 0, fmt.Errorf("failed to count account email credentials: missing required parameter: table_name=%q", "empty")
    }
    if err := option.Validate(); err != nil {
        return 0, fmt.Errorf("failed to count account email credentials: %w", err)
    }

    query, args := option.BuildQuery("SELECT COUNT(*) FROM " + s.tableName)

    var result int64
    if err := s.executor.QueryRow(query, args...).Scan(&result); err != nil {
        return 0, fmt.Errorf("failed to count account email credentials: %w", err)
    }

    return result, nil
}


//
// Update account email credential.
//
// Version:
//   - 2026-05-04: Added.
//
func (s *CredentialStore) Update(option *CredentialUpdateOption) error {
    if s == nil {
        return fmt.Errorf("failed to update account email credential: missing required parameter: credential_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to update account email credential: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update account email credential: missing required parameter: table_name=%q", "empty")
    }
    if err := option.Validate(); err != nil {
        return fmt.Errorf("failed to update account email credential: %w", err)
    }

    assignments := make([]string, 0, 5)
    args := make([]any, 0, 6)

    if option.AccountID != nil {
        assignments = append(assignments, ColAccountID + " = ?")
        args = append(args, *option.AccountID)
    }

    if option.Email != nil {
        assignments = append(assignments, ColEmail + " = ?")
        args = append(args, *option.Email)
    }

    if option.Status != nil {
        assignments = append(assignments, ColStatus + " = ?")
        args = append(args, *option.Status)
    }

    if option.LastUsedAt != nil {
        assignments = append(assignments, ColLastUsedAt + " = ?")
        args = append(args, *option.LastUsedAt)
    }

    if option.MetaData != nil {
        assignments = append(assignments, ColMetaData + " = ?")
        args = append(args, *option.MetaData)
    }

    if len(assignments) == 0 {
        return fmt.Errorf("failed to update account email credential: invalid parameter: assignments=%q", "empty")
    }

    args = append(args, option.ID)

    query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?;", s.tableName, strings.Join(assignments, ", "), ColID)

    if _, err := s.executor.Exec(query, args...); err != nil {
        return fmt.Errorf("failed to update account email credential: %w", err)
    }

    return nil
}


//
// Delete account email credential by ID.
//
// Version:
//   - 2026-05-04: Added.
//
func (s *CredentialStore) DeleteByID(id uint64) error {
    if s == nil {
        return fmt.Errorf("failed to delete account email credential by id: missing required parameter: credential_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to delete account email credential by id: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to delete account email credential by id: missing required parameter: table_name=%q", "empty")
    }
    if id == 0 {
        return fmt.Errorf("failed to delete account email credential by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", s.tableName, ColID)

    if _, err := s.executor.Exec(query, id); err != nil {
        return fmt.Errorf("failed to delete account email credential by id: %w", err)
    }

    return nil
}


//      
// Build query.
//  
// Version:
//   - 2025-05-03: Added.
//  
func (o *CredentialSelectOption) BuildQuery(selectFromClause string) (string, []any) {
    // Guard.
    if o == nil {
        return selectFromClause, nil
    }

    var query strings.Builder
    query.WriteString(selectFromClause)

    conditions := make([]string, 0, 7)
    args := make([]any, 0, 9)
        
    if o.ID != nil {
        conditions = append(conditions, ColID + " = ?")
        args = append(args, *o.ID)
    }
    if o.AccountID != nil {
        conditions = append(conditions, ColAccountID + " = ?")
        args = append(args, *o.AccountID)
    }
    if o.Email != nil {
        conditions = append(conditions, ColEmail + " = ?")
        args = append(args, *o.Email)
    }
    if o.EmailLike != nil {
        conditions = append(conditions, ColEmail + " LIKE ?")
        args = append(args, "%" + *o.EmailLike + "%")
    }
    if o.Status != nil {
        conditions = append(conditions, ColStatus + " = ?")
        args = append(args, *o.Status)
    }
    if o.LastUsedAtGTE != nil {
        conditions = append(conditions, ColLastUsedAt + " >= ?")
        args = append(args, *o.LastUsedAtGTE)
    }
    if o.LastUsedAtLTE != nil {
        conditions = append(conditions, ColLastUsedAt + " <= ?")
        args = append(args, *o.LastUsedAtLTE)
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
//   - 2025-05-02: Added.
//
func (o *CredentialSelectOption) Validate() error {
    // Guard.
    if o == nil {
        return fmt.Errorf("missing required parameter: credential_select_option=null")
    }

    if o.ID != nil {
        if err := ValidateCredentialID(*o.ID); err != nil {
            return err
        }
    }

    if o.AccountID != nil {
        if err := ValidateCredentialAccountID(*o.AccountID); err != nil {
            return err
        }
    }

    if o.Email != nil {
        if err := ValidateCredentialEmail(*o.Email); err != nil {
            return err
        }
    }

    if o.EmailLike != nil {
        if err := ValidateCredentialEmail(*o.EmailLike); err != nil {
            return err
        }
    }

    if o.Status != nil {
        if err := (&Credential{Status: *o.Status}).ValidateStatus(); err != nil {
            return err
        }
    }

    if o.LastUsedAtGTE != nil && o.LastUsedAtLTE != nil && o.LastUsedAtGTE.After(*o.LastUsedAtLTE) {
        return fmt.Errorf("invalid parameter: last_used_at_gte=%s last_used_at_lte=%s", *o.LastUsedAtGTE, *o.LastUsedAtLTE)
    }

    if o.OrderBy != "" {
        switch o.OrderBy {
        case ColID,
            ColAccountID,
            ColEmail,
            ColStatus,
            ColLastUsedAt,
            ColCreatedAt,
            ColUpdatedAt:
        default:
            return fmt.Errorf("invalid parameter: order_by=%s", o.OrderBy)
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


//
// Validate account email credential update option.
//
// Version:
//   - 2025-05-03: Added.
//
func (o *CredentialUpdateOption) Validate() error {
    // Guard.
    if o == nil {
        return fmt.Errorf("missing required parameter: credential_update_option=null")
    }

    if o.ID == 0 {
        return fmt.Errorf("invalid parameter: id=0")
    }

    if o.AccountID != nil {
        if err := ValidateCredentialAccountID(*o.AccountID); err != nil {
            return err
        }
    }

    if o.Email != nil {
        if err := ValidateCredentialEmail(*o.Email); err != nil {
            return err
        }
    }

    if o.Status != nil {
        c := Credential{
            Status: *o.Status,
        }
        if err := c.ValidateStatus(); err != nil {
            return err
        }
    }

    if o.LastUsedAt != nil {
        c := Credential{
            LastUsedAt: o.LastUsedAt,
        }
        if err := c.ValidateLastUsedAt(); err != nil {
            return err
        }
    }

    if o.MetaData != nil {
        c := Credential{
            MetaData: o.MetaData,
        }
        if err := c.ValidateMetaData(); err != nil {
            return err
        }
    }

    return nil
}
