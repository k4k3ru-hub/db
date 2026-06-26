//
// deposit_address.go
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
    DefaultDepositAddressTableName = "onchain_deposit_addresses"
)

var (
    depositAddressIDCounter = &helper.IdCounter{}
)


type DepositAddress struct {
    ID                  uint64               `json:"id,string"`
    OwnerRef            string               `json:"ownerRef"`
    Status              DepositAddressStatus `json:"status"`
    ChainFamily         ChainFamily          `json:"chainFamily"`
    Address             string               `json:"address"`
    SecretProviderKind  string               `json:"secretProviderKind"`
    SecretKeyVersion    string               `json:"secretKeyVersion"`
    EncryptedPrivateKey string               `json:"encryptedPrivateKey"`
    CreatedAt           time.Time            `json:"createdAt,omitempty"`
    UpdatedAt           time.Time            `json:"updatedAt,omitempty"`
}

type DepositAddressStore struct {
    tableName string
}

type DepositAddressInsertParams struct {
    ID                  uint64               `json:"id,string"`
    OwnerRef            string               `json:"ownerRef"`
    Status              DepositAddressStatus `json:"status"`
    ChainFamily         ChainFamily          `json:"chainFamily"`
    Address             string               `json:"address"`
    SecretProviderKind  string               `json:"secretProviderKind"`
    SecretKeyVersion    string               `json:"secretKeyVersion"`
    EncryptedPrivateKey string               `json:"encryptedPrivateKey"`
    CreatedAt           time.Time            `json:"createdAt,omitempty"`
    UpdatedAt           time.Time            `json:"updatedAt,omitempty"`
    Ignore              bool                 `json:"ignore"`
}

type DepositAddressSelectParams struct {
    ID          *uint64               `json:"id,string,omitempty"`
    OwnerRef    *string               `json:"ownerRef,omitempty"`
    Status      *DepositAddressStatus `json:"status,omitempty"`
    ChainFamily *ChainFamily          `json:"chainFamily,omitempty"`
    Address     *string               `json:"address,omitempty"`
    OrderBy     string                `json:"orderBy"`
    OrderByDesc bool                  `json:"orderByDesc"`
    Limit       int                   `json:"limit"`
    Offset      int                   `json:"offset"`
}

type DepositAddressUpdateParams struct {
    ID     uint64                `json:"id,string"`
    Status *DepositAddressStatus `json:"status,omitempty"`
}


//
// Generate new onchain deposit address ID.
//
// Version:
//   - 2026-05-25: Added.
//
func GenerateDepositAddressID() uint64 {
    return depositAddressIDCounter.GenerateID()
}


//
// Create new onchain deposit address store.
//
// Version:
//   - 2026-05-25: Added.
//
func NewDepositAddressStore(tableName string) (*DepositAddressStore, error) {
    // Guard.
    tableName = strings.TrimSpace(tableName)
    if tableName == "" {
        return nil, fmt.Errorf("failed to create onchain deposit address store: missing required parameter: table_name=%q", "empty")
    }

    return &DepositAddressStore{
        tableName: tableName,
    }, nil
}


//
// Validate onchain deposit address ID.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateDepositAddressID(id uint64) error {
    if id == 0 {
        return fmt.Errorf("invalid parameter: id=0")
    }
    return nil
}


//
// Validate onchain deposit address ID.
//
// Version:
//   - 2026-05-25: Added.
//
func (da *DepositAddress) ValidateID() error {
    if da == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateDepositAddressID(da.ID)
}


//
// Validate onchain deposit address owner ref.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateDepositAddressOwnerRef(ownerRef string) error {
    s := strings.TrimSpace(ownerRef)
    if s == "" {
        return fmt.Errorf("invalid parameter: owner_ref=%q", "empty")
    }
    if len(s) > 128 {
        return fmt.Errorf("invalid parameter: owner_ref=%q", "too long")
    }
    return nil
}


//
// Validate onchain deposit address owner ref.
//
// Version:
//   - 2026-05-25: Added.
//
func (da *DepositAddress) ValidateOwnerRef() error {
    if da == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateDepositAddressOwnerRef(da.OwnerRef)
}


//
// Validate onchain deposit address status.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateDepositAddressStatus(status DepositAddressStatus) error {
    if !status.IsValid() {
        return fmt.Errorf("invalid parameter: status=%d", status)
    }
    return nil
}


//
// Validate onchain deposit address status.
//
// Version:
//   - 2026-05-25: Added.
//
func (w *DepositAddress) ValidateStatus() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateDepositAddressStatus(w.Status)
}


//
// Validate onchain deposit address chain family.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateDepositAddressChainFamily(f ChainFamily) error {
    if err := f.Validate(); err != nil {
        return err
    }
    return nil
}


