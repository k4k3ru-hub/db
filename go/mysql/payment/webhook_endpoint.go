//
// webhook_endpoint.go
//
package payment

import (
    "database/sql"
    "fmt"
    "strings"
    "time"

    _ "github.com/go-sql-driver/mysql"

    "github.com/k4k3ru-hub/db/go/mysql/helper"
)

const (
    DefaultWebhookEndpointTableName = "payment_webhook_endpoints"
)

var (
    webhookEndpointIDCounter = &helper.IdCounter{}
)


type WebhookEndpoint struct {
    ID                 uint64                    
    AccountID          uint64                    
    Name               string                    
    URL                string                    
    EncryptedSecretKey string                    
    KMSProviderKind string                    
    KMSKeyVersion   string                    
    SignatureAlgorithm WebhookSignatureAlgorithm 
    CreatedAt          time.Time                 
    UpdatedAt          time.Time                 
}

type WebhookEndpointStore struct {
    tableName string
}


type WebhookEndpointInsertParams struct {
    ID                 uint64                    
    AccountID          uint64                    
    Name               string                    
    URL                string                    
    EncryptedSecretKey string                    
    KMSProviderKind    string                    
    KMSKeyVersion      string                    
    SignatureAlgorithm WebhookSignatureAlgorithm 
    CreatedAt          time.Time                 
    UpdatedAt          time.Time                 
    Ignore             bool                      
}

type WebhookEndpointSelectParams struct {
    ID          *uint64
    AccountID   *uint64
    Name        *string
    OrderBy     string 
    OrderByDesc bool   
    Limit       int    
    Offset      int    
}

type WebhookEndpointUpdateParams struct {
    AccountID          *uint64
    Name               *string
    URL                *string
    EncryptedSecretKey *string
    KMSProviderKind    *string
    KMSKeyVersion      *string
    SignatureAlgorithm *WebhookSignatureAlgorithm
}


//
// Generate payment webhook endpoint ID.
//
// Version:
//   - 2026-06-01: Added.
//
func GenerateWebhookEndpointID() uint64 {
    return webhookEndpointIDCounter.GenerateID()
}


//
// Create payment webhook endpoint store.
//
// Version:
//   - 2026-06-01: Added.
//
func NewWebhookEndpointStore(tableName string) (*WebhookEndpointStore, error) {
    // Guard.
    tableName = strings.TrimSpace(tableName)
    if tableName == "" {
        return nil, fmt.Errorf("failed to create new payment webhook endpoint store: missing required parameter: table_name=%q", "empty")
    }

    return &WebhookEndpointStore{
        tableName: tableName,
    }, nil
}


//
// Validate payment webhook endpoint ID.
//
// Version:
//   - 2026-06-01: Added.
//
func ValidateWebhookEndpointID(id uint64) error {
    if id == 0 {
        return fmt.Errorf("invalid parameter: id=0")
    }
    return nil
}


//
// Validate payment webhook endpoint ID.
//
// Version:
//   - 2026-06-01: Added.
//
func (e *WebhookEndpoint) ValidateID() error {
    if e == nil {
        return fmt.Errorf("missing required parameter: payment_webhook_endpoint=null")
    }
    return ValidateWebhookEndpointID(e.ID)
}


//
// Validate payment webhook endpoint account ID.
//
// Version:
//   - 2026-06-01: Added.
//
func ValidateWebhookEndpointAccountID(accountID uint64) error {
    if accountID == 0 {
        return fmt.Errorf("invalid parameter: account_id=0")
    }
    return nil
}


//
// Validate payment webhook endpoint account ID.
//
// Version:
//   - 2026-06-01: Added.
//
func (e *WebhookEndpoint) ValidateAccountID() error {
    if e == nil {
        return fmt.Errorf("missing required parameter: payment_webhook_endpoint=null")
    }
    return ValidateWebhookEndpointAccountID(e.AccountID)
}


