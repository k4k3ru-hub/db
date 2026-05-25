//
// account_wallet.go
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
    DefaultAccountWalletTableName = "onchain_account_wallets"
)

var (
    accountWalletIDCounter = &helper.IdCounter{}
)


type AccountWallet struct {
    ID                  uint64              `json:"id,string"`
    AccountID           uint64              `json:"accountId,string"`
    Status              AccountWalletStatus `json:"status"`
    Chain               Chain               `json:"chain"`
    Network             Network             `json:"network"`
    Address             string              `json:"address"`
    ProviderKind        string              `json:"providerKind"`
    KeyVersion          string              `json:"keyVersion"`
    EncryptedPrivateKey string              `json:"encryptedPrivateKey"`
    CreatedAt           time.Time           `json:"createdAt,omitempty"`
    UpdatedAt           time.Time           `json:"updatedAt,omitempty"`
}

type AccountWalletStore struct {
    executor  helper.Executor
    tableName string
}

type AccountWalletInsertParams struct {
    ID                  uint64              `json:"id,string"`
    AccountID           uint64              `json:"accountId,string"`
    Status              AccountWalletStatus `json:"status"`
    Chain               Chain               `json:"chain"`
    Network             Network             `json:"network"`
    Address             string              `json:"address"`
    ProviderKind        string              `json:"providerKind"`
    KeyVersion          string              `json:"keyVersion"`
    EncryptedPrivateKey string              `json:"encryptedPrivateKey"`
    CreatedAt           time.Time           `json:"createdAt,omitempty"`
    UpdatedAt           time.Time           `json:"updatedAt,omitempty"`
    Ignore              bool                `json:"ignore"`
}

type AccountWalletSelectParams struct {
    ID          *uint64              `json:"id,string,omitempty"`
    AccountID   *uint64              `json:"accountId,string,omitempty"`
    Status      *AccountWalletStatus `json:"status,omitempty"`
    Chain       *Chain               `json:"chain,omitempty"`
    Network     *Network             `json:"network,omitempty"`
    Address     *string              `json:"address,omitempty"`
    OrderBy     string               `json:"orderBy"`
    OrderByDesc bool                 `json:"orderByDesc"`
    Limit       int                  `json:"limit"`
    Offset      int                  `json:"offset"`
}

type AccountWalletUpdateParams struct {
    ID     uint64               `json:"id,string"`
    Status *AccountWalletStatus `json:"status,omitempty"`
}


//
// Generate new onchain account wallet ID.
//
// Version:
//   - 2026-05-25: Added.
//
func GenerateAccountWalletID() uint64 {
    return accountWalletIDCounter.GenerateID()
}


//
// Create new onchain account wallet store.
//
// Version:
//   - 2026-05-25: Added.
//
func NewAccountWalletStore(executor helper.Executor, tableName string) (*AccountWalletStore, error) {
    // Guard.
    if executor == nil {
        return nil, fmt.Errorf("failed to create onchain account wallet store: missing required parameter: executor=null")
    }
    if tableName == "" {
        return nil, fmt.Errorf("failed to create onchain account wallet store: missing required parameter: table_name=%q", "empty")
    }

    return &AccountWalletStore{
        executor:  executor,
        tableName: tableName,
    }, nil
}


//
// Validate onchain account wallet ID.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateAccountWalletID(id uint64) error {
    if id == 0 {
        return fmt.Errorf("invalid parameter: id=0")
    }
    return nil
}


//
// Validate onchain account wallet ID.
//
// Version:
//   - 2026-05-25: Added.
//
func (w *AccountWallet) ValidateID() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateAccountWalletID(w.ID)
}


//
// Validate onchain account wallet account ID.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateAccountWalletAccountID(accountID uint64) error {
    if accountID == 0 {
        return fmt.Errorf("invalid parameter: account_id=0")
    }
    return nil
}


//
// Validate onchain account wallet account ID.
//
// Version:
//   - 2026-05-25: Added.
//
func (w *AccountWallet) ValidateAccountID() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateAccountWalletID(w.AccountID)
}


//
// Validate onchain account wallet status.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateAccountWalletStatus(status AccountWalletStatus) error {
    if !status.IsValid() {
        return fmt.Errorf("invalid parameter: status=%d", status)
    }
    return nil
}


//
// Validate onchain account wallet status.
//
// Version:
//   - 2026-05-25: Added.
//
func (w *AccountWallet) ValidateStatus() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateAccountWalletStatus(w.Status)
}


