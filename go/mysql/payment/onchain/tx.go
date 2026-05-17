//
// tx.go
//
package onchain

import (
    "database/sql"
    "fmt"
    "strings"
    "time"
    "unicode/utf8"

    _ "github.com/go-sql-driver/mysql"

    "github.com/k4k3ru-hub/db/go/mysql/helper"
)

const (
    DefaultTxTableName = "payment_onchain_txs"
)

var (
    txIDCounter = &helper.IdCounter{}
)

type Tx struct {
    ID            uint64     `json:"id,string"`
    RequestID     *uint64    `json:"requestId,string,omitempty"`
    AccountID     uint64     `json:"accountId,string"`
    Status        TxStatus   `json:"status"`
    Chain         Chain      `json:"chain"`
    Network       Network    `json:"network"`
    Asset         Asset      `json:"asset"`
    TxHash        string     `json:"txHash"`
    FromAddress   string     `json:"fromAddress,omitempty"`
    ToAddress     string     `json:"toAddress"`
    Amount        string     `json:"amount"`
    Confirmations uint64     `json:"confirmations,string"`
    DetectedAt    time.Time  `json:"detectedAt"`
    ConfirmedAt   *time.Time `json:"confirmedAt,omitempty"`
    CreatedAt     time.Time  `json:"createdAt,omitempty"`
    UpdatedAt     time.Time  `json:"updatedAt,omitempty"`
}

type TxStore struct {
    executor         helper.Executor
    tableName        string
    requestTableName string
}

type TxInsertParams struct {
    ID            uint64     `json:"id,string"`
    RequestID     *uint64    `json:"requestId,string,omitempty"`
    AccountID     uint64     `json:"accountId,string"`
    Status        TxStatus   `json:"status"`
    Chain         Chain      `json:"chain"`
    Network       Network    `json:"network"`
    Asset         Asset      `json:"asset"`
    TxHash        string     `json:"txHash"`
    FromAddress   string     `json:"fromAddress,omitempty"`
    ToAddress     string     `json:"toAddress"`
    Amount        string     `json:"amount"`
    Confirmations uint64     `json:"confirmations,string"`
    DetectedAt    time.Time  `json:"detectedAt"`
    ConfirmedAt   *time.Time `json:"confirmedAt,omitempty"`
    CreatedAt     time.Time  `json:"createdAt"`
    UpdatedAt     time.Time  `json:"updatedAt"`
    Ignore        bool       `json:"ignore"`
}

type TxSelectParams struct {
    RequestID   *uint64   `json:"requestId,string,omitempty"`
    AccountID   *uint64   `json:"accountId,string,omitempty"`
    Status      *TxStatus `json:"status,omitempty"`
    Chain       *Chain    `json:"chain,omitempty"`
    Network     *Network  `json:"network,omitempty"`
    Asset       *Asset    `json:"asset,omitempty"`
    TxHash      *string   `json:"txHash,omitempty"`
    ToAddress   *string   `json:"toAddress,omitempty"`
    OrderBy     string    `json:"orderBy"`
    OrderByDesc bool      `json:"orderByDesc"`
    Limit       int       `json:"limit"`
    Offset      int       `json:"offset"`
}

type TxUpdateParams struct {
    ID            uint64     `json:"id,string"`
    RequestID     *uint64    `json:"requestId,string,omitempty"`
    Status        *TxStatus  `json:"status,omitempty"`
    Confirmations *uint64    `json:"confirmations,string,omitempty"`
    ConfirmedAt   *time.Time `json:"confirmedAt,omitempty"`
}


//
// Generate new payment onchain tx ID.
//
// Version:
//   - 2026-05-16: Added.
//
func GenerateTxID() uint64 {
    return txIDCounter.GenerateID()
}


//
// Create new payment onchain tx store.
//
// Version:
//   - 2026-05-16: Added.
//
func NewTxStore(executor helper.Executor, tableName, requestTableName string) (*TxStore, error) {
    if executor == nil {
        return nil, fmt.Errorf("failed to create payment onchain tx store: missing required parameter: executor=null")
    }
    if tableName == "" {
        return nil, fmt.Errorf("failed to create payment onchain tx store: missing required parameter: table_name=%q", "empty")
    }
    if requestTableName == "" {
        return nil, fmt.Errorf("failed to create payment onchain tx store: missing required parameter: request_table_name=%q", "empty")
    }

    return &TxStore{
        executor:         executor,
        tableName:        tableName,
        requestTableName: requestTableName,
    }, nil
}


