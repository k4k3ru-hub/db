//
// intent.go
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
    DefaultIntentTableName = "payment_onchain_intents"

    DefaultIntentExpiresIn = 30 * time.Minute
)

var (
    intentIDCounter = &helper.IdCounter{}
)


type Intent struct {
    ID                uint64        `json:"id,string"`
    AccountID         uint64        `json:"accountId,string"`
    Status            IntentStatus  `json:"status"`
    Chain             Chain         `json:"chain"`
    Network           Network       `json:"network"`
    Token             Token         `json:"token"`
    RecipientID       uint64        `json:"recipientId,string"`
    RecipientKind     RecipientKind `json:"recipientKind"`
    RecipientAddress  string        `json:"recipientAddress"`
    Amount            string        `json:"amount"`
    Memo              *string       `json:"memo,omitempty"`
    ExpiresAt         time.Time     `json:"expiresAt"`
    WebhookEndpointID *uint64       `json:"webhookEndpointId,omitempty"`
    Metadata          helper.JSON   `json:"metadata,omitempty"`
    CreatedAt         time.Time     `json:"createdAt,omitempty"`
    UpdatedAt         time.Time     `json:"updatedAt,omitempty"`
}

type IntentStore struct {
    tableName string
}

type IntentInsertParams struct {
    ID                uint64        `json:"id,string"`
    AccountID         uint64        `json:"accountId,string"`
    Status            IntentStatus  `json:"status"`
    Chain             Chain         `json:"chain"`
    Network           Network       `json:"network"`
    Token             Token         `json:"token"`
    RecipientID       uint64        `json:"recipientId,string"`
    RecipientKind     RecipientKind `json:"recipientKind"`
    RecipientAddress  string        `json:"recipientAddress"`
    Amount            string        `json:"amount"`
    Memo              *string       `json:"memo,omitempty"`
    ExpiresAt         time.Time     `json:"expiresAt"`
    WebhookEndpointID *uint64       `json:"webhookEndpointId,omitempty"`
    Metadata          helper.JSON   `json:"metadata,omitempty"`
    CreatedAt         time.Time     `json:"createdAt,omitempty"`
    UpdatedAt         time.Time     `json:"updatedAt,omitempty"`
    Ignore            bool          `json:"ignore"`
}

type IntentSelectParams struct {
    AccountID        *uint64        `json:"accountId,string,omitempty"`
    Status           *IntentStatus  `json:"status,omitempty"`
    Chain            *Chain         `json:"chain,omitempty"`
    Network          *Network       `json:"network,omitempty"`
    Token            *Token         `json:"token,omitempty"`
    RecipientID      *uint64        `json:"recipientId,string,omitempty"`
    RecipientAddress *string        `json:"recipientAddress,omitempty"`
    OrderBy          string         `json:"orderBy"`
    OrderByDesc      bool           `json:"orderByDesc"`
    Limit            int            `json:"limit"`
    Offset           int            `json:"offset"`
}

type IntentUpdateParams struct {
    ID               uint64         `json:"id,string"`
    Status           *IntentStatus  `json:"status,omitempty"`
    RecipientAddress *string        `json:"recipientAddress,omitempty"`
    Amount           *string        `json:"amount,omitempty"`
    Memo             *string        `json:"memo,omitempty"`
    MemoSetNull      bool           `json:"memoSetNull"`
    ExpiresAt        *time.Time     `json:"expiresAt,omitempty"`
}


//
// Generate default payment onchain intent expires at.
//
// Version:
//   - 2026-05-19: Added.
//
func DefaultExpiresAt() time.Time {
    return time.Now().UTC().Add(DefaultIntentExpiresIn)
}


//
// Generate new payment onchain intent ID.
//
// Version:
//   - 2026-05-16: Added.
//
func GenerateIntentID() uint64 {
    return intentIDCounter.GenerateID()
}


//
// Create new payment onchain intent store.
//
// Version:
//   - 2026-05-16: Added.
//
func NewIntentStore(tableName string) (*IntentStore, error) {
    // Guard.
    tableName = strings.TrimSpace(tableName)
    if tableName == "" {
        return nil, fmt.Errorf("failed to create new payment onchain intent store: missing required parameter: table_name=%q", "empty")
    }

    return &IntentStore{
        tableName: tableName,
    }, nil
}


