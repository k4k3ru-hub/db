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

    "github.com/k4k3ru-hub/db/go/mysql/helper"
)

const (
    DefaultRequestTableName = "payment_onchain_requests"

    DefaultRequestExpiresIn = 30 * time.Minute
)

var (
    requestIDCounter = &helper.IdCounter{}
)


type Request struct {
    ID               uint64        `json:"id,string"`
    AccountID        uint64        `json:"accountId,string"`
    Status           RequestStatus `json:"status"`
    Chain            Chain         `json:"chain"`
    Network          Network       `json:"network"`
    Symbol           Symbol        `json:"symbol"`
    RecipientAddress string        `json:"recipientAddress"`
    Amount           string        `json:"amount"`
    Memo             *string       `json:"memo,omitempty"`
    ExpiresAt        time.Time     `json:"expiresAt"`
    CreatedAt        time.Time     `json:"createdAt,omitempty"`
    UpdatedAt        time.Time     `json:"updatedAt,omitempty"`
}

type RequestStore struct {
    executor  helper.Executor
    tableName string
}

type RequestInsertParams struct {
    ID        uint64        `json:"id,string"`
    AccountID uint64        `json:"accountId,string"`
    Status    RequestStatus `json:"status"`
    Chain     Chain         `json:"chain"`
    Network   Network       `json:"network"`
    Symbol     Symbol         `json:"symbol"`
    RecipientAddress   string        `json:"recipientAddress"`
    Amount    string        `json:"amount"`
    Memo      *string       `json:"memo,omitempty"`
    ExpiresAt time.Time     `json:"expiresAt"`
    CreatedAt time.Time     `json:"createdAt,omitempty"`
    UpdatedAt time.Time     `json:"updatedAt,omitempty"`
    Ignore    bool          `json:"ignore"`
}

type RequestSelectParams struct {
    AccountID   *uint64        `json:"accountId,string,omitempty"`
    Status      *RequestStatus `json:"status,omitempty"`
    Chain       *Chain         `json:"chain,omitempty"`
    Network     *Network       `json:"network,omitempty"`
    Symbol       *Symbol         `json:"symbol,omitempty"`
    RecipientAddress     *string        `json:"recipientAddress,omitempty"`
    OrderBy     string         `json:"orderBy"`
    OrderByDesc bool           `json:"orderByDesc"`
    Limit       int            `json:"limit"`
    Offset      int            `json:"offset"`
}

type RequestUpdateParams struct {
    ID          uint64         `json:"id,string"`
    Status      *RequestStatus `json:"status,omitempty"`
    RecipientAddress     *string        `json:"recipientAddress,omitempty"`
    Amount      *string        `json:"amount,omitempty"`
    Memo        *string        `json:"memo,omitempty"`
    MemoSetNull bool           `json:"memoSetNull"`
    ExpiresAt   *time.Time     `json:"expiresAt,omitempty"`
}


//
// Generate default payment onchain request expires at.
//
// Version:
//   - 2026-05-19: Added.
//
func DefaultExpiresAt() time.Time {
    return time.Now().UTC().Add(DefaultRequestExpiresIn)
}


//
// Generate new payment onchain request ID.
//
// Version:
//   - 2026-05-16: Added.
//
func GenerateRequestID() uint64 {
    return requestIDCounter.GenerateID()
}


//
// Create new payment onchain request store.
//
// Version:
//   - 2026-05-16: Added.
//
func NewRequestStore(executor helper.Executor, tableName string) (*RequestStore, error) {
    // Guard.
    if executor == nil {
        return nil, fmt.Errorf("failed to create payment onchain request store: missing required parameter: executor=null")
    }
    if tableName == "" {
        return nil, fmt.Errorf("failed to create payment onchain request store: missing required parameter: table_name=%q", "empty")
    }

    return &RequestStore{
        executor:  executor,
        tableName: tableName,
    }, nil
}


