//
// usage_ledger.go
//
package app

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "strings"
    "time"

    _ "github.com/go-sql-driver/mysql"

    "github.com/k4k3ru-hub/db/go/mysql/account"
    "github.com/k4k3ru-hub/db/go/mysql/helper"
)

const (
    CreditBalanceTickScale uint64 = 1_000_000

    DefaultUsageLedgerTableName = "account_app_usage_ledger"
)

var (
)


//
// UsageLedger.
//
// Version:
//   - 2026-05-01: Added.
//
type UsageLedger struct {
    AccountID          uint64            `json:"accountId,string"`
    Status             UsageLedgerStatus `json:"status,string"`
    CreditBalanceTicks uint64            `json:"creditBalanceTicks,string"`
    BonusBalanceTicks  uint64            `json:"bonusBalanceTicks,string"`
    CreditExpiresAt    *time.Time        `json:"creditExpiresAt"`
    BonusExpiresAt     *time.Time        `json:"bonusExpiresAt"`
    MetaData           *string           `json:"metaData"`
    CreatedAt          time.Time         `json:"createdAt"`
    UpdatedAt          time.Time         `json:"updatedAt"`
}


//
// UsageLedgerStore.
//
// Version:
//   - 2026-05-01: Added.
//
type UsageLedgerStore struct {
    executor         helper.Executor
    tableName        string
    accountTableName string
}



type UsageLedgerInsertParams struct {
    AccountID          uint64            `json:"accountId,string"`
    Status             UsageLedgerStatus `json:"status,string"`
    CreditBalanceTicks uint64            `json:"creditBalanceTicks,string"`
    BonusBalanceTicks  uint64            `json:"bonusBalanceTicks,string"`
    CreditExpiresAt    *time.Time        `json:"creditExpiresAt"`
    BonusExpiresAt     *time.Time        `json:"bonusExpiresAt"`
    MetaData           *string           `json:"metaData"`
    CreatedAt          time.Time         `json:"createdAt"`
    UpdatedAt          time.Time         `json:"updatedAt"`
    Ignore             bool              `json:"ignore"`
}


//
// UsageLedgerSelectOption.
//
// Version:
//   - 2026-05-01: Added.
//
type UsageLedgerSelectOption struct {
    AccountID             *uint64            `json:"accountId,string,omitempty"`
    Status                *UsageLedgerStatus `json:"status,omitempty"`
    CreditBalanceTicksGTE *uint64            `json:"creditBalanceTicksGte,string,omitempty"`
    CreditBalanceTicksLTE *uint64            `json:"creditBalanceTicksLte,string,omitempty"`
    BonusBalanceTicksGTE  *uint64            `json:"bonusBalanceTicksGte,string,omitempty"`
    BonusBalanceTicksLTE  *uint64            `json:"bonusBalanceTicksLte,string,omitempty"`
    CreditExpiresAtGTE    *time.Time         `json:"creditExpiresAtGte,omitempty"`
    CreditExpiresAtLTE    *time.Time         `json:"creditExpiresAtLte,omitempty"`
    BonusExpiresAtGTE     *time.Time         `json:"bonusExpiresAtGte,omitempty"`
    BonusExpiresAtLTE     *time.Time         `json:"bonusExpiresAtLte,omitempty"`
    OrderBy               string             `json:"orderBy"`
    OrderByDesc           bool               `json:"orderByDesc"`
    Limit                 int                `json:"limit"`
    Offset                int                `json:"offset"`
}


//
// UsageLedgerUpdateOption.
//
// Version:
//   - 2026-05-01: Added.
//
type UsageLedgerUpdateOption struct {
    AccountID          uint64             `json:"accountId,string"`
    Status             *UsageLedgerStatus `json:"status,omitempty"`
    CreditBalanceTicks *uint64            `json:"creditBalanceTicks,string,omitempty"`
    BonusBalanceTicks  *uint64            `json:"bonusBalanceTicks,string,omitempty"`
    CreditExpiresAt    *time.Time         `json:"creditExpiresAt,omitempty"`
    BonusExpiresAt     *time.Time         `json:"bonusExpiresAt,omitempty"`
    MetaData           *string            `json:"metaData,omitempty"`
}


