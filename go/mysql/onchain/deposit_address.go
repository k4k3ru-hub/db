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
    accountWalletIDCounter = &helper.IdCounter{}
)


type DepositAddress struct {
    ID                  uint64               `json:"id,string"`
    OwnerID             uint64               `json:"ownerId,string"`
    Status              DepositAddressStatus `json:"status"`
    Chain               Chain                `json:"chain"`
    Network             Network              `json:"network"`
    Address             string               `json:"address"`
    SecretProviderKind  string               `json:"secretProviderKind"`
    SecretKeyVersion    string               `json:"secretKeyVersion"`
    EncryptedPrivateKey string               `json:"encryptedPrivateKey"`
    CreatedAt           time.Time            `json:"createdAt,omitempty"`
    UpdatedAt           time.Time            `json:"updatedAt,omitempty"`
}

type DepositAddressStore struct {
    executor  helper.Executor
    tableName string
}

type DepositAddressInsertParams struct {
    ID                  uint64               `json:"id,string"`
    OwnerID             uint64               `json:"ownerId,string"`
    Status              DepositAddressStatus `json:"status"`
    Chain               Chain                `json:"chain"`
    Network             Network              `json:"network"`
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
    OwnerID     *uint64               `json:"ownerId,string,omitempty"`
    Status      *DepositAddressStatus `json:"status,omitempty"`
    Chain       *Chain                `json:"chain,omitempty"`
    Network     *Network              `json:"network,omitempty"`
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
    return accountWalletIDCounter.GenerateID()
}


//
// Create new onchain deposit address store.
//
// Version:
//   - 2026-05-25: Added.
//
func NewDepositAddressStore(executor helper.Executor, tableName string) (*DepositAddressStore, error) {
    // Guard.
    if executor == nil {
        return nil, fmt.Errorf("failed to create onchain deposit address store: missing required parameter: executor=null")
    }
    if tableName == "" {
        return nil, fmt.Errorf("failed to create onchain deposit address store: missing required parameter: table_name=%q", "empty")
    }

    return &DepositAddressStore{
        executor:  executor,
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
func (w *DepositAddress) ValidateID() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateDepositAddressID(w.ID)
}


//
// Validate onchain deposit address owner ID.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateDepositAddressOwnerID(ownerID uint64) error {
    if ownerID == 0 {
        return fmt.Errorf("invalid parameter: owner_id=0")
    }
    return nil
}


//
// Validate onchain deposit address owner ID.
//
// Version:
//   - 2026-05-25: Added.
//
func (w *DepositAddress) ValidateOwnerID() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateDepositAddressID(w.OwnerID)
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
// Validate onchain deposit address chain.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateDepositAddressChain(c Chain) error {
    if err := c.Validate(); err != nil {
        return err
    }
    return nil
}


//
// Validate onchain deposit address chain.
//
// Version:
//   - 2026-05-25: Added.
//
func (w *DepositAddress) ValidateChain() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateDepositAddressChain(w.Chain)
}


//
// Validate onchain deposit address network.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateDepositAddressNetwork(n Network) error {
    if err := n.Validate(); err != nil {
        return err
    }
    return nil
}