//
// Validate payment onchain request ID.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateRequestID(id uint64) error {
    if id == 0 {
        return fmt.Errorf("invalid parameter: id=0")
    }
    return nil
}


//
// Validate payment onchain request ID.
//
// Version:
//   - 2026-05-16: Added.
//
func (r *Request) ValidateID() error {
    if r == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_request=null")
    }
    return ValidateRequestID(r.ID)
}


//
// Validate payment onchain request account ID.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateRequestAccountID(accountID uint64) error {
    if accountID == 0 {
        return fmt.Errorf("invalid parameter: account_id=0")
    }
    return nil
}


//
// Validate payment onchain request account ID.
//
// Version:
//   - 2026-05-16: Added.
//
func (r *Request) ValidateAccountID() error {
    if r == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_request=null")
    }
    return ValidateRequestAccountID(r.AccountID)
}


//
// Validate payment onchain request status.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateRequestStatus(s RequestStatus) error {
    if !s.IsValid() {
        return fmt.Errorf("invalid parameter: status=%d", s)
    }
    return nil
}


//
// Validate payment onchain request status.
//
// Version:
//   - 2026-05-16: Added.
//
func (r *Request) ValidateStatus() error {
    if r == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_request=null")
    }
    return ValidateRequestStatus(r.Status)
}


//
// Validate payment onchain request chain.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateRequestChain(c Chain) error {
    if !c.IsValid() {
        return fmt.Errorf("invalid parameter: chain=%q", helper.TruncateRunes(string(c), 32))
    }
    return nil
}


//
// Validate payment onchain request chain.
//
// Version:
//   - 2026-05-16: Added.
//
func (r *Request) ValidateChain() error {
    if r == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_request=null")
    }
    return ValidateRequestChain(r.Chain)
}


//
// Validate payment onchain request network.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateRequestNetwork(n Network) error {
    if !n.IsValid() {
        return fmt.Errorf("invalid parameter: network=%q",  helper.TruncateRunes(string(n), 32))
    }
    return nil
}


//
// Validate payment onchain request network.
//
// Version:
//   - 2026-05-16: Added.
//
func (r *Request) ValidateNetwork() error {
    if r == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_request=null")
    }
    return ValidateRequestNetwork(r.Network)
}


//
// Validate payment onchain request symbol.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateRequestSymbol(s Symbol) error {
    return s.Validate()
}


//
// Validate payment onchain request symbol.
//
// Version:
//   - 2026-05-16: Added.
//
func (r *Request) ValidateSymbol() error {
    if r == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_request=null")
    }
    return ValidateRequestSymbol(r.Symbol)
}


//
// Validate payment onchain request recipient address.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateRequestRecipientAddress(recipientAddress string) error {
    s := strings.TrimSpace(recipientAddress)
    if s == "" {
        return fmt.Errorf("invalid parameter: recipient_address=%q", "empty")
    }
    if len(s) > 255 {
        return fmt.Errorf("invalid parameter: recipient_address=%q", "too long")
    }
    return nil
}


//
// Validate payment onchain request recipient address.
//
// Version:
//   - 2026-05-16: Added.
//
func (r *Request) ValidateRecipientAddress() error {
    if r == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_request=null")
    }
    return ValidateRequestRecipientAddress(r.RecipientAddress)
}


//
// Validate payment onchain request amount.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateRequestAmount(amount string) error {
    if amount == "" {
        return fmt.Errorf("invalid parameter: amount=%q", "empty")
    }
    if utf8.RuneCountInString(amount) > 78 {
        return fmt.Errorf("invalid parameter: amount=%q", helper.TruncateRunes(amount, 78))
    }
    return nil
}


//
// Validate payment onchain request amount.
//
// Version:
//   - 2026-05-16: Added.
//
func (r *Request) ValidateAmount() error {
    if r == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_request=null")
    }
    return ValidateRequestAmount(r.Amount)
}


