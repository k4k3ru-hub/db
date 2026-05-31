//
// request_tx.go
//
package onchain

import (
    "database/sql"
    "fmt"
    "math/big"
    "strings"
    "time"
    "unicode/utf8"

    _ "github.com/go-sql-driver/mysql"

    "github.com/k4k3ru-hub/db/go/mysql/helper"
)

const (
    DefaultRequestTxTableName = "payment_onchain_request_txs"
)

var (
    requestTxIDCounter = &helper.IdCounter{}
)

type RequestTx struct {
    ID          uint64    `json:"id,string"`
    RequestID   uint64    `json:"requestId,string"`
    AccountID   uint64    `json:"accountId,string"`
    Chain       Chain     `json:"chain"`
    Network     Network   `json:"network"`
    Token       Token     `json:"token"`
    BlockNumber uint64    `json:"blockNumber"`
    TxHash      string    `json:"txHash"`
    FromAddress string    `json:"fromAddress"`
    ToAddress   string    `json:"toAddress"`
    Amount      string    `json:"amount"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}

type RequestTxStore struct {
    executor         helper.Executor
    tableName        string
    requestTableName string
}

type RequestTxInsertParams struct {
    ID          uint64    `json:"id,string"`
    RequestID   uint64    `json:"requestId,string"`
    AccountID   uint64    `json:"accountId,string"`
    Chain       Chain     `json:"chain"`
    Network     Network   `json:"network"`
    Token       Token     `json:"token"`
    BlockNumber uint64    `json:"blockNumber"`
    TxHash      string    `json:"txHash"`
    FromAddress string    `json:"fromAddress"`
    ToAddress   string    `json:"toAddress"`
    Amount      string    `json:"amount"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
    Ignore      bool      `json:"ignore"`
}

type RequestTxSelectParams struct {
    RequestID   *uint64   `json:"requestId,string,omitempty"`
    AccountID   *uint64   `json:"accountId,string,omitempty"`
    Chain       *Chain    `json:"chain,omitempty"`
    Network     *Network  `json:"network,omitempty"`
    Token       *Token    `json:"token,omitempty"`
    TxHash      *string   `json:"txHash,omitempty"`
    ToAddress   *string   `json:"toAddress,omitempty"`
    OrderBy     string    `json:"orderBy"`
    OrderByDesc bool      `json:"orderByDesc"`
    Limit       int       `json:"limit"`
    Offset      int       `json:"offset"`
}

type RequestTxUpdateParams struct {
    ID            uint64     `json:"id,string"`
    RequestID     *uint64    `json:"requestId,string,omitempty"`
}


//
// Generate payment onchain request tx ID.
//
// Version:
//   - 2026-05-16: Added.
//
func GenerateRequestTxID() uint64 {
    return requestTxIDCounter.GenerateID()
}


//
// Create new payment onchain request tx store.
//
// Version:
//   - 2026-05-16: Added.
//
func NewRequestTxStore(executor helper.Executor, tableName, requestTableName string) (*RequestTxStore, error) {
    if executor == nil {
        return nil, fmt.Errorf("failed to create payment onchain request tx store: missing required parameter: executor=null")
    }
    if tableName == "" {
        return nil, fmt.Errorf("failed to create payment onchain request tx store: missing required parameter: table_name=%q", "empty")
    }
    if requestTableName == "" {
        return nil, fmt.Errorf("failed to create payment onchain request tx store: missing required parameter: request_table_name=%q", "empty")
    }

    return &RequestTxStore{
        executor:         executor,
        tableName:        tableName,
        requestTableName: requestTableName,
    }, nil
}


//
// Validate payment onchain request tx ID.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateRequestTxID(id uint64) error {
    if id == 0 {
        return fmt.Errorf("invalid parameter: id=0")
    }
    return nil
}


//
// Validate payment onchain request tx ID.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *RequestTx) ValidateID() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_request_tx=null")
    }
    return ValidateRequestTxID(t.ID)
}


