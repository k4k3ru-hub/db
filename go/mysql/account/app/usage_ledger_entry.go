//
// usage_ledger_entry.go
//
package app

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "strings"
    "time"
    "unicode/utf8"

    _ "github.com/go-sql-driver/mysql"

    "github.com/k4k3ru-hub/db/go/mysql/account"
    "github.com/k4k3ru-hub/db/go/mysql/helper"
)

const (
    DefaultUsageLedgerEntryTableName = "account_app_usage_ledger_entries"
)

var (
    usageLedgerEntryIDCounter = &helper.IdCounter{}
)


//
// UsageLedgerEntry.
//
// Version:
//   - 2026-05-02: Added.
//
type UsageLedgerEntry struct {
    ID               uint64               `json:"id,string"`
    AccountID        uint64               `json:"accountId,string"`
    Type             UsageLedgerEntryType `json:"type"`
    CreditDeltaTicks int64                `json:"creditDeltaTicks,string"`
    BonusDeltaTicks  int64                `json:"bonusDeltaTicks,string"`
    Description      *string              `json:"description,omitempty"`
    MetaData         *string              `json:"metaData,omitempty"`
    CreatedAt        time.Time            `json:"createdAt"`
    UpdatedAt        time.Time            `json:"updatedAt"`
}


//
// UsageLedgerEntryStore.
//
// Version:
//   - 2026-05-02: Added.
//
type UsageLedgerEntryStore struct {
    db               *sql.DB
    tableName        string
    accountTableName string
}


//  
// UsageLedgerEntrySelectOption.
//  
// Version:
//   - 2026-05-02: Added.
//
type UsageLedgerEntrySelectOption struct {
    ID                   *uint64               `json:"id,string,omitempty"`
    AccountID            *uint64               `json:"accountId,string,omitempty"`
    Type                 *UsageLedgerEntryType `json:"type,omitempty"`
    CreditDeltaTicksGTE  *int64                `json:"creditDeltaTicksGte,string,omitempty"`
    CreditDeltaTicksLTE  *int64                `json:"creditDeltaTicksLte,string,omitempty"`
    BonusDeltaTicksGTE   *int64                `json:"bonusDeltaTicksGte,string,omitempty"`
    BonusDeltaTicksLTE   *int64                `json:"bonusDeltaTicksLte,string,omitempty"`
    CreatedAtGTE         *time.Time            `json:"createdAtGte,omitempty"`
    CreatedAtLTE         *time.Time            `json:"createdAtLte,omitempty"`
    OrderBy              string                `json:"orderBy"`
    OrderByDesc          bool                  `json:"orderByDesc"`
    Limit                int                   `json:"limit"`
    Offset               int                   `json:"offset"`
}


//
// UsageLedgerEntryUpdateOption.
//  
// Version:
//   - 2026-05-02: Added.
//
type UsageLedgerEntryUpdateOption struct {
    ID               uint64                `json:"id,string"`
    Type             *UsageLedgerEntryType `json:"type,omitempty"`
    CreditDeltaTicks *int64                `json:"creditDeltaTicks,string,omitempty"`
    BonusDeltaTicks  *int64                `json:"bonusDeltaTicks,string,omitempty"`
    Description      *string               `json:"description,omitempty"`
    MetaData         *string               `json:"metaData,omitempty"`
}


//
// Generate user ledger entry ID.
//  
// Version:
//   - 2026-05-02: Added.
//
func GenerateUsageLedgerEntryID() uint64 {
    return usageLedgerEntryIDCounter.GenerateID()
}


//
// Create new usage ledger entry store.
//
// Version:
//   - 2026-05-02: Added.
//
func NewUsageLedgerEntryStore(db *sql.DB, tableName, accountTableName string) (*UsageLedgerEntryStore, error) {
    if db == nil {
        return nil, fmt.Errorf("failed to create account app usage ledger entry store: missing required parameter: db=null")
    }
    if tableName == "" {
        return nil, fmt.Errorf("failed to create account app usage ledger entry store: missing required parameter: table_name=%q", "empty")
    }
    if accountTableName == "" {
        return nil, fmt.Errorf("failed to create account app usage ledger entry store: missing required parameter: account_table_name=empty")
    }

    return &UsageLedgerEntryStore{
        db:               db,
        tableName:        tableName,
        accountTableName: accountTableName,
    }, nil
}