//
// Create new usage ledger store.
//
// Version:
//   - 2026-05-01: Added.
//
func NewUsageLedgerStore(executor helper.Executor, tableName, accountTableName string) (*UsageLedgerStore, error) {
    if executor == nil {
        return nil, fmt.Errorf("failed to create account app usage ledger store: missing required parameter: executor=null")
    }
    if tableName == "" {
        return nil, fmt.Errorf("failed to create account app usage ledger store: missing required parameter: table_name=%q", "empty")
    }
    if accountTableName == "" {
        return nil, fmt.Errorf("failed to create account app usage ledger store: missing required parameter: account_table_name=empty")
    }

    return &UsageLedgerStore{
        executor:                 executor,
        tableName:          tableName,
        accountTableName:   accountTableName,
    }, nil
}



//
// Validate account APP usage ledger account ID.
//
// Version:
//   - 2026-05-12: Added.
//
func ValidateUsageLedgerAccountID(accountID uint64) error {
    if accountID == 0 {
        return fmt.Errorf("invalid parameter: account_id=0")
    }
    return nil
}


//
// Validate account APP usage ledger account ID.
//
// Version:
//   - 2026-05-12: Added.
//
func (l *UsageLedger) ValidateAccountID() error {
    if l == nil {
        return fmt.Errorf("missing required parameter: account_app_usage_ledger=null")
    }
    return ValidateUsageLedgerAccountID(l.AccountID)
}


//
// Validate account APP usage ledger status.
//
// Version:
//   - 2026-05-12: Added.
//
func ValidateUsageLedgerStatus(status UsageLedgerStatus) error {
    if !status.IsValid() {
        return fmt.Errorf("invalid parameter: status=%d", status)
    }
    return nil
}


//
// Validate account APP usage ledger status.
//
// Version:
//   - 2026-05-12: Added.
//
func (l *UsageLedger) ValidateStatus() error {
    if l == nil {
        return fmt.Errorf("missing required parameter: account_app_usage_ledger=null")
    }
    return ValidateUsageLedgerStatus(l.Status)
}


//
// Validate account APP usage ledger credit expires at.
//
// Version:
//   - 2026-05-12: Added.
//
func ValidateUsageLedgerCreditExpiresAt(creditExpiresAt *time.Time) error {
    if creditExpiresAt == nil {
        return nil
    }
    if (*creditExpiresAt).IsZero() {
        return fmt.Errorf("invalid parameter: credit_expires_at=%q", "empty")
    }
    return nil
}


//
// Validate account APP usage ledger credit expires at.
//
// Version:
//   - 2026-05-12: Added.
//
func (l *UsageLedger) ValidateCreditExpiresAt() error {
    if l == nil {
        return fmt.Errorf("missing required parameter: account_app_usage_ledger=null")
    }
    return ValidateUsageLedgerCreditExpiresAt(l.CreditExpiresAt)
}


//
// Validate account APP usage ledger bonus expires at.
//  
// Version:
//   - 2026-05-12: Added.
//
func ValidateUsageLedgerBonusExpiresAt(bonusExpiresAt *time.Time) error {
    if bonusExpiresAt == nil {
        return nil
    }
    if (*bonusExpiresAt).IsZero() {
        return fmt.Errorf("invalid parameter: bonus_expires_at=%q", "empty")
    }
    return nil
}


//
// Validate account APP usage ledger bonus expires at.
//
// Version:
//   - 2026-05-12: Added.
//
func (l *UsageLedger) ValidateBonusExpiresAt() error {
    if l == nil {
        return fmt.Errorf("missing required parameter: account_app_usage_ledger=null")
    }
    return ValidateUsageLedgerBonusExpiresAt(l.BonusExpiresAt)
}