//
// Validate payment onchain tx ID.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateTxID(id uint64) error {
    if id == 0 {
        return fmt.Errorf("invalid parameter: id=0")
    }
    return nil
}


//
// Validate payment onchain tx ID.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *Tx) ValidateID() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_tx=null")
    }
    return ValidateTxID(t.ID)
}


//
// Validate payment onchain tx request ID.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateTxRequestID(requestID *uint64) error {
    if requestID == nil {
        return nil
    }
    if *requestID == 0 {
        return fmt.Errorf("invalid parameter: request_id=0")
    }
    return nil
}


//
// Validate payment onchain tx request ID.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *Tx) ValidateRequestID() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_tx=null")
    }
    return ValidateTxRequestID(t.RequestID)
}


//
// Validate payment onchain tx account ID.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateTxAccountID(accountID uint64) error {
    if accountID == 0 {
        return fmt.Errorf("invalid parameter: account_id=0")
    }
    return nil
}


//
// Validate payment onchain tx account ID.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *Tx) ValidateAccountID() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_tx=null")
    }
    return ValidateTxAccountID(t.AccountID)
}


//
// Validate payment onchain tx status.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateTxStatus(s TxStatus) error {
    if err := s.Validate(); err != nil {
        return err
    }
    return nil
}


//
// Validate payment onchain tx status.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *Tx) ValidateStatus() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_tx=null")
    }
    return ValidateTxStatus(t.Status)
}


//
// Validate payment onchain tx chain.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateTxChain(c Chain) error {
    if err := c.Validate(); err != nil {
        return err
	}
    return nil
}


//
// Validate payment onchain tx chain.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *Tx) ValidateChain() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_tx=null")
    }
    return ValidateTxChain(t.Chain)
}


//
// Validate payment onchain tx network.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateTxNetwork(n Network) error {
    if err := n.Validate(); err != nil {
        return err
    }
    return nil
}


//
// Validate payment onchain tx network.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *Tx) ValidateNetwork() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_tx=null")
    }
    return ValidateTxNetwork(t.Network)
}


//
// Validate payment onchain tx asset.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateTxAsset(a Asset) error {
    if err := a.Validate(); err != nil {
        return err
    }
    return nil
}


//
// Validate payment onchain tx asset.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *Tx) ValidateAsset() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_tx=null")
    }
    return ValidateTxAsset(t.Asset)
}


//
// Validate payment onchain tx tx hash.
//  
// Version:
//   - 2026-05-16: Added.
//
func ValidateTxTxHash(txHash string) error {
    if txHash == "" {
        return fmt.Errorf("invalid parameter: tx_hash=%q", "empty")
    }
    if utf8.RuneCountInString(txHash) > 255 {
        return fmt.Errorf("invalid parameter: tx_hash=%q", helper.TruncateRunes(txHash, 255))
    }
    return nil
}


//
// Validate payment onchain tx tx hash.
//  
// Version:
//   - 2026-05-16: Added.
//
func (t *Tx) ValidateTxHash() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_tx=null")
    }
    return ValidateTxTxHash(t.TxHash)
}


//
// Validate payment onchain tx from address.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateTxFromAddress(fromAddress string) error {
    if fromAddress == "" {
        return nil
    }
    if utf8.RuneCountInString(fromAddress) > 255 {
        return fmt.Errorf("invalid parameter: from_address=%q", helper.TruncateRunes(fromAddress, 255))
    }
    return nil
}


//
// Validate payment onchain tx from address.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *Tx) ValidateFromAddress() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_tx=null")
    }
    return ValidateTxFromAddress(t.FromAddress)
}


//
// Validate payment onchain tx to address.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateTxToAddress(toAddress string) error {
    if toAddress == "" {
        return fmt.Errorf("invalid parameter: to_address=empty")
    }
    if utf8.RuneCountInString(toAddress) > 255 {
        return fmt.Errorf("invalid parameter: to_address=%q", helper.TruncateRunes(toAddress, 255))
    }
    return nil
}


