//
// wallet.go
//
package onchain

import (
    "database/sql"
    "fmt"
    "strings"
    "time"

    _ "github.com/go-sql-driver/mysql"

    "github.com/k4k3ru-hub/db/go/mysql/account"
    "github.com/k4k3ru-hub/db/go/mysql/helper"
)

const (
    DefaultWalletTableName = "account_onchain_wallets"
)

var (
    walletIDCounter = &helper.IdCounter{}
)


type Wallet struct {
    ID                  uint64       `json:"id,string"`
    AccountID           uint64       `json:"accountId,string"`
    Status              WalletStatus `json:"status"`
    Chain               Chain        `json:"chain"`
    Network             Network      `json:"network"`
    Address             string       `json:"address"`
    ProviderKind        string       `json:"providerKind"`
    KeyVersion          string       `json:"keyVersion"`
    EncryptedPrivateKey string       `json:"encryptedPrivateKey"`
    CreatedAt           time.Time    `json:"createdAt,omitempty"`
    UpdatedAt           time.Time    `json:"updatedAt,omitempty"`
}

type WalletStore struct {
    executor         helper.Executor
    tableName        string
    accountTableName string
}

type WalletInsertParams struct {
    ID                  uint64       `json:"id,string"`
    AccountID           uint64       `json:"accountId,string"`
    Status              WalletStatus `json:"status"`
    Chain               Chain        `json:"chain"`
    Network             Network      `json:"network"`
    Address             string       `json:"address"`
    ProviderKind        string       `json:"providerKind"`
    KeyVersion          string       `json:"keyVersion"`
    EncryptedPrivateKey string       `json:"encryptedPrivateKey"`
    CreatedAt           time.Time    `json:"createdAt,omitempty"`
    UpdatedAt           time.Time    `json:"updatedAt,omitempty"`
    Ignore              bool         `json:"ignore"`
}

type WalletSelectParams struct {
    ID          *uint64       `json:"id,string,omitempty"`
    AccountID   *uint64       `json:"accountId,string,omitempty"`
    Status      *WalletStatus `json:"status,omitempty"`
    Chain       *Chain        `json:"chain,omitempty"`
    Network     *Network      `json:"network,omitempty"`
    Address     *string       `json:"address,omitempty"`
    OrderBy     string        `json:"orderBy"`
    OrderByDesc bool          `json:"orderByDesc"`
    Limit       int           `json:"limit"`
    Offset      int           `json:"offset"`
}

type WalletUpdateParams struct {
    ID     uint64        `json:"id,string"`
    Status *WalletStatus `json:"status,omitempty"`
}


//
// Generate new account onchain wallet ID.
//
// Version:
//   - 2026-05-25: Added.
//
func GenerateWalletID() uint64 {
    return walletIDCounter.GenerateID()
}


//
// Create new account onchain wallet store.
//
// Version:
//   - 2026-05-25: Added.
//
func NewWalletStore(executor helper.Executor, tableName, accountTableName string) (*WalletStore, error) {
    // Guard.
    if executor == nil {
        return nil, fmt.Errorf("failed to create account onchain wallet store: missing required parameter: executor=null")
    }
    if tableName == "" {
        return nil, fmt.Errorf("failed to create account onchain wallet store: missing required parameter: table_name=%q", "empty")
    }
    if accountTableName == "" {
        return nil, fmt.Errorf("failed to create account onchain wallet store: missing required parameter: account_table_name=%q", "empty")
    }

    return &WalletStore{
        executor:         executor,
        tableName:        tableName,
        accountTableName: accountTableName,
    }, nil
}


//
// Validate account onchain wallet ID.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateWalletID(id uint64) error {
    if id == 0 {
        return fmt.Errorf("invalid parameter: id=0")
    }
    return nil
}


//
// Validate account onchain wallet ID.
//
// Version:
//   - 2026-05-25: Added.
//
func (w *Wallet) ValidateID() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateWalletID(w.ID)
}


//
// Validate account onchain wallet account ID.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateWalletAccountID(accountID uint64) error {
    if accountID == 0 {
        return fmt.Errorf("invalid parameter: account_id=0")
    }
    return nil
}


//
// Validate account onchain wallet account ID.
//
// Version:
//   - 2026-05-25: Added.
//
func (w *Wallet) ValidateAccountID() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateWalletID(w.AccountID)
}


//
// Validate account onchain wallet status.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateWalletStatus(status WalletStatus) error {
    if !status.IsValid() {
        return fmt.Errorf("invalid parameter: status=%d", status)
    }
    return nil
}


//
// Validate account onchain wallet status.
//
// Version:
//   - 2026-05-25: Added.
//
func (w *Wallet) ValidateStatus() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateWalletStatus(w.Status)
}


//
// Validate account onchain wallet chain.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateWalletChain(c Chain) error {
    if err := c.Validate(); err != nil {
        return err
    }
    return nil
}


//
// Validate account onchain wallet chain.
//
// Version:
//   - 2026-05-25: Added.
//
func (w *Wallet) ValidateChain() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateWalletChain(w.Chain)
}