//
// Validate onchain deposit address network.
//
// Version:
//   - 2026-05-25: Added.
//
func (w *DepositAddress) ValidateNetwork() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateDepositAddressNetwork(w.Network)
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
// Create onchain deposit addresss table.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *DepositAddressStore) CreateTable() error {
    if s == nil {
        return fmt.Errorf("failed to create onchain deposit addresss table: missing required parameter: wallet_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to create onchain deposit addresss table: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to create onchain deposit addresss table: missing required parameter: table_name=%q", "empty")
    }

    query := fmt.Sprintf(
        `CREATE TABLE IF NOT EXISTS %s (
            %s BIGINT UNSIGNED NOT NULL COMMENT 'ID',
            %s BIGINT UNSIGNED NOT NULL COMMENT 'Account ID',
            %s TINYINT UNSIGNED NOT NULL COMMENT 'Status',
            %s VARCHAR(64) NOT NULL COMMENT 'Chain',
            %s VARCHAR(64) NOT NULL COMMENT 'Network',
            %s VARCHAR(255) NOT NULL COMMENT 'Address',
            %s VARCHAR(32) NOT NULL COMMENT 'Secret Provider kind',
            %s VARCHAR(32) NOT NULL COMMENT 'Secret key version',
            %s TEXT NOT NULL COMMENT 'Encrypted private key',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
            PRIMARY KEY (%s),
            UNIQUE KEY uq_account_onchain_wallets_chain_network_address (%s, %s, %s),
            KEY idx_account_onchain_wallets_account_chain_network_status (%s, %s, %s, %s)
        );`,
        s.tableName,
        ColID,
        ColOwnerID,
        ColStatus,
        ColChain,
        ColNetwork,
        ColAddress,
        ColSecretProviderKind,
        ColSecretKeyVersion,
        ColEncryptedPrivateKey,
        ColCreatedAt,
        ColUpdatedAt,
        ColID,
        ColChain, ColNetwork, ColAddress,
        ColOwnerID, ColChain, ColNetwork, ColStatus,
    )

    if _, err := s.executor.Exec(query); err != nil {
        return fmt.Errorf("failed to create onchain deposit addresss table: %w", err)
    }

    return nil
}


//
// Insert onchain deposit address.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *DepositAddressStore) Insert(p *DepositAddressInsertParams) error {
    if s == nil {
        return fmt.Errorf("failed to insert onchain deposit address: missing required parameter: wallet_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to insert onchain deposit address: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to insert onchain deposit address: missing required parameter: table_name=%q", "empty")
    }
    if p == nil {
        return fmt.Errorf("failed to insert onchain deposit address: missing required parameter: wallet_insert_params=null")
    }
    if err := ValidateDepositAddressOwnerID(p.OwnerID); err != nil {
        return fmt.Errorf("failed to insert onchain deposit address: %w", err)
    }
    if err := ValidateDepositAddressStatus(p.Status); err != nil {
        return fmt.Errorf("failed to insert onchain deposit address: %w", err)
    }
    if err := ValidateDepositAddressChain(p.Chain); err != nil {
        return fmt.Errorf("failed to insert onchain deposit address: %w", err)
    }
    if err := ValidateDepositAddressNetwork(p.Network); err != nil {
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
        "%s INTO %s (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
        queryPrefix,
        s.tableName,
        ColID,
        ColOwnerID,
        ColStatus,
        ColChain,
        ColNetwork,
        ColAddress,
        ColSecretProviderKind,
        ColSecretKeyVersion,
        ColEncryptedPrivateKey,
        ColCreatedAt,
        ColUpdatedAt,
    )

    if _, err := s.executor.Exec(
        query,
        p.ID,
        p.OwnerID,
        p.Status,
        p.Chain,
        p.Network,
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
// Select onchain deposit addresss.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *DepositAddressStore) Select(p *DepositAddressSelectParams) ([]*DepositAddress, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select onchain deposit addresss: missing required parameter: wallet_store=null")
    }
    if s.executor == nil {
        return nil, fmt.Errorf("failed to select onchain deposit addresss: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select onchain deposit addresss: missing required parameter: table_name=%q", "empty")
    }
    if err := p.Validate(); err != nil {
        return nil, fmt.Errorf("failed to select onchain deposit addresss: %w", err)
    }

    query, args := p.BuildQuery("SELECT * FROM " + s.tableName)

    // Execute.
    rows, err := s.executor.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to select onchain deposit addresss: %w", err)
    }
    defer rows.Close()

    var result []*DepositAddress
    for rows.Next() {
        row := &DepositAddress{}
        if err := rows.Scan(
            &row.ID,
            &row.OwnerID,
            &row.Status,
            &row.Chain,
            &row.Network,
            &row.Address,
            &row.SecretProviderKind,
            &row.SecretKeyVersion,
            &row.EncryptedPrivateKey,
            &row.CreatedAt,
            &row.UpdatedAt,
        ); err != nil {
            return nil, fmt.Errorf("failed to select onchain deposit addresss: %w", err)
        }

        result = append(result, row)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to select onchain deposit addresss: %w", err)
    }

    return result, nil
}