//
// Validate payment onchain request tx request ID.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateRequestTxRequestID(requestID uint64) error {
    if requestID == 0 {
        return fmt.Errorf("invalid parameter: request_id=0")
    }
    return nil
}


//
// Validate payment onchain request tx request ID.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *RequestTx) ValidateRequestID() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_request_tx=null")
    }
    return ValidateRequestTxRequestID(t.RequestID)
}


//
// Validate payment onchain request tx account ID.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateRequestTxAccountID(accountID uint64) error {
    if accountID == 0 {
        return fmt.Errorf("invalid parameter: account_id=0")
    }
    return nil
}


//
// Validate payment onchain request tx account ID.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *RequestTx) ValidateAccountID() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_request_tx=null")
    }
    return ValidateRequestTxAccountID(t.AccountID)
}


//
// Validate payment onchain request tx chain.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateRequestTxChain(c Chain) error {
    if err := c.Validate(); err != nil {
        return err
	}
    return nil
}


//
// Validate payment onchain request tx chain.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *RequestTx) ValidateChain() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_request_tx=null")
    }
    return ValidateRequestTxChain(t.Chain)
}


//
// Validate payment onchain request tx network.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateRequestTxNetwork(n Network) error {
    if err := n.Validate(); err != nil {
        return err
    }
    return nil
}


//
// Validate payment onchain request tx network.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *RequestTx) ValidateNetwork() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_request_tx=null")
    }
    return ValidateRequestTxNetwork(t.Network)
}


//
// Validate payment onchain request tx token.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateRequestTxToken(t Token) error {
    if err := t.Validate(); err != nil {
        return err
    }
    return nil
}


//
// Validate payment onchain request tx token.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *RequestTx) ValidateToken() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_request_tx=null")
    }
    return ValidateRequestTxToken(t.Token)
}


//
// Validate payment onchain request tx block number.
//
// Version:
//   - 2026-05-27: Added.
//
func ValidateRequestTxBlockNumber(blockNumber uint64) error {
    if blockNumber == 0 {
        return fmt.Errorf("invalid parameter: block_number=0")
    }
    return nil
}


//
// Validate payment onchain request tx block number.
//
// Version:
//   - 2026-05-27: Added.
//
func (t *RequestTx) ValidateBlockNumber() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_request_tx=null")
    }
    return ValidateRequestTxBlockNumber(t.BlockNumber)
}


//
// Validate payment onchain request tx tx hash.
//  
// Version:
//   - 2026-05-16: Added.
//
func ValidateRequestTxTxHash(txHash string) error {
    if txHash == "" {
        return fmt.Errorf("invalid parameter: tx_hash=%q", "empty")
    }
    if utf8.RuneCountInString(txHash) > 255 {
        return fmt.Errorf("invalid parameter: tx_hash=%q", helper.TruncateRunes(txHash, 255))
    }
    return nil
}


//
// Validate payment onchain request tx tx hash.
//  
// Version:
//   - 2026-05-16: Added.
//
func (t *RequestTx) ValidateRequestTxTxHash() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_request_tx=null")
    }
    return ValidateRequestTxTxHash(t.TxHash)
}


//
// Validate payment onchain request tx from address.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateRequestTxFromAddress(fromAddress string) error {
    if fromAddress == "" {
        return nil
    }
    if utf8.RuneCountInString(fromAddress) > 255 {
        return fmt.Errorf("invalid parameter: from_address=%q", helper.TruncateRunes(fromAddress, 255))
    }
    return nil
}


//
// Validate payment onchain request tx from address.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *RequestTx) ValidateFromAddress() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_request_tx=null")
    }
    return ValidateRequestTxFromAddress(t.FromAddress)
}


//
// Validate payment onchain request tx to address.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateRequestTxToAddress(toAddress string) error {
    if toAddress == "" {
        return fmt.Errorf("invalid parameter: to_address=empty")
    }
    if utf8.RuneCountInString(toAddress) > 255 {
        return fmt.Errorf("invalid parameter: to_address=%q", helper.TruncateRunes(toAddress, 255))
    }
    return nil
}


