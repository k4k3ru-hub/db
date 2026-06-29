//
// credential.go
//
package api

import (
    "crypto/rand"
    "database/sql"
    "encoding/base64"
    "fmt"
    "strings"
    "time"
    "unicode/utf8"
    _ "github.com/go-sql-driver/mysql"

    "github.com/k4k3ru-hub/db/go/mysql/account"
    "github.com/k4k3ru-hub/db/go/mysql/helper"
)


const (
    DefaultCredentialTableName = "account_api_credentials"
)


var (
    credentialIDCounter = &helper.IdCounter{}
)


type Credential struct {
    ID                 uint64
    AccountID          uint64
    Status             CredentialStatus
    Name               string
    APIKey             string
    SignatureAlgorithm *CredentialSignatureAlgorithm
    PublicKey          *string
    SecretProviderRef  *string
    SecretRef          *string
    Scopes             CredentialScopes
    ExpiresAt          *time.Time
    CreatedAt          time.Time
    UpdatedAt          time.Time
}

type CredentialStore struct {
    tableName        string
    accountTableName string
}

type CredentialInsertParams struct {
    ID                 uint64
    AccountID          uint64
    Status             CredentialStatus
    Name               string
    APIKey             string
    SignatureAlgorithm *CredentialSignatureAlgorithm
    PublicKey          *string
    SecretProviderRef  *string
    SecretRef          *string
    Scopes             CredentialScopes
    ExpiresAt          *time.Time
    CreatedAt          time.Time
    UpdatedAt          time.Time
    Ignore             bool
}

type CredentialSelectParams struct {
    ID           *uint64
    AccountID    *uint64
    Status       *CredentialStatus
    NameLike     *string
    ExpiresAtGTE *time.Time
    ExpiresAtLTE *time.Time
    OrderBy      string
    OrderByDesc  bool
    Limit        int
    Offset       int
}

type CredentialUpdateParams struct {
    AccountID                 *uint64
    Status                    *CredentialStatus
    Name                      *string
    APIKey                    *string
    SignatureAlgorithm        *CredentialSignatureAlgorithm
    PublicKey                 *string
    SecretProviderRef         *string
    SecretRef                 *string
    Scopes                    *CredentialScopes
    ExpiresAt                 *time.Time
    SetNullSignatureAlgorithm bool
    SetNullPublicKey          bool
    SetNullSecretProviderRef  bool
    SetNullSecretRef          bool
}


//
// Generate account API credential ID.
//
// Version:
//   - 2026-05-09: Added.
//
func GenerateCredentialID() uint64 {
    return credentialIDCounter.GenerateID()
}


//
// Generate account API credential API key.
//
// Version:
//   - 2026-06-03: Added.
//
func GenerateCredentialAPIKey(prefix string) (string, error) {
    if prefix == "" {
        return "", fmt.Errorf("failed to generate api key: missing required parameter: prefix=%q", "empty")
    }

    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return "", fmt.Errorf("failed to generate api key: %w", err)
    }

    return prefix + base64.RawURLEncoding.EncodeToString(b), nil
}


//
// Create new account API credential store.
//
// Version:
//   - 2026-05-09: Added.
//
func NewCredentialStore(tableName, accountTableName string) (*CredentialStore, error) {
    // Guard.
    tableName = strings.TrimSpace(tableName)
    if tableName == "" {
        return nil, fmt.Errorf("failed to create account api credential store: missing required parameter: table_name=%q", "empty")
    }
    accountTableName = strings.TrimSpace(accountTableName)
    if accountTableName == "" {
        return nil, fmt.Errorf("failed to create account api credential store: missing required parameter: account_table_name=%q", "empty")
    }

    return &CredentialStore{
        tableName:        tableName,
        accountTableName: accountTableName,
    }, nil
}


//
// Validate account API credential ID.
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
// Validate account API credential ID.
//
// Version:
//   - 2026-05-09: Added.
//
func (c *Credential) ValidateID() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: account_api_credential=null")
    }
    return ValidateCredentialID(c.ID)
}