//
// Select onchain deposit address by ID.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *DepositAddressStore) SelectByID(id uint64) (*DepositAddress, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select onchain deposit address by id: missing required parameter: wallet_store=null")
    }
    if s.executor == nil {
        return nil, fmt.Errorf("failed to select onchain deposit address by id: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select onchain deposit address by id: missing required parameter: table_name=%q", "empty")
    }
    if id == 0 {
        return nil, fmt.Errorf("failed to select onchain deposit address by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColID)

    row := s.executor.QueryRow(query, id)

    result := &DepositAddress{}
    err := row.Scan(
        &result.ID,
        &result.OwnerID,
        &result.Status,
        &result.Chain,
        &result.Network,
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
// Count onchain deposit addresss.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *DepositAddressStore) Count(p *DepositAddressSelectParams) (int64, error) {
    if s == nil {
        return 0, fmt.Errorf("failed to count onchain deposit addresss: missing required parameter: wallet_store=null")
    }
    if s.executor == nil {
        return 0, fmt.Errorf("failed to count onchain deposit addresss: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return 0, fmt.Errorf("failed to count onchain deposit addresss: missing required parameter: table_name=%q", "empty")
    }
    if err := p.Validate(); err != nil {
        return 0, fmt.Errorf("failed to count onchain deposit addresss: %w", err)
    }

    query, args := p.BuildQuery("SELECT COUNT(*) FROM " + s.tableName)

    var result int64
    if err := s.executor.QueryRow(query, args...).Scan(&result); err != nil {
        return 0, fmt.Errorf("failed to count onchain deposit addresss: %w", err)
    }

    return result, nil
}


//
// Delete onchain deposit address by ID.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *DepositAddressStore) DeleteByID(id uint64) error {
    if s == nil {
        return fmt.Errorf("failed to delete onchain deposit address by id: missing required parameter: wallet_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to delete onchain deposit address by id: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to delete onchain deposit address by id: missing required parameter: table_name=%q", "empty")
    }
    if id == 0 {
        return fmt.Errorf("failed to delete onchain deposit address by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", s.tableName, ColID)

    if _, err := s.executor.Exec(query, id); err != nil {
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

func (s *DepositAddressStore) Update(p *DepositAddressUpdateParams) error {
    if s == nil {
        return fmt.Errorf("failed to update onchain deposit address: missing required parameter: wallet_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to update onchain deposit address: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update onchain deposit address: missing required parameter: table_name=%q", "empty")
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

    if _, err := s.executor.Exec(query, args...); err != nil {
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

    conditions := make([]string, 0, 6)
    args := make([]any, 0, 8)

    if p.ID != nil {
        conditions = append(conditions, ColID + " = ?")
        args = append(args, *p.ID)
    }
    if p.OwnerID != nil {
        conditions = append(conditions, ColOwnerID + " = ?")
        args = append(args, *p.OwnerID)
    }
    if p.Status != nil {
        conditions = append(conditions, ColStatus + " = ?")
        args = append(args, *p.Status)
    }
    if p.Chain != nil {
        conditions = append(conditions, ColChain + " = ?")
        args = append(args, *p.Chain)
    }
    if p.Network != nil {
        conditions = append(conditions, ColNetwork + " = ?")
        args = append(args, *p.Network)
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
    if p.OwnerID != nil {
        if err := ValidateDepositAddressOwnerID(*p.OwnerID); err != nil {
            return err
        }
    }
    if p.Status != nil {
        if err := ValidateDepositAddressStatus(*p.Status); err != nil {
            return err
        }
    }
    if p.Chain != nil {
        if err := ValidateDepositAddressChain(*p.Chain); err != nil {
            return err
        }
    }
    if p.Network != nil {
        if err := ValidateDepositAddressNetwork(*p.Network); err != nil {
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
            ColOwnerID,
            ColStatus,
            ColChain,
            ColNetwork,
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