//
// Validate payment onchain request tx to address.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *RequestTx) ValidateToAddress() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_request_tx=null")
    }
    return ValidateRequestTxToAddress(t.ToAddress)
}


//
// Validate payment onchain request tx amount.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateRequestTxAmount(amount string) error {
    if amount == "" {
        return fmt.Errorf("invalid parameter: amount=%q", "empty")
    }
    if utf8.RuneCountInString(amount) > 78 {
        return fmt.Errorf("invalid parameter: amount=%q", helper.TruncateRunes(amount, 78))
    }
    return nil
}


//
// Validate payment onchain request tx amount.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *RequestTx) ValidateAmount() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_request_tx=null")
    }
    return ValidateRequestTxAmount(t.Amount)
}


//
// Create payment onchain request txs table.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *RequestTxStore) CreateTable() error {
    if s == nil {
        return fmt.Errorf("failed to create payment onchain request txs table: missing required parameter: request_tx_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to create payment onchain request txs table: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to create payment onchain request txs table: missing required parameter: table_name=%q", "empty")
    }
    if s.requestTableName == "" {
        return fmt.Errorf("failed to create payment onchain request txs table: missing required parameter: request_table_name=%q", "empty")
    }

    query := fmt.Sprintf(
        `CREATE TABLE IF NOT EXISTS %s (
            %s BIGINT UNSIGNED NOT NULL COMMENT 'ID',
            %s BIGINT UNSIGNED NOT NULL COMMENT 'Request ID',
            %s BIGINT UNSIGNED NOT NULL COMMENT 'Account ID',
            %s VARCHAR(64) NOT NULL COMMENT 'Chain',
            %s VARCHAR(64) NOT NULL COMMENT 'Network',
            %s VARCHAR(64) NOT NULL COMMENT 'Token',
            %s BIGINT UNSIGNED NOT NULL COMMENT 'Block number',
            %s VARCHAR(255) NOT NULL COMMENT 'Transaction hash',
            %s VARCHAR(255) NOT NULL COMMENT 'From address',
            %s VARCHAR(255) NOT NULL COMMENT 'Monitored recipient address',
            %s VARCHAR(78) NOT NULL COMMENT 'Amount',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
            PRIMARY KEY (%s),
            UNIQUE KEY uk_payment_onchain_request_txs_chain_network_tx_hash (%s, %s, %s),
            KEY idx_payment_onchain_request_txs_request_id (%s),
            KEY idx_payment_onchain_request_txs_account_id (%s),
            KEY idx_payment_onchain_request_txs_chain_network_to_address (%s, %s, %s),
            KEY idx_payment_onchain_request_txs_chain_network_token (%s, %s, %s),
            CONSTRAINT fk_payment_onchain_request_txs_request_id FOREIGN KEY (%s) REFERENCES %s (%s) ON DELETE CASCADE ON UPDATE CASCADE
        );`,
        s.tableName,
        ColID,
        ColRequestID,
        ColAccountID,
        ColChain,
        ColNetwork,
        ColToken,
        ColBlockNumber,
        ColTxHash,
        ColFromAddress,
        ColToAddress,
        ColAmount,
        ColCreatedAt,
        ColUpdatedAt,
        ColID,
        ColChain, ColNetwork, ColTxHash,
        ColRequestID,
        ColAccountID,
        ColChain, ColNetwork, ColToAddress,
        ColChain, ColNetwork, ColToken,
        ColRequestID, s.requestTableName, ColID,
    )

    if _, err := s.executor.Exec(query); err != nil {
        return fmt.Errorf("failed to create payment onchain request txs table: %w", err)
    }

    return nil
}