//
// Validate account APP usage ledger entry ID.
//
// Version:
//   - 2026-05-12: Added.
//
func ValidateUsageLedgerEntryID(id uint64) error {
    if id == 0 {
        return fmt.Errorf("invalid parameter: id=0")
    }
    return nil
}


//
// Validate account APP usage ledger entry ID.
//
// Version:
//   - 2026-05-12: Added.
//
func (e *UsageLedgerEntry) ValidateID() error {
    if e == nil {
        return fmt.Errorf("missing required parameter: account_app_usage_ledger_entry=null")
    }
    return ValidateUsageLedgerEntryID(e.ID)
}


//
// Validate account APP usage ledger entry account ID.
//
// Version:
//   - 2026-05-12: Added.
//
func ValidateUsageLedgerEntryAccountID(accountID uint64) error {
    if accountID == 0 {
        return fmt.Errorf("invalid parameter: account_id=0")
    }
    return nil
}


//
// Validate account APP usage ledger entry account ID.
//
// Version:
//   - 2026-05-12: Added.
//
func (e *UsageLedgerEntry) ValidateAccountID() error {
    if e == nil {
        return fmt.Errorf("missing required parameter: account_app_usage_ledger_entry=null")
    }
    return ValidateUsageLedgerEntryAccountID(e.AccountID)
}


//
// Validate account APP usage ledger entry type.
//
// Version:
//   - 2026-05-12: Added.
//
func ValidateUsageLedgerEntryType(t UsageLedgerEntryType) error {
    if !t.IsValid() {
        return fmt.Errorf("invalid parameter: type=%d", t)
    }
    return nil
}


//
// Validate account APP usage ledger entry type.
//
// Version:
//   - 2026-05-12: Added.
//
func (e *UsageLedgerEntry) ValidateType() error {
    if e == nil {
        return fmt.Errorf("missing required parameter: account_app_usage_ledger_entry=null")
    }
    return ValidateUsageLedgerEntryType(e.Type)
}


//
// Validate account APP usage ledger entry description.
//
// Version:
//   - 2026-05-12: Added.
//
func ValidateUsageLedgerEntryDescription(description *string) error {
    if description == nil {
        return nil
    }
    if utf8.RuneCountInString(*description) > 255 {
        return fmt.Errorf("invalid parameter: max_length=255 description=%q", "too long")
    }
    return nil
}


//
// Validate account APP usage ledger entry description.
//
// Version:
//   - 2026-05-12: Added.
//
func (e *UsageLedgerEntry) ValidateDescription() error {
    if e == nil {
        return fmt.Errorf("missing required parameter: account_app_usage_ledger_entry=null")
    }
    return ValidateUsageLedgerEntryDescription(e.Description)
}


//
// Validate account APP usage ledger entry meta data.
//  
// Version:
//   - 2026-05-12: Added.
//
func ValidateUsageLedgerEntryMetaData(metaData *string) error {
    if metaData == nil {
        return nil
    }
    if len([]byte(*metaData)) > 4096 {
        return fmt.Errorf("invalid parameter: max_size=4096 meta_data=%q", "too long")
    }
    if !json.Valid([]byte(*metaData)) {
        return fmt.Errorf("invalid parameter: meta_data=%q", helper.TruncateRunes(*metaData, 1024))
    }
    return nil
}


//
// Validate account APP usage ledger entry meta data.
//
// Version:
//   - 2026-05-12: Added.
//
func (e *UsageLedgerEntry) ValidateMetaData() error {
    if e == nil {
        return fmt.Errorf("missing required parameter: account_app_usage_ledger_entry=null")
    }
    return ValidateUsageLedgerEntryMetaData(e.MetaData)
}


