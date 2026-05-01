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

    "github.com/k4k3ru-hub/db/go/mysql/account"
    "github.com/k4k3ru-hub/db/go/mysql/helper"
)

const (
    DefaultTxTableName = "account_payment_onchain_txs"
)

var (
    txIDCounter = &helper.IdCounter{}
)

type Tx struct {
    ID            uint64     `json:"id,string"`
    AccountID     uint64     `json:"accountId,string"`
    RequestID     *uint64    `json:"requestId,string,omitempty"`
    Status        TxStatus   `json:"status"`
    Chain         Chain      `json:"chain"`
    Network       Network    `json:"network"`
    Asset         string     `json:"asset"`
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
    db               *sql.DB
    tableName        string
    accountTableName string
    requestTableName string
}

type TxSelectOption struct {
    AccountID   *uint64   `json:"accountId,string,omitempty"`
    RequestID   *uint64   `json:"requestId,string,omitempty"`
    Status      *TxStatus `json:"status,omitempty"`
    Chain       *Chain    `json:"chain,omitempty"`
    Network     *Network  `json:"network,omitempty"`
    Asset       *string   `json:"asset,omitempty"`
    TxHash      *string   `json:"txHash,omitempty"`
    ToAddress   *string   `json:"toAddress,omitempty"`
    OrderBy     string    `json:"orderBy"`
    OrderByDesc bool      `json:"orderByDesc"`
    Limit       int       `json:"limit"`
    Offset      int       `json:"offset"`
}

type TxUpdateOption struct {
    ID            uint64     `json:"id,string"`
    RequestID     *uint64    `json:"requestId,string,omitempty"`
    Status        *TxStatus  `json:"status,omitempty"`
    Confirmations *uint64    `json:"confirmations,string,omitempty"`
    ConfirmedAt   *time.Time `json:"confirmedAt,omitempty"`
}

func GenerateTxID() uint64 {
    return txIDCounter.GenerateID()
}

func NewTxStore(db *sql.DB, tableName, accountTableName, requestTableName string) (*TxStore, error) {
    if db == nil {
        return nil, fmt.Errorf("failed to create account payment onchain tx store: missing required parameter: db=null")
    }
    if tableName == "" {
        return nil, fmt.Errorf("failed to create account payment onchain tx store: missing required parameter: table_name=empty")
    }
    if accountTableName == "" {
        return nil, fmt.Errorf("failed to create account payment onchain tx store: missing required parameter: account_table_name=empty")
    }
    if requestTableName == "" {
        return nil, fmt.Errorf("failed to create account payment onchain tx store: missing required parameter: request_table_name=empty")
    }

    return &TxStore{
        db:        db,
        tableName: tableName,
        accountTableName: accountTableName,
        requestTableName: requestTableName,
    }, nil
}

func (t *Tx) ValidateRequestID() error {
    if t == nil {
        return fmt.Errorf("invalid parameter: account_payment_onchain_tx=null")
    }
    if t.RequestID == nil {
        return nil
    }
    if *t.RequestID == 0 {
        return fmt.Errorf("invalid parameter: request_id=0")
    }
    return nil
}

func (t *Tx) ValidateAccountID() error {
    if t == nil {
        return fmt.Errorf("invalid parameter: account_payment_onchain_tx=null")
    }
    if t.AccountID == 0 {
        return fmt.Errorf("invalid parameter: account_id=0")
    }
    return nil
}

func (t *Tx) ValidateStatus() error {
    if t == nil {
        return fmt.Errorf("invalid parameter: account_payment_onchain_tx=null")
    }
    if !t.Status.IsValid() {
        return fmt.Errorf("invalid parameter: status=%d", t.Status)
    }
    return nil
}

func (t *Tx) ValidateChain() error {
    if t == nil {
        return fmt.Errorf("invalid parameter: account_payment_onchain_tx=null")
    }
    if !t.Chain.IsValid() {
		return fmt.Errorf("invalid parameter: chain=%s", t.Chain)
	}
    return nil
}

func (t *Tx) ValidateNetwork() error {
    if t == nil {
        return fmt.Errorf("invalid parameter: account_payment_onchain_tx=null")
    }
    if !t.Network.IsValid() {
        return fmt.Errorf("invalid parameter: network=%s", t.Network)
    }
    return nil
}

