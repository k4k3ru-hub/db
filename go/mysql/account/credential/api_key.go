//
// api_key.go
//
package credential

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "strings"
    "time"
    "unicode/utf8"
    _ "github.com/go-sql-driver/mysql"

    myHelper "github.com/k4k3ru-hub/db/go/mysql/helper"
)


const (
    DefaultAPIKeyCredentialTableName = "account_api_key_credentials"
)


var (
    apiKeyCredentialIDCounter = &myHelper.IdCounter{}
)


type APIKeyCredential struct {
    ID          uint64     `json:"id,string"`
    AccountID   uint64     `json:"accountId,string"`
    Status      Status     `json:"status"`
    Name        string     `json:"name"`
    PublicToken string     `json:"publicToken"`
    SecretToken string     `json:"secretToken"`
    ExpiresAt   time.Time  `json:"expiresAt"`
    Scopes      []string   `json:"scopes"`
    LastUsedAt  *time.Time `json:"lastUsedAt,omitempty"`
    CreatedAt   time.Time  `json:"createdAt,omitempty"`
    UpdatedAt   time.Time  `json:"updatedAt,omitempty"`
}


type APIKeyCredentialStore struct {
    db               *sql.DB
    tableName        string
    accountTableName string
}


type APIKeyCredentialSelectOption struct {
    AccountID   *uint64 `json:"accountId,string,omitempty"`
    Status      *Status `json:"status,omitempty"`
    NameLike    *string `json:"nameLike,omitempty"`
    OrderBy     string  `json:"orderBy"`
    OrderByDesc bool    `json:"orderByDesc"`
    Limit       int     `json:"limit"`
    Offset      int     `json:"offset"`
}


type APIKeyCredentialUpdateOption struct {  
    ID          uint64     `json:"id,string"`
    AccountID   *uint64    `json:"accountId,string,omitempty"`
    Status      *Status    `json:"status,omitempty"`
    Name        *string    `json:"name,omitempty"`
    PublicToken *string    `json:"publicToken,omitempty"`
    SecretToken *string    `json:"secretToken,omitempty"`
    ExpiresAt   *time.Time `json:"expiresAt,omitempty"`
    Scopes      []string   `json:"scopes,omitempty"`
    LastUsedAt  *time.Time `json:"lastUsedAt,omitempty"`
}


type APIKeyCredentialDeleteOption struct {
    ID uint64 `json:"id,string"`
}


//
// Generate account API key credential ID.
//  
func GenerateAPIKeyCredentialID() uint64 {
    return apiKeyCredentialIDCounter.GenerateID()
}


//
// Validate account ID.
//
func (c *APIKeyCredential) ValidateAccountID() error {
    // Guard.
    if c == nil {
        return fmt.Errorf("invalid parameter: api_key_credential=null")
    }

    if c.AccountID == 0 {
        return fmt.Errorf("invalid parameter: account_id=0")
    }

    return nil
}


//
// Validate expires at.
//
func (c *APIKeyCredential) ValidateExpiresAt() error {
    // Guard.
    if c == nil {
        return fmt.Errorf("invalid parameter: api_key_credential=null")
    }

    if c.ExpiresAt.IsZero() {
        return fmt.Errorf("invalid parameter: expires_at=empty")
    }

    now := time.Now()

    if !c.ExpiresAt.After(now) {
        return fmt.Errorf("invalid parameter: expires_at=%s", c.ExpiresAt.Format(time.RFC3339))
    }

    return nil
}


//
// Validate used at.
//
func (c *APIKeyCredential) ValidateLastUsedAt() error {
    // Guard.
    if c == nil {
        return fmt.Errorf("invalid parameter: api_key_credential=null")
    }

    if c.LastUsedAt == nil {
        return nil
    }

    if c.LastUsedAt.IsZero() {
        return fmt.Errorf("invalid parameter: used_at=empty")
    }

    now := time.Now()

    if c.LastUsedAt.After(now) {
        return fmt.Errorf("invalid parameter: used_at=%s", c.LastUsedAt.Format(time.RFC3339))
    }

    return nil
}


