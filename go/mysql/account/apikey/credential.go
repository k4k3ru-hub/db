//
// credential.go
//
package apikey

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "strings"
    "time"
    "unicode/utf8"
    _ "github.com/go-sql-driver/mysql"

    myAccount "github.com/k4k3ru-hub/db/go/mysql/account"
    myHelper  "github.com/k4k3ru-hub/db/go/mysql/helper"
)


const (
    DefaultCredentialTableName = "account_api_key_credentials"
)


var (
    credentialIDCounter = &myHelper.IdCounter{}
)


//
// Credential.
//
// Version:
//   - 2026-05-09: Added.
//
type Credential struct {
    ID          uint64           `json:"id,string"`
    AccountID   uint64           `json:"accountId,string"`
    Status      CredentialStatus `json:"status"`
    Name        string           `json:"name"`
    PublicToken string           `json:"publicToken"`
    SecretToken string           `json:"secretToken"`
    ExpiresAt   time.Time        `json:"expiresAt"`
    Scopes      []string         `json:"scopes"`
    LastUsedAt  *time.Time       `json:"lastUsedAt,omitempty"`
    CreatedAt   time.Time        `json:"createdAt"`
    UpdatedAt   time.Time        `json:"updatedAt"`
}


//
// CredentialStore.
//
// Version:
//   - 2026-05-09: Added.
//
type CredentialStore struct {
    db               *sql.DB
    tableName        string
    accountTableName string
}


//
// CredentialSelectOption.
//
// Version:
//   - 2026-05-09: Added.
//
type CredentialSelectOption struct {
    ID           *uint64           `json:"id,string,omitempty"`
    AccountID    *uint64           `json:"accountId,string,omitempty"`
    Status       *CredentialStatus `json:"status,omitempty"`
    NameLike     *string           `json:"nameLike,omitempty"`
    ExpiresAtGTE *time.Time        `json:"idGte,omitempty"`
    ExpiresAtLTE *time.Time        `json:"idLte,omitempty"`
    OrderBy      string            `json:"orderBy"`
    OrderByDesc  bool              `json:"orderByDesc"`
    Limit        int               `json:"limit"`
    Offset       int               `json:"offset"`
}


//
// CredentialUpdateOption.
//
// Version:
//   - 2026-05-09: Added.
//
type CredentialUpdateOption struct {
    ID          uint64            `json:"id,string"`
    AccountID   *uint64           `json:"accountId,string,omitempty"`
    Status      *CredentialStatus `json:"status,omitempty"`
    Name        *string           `json:"name,omitempty"`
    PublicToken *string           `json:"publicToken,omitempty"`
    SecretToken *string           `json:"secretToken,omitempty"`
    ExpiresAt   *time.Time        `json:"expiresAt,omitempty"`
    Scopes      []string          `json:"scopes,omitempty"`
    LastUsedAt  *time.Time        `json:"lastUsedAt,omitempty"`
}


//
// Validate account API key credential ID.
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
// Validate account API key credential ID.
//
// Version:
//   - 2026-05-09: Added.
//
func (c *Credential) ValidateID() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: account_api_key_credential=null")
    }
    return ValidateCredentialID(c.ID)
}


//
// Validate account API key credential account ID.
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
// Validate account API key credential account ID.
//
// Version:
//   - 2026-05-09: Added.
//
func (c *Credential) ValidateAccountID() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: account_api_key_credential=null")
    }
    return ValidateCredentialID(c.AccountID)
}


//
// Validate account API key credential status.
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
// Validate account API key credential status.
//
// Version:
//   - 2026-05-09: Added.
//
func (c *Credential) ValidateStatus() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: account_api_key_credential=null")
    }
    return ValidateCredentialStatus(c.Status)
}


//
// Validate account API key credential name.
//
// Version:
//   - 2026-05-09: Added.
//
func ValidateCredentialName(name string) error {
    if name == "" {
        return fmt.Errorf("invalid parameter: name=%q", "empty")
    }
    if utf8.RuneCountInString(name) > 64 {
        return fmt.Errorf("invalid parameter: max_length=64 name=%q", "too long")
    }
    return nil
}


