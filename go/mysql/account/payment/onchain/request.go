//
// request.go
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
	DefaultRequestTableName = "payment_onchain_requests"
)

var (
	requestIDCounter = &helper.IdCounter{}
)

type Request struct {
	ID        uint64        `json:"id,string"`
	AccountID uint64        `json:"accountId,string"`
	Status    RequestStatus `json:"status"`
	Chain     Chain         `json:"chain"`
	Network   Network       `json:"network"`
	Asset     string        `json:"asset"`
	Address   string        `json:"address"`
	Amount    string        `json:"amount"`
	Memo      string        `json:"memo,omitempty"`
	ExpiresAt time.Time     `json:"expiresAt"`
	CreatedAt time.Time     `json:"createdAt,omitempty"`
	UpdatedAt time.Time     `json:"updatedAt,omitempty"`
}

type RequestStore struct {
	db               *sql.DB
	tableName        string
	accountTableName string
}

type RequestSelectOption struct {
	AccountID   *uint64        `json:"accountId,string,omitempty"`
	Status      *RequestStatus `json:"status,omitempty"`
	Chain       *Chain         `json:"chain,omitempty"`
	Network     *Network       `json:"network,omitempty"`
	Asset       *string        `json:"asset,omitempty"`
	Address     *string        `json:"address,omitempty"`
	OrderBy     string         `json:"orderBy"`
	OrderByDesc bool           `json:"orderByDesc"`
	Limit       int            `json:"limit"`
	Offset      int            `json:"offset"`
}

type RequestUpdateOption struct {
	ID        uint64         `json:"id,string"`
	Status    *RequestStatus `json:"status,omitempty"`
	Address   *string        `json:"address,omitempty"`
	Amount    *string        `json:"amount,omitempty"`
	Memo      *string        `json:"memo,omitempty"`
	ExpiresAt *time.Time     `json:"expiresAt,omitempty"`
}

func GenerateRequestID() uint64 {
	return requestIDCounter.GenerateID()
}

func NewRequestStore(db *sql.DB, tableName, accountTableName string) (*RequestStore, error) {
	if db == nil {
		return nil, fmt.Errorf("failed to create payment onchain request store: missing required parameter: db=null")
	}
	if tableName == "" {
		return nil, fmt.Errorf("failed to create payment onchain request store: missing required parameter: table_name=empty")
	}
	if accountTableName == "" {
		return nil, fmt.Errorf("failed to create payment onchain request store: missing required parameter: account_table_name=empty")
	}

	return &RequestStore{
		db:               db,
		tableName:        tableName,
		accountTableName: accountTableName,
	}, nil
}

func (r *Request) ValidateAccountID() error {
	if r == nil {
		return fmt.Errorf("invalid parameter: payment_onchain_request=null")
	}
	if r.AccountID == 0 {
		return fmt.Errorf("invalid parameter: account_id=0")
	}
	return nil
}

func (r *Request) ValidateStatus() error {
	if r == nil {
		return fmt.Errorf("invalid parameter: payment_onchain_request=null")
	}
	if !r.Status.IsValid() {
		return fmt.Errorf("invalid parameter: status=%d", r.Status)
	}
	return nil
}

func (r *Request) ValidateChain() error {
	if r == nil {
		return fmt.Errorf("invalid parameter: payment_onchain_request=null")
	}
	if !r.Chain.IsValid() {
        return fmt.Errorf("invalid parameter: chain=%s", r.Chain)
	}
	return nil
}

func (r *Request) ValidateNetwork() error {
	if r == nil {
		return fmt.Errorf("invalid parameter: payment_onchain_request=null")
	}
	if !r.Network.IsValid() {
        return fmt.Errorf("invalid parameter: network=%s", r.Network)
	}
	return nil
}

func (r *Request) ValidateAsset() error {
	if r == nil {
		return fmt.Errorf("invalid parameter: payment_onchain_request=null")
	}
	if r.Asset == "" {
		return fmt.Errorf("invalid parameter: asset=empty")
	}
	if utf8.RuneCountInString(r.Asset) > 64 {
		return fmt.Errorf("invalid parameter: asset=%q", helper.TruncateRunes(r.Asset, 64))
	}
	return nil
}

func (r *Request) ValidateAddress() error {
	if r == nil {
		return fmt.Errorf("invalid parameter: payment_onchain_request=null")
	}
	if r.Address == "" {
		return fmt.Errorf("invalid parameter: address=empty")
	}
	if utf8.RuneCountInString(r.Address) > 254 {
		return fmt.Errorf("invalid parameter: address=%q", helper.TruncateRunes(r.Address, 254))
	}
	return nil
}

