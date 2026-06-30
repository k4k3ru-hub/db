//
// recipient.go
//
package onchain

import (
    "database/sql"
    "fmt"
    "strings"
    "time"

    _ "github.com/go-sql-driver/mysql"

    "github.com/k4k3ru-hub/db/go/mysql/helper"
)

const (
    DefaultRecipientTableName = "payment_onchain_recipients"
)

var (
    recipientIDCounter = &helper.IdCounter{}
)


type Recipient struct {
    ID                  uint64         
    Status              RecipientStatus
    Kind                RecipientKind  
    ChainFamily         ChainFamily    
    Address             string         
    EncryptedPrivateKey *string        
    SecretProviderRef   *string
    CreatedAt           time.Time      
    UpdatedAt           time.Time      
}

type RecipientStore struct {
    tableName string
}

type RecipientInsertParams struct {
    ID                  uint64         
    Status              RecipientStatus
    Kind                RecipientKind  
    ChainFamily         ChainFamily    
    Address             string         
    EncryptedPrivateKey *string        
    SecretProviderRef   *string
    CreatedAt           time.Time
    UpdatedAt           time.Time
    Ignore              bool     
}

type RecipientSelectParams struct {
    ID          *uint64
    Status      *RecipientStatus
    Kind        *RecipientKind
    ChainFamily *ChainFamily
    Address     *string
    OrderBy     string
    OrderByDesc bool
    Limit       int
    Offset      int
}

type RecipientUpdateParams struct {
    ID     uint64           `json:"id,string"`
    Status *RecipientStatus `json:"status,omitempty"`
}


//
// Generate new payment onchain recipient ID.
//
// Version:
//   - 2026-05-25: Added.
//
func GenerateRecipientID() uint64 {
    return recipientIDCounter.GenerateID()
}


//
// Create new payment onchain recipient store.
//
// Version:
//   - 2026-05-25: Added.
//
func NewRecipientStore(tableName string) (*RecipientStore, error) {
    // Guard.
    tableName = strings.TrimSpace(tableName)
    if tableName == "" {
        return nil, fmt.Errorf("failed to create payment onchain recipient store: missing required parameter: table_name=%q", "empty")
    }

    return &RecipientStore{
        tableName: tableName,
    }, nil
}


//
// Validate payment onchain recipient ID.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateRecipientID(id uint64) error {
    if id == 0 {
        return fmt.Errorf("invalid parameter: id=0")
    }
    return nil
}


//
// Validate payment onchain recipient ID.
//
// Version:
//   - 2026-05-25: Added.
//
func (da *Recipient) ValidateID() error {
    if da == nil {
        return fmt.Errorf("missing required parameter: recipient=null")
    }
    return ValidateRecipientID(da.ID)
}


//
// Validate payment onchain recipient status.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateRecipientStatus(status RecipientStatus) error {
    if !status.IsValid() {
        return fmt.Errorf("invalid parameter: status=%d", status)
    }
    return nil
}


//
// Validate payment onchain recipient status.
//
// Version:
//   - 2026-05-25: Added.
//
func (w *Recipient) ValidateStatus() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: recipient=null")
    }
    return ValidateRecipientStatus(w.Status)
}


//
// Validate payment onchain recipient kind.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateRecipientKind(kind RecipientKind) error {
    if err := kind.Validate(); err != nil {
        return err
    }
    return nil
}


//
// Validate payment onchain recipient kind.
//
// Version:
//   - 2026-05-25: Added.
//
func (r *Recipient) ValidateKind() error {
    if r == nil {
        return fmt.Errorf("missing required parameter: recipient=null")
    }
    return ValidateRecipientKind(r.Kind)
}


//
// Validate payment onchain recipient chain family.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateRecipientChainFamily(f ChainFamily) error {
    if err := f.Validate(); err != nil {
        return err
    }
    return nil
}


//
// Validate payment onchain recipient chain family.
//
// Version:
//   - 2026-05-25: Added.
//
func (a *Recipient) ValidateChainFamily() error {
    if a == nil {
        return fmt.Errorf("missing required parameter: deposit_address=null")
    }
    return ValidateRecipientChainFamily(a.ChainFamily)
}


//
// Validate payment onchain recipient address.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateRecipientAddress(address string) error {
    s := strings.TrimSpace(address)
    if s == "" {
        return fmt.Errorf("invalid parameter: address=%q", "empty")
    }
    if len(s) > 255 {
        return fmt.Errorf("invalid parameter: address=%q", "too long")
    }
    return nil
}