//
// Validate payment onchain intent ID.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateIntentID(id uint64) error {
    if id == 0 {
        return fmt.Errorf("invalid parameter: id=0")
    }
    return nil
}


//
// Validate payment onchain intent ID.
//
// Version:
//   - 2026-05-16: Added.
//
func (i *Intent) ValidateID() error {
    if i == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_intent=null")
    }
    return ValidateIntentID(i.ID)
}


//
// Validate payment onchain intent account ID.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateIntentAccountID(accountID uint64) error {
    if accountID == 0 {
        return fmt.Errorf("invalid parameter: account_id=0")
    }
    return nil
}


//
// Validate payment onchain intent account ID.
//
// Version:
//   - 2026-05-16: Added.
//
func (i *Intent) ValidateAccountID() error {
    if i == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_intent=null")
    }
    return ValidateIntentAccountID(i.AccountID)
}


//
// Validate payment onchain intent status.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateIntentStatus(s IntentStatus) error {
    if !s.IsValid() {
        return fmt.Errorf("invalid parameter: status=%d", s)
    }
    return nil
}


//
// Validate payment onchain intent status.
//
// Version:
//   - 2026-05-16: Added.
//
func (i *Intent) ValidateStatus() error {
    if i == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_intent=null")
    }
    return ValidateIntentStatus(i.Status)
}


//
// Validate payment onchain intent chain.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateIntentChain(c Chain) error {
    if !c.IsValid() {
        return fmt.Errorf("invalid parameter: chain=%q", helper.TruncateRunes(string(c), 32))
    }
    return nil
}


//
// Validate payment onchain intent chain.
//
// Version:
//   - 2026-05-16: Added.
//
func (i *Intent) ValidateChain() error {
    if i == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_intent=null")
    }
    return ValidateIntentChain(i.Chain)
}


//
// Validate payment onchain intent network.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateIntentNetwork(n Network) error {
    if !n.IsValid() {
        return fmt.Errorf("invalid parameter: network=%q",  helper.TruncateRunes(string(n), 32))
    }
    return nil
}


//
// Validate payment onchain intent network.
//
// Version:
//   - 2026-05-16: Added.
//
func (i *Intent) ValidateNetwork() error {
    if i == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_intent=null")
    }
    return ValidateIntentNetwork(i.Network)
}


//
// Validate payment onchain intent token.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateIntentToken(t Token) error {
    return t.Validate()
}


//
// Validate payment onchain intent token.
//
// Version:
//   - 2026-05-16: Added.
//
func (i *Intent) ValidateToken() error {
    if i == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_intent=null")
    }
    return ValidateIntentToken(i.Token)
}


//
// Validate payment onchain intent recipient ID.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateIntentRecipientID(recipientID uint64) error {
    if recipientID == 0 {
        return fmt.Errorf("invalid parameter: recipient_id=0")
    }
    return nil
}


//
// Validate payment onchain intent recipient ID.
//
// Version:
//   - 2026-05-16: Added.
//
func (i *Intent) ValidateRecipientID() error {
    if i == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_intent=null")
    }
    return ValidateIntentRecipientID(i.RecipientID)
}


//
// Validate payment onchain intent recipient kind.
//
// Version:
//   - 2026-05-31: Added.
//
func ValidateIntentRecipientKind(kind RecipientKind) error {
    return kind.Validate()
}


//
// Validate payment onchain intent recipient kind.
//
// Version:
//   - 2026-05-31: Added.
//
func (i *Intent) ValidateRecipientKind() error {
    if i == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_intent=null")
    }
    return ValidateIntentRecipientKind(i.RecipientKind)
}


//
// Validate payment onchain intent recipient address.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateIntentRecipientAddress(recipientAddress string) error {
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
// Validate payment onchain intent recipient address.
//
// Version:
//   - 2026-05-16: Added.
//
func (r *Intent) ValidateRecipientAddress() error {
    if r == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_intent=null")
    }
    return ValidateIntentRecipientAddress(r.RecipientAddress)
}