func (t *Tx) ValidateAsset() error {
    if t == nil {
        return fmt.Errorf("invalid parameter: account_payment_onchain_tx=null")
    }
    if t.Asset == "" {
        return fmt.Errorf("invalid parameter: asset=empty")
    }
    if utf8.RuneCountInString(t.Asset) > 64 {
        return fmt.Errorf("invalid parameter: asset=%q", helper.TruncateRunes(t.Asset, 64))
    }
    return nil
}

func (t *Tx) ValidateTxHash() error {
    if t == nil {
        return fmt.Errorf("invalid parameter: account_payment_onchain_tx=null")
    }
    if t.TxHash == "" {
        return fmt.Errorf("invalid parameter: tx_hash=empty")
    }
    if utf8.RuneCountInString(t.TxHash) > 254 {
        return fmt.Errorf("invalid parameter: tx_hash=%q", helper.TruncateRunes(t.TxHash, 254))
    }
    return nil
}

func (t *Tx) ValidateFromAddress() error {
    if t == nil {
        return fmt.Errorf("invalid parameter: account_payment_onchain_tx=null")
    }
    if t.FromAddress == "" {
        return nil
    }
    if utf8.RuneCountInString(t.FromAddress) > 254 {
        return fmt.Errorf("invalid parameter: from_address=%q", helper.TruncateRunes(t.FromAddress, 254))
    }
    return nil
}

func (t *Tx) ValidateToAddress() error {
    if t == nil {
        return fmt.Errorf("invalid parameter: account_payment_onchain_tx=null")
    }
    if t.ToAddress == "" {
        return fmt.Errorf("invalid parameter: to_address=empty")
    }
    if utf8.RuneCountInString(t.ToAddress) > 254 {
        return fmt.Errorf("invalid parameter: to_address=%q", helper.TruncateRunes(t.ToAddress, 254))
    }
    return nil
}

func (t *Tx) ValidateAmount() error {
    if t == nil {
        return fmt.Errorf("invalid parameter: account_payment_onchain_tx=null")
    }
    if t.Amount == "" {
        return fmt.Errorf("invalid parameter: amount=empty")
    }
    if utf8.RuneCountInString(t.Amount) > 78 {
        return fmt.Errorf("invalid parameter: amount=%q", helper.TruncateRunes(t.Amount, 78))
    }
    return nil
}

func (t *Tx) ValidateDetectedAt() error {
    if t == nil {
        return fmt.Errorf("invalid parameter: account_payment_onchain_tx=null")
    }
    if t.DetectedAt.IsZero() {
        return fmt.Errorf("invalid parameter: detected_at=empty")
    }
    return nil
}

func (t *Tx) ValidateConfirmedAt() error {
    if t == nil {
        return fmt.Errorf("invalid parameter: account_payment_onchain_tx=null")
    }
    if t.ConfirmedAt == nil {
        return nil
    }
    if t.ConfirmedAt.IsZero() {
        return fmt.Errorf("invalid parameter: confirmed_at=empty")
    }
    return nil
}