//
// Validate payment onchain tx to address.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *Tx) ValidateToAddress() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_tx=null")
    }
    return ValidateTxToAddress(t.ToAddress)
}


//
// Validate payment onchain tx amount.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateTxAmount(amount string) error {
    if amount == "" {
        return fmt.Errorf("invalid parameter: amount=%q", "empty")
    }
    if utf8.RuneCountInString(amount) > 78 {
        return fmt.Errorf("invalid parameter: amount=%q", helper.TruncateRunes(amount, 78))
    }
    return nil
}


//
// Validate payment onchain tx amount.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *Tx) ValidateAmount() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_tx=null")
    }
    return ValidateTxAmount(t.Amount)
}


//
// Validate payment onchain tx detected at.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateTxDetectedAt(detectedAt time.Time) error {
    if detectedAt.IsZero() {
        return fmt.Errorf("invalid parameter: detected_at=%q", "empty")
    }
    return nil
}


//
// Validate payment onchain tx detected at.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *Tx) ValidateDetectedAt() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_tx=null")
    }
    return ValidateTxDetectedAt(t.DetectedAt)
}


//
// Validate payment onchain tx confirmed at.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateTxConfirmedAt(confirmedAt *time.Time) error {
    if confirmedAt == nil {
        return nil
    }
    if confirmedAt.IsZero() {
        return fmt.Errorf("invalid parameter: confirmed_at=%q", "empty")
    }
    return nil
}


//
// Validate payment onchain tx confirmed at.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *Tx) ValidateConfirmedAt() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_tx=null")
    }
    return ValidateTxConfirmedAt(t.ConfirmedAt)
}


//
// Create payment onchain txs table.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *TxStore) CreateTable() error {
    if s == nil {
        return fmt.Errorf("failed to create payment onchain txs table: missing required parameter: tx_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to create payment onchain txs table: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to create payment onchain txs table: missing required parameter: table_name=%q", "empty")
    }
    if s.requestTableName == "" {
        return fmt.Errorf("failed to create payment onchain txs table: missing required parameter: request_table_name=%q", "empty")
    }

    query := fmt.Sprintf(
        `CREATE TABLE IF NOT EXISTS %s (
            %s BIGINT UNSIGNED NOT NULL COMMENT 'ID',
            %s BIGINT UNSIGNED NULL COMMENT 'Request ID',
            %s BIGINT UNSIGNED NOT NULL COMMENT 'Account ID',
            %s TINYINT UNSIGNED NOT NULL COMMENT 'Status',
            %s VARCHAR(64) NOT NULL COMMENT 'Chain',
            %s VARCHAR(64) NOT NULL COMMENT 'Network',
            %s VARCHAR(64) NOT NULL COMMENT 'Asset',
            %s VARCHAR(255) NOT NULL COMMENT 'Transaction hash',
            %s VARCHAR(255) NULL COMMENT 'Source address if available',
            %s VARCHAR(255) NOT NULL COMMENT 'Monitored recipient address',
            %s VARCHAR(78) NOT NULL COMMENT 'Amount',
            %s BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Confirmations',
            %s DATETIME NOT NULL COMMENT 'First detected at',
            %s DATETIME NULL COMMENT 'Reached required confirmations at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
            PRIMARY KEY (%s),
            UNIQUE KEY uk_payment_onchain_txs_chain_network_tx_hash (%s, %s, %s),
            KEY idx_payment_onchain_txs_request_id (%s),
            KEY idx_payment_onchain_txs_account_id (%s),
            KEY idx_payment_onchain_txs_status (%s),
            KEY idx_payment_onchain_txs_to_address (%s),
            KEY idx_payment_onchain_txs_chain_network_asset (%s, %s, %s),
            CONSTRAINT fk_payment_onchain_txs_request_id FOREIGN KEY (%s) REFERENCES %s (%s) ON DELETE SET NULL ON UPDATE CASCADE;
        `,
        s.tableName,
        ColID,
        ColRequestID,
        ColAccountID,
        ColStatus,
        ColChain,
        ColNetwork,
        ColAsset,
        ColTxHash,
        ColFromAddress,
        ColToAddress,
        ColAmount,
        ColConfirmations,
        ColDetectedAt,
        ColConfirmedAt,
        ColCreatedAt,
        ColUpdatedAt,
        ColID,
        ColChain, ColNetwork, ColTxHash,
        ColRequestID,
        ColAccountID,
        ColStatus,
        ColToAddress,
        ColChain, ColNetwork, ColAsset,
        ColRequestID, s.requestTableName, ColID,
    )

    if _, err := s.executor.Exec(query); err != nil {
        return fmt.Errorf("failed to create payment onchain txs table: %w", err)
    }

    return nil
}