//
// Validate payment onchain intent amount.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateIntentAmount(amount string) error {
    if amount == "" {
        return fmt.Errorf("invalid parameter: amount=%q", "empty")
    }
    if utf8.RuneCountInString(amount) > 78 {
        return fmt.Errorf("invalid parameter: amount=%q", helper.TruncateRunes(amount, 78))
    }
    return nil
}


//
// Validate payment onchain intent amount.
//
// Version:
//   - 2026-05-16: Added.
//
func (r *Intent) ValidateAmount() error {
    if r == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_intent=null")
    }
    return ValidateIntentAmount(r.Amount)
}


//
// Validate payment onchain intent memo.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateIntentMemo(memo *string) error {
    if memo == nil {
        return nil
    }
    if utf8.RuneCountInString(*memo) > 255 {
        return fmt.Errorf("invalid parameter: memo=%q", helper.TruncateRunes(*memo, 255))
    }
    return nil
}


//
// Validate payment onchain intent memo.
//
// Version:
//   - 2026-05-16: Added.
//
func (r *Intent) ValidateMemo() error {
    if r == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_intent=null")
    }
    return ValidateIntentMemo(r.Memo)
}


//
// Validate payment onchain intent expires at.
//
// Version:
//   - 2026-05-16: Added.
//
func ValidateIntentExpiresAt(expiresAt time.Time) error {
    if expiresAt.IsZero() {
        return fmt.Errorf("invalid parameter: expires_at=%q", "empty")
    }
    return nil
}


//
// Validate payment onchain intent expires at.
//
// Version:
//   - 2026-05-16: Added.
//
func (r *Intent) ValidateExpiresAt() error {
    if r == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_intent=null")
    }
    return ValidateIntentExpiresAt(r.ExpiresAt)
}


//
// Validate payment onchain intent webhook endpoint ID.
//
// Version:
//   - 2026-05-31: Added.
//
func ValidateIntentWebhookEndpointID(webhookEndpointID *uint64) error {
    if webhookEndpointID == nil {
        return nil
    }
    if *webhookEndpointID == 0 {
        return fmt.Errorf("invalid parameter: webhook_endpoint_id=0")
    }
    return nil
}


//
// Validate payment onchain intent webhook endpoint ID.
//
// Version:
//   - 2026-05-31: Added.
//
func (r *Intent) ValidateWebhookEndpointID() error {
    if r == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_intent=null")
    }
    return ValidateIntentWebhookEndpointID(r.WebhookEndpointID)
}


//
// Validate payment onchain intent metadata.
//
// Version:
//   - 2026-05-31: Added.
//
func ValidateIntentMetadata(metadata helper.JSON) error {
    if metadata == nil {
        return nil
    }

    if err := metadata.Validate(); err != nil {
        return  fmt.Errorf("invalid parameter: metadata: %w", err)
    }

    return nil
}


//
// Validate payment onchain intent metadata.
//
// Version:
//   - 2026-05-31: Added.
//
func (r *Intent) ValidateMetadata() error {
    if r == nil {
        return fmt.Errorf("missing required parameter: payment_onchain_intent=null")
    }
    return ValidateIntentMetadata(r.Metadata)
}