func (s *TxStore) CreateTable() error {
    if s == nil {
        return fmt.Errorf("failed to create account payment onchain txs table: missing required parameter: tx_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to create account payment onchain txs table: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to create account payment onchain txs table: missing required parameter: table_name=empty")
    }
    if s.accountTableName == "" {
        return fmt.Errorf("failed to create account payment onchain txs table: missing required parameter: account_table_name=empty")
    }
    if s.requestTableName == "" {
        return fmt.Errorf("failed to create account payment onchain txs table: missing required parameter: request_table_name=empty")
    }

    query := fmt.Sprintf(
        `CREATE TABLE IF NOT EXISTS %s (
            %s BIGINT UNSIGNED NOT NULL COMMENT 'ID',
            %s BIGINT UNSIGNED NOT NULL COMMENT 'Account ID',
            %s BIGINT UNSIGNED NULL COMMENT 'Request ID',
            %s TINYINT UNSIGNED NOT NULL COMMENT 'Status',
            %s VARCHAR(64) NOT NULL COMMENT 'Chain',
            %s VARCHAR(64) NOT NULL COMMENT 'Network',
            %s VARCHAR(64) NOT NULL COMMENT 'Asset',
            %s VARCHAR(254) NOT NULL COMMENT 'Transaction hash',
            %s VARCHAR(254) NULL COMMENT 'From address',
            %s VARCHAR(254) NOT NULL COMMENT 'To address',
            %s VARCHAR(78) NOT NULL COMMENT 'Amount',
            %s BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Confirmations',
            %s DATETIME NOT NULL COMMENT 'Detected at',
            %s DATETIME NULL COMMENT 'Confirmed at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
            PRIMARY KEY (%s),
            UNIQUE KEY uk_account_payment_onchain_txs_chain_network_tx_hash (%s, %s, %s),
            KEY idx_account_payment_onchain_txs_request_id (%s),
            KEY idx_account_payment_onchain_txs_account_id (%s),
            KEY idx_account_payment_onchain_txs_status (%s),
            KEY idx_account_payment_onchain_txs_to_address (%s),
            KEY idx_account_payment_onchain_txs_chain_network_asset (%s, %s, %s),
            CONSTRAINT fk_account_payment_onchain_txs_account_id FOREIGN KEY (%s) REFERENCES %s (%s) ON DELETE CASCADE ON UPDATE CASCADE),
            CONSTRAINT fk_account_payment_onchain_txs_request_id FOREIGN KEY (%s) REFERENCES %s (%s) ON DELETE SET NULL ON UPDATE CASCADE;
        `,
        s.tableName,
        ColID,
        ColAccountID,
        ColRequestID,
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
        ColAccountID, s.accountTableName, account.ColID,
        ColRequestID, s.requestTableName, ColID,
    )

    if _, err := s.db.Exec(query); err != nil {
        return fmt.Errorf("failed to create account payment onchain txs table: %w", err)
    }

    return nil
}