//
// Validate onchain account wallet chain.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateAccountWalletChain(c Chain) error {
    if err := c.Validate(); err != nil {
        return err
    }
    return nil
}


//
// Validate onchain account wallet chain.
//
// Version:
//   - 2026-05-25: Added.
//
func (w *AccountWallet) ValidateChain() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateAccountWalletChain(w.Chain)
}


//
// Validate onchain account wallet network.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateAccountWalletNetwork(n Network) error {
    if err := n.Validate(); err != nil {
        return err
    }
    return nil
}


//
// Validate onchain account wallet network.
//
// Version:
//   - 2026-05-25: Added.
//
func (w *AccountWallet) ValidateNetwork() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateAccountWalletNetwork(w.Network)
}


//
// Validate onchain account wallet address.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateAccountWalletAddress(address string) error {
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
// Validate onchain account wallet address.
//
// Version:
//   - 2026-05-25: Added.
//
func (w *AccountWallet) ValidateAddress() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateAccountWalletAddress(w.Address)
}


//
// Validate onchain account wallet provider kind.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateAccountWalletProviderKind(providerKind string) error {
    s := strings.TrimSpace(providerKind)
    if s == "" {
        return fmt.Errorf("invalid parameter: provider_kind=%q", "empty")
    }
    if len(s) > 32 {
        return fmt.Errorf("invalid parameter: provider_kind=%q", "too long")
    }
    return nil
}


//
// Validate onchain account wallet provider kind.
//
// Version:
//   - 2026-05-25: Added.
//
func (w *AccountWallet) ValidateProviderKind() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateAccountWalletProviderKind(w.ProviderKind)
}


//
// Validate onchain account wallet key version.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateAccountWalletKeyVersion(keyVersion string) error {
    s := strings.TrimSpace(keyVersion)
    if s == "" {        
        return fmt.Errorf("invalid parameter: key_version=%q", "empty")
    }
    if len(s) > 32 {
        return fmt.Errorf("invalid parameter: key_version=%q", "too long")
    }
    return nil
}   
    
    
//  
// Validate onchain account wallet key version.
//  
// Version:
//   - 2026-05-25: Added.
//
func (w *AccountWallet) ValidateKeyVersion() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateAccountWalletKeyVersion(w.KeyVersion)
}


//
// Validate onchain account wallet encrypted private key.
//
// Version:
//   - 2026-05-25: Added.
//  
func ValidateAccountWalletEncryptedPrivateKey(encryptedPrivateKey string) error {
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
// Validate onchain account wallet encrypted private key.
//  
// Version:
//   - 2026-05-25: Added.
//
func (w *AccountWallet) ValidateEncryptedPrivateKey() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateAccountWalletEncryptedPrivateKey(w.EncryptedPrivateKey)
}