//
// Create payment onchain intents table.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *IntentStore) CreateTable(executor helper.Executor) error {
    if s == nil {
        return fmt.Errorf("failed to create payment onchain intents table: missing required parameter: request_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to create payment onchain intents table: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to create payment onchain intents table: missing required parameter: executor=null")
    }

    query := fmt.Sprintf(
        `CREATE TABLE IF NOT EXISTS %s (
            %s BIGINT UNSIGNED NOT NULL COMMENT 'ID',
            %s BIGINT UNSIGNED NOT NULL COMMENT 'Account ID',
            %s TINYINT UNSIGNED NOT NULL COMMENT 'Status',
            %s VARCHAR(64) NOT NULL COMMENT 'Chain',
            %s VARCHAR(64) NOT NULL COMMENT 'Network',
            %s VARCHAR(64) NOT NULL COMMENT 'Token',
            %s BIGINT UNSIGNED NOT NULL COMMENT 'Recipient ID',
            %s VARCHAR(16) NOT NULL COMMENT 'Recipient Kind',
            %s VARCHAR(255) NOT NULL COMMENT 'Recipient address',
            %s VARCHAR(78) NOT NULL COMMENT 'Intended payment amount in token smallest units as decimal string',
            %s VARCHAR(255) NULL COMMENT 'Memo',
            %s DATETIME NOT NULL COMMENT 'Expires at',
            %s BIGINT UNSIGNED NULL COMMENT 'Webhook Endpoint ID',
            %s JSON NULL COMMENT 'Metadata JSON',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
            PRIMARY KEY (%s),
            KEY idx_payment_onchain_intents_account_id_status (%s, %s),
            KEY idx_payment_onchain_intents_rec_add_cha_net_tok_sta (%s, %s, %s, %s, %s)
        );`,
        s.tableName,
        ColID,
        ColAccountID,
        ColStatus,
        ColChain,
        ColNetwork,
        ColToken,
        ColRecipientID,
        ColRecipientKind,
        ColRecipientAddress,
        ColAmount,
        ColMemo,
        ColExpiresAt,
        ColWebhookEndpointID,
        ColMetadata,
        ColCreatedAt,
        ColUpdatedAt,
        ColID,
        ColAccountID, ColStatus,
        ColRecipientAddress, ColChain, ColNetwork, ColToken, ColStatus,
    )

    if _, err := executor.Exec(query); err != nil {
        return fmt.Errorf("failed to create payment onchain intents table: %w", err)
    }

    return nil
}


//
// Insert payment onchain intent.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *IntentStore) Insert(executor helper.Executor, p *IntentInsertParams) error {
    if s == nil {
        return fmt.Errorf("failed to insert payment onchain intent: missing required parameter: request_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to insert payment onchain intent: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to insert payment onchain intent: missing required parameter: executor=null")
    }
    if p == nil {
        return fmt.Errorf("failed to insert payment onchain intent: missing required parameter: request_insert_params=null")
    }
    if err := ValidateIntentAccountID(p.AccountID); err != nil {
        return fmt.Errorf("failed to insert payment onchain intent: %w", err)
    }
    if err := ValidateIntentStatus(p.Status); err != nil {
        return fmt.Errorf("failed to insert payment onchain intent: %w", err)
    }
    if err := ValidateIntentChain(p.Chain); err != nil {
        return fmt.Errorf("failed to insert payment onchain intent: %w", err)
    }
    if err := ValidateIntentNetwork(p.Network); err != nil {
        return fmt.Errorf("failed to insert payment onchain intent: %w", err)
    }
    if err := ValidateIntentToken(p.Token); err != nil {
        return fmt.Errorf("failed to insert payment onchain intent: %w", err)
    }
    if err := ValidateIntentRecipientID(p.RecipientID); err != nil {
        return fmt.Errorf("failed to insert payment onchain intent: %w", err)
    }
    if err := ValidateIntentRecipientKind(p.RecipientKind); err != nil {
        return fmt.Errorf("failed to insert payment onchain intent: %w", err)
    }
    if err := ValidateIntentRecipientAddress(p.RecipientAddress); err != nil {
        return fmt.Errorf("failed to insert payment onchain intent: %w", err)
    }
    if err := ValidateIntentAmount(p.Amount); err != nil {
        return fmt.Errorf("failed to insert payment onchain intent: %w", err)
    }
    if err := ValidateIntentMemo(p.Memo); err != nil {
        return fmt.Errorf("failed to insert payment onchain intent: %w", err)
    }
    if err := ValidateIntentExpiresAt(p.ExpiresAt); err != nil {
        return fmt.Errorf("failed to insert payment onchain intent: %w", err)
    }
    if err := ValidateIntentWebhookEndpointID(p.WebhookEndpointID); err != nil {
        return fmt.Errorf("failed to insert payment onchain intent: %w", err)
    }
    if err := ValidateIntentMetadata(p.Metadata); err != nil {
        return fmt.Errorf("failed to insert payment onchain intent: %w", err)
    }

    if p.ID == 0 {
        p.ID = GenerateIntentID()
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
        ColAccountID,
        ColStatus,
        ColChain,
        ColNetwork,
        ColToken,
        ColRecipientID,
        ColRecipientKind,
        ColRecipientAddress,
        ColAmount,
        ColMemo,
        ColExpiresAt,
        ColWebhookEndpointID,
        ColMetadata,
        ColCreatedAt,
        ColUpdatedAt,
    )

    if _, err := executor.Exec(
        query,
        p.ID,
        p.AccountID,
        p.Status,
        p.Chain,
        p.Network,
        p.Token,
        p.RecipientID,
        p.RecipientKind,
        p.RecipientAddress,
        p.Amount,
        p.Memo,
        p.ExpiresAt,
        p.WebhookEndpointID,
        p.Metadata,
        p.CreatedAt,
        p.UpdatedAt,
    ); err != nil {
        return fmt.Errorf("failed to insert payment onchain intent: %w", err)
    }

    return nil
}