//
// Validate status.
//
func (c *APIKeyCredential) ValidateStatus() error {
    // Guard.
    if c == nil {
        return fmt.Errorf("invalid parameter: api_key_credential=null")
    }

    if !c.Status.IsValid() {
        return fmt.Errorf("invalid parameter: status=%d", c.Status)
    }

    return nil
}


//
// Validate name.
//
func (c *APIKeyCredential) ValidateName() error {
    // Guard.
    if c == nil {
        return fmt.Errorf("invalid parameter: api_key_credential=null")
    }

    if c.Name == "" {
        return fmt.Errorf("invalid parameter: name=empty")
    }
    if utf8.RuneCountInString(c.Name) > 64 {
        return fmt.Errorf("invalid parameter: name=%q", myHelper.TruncateRunes(c.Name, 64))
    }

    return nil
}


//
// Validate public token.
//
func (c *APIKeyCredential) ValidatePublicToken() error {
    // Guard.
    if c == nil {
        return fmt.Errorf("invalid parameter: api_key_credential=null")
    }

    if c.PublicToken == "" {
        return fmt.Errorf("invalid parameter: public_token=empty")
    }
    if utf8.RuneCountInString(c.PublicToken) > 128 {
        return fmt.Errorf("invalid parameter: public_token=%q", myHelper.TruncateRunes(c.PublicToken, 128))
    }

    return nil
}


//
// Validate secret token.
//
func (c *APIKeyCredential) ValidateSecretToken() error {
    // Guard.
    if c == nil {
        return fmt.Errorf("invalid parameter: api_key_credential=null")
    }

    if c.SecretToken == "" {
        return fmt.Errorf("invalid parameter: secret_token=empty")
    }
    if utf8.RuneCountInString(c.SecretToken) > 254 {
        return fmt.Errorf("invalid parameter: secret_token=%q", myHelper.TruncateRunes(c.SecretToken, 254))
    }

    return nil
}


//
// Validate scopes.
//
func (c *APIKeyCredential) ValidateScopes() error {
    if c == nil {
        return fmt.Errorf("invalid parameter: api_key_credential=null")
    }

    if len(c.Scopes) == 0 {
        return nil
    }

    seen := make(map[string]struct{}, len(c.Scopes))

    for _, scope := range c.Scopes {
        if scope == "" {
            return fmt.Errorf("invalid parameter: scope=empty")
        }
        if utf8.RuneCountInString(scope) > 254 {
            return fmt.Errorf("invalid parameter: scope=%q", myHelper.TruncateRunes(scope, 254))
        }
        if _, ok := seen[scope]; ok {
            return fmt.Errorf("invalid parameter: scope=%q", myHelper.TruncateRunes(scope, 254))
        }
        seen[scope] = struct{}{}
    }

    b, err := json.Marshal(c.Scopes)
    if err != nil {
        return fmt.Errorf("invalid parameter: %w", err)
    }
    if len(b) > 4096 {
        return fmt.Errorf("invalid parameter: scopes_json_size=%d", len(b))
    }

    return nil
}


//
// Create new account API key credential store.
//
// Version:
//   - 2026-04-30: Added.
//
func NewAPIKeyCredentialStore(db *sql.DB, tableName, accountTableName string) (*APIKeyCredentialStore, error) {
    // Guard.
    if db == nil {
        return nil, fmt.Errorf("failed to create account api key credential store: missing required parameter: db=null")
    }
    if tableName == "" {
        return nil, fmt.Errorf("failed to create account api key credential store: missing required parameter: table_name=empty")
    }
    if accountTableName == "" {
        return nil, fmt.Errorf("failed to create account api key credential store: missing required parameter: account_table_name=empty")
    }

    return &APIKeyCredentialStore{
        db:               db,
        tableName:        tableName,
        accountTableName: accountTableName,
    }, nil
}