//
// Validate payment onchain request memo.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateRequestMemo(memo *string) error {
    if memo == nil {
        return nil
    }
    if utf8.RuneCountInString(*memo) > 255 {
        return fmt.Errorf("invalid parameter: memo=%q", helper.TruncateRunes(*memo, 255))
    }
    return nil
}


//
// Validate payment onchain request memo.
//
// Version:
//   - 2026-05-16: Added.
//
func (r *Request) ValidateMemo() error {
    if r == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_request=null")
    }
    return ValidateRequestMemo(r.Memo)
}


//
// Validate payment onchain request expires at.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateRequestExpiresAt(expiresAt time.Time) error {
    if expiresAt.IsZero() {
        return fmt.Errorf("invalid parameter: expires_at=%q", "empty")
    }
    return nil
}


//
// Validate payment onchain request expires at.
//
// Version:
//   - 2026-05-16: Added.
//
func (r *Request) ValidateExpiresAt() error {
    if r == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_request=null")
    }
    return ValidateRequestExpiresAt(r.ExpiresAt)
}


//
// Create payment onchain requests table.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *RequestStore) CreateTable() error {
    if s == nil {
        return fmt.Errorf("failed to create payment onchain requests table: missing required parameter: request_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to create payment onchain requests table: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to create payment onchain requests table: missing required parameter: table_name=%q", "empty")
    }

    query := fmt.Sprintf(
        `CREATE TABLE IF NOT EXISTS %s (
            %s BIGINT UNSIGNED NOT NULL COMMENT 'ID',
            %s BIGINT UNSIGNED NOT NULL COMMENT 'Account ID',
            %s TINYINT UNSIGNED NOT NULL COMMENT 'Status',
            %s VARCHAR(64) NOT NULL COMMENT 'Chain',
            %s VARCHAR(64) NOT NULL COMMENT 'Network',
            %s VARCHAR(64) NOT NULL COMMENT 'Symbol',
            %s VARCHAR(255) NOT NULL COMMENT 'Deposit recipient address',
            %s VARCHAR(78) NOT NULL COMMENT 'Requested payment amount in token smallest units as decimal string',
            %s VARCHAR(255) NULL COMMENT 'Memo',
            %s DATETIME NOT NULL COMMENT 'Expires at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
            PRIMARY KEY (%s),
            KEY idx_payment_onchain_requests_account_id (%s),
            KEY idx_payment_onchain_requests_status (%s),
            KEY idx_payment_onchain_requests_recipient_address (%s),
            KEY idx_payment_onchain_requests_chain_network_symbol (%s, %s, %s)
        );`,
        s.tableName,
        ColID,
        ColAccountID,
        ColStatus,
        ColChain,
        ColNetwork,
        ColSymbol,
        ColRecipientAddress,
        ColAmount,
        ColMemo,
        ColExpiresAt,
        ColCreatedAt,
        ColUpdatedAt,
        ColID,
        ColAccountID,
        ColStatus,
        ColRecipientAddress,
        ColChain, ColNetwork, ColSymbol,
    )

    if _, err := s.executor.Exec(query); err != nil {
        return fmt.Errorf("failed to create payment onchain requests table: %w", err)
    }

    return nil
}