//
// Select payment onchain intents.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *IntentStore) Select(executor helper.Executor, p *IntentSelectParams) ([]*Intent, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select payment onchain intents: missing required parameter: request_store=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select payment onchain intents: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return nil, fmt.Errorf("failed to select payment onchain intents: missing required parameter: executor=null")
    }
    if err := p.Validate(); err != nil {
        return nil, fmt.Errorf("failed to select payment onchain intents: %w", err)
    }

    query, args := p.BuildQuery("SELECT * FROM " + s.tableName)

    // Execute.
    rows, err := executor.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to select payment onchain intents: %w", err)
    }
    defer rows.Close()

    var result []*Intent
    for rows.Next() {
        row := &Intent{}
        if err := rows.Scan(
            &row.ID,
            &row.AccountID,
            &row.Status,
            &row.Chain,
            &row.Network,
            &row.Token,
            &row.RecipientID,
            &row.RecipientKind,
            &row.RecipientAddress,
            &row.Amount,
            &row.Memo,
            &row.ExpiresAt,
            &row.WebhookEndpointID,
            &row.Metadata,
            &row.CreatedAt,
            &row.UpdatedAt,
        ); err != nil {
            return nil, fmt.Errorf("failed to select payment onchain intents: %w", err)
        }

        result = append(result, row)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to select payment onchain intents: %w", err)
    }

    return result, nil
}