//
// Count account API key credentials.
//
// Version:
//   - 2026-04-30: Added.
//
func (s *APIKeyCredentialStore) Count(option *APIKeyCredentialSelectOption) (int64, error) {
    // Guard.
    if s == nil {
        return 0, fmt.Errorf("failed to count account api key credentials: missing required parameter: api_key_credential_store=null")
    }
    if s.db == nil {
        return 0, fmt.Errorf("failed to count account api key credentials: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return 0, fmt.Errorf("failed to count account api key credentials: missing required parameter: table_name=empty")
    }
    if option == nil {
        return 0, fmt.Errorf("failed to count account api key credentials: missing required parameter: select_option=null")
    }

    // Build conditions.
    conditions, args := buildAPIKeyCredentialSelectCondition(option)

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
        return 0, fmt.Errorf("failed to count account api key credentials: %w", err)
    }

    return result, nil
}


//
// Create account API key credentials table.
//
// Version:
//   - 2026-04-30: Added.
//
func (s *APIKeyCredentialStore) CreateTable() error {
    // Guard.
    if s == nil {
        return fmt.Errorf("failed to create account api key credentials table: missing required parameter: api_key_credential_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to create account api key credentials table: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to create account api key credentials table: missing required parameter: table_name=empty")
    }
    if s.accountTableName == "" {
        return fmt.Errorf("failed to create account api key credentials table: missing required parameter: account_table_name=empty")
    }

    // Generate a CREATE TABLE query.
    query := fmt.Sprintf(
        `CREATE TABLE IF NOT EXISTS %s (
            %s BIGINT UNSIGNED NOT NULL COMMENT 'ID',
            %s BIGINT UNSIGNED NOT NULL COMMENT 'Account ID',
            %s TINYINT UNSIGNED NOT NULL COMMENT 'Status',
            %s VARCHAR(64) NOT NULL COMMENT 'Name',
            %s VARCHAR(128) NOT NULL COMMENT 'Public token',
            %s VARCHAR(254) NOT NULL COMMENT 'Secret token',
            %s DATETIME NOT NULL COMMENT 'Expires at',
            %s JSON NULL COMMENT 'Scopes',
            %s DATETIME NULL COMMENT 'Last used at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
            PRIMARY KEY (%s),
            UNIQUE KEY uk_account_api_key_credentials_account_id_name (%s, %s),
            KEY idx_account_api_key_credentials_account_id (%s),
            KEY idx_account_api_key_credentials_status (%s),
            KEY idx_account_api_key_credentials_public_token (%s),
            CONSTRAINT fk_account_api_key_credentials_account_id FOREIGN KEY (%s) REFERENCES %s (%s) ON DELETE CASCADE ON UPDATE CASCADE);
        `,
        s.tableName,
        ColID,
        ColAccountID,
        ColStatus,
        ColName,
        ColPublicToken,
        ColSecretToken,
        ColExpiresAt,
        ColScopes,
        ColLastUsedAt,
        ColCreatedAt,
        ColUpdatedAt,
        ColID,
        ColAccountID, ColName,
        ColAccountID,
        ColStatus,
        ColPublicToken,
        ColAccountID, s.accountTableName, ColID,
    )

    // Execute the query.
    if _, err := s.db.Exec(query); err != nil {
        return fmt.Errorf("failed to create account api key credentials table: %w", err)
    }

    return nil
}


//
// Delete account API key credential by ID.
//
// Version:
//   - 2026-04-30: Added.
//
func (s *APIKeyCredentialStore) DeleteByID(id uint64) error {
    // Guard.
    if s == nil {
        return fmt.Errorf("failed to delete account api key credential by id: missing required parameter: api_key_credential_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to delete account api key credential by id: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to delete account api key credential by id: missing required parameter: table_name=empty")
    }
    if id == 0 {
        return fmt.Errorf("failed to delete account api key credential by id: invalid parameter: id=0")
    }

    // Generate a DELETE query.
    query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", s.tableName, ColID)

    // Execute.
    if _, err := s.db.Exec(query, id); err != nil {
        return fmt.Errorf("failed to delete account api key credential by id: %w", err)
    }

    return nil
}