//
// Validate account APP usage ledger meta data.
//
// Version:
//   - 2026-05-12: Added.
//
func ValidateUsageLedgerMetaData(meta *string) error {
    if meta == nil {
        return nil
    }
    if len([]byte(*meta)) > 4096 {
        return fmt.Errorf("invalid parameter: max_size=4096 meta_data=%q", "too long")
    }
    if !json.Valid([]byte(*meta)) {
        return fmt.Errorf("invalid parameter: meta_data=%q", helper.TruncateRunes(*meta, 1024))
    }
    return nil
}


//
// Validate account APP usage ledger meta data.
//
// Version:
//   - 2026-05-12: Added.
//
func (l *UsageLedger) ValidateMetaData() error {
    if l == nil {
        return fmt.Errorf("missing required parameter: account_app_usage_ledger=null")
    }
    return ValidateUsageLedgerMetaData(l.MetaData)
}


//
// Create account app usage ledger table.
//
// Version:
//   - 2026-05-01: Added.
//
func (s *UsageLedgerStore) CreateTable() error {
    if s == nil {
        return fmt.Errorf("failed to create account app usage ledger table: missing required parameter: usage_ledger_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to create account app usage ledger table: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to create account app usage ledger table: missing required parameter: table_name=%q", "empty")
    }
    if s.accountTableName == "" {
        return fmt.Errorf("failed to create account app usage ledger table: missing required parameter: account_table_name=%q", "empty")
    }

    query := fmt.Sprintf(
        `CREATE TABLE IF NOT EXISTS %s (
            %s BIGINT UNSIGNED NOT NULL COMMENT 'Account ID',
            %s TINYINT UNSIGNED NOT NULL COMMENT 'Status',
            %s BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Credit balance ticks',
            %s BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Bonus balance ticks',
            %s DATETIME NULL COMMENT 'Credit expires at',
            %s DATETIME NULL COMMENT 'Bonus expires at',
            %s JSON NULL COMMENT 'Meta data',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
            PRIMARY KEY (%s),
            KEY idx_account_app_usage_ledger_status (%s),
            CONSTRAINT fk_account_app_usage_ledger_account_id FOREIGN KEY (%s) REFERENCES %s (%s) ON DELETE CASCADE ON UPDATE CASCADE
        );`,
        s.tableName,
        ColAccountID,
        ColStatus,
        ColCreditBalanceTicks,
        ColBonusBalanceTicks,
        ColCreditExpiresAt,
        ColBonusExpiresAt,
        ColMetaData,
        ColCreatedAt,
        ColUpdatedAt,
        ColAccountID,
        ColStatus,
        ColAccountID, s.accountTableName, account.ColID,
    )

    if _, err := s.executor.Exec(query); err != nil {
        return fmt.Errorf("failed to create account app usage ledger table: %w", err)
    }

    return nil
}