//
// Validate account API key credential name.
//
// Version:
//   - 2026-05-09: Added.
//
func (c *Credential) ValidateName() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: account_api_key_credential=null")
    }
    return ValidateCredentialName(c.Name)
}


//
// Validate account API key credential public token.
//
// Version:
//   - 2026-05-09: Added.
//
func ValidateCredentialPublicToken(publicToken string) error {
    if publicToken == "" {
        return fmt.Errorf("invalid parameter: public_token=%q", "empty")
    }
    if utf8.RuneCountInString(publicToken) > 128 {
        return fmt.Errorf("invalid parameter: max_length=64 public_token=%q", "too long")
    }
    return nil
}


//
// Validate account API key credential public token.
//
// Version:
//   - 2026-05-09: Added.
//
func (c *Credential) ValidatePublicToken() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: account_api_key_credential=null")
    }
    return ValidateCredentialPublicToken(c.PublicToken)
}


//
// Validate account API key credential secret token.
//
// Version:
//   - 2026-05-09: Added.
//
func ValidateCredentialSecretToken(secretToken string) error {
    if secretToken == "" {
        return fmt.Errorf("invalid parameter: secret_token=%q", "empty")
    }
    if utf8.RuneCountInString(secretToken) > 255 {
        return fmt.Errorf("invalid parameter: max_length=64 secret_token=%q", "too long")
    }
    return nil
}


//
// Validate account API key credential secret token.
//
// Version:
//   - 2026-05-09: Added.
//
func (c *Credential) ValidateSecretToken() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: account_api_key_credential=null")
    }
    return ValidateCredentialSecretToken(c.SecretToken)
}


//
// Validate account API key credential expires at.
//
// Version:
//   - 2026-05-09: Added.
//
func ValidateCredentialExpiresAt(expiresAt time.Time) error {
    if expiresAt.IsZero() {
        return fmt.Errorf("invalid parameter: expires_at=%q", "empty")
    }
    return nil
}


//
// Validate account API key credential expires at.
//
// Version:
//   - 2026-05-09: Added.
//
func (c *Credential) ValidateExpiresAt() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: account_api_key_credential=null")
    }
    return ValidateCredentialExpiresAt(c.ExpiresAt)
}


//
// Validate account API key credential scopes.
//
// Version:
//   - 2026-05-09: Added.
//
func ValidateCredentialScopes(scopes []string) error {
    if len(scopes) == 0 {
        return nil
    }

    seen := make(map[string]struct{}, len(scopes))

    for _, scope := range scopes {
        if scope == "" {
            return fmt.Errorf("invalid parameter: scope=%q", "empty")
        }
        if _, ok := seen[scope]; ok {
            return fmt.Errorf("invalid parameter: scope=%q", "duplicate")
        }
        seen[scope] = struct{}{}
    }

    b, err := json.Marshal(scopes)
    if err != nil {
        return fmt.Errorf("invalid parameter: %w", err)
    }
    if len(b) > 4096 {
        return fmt.Errorf("invalid parameter: max_size=4096 scopes=%q", "too large")
    }

    return nil
    
}


//
// Validate account API key credential scopes.
//
// Version:
//   - 2026-05-09: Added.
//
func (c *Credential) ValidateScopes() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: account_api_key_credential=null")
    }
    return ValidateCredentialScopes(c.Scopes)
}


//
// Validate account API key credential last used at.
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
    return nil
}


//
// Validate account API key credential last used at.
//
// Version:
//   - 2026-05-09: Added.
//
func (c *Credential) ValidateLastUsedAt() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: account=null")
    }
    return ValidateCredentialLastUsedAt(c.LastUsedAt)
}


//
// Generate account API key credential ID.
//
// Version:
//   - 2026-05-09: Added.
//
func GenerateCredentialID() uint64 {
    return credentialIDCounter.GenerateID()
}


//
// Create new account API key credential store.
//
// Version:
//   - 2026-05-09: Added.
//
func NewCredentialStore(db *sql.DB, tableName, accountTableName string) (*CredentialStore, error) {
    // Guard.
    if db == nil {
        return nil, fmt.Errorf("failed to create account api key credential store: missing required parameter: db=null")
    }
    if tableName == "" {
        return nil, fmt.Errorf("failed to create account api key credential store: missing required parameter: table_name=%q", "empty")
    }

    return &CredentialStore{
        db:               db,
        tableName:        tableName,
        accountTableName: accountTableName,
    }, nil
}