//
// Validate onchain deposit address chain family.
//
// Version:
//   - 2026-05-25: Added.
//
func (a *DepositAddress) ValidateChainFamily() error {
    if a == nil {
        return fmt.Errorf("missing required parameter: deposit_address=null")
    }
    return ValidateDepositAddressChainFamily(a.ChainFamily)
}


//
// Validate onchain deposit address address.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateDepositAddressAddress(address string) error {
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
// Validate onchain deposit address address.
//
// Version:
//   - 2026-05-25: Added.
//
func (w *DepositAddress) ValidateAddress() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateDepositAddressAddress(w.Address)
}


//
// Validate onchain deposit address secret provider kind.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateDepositAddressSecretProviderKind(secretProviderKind string) error {
    s := strings.TrimSpace(secretProviderKind)
    if s == "" {
        return fmt.Errorf("invalid parameter: secret_provider_kind=%q", "empty")
    }
    if len(s) > 32 {
        return fmt.Errorf("invalid parameter: secret_provider_kind=%q", "too long")
    }
    return nil
}


//
// Validate onchain deposit address secret provider kind.
//
// Version:
//   - 2026-05-25: Added.
//
func (w *DepositAddress) ValidateSecretProviderKind() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateDepositAddressSecretProviderKind(w.SecretProviderKind)
}


//
// Validate onchain deposit address secret key version.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateDepositAddressSecretKeyVersion(secretKeyVersion string) error {
    s := strings.TrimSpace(secretKeyVersion)
    if s == "" {        
        return fmt.Errorf("invalid parameter: secret_key_version=%q", "empty")
    }
    if len(s) > 32 {
        return fmt.Errorf("invalid parameter: secret_key_version=%q", "too long")
    }
    return nil
}   
    
    
//  
// Validate onchain deposit address secret key version.
//  
// Version:
//   - 2026-05-25: Added.
//
func (w *DepositAddress) ValidateSecretKeyVersion() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateDepositAddressSecretKeyVersion(w.SecretKeyVersion)
}


//
// Validate onchain deposit address encrypted private key.
//
// Version:
//   - 2026-05-25: Added.
//  
func ValidateDepositAddressEncryptedPrivateKey(encryptedPrivateKey string) error {
    s := strings.TrimSpace(encryptedPrivateKey)
    if s == "" {        
        return fmt.Errorf("invalid parameter: encrypted_private_key=%q", "empty")
    }
    if len(s) > 1024 {
        return fmt.Errorf("invalid parameter: encrypted_private_key=%q", "too long")
    }
    return nil          
}   
    

//  
// Validate onchain deposit address encrypted private key.
//  
// Version:
//   - 2026-05-25: Added.
//
func (w *DepositAddress) ValidateEncryptedPrivateKey() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateDepositAddressEncryptedPrivateKey(w.EncryptedPrivateKey)
}


//
// Create onchain deposit addresses table.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *DepositAddressStore) CreateTable(executor helper.Executor) error {
    if s == nil {
        return fmt.Errorf("failed to create onchain deposit addresses table: missing required parameter: wallet_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to create onchain deposit addresses table: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to create onchain deposit addresses table: missing required parameter: executor=null")
    }

    query := fmt.Sprintf(
        `CREATE TABLE IF NOT EXISTS %s (
            %s BIGINT UNSIGNED NOT NULL COMMENT 'ID',
            %s VARCHAR(128) NOT NULL COMMENT 'Owner reference',
            %s TINYINT UNSIGNED NOT NULL COMMENT 'Status',
            %s VARCHAR(16) NOT NULL COMMENT 'ChainFamily',
            %s VARCHAR(255) NOT NULL COMMENT 'Address',
            %s VARCHAR(32) NOT NULL COMMENT 'Secret Provider kind',
            %s VARCHAR(32) NOT NULL COMMENT 'Secret key version',
            %s TEXT NOT NULL COMMENT 'Encrypted private key',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
            PRIMARY KEY (%s),
            UNIQUE KEY uq_onchain_deposit_addresses_cha_fam_add (%s, %s),
            KEY idx_onchain_deposit_addresses_own_ref_cha_fam_sta (%s, %s, %s)
        );`,
        s.tableName,
        ColID,
        ColOwnerRef,
        ColStatus,
        ColChainFamily,
        ColAddress,
        ColSecretProviderKind,
        ColSecretKeyVersion,
        ColEncryptedPrivateKey,
        ColCreatedAt,
        ColUpdatedAt,
        ColID,
        ColChainFamily, ColAddress,
        ColOwnerRef, ColChainFamily, ColStatus,
    )

    if _, err := executor.Exec(query); err != nil {
        return fmt.Errorf("failed to create onchain deposit addresses table: %w", err)
    }

    return nil
}