//
// Insert payment onchain request tx.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *RequestTxStore) Insert(p *RequestTxInsertParams) error {
    if s == nil {
        return fmt.Errorf("failed to insert payment onchain request tx: missing required parameter: request_tx_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to insert payment onchain request tx: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to insert payment onchain request tx: missing required parameter: table_name=%q", "empty")
    }
    if p == nil {
        return fmt.Errorf("failed to insert payment onchain request tx: missing required parameter: tx_insert_params=null")
    }
    if err := ValidateRequestTxAccountID(p.AccountID); err != nil {
        return fmt.Errorf("failed to insert payment onchain request tx: %w", err)
    }
    if err := ValidateRequestTxRequestID(p.RequestID); err != nil {
        return fmt.Errorf("failed to insert payment onchain request tx: %w", err)
    }
    if err := ValidateRequestTxChain(p.Chain); err != nil {
        return fmt.Errorf("failed to insert payment onchain request tx: %w", err)
    }
    if err := ValidateRequestTxNetwork(p.Network); err != nil {
        return fmt.Errorf("failed to insert payment onchain request tx: %w", err)
    }
    if err := ValidateRequestTxToken(p.Token); err != nil {
        return fmt.Errorf("failed to insert payment onchain request tx: %w", err)
    }
    if err := ValidateRequestTxBlockNumber(p.BlockNumber); err != nil {
        return fmt.Errorf("failed to insert payment onchain request tx: %w", err)
    }
    if err := ValidateRequestTxTxHash(p.TxHash); err != nil {
        return fmt.Errorf("failed to insert payment onchain request tx: %w", err)
    }
    if err := ValidateRequestTxFromAddress(p.FromAddress); err != nil {
        return fmt.Errorf("failed to insert payment onchain request tx: %w", err)
    }
    if err := ValidateRequestTxToAddress(p.ToAddress); err != nil {
        return fmt.Errorf("failed to insert payment onchain request tx: %w", err)
    }
    if err := ValidateRequestTxAmount(p.Amount); err != nil {
        return fmt.Errorf("failed to insert payment onchain request tx: %w", err)
    }

    if p.ID == 0 {
        p.ID = GenerateRequestTxID()
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
        "%s INTO %s (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
        queryPrefix,
        s.tableName,
        ColID,
        ColRequestID,
        ColAccountID,
        ColChain,
        ColNetwork,
        ColToken,
        ColBlockNumber,
        ColTxHash,
        ColFromAddress,
        ColToAddress,
        ColAmount,
        ColCreatedAt,
        ColUpdatedAt,
    )

    if _, err := s.executor.Exec(
        query,
        p.ID,
        p.RequestID,
        p.AccountID,
        p.Chain,
        p.Network,
        p.Token,
        p.BlockNumber,
        p.TxHash,
        p.FromAddress,
        p.ToAddress,
        p.Amount,
        p.CreatedAt,
        p.UpdatedAt,
    ); err != nil {
        return fmt.Errorf("failed to insert payment onchain request tx: %w", err)
    }

    return nil
}


