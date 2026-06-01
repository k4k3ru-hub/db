//
// intent_transfer.go
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
    DefaultIntentTransferTableName = "payment_onchain_intent_transfers"
)

var (
    intentTransferIDCounter = &helper.IdCounter{}
)

type IntentTransfer struct {
    ID          uint64    `json:"id,string"`
    IntentID    uint64    `json:"intentId,string"`
    OwnerRef    string    `json:"ownerRef"`
    Chain       Chain     `json:"chain"`
    Network     Network   `json:"network"`
    Token       Token     `json:"token"`
    BlockNumber uint64    `json:"blockNumber"`
    TxHash      string    `json:"txHash"`
    EventIndex  uint64    `json:"eventIndex"`
    FromAddress string    `json:"fromAddress"`
    ToAddress   string    `json:"toAddress"`
    Amount      string    `json:"amount"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}

type IntentTransferStore struct {
    executor         helper.Executor
    tableName        string
    intentTableName string
}

type IntentTransferInsertParams struct {
    ID          uint64    `json:"id,string"`
    IntentID    uint64    `json:"intentId,string"`
    OwnerRef    string    `json:"ownerRef"`
    Chain       Chain     `json:"chain"`
    Network     Network   `json:"network"`
    Token       Token     `json:"token"`
    BlockNumber uint64    `json:"blockNumber"`
    TxHash      string    `json:"txHash"`
    EventIndex  uint64    `json:"eventIndex"`
    FromAddress string    `json:"fromAddress"`
    ToAddress   string    `json:"toAddress"`
    Amount      string    `json:"amount"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
    Ignore      bool      `json:"ignore"`
}

type IntentTransferSelectParams struct {
    IntentID    *uint64   `json:"intentId,string,omitempty"`
    OwnerRef    *string   `json:"ownerRef,omitempty"`
    Chain       *Chain    `json:"chain,omitempty"`
    Network     *Network  `json:"network,omitempty"`
    Token       *Token    `json:"token,omitempty"`
    TxHash      *string   `json:"txHash,omitempty"`
    EventIndex  *uint64   `json:"eventIndex,omitempty"`
    ToAddress   *string   `json:"toAddress,omitempty"`
    OrderBy     string    `json:"orderBy"`
    OrderByDesc bool      `json:"orderByDesc"`
    Limit       int       `json:"limit"`
    Offset      int       `json:"offset"`
}

type IntentTransferUpdateParams struct {
    ID        uint64 `json:"id,string"`
    IntentID *uint64 `json:"intentId,string,omitempty"`
}


//
// Generate payment onchain intent transfer ID.
//
// Version:
//   - 2026-05-16: Added.
//
func GenerateIntentTransferID() uint64 {
    return intentTransferIDCounter.GenerateID()
}


//
// Create new payment onchain intent transfer store.
//
// Version:
//   - 2026-05-16: Added.
//
func NewIntentTransferStore(executor helper.Executor, tableName, intentTableName string) (*IntentTransferStore, error) {
    if executor == nil {
        return nil, fmt.Errorf("failed to create payment onchain intent transfer store: missing required parameter: executor=null")
    }
    if tableName == "" {
        return nil, fmt.Errorf("failed to create payment onchain intent transfer store: missing required parameter: table_name=%q", "empty")
    }
    if intentTableName == "" {
        return nil, fmt.Errorf("failed to create payment onchain intent transfer store: missing required parameter: request_table_name=%q", "empty")
    }

    return &IntentTransferStore{
        executor:        executor,
        tableName:       tableName,
        intentTableName: intentTableName,
    }, nil
}


//
// Validate payment onchain intent transfer ID.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateIntentTransferID(id uint64) error {
    if id == 0 {
        return fmt.Errorf("invalid parameter: id=0")
    }
    return nil
}


//
// Validate payment onchain intent transfer ID.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *IntentTransfer) ValidateID() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_intent_transfer=null")
    }
    return ValidateIntentTransferID(t.ID)
}


//
// Validate payment onchain intent transfer intent ID.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateIntentTransferIntentID(intentID uint64) error {
    if intentID == 0 {
        return fmt.Errorf("invalid parameter: request_id=0")
    }
    return nil
}