//
// Validate payment onchain recipient address.
//
// Version:
//   - 2026-05-25: Added.
//
func (w *Recipient) ValidateAddress() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: recipient=null")
    }
    return ValidateRecipientAddress(w.Address)
}


//
// Validate payment onchain recipient encrypted private key.
//
// Version:
//   - 2026-05-25: Added.
//  
func ValidateRecipientEncryptedPrivateKey(encryptedPrivateKey *string) error {
    if encryptedPrivateKey == nil {
        return nil
    }
    s := strings.TrimSpace(*encryptedPrivateKey)
    if s == "" {        
        return fmt.Errorf("invalid parameter: encrypted_private_key=%q", "empty")
    }
    if len(s) > 1024 {
        return fmt.Errorf("invalid parameter: encrypted_private_key=%q", "too long")
    }
    return nil          
}   
    

//  
// Validate payment onchain recipient encrypted private key.
//  
// Version:
//   - 2026-05-25: Added.
//
func (r *Recipient) ValidateEncryptedPrivateKey() error {
    if r == nil {
        return fmt.Errorf("missing required parameter: recipient=null")
    }
    return ValidateRecipientEncryptedPrivateKey(r.EncryptedPrivateKey)
}


//
// Validate payment onchain recipient secret provider ref.
//
// Version:
//   - 2026-06-30: Added.
//
func ValidateRecipientSecretProviderRef(secretProviderRef *string) error {
    if secretProviderRef == nil {
        return nil
    }
    if *secretProviderRef == "" {
        return fmt.Errorf("invalid parameter: secret_provider_ref=%q", "empty")
    }
    if len(*secretProviderRef) > 128 {
        return fmt.Errorf("invalid parameter: secret_provider_ref=%q max_length=128", "too long")
    }
    return nil
}


//
// Validate payment onchain recipient secret provider ref.
//
// Version:
//   - 2026-06-30: Added.
//
func (r *Recipient) ValidateSecretProviderRef() error {
    if r == nil {
        return fmt.Errorf("missing required parameter: recipient=null")
    }
    return ValidateRecipientSecretProviderRef(r.SecretProviderRef)
}


//
// Create payment onchain recipients table.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *RecipientStore) CreateTable(executor helper.Executor) error {
    if s == nil {
        return fmt.Errorf("failed to create payment onchain recipients table: missing required parameter: wallet_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to create payment onchain recipients table: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to create payment onchain recipients table: missing required parameter: executor=null")
    }

    query := fmt.Sprintf(
        `CREATE TABLE IF NOT EXISTS %s (
            %s BIGINT UNSIGNED NOT NULL COMMENT 'ID',
            %s TINYINT UNSIGNED NOT NULL COMMENT 'Status',
            %s TINYINT UNSIGNED NOT NULL COMMENT 'Kind',
            %s VARCHAR(16) NOT NULL COMMENT 'ChainFamily',
            %s VARCHAR(255) NOT NULL COMMENT 'Address',
            %s TEXT NULL COMMENT 'Encrypted private key',
            %s VARCHAR(128) NULL COMMENT 'Secret provider ref',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
            PRIMARY KEY (%s),
            UNIQUE KEY uq_payment_onchain_recipients_cha_fam_add (%s, %s),
            KEY idx_payment_onchain_recipients_sta_cha_fam (%s, %s, %s)
        );`,
        s.tableName,
        ColID,
        ColStatus,
        ColKind,
        ColChainFamily,
        ColAddress,
        ColEncryptedPrivateKey,
        ColSecretProviderRef,
        ColCreatedAt,
        ColUpdatedAt,
        ColID,
        ColChainFamily, ColAddress,
        ColStatus, ColChainFamily, ColAddress,
    )

    if _, err := executor.Exec(query); err != nil {
        return fmt.Errorf("failed to create payment onchain recipients table: %w", err)
    }

    return nil
}