//
// Create onchain account wallets table.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *AccountWalletStore) CreateTable() error {
    if s == nil {
        return fmt.Errorf("failed to create onchain account wallets table: missing required parameter: wallet_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to create onchain account wallets table: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to create onchain account wallets table: missing required parameter: table_name=%q", "empty")
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
            %s VARCHAR(32) NOT NULL COMMENT 'Secret Provider key version',
            %s TEXT NOT NULL COMMENT 'Encrypted private key',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
            PRIMARY KEY (%s),
            UNIQUE KEY uq_account_onchain_wallets_chain_network_address (%s, %s, %s),
            KEY idx_account_onchain_wallets_account_chain_network_status (%s, %s, %s, %s)
        );`,
        s.tableName,
        ColID,
        ColAccountID,
        ColStatus,
        ColChain,
        ColNetwork,
        ColAddress,
        ColProviderKind,
        ColKeyVersion,
        ColEncryptedPrivateKey,
        ColCreatedAt,
        ColUpdatedAt,
        ColID,
        ColChain, ColNetwork, ColAddress,
        ColAccountID, ColChain, ColNetwork, ColStatus,
    )

    if _, err := s.executor.Exec(query); err != nil {
        return fmt.Errorf("failed to create onchain account wallets table: %w", err)
    }

    return nil
}


//
// Insert onchain account wallet.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *AccountWalletStore) Insert(p *AccountWalletInsertParams) error {
    if s == nil {
        return fmt.Errorf("failed to insert onchain account wallet: missing required parameter: wallet_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to insert onchain account wallet: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to insert onchain account wallet: missing required parameter: table_name=%q", "empty")
    }
    if p == nil {
        return fmt.Errorf("failed to insert onchain account wallet: missing required parameter: wallet_insert_params=null")
    }
    if err := ValidateAccountWalletAccountID(p.AccountID); err != nil {
        return fmt.Errorf("failed to insert onchain account wallet: %w", err)
    }
    if err := ValidateAccountWalletStatus(p.Status); err != nil {
        return fmt.Errorf("failed to insert onchain account wallet: %w", err)
    }
    if err := ValidateAccountWalletChain(p.Chain); err != nil {
        return fmt.Errorf("failed to insert onchain account wallet: %w", err)
    }
    if err := ValidateAccountWalletNetwork(p.Network); err != nil {
        return fmt.Errorf("failed to insert onchain account wallet: %w", err)
    }
    if err := ValidateAccountWalletAddress(p.Address); err != nil {
        return fmt.Errorf("failed to insert onchain account wallet: %w", err)
    }
    if err := ValidateAccountWalletProviderKind(p.ProviderKind); err != nil {
        return fmt.Errorf("failed to insert onchain account wallet: %w", err)
    }
    if err := ValidateAccountWalletKeyVersion(p.KeyVersion); err != nil {
        return fmt.Errorf("failed to insert onchain account wallet: %w", err)
    }
    if err := ValidateAccountWalletEncryptedPrivateKey(p.EncryptedPrivateKey); err != nil {
        return fmt.Errorf("failed to insert onchain account wallet: %w", err)
    }

    if p.ID == 0 {
        p.ID = GenerateAccountWalletID()
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
        ColAccountID,
        ColStatus,
        ColChain,
        ColNetwork,
        ColAddress,
        ColProviderKind,
        ColKeyVersion,
        ColEncryptedPrivateKey,
        ColCreatedAt,
        ColUpdatedAt,
    )

    if _, err := s.executor.Exec(
        query,
        p.ID,
        p.AccountID,
        p.Status,
        p.Chain,
        p.Network,
        p.Address,
        p.ProviderKind,
        p.KeyVersion,
        p.EncryptedPrivateKey,
        p.CreatedAt,
        p.UpdatedAt,
    ); err != nil {
        return fmt.Errorf("failed to insert onchain account wallet: %w", err)
    }

    return nil
}


//
// Select onchain account wallets.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *AccountWalletStore) Select(p *AccountWalletSelectParams) ([]*AccountWallet, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select onchain account wallets: missing required parameter: wallet_store=null")
    }
    if s.executor == nil {
        return nil, fmt.Errorf("failed to select onchain account wallets: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select onchain account wallets: missing required parameter: table_name=%q", "empty")
    }
    if err := p.Validate(); err != nil {
        return nil, fmt.Errorf("failed to select onchain account wallets: %w", err)
    }

    query, args := p.BuildQuery("SELECT * FROM " + s.tableName)

    // Execute.
    rows, err := s.executor.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to select onchain account wallets: %w", err)
    }
    defer rows.Close()

    var result []*AccountWallet
    for rows.Next() {
        row := &AccountWallet{}
        if err := rows.Scan(
            &row.ID,
            &row.AccountID,
            &row.Status,
            &row.Chain,
            &row.Network,
            &row.Address,
            &row.ProviderKind,
            &row.KeyVersion,
            &row.EncryptedPrivateKey,
            &row.CreatedAt,
            &row.UpdatedAt,
        ); err != nil {
            return nil, fmt.Errorf("failed to select onchain account wallets: %w", err)
        }

        result = append(result, row)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to select onchain account wallets: %w", err)
    }

    return result, nil
}


//
// Select onchain account wallet by ID.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *AccountWalletStore) SelectByID(id uint64) (*AccountWallet, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select onchain account wallet by id: missing required parameter: wallet_store=null")
    }
    if s.executor == nil {
        return nil, fmt.Errorf("failed to select onchain account wallet by id: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select onchain account wallet by id: missing required parameter: table_name=%q", "empty")
    }
    if id == 0 {
        return nil, fmt.Errorf("failed to select onchain account wallet by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColID)

    row := s.executor.QueryRow(query, id)

    result := &AccountWallet{}
    err := row.Scan(
        &result.ID,
        &result.AccountID,
        &result.Status,
        &result.Chain,
        &result.Network,
        &result.Address,
        &result.ProviderKind,
        &result.KeyVersion,
        &result.EncryptedPrivateKey,
        &result.CreatedAt,
        &result.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to select onchain account wallet by id: %w", err)
    }

    return result, nil
}


//
// Count onchain account wallets.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *AccountWalletStore) Count(p *AccountWalletSelectParams) (int64, error) {
    if s == nil {
        return 0, fmt.Errorf("failed to count onchain account wallets: missing required parameter: wallet_store=null")
    }
    if s.executor == nil {
        return 0, fmt.Errorf("failed to count onchain account wallets: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return 0, fmt.Errorf("failed to count onchain account wallets: missing required parameter: table_name=%q", "empty")
    }
    if err := p.Validate(); err != nil {
        return 0, fmt.Errorf("failed to count onchain account wallets: %w", err)
    }

    query, args := p.BuildQuery("SELECT COUNT(*) FROM " + s.tableName)

    var result int64
    if err := s.executor.QueryRow(query, args...).Scan(&result); err != nil {
        return 0, fmt.Errorf("failed to count onchain account wallets: %w", err)
    }

    return result, nil
}


//
// Delete onchain account wallet by ID.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *AccountWalletStore) DeleteByID(id uint64) error {
    if s == nil {
        return fmt.Errorf("failed to delete onchain account wallet by id: missing required parameter: wallet_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to delete onchain account wallet by id: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to delete onchain account wallet by id: missing required parameter: table_name=%q", "empty")
    }
    if id == 0 {
        return fmt.Errorf("failed to delete onchain account wallet by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", s.tableName, ColID)

    if _, err := s.executor.Exec(query, id); err != nil {
        return fmt.Errorf("failed to delete onchain account wallet by id: %w", err)
    }

    return nil
}


//
// Update onchain account wallet by ID.
//
// Version:
//   - 2026-05-25: Added.
//

func (s *AccountWalletStore) Update(p *AccountWalletUpdateParams) error {
    if s == nil {
        return fmt.Errorf("failed to update onchain account wallet: missing required parameter: wallet_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to update onchain account wallet: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update onchain account wallet: missing required parameter: table_name=%q", "empty")
    }
    if p == nil {
        return fmt.Errorf("failed to update onchain account wallet: missing required parameter: wallet_update_params=null")
    }
    if p.ID == 0 {
        return fmt.Errorf("failed to update onchain account wallet: invalid parameter: id=0")
    }

    assignments := make([]string, 0, 1)
    args := make([]any, 0, 2)

    if p.Status != nil {
        if err := ValidateAccountWalletStatus(*p.Status); err != nil {
            return fmt.Errorf("failed to update onchain account wallet: %w", err)
        }
        assignments = append(assignments, ColStatus+" = ?")
        args = append(args, *p.Status)
    }

    if len(assignments) == 0 {
        return fmt.Errorf("failed to update onchain account wallet: invalid parameter: assignments=%q", "empty")
    }

    args = append(args, p.ID)

    query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?;", s.tableName, strings.Join(assignments, ", "), ColID)

    if _, err := s.executor.Exec(query, args...); err != nil {
        return fmt.Errorf("failed to update onchain account wallet: %w", err)
    }

    return nil
}


//
// Build query.
//
// Version:
//   - 2025-05-16: Added.
//
func (p *AccountWalletSelectParams) BuildQuery(selectFromClause string) (string, []any) {
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
    if p.AccountID != nil {
        conditions = append(conditions, ColAccountID + " = ?")
        args = append(args, *p.AccountID)
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
func (p *AccountWalletSelectParams) Validate() error {
    // Guard.
    if p == nil {
        return fmt.Errorf("missing required parameter: wallet_select_params=null")
    }

    if p.ID != nil {
        if err := ValidateAccountWalletID(*p.ID); err != nil {
            return err
        }
    }
    if p.AccountID != nil {
        if err := ValidateAccountWalletAccountID(*p.AccountID); err != nil {
            return err
        }
    }
    if p.Status != nil {
        if err := ValidateAccountWalletStatus(*p.Status); err != nil {
            return err
        }
    }
    if p.Chain != nil {
        if err := ValidateAccountWalletChain(*p.Chain); err != nil {
            return err
        }
    }
    if p.Network != nil {
        if err := ValidateAccountWalletNetwork(*p.Network); err != nil {
            return err
        }
    }
    if p.Address != nil {
        if err := ValidateAccountWalletAddress(*p.Address); err != nil {
            return err
        }
    }

    if p.OrderBy != "" {
        switch p.OrderBy {
        case ColID,
            ColAccountID,
            ColStatus,
            ColChain,
            ColNetwork,
            ColAddress,
            ColProviderKind,
            ColKeyVersion,
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