//
// Validate payment onchain intent transfer intent ID.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *IntentTransfer) ValidateIntentID() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_intent_transfer=null")
    }
    return ValidateIntentTransferIntentID(t.IntentID)
}


//
// Validate payment onchain intent transfer owner ref.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateIntentTransferOwnerRef(ownerRef string) error {
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
// Validate payment onchain intent transfer owner ref.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *IntentTransfer) ValidateOwnerRef() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_intent_transfer=null")
    }
    return ValidateIntentTransferOwnerRef(t.OwnerRef)
}


//
// Validate payment onchain intent transfer chain.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateIntentTransferChain(c Chain) error {
    if err := c.Validate(); err != nil {
        return err
	}
    return nil
}


//
// Validate payment onchain intent transfer chain.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *IntentTransfer) ValidateChain() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_intent_transfer=null")
    }
    return ValidateIntentTransferChain(t.Chain)
}


//
// Validate payment onchain intent transfer network.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateIntentTransferNetwork(n Network) error {
    if err := n.Validate(); err != nil {
        return err
    }
    return nil
}


//
// Validate payment onchain intent transfer network.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *IntentTransfer) ValidateNetwork() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_intent_transfer=null")
    }
    return ValidateIntentTransferNetwork(t.Network)
}


//
// Validate payment onchain intent transfer token.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateIntentTransferToken(t Token) error {
    if err := t.Validate(); err != nil {
        return err
    }
    return nil
}


//
// Validate payment onchain intent transfer token.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *IntentTransfer) ValidateToken() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_intent_transfer=null")
    }
    return ValidateIntentTransferToken(t.Token)
}


//
// Validate payment onchain intent transfer block number.
//
// Version:
//   - 2026-05-27: Added.
//
func ValidateIntentTransferBlockNumber(blockNumber uint64) error {
    if blockNumber == 0 {
        return fmt.Errorf("invalid parameter: block_number=0")
    }
    return nil
}


//
// Validate payment onchain intent transfer block number.
//
// Version:
//   - 2026-05-27: Added.
//
func (t *IntentTransfer) ValidateBlockNumber() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_intent_transfer=null")
    }
    return ValidateIntentTransferBlockNumber(t.BlockNumber)
}


//
// Validate payment onchain intent transfer tx hash.
//  
// Version:
//   - 2026-05-16: Added.
//
func ValidateIntentTransferTxHash(txHash string) error {
    if txHash == "" {
        return fmt.Errorf("invalid parameter: tx_hash=%q", "empty")
    }
    if utf8.RuneCountInString(txHash) > 255 {
        return fmt.Errorf("invalid parameter: tx_hash=%q", helper.TruncateRunes(txHash, 255))
    }
    return nil
}


//
// Validate payment onchain intent transfer tx hash.
//  
// Version:
//   - 2026-05-16: Added.
//
func (t *IntentTransfer) ValidateIntentTransferTxHash() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_intent_transfer=null")
    }
    return ValidateIntentTransferTxHash(t.TxHash)
}


//
// Validate payment onchain intent transfer event index.
//
// Version:
//   - 2026-05-31: Added.
//
func ValidateIntentTransferEventIndex(eventIndex uint64) error {
    if eventIndex == 0 {
        return fmt.Errorf("invalid parameter: event_index=0")
    }
    return nil
}


//
// Validate payment onchain intent transfer event index.
//
// Version:
//   - 2026-05-31: Added.
//
func (t *IntentTransfer) ValidateEventIndex() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_intent_transfer=null")
    }
    return ValidateIntentTransferEventIndex(t.EventIndex)
}


//
// Validate payment onchain intent transfer from address.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateIntentTransferFromAddress(fromAddress string) error {
    if fromAddress == "" {
        return nil
    }
    if utf8.RuneCountInString(fromAddress) > 255 {
        return fmt.Errorf("invalid parameter: from_address=%q", helper.TruncateRunes(fromAddress, 255))
    }
    return nil
}