//
// Insert payment onchain recipient.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *RecipientStore) Insert(executor helper.Executor, p *RecipientInsertParams) error {
    if s == nil {
        return fmt.Errorf("failed to insert payment onchain recipient: missing required parameter: wallet_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to insert payment onchain recipient: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to insert payment onchain recipient: missing required parameter: executor=null")
    }
    if p == nil {
        return fmt.Errorf("failed to insert payment onchain recipient: missing required parameter: wallet_insert_params=null")
    }
    if err := ValidateRecipientStatus(p.Status); err != nil {
        return fmt.Errorf("failed to insert payment onchain recipient: %w", err)
    }
    if err := ValidateRecipientKind(p.Kind); err != nil {
        return fmt.Errorf("failed to insert payment onchain recipient: %w", err)
    }
    if err := ValidateRecipientChainFamily(p.ChainFamily); err != nil {
        return fmt.Errorf("failed to insert payment onchain recipient: %w", err)
    }
    if err := ValidateRecipientAddress(p.Address); err != nil {
        return fmt.Errorf("failed to insert payment onchain recipient: %w", err)
    }
    if err := ValidateRecipientEncryptedPrivateKey(p.EncryptedPrivateKey); err != nil {
        return fmt.Errorf("failed to insert payment onchain recipient: %w", err)
    }
    if err := ValidateRecipientSecretProviderRef(p.SecretProviderRef); err != nil {
        return fmt.Errorf("failed to insert payment onchain recipient: %w", err)
    }

    if p.ID == 0 {
        p.ID = GenerateRecipientID()
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
        "%s INTO %s (%s, %s, %s, %s, %s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);",
        queryPrefix,
        s.tableName,
        ColID,
        ColStatus,
        ColKind,
        ColChainFamily,
        ColAddress,
        ColEncryptedPrivateKey,
        ColSecretProviderRef,
        ColCreatedAt,
        ColUpdatedAt,
    )

    if _, err := executor.Exec(
        query,
        p.ID,
        p.Status,
        p.Kind,
        p.ChainFamily,
        p.Address,
        p.EncryptedPrivateKey,
        p.SecretProviderRef,
        p.CreatedAt,
        p.UpdatedAt,
    ); err != nil {
        return fmt.Errorf("failed to insert payment onchain recipient: %w", err)
    }

    return nil
}


//
// Select payment onchain recipients.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *RecipientStore) Select(executor helper.Executor, p *RecipientSelectParams) ([]*Recipient, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select payment onchain recipients: missing required parameter: wallet_store=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select payment onchain recipients: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return nil, fmt.Errorf("failed to select payment onchain recipients: missing required parameter: executor=null")
    }
    if err := p.Validate(); err != nil {
        return nil, fmt.Errorf("failed to select payment onchain recipients: %w", err)
    }

    query, args := p.BuildQuery("SELECT * FROM " + s.tableName)

    // Execute.
    rows, err := executor.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to select payment onchain recipients: %w", err)
    }
    defer rows.Close()

    var result []*Recipient
    for rows.Next() {
        row := &Recipient{}
        if err := rows.Scan(
            &row.ID,
            &row.Status,
            &row.Kind,
            &row.ChainFamily,
            &row.Address,
            &row.EncryptedPrivateKey,
            &row.SecretProviderRef,
            &row.CreatedAt,
            &row.UpdatedAt,
        ); err != nil {
            return nil, fmt.Errorf("failed to select payment onchain recipients: %w", err)
        }

        result = append(result, row)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to select payment onchain recipients: %w", err)
    }

    return result, nil
}