//
// Validate payment webhook endpoint name.
//
// Version:
//   - 2026-06-01: Added.
//
func ValidateWebhookEndpointName(name string) error {
    s := strings.TrimSpace(name)
    if s == "" {
        return fmt.Errorf("invalid parameter: name=%q", "empty")
    }
    if len(s) > 64 {
        return fmt.Errorf("invalid parameter: name=%q max_length=64", "too long")
    }
    return nil
}


//
// Validate payment webhook endpoint name.
//
// Version:
//   - 2026-06-01: Added.
//
func (e *WebhookEndpoint) ValidateName() error {
    if e == nil {
        return fmt.Errorf("missing required parameter: payment_webhook_endpoint=null")
    }
    return ValidateWebhookEndpointName(e.Name)
}


//
// Validate payment webhook endpoint URL.
//
// Version:
//   - 2026-06-01: Added.
//
func ValidateWebhookEndpointURL(url string) error {
    s := strings.TrimSpace(url)
    if s == "" {
        return fmt.Errorf("invalid parameter: url=%q", "empty")
    }
    if len(s) > 2048 {
        return fmt.Errorf("invalid parameter: url=%q max_length=2048", "too long")
    }
    return nil
}


//
// Validate payment webhook endpoint url.
//
// Version:
//   - 2026-06-01: Added.
//
func (e *WebhookEndpoint) ValidateURL() error {
    if e == nil {
        return fmt.Errorf("missing required parameter: payment_webhook_endpoint=null")
    }
    return ValidateWebhookEndpointURL(e.URL)
}


//
// Validate payment webhook endpoint encrypted secret key.
//
// Version:
//   - 2026-06-01: Added.
//
func ValidateWebhookEndpointEncryptedSecretKey(encryptedSecretKey string) error {
    s := strings.TrimSpace(encryptedSecretKey)
    if s == "" {
        return fmt.Errorf("invalid parameter: encrypted_secret_key=%q", "empty")
    }
    if len(s) > 4096 {
        return fmt.Errorf("invalid parameter: encrypted_secret_key=%q max_length=4096", "too long")
    }
    return nil
}


//
// Validate payment webhook endpoint encrypted secret key.
//
// Version:
//   - 2026-06-01: Added.
//
func (e *WebhookEndpoint) ValidateEncryptedSecretKey() error {
    if e == nil {
        return fmt.Errorf("missing required parameter: payment_webhook_endpoint=null")
    }
    return ValidateWebhookEndpointEncryptedSecretKey(e.EncryptedSecretKey)
}


//
// Validate payment webhook endpoint KMS provider kind.
//
// Version:
//   - 2026-06-01: Added.
//
func ValidateWebhookEndpointKMSProviderKind(kmsProviderKind string) error {
    s := strings.TrimSpace(kmsProviderKind)
    if s == "" {
        return fmt.Errorf("invalid parameter: kms_provider_kind=%q", "empty")
    }
    if len(s) > 32 {
        return fmt.Errorf("invalid parameter: kms_provider_kind=%q", "too long")
    }
    return nil
}


//
// Validate payment webhook endpoint KMS provider kind.
//
// Version:
//   - 2026-06-01: Added.
//
func (e *WebhookEndpoint) ValidateKMSProviderKind() error {
    if e == nil {
        return fmt.Errorf("missing required parameter: payment_webhook_endpoint=null")
    }
    return ValidateWebhookEndpointKMSProviderKind(e.KMSProviderKind)
}


//
// Validate payment webhook endpoint KMS key version.
//
// Version:
//   - 2026-06-01: Added.
//
func ValidateWebhookEndpointKMSKeyVersion(kmsKeyVersion string) error {
    s := strings.TrimSpace(kmsKeyVersion)
    if s == "" {
        return fmt.Errorf("invalid parameter: KMS_key_version=%q", "empty")
    }
    if len(s) > 32 {
        return fmt.Errorf("invalid parameter: KMS_key_version=%q", "too long")
    }
    return nil
}