//
// Validate payment onchain intent transfer from address.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *IntentTransfer) ValidateFromAddress() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_intent_transfer=null")
    }
    return ValidateIntentTransferFromAddress(t.FromAddress)
}


//
// Validate payment onchain intent transfer to address.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateIntentTransferToAddress(toAddress string) error {
    if toAddress == "" {
        return fmt.Errorf("invalid parameter: to_address=empty")
    }
    if utf8.RuneCountInString(toAddress) > 255 {
        return fmt.Errorf("invalid parameter: to_address=%q", helper.TruncateRunes(toAddress, 255))
    }
    return nil
}


//
// Validate payment onchain intent transfer to address.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *IntentTransfer) ValidateToAddress() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_intent_transfer=null")
    }
    return ValidateIntentTransferToAddress(t.ToAddress)
}


//
// Validate payment onchain intent transfer amount.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateIntentTransferAmount(amount string) error {
    if amount == "" {
        return fmt.Errorf("invalid parameter: amount=%q", "empty")
    }
    if utf8.RuneCountInString(amount) > 78 {
        return fmt.Errorf("invalid parameter: amount=%q", helper.TruncateRunes(amount, 78))
    }
    return nil
}


//
// Validate payment onchain intent transfer amount.
//
// Version:
//   - 2026-05-16: Added.
//
func (t *IntentTransfer) ValidateAmount() error {
    if t == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_intent_transfer=null")
    }
    return ValidateIntentTransferAmount(t.Amount)
}


//
// Create payment onchain intent transfers table.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *IntentTransferStore) CreateTable() error {
    if s == nil {
        return fmt.Errorf("failed to create payment onchain intent transfers table: missing required parameter: request_tx_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to create payment onchain intent transfers table: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to create payment onchain intent transfers table: missing required parameter: table_name=%q", "empty")
    }
    if s.intentTableName == "" {
        return fmt.Errorf("failed to create payment onchain intent transfers table: missing required parameter: request_table_name=%q", "empty")
    }

    query := fmt.Sprintf(
        `CREATE TABLE IF NOT EXISTS %s (
            %s BIGINT UNSIGNED NOT NULL COMMENT 'ID',
            %s BIGINT UNSIGNED NOT NULL COMMENT 'Intent ID',
            %s VARCHAR(128) NOT NULL COMMENT 'Owner Ref',
            %s VARCHAR(64) NOT NULL COMMENT 'Chain',
            %s VARCHAR(64) NOT NULL COMMENT 'Network',
            %s VARCHAR(64) NOT NULL COMMENT 'Token',
            %s BIGINT UNSIGNED NOT NULL COMMENT 'Block number',
            %s VARCHAR(255) NOT NULL COMMENT 'Transaction hash',
            %s BIGINT UNSIGNED NOT NULL COMMENT 'Event index',
            %s VARCHAR(255) NOT NULL COMMENT 'From address',
            %s VARCHAR(255) NOT NULL COMMENT 'Monitored recipient address',
            %s VARCHAR(78) NOT NULL COMMENT 'Amount',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
            PRIMARY KEY (%s),
            UNIQUE KEY uk_payment_onchain_intent_transfers_chain_network_tx_hash_event_index (%s, %s, %s, %s),
            KEY idx_payment_onchain_intent_transfers_intent_id (%s),
            KEY idx_payment_onchain_intent_transfers_cha_net_to_add (%s, %s, %s),
            KEY idx_payment_onchain_intent_transfers_cha_net_tok (%s, %s, %s),
            CONSTRAINT fk_payment_onchain_intent_transfers_intent_id FOREIGN KEY (%s) REFERENCES %s (%s) ON DELETE CASCADE ON UPDATE CASCADE
        );`,
        s.tableName,
        ColID,
        ColIntentID,
        ColOwnerRef,
        ColChain,
        ColNetwork,
        ColToken,
        ColBlockNumber,
        ColTxHash,
        ColEventIndex,
        ColFromAddress,
        ColToAddress,
        ColAmount,
        ColCreatedAt,
        ColUpdatedAt,
        ColID,
        ColChain, ColNetwork, ColTxHash, ColEventIndex,
        ColIntentID,
        ColChain, ColNetwork, ColToAddress,
        ColChain, ColNetwork, ColToken,
        ColIntentID, s.intentTableName, ColID,
    )

    if _, err := s.executor.Exec(query); err != nil {
        return fmt.Errorf("failed to create payment onchain intent transfers table: %w", err)
    }

    return nil
}