//
// Validate account onchain wallet network.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateWalletNetwork(n Network) error {
    if err := n.Validate(); err != nil {
        return err
    }
    return nil
}


//
// Validate account onchain wallet network.
//
// Version:
//   - 2026-05-25: Added.
//
func (w *Wallet) ValidateNetwork() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateWalletNetwork(w.Network)
}


//
// Validate account onchain wallet address.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateWalletAddress(address string) error {
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
// Validate account onchain wallet address.
//
// Version:
//   - 2026-05-25: Added.
//
func (w *Wallet) ValidateAddress() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateWalletAddress(w.Address)
}


//
// Validate account onchain wallet provider kind.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateWalletProviderKind(providerKind string) error {
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
// Validate account onchain wallet provider kind.
//
// Version:
//   - 2026-05-25: Added.
//
func (w *Wallet) ValidateProviderKind() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateWalletProviderKind(w.ProviderKind)
}


//
// Validate account onchain wallet key version.
//
// Version:
//   - 2026-05-25: Added.
//
func ValidateWalletKeyVersion(keyVersion string) error {
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
// Validate account onchain wallet key version.
//  
// Version:
//   - 2026-05-25: Added.
//
func (w *Wallet) ValidateKeyVersion() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateWalletKeyVersion(w.KeyVersion)
}


//
// Validate account onchain wallet encrypted private key.
//
// Version:
//   - 2026-05-25: Added.
//  
func ValidateWalletEncryptedPrivateKey(encryptedPrivateKey string) error {
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
// Validate account onchain wallet encrypted private key.
//  
// Version:
//   - 2026-05-25: Added.
//
func (w *Wallet) ValidateEncryptedPrivateKey() error {
    if w == nil {
        return fmt.Errorf("missing required parameter: wallet=null")
    }
    return ValidateWalletEncryptedPrivateKey(w.EncryptedPrivateKey)
}


//
// Create account onchain wallets table.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *WalletStore) CreateTable() error {
    if s == nil {
        return fmt.Errorf("failed to create account onchain wallets table: missing required parameter: wallet_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to create account onchain wallets table: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to create account onchain wallets table: missing required parameter: table_name=%q", "empty")
    }
    if s.accountTableName == "" {
        return fmt.Errorf("failed to create account onchain wallets table: missing required parameter: account_table_name=%q", "empty")
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
            KEY idx_account_onchain_wallets_account_chain_network_status (%s, %s, %s, %s),
            CONSTRAINT fk_account_onchain_wallets_account_id FOREIGN KEY (%s) REFERENCES %s (%s) ON DELETE CASCADE ON UPDATE CASCADE
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
        ColAccountID, s.accountTableName, account.ColID,
    )

    if _, err := s.executor.Exec(query); err != nil {
        return fmt.Errorf("failed to create account onchain wallets table: %w", err)
    }

    return nil
}


//
// Insert account onchain wallet.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *WalletStore) Insert(p *WalletInsertParams) error {
    if s == nil {
        return fmt.Errorf("failed to insert account onchain wallet: missing required parameter: wallet_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to insert account onchain wallet: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to insert account onchain wallet: missing required parameter: table_name=%q", "empty")
    }
    if p == nil {
        return fmt.Errorf("failed to insert account onchain wallet: missing required parameter: wallet_insert_params=null")
    }
    if err := ValidateWalletAccountID(p.AccountID); err != nil {
        return fmt.Errorf("failed to insert account onchain wallet: %w", err)
    }
    if err := ValidateWalletStatus(p.Status); err != nil {
        return fmt.Errorf("failed to insert account onchain wallet: %w", err)
    }
    if err := ValidateWalletChain(p.Chain); err != nil {
        return fmt.Errorf("failed to insert account onchain wallet: %w", err)
    }
    if err := ValidateWalletNetwork(p.Network); err != nil {
        return fmt.Errorf("failed to insert account onchain wallet: %w", err)
    }
    if err := ValidateWalletAddress(p.Address); err != nil {
        return fmt.Errorf("failed to insert account onchain wallet: %w", err)
    }
    if err := ValidateWalletProviderKind(p.ProviderKind); err != nil {
        return fmt.Errorf("failed to insert account onchain wallet: %w", err)
    }
    if err := ValidateWalletKeyVersion(p.KeyVersion); err != nil {
        return fmt.Errorf("failed to insert account onchain wallet: %w", err)
    }
    if err := ValidateWalletEncryptedPrivateKey(p.EncryptedPrivateKey); err != nil {
        return fmt.Errorf("failed to insert account onchain wallet: %w", err)
    }

    if p.ID == 0 {
        p.ID = GenerateWalletID()
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
        return fmt.Errorf("failed to insert account onchain wallet: %w", err)
    }

    return nil
}