//
// Validate payment webhook endpoint KMS key version.
//
// Version:
//   - 2026-06-01: Added.
//
func (e *WebhookEndpoint) ValidateKMSKeyVersion() error {
    if e == nil {
        return fmt.Errorf("missing required parameter: payment_webhook_endpoint=null")
    }
    return ValidateWebhookEndpointKMSKeyVersion(e.KMSKeyVersion)
}


//
// Validate payment webhook endpoint signature algorithm.
//
// Version:
//   - 2026-06-01: Added.
//
func ValidateWebhookEndpointSignatureAlgorithm(s WebhookSignatureAlgorithm) error {
    if err := s.Validate(); err != nil {
        return err
    }
    return nil
}


//
// Validate payment webhook endpoint signature algorithm.
//
// Version:
//   - 2026-06-01: Added.
//
func (e *WebhookEndpoint) ValidateSignatureAlgorithm() error {
    if e == nil {
        return fmt.Errorf("missing required parameter: payment_webhook_endpoint=null")
    }
    return ValidateWebhookEndpointSignatureAlgorithm(e.SignatureAlgorithm)
}


//
// Create payment webhook endpoints table.
//
// Version:
//   - 2026-06-01: Added.
//
func (s *WebhookEndpointStore) CreateTable(executor helper.Executor) error {
    if s == nil {
        return fmt.Errorf("failed to create payment webhook endpoints table: missing required parameter: webhook_endpoint_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to create payment webhook endpoints table: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to create payment webhook endpoints table: missing required parameter: executor=null")
    }

    query := fmt.Sprintf(
        `CREATE TABLE IF NOT EXISTS %s (
            %s BIGINT UNSIGNED NOT NULL COMMENT 'ID',
            %s BIGINT UNSIGNED NOT NULL COMMENT 'Account ID',
            %s VARCHAR(64) NOT NULL COMMENT 'Name',
            %s VARCHAR(2048) NOT NULL COMMENT 'Webhook endpoint URL',
            %s TEXT NOT NULL COMMENT 'Encrypted secret',
            %s VARCHAR(32) NOT NULL COMMENT 'KMS provider kind',
            %s VARCHAR(32) NOT NULL COMMENT 'KMS key version',
            %s VARCHAR(32) NOT NULL COMMENT 'Signature algorithm',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
            PRIMARY KEY (%s),
            UNIQUE KEY uq_payment_webhook_endpoints_account_id_name (%s, %s),
            KEY idx_payment_webhook_endpoints_account_id (%s)
        );`,
        s.tableName,
        ColID,
        ColAccountID,
        ColName,
        ColURL,
        ColEncryptedSecretKey,
        ColKMSProviderKind,
        ColKMSKeyVersion,
        ColSignatureAlgorithm,
        ColCreatedAt,
        ColUpdatedAt,
        ColID,
        ColAccountID, ColName,
        ColAccountID,
    )

    if _, err := executor.Exec(query); err != nil {
        return fmt.Errorf("failed to create payment webhook endpoints table: %w", err)
    }

    return nil
}