//
// Select payment onchain intent by ID.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *IntentStore) SelectByID(executor helper.Executor, id uint64) (*Intent, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select payment onchain intent by id: missing required parameter: request_store=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select payment onchain intent by id: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return nil, fmt.Errorf("failed to select payment onchain intent by id: missing required parameter: executor=null")
    }
    if id == 0 {
        return nil, fmt.Errorf("failed to select payment onchain intent by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColID)

    row := executor.QueryRow(query, id)

    result := &Intent{}
    err := row.Scan(
        &result.ID,
        &result.AccountID,
        &result.Status,
        &result.Chain,
        &result.Network,
        &result.Token,
        &result.RecipientID,
        &result.RecipientKind,
        &result.RecipientAddress,
        &result.Amount,
        &result.Memo,
        &result.ExpiresAt,
        &result.WebhookEndpointID,
        &result.Metadata,
        &result.CreatedAt,
        &result.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to select payment onchain intent by id: %w", err)
    }

    return result, nil
}


//
// Count payment onchain intents.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *IntentStore) Count(executor helper.Executor, p *IntentSelectParams) (int64, error) {
    if s == nil {
        return 0, fmt.Errorf("failed to count payment onchain intents: missing required parameter: request_store=null")
    }
    if s.tableName == "" {
        return 0, fmt.Errorf("failed to count payment onchain intents: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return 0, fmt.Errorf("failed to count payment onchain intents: missing required parameter: executor=null")
    }
    if err := p.Validate(); err != nil {
        return 0, fmt.Errorf("failed to count payment onchain intents: %w", err)
    }

    query, args := p.BuildQuery("SELECT COUNT(*) FROM " + s.tableName)

    var result int64
    if err := executor.QueryRow(query, args...).Scan(&result); err != nil {
        return 0, fmt.Errorf("failed to count payment onchain intents: %w", err)
    }

    return result, nil
}


//
// Delete payment onchain intent by ID.
//
// Version:
//   - 2026-05-16: Added.
//
func (s *IntentStore) DeleteByID(executor helper.Executor, id uint64) error {
    if s == nil {
        return fmt.Errorf("failed to delete payment onchain intent by id: missing required parameter: request_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to delete payment onchain intent by id: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to delete payment onchain intent by id: missing required parameter: executor=null")
    }
    if id == 0 {
        return fmt.Errorf("failed to delete payment onchain intent by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", s.tableName, ColID)

    if _, err := executor.Exec(query, id); err != nil {
        return fmt.Errorf("failed to delete payment onchain intent by id: %w", err)
    }

    return nil
}


//
// Update payment onchain intent status by ID.
//
func (s *IntentStore) UpdateStatusByID(executor helper.Executor, id uint64, status IntentStatus) error {
    if s == nil {
        return fmt.Errorf("failed to update payment onchain intent status by id: missing required parameter: request_store=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update payment onchain intent status by id: missing required parameter: table_name=%q", "empty")
    }
    if executor == nil {
        return fmt.Errorf("failed to update payment onchain intent status by id: missing required parameter: executor=null")
    }
    if id == 0 {
        return fmt.Errorf("failed to update payment onchain intent status by id: invalid parameter: id=0")
    }
    if err := status.Validate(); err != nil {
        return fmt.Errorf("failed to update payment onchain intent status by id: %w", err)
    }

    query := fmt.Sprintf(
        "UPDATE %s SET %s = ? WHERE %s = ?;",
        s.tableName,
        ColStatus,
        ColID,
    )

    args := make([]any, 0, 2)
    args = append(args, status, id)

    if _, err := executor.Exec(query, args...); err != nil {
        return fmt.Errorf("failed to update payment onchain intent status by id: %w", err)
    }

    return nil
}


//
// Build query.
//
// Version:
//   - 2025-05-16: Added.
//
func (p *IntentSelectParams) BuildQuery(selectFromClause string) (string, []any) {
    // Guard.
    if p == nil {
        return selectFromClause, nil
    }

    var query strings.Builder
    query.WriteString(selectFromClause)

    conditions := make([]string, 0, 7)
    args := make([]any, 0, 9)

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
    if p.Token != nil {
        conditions = append(conditions, ColToken + " = ?")
        args = append(args, *p.Token)
    }
    if p.RecipientID != nil {
        conditions = append(conditions, ColRecipientID + " = ?")
        args = append(args, *p.RecipientID)
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
// Validate payment onchain intent select params.
//
// Version:
//   - 2025-05-16: Added.
//
func (p *IntentSelectParams) Validate() error {
    // Guard.
    if p == nil {
        return nil
    }

    if p.AccountID != nil {
        if err := ValidateIntentAccountID(*p.AccountID); err != nil {
            return err
        }
    }
    if p.Status != nil {
        if err := ValidateIntentStatus(*p.Status); err != nil {
            return err
        }
    }
    if p.Chain != nil {
        if err := ValidateIntentChain(*p.Chain); err != nil {
            return err
        }
    }
    if p.Network != nil {
        if err := ValidateIntentNetwork(*p.Network); err != nil {
            return err
        }
    }
    if p.Token != nil {
        if err := ValidateIntentToken(*p.Token); err != nil {
            return err
        }
    }
    if p.RecipientID != nil {
        if err := ValidateIntentRecipientID(*p.RecipientID); err != nil {
            return err
        }
    }
    if p.RecipientAddress != nil {
        if err := ValidateIntentRecipientAddress(*p.RecipientAddress); err != nil {
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
            ColToken,
            ColRecipientID,
            ColRecipientKind,
            ColRecipientAddress,
            ColWebhookEndpointID,
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