//
// Upsert payment onchain request tx.
//
// Version:
//   - 2026-05-27: Added.
//
func (s *RequestTxStore) Upsert(p *RequestTxInsertParams) error {
    if s == nil {
        return fmt.Errorf("failed to upsert payment onchain request tx: missing required parameter: request_tx_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to upsert payment onchain request tx: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to upsert payment onchain request tx: missing required parameter: table_name=%q", "empty")
    }
    if p == nil {
        return fmt.Errorf("failed to upsert payment onchain request tx: missing required parameter: tx_insert_params=null")
    }
    if err := ValidateRequestTxAccountID(p.AccountID); err != nil {
        return fmt.Errorf("failed to upsert payment onchain request tx: %w", err)
    }
    if err := ValidateRequestTxRequestID(p.RequestID); err != nil {
        return fmt.Errorf("failed to upsert payment onchain request tx: %w", err)
    }
    if err := ValidateRequestTxChain(p.Chain); err != nil {
        return fmt.Errorf("failed to upsert payment onchain request tx: %w", err)
    }
    if err := ValidateRequestTxNetwork(p.Network); err != nil {
        return fmt.Errorf("failed to upsert payment onchain request tx: %w", err)
    }
    if err := ValidateRequestTxToken(p.Token); err != nil {
        return fmt.Errorf("failed to upsert payment onchain request tx: %w", err)
    }
    if err := ValidateRequestTxBlockNumber(p.BlockNumber); err != nil {
        return fmt.Errorf("failed to upsert payment onchain request tx: %w", err)
    }
    if err := ValidateRequestTxTxHash(p.TxHash); err != nil {
        return fmt.Errorf("failed to upsert payment onchain request tx: %w", err)
    }
    if err := ValidateRequestTxFromAddress(p.FromAddress); err != nil {
        return fmt.Errorf("failed to upsert payment onchain request tx: %w", err)
    }
    if err := ValidateRequestTxToAddress(p.ToAddress); err != nil {
        return fmt.Errorf("failed to upsert payment onchain request tx: %w", err)
    }
    if err := ValidateRequestTxAmount(p.Amount); err != nil {
        return fmt.Errorf("failed to upsert payment onchain request tx: %w", err)
    }

    if p.ID == 0 {
        p.ID = GenerateRequestTxID()
    }

    now := time.Now()
    if p.CreatedAt.IsZero() {
        p.CreatedAt = now
    }
    if p.UpdatedAt.IsZero() {
        p.UpdatedAt = now
    }

    query := fmt.Sprintf(
        `INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
         ON DUPLICATE KEY UPDATE
            %s = VALUES(%s),
            %s = VALUES(%s),
            %s = VALUES(%s),
            %s = VALUES(%s),
            %s = VALUES(%s),
            %s = VALUES(%s),
            %s = VALUES(%s),
            %s = VALUES(%s),
            %s = VALUES(%s);`,
        s.tableName,
        ColID,
        ColRequestID,
        ColAccountID,
        ColChain,
        ColNetwork,
        ColToken,
        ColBlockNumber,
        ColTxHash,
        ColFromAddress,
        ColToAddress,
        ColAmount,
        ColCreatedAt,
        ColUpdatedAt,
        ColRequestID, ColRequestID,
        ColAccountID, ColAccountID,
        ColChain, ColChain,
        ColNetwork, ColNetwork,
        ColToken, ColToken,
        ColBlockNumber, ColBlockNumber,
        ColFromAddress, ColFromAddress,
        ColToAddress, ColToAddress,
        ColUpdatedAt, ColUpdatedAt,
    )

    if _, err := s.executor.Exec(
        query,
        p.ID,
        p.RequestID,
        p.AccountID,
        p.Chain,
        p.Network,
        p.Token,
        p.BlockNumber,
        p.TxHash,
        p.FromAddress,
        p.ToAddress,
        p.Amount,
        p.CreatedAt,
        p.UpdatedAt,
    ); err != nil {
        return fmt.Errorf("failed to upsert payment onchain request tx: %w", err)
    }

    return nil
}


//
// Select payment onchain request txs.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *RequestTxStore) Select(p *RequestTxSelectParams) ([]*RequestTx, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select payment onchain request txs: missing required parameter: request_tx_store=null")
    }
    if s.executor == nil {
        return nil, fmt.Errorf("failed to select payment onchain request txs: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select payment onchain request txs: missing required parameter: table_name=%q", "empty")
    }
    if err := p.Validate(); err != nil {
        return nil, fmt.Errorf("failed to select payment onchain request txs: %w", err)
    }

    query, args := p.BuildQuery("SELECT * FROM " + s.tableName)

    // Execute.
    rows, err := s.executor.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to select payment onchain request txs: %w", err)
    }
    defer rows.Close()

    var result []*RequestTx
    for rows.Next() {
        row := &RequestTx{}
        err := rows.Scan(
            &row.ID,
            &row.RequestID,
            &row.AccountID,
            &row.Chain,
            &row.Network,
            &row.Token,
            &row.BlockNumber,
            &row.TxHash,
            &row.FromAddress,
            &row.ToAddress,
            &row.Amount,
            &row.CreatedAt,
            &row.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to select payment onchain request txs: %w", err)
        }

        result = append(result, row)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to select payment onchain request txs: %w", err)
    }

    return result, nil
}