//
// Insert payment webhook endpoint.
//
// Version:
//   - 2026-06-01: Added.
//
func (s *WebhookEndpointStore) Insert(executor helper.Executor, p *WebhookEndpointInsertParams) error {
    if s == nil {
        return fmt.Errorf("failed to insert payment webhook endpoint: missing required parameter: webhook_endpoint_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to insert payment webhook endpoint: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to insert payment webhook endpoint: missing required parameter: executor=null")
    }
    if p == nil {
        return fmt.Errorf("failed to insert payment webhook endpoint: missing required parameter: webhook_endpoint_insert_params=null")
    }
    if err := ValidateWebhookEndpointAccountID(p.AccountID); err != nil {
        return fmt.Errorf("failed to insert payment webhook endpoint: %w", err)
    }
    if err := ValidateWebhookEndpointName(p.Name); err != nil {
        return fmt.Errorf("failed to insert payment webhook endpoint: %w", err)
    }
    if err := ValidateWebhookEndpointURL(p.URL); err != nil {
        return fmt.Errorf("failed to insert payment webhook endpoint: %w", err)
    }
    if err := ValidateWebhookEndpointEncryptedSecretKey(p.EncryptedSecretKey); err != nil {
        return fmt.Errorf("failed to insert payment webhook endpoint: %w", err)
    }
    if err := ValidateWebhookEndpointKMSProviderKind(p.KMSProviderKind); err != nil {
        return fmt.Errorf("failed to insert payment webhook endpoint: %w", err)
    }
    if err := ValidateWebhookEndpointKMSKeyVersion(p.KMSKeyVersion); err != nil {
        return fmt.Errorf("failed to insert payment webhook endpoint: %w", err)
    }
    if err := ValidateWebhookEndpointSignatureAlgorithm(p.SignatureAlgorithm); err != nil {
        return fmt.Errorf("failed to insert payment webhook endpoint: %w", err)
    }

    if p.ID == 0 {
        p.ID = GenerateWebhookEndpointID()
    }

    now := time.Now()
    if p.CreatedAt.IsZero() {
        p.CreatedAt = now
    }
    if p.UpdatedAt.IsZero() {
        p.UpdatedAt = now
    }

    queryPrefix := "INSERT"
    if p.Ignore {
        queryPrefix = "INSERT IGNORE"
    }

    query := fmt.Sprintf(
        "%s INTO %s (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
        queryPrefix,
        s.tableName,
        ColID,
        ColAccountID,
        ColName,
        ColURL,
        ColEncryptedSecretKey,
        ColKMSProviderKind,
        ColKMSKeyVersion,
        ColSignatureAlgorithm,
        ColCreatedAt,
        ColUpdatedAt,
    )

    if _, err := executor.Exec(
        query,
        p.ID,
        p.AccountID,
        p.Name,
        p.URL,
        p.EncryptedSecretKey,
        p.KMSProviderKind,
        p.KMSKeyVersion,
        p.SignatureAlgorithm,
        p.CreatedAt,
        p.UpdatedAt,
    ); err != nil {
        return fmt.Errorf("failed to insert payment webhook endpoint: %w", err)
    }

    return nil
}


//
// Select payment webhook endpoints.
//
// Version:
//   - 2026-06-01: Added.
//
func (s *WebhookEndpointStore) Select(executor helper.Executor, p *WebhookEndpointSelectParams) ([]*WebhookEndpoint, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select payment webhook endpoints: missing required parameter: webhook_endpoint_store=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select payment webhook endpoints: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return nil, fmt.Errorf("failed to select payment webhook endpoints: missing required parameter: executor=null")
    }
    if err := p.Validate(); err != nil {
        return nil, fmt.Errorf("failed to select payment webhook endpoints: %w", err)
    }

    query, args := p.BuildQuery("SELECT * FROM " + s.tableName)

    // Execute.
    rows, err := executor.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to select payment webhook endpoints: %w", err)
    }
    defer rows.Close()

    var result []*WebhookEndpoint
    for rows.Next() {
        row := &WebhookEndpoint{}
        if err := rows.Scan(
            &row.ID,
            &row.AccountID,
            &row.Name,
            &row.URL,
            &row.EncryptedSecretKey,
            &row.KMSProviderKind,
            &row.KMSKeyVersion,
            &row.SignatureAlgorithm,
            &row.CreatedAt,
            &row.UpdatedAt,
        ); err != nil {
            return nil, fmt.Errorf("failed to select payment webhook endpoints: %w", err)
        }

        result = append(result, row)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to select payment webhook endpoints: %w", err)
    }

    return result, nil
}