//
// Insert payment onchain request.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *RequestStore) Insert(p *RequestInsertParams) error {
    if s == nil {
        return fmt.Errorf("failed to insert payment onchain request: missing required parameter: request_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to insert payment onchain request: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to insert payment onchain request: missing required parameter: table_name=%q", "empty")
    }
    if p == nil {
        return fmt.Errorf("failed to insert payment onchain request: missing required parameter: request_insert_params=null")
    }
    if err := ValidateRequestAccountID(p.AccountID); err != nil {
        return fmt.Errorf("failed to insert payment onchain request: %w", err)
    }
    if err := ValidateRequestStatus(p.Status); err != nil {
        return fmt.Errorf("failed to insert payment onchain request: %w", err)
    }
    if err := ValidateRequestChain(p.Chain); err != nil {
        return fmt.Errorf("failed to insert payment onchain request: %w", err)
    }
    if err := ValidateRequestNetwork(p.Network); err != nil {
        return fmt.Errorf("failed to insert payment onchain request: %w", err)
    }
    if err := ValidateRequestSymbol(p.Symbol); err != nil {
        return fmt.Errorf("failed to insert payment onchain request: %w", err)
    }
    if err := ValidateRequestRecipientAddress(p.RecipientAddress); err != nil {
        return fmt.Errorf("failed to insert payment onchain request: %w", err)
    }
    if err := ValidateRequestAmount(p.Amount); err != nil {
        return fmt.Errorf("failed to insert payment onchain request: %w", err)
    }
    if err := ValidateRequestMemo(p.Memo); err != nil {
        return fmt.Errorf("failed to insert payment onchain request: %w", err)
    }
    if err := ValidateRequestExpiresAt(p.ExpiresAt); err != nil {
        return fmt.Errorf("failed to insert payment onchain request: %w", err)
    }

    if p.ID == 0 {
        p.ID = GenerateRequestID()
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
        "%s INTO %s (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
        queryPrefix,
        s.tableName,
        ColID,
        ColAccountID,
        ColStatus,
        ColChain,
        ColNetwork,
        ColSymbol,
        ColRecipientAddress,
        ColAmount,
        ColMemo,
        ColExpiresAt,
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
        p.Symbol,
        p.RecipientAddress,
        p.Amount,
        p.Memo,
        p.ExpiresAt,
        p.CreatedAt,
        p.UpdatedAt,
    ); err != nil {
        return fmt.Errorf("failed to insert payment onchain request: %w", err)
    }

    return nil
}


//
// Select payment onchain requests.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *RequestStore) Select(p *RequestSelectParams) ([]*Request, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select payment onchain requests: missing required parameter: request_store=null")
    }
    if s.executor == nil {
        return nil, fmt.Errorf("failed to select payment onchain requests: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select payment onchain requests: missing required parameter: table_name=%q", "empty")
    }
    if err := p.Validate(); err != nil {
        return nil, fmt.Errorf("failed to select payment onchain requests: %w", err)
    }

    query, args := p.BuildQuery("SELECT * FROM " + s.tableName)

    // Execute.
    rows, err := s.executor.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to select payment onchain requests: %w", err)
    }
    defer rows.Close()

    var result []*Request
    for rows.Next() {
        row := &Request{}
        if err := rows.Scan(
            &row.ID,
            &row.AccountID,
            &row.Status,
            &row.Chain,
            &row.Network,
            &row.Symbol,
            &row.RecipientAddress,
            &row.Amount,
            &row.Memo,
            &row.ExpiresAt,
            &row.CreatedAt,
            &row.UpdatedAt,
        ); err != nil {
            return nil, fmt.Errorf("failed to select payment onchain requests: %w", err)
        }

        result = append(result, row)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to select payment onchain requests: %w", err)
    }

    return result, nil
}