//
// Insert onchain deposit address.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *DepositAddressStore) Insert(executor helper.Executor, p *DepositAddressInsertParams) error {
    if s == nil {
        return fmt.Errorf("failed to insert onchain deposit address: missing required parameter: wallet_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to insert onchain deposit address: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to insert onchain deposit address: missing required parameter: executor=null")
    }
    if p == nil {
        return fmt.Errorf("failed to insert onchain deposit address: missing required parameter: wallet_insert_params=null")
    }
    if err := ValidateDepositAddressOwnerRef(p.OwnerRef); err != nil {
        return fmt.Errorf("failed to insert onchain deposit address: %w", err)
    }
    if err := ValidateDepositAddressStatus(p.Status); err != nil {
        return fmt.Errorf("failed to insert onchain deposit address: %w", err)
    }
    if err := ValidateDepositAddressChainFamily(p.ChainFamily); err != nil {
        return fmt.Errorf("failed to insert onchain deposit address: %w", err)
    }
    if err := ValidateDepositAddressAddress(p.Address); err != nil {
        return fmt.Errorf("failed to insert onchain deposit address: %w", err)
    }
    if err := ValidateDepositAddressSecretProviderKind(p.SecretProviderKind); err != nil {
        return fmt.Errorf("failed to insert onchain deposit address: %w", err)
    }
    if err := ValidateDepositAddressSecretKeyVersion(p.SecretKeyVersion); err != nil {
        return fmt.Errorf("failed to insert onchain deposit address: %w", err)
    }
    if err := ValidateDepositAddressEncryptedPrivateKey(p.EncryptedPrivateKey); err != nil {
        return fmt.Errorf("failed to insert onchain deposit address: %w", err)
    }

    if p.ID == 0 {
        p.ID = GenerateDepositAddressID()
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
        ColOwnerRef,
        ColStatus,
        ColChainFamily,
        ColAddress,
        ColSecretProviderKind,
        ColSecretKeyVersion,
        ColEncryptedPrivateKey,
        ColCreatedAt,
        ColUpdatedAt,
    )

    if _, err := executor.Exec(
        query,
        p.ID,
        p.OwnerRef,
        p.Status,
        p.ChainFamily,
        p.Address,
        p.SecretProviderKind,
        p.SecretKeyVersion,
        p.EncryptedPrivateKey,
        p.CreatedAt,
        p.UpdatedAt,
    ); err != nil {
        return fmt.Errorf("failed to insert onchain deposit address: %w", err)
    }

    return nil
}


//
// Select onchain deposit addresses.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *DepositAddressStore) Select(executor helper.Executor, p *DepositAddressSelectParams) ([]*DepositAddress, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select onchain deposit addresses: missing required parameter: wallet_store=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select onchain deposit addresses: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return nil, fmt.Errorf("failed to select onchain deposit addresses: missing required parameter: executor=null")
    }
    if err := p.Validate(); err != nil {
        return nil, fmt.Errorf("failed to select onchain deposit addresses: %w", err)
    }

    query, args := p.BuildQuery("SELECT * FROM " + s.tableName)

    // Execute.
    rows, err := executor.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to select onchain deposit addresses: %w", err)
    }
    defer rows.Close()

    var result []*DepositAddress
    for rows.Next() {
        row := &DepositAddress{}
        if err := rows.Scan(
            &row.ID,
            &row.OwnerRef,
            &row.Status,
            &row.ChainFamily,
            &row.Address,
            &row.SecretProviderKind,
            &row.SecretKeyVersion,
            &row.EncryptedPrivateKey,
            &row.CreatedAt,
            &row.UpdatedAt,
        ); err != nil {
            return nil, fmt.Errorf("failed to select onchain deposit addresses: %w", err)
        }

        result = append(result, row)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to select onchain deposit addresses: %w", err)
    }

    return result, nil
}