//
// Select payment webhook endpoint by ID.
//
// Version:
//   - 2026-06-01: Added.
//
func (s *WebhookEndpointStore) SelectByID(executor helper.Executor, id uint64) (*WebhookEndpoint, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select payment webhook endpoint by id: missing required parameter: webhook_endpoint_store=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select payment webhook endpoint by id: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return nil, fmt.Errorf("failed to select payment webhook endpoint by id: missing required parameter: executor=null")
    }
    if id == 0 {
        return nil, fmt.Errorf("failed to select payment webhook endpoint by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColID)

    row := executor.QueryRow(query, id)

    result := &WebhookEndpoint{}
    err := row.Scan(
        &result.ID,
        &result.AccountID,
        &result.Name,
        &result.URL,
        &result.EncryptedSecretKey,
        &result.KMSProviderKind,
        &result.KMSKeyVersion,
        &result.SignatureAlgorithm,
        &result.CreatedAt,
        &result.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to select payment webhook endpoint by id: %w", err)
    }

    return result, nil
}


//
// Select payment webhook endpoint by account ID and name.
//
// Version:
//   - 2026-06-01: Added.
//
func (s *WebhookEndpointStore) SelectByAccountIDAndName(executor helper.Executor, accountID uint64, name string) (*WebhookEndpoint, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select payment webhook endpoint by account id and name: missing required parameter: webhook_endpoint_store=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select payment webhook endpoint by account id and name: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return nil, fmt.Errorf("failed to select payment webhook endpoint by account id and name: missing required parameter: executor=null")
    }
    if accountID == 0 {
        return nil, fmt.Errorf("failed to select payment webhook endpoint by account id and name: missing required parameter: account_id=0")
    }
    if name == "" {
        return nil, fmt.Errorf("failed to select payment webhook endpoint by account id and name: missing required parameter: name=%q", "empty")
    }

    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? AND %s = ? LIMIT 1;", s.tableName, ColAccountID, ColName)

    row := executor.QueryRow(query, accountID, name)

    result := &WebhookEndpoint{}
    err := row.Scan(
        &result.ID,
        &result.AccountID,
        &result.Name,
        &result.URL,
        &result.EncryptedSecretKey,
        &result.KMSProviderKind,
        &result.KMSKeyVersion,
        &result.SignatureAlgorithm,
        &result.CreatedAt,
        &result.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to select payment webhook endpoint by account id and name: %w", err)
    }

    return result, nil
}