//
// Insert payment onchain intent transfer.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *IntentTransferStore) Insert(p *IntentTransferInsertParams) error {
    if s == nil {
        return fmt.Errorf("failed to insert payment onchain intent transfer: missing required parameter: request_tx_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to insert payment onchain intent transfer: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to insert payment onchain intent transfer: missing required parameter: table_name=%q", "empty")
    }
    if p == nil {
        return fmt.Errorf("failed to insert payment onchain intent transfer: missing required parameter: tx_insert_params=null")
    }
    if err := ValidateIntentTransferOwnerRef(p.OwnerRef); err != nil {
        return fmt.Errorf("failed to insert payment onchain intent transfer: %w", err)
    }
    if err := ValidateIntentTransferIntentID(p.IntentID); err != nil {
        return fmt.Errorf("failed to insert payment onchain intent transfer: %w", err)
    }
    if err := ValidateIntentTransferChain(p.Chain); err != nil {
        return fmt.Errorf("failed to insert payment onchain intent transfer: %w", err)
    }
    if err := ValidateIntentTransferNetwork(p.Network); err != nil {
        return fmt.Errorf("failed to insert payment onchain intent transfer: %w", err)
    }
    if err := ValidateIntentTransferToken(p.Token); err != nil {
        return fmt.Errorf("failed to insert payment onchain intent transfer: %w", err)
    }
    if err := ValidateIntentTransferBlockNumber(p.BlockNumber); err != nil {
        return fmt.Errorf("failed to insert payment onchain intent transfer: %w", err)
    }
    if err := ValidateIntentTransferTxHash(p.TxHash); err != nil {
        return fmt.Errorf("failed to insert payment onchain intent transfer: %w", err)
    }
    if err := ValidateIntentTransferEventIndex(p.EventIndex); err != nil {
        return fmt.Errorf("failed to insert payment onchain intent transfer: %w", err)
    }
    if err := ValidateIntentTransferFromAddress(p.FromAddress); err != nil {
        return fmt.Errorf("failed to insert payment onchain intent transfer: %w", err)
    }
    if err := ValidateIntentTransferToAddress(p.ToAddress); err != nil {
        return fmt.Errorf("failed to insert payment onchain intent transfer: %w", err)
    }
    if err := ValidateIntentTransferAmount(p.Amount); err != nil {
        return fmt.Errorf("failed to insert payment onchain intent transfer: %w", err)
    }

    if p.ID == 0 {
        p.ID = GenerateIntentTransferID()
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
        "%s INTO %s (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
        queryPrefix,
        s.tableName,
        ColID,
        ColIntentID,
        ColOwnerRef,
        ColChain,
        ColNetwork,
        ColToken,
        ColBlockNumber,
        ColTxHash,
        ColEventIndex,
        ColFromAddress,
        ColToAddress,
        ColAmount,
        ColCreatedAt,
        ColUpdatedAt,
    )

    if _, err := s.executor.Exec(
        query,
        p.ID,
        p.IntentID,
        p.OwnerRef,
        p.Chain,
        p.Network,
        p.Token,
        p.BlockNumber,
        p.TxHash,
        p.EventIndex,
        p.FromAddress,
        p.ToAddress,
        p.Amount,
        p.CreatedAt,
        p.UpdatedAt,
    ); err != nil {
        return fmt.Errorf("failed to insert payment onchain intent transfer: %w", err)
    }

    return nil
}