//
// Select onchain deposit address by ID.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *DepositAddressStore) SelectByID(executor helper.Executor, id uint64) (*DepositAddress, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select onchain deposit address by id: missing required parameter: wallet_store=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select onchain deposit address by id: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return nil, fmt.Errorf("failed to select onchain deposit address by id: missing required parameter: executor=null")
    }
    if id == 0 {
        return nil, fmt.Errorf("failed to select onchain deposit address by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColID)

    row := executor.QueryRow(query, id)

    result := &DepositAddress{}
    err := row.Scan(
        &result.ID,
        &result.OwnerRef,
        &result.Status,
        &result.ChainFamily,
        &result.Address,
        &result.SecretProviderKind,
        &result.SecretKeyVersion,
        &result.EncryptedPrivateKey,
        &result.CreatedAt,
        &result.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to select onchain deposit address by id: %w", err)
    }

    return result, nil
}


//
// Count onchain deposit addresses.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *DepositAddressStore) Count(executor helper.Executor, p *DepositAddressSelectParams) (int64, error) {
    if s == nil {
        return 0, fmt.Errorf("failed to count onchain deposit addresses: missing required parameter: wallet_store=null")
    }
    if s.tableName == "" {
        return 0, fmt.Errorf("failed to count onchain deposit addresses: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return 0, fmt.Errorf("failed to count onchain deposit addresses: missing required parameter: executor=null")
    }
    if err := p.Validate(); err != nil {
        return 0, fmt.Errorf("failed to count onchain deposit addresses: %w", err)
    }

    query, args := p.BuildQuery("SELECT COUNT(*) FROM " + s.tableName)

    var result int64
    if err := executor.QueryRow(query, args...).Scan(&result); err != nil {
        return 0, fmt.Errorf("failed to count onchain deposit addresses: %w", err)
    }

    return result, nil
}


//
// Delete onchain deposit address by ID.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *DepositAddressStore) DeleteByID(executor helper.Executor, id uint64) error {
    if s == nil {
        return fmt.Errorf("failed to delete onchain deposit address by id: missing required parameter: wallet_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to delete onchain deposit address by id: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to delete onchain deposit address by id: missing required parameter: executor=null")
    }
    if id == 0 {
        return fmt.Errorf("failed to delete onchain deposit address by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", s.tableName, ColID)

    if _, err := executor.Exec(query, id); err != nil {
        return fmt.Errorf("failed to delete onchain deposit address by id: %w", err)
    }

    return nil
}


//
// Update onchain deposit address by ID.
//
// Version:
//   - 2026-05-25: Added.
//

func (s *DepositAddressStore) Update(executor helper.Executor, p *DepositAddressUpdateParams) error {
    if s == nil {
        return fmt.Errorf("failed to update onchain deposit address: missing required parameter: wallet_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update onchain deposit address: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to update onchain deposit address: missing required parameter: executor=null")
    }
    if p == nil {
        return fmt.Errorf("failed to update onchain deposit address: missing required parameter: wallet_update_params=null")
    }
    if p.ID == 0 {
        return fmt.Errorf("failed to update onchain deposit address: invalid parameter: id=0")
    }

    assignments := make([]string, 0, 1)
    args := make([]any, 0, 2)

    if p.Status != nil {
        if err := ValidateDepositAddressStatus(*p.Status); err != nil {
            return fmt.Errorf("failed to update onchain deposit address: %w", err)
        }
        assignments = append(assignments, ColStatus+" = ?")
        args = append(args, *p.Status)
    }

    if len(assignments) == 0 {
        return fmt.Errorf("failed to update onchain deposit address: invalid parameter: assignments=%q", "empty")
    }

    args = append(args, p.ID)

    query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?;", s.tableName, strings.Join(assignments, ", "), ColID)

    if _, err := executor.Exec(query, args...); err != nil {
        return fmt.Errorf("failed to update onchain deposit address: %w", err)
    }

    return nil
}


//
// Build query.
//
// Version:
//   - 2025-05-16: Added.
//
func (p *DepositAddressSelectParams) BuildQuery(selectFromClause string) (string, []any) {
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
    if p.OwnerRef != nil {
        conditions = append(conditions, ColOwnerRef + " = ?")
        args = append(args, *p.OwnerRef)
    }
    if p.Status != nil {
        conditions = append(conditions, ColStatus + " = ?")
        args = append(args, *p.Status)
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
func (p *DepositAddressSelectParams) Validate() error {
    // Guard.
    if p == nil {
        return fmt.Errorf("missing required parameter: wallet_select_params=null")
    }

    if p.ID != nil {
        if err := ValidateDepositAddressID(*p.ID); err != nil {
            return err
        }
    }
    if p.OwnerRef != nil {
        if err := ValidateDepositAddressOwnerRef(*p.OwnerRef); err != nil {
            return err
        }
    }
    if p.Status != nil {
        if err := ValidateDepositAddressStatus(*p.Status); err != nil {
            return err
        }
    }
    if p.ChainFamily != nil {
        if err := ValidateDepositAddressChainFamily(*p.ChainFamily); err != nil {
            return err
        }
    }
    if p.Address != nil {
        if err := ValidateDepositAddressAddress(*p.Address); err != nil {
            return err
        }
    }

    if p.OrderBy != "" {
        switch p.OrderBy {
        case ColID,
            ColOwnerRef,
            ColStatus,
            ColChainFamily,
            ColAddress,
            ColSecretProviderKind,
            ColSecretKeyVersion,
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