//
// Insert payment onchain tx.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *TxStore) Insert(p *TxInsertParams) error {
    if s == nil {
        return fmt.Errorf("failed to insert payment onchain tx: missing required parameter: tx_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to insert payment onchain tx: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to insert payment onchain tx: missing required parameter: table_name=%q", "empty")
    }
    if p == nil {
        return fmt.Errorf("failed to insert payment onchain tx: missing required parameter: tx_insert_params=null")
    }
    if err := ValidateTxAccountID(p.AccountID); err != nil {
        return fmt.Errorf("failed to insert payment onchain tx: %w", err)
    }
    if err := ValidateTxRequestID(p.RequestID); err != nil {
        return fmt.Errorf("failed to insert payment onchain tx: %w", err)
    }
    if err := ValidateTxStatus(p.Status); err != nil {
        return fmt.Errorf("failed to insert payment onchain tx: %w", err)
    }
    if err := ValidateTxChain(p.Chain); err != nil {
        return fmt.Errorf("failed to insert payment onchain tx: %w", err)
    }
    if err := ValidateTxNetwork(p.Network); err != nil {
        return fmt.Errorf("failed to insert payment onchain tx: %w", err)
    }
    if err := ValidateTxAsset(p.Asset); err != nil {
        return fmt.Errorf("failed to insert payment onchain tx: %w", err)
    }
    if err := ValidateTxTxHash(p.TxHash); err != nil {
        return fmt.Errorf("failed to insert payment onchain tx: %w", err)
    }
    if err := ValidateTxFromAddress(p.FromAddress); err != nil {
        return fmt.Errorf("failed to insert payment onchain tx: %w", err)
    }
    if err := ValidateTxToAddress(p.ToAddress); err != nil {
        return fmt.Errorf("failed to insert payment onchain tx: %w", err)
    }
    if err := ValidateTxAmount(p.Amount); err != nil {
        return fmt.Errorf("failed to insert payment onchain tx: %w", err)
    }
    if err := ValidateTxDetectedAt(p.DetectedAt); err != nil {
        return fmt.Errorf("failed to insert payment onchain tx: %w", err)
    }
    if err := ValidateTxConfirmedAt(p.ConfirmedAt); err != nil {
        return fmt.Errorf("failed to insert payment onchain tx: %w", err)
    }

    if p.ID == 0 {
        p.ID = GenerateTxID()
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
        "%s INTO %s (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
        queryPrefix,
        s.tableName,
        ColID,
        ColRequestID,
        ColAccountID,
        ColStatus,
        ColChain,
        ColNetwork,
        ColAsset,
        ColTxHash,
        ColFromAddress,
        ColToAddress,
        ColAmount,
        ColConfirmations,
        ColDetectedAt,
        ColConfirmedAt,
        ColCreatedAt,
        ColUpdatedAt,
    )

    if _, err := s.executor.Exec(
        query,
        p.ID,
        p.RequestID,
        p.AccountID,
        p.Status,
        p.Chain,
        p.Network,
        p.Asset,
        p.TxHash,
        p.FromAddress,
        p.ToAddress,
        p.Amount,
        p.Confirmations,
        p.DetectedAt,
        p.ConfirmedAt,
        p.CreatedAt,
        p.UpdatedAt,
    ); err != nil {
        return fmt.Errorf("failed to insert payment onchain tx: %w", err)
    }

    return nil
}