func (s *TxStore) Insert(row *Tx) error {
    if s == nil {
        return fmt.Errorf("failed to insert account payment onchain tx: missing required parameter: tx_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to insert account payment onchain tx: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to insert account payment onchain tx: missing required parameter: table_name=empty")
    }
    if row == nil {
        return fmt.Errorf("failed to insert account payment onchain tx: missing required parameter: tx=null")
    }
    if err := row.ValidateAccountID(); err != nil {
        return fmt.Errorf("failed to insert account payment onchain tx: %w", err)
    }
    if err := row.ValidateRequestID(); err != nil {
        return fmt.Errorf("failed to insert account payment onchain tx: %w", err)
    }
    if err := row.ValidateStatus(); err != nil {
        return fmt.Errorf("failed to insert account payment onchain tx: %w", err)
    }
    if err := row.ValidateChain(); err != nil {
        return fmt.Errorf("failed to insert account payment onchain tx: %w", err)
    }
    if err := row.ValidateNetwork(); err != nil {
        return fmt.Errorf("failed to insert account payment onchain tx: %w", err)
    }
    if err := row.ValidateAsset(); err != nil {
        return fmt.Errorf("failed to insert account payment onchain tx: %w", err)
    }
    if err := row.ValidateTxHash(); err != nil {
        return fmt.Errorf("failed to insert account payment onchain tx: %w", err)
    }
    if err := row.ValidateFromAddress(); err != nil {
        return fmt.Errorf("failed to insert account payment onchain tx: %w", err)
    }
    if err := row.ValidateToAddress(); err != nil {
        return fmt.Errorf("failed to insert account payment onchain tx: %w", err)
    }
    if err := row.ValidateAmount(); err != nil {
        return fmt.Errorf("failed to insert account payment onchain tx: %w", err)
    }
    if err := row.ValidateDetectedAt(); err != nil {
        return fmt.Errorf("failed to insert account payment onchain tx: %w", err)
    }
    if err := row.ValidateConfirmedAt(); err != nil {
        return fmt.Errorf("failed to insert account payment onchain tx: %w", err)
    }

    if row.ID == 0 {
        row.ID = GenerateTxID()
    }

    now := time.Now()
    if row.CreatedAt.IsZero() {
        row.CreatedAt = now
    }
    if row.UpdatedAt.IsZero() {
        row.UpdatedAt = now
    }

    query := fmt.Sprintf(
        "INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
        s.tableName,
        ColID,
        ColAccountID,
        ColRequestID,
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

    if _, err := s.db.Exec(
        query,
        row.ID,
        row.AccountID,
        row.RequestID,
        row.Status,
        row.Chain,
        row.Network,
        row.Asset,
        row.TxHash,
        row.FromAddress,
        row.ToAddress,
        row.Amount,
        row.Confirmations,
        row.DetectedAt,
        row.ConfirmedAt,
        row.CreatedAt,
        row.UpdatedAt,
    ); err != nil {
        return fmt.Errorf("failed to insert account payment onchain tx: %w", err)
    }

    return nil
}

func (s *TxStore) SelectByID(id uint64) (*Tx, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select account payment onchain tx by id: missing required parameter: tx_store=null")
    }
    if s.db == nil {
        return nil, fmt.Errorf("failed to select account payment onchain tx by id: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account payment onchain tx by id: missing required parameter: table_name=empty")
    }
    if id == 0 {
        return nil, fmt.Errorf("failed to select account payment onchain tx by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColID)

    result := &Tx{}
    err := s.db.QueryRow(query, id).Scan(
        &result.ID,
        &result.AccountID,
        &result.RequestID,
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
        return nil, fmt.Errorf("failed to select account payment onchain tx by id: %w", err)
    }

    return result, nil
}

func (s *TxStore) Select(option *TxSelectOption) ([]*Tx, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select account payment onchain txs: missing required parameter: tx_store=null")
    }
    if s.db == nil {
        return nil, fmt.Errorf("failed to select account payment onchain txs: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account payment onchain txs: missing required parameter: table_name=empty")
    }
    if option == nil {
        return nil, fmt.Errorf("failed to select account payment onchain txs: missing required parameter: select_option=null")
    }

    conditions, args := buildTxSelectCondition(option)

    var query strings.Builder
    query.WriteString("SELECT * FROM ")
    query.WriteString(s.tableName)

    if len(conditions) > 0 {
        query.WriteString(" WHERE ")
        query.WriteString(strings.Join(conditions, " AND "))
    }

    if option.OrderBy != "" {
        query.WriteString(" ORDER BY ")
        query.WriteString(option.OrderBy)
        if option.OrderByDesc {
            query.WriteString(" DESC")
        }
    }

    if option.Limit > 0 {
        if option.Offset < 0 {
            return nil, fmt.Errorf("failed to select account payment onchain txs: invalid parameter: offset=%d", option.Offset)
        }
        query.WriteString(" LIMIT ? OFFSET ?")
        args = append(args, option.Limit, option.Offset)
    }

    rows, err := s.db.Query(query.String(), args...)
    if err != nil {
        return nil, fmt.Errorf("failed to select account payment onchain txs: %w", err)
    }
    defer rows.Close()

    var result []*Tx
    for rows.Next() {
        row := &Tx{}
        err := rows.Scan(
            &row.ID,
            &row.AccountID,
            &row.RequestID,
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
            return nil, fmt.Errorf("failed to select account payment onchain txs: %w", err)
        }

        result = append(result, row)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to select account payment onchain txs: %w", err)
    }

    return result, nil
}

func (s *TxStore) Count(option *TxSelectOption) (int64, error) {
    if s == nil {
        return 0, fmt.Errorf("failed to count account payment onchain txs: missing required parameter: tx_store=null")
    }
    if s.db == nil {
        return 0, fmt.Errorf("failed to count account payment onchain txs: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return 0, fmt.Errorf("failed to count account payment onchain txs: missing required parameter: table_name=empty")
    }
    if option == nil {
        return 0, fmt.Errorf("failed to count account payment onchain txs: missing required parameter: select_option=null")
    }

    conditions, args := buildTxSelectCondition(option)

    var query strings.Builder
    query.WriteString("SELECT COUNT(*) FROM ")
    query.WriteString(s.tableName)

    if len(conditions) > 0 {
        query.WriteString(" WHERE ")
        query.WriteString(strings.Join(conditions, " AND "))
    }

    var result int64
    if err := s.db.QueryRow(query.String(), args...).Scan(&result); err != nil {
        return 0, fmt.Errorf("failed to count account payment onchain txs: %w", err)
    }

    return result, nil
}

func (s *TxStore) DeleteByID(id uint64) error {
    if s == nil {
        return fmt.Errorf("failed to delete account payment onchain tx by id: missing required parameter: tx_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to delete account payment onchain tx by id: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to delete account payment onchain tx by id: missing required parameter: table_name=empty")
    }
    if id == 0 {
        return fmt.Errorf("failed to delete account payment onchain tx by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", s.tableName, ColID)

    if _, err := s.db.Exec(query, id); err != nil {
        return fmt.Errorf("failed to delete account payment onchain tx by id: %w", err)
    }

    return nil
}

func (s *TxStore) Update(option *TxUpdateOption) error {
    if s == nil {
        return fmt.Errorf("failed to update account payment onchain tx: missing required parameter: tx_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to update account payment onchain tx: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update account payment onchain tx: missing required parameter: table_name=empty")
    }
    if option == nil {
        return fmt.Errorf("failed to update account payment onchain tx: missing required parameter: option=null")
    }
    if option.ID == 0 {
        return fmt.Errorf("failed to update account payment onchain tx: invalid parameter: id=0")
    }

    assignments := make([]string, 0, 4)
    args := make([]any, 0, 4)

    if option.RequestID != nil {
        if *option.RequestID == 0 {
            return fmt.Errorf("failed to update account payment onchain tx: invalid parameter: request_id=0")
        }
        assignments = append(assignments, ColRequestID+" = ?")
        args = append(args, *option.RequestID)
    }

    if option.Status != nil {
        t := Tx{Status: *option.Status}
        if err := t.ValidateStatus(); err != nil {
            return fmt.Errorf("failed to update account payment onchain tx: %w", err)
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
            return fmt.Errorf("failed to update account payment onchain tx: %w", err)
        }
        assignments = append(assignments, ColConfirmedAt+" = ?")
        args = append(args, *option.ConfirmedAt)
    }

    if len(assignments) == 0 {
        return fmt.Errorf("failed to update account payment onchain tx: invalid parameter: assignments=empty")
    }

    args = append(args, option.ID)

    query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?;", s.tableName, strings.Join(assignments, ", "), ColID)

    if _, err := s.db.Exec(query, args...); err != nil {
        return fmt.Errorf("failed to update account payment onchain tx: %w", err)
    }

    return nil
}

func (s *TxStore) UpdateConfirmed(id uint64, confirmations uint64) error {
    if id == 0 {
        return fmt.Errorf("failed to update account payment onchain tx confirmed: invalid parameter: id=0")
    }

    status := TxStatusConfirmed
    now := time.Now()

    return s.Update(&TxUpdateOption{
        ID:            id,
        Status:        &status,
        Confirmations: &confirmations,
        ConfirmedAt:   &now,
    })
}

func buildTxSelectCondition(option *TxSelectOption) ([]string, []any) {
    conditions := make([]string, 0, 8)
    args := make([]any, 0, 8)

    if option != nil {
        if option.AccountID != nil {
            conditions = append(conditions, ColAccountID+" = ?")
            args = append(args, *option.AccountID)
        }
        if option.RequestID != nil {
            conditions = append(conditions, ColRequestID+" = ?")
            args = append(args, *option.RequestID)
        }
        if option.Status != nil {
            conditions = append(conditions, ColStatus+" = ?")
            args = append(args, *option.Status)
        }
        if option.Chain != nil {
            conditions = append(conditions, ColChain+" = ?")
            args = append(args, *option.Chain)
        }
        if option.Network != nil {
            conditions = append(conditions, ColNetwork+" = ?")
            args = append(args, *option.Network)
        }
        if option.Asset != nil {
            conditions = append(conditions, ColAsset+" = ?")
            args = append(args, *option.Asset)
        }
        if option.TxHash != nil {
            conditions = append(conditions, ColTxHash+" = ?")
            args = append(args, *option.TxHash)
        }
        if option.ToAddress != nil {
            conditions = append(conditions, ColToAddress+" = ?")
            args = append(args, *option.ToAddress)
        }
    }

    return conditions, args
}