func (r *Request) ValidateAmount() error {
	if r == nil {
		return fmt.Errorf("invalid parameter: payment_onchain_request=null")
	}
	if r.Amount == "" {
		return fmt.Errorf("invalid parameter: amount=empty")
	}
	if utf8.RuneCountInString(r.Amount) > 78 {
		return fmt.Errorf("invalid parameter: amount=%q", helper.TruncateRunes(r.Amount, 78))
	}
	return nil
}

func (r *Request) ValidateMemo() error {
	if r == nil {
		return fmt.Errorf("invalid parameter: payment_onchain_request=null")
	}
	if r.Memo == "" {
		return nil
	}
	if utf8.RuneCountInString(r.Memo) > 254 {
		return fmt.Errorf("invalid parameter: memo=%q", helper.TruncateRunes(r.Memo, 254))
	}
	return nil
}

func (r *Request) ValidateExpiresAt() error {
	if r == nil {
		return fmt.Errorf("invalid parameter: payment_onchain_request=null")
	}
	if r.ExpiresAt.IsZero() {
		return fmt.Errorf("invalid parameter: expires_at=empty")
	}
	return nil
}

func (s *RequestStore) CreateTable() error {
	if s == nil {
		return fmt.Errorf("failed to create payment onchain requests table: missing required parameter: request_store=null")
	}
	if s.db == nil {
		return fmt.Errorf("failed to create payment onchain requests table: missing required parameter: db=null")
	}
	if s.tableName == "" {
		return fmt.Errorf("failed to create payment onchain requests table: missing required parameter: table_name=empty")
	}
	if s.accountTableName == "" {
		return fmt.Errorf("failed to create payment onchain requests table: missing required parameter: account_table_name=empty")
	}

	query := fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s (
			%s BIGINT UNSIGNED NOT NULL COMMENT 'ID',
			%s BIGINT UNSIGNED NOT NULL COMMENT 'Account ID',
			%s TINYINT UNSIGNED NOT NULL COMMENT 'Status',
			%s VARCHAR(64) NOT NULL COMMENT 'Chain',
			%s VARCHAR(64) NOT NULL COMMENT 'Network',
			%s VARCHAR(64) NOT NULL COMMENT 'Asset',
			%s VARCHAR(254) NOT NULL COMMENT 'Deposit address',
			%s VARCHAR(78) NOT NULL COMMENT 'Requested amount',
			%s VARCHAR(254) NULL COMMENT 'Memo',
			%s DATETIME NOT NULL COMMENT 'Expires at',
			%s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
			%s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
			PRIMARY KEY (%s),
			KEY idx_payment_onchain_requests_account_id (%s),
			KEY idx_payment_onchain_requests_status (%s),
			KEY idx_payment_onchain_requests_address (%s),
			KEY idx_payment_onchain_requests_chain_network_asset (%s, %s, %s),
			CONSTRAINT fk_payment_onchain_requests_account_id FOREIGN KEY (%s) REFERENCES %s (%s) ON DELETE CASCADE ON UPDATE CASCADE);
		`,
		s.tableName,
		ColID,
		ColAccountID,
		ColStatus,
		ColChain,
		ColNetwork,
		ColAsset,
		ColAddress,
		ColAmount,
		ColMemo,
		ColExpiresAt,
		ColCreatedAt,
		ColUpdatedAt,
		ColID,
		ColAccountID,
		ColStatus,
		ColAddress,
		ColChain, ColNetwork, ColAsset,
		ColAccountID, s.accountTableName, account.ColID,
	)

	if _, err := s.db.Exec(query); err != nil {
		return fmt.Errorf("failed to create payment onchain requests table: %w", err)
	}

	return nil
}

func (s *RequestStore) Insert(row *Request) error {
	if s == nil {
		return fmt.Errorf("failed to insert payment onchain request: missing required parameter: request_store=null")
	}
	if s.db == nil {
		return fmt.Errorf("failed to insert payment onchain request: missing required parameter: db=null")
	}
	if s.tableName == "" {
		return fmt.Errorf("failed to insert payment onchain request: missing required parameter: table_name=empty")
	}
	if row == nil {
		return fmt.Errorf("failed to insert payment onchain request: missing required parameter: request=null")
	}
	if err := row.ValidateAccountID(); err != nil {
		return fmt.Errorf("failed to insert payment onchain request: %w", err)
	}
	if err := row.ValidateStatus(); err != nil {
		return fmt.Errorf("failed to insert payment onchain request: %w", err)
	}
	if err := row.ValidateChain(); err != nil {
		return fmt.Errorf("failed to insert payment onchain request: %w", err)
	}
	if err := row.ValidateNetwork(); err != nil {
		return fmt.Errorf("failed to insert payment onchain request: %w", err)
	}
	if err := row.ValidateAsset(); err != nil {
		return fmt.Errorf("failed to insert payment onchain request: %w", err)
	}
	if err := row.ValidateAddress(); err != nil {
		return fmt.Errorf("failed to insert payment onchain request: %w", err)
	}
	if err := row.ValidateAmount(); err != nil {
		return fmt.Errorf("failed to insert payment onchain request: %w", err)
	}
	if err := row.ValidateMemo(); err != nil {
		return fmt.Errorf("failed to insert payment onchain request: %w", err)
	}
	if err := row.ValidateExpiresAt(); err != nil {
		return fmt.Errorf("failed to insert payment onchain request: %w", err)
	}

	if row.ID == 0 {
		row.ID = GenerateRequestID()
	}

	now := time.Now()
	if row.CreatedAt.IsZero() {
		row.CreatedAt = now
	}
	if row.UpdatedAt.IsZero() {
		row.UpdatedAt = now
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
		s.tableName,
		ColID,
		ColAccountID,
		ColStatus,
		ColChain,
		ColNetwork,
		ColAsset,
		ColAddress,
		ColAmount,
		ColMemo,
		ColExpiresAt,
		ColCreatedAt,
		ColUpdatedAt,
	)

	if _, err := s.db.Exec(
		query,
		row.ID,
		row.AccountID,
		row.Status,
		row.Chain,
		row.Network,
		row.Asset,
		row.Address,
		row.Amount,
		row.Memo,
		row.ExpiresAt,
		row.CreatedAt,
		row.UpdatedAt,
	); err != nil {
		return fmt.Errorf("failed to insert payment onchain request: %w", err)
	}

	return nil
}

func (s *RequestStore) SelectByID(id uint64) (*Request, error) {
	if s == nil {
		return nil, fmt.Errorf("failed to select payment onchain request by id: missing required parameter: request_store=null")
	}
	if s.db == nil {
		return nil, fmt.Errorf("failed to select payment onchain request by id: missing required parameter: db=null")
	}
	if s.tableName == "" {
		return nil, fmt.Errorf("failed to select payment onchain request by id: missing required parameter: table_name=empty")
	}
	if id == 0 {
		return nil, fmt.Errorf("failed to select payment onchain request by id: invalid parameter: id=0")
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColID)

	row := s.db.QueryRow(query, id)

	result := &Request{}
	err := row.Scan(
		&result.ID,
		&result.AccountID,
		&result.Status,
		&result.Chain,
		&result.Network,
		&result.Asset,
		&result.Address,
		&result.Amount,
		&result.Memo,
		&result.ExpiresAt,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to select payment onchain request by id: %w", err)
	}

	return result, nil
}

func buildRequestSelectCondition(option *RequestSelectOption) ([]string, []any) {
	conditions := make([]string, 0, 6)
	args := make([]any, 0, 6)

	if option != nil {
		if option.AccountID != nil {
			conditions = append(conditions, ColAccountID+" = ?")
			args = append(args, *option.AccountID)
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
		if option.Address != nil {
			conditions = append(conditions, ColAddress+" = ?")
			args = append(args, *option.Address)
		}
	}

	return conditions, args
}

func (s *RequestStore) Count(option *RequestSelectOption) (int64, error) {
	if s == nil {
		return 0, fmt.Errorf("failed to count payment onchain requests: missing required parameter: request_store=null")
	}
	if s.db == nil {
		return 0, fmt.Errorf("failed to count payment onchain requests: missing required parameter: db=null")
	}
	if s.tableName == "" {
		return 0, fmt.Errorf("failed to count payment onchain requests: missing required parameter: table_name=empty")
	}
	if option == nil {
		return 0, fmt.Errorf("failed to count payment onchain requests: missing required parameter: select_option=null")
	}

	conditions, args := buildRequestSelectCondition(option)

	var query strings.Builder
	query.WriteString("SELECT COUNT(*) FROM ")
	query.WriteString(s.tableName)

	if len(conditions) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(conditions, " AND "))
	}

	var result int64
	if err := s.db.QueryRow(query.String(), args...).Scan(&result); err != nil {
		return 0, fmt.Errorf("failed to count payment onchain requests: %w", err)
	}

	return result, nil
}