//
// Select payment onchain request tx by ID.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *RequestTxStore) SelectByID(id uint64) (*RequestTx, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select payment onchain request tx by id: missing required parameter: request_tx_store=null")
    }
    if s.executor == nil {
        return nil, fmt.Errorf("failed to select payment onchain request tx by id: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select payment onchain request tx by id: missing required parameter: table_name=%q", "empty")
    }
    if id == 0 {
        return nil, fmt.Errorf("failed to select payment onchain request tx by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColID)

    result := &RequestTx{}
    err := s.executor.QueryRow(query, id).Scan(
        &result.ID,
        &result.RequestID,
        &result.AccountID,
        &result.Chain,
        &result.Network,
        &result.Token,
        &result.BlockNumber,
        &result.TxHash,
        &result.FromAddress,
        &result.ToAddress,
        &result.Amount,
        &result.CreatedAt,
        &result.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to select payment onchain request tx by id: %w", err)
    }

    return result, nil
}


//
// Count payment onchain request txs.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *RequestTxStore) Count(p *RequestTxSelectParams) (int64, error) {
    if s == nil {
        return 0, fmt.Errorf("failed to count payment onchain request txs: missing required parameter: request_tx_store=null")
    }
    if s.executor == nil {
        return 0, fmt.Errorf("failed to count payment onchain request txs: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return 0, fmt.Errorf("failed to count payment onchain request txs: missing required parameter: table_name=%q", "empty")
    }
    if err := p.Validate(); err != nil {
        return 0, fmt.Errorf("failed to count payment onchain request txs: %w", err)
    }

    query, args := p.BuildQuery("SELECT COUNT(*) FROM " + s.tableName)

    // Execute.
    var result int64
    if err := s.executor.QueryRow(query, args...).Scan(&result); err != nil {
        return 0, fmt.Errorf("failed to count payment onchain request txs: %w", err)
    }

    return result, nil
}


//
// Delete payment onchain request tx by ID.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *RequestTxStore) DeleteByID(id uint64) error {
    if s == nil {
        return fmt.Errorf("failed to delete payment onchain request tx by id: missing required parameter: request_tx_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to delete payment onchain request tx by id: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to delete payment onchain request tx by id: missing required parameter: table_name=%q", "empty")
    }
    if id == 0 {
        return fmt.Errorf("failed to delete payment onchain request tx by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", s.tableName, ColID)

    if _, err := s.executor.Exec(query, id); err != nil {
        return fmt.Errorf("failed to delete payment onchain request tx by id: %w", err)
    }

    return nil
}



func (s *RequestTxStore) Update(option *RequestTxUpdateParams) error {
    if s == nil {
        return fmt.Errorf("failed to update payment onchain request tx: missing required parameter: request_tx_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to update payment onchain request tx: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update payment onchain request tx: missing required parameter: table_name=%q", "empty")
    }
    if option == nil {
        return fmt.Errorf("failed to update payment onchain request tx: missing required parameter: option=null")
    }
    if option.ID == 0 {
        return fmt.Errorf("failed to update payment onchain request tx: invalid parameter: id=0")
    }

    assignments := make([]string, 0, 4)
    args := make([]any, 0, 4)

    if option.RequestID != nil {
        if *option.RequestID == 0 {
            return fmt.Errorf("failed to update payment onchain request tx: invalid parameter: request_id=0")
        }
        assignments = append(assignments, ColRequestID+" = ?")
        args = append(args, *option.RequestID)
    }

    if len(assignments) == 0 {
        return fmt.Errorf("failed to update payment onchain request tx: invalid parameter: assignments=empty")
    }

    args = append(args, option.ID)

    query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?;", s.tableName, strings.Join(assignments, ", "), ColID)

    if _, err := s.executor.Exec(query, args...); err != nil {
        return fmt.Errorf("failed to update payment onchain request tx: %w", err)
    }

    return nil
}


//
// Sum confirmed amount by request ID.
//
// Version:
//   - 2026-05-30: Added.
//
func (s *RequestTxStore) SumConfirmedAmountByRequestID(requestID uint64, latestBlockNumber uint64, requiredConfirmations uint64) (*big.Int, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to sum confirmed amount by request id: missing required parameter: request_tx_store=null")
    }
    if s.executor == nil {
        return nil, fmt.Errorf("failed to sum confirmed amount by request id: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to sum confirmed amount by request id: missing required parameter: table_name=%q", "empty")
    }
    if requestID == 0 {
        return nil, fmt.Errorf("failed to sum confirmed amount by request id: invalid parameter: request_id=0")
    }
    if latestBlockNumber == 0 {
        return nil, fmt.Errorf("failed to sum confirmed amount by request id: invalid parameter: latest_block_number=0")
    }
    if requiredConfirmations == 0 {
        return nil, fmt.Errorf("failed to sum confirmed amount by request id: invalid parameter: required_confirmations=0")
    }

    query := fmt.Sprintf(
        "SELECT %s FROM %s WHERE %s = ? AND ? >= %s AND (? - %s + 1) >= ?;",
        ColAmount,
        s.tableName,
        ColRequestID,
        ColBlockNumber,
        ColBlockNumber,
    )

    args := make([]any, 0, 4)
    args = append(args, requestID, latestBlockNumber, latestBlockNumber, requiredConfirmations)

    // Execute.
    rows, err := s.executor.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to sum confirmed amount by request id: %w", err)
    }
    defer rows.Close()

    result := new(big.Int)

    for rows.Next() {
        var amountStr string
        if err := rows.Scan(&amountStr); err != nil {
            return nil, fmt.Errorf("failed to sum confirmed amount by request id: %w", err)
        }

        amount, ok := new(big.Int).SetString(amountStr, 10)
        if !ok {
            return nil, fmt.Errorf("invalid amount: %q", amountStr)
        }

        result.Add(result, amount)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to sum confirmed amount by request id: %w", err)
    }

    return result, nil
}