//
// Insert an account API key credential.
//
// Version:
//   - 2026-04-30: Added.
//
func (s *APIKeyCredentialStore) Insert(row *APIKeyCredential) error {
    // Guard.
    if s == nil {
        return fmt.Errorf("failed to insert account api key credential: missing required parameter: api_key_credential_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to insert account api key credential: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to insert account api key credential: missing required parameter: table_name=empty")
    }
    if row == nil {
        return fmt.Errorf("failed to insert account api key credential: missing required parameter: account_api_key_credential=null table=%q", s.tableName)
    }
    if err := row.ValidateAccountID(); err != nil {
        return fmt.Errorf("failed to insert account api key credential: %w table=%q", err, s.tableName)
    }
    if err := row.ValidateStatus(); err != nil {
        return fmt.Errorf("failed to insert account api key credential: %w table=%q", err, s.tableName)
    }
    if err := row.ValidateName(); err != nil {
        return fmt.Errorf("failed to insert account api key credential: %w table=%q", err, s.tableName)
    }
    if err := row.ValidatePublicToken(); err != nil {
        return fmt.Errorf("failed to insert account api key credential: %w table=%q", err, s.tableName)
    }
    if err := row.ValidateSecretToken(); err != nil {
        return fmt.Errorf("failed to insert account api key credential: %w table=%q", err, s.tableName)
    }
    if err := row.ValidateExpiresAt(); err != nil {
        return fmt.Errorf("failed to insert account api key credential: %w table=%q", err, s.tableName)
    }
    if err := row.ValidateScopes(); err != nil {
        return fmt.Errorf("failed to insert account api key credential: %w table=%q", err, s.tableName)
    }
    if err := row.ValidateLastUsedAt(); err != nil {
        return fmt.Errorf("failed to insert account api key credential: %w table=%q", err, s.tableName)
    }

    // Generate an INSERT query.
    query := fmt.Sprintf(
        "INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
        s.tableName,
        ColID,
        ColAccountID,
        ColStatus,
        ColName,
        ColPublicToken,
        ColSecretToken,
        ColExpiresAt,
        ColScopes,
        ColLastUsedAt,
        ColCreatedAt,
        ColUpdatedAt,
    )

    if row.ID == 0 {
        row.ID = GenerateAPIKeyCredentialID()
    }

    var scopes any
    if len(row.Scopes) > 0 {
        b, err := json.Marshal(row.Scopes)
        if err != nil {
            return fmt.Errorf("failed to insert account api key credential: %w", err)
        }
        scopes = string(b)
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
        row.AccountID,
        row.Status,
        row.Name,
        row.PublicToken,
        row.SecretToken,
        row.ExpiresAt,
        scopes,
        row.LastUsedAt,
        row.CreatedAt,
        row.UpdatedAt,
    ); err != nil {
        return fmt.Errorf("failed to insert account api key credential: %w", err)
    }

    return nil
}


//
// Select account API key credentials.
//
// Version:
//   - 2026-04-30: Added.
//
func (s *APIKeyCredentialStore) Select(option *APIKeyCredentialSelectOption) ([]*APIKeyCredential, error) {
    // Guard.
    if s == nil {
        return nil, fmt.Errorf("failed to select account api key credentials: missing required parameter: api_key_credential_store=null")
    }
    if s.db == nil {
        return nil, fmt.Errorf("failed to select account api key credentials: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account api key credentials: missing required parameter: table_name=empty")
    }
    if option == nil {
        return nil, fmt.Errorf("failed to select account api key credentials: missing required parameter: select_option=null")
    }

    // Build conditions.
    conditions, args := buildAPIKeyCredentialSelectCondition(option)

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
            return nil, fmt.Errorf("failed to select account api key credentials: invalid parameter: offset=%d", option.Offset)
        }
        query.WriteString(" LIMIT ? OFFSET ?")
        args = append(args, option.Limit, option.Offset)
    }

    // Execute.
    rows, err := s.db.Query(query.String(), args...)
    if err != nil {
        return nil, fmt.Errorf("failed to select account api key credentials: %w", err)
    }

    defer rows.Close()

    // Scan.
    var result []*APIKeyCredential
    for rows.Next() {
        row := &APIKeyCredential{}
        var scopes sql.NullString
        err := rows.Scan(
            &row.ID,
            &row.AccountID,
            &row.Status,
            &row.Name,
            &row.PublicToken,
            &row.SecretToken,
            &row.ExpiresAt,
            &scopes,
            &row.LastUsedAt,
            &row.CreatedAt,
            &row.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to select account api key credentials: %w", err)
        }

        if scopes.Valid && scopes.String != "" {
            if err := json.Unmarshal([]byte(scopes.String), &row.Scopes); err != nil {
                return nil, fmt.Errorf("failed to select account api key credentials: %w", err)
            }
        }

        result = append(result, row)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to select account api key credentials: %w", err)
    }

    return result, nil
}