//
// Select account onchain wallets.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *WalletStore) Select(p *WalletSelectParams) ([]*Wallet, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select account onchain wallets: missing required parameter: wallet_store=null")
    }
    if s.executor == nil {
        return nil, fmt.Errorf("failed to select account onchain wallets: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account onchain wallets: missing required parameter: table_name=%q", "empty")
    }
    if err := p.Validate(); err != nil {
        return nil, fmt.Errorf("failed to select account onchain wallets: %w", err)
    }

    query, args := p.BuildQuery("SELECT * FROM " + s.tableName)

    // Execute.
    rows, err := s.executor.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to select account onchain wallets: %w", err)
    }
    defer rows.Close()

    var result []*Wallet
    for rows.Next() {
        row := &Wallet{}
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
            return nil, fmt.Errorf("failed to select account onchain wallets: %w", err)
        }

        result = append(result, row)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to select account onchain wallets: %w", err)
    }

    return result, nil
}


//
// Select account onchain wallet by ID.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *WalletStore) SelectByID(id uint64) (*Wallet, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select account onchain wallet by id: missing required parameter: wallet_store=null")
    }
    if s.executor == nil {
        return nil, fmt.Errorf("failed to select account onchain wallet by id: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account onchain wallet by id: missing required parameter: table_name=%q", "empty")
    }
    if id == 0 {
        return nil, fmt.Errorf("failed to select account onchain wallet by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColID)

    row := s.executor.QueryRow(query, id)

    result := &Wallet{}
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
        return nil, fmt.Errorf("failed to select account onchain wallet by id: %w", err)
    }

    return result, nil
}


//
// Count account onchain wallets.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *WalletStore) Count(p *WalletSelectParams) (int64, error) {
    if s == nil {
        return 0, fmt.Errorf("failed to count account onchain wallets: missing required parameter: wallet_store=null")
    }
    if s.executor == nil {
        return 0, fmt.Errorf("failed to count account onchain wallets: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return 0, fmt.Errorf("failed to count account onchain wallets: missing required parameter: table_name=%q", "empty")
    }
    if err := p.Validate(); err != nil {
        return 0, fmt.Errorf("failed to count account onchain wallets: %w", err)
    }

    query, args := p.BuildQuery("SELECT COUNT(*) FROM " + s.tableName)

    var result int64
    if err := s.executor.QueryRow(query, args...).Scan(&result); err != nil {
        return 0, fmt.Errorf("failed to count account onchain wallets: %w", err)
    }

    return result, nil
}


//
// Delete account onchain wallet by ID.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *WalletStore) DeleteByID(id uint64) error {
    if s == nil {
        return fmt.Errorf("failed to delete account onchain wallet by id: missing required parameter: wallet_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to delete account onchain wallet by id: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to delete account onchain wallet by id: missing required parameter: table_name=%q", "empty")
    }
    if id == 0 {
        return fmt.Errorf("failed to delete account onchain wallet by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", s.tableName, ColID)

    if _, err := s.executor.Exec(query, id); err != nil {
        return fmt.Errorf("failed to delete account onchain wallet by id: %w", err)
    }

    return nil
}


//
// Update account onchain wallet by ID.
//
// Version:
//   - 2026-05-25: Added.
//

func (s *WalletStore) Update(p *WalletUpdateParams) error {
    if s == nil {
        return fmt.Errorf("failed to update account onchain wallet: missing required parameter: wallet_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to update account onchain wallet: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update account onchain wallet: missing required parameter: table_name=%q", "empty")
    }
    if p == nil {
        return fmt.Errorf("failed to update account onchain wallet: missing required parameter: wallet_update_params=null")
    }
    if p.ID == 0 {
        return fmt.Errorf("failed to update account onchain wallet: invalid parameter: id=0")
    }

    assignments := make([]string, 0, 1)
    args := make([]any, 0, 2)

    if p.Status != nil {
        if err := ValidateWalletStatus(*p.Status); err != nil {
            return fmt.Errorf("failed to update account onchain wallet: %w", err)
        }
        assignments = append(assignments, ColStatus+" = ?")
        args = append(args, *p.Status)
    }

    if len(assignments) == 0 {
        return fmt.Errorf("failed to update account onchain wallet: invalid parameter: assignments=%q", "empty")
    }

    args = append(args, p.ID)

    query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?;", s.tableName, strings.Join(assignments, ", "), ColID)

    if _, err := s.executor.Exec(query, args...); err != nil {
        return fmt.Errorf("failed to update account onchain wallet: %w", err)
    }

    return nil
}


//
// Build query.
//
// Version:
//   - 2025-05-16: Added.
//
func (p *WalletSelectParams) BuildQuery(selectFromClause string) (string, []any) {
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
func (p *WalletSelectParams) Validate() error {
    // Guard.
    if p == nil {
        return fmt.Errorf("missing required parameter: wallet_select_params=null")
    }

    if p.ID != nil {
        if err := ValidateWalletID(*p.ID); err != nil {
            return err
        }
    }
    if p.AccountID != nil {
        if err := ValidateWalletAccountID(*p.AccountID); err != nil {
            return err
        }
    }
    if p.Status != nil {
        if err := ValidateWalletStatus(*p.Status); err != nil {
            return err
        }
    }
    if p.Chain != nil {
        if err := ValidateWalletChain(*p.Chain); err != nil {
            return err
        }
    }
    if p.Network != nil {
        if err := ValidateWalletNetwork(*p.Network); err != nil {
            return err
        }
    }
    if p.Address != nil {
        if err := ValidateWalletAddress(*p.Address); err != nil {
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