//
// Build query.
//
// Version:
//   - 2025-05-16: Added.
//
func (p *RequestTxSelectParams) BuildQuery(selectFromClause string) (string, []any) {
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
    if p.Chain != nil {
        conditions = append(conditions, ColChain + " = ?")
        args = append(args, *p.Chain)
    }
    if p.Network != nil {
        conditions = append(conditions, ColNetwork + " = ?")
        args = append(args, *p.Network)
    }
    if p.Token != nil {
        conditions = append(conditions, ColToken + " = ?")
        args = append(args, *p.Token)
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
func (p *RequestTxSelectParams) Validate() error {
    // Guard.
    if p == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_request_tx_select_params=null")
    }

    if p.RequestID != nil {
        if err := ValidateRequestTxRequestID(*p.RequestID); err != nil {
            return err
        }
    }
    if p.AccountID != nil {
        if err := ValidateRequestTxAccountID(*p.AccountID); err != nil {
            return err
        }
    }
    if p.Chain != nil {
        if err := ValidateRequestTxChain(*p.Chain); err != nil {
            return err
        }
    }
    if p.Network != nil {
        if err := ValidateRequestTxNetwork(*p.Network); err != nil {
            return err
        }
    }
    if p.Token != nil {
        if err := ValidateRequestTxToken(*p.Token); err != nil {
            return err
        }
    }
    if p.TxHash != nil {
        if err := ValidateRequestTxTxHash(*p.TxHash); err != nil {
            return err
        }
    }
    if p.ToAddress != nil {
        if err := ValidateRequestTxToAddress(*p.ToAddress); err != nil {
            return err
        }
    }

    if p.OrderBy != "" {
        switch p.OrderBy {
        case ColID,
            ColRequestID,
            ColAccountID,
            ColChain,
            ColNetwork,
            ColToken,
            ColBlockNumber,
            ColTxHash,
            ColFromAddress,
            ColToAddress,
            ColAmount,
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