//
// Create account app usage ledger entry table.
//
// Version:
//   - 2026-05-02: Added.
//
func (s *UsageLedgerEntryStore) CreateTable() error {
    if s == nil {
        return fmt.Errorf("failed to create account app usage ledger entries table: missing required parameter: usage_ledger_entry_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to create account app usage ledger entries table: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to create account app usage ledger entries table: missing required parameter: table_name=%q", "empty")
    }
    if s.accountTableName == "" {
        return fmt.Errorf("failed to create account app usage ledger entry table: missing required parameter: account_table_name=%q", "empty")
    }

    query := fmt.Sprintf(
        `CREATE TABLE IF NOT EXISTS %s (
            %s BIGINT UNSIGNED NOT NULL COMMENT 'ID',
            %s BIGINT UNSIGNED NOT NULL COMMENT 'Account ID',
            %s TINYINT UNSIGNED NOT NULL COMMENT 'Type',
            %s BIGINT NOT NULL DEFAULT 0 COMMENT 'Credit delta ticks',
            %s BIGINT NOT NULL DEFAULT 0 COMMENT 'Bonus delta ticks',
            %s VARCHAR(255) NULL COMMENT 'Description',
            %s JSON NULL COMMENT 'Meta data',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
            %s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
            PRIMARY KEY (%s),
            KEY idx_account_app_usage_ledger_entry_account_id (%s),
            KEY idx_account_app_usage_ledger_entry_type (%s),
            KEY idx_account_app_usage_ledger_entry_created_at (%s),
            CONSTRAINT fk_account_app_usage_ledger_entries_account_id FOREIGN KEY (%s) REFERENCES %s (%s) ON DELETE CASCADE ON UPDATE CASCADE
        );`,
        s.tableName,
        ColID,
        ColAccountID,
        ColType,
        ColCreditDeltaTicks,
        ColBonusDeltaTicks,
        ColDescription,
        ColMetaData,
        ColCreatedAt,
        ColUpdatedAt,
        ColID,
        ColAccountID,
        ColType,
        ColCreatedAt,
        ColAccountID, s.accountTableName, account.ColID,
    )

    if _, err := s.db.Exec(query); err != nil {
        return fmt.Errorf("failed to create account app usage ledger entries table: %w", err)
    }

    return nil
}


//
// Insert account app usage ledger entry.
//
// Version:
//   - 2026-05-02: Added.
//
func (s *UsageLedgerEntryStore) Insert(row *UsageLedgerEntry) error {
    if s == nil {
        return fmt.Errorf("failed to insert account app usage ledger entry: missing required parameter: usage_ledger_entry_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to insert account app usage ledger entry: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to insert account app usage ledger entry: missing required parameter: table_name=%q", "empty")
    }
    if row == nil {
        return fmt.Errorf("failed to insert account app usage ledger entry: missing required parameter: usage_ledger_entry=null")
    }
    if err := row.ValidateAccountID(); err != nil {
        return fmt.Errorf("failed to insert account app usage ledger entry: %w", err)
    }
    if err := row.ValidateType(); err != nil {
        return fmt.Errorf("failed to insert account app usage ledger entry: %w", err)
    }
    if err := row.ValidateDescription(); err != nil {
        return fmt.Errorf("failed to insert account app usage ledger entry: %w", err)
    }
    if err := row.ValidateMetaData(); err != nil {
        return fmt.Errorf("failed to insert account app usage ledger entry: %w", err)
    }

    if row.ID == 0 {
        row.ID = GenerateUsageLedgerEntryID()
    }

    now := time.Now()
    if row.CreatedAt.IsZero() {
        row.CreatedAt = now
    }
    if row.UpdatedAt.IsZero() {
        row.UpdatedAt = now
    }

    query := fmt.Sprintf(
        "INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);",
        s.tableName,
        ColID,
        ColAccountID,
        ColType,
        ColCreditDeltaTicks,
        ColBonusDeltaTicks,
        ColDescription,
        ColMetaData,
        ColCreatedAt,
        ColUpdatedAt,
    )

    if _, err := s.db.Exec(
        query,
        row.ID,
        row.AccountID,
        row.Type,
        row.CreditDeltaTicks,
        row.BonusDeltaTicks,
        row.Description,
        row.MetaData,
        row.CreatedAt,
        row.UpdatedAt,
    ); err != nil {
        return fmt.Errorf("failed to insert account app usage ledger entry: %w", err)
    }

    return nil
}