//
// Upsert payment onchain intent transfer.
//
// Version:
//   - 2026-05-27: Added.
//
func (s *IntentTransferStore) Upsert(p *IntentTransferInsertParams) error {
    if s == nil {
        return fmt.Errorf("failed to upsert payment onchain intent transfer: missing required parameter: request_tx_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to upsert payment onchain intent transfer: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to upsert payment onchain intent transfer: missing required parameter: table_name=%q", "empty")
    }
    if p == nil {
        return fmt.Errorf("failed to upsert payment onchain intent transfer: missing required parameter: tx_insert_params=null")
    }
    if err := ValidateIntentTransferOwnerRef(p.OwnerRef); err != nil {
        return fmt.Errorf("failed to upsert payment onchain intent transfer: %w", err)
    }
    if err := ValidateIntentTransferIntentID(p.IntentID); err != nil {
        return fmt.Errorf("failed to upsert payment onchain intent transfer: %w", err)
    }
    if err := ValidateIntentTransferChain(p.Chain); err != nil {
        return fmt.Errorf("failed to upsert payment onchain intent transfer: %w", err)
    }
    if err := ValidateIntentTransferNetwork(p.Network); err != nil {
        return fmt.Errorf("failed to upsert payment onchain intent transfer: %w", err)
    }
    if err := ValidateIntentTransferToken(p.Token); err != nil {
        return fmt.Errorf("failed to upsert payment onchain intent transfer: %w", err)
    }
    if err := ValidateIntentTransferBlockNumber(p.BlockNumber); err != nil {
        return fmt.Errorf("failed to upsert payment onchain intent transfer: %w", err)
    }
    if err := ValidateIntentTransferTxHash(p.TxHash); err != nil {
        return fmt.Errorf("failed to upsert payment onchain intent transfer: %w", err)
    }
    if err := ValidateIntentTransferEventIndex(p.EventIndex); err != nil {
        return fmt.Errorf("failed to upsert payment onchain intent transfer: %w", err)
    }
    if err := ValidateIntentTransferFromAddress(p.FromAddress); err != nil {
        return fmt.Errorf("failed to upsert payment onchain intent transfer: %w", err)
    }
    if err := ValidateIntentTransferToAddress(p.ToAddress); err != nil {
        return fmt.Errorf("failed to upsert payment onchain intent transfer: %w", err)
    }
    if err := ValidateIntentTransferAmount(p.Amount); err != nil {
        return fmt.Errorf("failed to upsert payment onchain intent transfer: %w", err)
    }

    if p.ID == 0 {
        p.ID = GenerateIntentTransferID()
    }

    now := time.Now()
    if p.CreatedAt.IsZero() {
        p.CreatedAt = now
    }
    if p.UpdatedAt.IsZero() {
        p.UpdatedAt = now
    }

    query := fmt.Sprintf(
        `INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
         ON DUPLICATE KEY UPDATE
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
        ColIntentID,
        ColOwnerRef,
        ColChain,
        ColNetwork,
        ColToken,
        ColBlockNumber,
        ColTxHash,
        ColEventIndex,
        ColFromAddress,
        ColToAddress,
        ColAmount,
        ColCreatedAt,
        ColUpdatedAt,
        ColIntentID, ColIntentID,
        ColOwnerRef, ColOwnerRef,
        ColToken, ColToken,
        ColBlockNumber, ColBlockNumber,
        ColFromAddress, ColFromAddress,
        ColToAddress, ColToAddress,
        ColAmount, ColAmount,
        ColUpdatedAt, ColUpdatedAt,
    )

    if _, err := s.executor.Exec(
        query,
        p.ID,
        p.IntentID,
        p.OwnerRef,
        p.Chain,
        p.Network,
        p.Token,
        p.BlockNumber,
        p.TxHash,
        p.EventIndex,
        p.FromAddress,
        p.ToAddress,
        p.Amount,
        p.CreatedAt,
        p.UpdatedAt,
    ); err != nil {
        return fmt.Errorf("failed to upsert payment onchain intent transfer: %w", err)
    }

    return nil
}


//
// Select payment onchain intent transfers.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *IntentTransferStore) Select(p *IntentTransferSelectParams) ([]*IntentTransfer, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select payment onchain intent transfers: missing required parameter: request_tx_store=null")
    }
    if s.executor == nil {
        return nil, fmt.Errorf("failed to select payment onchain intent transfers: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select payment onchain intent transfers: missing required parameter: table_name=%q", "empty")
    }
    if err := p.Validate(); err != nil {
        return nil, fmt.Errorf("failed to select payment onchain intent transfers: %w", err)
    }

    query, args := p.BuildQuery("SELECT * FROM " + s.tableName)

    // Execute.
    rows, err := s.executor.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to select payment onchain intent transfers: %w", err)
    }
    defer rows.Close()

    var result []*IntentTransfer
    for rows.Next() {
        row := &IntentTransfer{}
        err := rows.Scan(
            &row.ID,
            &row.IntentID,
            &row.OwnerRef,
            &row.Chain,
            &row.Network,
            &row.Token,
            &row.BlockNumber,
            &row.TxHash,
            &row.EventIndex,
            &row.FromAddress,
            &row.ToAddress,
            &row.Amount,
            &row.CreatedAt,
            &row.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to select payment onchain intent transfers: %w", err)
        }

        result = append(result, row)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to select payment onchain intent transfers: %w", err)
    }

    return result, nil
}


//
// Select payment onchain intent transfer by ID.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *IntentTransferStore) SelectByID(id uint64) (*IntentTransfer, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select payment onchain intent transfer by id: missing required parameter: request_tx_store=null")
    }
    if s.executor == nil {
        return nil, fmt.Errorf("failed to select payment onchain intent transfer by id: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select payment onchain intent transfer by id: missing required parameter: table_name=%q", "empty")
    }
    if id == 0 {
        return nil, fmt.Errorf("failed to select payment onchain intent transfer by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColID)

    result := &IntentTransfer{}
    err := s.executor.QueryRow(query, id).Scan(
        &result.ID,
        &result.IntentID,
        &result.OwnerRef,
        &result.Chain,
        &result.Network,
        &result.Token,
        &result.BlockNumber,
        &result.TxHash,
        &result.EventIndex,
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
        return nil, fmt.Errorf("failed to select payment onchain intent transfer by id: %w", err)
    }

    return result, nil
}


//
// Count payment onchain intent transfers.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *IntentTransferStore) Count(p *IntentTransferSelectParams) (int64, error) {
    if s == nil {
        return 0, fmt.Errorf("failed to count payment onchain intent transfers: missing required parameter: request_tx_store=null")
    }
    if s.executor == nil {
        return 0, fmt.Errorf("failed to count payment onchain intent transfers: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return 0, fmt.Errorf("failed to count payment onchain intent transfers: missing required parameter: table_name=%q", "empty")
    }
    if err := p.Validate(); err != nil {
        return 0, fmt.Errorf("failed to count payment onchain intent transfers: %w", err)
    }

    query, args := p.BuildQuery("SELECT COUNT(*) FROM " + s.tableName)

    // Execute.
    var result int64
    if err := s.executor.QueryRow(query, args...).Scan(&result); err != nil {
        return 0, fmt.Errorf("failed to count payment onchain intent transfers: %w", err)
    }

    return result, nil
}


//
// Delete payment onchain intent transfer by ID.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *IntentTransferStore) DeleteByID(id uint64) error {
    if s == nil {
        return fmt.Errorf("failed to delete payment onchain intent transfer by id: missing required parameter: request_tx_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to delete payment onchain intent transfer by id: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to delete payment onchain intent transfer by id: missing required parameter: table_name=%q", "empty")
    }
    if id == 0 {
        return fmt.Errorf("failed to delete payment onchain intent transfer by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", s.tableName, ColID)

    if _, err := s.executor.Exec(query, id); err != nil {
        return fmt.Errorf("failed to delete payment onchain intent transfer by id: %w", err)
    }

    return nil
}



func (s *IntentTransferStore) Update(option *IntentTransferUpdateParams) error {
    if s == nil {
        return fmt.Errorf("failed to update payment onchain intent transfer: missing required parameter: request_tx_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to update payment onchain intent transfer: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update payment onchain intent transfer: missing required parameter: table_name=%q", "empty")
    }
    if option == nil {
        return fmt.Errorf("failed to update payment onchain intent transfer: missing required parameter: option=null")
    }
    if option.ID == 0 {
        return fmt.Errorf("failed to update payment onchain intent transfer: invalid parameter: id=0")
    }

    assignments := make([]string, 0, 4)
    args := make([]any, 0, 4)

    if option.IntentID != nil {
        if *option.IntentID == 0 {
            return fmt.Errorf("failed to update payment onchain intent transfer: invalid parameter: request_id=0")
        }
        assignments = append(assignments, ColIntentID+" = ?")
        args = append(args, *option.IntentID)
    }

    if len(assignments) == 0 {
        return fmt.Errorf("failed to update payment onchain intent transfer: invalid parameter: assignments=empty")
    }

    args = append(args, option.ID)

    query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?;", s.tableName, strings.Join(assignments, ", "), ColID)

    if _, err := s.executor.Exec(query, args...); err != nil {
        return fmt.Errorf("failed to update payment onchain intent transfer: %w", err)
    }

    return nil
}


//
// Sum confirmed amount by intent ID.
//
// Version:
//   - 2026-05-30: Added.
//
func (s *IntentTransferStore) SumConfirmedAmountByIntentID(intentID uint64, latestBlockNumber uint64, requiredConfirmations uint64) (*big.Int, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to sum confirmed amount by request id: missing required parameter: request_tx_store=null")
    }
    if s.executor == nil {
        return nil, fmt.Errorf("failed to sum confirmed amount by request id: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to sum confirmed amount by request id: missing required parameter: table_name=%q", "empty")
    }
    if intentID == 0 {
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
        ColIntentID,
        ColBlockNumber,
        ColBlockNumber,
    )

    args := make([]any, 0, 4)
    args = append(args, intentID, latestBlockNumber, latestBlockNumber, requiredConfirmations)

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
func (p *IntentTransferSelectParams) BuildQuery(selectFromClause string) (string, []any) {
    // Guard.
    if p == nil {
        return selectFromClause, nil
    }

    var query strings.Builder
    query.WriteString(selectFromClause)

    conditions := make([]string, 0, 8)
    args := make([]any, 0, 10)

    if p.IntentID != nil {
        conditions = append(conditions, ColIntentID + " = ?")
        args = append(args, *p.IntentID)
    }
    if p.OwnerRef != nil {
        conditions = append(conditions, ColOwnerRef + " = ?")
        args = append(args, *p.OwnerRef)
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
    if p.EventIndex != nil {
        conditions = append(conditions, ColEventIndex + " = ?")
        args = append(args, *p.EventIndex)
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
func (p *IntentTransferSelectParams) Validate() error {
    // Guard.
    if p == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_intent_transfer_select_params=null")
    }

    if p.IntentID != nil {
        if err := ValidateIntentTransferIntentID(*p.IntentID); err != nil {
            return err
        }
    }
    if p.OwnerRef != nil {
        if err := ValidateIntentTransferOwnerRef(*p.OwnerRef); err != nil {
            return err
        }
    }
    if p.Chain != nil {
        if err := ValidateIntentTransferChain(*p.Chain); err != nil {
            return err
        }
    }
    if p.Network != nil {
        if err := ValidateIntentTransferNetwork(*p.Network); err != nil {
            return err
        }
    }
    if p.Token != nil {
        if err := ValidateIntentTransferToken(*p.Token); err != nil {
            return err
        }
    }
    if p.TxHash != nil {
        if err := ValidateIntentTransferTxHash(*p.TxHash); err != nil {
            return err
        }
    }
    if p.EventIndex != nil {
        if err := ValidateIntentTransferEventIndex(*p.EventIndex); err != nil {
            return err
        }
    }
    if p.ToAddress != nil {
        if err := ValidateIntentTransferToAddress(*p.ToAddress); err != nil {
            return err
        }
    }

    if p.OrderBy != "" {
        switch p.OrderBy {
        case ColID,
            ColIntentID,
            ColOwnerRef,
            ColChain,
            ColNetwork,
            ColToken,
            ColBlockNumber,
            ColTxHash,
            ColEventIndex,
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