//
// Validate account API credential account ID.
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
// Validate account API credential account ID.
//
// Version:
//   - 2026-05-09: Added.
//
func (c *Credential) ValidateAccountID() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: account_api_credential=null")
    }
    return ValidateCredentialID(c.AccountID)
}


//
// Validate account API credential status.
//
// Version:
//   - 2026-05-09: Added.
//
func ValidateCredentialStatus(status CredentialStatus) error {
    if err := status.Validate(); err != nil {
        return err
    }
    return nil
}


//
// Validate account API credential status.
//
// Version:
//   - 2026-05-09: Added.
//
func (c *Credential) ValidateStatus() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: account_api_credential=null")
    }
    return ValidateCredentialStatus(c.Status)
}


//
// Validate account API credential name.
//
// Version:
//   - 2026-05-09: Added.
//
func ValidateCredentialName(name string) error {
    if name == "" {
        return fmt.Errorf("invalid parameter: name=%q", "empty")
    }
    if utf8.RuneCountInString(name) > 64 {
        return fmt.Errorf("invalid parameter: name=%q max_length=64", "too long")
    }
    return nil
}


//
// Validate account API credential name.
//
// Version:
//   - 2026-05-09: Added.
//
func (c *Credential) ValidateName() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: account_api_credential=null")
    }
    return ValidateCredentialName(c.Name)
}


//
// Validate account API credential API key.
//
// Version:
//   - 2026-06-03: Added.
//
func ValidateCredentialAPIKey(apiKey string) error {
    if apiKey == "" {
        return fmt.Errorf("invalid parameter: api_key=%q", "empty")
    }
    if utf8.RuneCountInString(apiKey) > 255 {
        return fmt.Errorf("invalid parameter: api_key=%q max_length=255", "too long")
    }
    return nil
}


//
// Validate account API credential API key.
//
// Version:
//   - 2026-06-03: Added.
//
func (c *Credential) ValidateAPIKey() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: account_api_credential=null")
    }
    return ValidateCredentialAPIKey(c.APIKey)
}


//
// Validate account API credential signature algorithm.
//
// Version:
//   - 2026-06-03: Added.
//
func ValidateCredentialSignatureAlgorithm(algorithm *CredentialSignatureAlgorithm) error {
    if algorithm == nil {
        return nil
    }
    if err := algorithm.Validate(); err != nil {
        return err
    }
    return nil
}


//
// Validate account API credential signature algorithm.
//
// Version:
//   - 2026-06-03: Added.
//
func (c *Credential) ValidateSignatureAlgorithm() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: account_api_credential=null")
    }
    return ValidateCredentialSignatureAlgorithm(c.SignatureAlgorithm)
}


//
// Validate account API credential public key.
//
// Version:
//   - 2026-06-03: Added.
//
func ValidateCredentialPublicKey(publicKey *string) error {
    if publicKey == nil {
        return nil
    }
    if *publicKey == "" {
        return fmt.Errorf("invalid parameter: public_key=%q", "empty")
    }
    if len(*publicKey) > 1024 {
        return fmt.Errorf("invalid parameter: public_key=%q max_length=1024", "too long")
    }
    return nil
}


//
// Validate account API credential public key.
//
// Version:
//   - 2026-06-03: Added.
//
func (c *Credential) ValidatePublicKey() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: account_api_credential=null")
    }
    return ValidateCredentialPublicKey(c.PublicKey)
}


//
// Validate account API credential secret provider ref.
//
// Version:
//   - 2026-06-03: Added.
//
func ValidateCredentialSecretProviderRef(secretProviderRef *string) error {
    if secretProviderRef == nil {
        return nil
    }
    if *secretProviderRef == "" {
        return fmt.Errorf("invalid parameter: secret_provider_ref=%q", "empty")
    }
    if len(*secretProviderRef) > 128 {
        return fmt.Errorf("invalid parameter: secret_provider_ref=%q max_length=32", "too long")
    }
    return nil
}