//
// Select payment onchain txs.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *TxStore) Select(p *TxSelectParams) ([]*Tx, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select payment onchain txs: missing required parameter: tx_store=null")
    }
    if s.executor == nil {
        return nil, fmt.Errorf("failed to select payment onchain txs: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select payment onchain txs: missing required parameter: table_name=%q", "empty")
    }
    if err := p.Validate(); err != nil {
        return nil, fmt.Errorf("failed to select payment onchain request txs: %w", err)
    }

    query, args := p.BuildQuery("SELECT * FROM " + s.tableName)

    // Execute.
    rows, err := s.executor.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to select payment onchain txs: %w", err)
    }
    defer rows.Close()

    var result []*Tx
    for rows.Next() {
        row := &Tx{}
        err := rows.Scan(
            &row.ID,
            &row.RequestID,
            &row.AccountID,
            &row.Status,
            &row.Chain,
            &row.Network,
            &row.Asset,
            &row.TxHash,
            &row.FromAddress,
            &row.ToAddress,
            &row.Amount,
            &row.Confirmations,
            &row.DetectedAt,
            &row.ConfirmedAt,
            &row.CreatedAt,
            &row.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to select payment onchain txs: %w", err)
        }

        result = append(result, row)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to select payment onchain txs: %w", err)
    }

    return result, nil
}