//
// Select account app usage ledger entry by ID.
//
// Version:
//   - 2026-05-02: Added.
//
func (s *UsageLedgerEntryStore) SelectByID(id uint64) (*UsageLedgerEntry, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select account app usage ledger entry by id: missing required parameter: usage_ledger_entry_store=null")
    }
    if s.db == nil {
        return nil, fmt.Errorf("failed to select account app usage ledger entry by id: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account app usage ledger entry by id: missing required parameter: table_name=%q", "empty")
    }
    if id == 0 {
        return nil, fmt.Errorf("failed to select account app usage ledger entry by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColID)

    result := &UsageLedgerEntry{}
    err := s.db.QueryRow(query, id).Scan(
        &result.ID,
        &result.AccountID,
        &result.Type,
        &result.CreditDeltaTicks,
        &result.BonusDeltaTicks,
        &result.Description,
        &result.MetaData,
        &result.CreatedAt,
        &result.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to select account app usage ledger entry by id: %w", err)
    }

    return result, nil
}


//
// Select account app usage ledger entries.
//
// Version:
//   - 2026-05-02: Added.
//
func (s *UsageLedgerEntryStore) Select(option *UsageLedgerEntrySelectOption) ([]*UsageLedgerEntry, error) {
    if s == nil {
        return nil, fmt.Errorf("failed to select account app usage ledger entries: missing required parameter: usage_ledger_entry_store=null")
    }
    if s.db == nil {
        return nil, fmt.Errorf("failed to select account app usage ledger entries: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return nil, fmt.Errorf("failed to select account app usage ledger entries: missing required parameter: table_name=%q", "empty")
    }
    if option == nil {
        return nil, fmt.Errorf("failed to select account app usage ledger entries: missing required parameter: select_option=null")
    }
    if err := option.Validate(); err != nil {
        return nil, fmt.Errorf("failed to select account app usage ledger entries: %w", err)
    }

    query, args := option.BuildQuery("SELECT * FROM " + s.tableName)

    rows, err := s.db.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to select account app usage ledger entries: %w", err)
    }
    defer rows.Close()

    var result []*UsageLedgerEntry
    for rows.Next() {
        row := &UsageLedgerEntry{}
        if err := rows.Scan(
            &row.ID,
            &row.AccountID,
            &row.Type,
            &row.CreditDeltaTicks,
            &row.BonusDeltaTicks,
            &row.Description,
            &row.MetaData,
            &row.CreatedAt,
            &row.UpdatedAt,
        ); err != nil {
            return nil, fmt.Errorf("failed to select account app usage ledger entries: %w", err)
        }

        result = append(result, row)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to select account app usage ledger entries: %w", err)
    }

    return result, nil
}


//
// Count account app usage ledger entries.
//
// Version:
//   - 2026-05-02: Added.
//
func (s *UsageLedgerEntryStore) Count(option *UsageLedgerEntrySelectOption) (int64, error) {
    if s == nil {
        return 0, fmt.Errorf("failed to count account app usage ledger entries: missing required parameter: usage_ledger_entry_store=null")
    }
    if s.db == nil {
        return 0, fmt.Errorf("failed to count account app usage ledger entries: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return 0, fmt.Errorf("failed to count account app usage ledger entries: missing required parameter: table_name=%q", "empty")
    }
    if option == nil {
        return 0, fmt.Errorf("failed to count account app usage ledger entries: missing required parameter: select_option=null")
    }
    if err := option.Validate(); err != nil {
        return 0, fmt.Errorf("failed to count account app usage ledger entries: %w", err)
    }

    query, args := option.BuildQuery("SELECT COUNT(*) FROM " + s.tableName)

    var result int64
    if err := s.db.QueryRow(query, args...).Scan(&result); err != nil {
        return 0, fmt.Errorf("failed to count account app usage ledger entries: %w", err)
    }

    return result, nil
}


//
// Update account app usage ledger entry.
//
// Version:
//   - 2026-05-02: Added.
//
func (s *UsageLedgerEntryStore) Update(option *UsageLedgerEntryUpdateOption) error {
    if s == nil {
        return fmt.Errorf("failed to update account app usage ledger entry: missing required parameter: usage_ledger_entry_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to update account app usage ledger entry: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to update account app usage ledger entry: missing required parameter: table_name=%q", "empty")
    }
    if err := option.Validate(); err != nil {
        return fmt.Errorf("failed to update account app usage ledger: %w", err)
    }

    assignments := make([]string, 0, 5)
    args := make([]any, 0, 6)

    if option.Type != nil {
        assignments = append(assignments, ColType + " = ?")
        args = append(args, *option.Type)
    }

    if option.CreditDeltaTicks != nil {
        assignments = append(assignments, ColCreditDeltaTicks + " = ?")
        args = append(args, *option.CreditDeltaTicks)
    }

    if option.BonusDeltaTicks != nil {
        assignments = append(assignments, ColBonusDeltaTicks + " = ?")
        args = append(args, *option.BonusDeltaTicks)
    }

    if option.Description != nil {
        assignments = append(assignments, ColDescription + " = ?")
        args = append(args, *option.Description)
    }

    if option.MetaData != nil {
        assignments = append(assignments, ColMetaData + " = ?")
        args = append(args, *option.MetaData)
    }

    if len(assignments) == 0 {
        return fmt.Errorf("failed to update account app usage ledger entry: invalid parameter: assignments=empty")
    }

    args = append(args, option.ID)

    query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?;", s.tableName, strings.Join(assignments, ", "), ColID)

    if _, err := s.db.Exec(query, args...); err != nil {
        return fmt.Errorf("failed to update account app usage ledger entry: %w", err)
    }

    return nil
}


//
// Delete account app usage ledger entry by ID.
//
// Version:
//   - 2026-05-02: Added.
//
func (s *UsageLedgerEntryStore) DeleteByID(id uint64) error {
    if s == nil {
        return fmt.Errorf("failed to delete account app usage ledger entry by id: missing required parameter: usage_ledger_entry_store=null")
    }
    if s.db == nil {
        return fmt.Errorf("failed to delete account app usage ledger entry by id: missing required parameter: db=null")
    }
    if s.tableName == "" {
        return fmt.Errorf("failed to delete account app usage ledger entry by id: missing required parameter: table_name=%q", "empty")
    }
    if id == 0 {
        return fmt.Errorf("failed to delete account app usage ledger entry by id: invalid parameter: id=0")
    }

    query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", s.tableName, ColID)

    if _, err := s.db.Exec(query, id); err != nil {
        return fmt.Errorf("failed to delete account app usage ledger entry by id: %w", err)
    }

    return nil
}


//      
// Build query.
//  
// Version:
//   - 2025-05-12: Added.
//
func (o *UsageLedgerEntrySelectOption) BuildQuery(selectFromClause string) (string, []any) {
    // Guard.
    if o == nil {
        return selectFromClause, nil
    }

    var query strings.Builder
    query.WriteString(selectFromClause)

    conditions := make([]string, 0, 10)
    args := make([]any, 0, 12)

    if o.ID != nil {
        conditions = append(conditions, ColID + " = ?")
        args = append(args, *o.ID)
    }
    if o.AccountID != nil {
        conditions = append(conditions, ColAccountID + " = ?")
        args = append(args, *o.AccountID)
    }
    if o.Type != nil {
        conditions = append(conditions, ColType + " = ?")
        args = append(args, *o.Type)
    }
    if o.CreditDeltaTicksGTE != nil {
        conditions = append(conditions, ColCreditDeltaTicks + " >= ?")
        args = append(args, *o.CreditDeltaTicksGTE)
    }
    if o.CreditDeltaTicksLTE != nil {
        conditions = append(conditions, ColCreditDeltaTicks + " <= ?")
        args = append(args, *o.CreditDeltaTicksLTE)
    }
    if o.BonusDeltaTicksGTE != nil {
        conditions = append(conditions, ColBonusDeltaTicks + " >= ?")
        args = append(args, *o.BonusDeltaTicksGTE)
    }
    if o.BonusDeltaTicksLTE != nil {
        conditions = append(conditions, ColBonusDeltaTicks + " <= ?")
        args = append(args, *o.BonusDeltaTicksLTE)
    }
    if o.CreatedAtGTE != nil {
        conditions = append(conditions, ColCreatedAt + " >= ?")
        args = append(args, *o.CreatedAtGTE)
    }
    if o.CreatedAtLTE != nil {
        conditions = append(conditions, ColCreatedAt + " <= ?")
        args = append(args, *o.CreatedAtLTE)
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
func (o *UsageLedgerEntrySelectOption) Validate() error {
    // Guard.
    if o == nil {
        return fmt.Errorf("missing required parameter: usage_ledger_entry_select_option=null")
    }

    if o.ID != nil {
        e := UsageLedgerEntry{
            ID: *o.ID,
        }
        if err := e.ValidateID(); err != nil {
            return err
        }
    }

    if o.AccountID != nil {
        e := UsageLedgerEntry{
            AccountID: *o.AccountID,
        }
        if err := e.ValidateAccountID(); err != nil {
            return err
        }
    }

    if o.Type != nil {
        e := UsageLedgerEntry{
            Type: *o.Type,
        }
        if err := e.ValidateType(); err != nil {
            return err
        }
    }

    if o.CreditDeltaTicksGTE != nil && o.CreditDeltaTicksLTE != nil && *o.CreditDeltaTicksGTE > *o.CreditDeltaTicksLTE {
        return fmt.Errorf("invalid parameter: credit_delta_ticks_gte=%d credit_delta_ticks_lte=%d", *o.CreditDeltaTicksGTE, *o.CreditDeltaTicksLTE)
    }

    if o.BonusDeltaTicksGTE != nil && o.BonusDeltaTicksLTE != nil && *o.BonusDeltaTicksGTE > *o.BonusDeltaTicksLTE {
        return fmt.Errorf("invalid parameter: bonus_delta_ticks_gte=%d bonus_delta_ticks_lte=%d", *o.BonusDeltaTicksGTE, *o.BonusDeltaTicksLTE)
    }

    if o.CreatedAtGTE != nil && o.CreatedAtLTE != nil && o.CreatedAtGTE.After(*o.CreatedAtLTE) {
        return fmt.Errorf("invalid parameter: created_at_gte=%s created_at_lte=%s", *o.CreatedAtGTE, *o.CreatedAtLTE)
    }

    if o.OrderBy != "" {
        switch o.OrderBy {
        case ColAccountID,
            ColStatus,
            ColType,
            ColCreditDeltaTicks,
            ColBonusDeltaTicks,
            ColCreatedAt,
            ColUpdatedAt:
        default:
            return fmt.Errorf("invalid parameter: order_by=%s", o.OrderBy)
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
// Validate account app usage ledger update option.
//
// Version:
//   - 2025-05-02: Added.
//
func (o *UsageLedgerEntryUpdateOption) Validate() error {
    // Guard.
    if o == nil {
        return fmt.Errorf("missing required parameter: usage_ledger_entry_update_option=null")
    }

    if o.ID == 0 {
        return fmt.Errorf("invalid parameter: id=0")
    }

    if o.Type != nil {
        e := UsageLedgerEntry{
            Type: *o.Type,
        }
        if err := e.ValidateType(); err != nil {
            return err
        }
    }

    if o.Description != nil {
        e := UsageLedgerEntry{
            Description: o.Description,
        }
        if err := e.ValidateDescription(); err != nil {
            return err
        }
    }

    if o.MetaData != nil {
        e := UsageLedgerEntry{
            MetaData: o.MetaData,
        }
        if err := e.ValidateMetaData(); err != nil {
            return err
        }
    }

    return nil
}