//
// Validate account API credential secret provider ref.
//
// Version:
//   - 2026-06-03: Added.
//
func (c *Credential) ValidateSecretProviderRef() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: account_api_credential=null")
    }
    return ValidateCredentialSecretProviderRef(c.SecretProviderRef)
}


//
// Validate account API credential secret ref.
//
// Version:
//   - 2026-06-03: Added.
//
func ValidateCredentialSecretRef(secretRef *string) error {
    if secretRef == nil {
        return nil
    }
    if *secretRef == "" {
        return fmt.Errorf("invalid parameter: secret_ref=%q", "empty")
    }
    if len(*secretRef) > 1024 {
        return fmt.Errorf("invalid parameter: secret_ref=%q max_length=1024", "too long")
    }
    return nil
}


//
// Validate account API credential secret ref.
//
// Version:
//   - 2026-06-03: Added.
//
func (c *Credential) ValidateSecretRef() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: account_api_credential=null")
    }
    return ValidateCredentialSecretRef(c.SecretRef)
}


//
// Validate account API credential scopes.
//
// Version:
//   - 2026-05-09: Added.
//
func ValidateCredentialScopes(s CredentialScopes) error {
    return s.Validate()
}


//
// Validate account API credential scopes.
//
// Version:
//   - 2026-05-09: Added.
//
func (c *Credential) ValidateScopes() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: account_api_credential=null")
    }
    return ValidateCredentialScopes(c.Scopes)
}


//
// Validate account API credential expires at.
//
// Version:
//   - 2026-05-09: Added.
//
func ValidateCredentialExpiresAt(expiresAt *time.Time) error {
    if expiresAt == nil {
        return nil
    }
    if expiresAt.IsZero() {
        return fmt.Errorf("invalid parameter: expires_at=%q", "empty")
    }
    return nil
}


//
// Validate account API credential expires at.
//
// Version:
//   - 2026-05-09: Added.
//
func (c *Credential) ValidateExpiresAt() error {
    if c == nil {
        return fmt.Errorf("missing required parameter: account_api_credential=null")
    }
    return ValidateCredentialExpiresAt(c.ExpiresAt)
}