//
// Select account API key credential by ID.
//
// Version:
//   - 2026-04-30: Added.
//
func (s *APIKeyCredentialStore) SelectByID(id uint64) (*APIKeyCredential, error) {
    // Guard.
    if s == nil {
        return nil, fmt.Errorf("failed to select account api key credential by id: missing required parameter: api_key_credential_store=null")
    }
    if s.db == nil {
        return nil, fmt.Errorf("failed to select account api key credential by id: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account api key credential by id: missing required parameter: table_name=empty")
    }
    if id == 0 {
        return nil, fmt.Errorf("failed to select account api key credential by id: invalid parameter: id=0")
    }

    // Generate a SELECT query.
    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColID)

    // Execute.
    row := s.db.QueryRow(query, id)

    // Scan.
    result := &APIKeyCredential{}
    var scopes sql.NullString
    err := row.Scan(
        &result.ID,
        &result.AccountID,
        &result.Status,
        &result.Name,
        &result.PublicToken,
        &result.SecretToken,
        &result.ExpiresAt,
        &scopes,
        &result.LastUsedAt,
        &result.CreatedAt,
        &result.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to select account api key credential by account id: %w", err)
    }

    if scopes.Valid && scopes.String != "" {
        if err := json.Unmarshal([]byte(scopes.String), &result.Scopes); err != nil {
            return nil, fmt.Errorf("failed to select account api key credential by id: %w", err)
        }
    }

    return result, nil
}