//
// Count account API key credentials.
//
// Version:
//   - 2026-05-09: Added.
//
func (s *CredentialStore) Count(option *CredentialSelectOption) (int64, error) {
    // Guard.
    if s == nil {
        return 0, fmt.Errorf("failed to count account api key credentials: missing required parameter: credential_store=null")
    }
    if s.db == nil {
        return 0, fmt.Errorf("failed to count account api key credentials: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return 0, fmt.Errorf("failed to count account api key credentials: missing required parameter: table_name=%q", "empty")
    }
    if err := option.Validate(); err != nil {
        return 0, fmt.Errorf("failed to count account api key credentials: %w", err)
    }

    query, args := option.BuildQuery("SELECT COUNT(*) FROM " + s.tableName)

    // Execute.
    var result int64
    err := s.db.QueryRow(query, args...).Scan(&result)
    if err != nil {
        return 0, fmt.Errorf("failed to count account api key credentials: %w", err)
    }

    return result, nil
}


//
// Create account API key credentials table.
//
// Version:
//   - 2026-05-09: Added.
//
func (s *CredentialStore) CreateTable() error {
    // Guard.
    if s == nil {
        return fmt.Errorf("failed to create account api key credentials table: missing required parameter: credential_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to create account api key credentials table: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to create account api key credentials table: missing required parameter: table_name=%q", "empty")
    }

     // Generate CREATE TABLE query.
    query := fmt.Sprintf(
        `CREATE TABLE IF NOT EXISTS %s (
            %s BIGINT UNSIGNED NOT NULL COMMENT 'ID',
            %s BIGINT UNSIGNED NOT NULL COMMENT 'Account ID',
            %s TINYINT UNSIGNED NOT NULL COMMENT 'Status',
            %s VARCHAR(64) NOT NULL COMMENT 'Name',
            %s VARCHAR(128) NOT NULL COMMENT 'Public token',
            %s VARCHAR(255) NOT NULL COMMENT 'Secret token',
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
        ColAccountID, s.accountTableName, myAccount.ColID,
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
//   - 2026-05-09: Added.
//
func (s *CredentialStore) DeleteByID(id uint64) error {
    // Guard.
    if s == nil {
        return fmt.Errorf("failed to delete account api key credential by id: missing required parameter: credential_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to delete account api key credential by id: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to delete account api key credential by id: missing required parameter: table_name=%q", "empty")
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
// Insert account API key credential.
//
// Version:
//   - 2026-05-09: Added.
//
func (s *CredentialStore) Insert(row *Credential) error {
    // Guard.
    if s == nil {
        return fmt.Errorf("failed to insert account api key credential: missing required parameter: credential_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to insert account api key credential: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to insert account api key credential: missing required parameter: table_name=%q", "empty")
    }
    if row == nil {
        return fmt.Errorf("failed to insert account api key credential: missing required parameter: account=null table=%q", s.tableName)
    }
    if err := row.ValidateAccountID(); err != nil {
        return fmt.Errorf("failed to insert account api key credential: %w", err)
    }
    if err := row.ValidateStatus(); err != nil {
        return fmt.Errorf("failed to insert account api key credential: %w", err)
    }
    if err := row.ValidateName(); err != nil {
        return fmt.Errorf("failed to insert account api key credential: %w", err)
    }
    if err := row.ValidatePublicToken(); err != nil {
        return fmt.Errorf("failed to insert account api key credential: %w", err)
    }
    if err := row.ValidateSecretToken(); err != nil {
        return fmt.Errorf("failed to insert account api key credential: %w", err)
    }
    if err := row.ValidateExpiresAt(); err != nil {
        return fmt.Errorf("failed to insert account api key credential: %w", err)
    }
    if err := row.ValidateScopes(); err != nil {
        return fmt.Errorf("failed to insert account api key credential: %w", err)
    }
    if err := row.ValidateLastUsedAt(); err != nil {
        return fmt.Errorf("failed to insert account api key credential: %w", err)
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
        row.ID = GenerateCredentialID()
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
//   - 2026-05-09: Added.
//
func (s *CredentialStore) Select(option *CredentialSelectOption) ([]*Credential, error) {
    // Guard.
    if s == nil {
        return nil, fmt.Errorf("failed to select account api key credentials: missing required parameter: credential_store=null")
    }
    if s.db == nil {
        return nil, fmt.Errorf("failed to select account api key credentials: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account api key credentials: missing required parameter: table_name=%q", "empty")
    }
    if err := option.Validate(); err != nil {
        return nil, fmt.Errorf("failed to select account api key credentials: %w", err)
    }

    query, args := option.BuildQuery("SELECT * FROM " + s.tableName)

    // Execute.
    rows, err := s.db.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to select account api key credentials: %w", err)
    }

    defer rows.Close()

    // Scan.
    var result []*Credential
    for rows.Next() {
        row := &Credential{}
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
//   - 2026-05-09: Added.
//
func (s *CredentialStore) SelectByID(id uint64) (*Credential, error) {
    // Guard.
    if s == nil {
        return nil, fmt.Errorf("failed to select account api key credential by id: missing required parameter: credential_store=null")
    }
    if s.db == nil {
        return nil, fmt.Errorf("failed to select account api key credential by id: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account api key credential by id: missing required parameter: table_name=%q", "empty")
    }
    if id == 0 {
        return nil, fmt.Errorf("failed to select account api key credential by id: invalid parameter: id=0")
    }

    // Generate a SELECT query.
    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColID)

    // Execute.
    row := s.db.QueryRow(query, id)

    // Scan.
    result := &Credential{}
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
        return nil, fmt.Errorf("failed to select account api key credential by id: %w", err)
    }

    if scopes.Valid && scopes.String != "" {
        if err := json.Unmarshal([]byte(scopes.String), &result.Scopes); err != nil {
            return nil, fmt.Errorf("failed to select account api key credential by id: %w", err)
        }
    }

    return result, nil
}


//
// Select account API key credential by account ID and name.
//
// Version:
//   - 2026-05-10: Added.
//
func (s *CredentialStore) SelectByAccountIDAndName(accountID uint64, name string) (*Credential, error) {
    // Guard.
    if s == nil {
        return nil, fmt.Errorf("failed to select account api key credential by account id and name: missing required parameter: credential_store=null")
    }
    if s.db == nil {
        return nil, fmt.Errorf("failed to select account api key credential by account id and name: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account api key credential by account id and name: missing required parameter: table_name=%q", "empty")
    }
    if accountID == 0 {
        return nil, fmt.Errorf("failed to select account api key credential by account id and name: invalid parameter: account_id=0")
    }
    if name == "" {
        return nil, fmt.Errorf("failed to select account api key credential by account id and name: invalid parameter: name=%q", "empty")
    }

    // Generate a SELECT query.
    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColID)

    // Execute.
    row := s.db.QueryRow(query, id)

    // Scan.
    result := &Credential{}
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
        return nil, fmt.Errorf("failed to select account api key credential by account id and name: %w", err)
    }

    if scopes.Valid && scopes.String != "" {
        if err := json.Unmarshal([]byte(scopes.String), &result.Scopes); err != nil {
            return nil, fmt.Errorf("failed to select account api key credential by account id and name: %w", err)
        }
    }

    return result, nil
}


//
// Update account API key credential by ID.
//
// Version:
//   - 2026-05-09: Added.
//
func (s *CredentialStore) Update(option *CredentialUpdateOption) error {
    // Guard.
    if s == nil {
        return fmt.Errorf("failed to update account api key credential by id: missing required parameter: credential_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to update account api key credential by id: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update account api key credential by id: missing required parameter: table_name=%q", "empty")
    }
    if option == nil {
        return fmt.Errorf("failed to update account api key credential by id: missing required parameter: option=null")
    }
    if option.ID == 0 {
        return fmt.Errorf("failed to update account api key credential by id: invalid parameter: id=0")
    }

    assignments := make([]string, 0, 8)
    args := make([]interface{}, 0, 9)

    if option.AccountID != nil {
        if err := ValidateCredentialAccountID(*option.AccountID); err != nil {
            return fmt.Errorf("failed to update account api key credential by id: %w", err)
        }
        assignments = append(assignments, ColAccountID + " = ?")
        args = append(args, *option.AccountID)
    }
    if option.Status != nil {
        if err := ValidateCredentialStatus(*option.Status); err != nil {
            return fmt.Errorf("failed to update account api key credential by id: %w", err)
        }
        assignments = append(assignments, ColStatus + " = ?")
        args = append(args, *option.Status)
    }
    if option.Name != nil {
        if err := ValidateCredentialName(*option.Name); err != nil {
            return fmt.Errorf("failed to update account api key credential by id: %w", err)
        }
        assignments = append(assignments, ColName + " = ?")
        args = append(args, *option.Name)
    }
    if option.PublicToken != nil {
        if err := ValidateCredentialPublicToken(*option.PublicToken); err != nil {
            return fmt.Errorf("failed to update account api key credential by id: %w", err)
        }
        assignments = append(assignments, ColPublicToken + " = ?")
        args = append(args, *option.PublicToken)
    }
    if option.SecretToken != nil {
        if err := ValidateCredentialSecretToken(*option.SecretToken); err != nil {
            return fmt.Errorf("failed to update account api key credential by id: %w", err)
        }
        assignments = append(assignments, ColSecretToken + " = ?")
        args = append(args, *option.SecretToken)
    }
    if option.ExpiresAt != nil {
        if err := ValidateCredentialExpiresAt(*option.ExpiresAt); err != nil {
            return fmt.Errorf("failed to update account api key credential by id: %w", err)
        }
        assignments = append(assignments, ColExpiresAt + " = ?")
        args = append(args, *option.ExpiresAt)
    }
    if len(option.Scopes) > 0 {
        if err := ValidateCredentialScopes(option.Scopes); err != nil {
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
        if err := ValidateCredentialLastUsedAt(option.LastUsedAt); err != nil {
            return fmt.Errorf("failed to update account api key credential by id: %w", err)
        }
        assignments = append(assignments, ColLastUsedAt + " = ?")
        args = append(args, *option.LastUsedAt)
    }

    if len(assignments) == 0 {
        return fmt.Errorf("failed to update account api key credential by id: invalid parameter: assignments=%q", "empty")
    }

    args = append(args, option.ID)

    // Generate a UPDATE query.
    query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?;", s.tableName, strings.Join(assignments, ", "), ColID)

    // Execute.
    if _, err := s.db.Exec(query, args...); err != nil {
        return fmt.Errorf("failed to update account api key credential by id: %w", err)
    }

    return nil
}


//
// Build query.
//
// Version:
//   - 2025-05-09: Added.
//
func (o *CredentialSelectOption) BuildQuery(selectFromClause string) (string, []any) {
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
    if o.AccountID != nil {
        conditions = append(conditions, ColAccountID + " = ?")
        args = append(args, *o.AccountID)
    }
    if o.Status != nil {
        conditions = append(conditions, ColStatus + " = ?")
        args = append(args, *o.Status)
    }
    if o.NameLike != nil {
        conditions = append(conditions, ColName + " = ?")
        args = append(args, "%" + *o.NameLike + "%")
    }
    if o.ExpiresAtGTE != nil {
        conditions = append(conditions, ColExpiresAt + " >= ?")
        args = append(args, *o.ExpiresAtGTE)
    }
    if o.ExpiresAtLTE != nil {
        conditions = append(conditions, ColExpiresAt + " >= ?")
        args = append(args, *o.ExpiresAtLTE)
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
// Validate account API key credential select option.
//
// Version:
//   - 2025-05-09: Added.
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
    if o.Status != nil {
        if err := ValidateCredentialStatus(*o.Status); err != nil {
            return err
        }
    }
    if o.NameLike != nil {
        if err := ValidateCredentialName(*o.NameLike); err != nil {
            return err
        }
    }
    if o.ExpiresAtGTE != nil && o.ExpiresAtLTE != nil && o.ExpiresAtGTE.After(*o.ExpiresAtLTE) {
        return fmt.Errorf("invalid parameter: expires_at_gte=%d expires_at_lte=%d", *o.ExpiresAtGTE, *o.ExpiresAtLTE)
    }

    if o.OrderBy != "" {
        switch o.OrderBy {
        case ColID,
            ColAccountID,
            ColStatus,
            ColName,
            ColExpiresAt,
            ColLastUsedAt,
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