//
// Count account API credentials.
//
// Version:
//   - 2026-05-09: Added.
//
func (s *CredentialStore) Count(executor helper.Executor, option *CredentialSelectParams) (int64, error) {
    // Guard.
    if s == nil {
        return 0, fmt.Errorf("failed to count account api credentials: missing required parameter: credential_store=null")
    }
    if s.tableName == "" {
        return 0, fmt.Errorf("failed to count account api credentials: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return 0, fmt.Errorf("failed to count account api credentials: missing required parameter: executor=null")
    }
    if err := option.Validate(); err != nil {
        return 0, fmt.Errorf("failed to count account api credentials: %w", err)
    }

    query, args := option.BuildQuery("SELECT COUNT(*) FROM " + s.tableName)

    // Execute query.
    var result int64
    err := executor.QueryRow(query, args...).Scan(&result)
    if err != nil {
        return 0, fmt.Errorf("failed to count account api credentials: %w", err)
    }

    return result, nil
}


//
// Create account API credentials table.
//
// Version:
//   - 2026-05-09: Added.
//
func (s *CredentialStore) CreateTable(executor helper.Executor) error {
    // Guard.
    if s == nil {
        return fmt.Errorf("failed to create account api credentials table: missing required parameter: credential_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to create account api credentials table: missing required parameter: table_name=%q", "empty")
    }
    if s.accountTableName == "" {
        return fmt.Errorf("failed to create account api credentials table: missing required parameter: account_table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to create account api credentials table: missing required parameter: executor=null")
    }

     // Generate CREATE TABLE query.
    query := fmt.Sprintf(
        `CREATE TABLE IF NOT EXISTS %s (
            %s BIGINT UNSIGNED NOT NULL COMMENT 'ID',
            %s BIGINT UNSIGNED NOT NULL COMMENT 'Account ID',
            %s TINYINT UNSIGNED NOT NULL COMMENT 'Status',
            %s VARCHAR(64) NOT NULL COMMENT 'Name',
            %s VARCHAR(255) NOT NULL COMMENT 'API key',
            %s TINYINT UNSIGNED NULL COMMENT 'Signature algorithm',
            %s TEXT NULL COMMENT 'public key',
            %s VARCHAR(128) NULL COMMENT 'Secret provider ref',
            %s TEXT NULL COMMENT 'Secret ref',
            %s JSON NOT NULL COMMENT 'Scopes',
            %s DATETIME NULL COMMENT 'Expires at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
            PRIMARY KEY (%s),
            UNIQUE KEY uk_account_api_credentials_api_key (%s),
            UNIQUE KEY uk_account_api_credentials_account_id_name (%s, %s),
            KEY idx_account_api_credentials_account_id (%s),
            CONSTRAINT fk_%s_acc_id FOREIGN KEY (%s) REFERENCES %s (%s) ON DELETE CASCADE ON UPDATE CASCADE);
        `,
        s.tableName,
        ColID,
        ColAccountID,
        ColStatus,
        ColName,
        ColAPIKey,
        ColSignatureAlgorithm,
        ColPublicKey,
        ColSecretProviderRef,
        ColSecretRef,
        ColScopes,
        ColExpiresAt,
        ColCreatedAt,
        ColUpdatedAt,
        ColID,
        ColAPIKey,
        ColAccountID, ColName,
        ColAccountID,
        s.tableName, ColAccountID, s.accountTableName, account.ColID,
    )

    // Execute the query.
    if _, err := executor.Exec(query); err != nil {
        return fmt.Errorf("failed to create account api credentials table: %w", err)
    }

    return nil
}


//
// Delete account API credential by ID.
//
// Version:
//   - 2026-05-09: Added.
//
func (s *CredentialStore) DeleteByID(executor helper.Executor, id uint64) error {
    // Guard.
    if s == nil {
        return fmt.Errorf("failed to delete account api credential by id: missing required parameter: credential_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to delete account api credential by id: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to delete account api credential by id: missing required parameter: executor=null")
    }
    if id == 0 {
        return fmt.Errorf("failed to delete account api credential by id: invalid parameter: id=0")
    }

    // Generate a DELETE query.
    query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", s.tableName, ColID)

    // Execute.
    if _, err := executor.Exec(query, id); err != nil {
        return fmt.Errorf("failed to delete account api credential by id: %w", err)
    }

    return nil
}


//
// Insert account API credential.
//
// Version:
//   - 2026-05-09: Added.
//
func (s *CredentialStore) Insert(executor helper.Executor, p *CredentialInsertParams) error {
    // Guard.
    if s == nil {
        return fmt.Errorf("failed to insert account api credential: missing required parameter: credential_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to insert account api credential: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to insert account api credential: missing required parameter: executor=null")
    }
    if p == nil {
        return fmt.Errorf("failed to insert account api credential: missing required parameter: credential_insert_params=null")
    }
    if err := ValidateCredentialAccountID(p.AccountID); err != nil {
        return fmt.Errorf("failed to insert account api credential: %w", err)
    }
    if err := ValidateCredentialStatus(p.Status); err != nil {
        return fmt.Errorf("failed to insert account api credential: %w", err)
    }
    if err := ValidateCredentialName(p.Name); err != nil {
        return fmt.Errorf("failed to insert account api credential: %w", err)
    }
    if err := ValidateCredentialAPIKey(p.APIKey); err != nil {
        return fmt.Errorf("failed to insert account api credential: %w", err)
    }
    if err := ValidateCredentialSignatureAlgorithm(p.SignatureAlgorithm); err != nil {
        return fmt.Errorf("failed to insert account api credential: %w", err)
    }
    if err := ValidateCredentialPublicKey(p.PublicKey); err != nil {
        return fmt.Errorf("failed to insert account api credential: %w", err)
    }
    if err := ValidateCredentialSecretProviderRef(p.SecretProviderRef); err != nil {
        return fmt.Errorf("failed to insert account api credential: %w", err)
    }
    if err := ValidateCredentialSecretRef(p.SecretRef); err != nil {
        return fmt.Errorf("failed to insert account api credential: %w", err)
    }
    if err := ValidateCredentialScopes(p.Scopes); err != nil {
        return fmt.Errorf("failed to insert account api credential: %w", err)
    }
    if err := ValidateCredentialExpiresAt(p.ExpiresAt); err != nil {
        return fmt.Errorf("failed to insert account api credential: %w", err)
    }

    // Generate an INSERT query.
    query := fmt.Sprintf(
        "INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
        s.tableName,
        ColID,
        ColAccountID,
        ColStatus,
        ColName,
        ColAPIKey,
        ColSignatureAlgorithm,
        ColPublicKey,
        ColSecretProviderRef,
        ColSecretRef,
        ColScopes,
        ColExpiresAt,
        ColCreatedAt,
        ColUpdatedAt,
    )

    if p.ID == 0 {
        p.ID = GenerateCredentialID()
    }

    now := time.Now()
    if p.CreatedAt.IsZero() {
        p.CreatedAt = now
    }
    if p.UpdatedAt.IsZero() {
        p.UpdatedAt = now
    }

    // Execute.
    if _, err := executor.Exec(
        query,
        p.ID,
        p.AccountID,
        p.Status,
        p.Name,
        p.APIKey,
        p.SignatureAlgorithm,
        p.PublicKey,
        p.SecretProviderRef,
        p.SecretRef,
        p.Scopes,
        p.ExpiresAt,
        p.CreatedAt,
        p.UpdatedAt,
    ); err != nil {
        return fmt.Errorf("failed to insert account api credential: %w", err)
    }

    return nil
}


//
// Select account API credentials.
//
// Version:
//   - 2026-05-09: Added.
//
func (s *CredentialStore) Select(executor helper.Executor, option *CredentialSelectParams) ([]*Credential, error) {
    // Guard.
    if s == nil {
        return nil, fmt.Errorf("failed to select account api credentials: missing required parameter: credential_store=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account api credentials: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return nil, fmt.Errorf("failed to select account api credentials: missing required parameter: executor=null")
    }
    if err := option.Validate(); err != nil {
        return nil, fmt.Errorf("failed to select account api credentials: %w", err)
    }

    query, args := option.BuildQuery("SELECT * FROM " + s.tableName)

    // Execute.
    rows, err := executor.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to select account api credentials: %w", err)
    }

    defer rows.Close()

    // Scan.
    var result []*Credential
    for rows.Next() {
        row := &Credential{}
        err := rows.Scan(
            &row.ID,
            &row.AccountID,
            &row.Status,
            &row.Name,
            &row.APIKey,
            &row.SignatureAlgorithm,
            &row.PublicKey,
            &row.SecretProviderRef,
            &row.SecretRef,
            &row.Scopes,
            &row.ExpiresAt,
            &row.CreatedAt,
            &row.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to select account api credentials: %w", err)
        }

        result = append(result, row)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to select account api credentials: %w", err)
    }

    return result, nil
}


//
// Select account API credential by ID.
//
// Version:
//   - 2026-05-09: Added.
//
func (s *CredentialStore) SelectByID(executor helper.Executor, id uint64) (*Credential, error) {
    // Guard.
    if s == nil {
        return nil, fmt.Errorf("failed to select account api credential by id: missing required parameter: credential_store=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account api credential by id: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return nil, fmt.Errorf("failed to select account api credential by id: missing required parameter: executor=null")
    }
    if id == 0 {
        return nil, fmt.Errorf("failed to select account api credential by id: invalid parameter: id=0")
    }

    // Generate a SELECT query.
    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColID)

    // Execute.
    row := executor.QueryRow(query, id)

    // Scan.
    result := &Credential{}
    err := row.Scan(
        &result.ID,
        &result.AccountID,
        &result.Status,
        &result.Name,
        &result.APIKey,
        &result.SignatureAlgorithm,
        &result.PublicKey,
        &result.SecretProviderRef,
        &result.SecretRef,
        &result.Scopes,
        &result.ExpiresAt,
        &result.CreatedAt,
        &result.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to select account api credential by id: %w", err)
    }

    return result, nil
}


//
// Select account API credential by API key.
//
// Version:
//   - 2026-06-03: Added.
//
func (s *CredentialStore) SelectByAPIKey(executor helper.Executor, apiKey string) (*Credential, error) {
    // Guard.
    if s == nil {
        return nil, fmt.Errorf("failed to select account api credential by api key: missing required parameter: credential_store=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account api credential by api key: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return nil, fmt.Errorf("failed to select account api credential by api key: missing required parameter: executor=null")
    }
    if apiKey == "" {
        return nil, fmt.Errorf("failed to select account api credential by api key: invalid parameter: api_key=%q", "empty")
    }

    // Generate SELECT query.
    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColAPIKey)

    // Execute query.
    row := executor.QueryRow(query, apiKey)

    // Scan.
    result := &Credential{}
    err := row.Scan(
        &result.ID,
        &result.AccountID,
        &result.Status,
        &result.Name,
        &result.APIKey,
        &result.SignatureAlgorithm,
        &result.PublicKey,
        &result.SecretProviderRef,
        &result.SecretRef,
        &result.Scopes,
        &result.ExpiresAt,
        &result.CreatedAt,
        &result.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to select account api credential by api key: %w", err)
    }

    return result, nil
}


//
// Select account API credential by account ID and name.
//
// Version:
//   - 2026-05-10: Added.
//
func (s *CredentialStore) SelectByAccountIDAndName(executor helper.Executor, accountID uint64, name string) (*Credential, error) {
    // Guard.
    if s == nil {
        return nil, fmt.Errorf("failed to select account api credential by account id and name: missing required parameter: credential_store=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account api credential by account id and name: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return nil, fmt.Errorf("failed to select account api credential by account id and name: missing required parameter: executor=null")
    }
    if accountID == 0 {
        return nil, fmt.Errorf("failed to select account api credential by account id and name: invalid parameter: account_id=0")
    }
    if name == "" {
        return nil, fmt.Errorf("failed to select account api credential by account id and name: invalid parameter: name=%q", "empty")
    }

    // Generate a SELECT query.
    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? AND %s = ? LIMIT 1;", s.tableName, ColAccountID, ColName)

    // Execute.
    row := executor.QueryRow(query, accountID, name)

    // Scan.
    result := &Credential{}
    err := row.Scan(
        &result.ID,
        &result.AccountID,
        &result.Status,
        &result.Name,
        &result.APIKey,
        &result.SignatureAlgorithm,
        &result.PublicKey,
        &result.SecretProviderRef,
        &result.SecretRef,
        &result.Scopes,
        &result.ExpiresAt,
        &result.CreatedAt,
        &result.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to select account api credential by account id and name: %w", err)
    }

    return result, nil
}


//
// Update account API credential by ID.
//
// Version:
//   - 2026-05-09: Added.
//
func (s *CredentialStore) UpdateByID(executor helper.Executor, params *CredentialUpdateParams, id uint64) error {
    // Guard.
    if s == nil {
        return fmt.Errorf("failed to update account api credential by id: missing required parameter: credential_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update account api credential by id: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to update account api credential by id: missing required parameter: executor=null")
    }
    if params == nil {
        return fmt.Errorf("failed to update account api credential by id: missing required parameter: credential_update_params=null")
    }
    if id == 0 {
        return fmt.Errorf("failed to update account api credential by id: invalid parameter: id=0")
    }

    assignments := make([]string, 0, 11)
    args := make([]any, 0, 12)

    if params.AccountID != nil {
        if err := ValidateCredentialAccountID(*params.AccountID); err != nil {
            return fmt.Errorf("failed to update account api credential by id: %w", err)
        }
        assignments = append(assignments, ColAccountID + " = ?")
        args = append(args, *params.AccountID)
    }
    if params.Status != nil {
        if err := ValidateCredentialStatus(*params.Status); err != nil {
            return fmt.Errorf("failed to update account api credential by id: %w", err)
        }
        assignments = append(assignments, ColStatus + " = ?")
        args = append(args, *params.Status)
    }
    if params.Name != nil {
        if err := ValidateCredentialName(*params.Name); err != nil {
            return fmt.Errorf("failed to update account api credential by id: %w", err)
        }
        assignments = append(assignments, ColName + " = ?")
        args = append(args, *params.Name)
    }
    if params.APIKey != nil {
        if err := ValidateCredentialAPIKey(*params.APIKey); err != nil {
            return fmt.Errorf("failed to update account api credential by id: %w", err)
        }
        assignments = append(assignments, ColAPIKey + " = ?")
        args = append(args, *params.APIKey)
    }
    if params.SetNullSignatureAlgorithm {
        assignments = append(assignments, ColSignatureAlgorithm + " = NULL")
    } else if params.SignatureAlgorithm != nil {
        if err := ValidateCredentialSignatureAlgorithm(params.SignatureAlgorithm); err != nil {
            return fmt.Errorf("failed to update account api credential by id: %w", err)
        }
        assignments = append(assignments, ColSignatureAlgorithm + " = ?")
        args = append(args, *params.SignatureAlgorithm)
    }
    if params.SetNullPublicKey {
        assignments = append(assignments, ColPublicKey + " = NULL")
    } else if params.PublicKey != nil {
        if err := ValidateCredentialSecretRef(params.PublicKey); err != nil {
            return fmt.Errorf("failed to update account api credential by id: %w", err)
        }
        assignments = append(assignments, ColPublicKey + " = ?")
        args = append(args, *params.PublicKey)
    }
    if params.SetNullSecretRef {
        assignments = append(assignments, ColSecretRef + " = NULL")
    } else if params.SecretRef != nil {
        if err := ValidateCredentialSecretRef(params.SecretRef); err != nil {
            return fmt.Errorf("failed to update account api credential by id: %w", err)
        }
        assignments = append(assignments, ColSecretRef + " = ?")
        args = append(args, *params.SecretRef)
    }
    if params.SetNullSecretProviderRef {
        assignments = append(assignments, ColSecretProviderRef + " = NULL")
    } else if params.SecretProviderRef != nil {
        if err := ValidateCredentialSecretProviderRef(params.SecretProviderRef); err != nil {
            return fmt.Errorf("failed to update account api credential by id: %w", err)
        }
        assignments = append(assignments, ColSecretProviderRef + " = ?")
        args = append(args, *params.SecretProviderRef)
    }
    if params.ExpiresAt != nil {
        if err := ValidateCredentialExpiresAt(params.ExpiresAt); err != nil {
            return fmt.Errorf("failed to update account api credential by id: %w", err)
        }
        assignments = append(assignments, ColExpiresAt + " = ?")
        args = append(args, *params.ExpiresAt)
    }
    if params.Scopes != nil {
        if err := ValidateCredentialScopes(*params.Scopes); err != nil {
            return fmt.Errorf("failed to update account api credential by id: %w", err)
        }
        assignments = append(assignments, ColScopes + " = ?")
        args = append(args, params.Scopes)
    }

    if len(assignments) == 0 {
        return fmt.Errorf("failed to update account api credential by id: invalid parameter: assignments=%q", "empty")
    }

    args = append(args, id)

    // Generate UPDATE query.
    query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?;", s.tableName, strings.Join(assignments, ", "), ColID)

    // Execute query.
    if _, err := executor.Exec(query, args...); err != nil {
        return fmt.Errorf("failed to update account api credential by id: %w", err)
    }

    return nil
}


//
// Build query.
//
// Version:
//   - 2025-05-09: Added.
//
func (o *CredentialSelectParams) BuildQuery(selectFromClause string) (string, []any) {
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
// Validate account API credential select option.
//
// Version:
//   - 2025-05-09: Added.
//
func (o *CredentialSelectParams) Validate() error {
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
            ColAPIKey,
            ColExpiresAt,
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