//
// Select payment onchain tx by ID.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *TxStore) SelectByID(id uint64) (*Tx, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select payment onchain tx by id: missing required parameter: tx_store=null")
    }
    if s.executor == nil {
        return nil, fmt.Errorf("failed to select payment onchain tx by id: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select payment onchain tx by id: missing required parameter: table_name=%q", "empty")
    }
    if id == 0 {
        return nil, fmt.Errorf("failed to select payment onchain tx by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColID)

    result := &Tx{}
    err := s.executor.QueryRow(query, id).Scan(
        &result.ID,
        &result.RequestID,
        &result.AccountID,
        &result.Status,
        &result.Chain,
        &result.Network,
        &result.Asset,
        &result.TxHash,
        &result.FromAddress,
        &result.ToAddress,
        &result.Amount,
        &result.Confirmations,
        &result.DetectedAt,
        &result.ConfirmedAt,
        &result.CreatedAt,
        &result.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to select payment onchain tx by id: %w", err)
    }

    return result, nil
}


//
// Count payment onchain txs.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *TxStore) Count(p *TxSelectParams) (int64, error) {
    if s == nil {
        return 0, fmt.Errorf("failed to count payment onchain txs: missing required parameter: tx_store=null")
    }
    if s.executor == nil {
        return 0, fmt.Errorf("failed to count payment onchain txs: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return 0, fmt.Errorf("failed to count payment onchain txs: missing required parameter: table_name=%q", "empty")
    }
    if err := p.Validate(); err != nil {
        return 0, fmt.Errorf("failed to count payment onchain txs: %w", err)
    }

    query, args := p.BuildQuery("SELECT COUNT(*) FROM " + s.tableName)

    // Execute.
    var result int64
    if err := s.executor.QueryRow(query, args...).Scan(&result); err != nil {
        return 0, fmt.Errorf("failed to count payment onchain txs: %w", err)
    }

    return result, nil
}


//
// Delete payment onchain tx by ID.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *TxStore) DeleteByID(id uint64) error {
    if s == nil {
        return fmt.Errorf("failed to delete payment onchain tx by id: missing required parameter: tx_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to delete payment onchain tx by id: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to delete payment onchain tx by id: missing required parameter: table_name=%q", "empty")
    }
    if id == 0 {
        return fmt.Errorf("failed to delete payment onchain tx by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", s.tableName, ColID)

    if _, err := s.executor.Exec(query, id); err != nil {
        return fmt.Errorf("failed to delete payment onchain tx by id: %w", err)
    }

    return nil
}



func (s *TxStore) Update(option *TxUpdateParams) error {
    if s == nil {
        return fmt.Errorf("failed to update payment onchain tx: missing required parameter: tx_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to update payment onchain tx: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update payment onchain tx: missing required parameter: table_name=%q", "empty")
    }
    if option == nil {
        return fmt.Errorf("failed to update payment onchain tx: missing required parameter: option=null")
    }
    if option.ID == 0 {
        return fmt.Errorf("failed to update payment onchain tx: invalid parameter: id=0")
    }

    assignments := make([]string, 0, 4)
    args := make([]any, 0, 4)

    if option.RequestID != nil {
        if *option.RequestID == 0 {
            return fmt.Errorf("failed to update payment onchain tx: invalid parameter: request_id=0")
        }
        assignments = append(assignments, ColRequestID+" = ?")
        args = append(args, *option.RequestID)
    }

    if option.Status != nil {
        t := Tx{Status: *option.Status}
        if err := t.ValidateStatus(); err != nil {
            return fmt.Errorf("failed to update payment onchain tx: %w", err)
        }
        assignments = append(assignments, ColStatus+" = ?")
        args = append(args, *option.Status)
    }

    if option.Confirmations != nil {
        assignments = append(assignments, ColConfirmations+" = ?")
        args = append(args, *option.Confirmations)
    }

    if option.ConfirmedAt != nil {
        t := Tx{ConfirmedAt: option.ConfirmedAt}
        if err := t.ValidateConfirmedAt(); err != nil {
            return fmt.Errorf("failed to update payment onchain tx: %w", err)
        }
        assignments = append(assignments, ColConfirmedAt+" = ?")
        args = append(args, *option.ConfirmedAt)
    }

    if len(assignments) == 0 {
        return fmt.Errorf("failed to update payment onchain tx: invalid parameter: assignments=empty")
    }

    args = append(args, option.ID)

    query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?;", s.tableName, strings.Join(assignments, ", "), ColID)

    if _, err := s.executor.Exec(query, args...); err != nil {
        return fmt.Errorf("failed to update payment onchain tx: %w", err)
    }

    return nil
}

func (s *TxStore) UpdateConfirmed(id uint64, confirmations uint64) error {
    if id == 0 {
        return fmt.Errorf("failed to update payment onchain tx confirmed: invalid parameter: id=0")
    }

    status := TxStatusConfirmed
    now := time.Now()

    return s.Update(&TxUpdateParams{
        ID:            id,
        Status:        &status,
        Confirmations: &confirmations,
        ConfirmedAt:   &now,
    })
}


//
// Build query.
//
// Version:
//   - 2025-05-16: Added.
//
func (p *TxSelectParams) BuildQuery(selectFromClause string) (string, []any) {
    // Guard.
    if p == nil {
        return selectFromClause, nil
    }

    var query strings.Builder
    query.WriteString(selectFromClause)

    conditions := make([]string, 0, 8)
    args := make([]any, 0, 10)

    if p.RequestID != nil {
        conditions = append(conditions, ColRequestID + " = ?")
        args = append(args, *p.RequestID)
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
    if p.Asset != nil {
        conditions = append(conditions, ColAsset + " = ?")
        args = append(args, *p.Asset)
    }
    if p.TxHash != nil {
        conditions = append(conditions, ColTxHash + " = ?")
        args = append(args, *p.TxHash)
    }
    if p.ToAddress != nil {
        conditions = append(conditions, ColToAddress + " = ?")
        args = append(args, *p.ToAddress)
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
func (p *TxSelectParams) Validate() error {
    // Guard.
    if p == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_tx_select_params=null")
    }

    if p.RequestID != nil {
        if err := ValidateTxRequestID(p.RequestID); err != nil {
            return err
        }
    }
    if p.AccountID != nil {
        if err := ValidateTxAccountID(*p.AccountID); err != nil {
            return err
        }
    }
    if p.Status != nil {
        if err := ValidateTxStatus(*p.Status); err != nil {
            return err
        }
    }
    if p.Chain != nil {
        if err := ValidateTxChain(*p.Chain); err != nil {
            return err
        }
    }
    if p.Network != nil {
        if err := ValidateTxNetwork(*p.Network); err != nil {
            return err
        }
    }
    if p.Asset != nil {
        if err := ValidateTxAsset(*p.Asset); err != nil {
            return err
        }
    }
    if p.TxHash != nil {
        if err := ValidateTxTxHash(*p.TxHash); err != nil {
            return err
        }
    }
    if p.ToAddress != nil {
        if err := ValidateTxToAddress(*p.ToAddress); err != nil {
            return err
        }
    }

    if p.OrderBy != "" {
        switch p.OrderBy {
        case ColID,
            ColRequestID,
            ColAccountID,
            ColStatus,
            ColChain,
            ColNetwork,
            ColAsset,
            ColTxHash,
            ColFromAddress,
            ColToAddress,
            ColAmount,
            ColConfirmations,
            ColDetectedAt,
            ColConfirmedAt,
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