//
// Update account API key credential.
//
// Version:
//   - 2026-04-30: Added.
//
func (s *APIKeyCredentialStore) Update(option *APIKeyCredentialUpdateOption) error {
    // Guard.
    if s == nil {
        return fmt.Errorf("failed to update account api key credential: missing required parameter: api_key_credential_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to update account api key credential: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update account api key credential: missing required parameter: table_name=empty")
    }
    if option == nil {
        return fmt.Errorf("failed to update account api key credential: missing required parameter: option=null")
    }
    if option.ID == 0 {
        return fmt.Errorf("failed to update account api key credential: invalid parameter: id=0")
    }

    assignments := make([]string, 0, 8)
    args := make([]interface{}, 0, 8)

    if option.AccountID != nil {
        c := APIKeyCredential{
            AccountID: *option.AccountID,
        }
        if err := c.ValidateAccountID(); err != nil {
            return fmt.Errorf("failed to update account api key credential: %w", err)
        }
        assignments = append(assignments, ColAccountID + " = ?")
        args = append(args, *option.AccountID)
    }

    if option.Status != nil {
        c := APIKeyCredential{
            Status: *option.Status,
        }
        if err := c.ValidateStatus(); err != nil {
            return fmt.Errorf("failed to update account api key credential: %w", err)
        }
        assignments = append(assignments, ColStatus + " = ?")
        args = append(args, *option.Status)
    }

    if option.Name != nil {
        c := APIKeyCredential{
            Name: *option.Name,
        }
        if err := c.ValidateName(); err != nil {
            return fmt.Errorf("failed to update account api key credential: %w", err)
        }
        assignments = append(assignments, ColName + " = ?")
        args = append(args, *option.Name)
    }

    if option.PublicToken != nil {
        c := APIKeyCredential{
            PublicToken: *option.PublicToken,
        }
        if err := c.ValidatePublicToken(); err != nil {
            return fmt.Errorf("failed to update account api key credential: %w", err)
        }
        assignments = append(assignments, ColPublicToken + " = ?")
        args = append(args, *option.PublicToken)
    }

    if option.SecretToken != nil {
        c := APIKeyCredential{
            SecretToken: *option.SecretToken,
        }
        if err := c.ValidateSecretToken(); err != nil {
            return fmt.Errorf("failed to update account api key credential: %w", err)
        }
        assignments = append(assignments, ColSecretToken + " = ?")
        args = append(args, *option.SecretToken)
    }

    if option.ExpiresAt != nil {
        c := APIKeyCredential{
            ExpiresAt: *option.ExpiresAt,
        }
        if err := c.ValidateExpiresAt(); err != nil {
            return fmt.Errorf("failed to update account api key credential: %w", err)
        }
        assignments = append(assignments, ColExpiresAt + " = ?")
        args = append(args, *option.ExpiresAt)
    }

    if len(option.Scopes) > 0 {
        c := APIKeyCredential{
            Scopes: option.Scopes,
        }
        if err := c.ValidateScopes(); err != nil {
            return fmt.Errorf("failed to update account api key credential: %w", err)
        }

        b, err := json.Marshal(option.Scopes)
        if err != nil {
            return fmt.Errorf("failed to update account api key credential: %w", err)
        }

        assignments = append(assignments, ColScopes + " = ?")
        args = append(args, string(b))
    }

    if option.LastUsedAt != nil {
        c := APIKeyCredential{
            LastUsedAt: option.LastUsedAt,
        }
        if err := c.ValidateLastUsedAt(); err != nil {
            return fmt.Errorf("failed to update account api key credential: %w", err)
        }
        assignments = append(assignments, ColLastUsedAt + " = ?")
        args = append(args, *option.LastUsedAt)
    }

    if len(assignments) == 0 {
        return fmt.Errorf("failed to update account api key credential: invalid parameter: assignments=empty")
    }

    args = append(args, option.ID)

    // Generate a UPDATE query.
    query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?;", s.tableName, strings.Join(assignments, ", "), ColID)

    // Execute.
    if _, err := s.db.Exec(query, args...); err != nil {
        return fmt.Errorf("failed to update account api key credential: %w", err)
    }

    return nil
}


//
// Update last used at.
//
// Version:
//   - 2026-04-30: Added.
//
func (s *APIKeyCredentialStore) UpdateLastUsedAt(id uint64) error {
    // Guard.
    if s == nil {
        return fmt.Errorf("failed to update account api key credential last used at: missing required parameter: api_key_credential_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to update account api key credential last used at: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update account api key credential last used at: missing required parameter: table_name=empty")
    }
    if id == 0 {
        return fmt.Errorf("failed to update account api key credential last used at: invalid parameter: id=0")
    }

    now := time.Now()

    return s.Update(&APIKeyCredentialUpdateOption{
        ID:         id,
        LastUsedAt: &now,
    })
}


//
// Build account api key credential select conditions.
//
// Version:
//   - 2026-04-30: Added.
//
func buildAPIKeyCredentialSelectCondition(option *APIKeyCredentialSelectOption) ([]string, []any) {
    conditions := make([]string, 0, 3)
    args := make([]any, 0, 3)

    if option != nil {
        if option.AccountID != nil {
            conditions = append(conditions, ColAccountID + " = ?")
            args = append(args, *option.AccountID)
        }
        if option.Status != nil {
            conditions = append(conditions, ColStatus + " = ?")
            args = append(args, *option.Status)
        }
        if option.NameLike != nil {
            conditions = append(conditions, ColName + " LIKE ?")
            args = append(args, "%" + *option.NameLike + "%")
        }
    }

    return conditions, args
}