//
// Select payment onchain request by ID.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *RequestStore) SelectByID(id uint64) (*Request, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select payment onchain request by id: missing required parameter: request_store=null")
    }
    if s.executor == nil {
        return nil, fmt.Errorf("failed to select payment onchain request by id: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select payment onchain request by id: missing required parameter: table_name=%q", "empty")
    }
    if id == 0 {
        return nil, fmt.Errorf("failed to select payment onchain request by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColID)

    row := s.executor.QueryRow(query, id)

    result := &Request{}
    err := row.Scan(
        &result.ID,
        &result.AccountID,
        &result.Status,
        &result.Chain,
        &result.Network,
        &result.Symbol,
        &result.RecipientAddress,
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


//
// Count payment onchain requests.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *RequestStore) Count(p *RequestSelectParams) (int64, error) {
    if s == nil {
        return 0, fmt.Errorf("failed to count payment onchain requests: missing required parameter: request_store=null")
    }
    if s.executor == nil {
        return 0, fmt.Errorf("failed to count payment onchain requests: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return 0, fmt.Errorf("failed to count payment onchain requests: missing required parameter: table_name=%q", "empty")
    }
    if err := p.Validate(); err != nil {
        return 0, fmt.Errorf("failed to count payment onchain requests: %w", err)
    }

    query, args := p.BuildQuery("SELECT COUNT(*) FROM " + s.tableName)

    var result int64
    if err := s.executor.QueryRow(query, args...).Scan(&result); err != nil {
        return 0, fmt.Errorf("failed to count payment onchain requests: %w", err)
    }

    return result, nil
}


//
// Delete payment onchain request by ID.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *RequestStore) DeleteByID(id uint64) error {
    if s == nil {
        return fmt.Errorf("failed to delete payment onchain request by id: missing required parameter: request_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to delete payment onchain request by id: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to delete payment onchain request by id: missing required parameter: table_name=%q", "empty")
    }
    if id == 0 {
        return fmt.Errorf("failed to delete payment onchain request by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", s.tableName, ColID)

    if _, err := s.executor.Exec(query, id); err != nil {
        return fmt.Errorf("failed to delete payment onchain request by id: %w", err)
    }

    return nil
}


//
// Update payment onchain request status by ID.
//
func (s *RequestStore) UpdateStatusByID(id uint64, status RequestStatus) error {
    if s == nil {
        return fmt.Errorf("failed to update payment onchain request status by id: missing required parameter: request_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to update payment onchain request status by id: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update payment onchain request status by id: missing required parameter: table_name=%q", "empty")
    }
    if id == 0 {
        return fmt.Errorf("failed to update payment onchain request status by id: invalid parameter: id=0")
    }
    if err := status.Validate(); err != nil {
        return fmt.Errorf("failed to update payment onchain request status by id: %w", err)
    }

    query := fmt.Sprintf(
        "UPDATE %s SET %s = ? WHERE %s = ?;",
        s.tableName,
        ColStatus,
        ColID,
    )

    args := make([]any, 0, 2)
    args = append(args, status, id)

    if _, err := s.executor.Exec(query, args...); err != nil {
        return fmt.Errorf("failed to update payment onchain request status by id: %w", err)
    }

    return nil
}


//
// Build query.
//
// Version:
//   - 2025-05-16: Added.
//
func (p *RequestSelectParams) BuildQuery(selectFromClause string) (string, []any) {
    // Guard.
    if p == nil {
        return selectFromClause, nil
    }

    var query strings.Builder
    query.WriteString(selectFromClause)

    conditions := make([]string, 0, 6)
    args := make([]any, 0, 8)

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
    if p.Symbol != nil {
        conditions = append(conditions, ColSymbol + " = ?")
        args = append(args, *p.Symbol)
    }
    if p.RecipientAddress != nil {
        conditions = append(conditions, ColRecipientAddress + " = ?")
        args = append(args, *p.RecipientAddress)
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
func (p *RequestSelectParams) Validate() error {
    // Guard.
    if p == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_request_select_params=null")
    }

    if p.AccountID != nil {
        if err := ValidateRequestAccountID(*p.AccountID); err != nil {
            return err
        }
    }
    if p.Status != nil {
        if err := ValidateRequestStatus(*p.Status); err != nil {
            return err
        }
    }
    if p.Chain != nil {
        if err := ValidateRequestChain(*p.Chain); err != nil {
            return err
        }
    }
    if p.Network != nil {
        if err := ValidateRequestNetwork(*p.Network); err != nil {
            return err
        }
    }
    if p.Symbol != nil {
        if err := ValidateRequestSymbol(*p.Symbol); err != nil {
            return err
        }
    }
    if p.RecipientAddress != nil {
        if err := ValidateRequestRecipientAddress(*p.RecipientAddress); err != nil {
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
            ColSymbol,
            ColRecipientAddress,
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