//
// Select payment onchain recipient by ID.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *RecipientStore) SelectByID(executor helper.Executor, id uint64) (*Recipient, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select payment onchain recipient by id: missing required parameter: wallet_store=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select payment onchain recipient by id: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return nil, fmt.Errorf("failed to select payment onchain recipient by id: missing required parameter: executor=null")
    }
    if id == 0 {
        return nil, fmt.Errorf("failed to select payment onchain recipient by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColID)

    row := executor.QueryRow(query, id)

    result := &Recipient{}
    err := row.Scan(
        &result.ID,
        &result.Status,
        &result.Kind,
        &result.ChainFamily,
        &result.Address,
        &result.EncryptedPrivateKey,
        &result.SecretProviderRef,
        &result.CreatedAt,
        &result.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to select payment onchain recipient by id: %w", err)
    }

    return result, nil
}


//
// Count payment onchain recipients.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *RecipientStore) Count(executor helper.Executor, p *RecipientSelectParams) (int64, error) {
    if s == nil {
        return 0, fmt.Errorf("failed to count payment onchain recipients: missing required parameter: wallet_store=null")
    }
    if s.tableName == "" {
        return 0, fmt.Errorf("failed to count payment onchain recipients: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return 0, fmt.Errorf("failed to count payment onchain recipients: missing required parameter: executor=null")
    }
    if err := p.Validate(); err != nil {
        return 0, fmt.Errorf("failed to count payment onchain recipients: %w", err)
    }

    query, args := p.BuildQuery("SELECT COUNT(*) FROM " + s.tableName)

    var result int64
    if err := executor.QueryRow(query, args...).Scan(&result); err != nil {
        return 0, fmt.Errorf("failed to count payment onchain recipients: %w", err)
    }

    return result, nil
}


//
// Delete payment onchain recipient by ID.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *RecipientStore) DeleteByID(executor helper.Executor, id uint64) error {
    if s == nil {
        return fmt.Errorf("failed to delete payment onchain recipient by id: missing required parameter: wallet_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to delete payment onchain recipient by id: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to delete payment onchain recipient by id: missing required parameter: executor=null")
    }
    if id == 0 {
        return fmt.Errorf("failed to delete payment onchain recipient by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", s.tableName, ColID)

    if _, err := executor.Exec(query, id); err != nil {
        return fmt.Errorf("failed to delete payment onchain recipient by id: %w", err)
    }

    return nil
}


//
// Update payment onchain recipient by ID.
//
// Version:
//   - 2026-05-25: Added.
//

func (s *RecipientStore) Update(executor helper.Executor, p *RecipientUpdateParams) error {
    if s == nil {
        return fmt.Errorf("failed to update payment onchain recipient: missing required parameter: wallet_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update payment onchain recipient: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to update payment onchain recipient: missing required parameter: executor=null")
    }
    if p == nil {
        return fmt.Errorf("failed to update payment onchain recipient: missing required parameter: wallet_update_params=null")
    }
    if p.ID == 0 {
        return fmt.Errorf("failed to update payment onchain recipient: invalid parameter: id=0")
    }

    assignments := make([]string, 0, 1)
    args := make([]any, 0, 2)

    if p.Status != nil {
        if err := ValidateRecipientStatus(*p.Status); err != nil {
            return fmt.Errorf("failed to update payment onchain recipient: %w", err)
        }
        assignments = append(assignments, ColStatus+" = ?")
        args = append(args, *p.Status)
    }

    if len(assignments) == 0 {
        return fmt.Errorf("failed to update payment onchain recipient: invalid parameter: assignments=%q", "empty")
    }

    args = append(args, p.ID)

    query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?;", s.tableName, strings.Join(assignments, ", "), ColID)

    if _, err := executor.Exec(query, args...); err != nil {
        return fmt.Errorf("failed to update payment onchain recipient: %w", err)
    }

    return nil
}


//
// Build query.
//
// Version:
//   - 2025-05-16: Added.
//
func (p *RecipientSelectParams) BuildQuery(selectFromClause string) (string, []any) {
    // Guard.
    if p == nil {
        return selectFromClause, nil
    }

    var query strings.Builder
    query.WriteString(selectFromClause)

    conditions := make([]string, 0, 5)
    args := make([]any, 0, 7)

    if p.ID != nil {
        conditions = append(conditions, ColID + " = ?")
        args = append(args, *p.ID)
    }
    if p.Status != nil {
        conditions = append(conditions, ColStatus + " = ?")
        args = append(args, *p.Status)
    }
    if p.Kind != nil {
        conditions = append(conditions, ColKind + " = ?")
        args = append(args, *p.Kind)
    }
    if p.ChainFamily != nil {
        conditions = append(conditions, ColChainFamily + " = ?")
        args = append(args, *p.ChainFamily)
    }
    if p.Address != nil {
        conditions = append(conditions, ColAddress + " = ?")
        args = append(args, *p.Address)
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
// Validate payment onchain request select params.
//
// Version:
//   - 2025-05-16: Added.
//
func (p *RecipientSelectParams) Validate() error {
    // Guard.
    if p == nil {
        return fmt.Errorf("missing required parameter: wallet_select_params=null")
    }

    if p.ID != nil {
        if err := ValidateRecipientID(*p.ID); err != nil {
            return err
        }
    }
    if p.Status != nil {
        if err := ValidateRecipientStatus(*p.Status); err != nil {
            return err
        }
    }
    if p.Kind != nil {
        if err := ValidateRecipientKind(*p.Kind); err != nil {
            return err
        }
    }
    if p.ChainFamily != nil {
        if err := ValidateRecipientChainFamily(*p.ChainFamily); err != nil {
            return err
        }
    }
    if p.Address != nil {
        if err := ValidateRecipientAddress(*p.Address); err != nil {
            return err
        }
    }

    if p.OrderBy != "" {
        switch p.OrderBy {
        case ColID,
            ColStatus,
            ColKind,
            ColChainFamily,
            ColAddress,
            ColEncryptedPrivateKey,
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