//
// Insert account app usage ledger.
//
// Version:
//   - 2026-05-01: Added.
//
func (s *UsageLedgerStore) Insert(p *UsageLedgerInsertParams) error {
    if s == nil {
        return fmt.Errorf("failed to insert account app usage ledger: missing required parameter: usage_ledger_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to insert account app usage ledger: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to insert account app usage ledger: missing required parameter: table_name=%q", "empty")
    }
    if p == nil {
        return fmt.Errorf("failed to insert account app usage ledger: missing required parameter: usage_ledger_insert_params=null")
    }
    if err := ValidateUsageLedgerAccountID(p.AccountID); err != nil {
        return fmt.Errorf("failed to insert account app usage ledger: %w", err)
    }
    if err := ValidateUsageLedgerStatus(p.Status); err != nil {
        return fmt.Errorf("failed to insert account app usage ledger: %w", err)
    }
    if err := ValidateUsageLedgerCreditExpiresAt(p.CreditExpiresAt); err != nil {
        return fmt.Errorf("failed to insert account app usage ledger: %w", err)
    }
    if err := ValidateUsageLedgerBonusExpiresAt(p.BonusExpiresAt); err != nil {
        return fmt.Errorf("failed to insert account app usage ledger: %w", err)
    }
    if err := ValidateUsageLedgerMetaData(p.MetaData); err != nil {
        return fmt.Errorf("failed to insert account app usage ledger: %w", err)
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
        "%s INTO %s (%s, %s, %s, %s, %s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);",
        queryPrefix,
        s.tableName,
        ColAccountID,
        ColStatus,
        ColCreditBalanceTicks,
        ColBonusBalanceTicks,
        ColCreditExpiresAt,
        ColBonusExpiresAt,
        ColMetaData,
        ColCreatedAt,
        ColUpdatedAt,
    )

    if _, err := s.executor.Exec(
        query,
        p.AccountID,
        p.Status,
        p.CreditBalanceTicks,
        p.BonusBalanceTicks,
        p.CreditExpiresAt,
        p.BonusExpiresAt,
        p.MetaData,
        p.CreatedAt,
        p.UpdatedAt,
    ); err != nil {
        return fmt.Errorf("failed to insert account app usage ledger: %w", err)
    }

    return nil
}


//
// Add bonus balance ticks.
//
// Version:
//   - 2026-05-13: Added.
//
func (s *UsageLedgerStore) AddBonusBalanceTicks(accountID uint64, delta uint64) error {
    if s == nil {
        return fmt.Errorf("failed to add account app usage ledger bonus balance ticks: missing required parameter: usage_ledger_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to add account app usage ledger bonus balance ticks: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to add account app usage ledger bonus balance ticks: missing required parameter: table_name=%q", "empty")
    }
    if err := ValidateUsageLedgerAccountID(accountID); err != nil {
        return fmt.Errorf("failed to add account app usage ledger bonus balance ticks: %w", err)
    }
    
    query := fmt.Sprintf(
        "UPDATE %s SET %s = %s + ? WHERE %s = ?;",
        s.tableName,
        ColBonusBalanceTicks,
        ColBonusBalanceTicks,
        ColAccountID,
    )

    if _, err := s.executor.Exec(query, delta, accountID); err != nil {
        return fmt.Errorf("failed to add account app usage ledger bonus balance ticks: %w", err)
    }

    return nil
}


//
// Select account app usage ledger by account ID.
//
// Version:
//   - 2026-05-01: Added.
//
func (s *UsageLedgerStore) SelectByAccountID(accountID uint64) (*UsageLedger, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select account app usage ledger by account id: missing required parameter: usage_ledger_store=null")
    }
    if s.executor == nil {
        return nil, fmt.Errorf("failed to select account app usage ledger by account id: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account app usage ledger by account id: missing required parameter: table_name=%q", "empty")
    }
    if accountID == 0 {
        return nil, fmt.Errorf("failed to select account app usage ledger by account id: invalid parameter: account_id=0")
    }

    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColAccountID)

    result := &UsageLedger{}
    err := s.executor.QueryRow(query, accountID).Scan(
        &result.AccountID,
        &result.Status,
        &result.CreditBalanceTicks,
        &result.BonusBalanceTicks,
        &result.CreditExpiresAt,
        &result.BonusExpiresAt,
        &result.MetaData,
        &result.CreatedAt,
        &result.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to select account app usage ledger by account id: %w", err)
    }

    return result, nil
}