//
// Count payment webhook endpoints.
//
// Version:
//   - 2026-06-01: Added.
//
func (s *WebhookEndpointStore) Count(executor helper.Executor, p *WebhookEndpointSelectParams) (int64, error) {
    if s == nil {
        return 0, fmt.Errorf("failed to count payment webhook endpoints: missing required parameter: webhook_endpoint_store=null")
    }
    if s.tableName == "" {
        return 0, fmt.Errorf("failed to count payment webhook endpoints: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return 0, fmt.Errorf("failed to count payment webhook endpoints: missing required parameter: executor=null")
    }
    if err := p.Validate(); err != nil {
        return 0, fmt.Errorf("failed to count payment webhook endpoints: %w", err)
    }

    query, args := p.BuildQuery("SELECT COUNT(*) FROM " + s.tableName)

    var result int64
    if err := executor.QueryRow(query, args...).Scan(&result); err != nil {
        return 0, fmt.Errorf("failed to count payment webhook endpoints: %w", err)
    }

    return result, nil
}


//
// Delete payment webhook endpoint by ID.
//
// Version:
//   - 2026-06-01: Added.
//
func (s *WebhookEndpointStore) DeleteByID(executor helper.Executor, id uint64) error {
    if s == nil {
        return fmt.Errorf("failed to delete payment webhook endpoint by id: missing required parameter: webhook_endpoint_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to delete payment webhook endpoint by id: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to delete payment webhook endpoint by id: missing required parameter: executor=null")
    }
    if id == 0 {
        return fmt.Errorf("failed to delete payment webhook endpoint by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", s.tableName, ColID)

    if _, err := executor.Exec(query, id); err != nil {
        return fmt.Errorf("failed to delete payment webhook endpoint by id: %w", err)
    }

    return nil
}


//
// Update payment webhook endpoint by ID.
//
func (s *WebhookEndpointStore) UpdateByID(executor helper.Executor, id uint64, p *WebhookEndpointUpdateParams) error {
    if s == nil {
        return fmt.Errorf("failed to update payment webhook endpoint by id: missing required parameter: request_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update payment webhook endpoint by id: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to update payment webhook endpoint by id: missing required parameter: executor=null")
    }
    if id == 0 {
        return fmt.Errorf("failed to update payment webhook endpoint by id: invalid parameter: id=0")
    }
    if p == nil {
        return fmt.Errorf("failed to update payment webhook endpoint by id: missing required parameter: webhook_endpoint_update_params=null")
    }

    // Validate webhook endpoint update params.
    if p.AccountID != nil {
        if err := ValidateWebhookEndpointAccountID(*p.AccountID); err != nil {
            return fmt.Errorf("failed to update payment webhook endpoint by id: %w", err)
        }
    }
    if p.Name != nil {
        if err := ValidateWebhookEndpointName(*p.Name); err != nil {
            return fmt.Errorf("failed to update payment webhook endpoint by id: %w", err)
        }
    }
    if p.URL != nil {
        if err := ValidateWebhookEndpointURL(*p.URL); err != nil {
            return fmt.Errorf("failed to update payment webhook endpoint by id: %w", err)
        }
    }
    if p.EncryptedSecretKey != nil {
        if err := ValidateWebhookEndpointEncryptedSecretKey(*p.EncryptedSecretKey); err != nil {
            return fmt.Errorf("failed to update payment webhook endpoint by id: %w", err)
        }
    }
    if p.KMSProviderKind != nil {
        if err := ValidateWebhookEndpointKMSProviderKind(*p.KMSProviderKind); err != nil {
            return fmt.Errorf("failed to update payment webhook endpoint by id: %w", err)
        }
    }
    if p.KMSKeyVersion != nil {
        if err := ValidateWebhookEndpointKMSKeyVersion(*p.KMSKeyVersion); err != nil {
            return fmt.Errorf("failed to update payment webhook endpoint by id: %w", err)
        }
    }
    if p.SignatureAlgorithm != nil {
        if err := ValidateWebhookEndpointSignatureAlgorithm(*p.SignatureAlgorithm); err != nil {
            return fmt.Errorf("failed to update payment webhook endpoint by id: %w", err)
        }
    }

    // Check whether unique key is conflict.
    if p.AccountID != nil || p.Name != nil {
        webhookEndpoint, err := s.SelectByID(executor, id)
        if err != nil {
            return fmt.Errorf("failed to update payment webhook endpoint by id: %w", err)
        }

        accountID := webhookEndpoint.AccountID
        if p.AccountID != nil {
            accountID = *p.AccountID
        }

        name := webhookEndpoint.Name
        if p.Name != nil {
            name = *p.Name
        }

        newWebhookEndpoint, err := s.SelectByAccountIDAndName(executor, accountID, name)
        if err != nil {
            return fmt.Errorf("failed to update payment webhook endpoint by id: %w", err)
        }

        if newWebhookEndpoint != nil {
            return fmt.Errorf("failed to update payment webhook endpoint by id: conflict: account_id=%d name=%q", accountID, name)
        }
    }

    assignments := make([]string, 0, 7)
    args := make([]any, 0, 8)

    if p.AccountID != nil {
        assignments = append(assignments, ColAccountID + " = ?")
        args = append(args, *p.AccountID)
    }
    if p.Name != nil {
        assignments = append(assignments, ColName + " = ?")
        args = append(args, *p.Name)
    }
    if p.URL != nil {
        assignments = append(assignments, ColURL + " = ?")
        args = append(args, *p.URL)
    }
    if p.EncryptedSecretKey != nil {
        assignments = append(assignments, ColEncryptedSecretKey + " = ?")
        args = append(args, *p.EncryptedSecretKey)
    }
    if p.KMSProviderKind != nil {
        assignments = append(assignments, ColKMSProviderKind + " = ?")
        args = append(args, *p.KMSProviderKind)
    }
    if p.KMSKeyVersion != nil {
        assignments = append(assignments, ColKMSKeyVersion + " = ?")
        args = append(args, *p.KMSKeyVersion)
    }
    if p.SignatureAlgorithm != nil {
        assignments = append(assignments, ColSignatureAlgorithm + " = ?")
        args = append(args, *p.SignatureAlgorithm)
    }

    if len(assignments) == 0 {
        return fmt.Errorf("failed to update payment webhook endpoint by id: invalid parameter: assignments=%q", "empty")
    }

    args = append(args, id)

    query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?;", s.tableName, strings.Join(assignments, ", "), ColID)

    if _, err := executor.Exec(query, args...); err != nil {
        return fmt.Errorf("failed to update payment webhook endpoint by id: %w", err)
    }

    return nil
}


//
// Build query.
//
// Version:
//   - 2025-06-01: Added.
//
func (p *WebhookEndpointSelectParams) BuildQuery(selectFromClause string) (string, []any) {
    // Guard.
    if p == nil {
        return selectFromClause, nil
    }

    var query strings.Builder
    query.WriteString(selectFromClause)

    conditions := make([]string, 0, 3)
    args := make([]any, 0, 5)

    if p.ID != nil {
        conditions = append(conditions, ColID + " = ?")
        args = append(args, *p.ID)
    }
    if p.AccountID != nil {
        conditions = append(conditions, ColAccountID + " = ?")
        args = append(args, *p.AccountID)
    }
    if p.Name != nil {
        conditions = append(conditions, ColName + " = ?")
        args = append(args, *p.Name)
    }

    if len(conditions) > 0 {
        query.WriteString(" WHERE ")
        query.WriteString(strings.Join(conditions, " AND "))
    }

    if p.OrderBy != "" {
        query.WriteString(" ORDER BY ")
        query.WriteString(p.OrderBy)
        if p.OrderByDesc {
            query.WriteString(" DESC")
        }
    }

    if p.Limit > 0 {
        query.WriteString(" LIMIT ? OFFSET ?")
        args = append(args, p.Limit, p.Offset)
    }

    return query.String(), args
}


//
// Validate payment webhook endpoint select params.
//
// Version:
//   - 2025-06-01: Added.
//
func (p *WebhookEndpointSelectParams) Validate() error {
    // Guard.
    if p == nil {
        return fmt.Errorf("missing required parameter: payment_webhook_endpoint_select_params=null")
    }

    if p.ID != nil {
        if err := ValidateWebhookEndpointID(*p.ID); err != nil {
            return err
        }
    }
    if p.AccountID != nil {
        if err := ValidateWebhookEndpointAccountID(*p.AccountID); err != nil {
            return err
        }
    }
    if p.Name != nil {
        if err := ValidateWebhookEndpointName(*p.Name); err != nil {
            return err
        }
    }

    if p.OrderBy != "" {
        switch p.OrderBy {
        case ColID,
            ColAccountID,
            ColName,
            ColURL,
            ColKMSProviderKind,
            ColKMSKeyVersion,
            ColSignatureAlgorithm,
            ColCreatedAt,
            ColUpdatedAt:
        default:
            return fmt.Errorf("invalid parameter: order_by=%q", helper.TruncateRunes(p.OrderBy, 32))
        }
    }

    if p.Limit < 0 {
        return fmt.Errorf("invalid parameter: limit=%d", p.Limit)
    }
    if p.Offset < 0 {  
        return fmt.Errorf("invalid parameter: offset=%d", p.Offset)
    }

    return nil
}