//
// Select account app usage ledger.
//
// Version:
//   - 2026-05-01: Added.
//
func (s *UsageLedgerStore) Select(option *UsageLedgerSelectOption) ([]*UsageLedger, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select account app usage ledger: missing required parameter: usage_ledger_store=null")
    }
    if s.executor == nil {
        return nil, fmt.Errorf("failed to select account app usage ledger: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account app usage ledger: missing required parameter: table_name=%q", "empty")
    }
    if err := option.Validate(); err != nil {
        return nil, fmt.Errorf("failed to select account app usage ledger: %w", err)
    }

    query, args := option.BuildQuery("SELECT * FROM " + s.tableName)

    // Execute.
    rows, err := s.executor.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to select account app usage ledger: %w", err)
    }
    defer rows.Close()

    var result []*UsageLedger
    for rows.Next() {
        row := &UsageLedger{}
        if err := rows.Scan(
            &row.AccountID,
            &row.Status,
            &row.CreditBalanceTicks,
            &row.BonusBalanceTicks,
            &row.CreditExpiresAt,
            &row.BonusExpiresAt,
            &row.MetaData,
            &row.CreatedAt,
            &row.UpdatedAt,
        ); err != nil {
            return nil, fmt.Errorf("failed to select account app usage ledger: %w", err)
        }

        result = append(result, row)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to select account app usage ledger: %w", err)
    }

    return result, nil
}


//  
// Count account app usage ledger.
//      
// Version:
//   - 2026-05-01: Added.
// 
func (s *UsageLedgerStore) Count(option *UsageLedgerSelectOption) (int64, error) {
    if s == nil {
        return 0, fmt.Errorf("failed to count account app usage ledger: missing required parameter: usage_ledger_store=null")
    }
    if s.executor == nil {
        return 0, fmt.Errorf("failed to count account app usage ledger: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return 0, fmt.Errorf("failed to count account app usage ledger: missing required parameter: table_name=%q", "empty")
    }
    if err := option.Validate(); err != nil {
        return 0, fmt.Errorf("failed to select account app usage ledger: %w", err)
    }

    query, args := option.BuildQuery("SELECT COUNT(*) FROM " + s.tableName)

    var result int64
    if err := s.executor.QueryRow(query, args...).Scan(&result); err != nil {
        return 0, fmt.Errorf("failed to count account app usage ledger: %w", err)
    }

    return result, nil
}


//
// Update account app usage ledger.
//
// Version:
//   - 2026-05-01: Added.
//
func (s *UsageLedgerStore) Update(option *UsageLedgerUpdateOption) error {
    if s == nil {
        return fmt.Errorf("failed to update account app usage ledger: missing required parameter: usage_ledger_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to update account app usage ledger: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update account app usage ledger: missing required parameter: table_name=%q", "empty")
    }
    if err := option.Validate(); err != nil {
        return fmt.Errorf("failed to update account app usage ledger: %w", err)
    }

    assignments := make([]string, 0, 6)
    args := make([]any, 0, 7)

    if option.Status != nil {
        assignments = append(assignments, ColStatus + " = ?")
        args = append(args, *option.Status)
    }

    if option.CreditBalanceTicks != nil {
        assignments = append(assignments, ColCreditBalanceTicks + " = ?")
        args = append(args, *option.CreditBalanceTicks)
    }

    if option.BonusBalanceTicks != nil {
        assignments = append(assignments, ColBonusBalanceTicks + " = ?")
        args = append(args, *option.BonusBalanceTicks)
    }

    if option.CreditExpiresAt != nil {
        assignments = append(assignments, ColCreditExpiresAt + " = ?")
        args = append(args, *option.CreditExpiresAt)
    }

    if option.BonusExpiresAt != nil {
        assignments = append(assignments, ColBonusExpiresAt + " = ?")
        args = append(args, *option.BonusExpiresAt)
    }

    if option.MetaData != nil {
        assignments = append(assignments, ColMetaData + " = ?")
        args = append(args, *option.MetaData)
    }

    if len(assignments) == 0 {
        return fmt.Errorf("failed to update account app usage ledger: invalid parameter: assignments=empty")
    }

    args = append(args, option.AccountID)

    query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?;", s.tableName, strings.Join(assignments, ", "), ColAccountID)

    if _, err := s.executor.Exec(query, args...); err != nil {
        return fmt.Errorf("failed to update account app usage ledger: %w", err)
    }

    return nil
}


//
// Delete account app usage ledger by account ID.
//
// Version:
//   - 2026-05-01: Added.
//
func (s *UsageLedgerStore) DeleteByAccountID(accountID uint64) error {
    if s == nil {
        return fmt.Errorf("failed to delete account app usage ledger by account id: missing required parameter: usage_ledger_store=null")
    }
    if s.executor == nil {
        return fmt.Errorf("failed to delete account app usage ledger by account id: missing required parameter: executor=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to delete account app usage ledger by account id: missing required parameter: table_name=%q", "empty")
    }
    if accountID == 0 {
        return fmt.Errorf("failed to delete account app usage ledger by account id: invalid parameter: account_id=0")
    }

    query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", s.tableName, ColAccountID)

    if _, err := s.executor.Exec(query, accountID); err != nil {
        return fmt.Errorf("failed to delete account app usage ledger by account id: %w", err)
    }

    return nil
}


//
// Build query.
//
// Version:
//   - 2025-05-12: Added.
//
func (o *UsageLedgerSelectOption) BuildQuery(selectFromClause string) (string, []any) {
    // Guard.
    if o == nil {
        return selectFromClause, nil
    }

    var query strings.Builder
    query.WriteString(selectFromClause)

    conditions := make([]string, 0, 10)
    args := make([]any, 0, 12)

    if o.AccountID != nil {
        conditions = append(conditions, ColAccountID + " = ?")
        args = append(args, *o.AccountID)
    }
    if o.Status != nil {
        conditions = append(conditions, ColStatus + " = ?")
        args = append(args, *o.Status)
    }
    if o.CreditBalanceTicksGTE != nil {
        conditions = append(conditions, ColCreditBalanceTicks + " >= ?")
        args = append(args, *o.CreditBalanceTicksGTE)
    }
    if o.CreditBalanceTicksLTE != nil {
        conditions = append(conditions, ColCreditBalanceTicks + " <= ?")
        args = append(args, *o.CreditBalanceTicksLTE)
    }
    if o.BonusBalanceTicksGTE != nil {
        conditions = append(conditions, ColBonusBalanceTicks + " >= ?")
        args = append(args, *o.BonusBalanceTicksGTE)
    }
    if o.BonusBalanceTicksLTE != nil {
        conditions = append(conditions, ColBonusBalanceTicks + " <= ?")
        args = append(args, *o.BonusBalanceTicksLTE)
    }
    if o.CreditExpiresAtGTE != nil {
        conditions = append(conditions, ColCreditExpiresAt + " >= ?")
        args = append(args, *o.CreditExpiresAtGTE)
    }
    if o.CreditExpiresAtLTE != nil {
        conditions = append(conditions, ColCreditExpiresAt + " <= ?")
        args = append(args, *o.CreditExpiresAtLTE)
    }
    if o.BonusExpiresAtGTE != nil {
        conditions = append(conditions, ColBonusExpiresAt + " >= ?")
        args = append(args, *o.BonusExpiresAtGTE)
    }
    if o.BonusExpiresAtLTE != nil {
        conditions = append(conditions, ColBonusExpiresAt + " <= ?")
        args = append(args, *o.BonusExpiresAtLTE)
    }

    if len(conditions) > 0 {
        query.WriteString(" WHERE ")
        query.WriteString(strings.Join(conditions, " AND "))
    }

    if o.OrderBy != "" {
        query.WriteString(" ORDER BY ")
        query.WriteString(o.OrderBy)
        if o.OrderByDesc {
            query.WriteString(" DESC")
        }
    }

    if o.Limit > 0 {
        query.WriteString(" LIMIT ? OFFSET ?")
        args = append(args, o.Limit, o.Offset)
    }

    return query.String(), args
}


//
// Validate account app usage ledger select option.
//
// Version:
//   - 2025-05-02: Added.
//
func (o *UsageLedgerSelectOption) Validate() error {
    // Guard.
    if o == nil {
        return fmt.Errorf("missing required parameter: usage_ledger_select_option=null")
    }

    if o.AccountID != nil {
        if err := ValidateUsageLedgerAccountID(*o.AccountID); err != nil {
            return err
        }
    }

    if o.Status != nil {
        if err := ValidateUsageLedgerStatus(*o.Status); err != nil {
            return err
        }
    }

    if o.CreditBalanceTicksGTE != nil && o.CreditBalanceTicksLTE != nil && *o.CreditBalanceTicksGTE > *o.CreditBalanceTicksLTE {
        return fmt.Errorf("invalid parameter: credit_balance_ticks_gte=%d credit_balance_ticks_lte=%d", *o.CreditBalanceTicksGTE, *o.CreditBalanceTicksLTE)
    }

    if o.BonusBalanceTicksGTE != nil && o.BonusBalanceTicksLTE != nil && *o.BonusBalanceTicksGTE > *o.BonusBalanceTicksLTE {
        return fmt.Errorf("invalid parameter: bonus_balance_ticks_gte=%d bonus_balance_ticks_lte=%d", *o.BonusBalanceTicksGTE, *o.BonusBalanceTicksLTE)
    }

    if o.CreditExpiresAtGTE != nil && o.CreditExpiresAtLTE != nil && o.CreditExpiresAtGTE.After(*o.CreditExpiresAtLTE) {
        return fmt.Errorf("invalid parameter: credit_expires_at_gte=%s credit_expires_at_lte=%s", *o.CreditExpiresAtGTE, *o.CreditExpiresAtLTE)
    }

    if o.BonusExpiresAtGTE != nil && o.BonusExpiresAtLTE != nil && o.BonusExpiresAtGTE.After(*o.BonusExpiresAtLTE) {
        return fmt.Errorf("invalid parameter: bonus_expires_at_gte=%s bonus_expires_at_lte=%s", *o.BonusExpiresAtGTE, *o.BonusExpiresAtLTE)
    }

    if o.OrderBy != "" {
        switch o.OrderBy {
        case ColAccountID,
            ColStatus,
            ColCreditBalanceTicks,
            ColBonusBalanceTicks,
            ColCreditExpiresAt,
            ColBonusExpiresAt,
            ColCreatedAt,
            ColUpdatedAt:
        default:
            return fmt.Errorf("invalid parameter: order_by=%q", o.OrderBy)
        }
    }

    if o.Limit < 0 {
        return fmt.Errorf("invalid parameter: limit=%d", o.Limit)
    }
    if o.Offset < 0 {   
        return fmt.Errorf("invalid parameter: offset=%d", o.Offset)
    }

    return nil
}


//
// Validate usage ledger update option.
//
// Version:
//   - 2025-05-02: Added.
//
func (o *UsageLedgerUpdateOption) Validate() error {
    // Guard.
    if o == nil {
        return fmt.Errorf("missing required parameter: usage_ledger_update_option=null")
    }

    if o.AccountID == 0 {
        return fmt.Errorf("invalid parameter: account_id=0")
    }

    if o.Status != nil {
        l := UsageLedger{
            Status: *o.Status,
        }
        if err := l.ValidateStatus(); err != nil {
            return err
        }
    }

    if o.CreditExpiresAt != nil {
        l := UsageLedger{
            CreditExpiresAt: o.CreditExpiresAt,
        }
        if err := l.ValidateCreditExpiresAt(); err != nil {
            return err
        }
    }

    if o.BonusExpiresAt != nil {
        l := UsageLedger{
            BonusExpiresAt: o.BonusExpiresAt,
        }
        if err := l.ValidateBonusExpiresAt(); err != nil {
            return err
        }
    }

    if o.MetaData != nil {
        l := UsageLedger{
            MetaData: o.MetaData,
        }
        if err := l.ValidateMetaData(); err != nil {
            return err
        }
    }

    return nil
}


