//
// product.go.go
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

	"github.com/k4k3ru-hub/db/go/mysql/helper"
)

const (
	DefaultProductTableName = "account_app_products"
)


var (
    productIDCounter = &helper.IdCounter{}
)


//
// Product.
//
// Version:
//   - 2026-05-02: Added.
//
type Product struct {
	ID             uint64        `json:"id,string"`
	Name           string        `json:"name"`
	Status         ProductStatus `json:"status"`
	Type           ProductType   `json:"type"`
	CreditTicks    uint64        `json:"creditTicks,string"`
	BonusTicks     uint64        `json:"bonusTicks,string"`
	PriceAmount    uint64        `json:"priceAmount,string"`
	PriceCurrency  PriceCurrency `json:"priceCurrency"`
	ExpiresInDays  uint32        `json:"expiresInDays"`
	PurchaseLimit  uint32        `json:"purchaseLimit"`
	Description    *string       `json:"description,omitempty"`
	MetaData       *string       `json:"metaData,omitempty"`
	CreatedAt      time.Time     `json:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt"`
}


//
// ProductStore.
//
// Version:
//   - 2026-05-02: Added.
//
type ProductStore struct {
	db        *sql.DB
	tableName string
}


//
// ProductSelectOption.
//
// Version:
//   - 2026-05-02: Added.
//
type ProductSelectOption struct {
	ID              *uint64        `json:"id,string,omitempty"`
	Name            *string        `json:"name,omitempty"`
	Status          *ProductStatus `json:"status,omitempty"`
	Type            *ProductType   `json:"type,omitempty"`
	CreditTicksGTE  *uint64        `json:"creditTicksGte,string,omitempty"`
	CreditTicksLTE  *uint64        `json:"creditTicksLte,string,omitempty"`
	BonusTicksGTE   *uint64        `json:"bonusTicksGte,string,omitempty"`
	BonusTicksLTE   *uint64        `json:"bonusTicksLte,string,omitempty"`
	PriceAmountGTE  *uint64        `json:"priceAmountGte,string,omitempty"`
	PriceAmountLTE  *uint64        `json:"priceAmountLte,string,omitempty"`
	PriceCurrency   *PriceCurrency `json:"priceCurrency,omitempty"`
	ExpiresInDays   *uint32        `json:"expiresInDays,omitempty"`
	PurchaseLimit   *uint32        `json:"purchaseLimit,omitempty"`
	CreatedAtGTE    *time.Time     `json:"createdAtGte,omitempty"`
	CreatedAtLTE    *time.Time     `json:"createdAtLte,omitempty"`
	OrderBy          string         `json:"orderBy"`
	OrderByDesc      bool           `json:"orderByDesc"`
	Limit            int            `json:"limit"`
	Offset           int            `json:"offset"`
}


//
// ProductUpdateOption.
//
// Version:
//   - 2026-05-02: Added.
//
type ProductUpdateOption struct {
	ID             uint64         `json:"id,string"`
	Name           *string        `json:"name,omitempty"`
	Status         *ProductStatus `json:"status,omitempty"`
	Type           *ProductType   `json:"type,omitempty"`
	CreditTicks    *uint64        `json:"creditTicks,string,omitempty"`
	BonusTicks     *uint64        `json:"bonusTicks,string,omitempty"`
	PriceAmount    *uint64        `json:"priceAmount,string,omitempty"`
	PriceCurrency  *PriceCurrency `json:"priceCurrency,omitempty"`
	ExpiresInDays  *uint32        `json:"expiresInDays,omitempty"`
	PurchaseLimit  *uint32        `json:"purchaseLimit,omitempty"`
	Description    *string        `json:"description,omitempty"`
	MetaData       *string        `json:"metaData,omitempty"`
}


//
// Generate product ID.
//
// Version:
//   - 2026-05-02: Added.
//
func GenerateProductID() uint64 {
	return productIDCounter.GenerateID()
}


//
// Create new product store.
//
// Version:
//   - 2026-05-02: Added.
//
func NewProductStore(db *sql.DB, tableName string) (*ProductStore, error) {
	if db == nil {
		return nil, fmt.Errorf("failed to create account app product store: missing required parameter: db=null")
	}
	if tableName == "" {
		return nil, fmt.Errorf("failed to create account app product store: missing required parameter: table_name=%q", "empty")
	}

	return &ProductStore{
		db:        db,
		tableName: tableName,
	}, nil
}


//
// Validate account product ID.
//
// Version:
//   - 2026-05-12: Added.
//
func ValidateProductID(id uint64) error {
    if id == 0 {
        return fmt.Errorf("invalid parameter: id=0")
    }
    return nil
}


//
// Validate account product ID.
//
// Version:
//   - 2026-05-12: Added.
//
func (p *Product) ValidateID() error {
	if p == nil {
		return fmt.Errorf("missing required parameter: account_product=null")
	}
	return ValidateProductID(p.ID)
}


//
// Validate account product name.
//
// Version:
//   - 2026-05-12: Added.
//
func ValidateProductName(name string) error {
    if strings.TrimSpace(name) == "" {
        return fmt.Errorf("invalid parameter: name=%q", "empty")
    }
    if utf8.RuneCountInString(name) > 64 {
        return fmt.Errorf("invalid parameter: max_length=64 name=%q", "too long")
    }
    return nil
}


//
// Validate account product name.
//
// Version:
//   - 2026-05-12: Added.
//
func (p *Product) ValidateName() error {
	if p == nil {
		return fmt.Errorf("missing required parameter: account_product=null")
	}
	return ValidateProductName(p.Name)
}


//
// Validate account product status.
//
// Version:
//   - 2026-05-12: Added.
//
func ValidateProductStatus(s ProductStatus) error {
    if !s.IsValid() {
        return fmt.Errorf("invalid parameter: status=%d", s)
    }
    return nil
}


//
// Validate account product status.
//
// Version:
//   - 2026-05-12: Added.
//
func (p *Product) ValidateStatus() error {
	if p == nil {
		return fmt.Errorf("missing required parameter: account_product=null")
	}
	return ValidateProductStatus(p.Status)
}


//
// Validate account product type.
//
// Version:
//   - 2026-05-12: Added.
//
func ValidateProductType(t ProductType) error {
    if !t.IsValid() {
        return fmt.Errorf("invalid parameter: type=%d", t)
    }
    return nil
}


//
// Validate account product type.
//
// Version:
//   - 2026-05-12: Added.
//
func (p *Product) ValidateType() error {
	if p == nil {
		return fmt.Errorf("missing required parameter: account_product=null")
	}
	return ValidateProductType(p.Type)
}


//
// Validate account product price currency.
//
// Version:
//   - 2026-05-12: Added.
//
func ValidateProductPriceCurrency(c PriceCurrency) error {
    if !c.IsValid() {
        return fmt.Errorf("invalid parameter: price_currency=%s", c)
    }
    return nil
}


//
// Validate account product price currency.
//
// Version:
//   - 2026-05-12: Added.
//
func (p *Product) ValidatePriceCurrency() error {
	if p == nil {
		return fmt.Errorf("missing required parameter: account_product=null")
	}
	return ValidateProductPriceCurrency(p.PriceCurrency)
}


//
// Validate account product description.
//
// Version:
//   - 2026-05-12: Added.
//
func ValidateProductDescription(description *string) error {
    if description == nil {
        return nil
    }
    if utf8.RuneCountInString(*description) > 255 {
        return fmt.Errorf("invalid parameter: max_length=255 description=%q", "too long")
    }
    return nil
}


//
// Validate account product description.
//
// Version:
//   - 2026-05-12: Added.
//
func (p *Product) ValidateDescription() error {
	if p == nil {
		return fmt.Errorf("missing required parameter: account_product=null")
	}
	return ValidateProductDescription(p.Description)
}


//
// Validate account product meta data.
//
// Version:
//   - 2026-05-12: Added.
//
func ValidateProductMetaData(metaData *string) error {
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
// Validate account product meta data.
//
// Version:
//   - 2026-05-12: Added.
//
func (p *Product) ValidateMetaData() error {
	if p == nil {
		return fmt.Errorf("missing required parameter: account_product=null")
	}
	return ValidateProductMetaData(p.MetaData)
}


//
// Create account app product table.
//
// Version:
//   - 2026-05-02: Added.
//
func (s *ProductStore) CreateTable() error {
	if s == nil {
		return fmt.Errorf("failed to create account app product table: missing required parameter: product_store=null")
	}
	if s.db == nil {
		return fmt.Errorf("failed to create account app product table: missing required parameter: db=null")
	}
	if s.tableName == "" {
		return fmt.Errorf("failed to create account app product table: missing required parameter: table_name=%q", "empty")
	}

	query := fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s (
			%s BIGINT UNSIGNED NOT NULL COMMENT 'ID',
			%s VARCHAR(128) NOT NULL COMMENT 'Name',
			%s TINYINT UNSIGNED NOT NULL COMMENT 'Status',
			%s TINYINT UNSIGNED NOT NULL COMMENT 'Type',
			%s BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Credit ticks',
			%s BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Bonus ticks',
			%s BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Price amount',
			%s VARCHAR(16) NOT NULL COMMENT 'Price currency',
			%s INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Expires in days',
			%s INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Purchase limit',
			%s VARCHAR(255) NULL COMMENT 'Description',
			%s JSON NULL COMMENT 'Meta data',
			%s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
			%s DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
			PRIMARY KEY (%s),
			UNIQUE KEY uk_account_app_products_name (%s),
			KEY idx_account_app_products_status (%s),
			KEY idx_account_app_products_type (%s),
			KEY idx_account_app_products_price_currency (%s)
		);`,
		s.tableName,
		ColID,
		ColName,
		ColStatus,
		ColType,
		ColCreditTicks,
		ColBonusTicks,
		ColPriceAmount,
		ColPriceCurrency,
		ColExpiresInDays,
		ColPurchaseLimit,
		ColDescription,
		ColMetaData,
		ColCreatedAt,
		ColUpdatedAt,
		ColID,
		ColName,
		ColStatus,
		ColType,
		ColPriceCurrency,
	)

	if _, err := s.db.Exec(query); err != nil {
		return fmt.Errorf("failed to create account app product table: %w", err)
	}

	return nil
}


//
// Insert account app product.
//
// Version:
//   - 2026-05-02: Added.
//
func (s *ProductStore) Insert(row *Product) error {
	if s == nil {
		return fmt.Errorf("failed to insert account app product: missing required parameter: product_store=null")
	}
	if s.db == nil {
		return fmt.Errorf("failed to insert account app product: missing required parameter: db=null")
	}
	if s.tableName == "" {
		return fmt.Errorf("failed to insert account app product: missing required parameter: table_name=%q", "empty")
	}
	if row == nil {
		return fmt.Errorf("failed to insert account app product: missing required parameter: product=null")
	}
	if err := row.ValidateName(); err != nil {
		return fmt.Errorf("failed to insert account app product: %w", err)
	}
	if err := row.ValidateStatus(); err != nil {
		return fmt.Errorf("failed to insert account app product: %w", err)
	}
	if err := row.ValidateType(); err != nil {
		return fmt.Errorf("failed to insert account app product: %w", err)
	}
	if err := row.ValidatePriceCurrency(); err != nil {
		return fmt.Errorf("failed to insert account app product: %w", err)
	}
	if err := row.ValidateDescription(); err != nil {
		return fmt.Errorf("failed to insert account app product: %w", err)
	}
	if err := row.ValidateMetaData(); err != nil {
		return fmt.Errorf("failed to insert account app product: %w", err)
	}

	if row.ID == 0 {
		row.ID = GenerateProductID()
	}

	now := time.Now()
	if row.CreatedAt.IsZero() {
		row.CreatedAt = now
	}
	if row.UpdatedAt.IsZero() {
		row.UpdatedAt = now
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
		s.tableName,
		ColID,
		ColName,
		ColStatus,
		ColType,
		ColCreditTicks,
		ColBonusTicks,
		ColPriceAmount,
		ColPriceCurrency,
		ColExpiresInDays,
		ColPurchaseLimit,
		ColDescription,
		ColMetaData,
		ColCreatedAt,
		ColUpdatedAt,
	)

	if _, err := s.db.Exec(
		query,
		row.ID,
		row.Name,
		row.Status,
		row.Type,
		row.CreditTicks,
		row.BonusTicks,
		row.PriceAmount,
		row.PriceCurrency,
		row.ExpiresInDays,
		row.PurchaseLimit,
		row.Description,
		row.MetaData,
		row.CreatedAt,
		row.UpdatedAt,
	); err != nil {
		return fmt.Errorf("failed to insert account app product: %w", err)
	}

	return nil
}


//
// Select account app product by ID.
//
// Version:
//   - 2026-05-02: Added.
//
func (s *ProductStore) SelectByID(id uint64) (*Product, error) {
	if s == nil {
		return nil, fmt.Errorf("failed to select account app product by id: missing required parameter: product_store=null")
	}
	if s.db == nil {
		return nil, fmt.Errorf("failed to select account app product by id: missing required parameter: db=null")
	}
	if s.tableName == "" {
		return nil, fmt.Errorf("failed to select account app product by id: missing required parameter: table_name=%q", "empty")
	}
	if id == 0 {
		return nil, fmt.Errorf("failed to select account app product by id: invalid parameter: id=0")
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1;", s.tableName, ColID)

	result := &Product{}
	err := s.db.QueryRow(query, id).Scan(
		&result.ID,
		&result.Name,
		&result.Status,
		&result.Type,
		&result.CreditTicks,
		&result.BonusTicks,
		&result.PriceAmount,
		&result.PriceCurrency,
		&result.ExpiresInDays,
		&result.PurchaseLimit,
		&result.Description,
		&result.MetaData,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to select account app product by id: %w", err)
	}

	return result, nil
}


//
// Select account app products.
//
// Version:
//   - 2026-05-02: Added.
//
func (s *ProductStore) Select(option *ProductSelectOption) ([]*Product, error) {
	if s == nil {
		return nil, fmt.Errorf("failed to select account app products: missing required parameter: product_store=null")
	}
	if s.db == nil {
		return nil, fmt.Errorf("failed to select account app products: missing required parameter: db=null")
	}
	if s.tableName == "" {
		return nil, fmt.Errorf("failed to select account app products: missing required parameter: table_name=%q", "empty")
	}
	if err := option.Validate(); err != nil {
		return nil, fmt.Errorf("failed to select account app products: %w", err)
	}

    query, args := option.BuildQuery("SELECT * FROM " + s.tableName)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to select account app products: %w", err)
	}
	defer rows.Close()

	var result []*Product
	for rows.Next() {
		row := &Product{}
		if err := rows.Scan(
			&row.ID,
			&row.Name,
			&row.Status,
			&row.Type,
			&row.CreditTicks,
			&row.BonusTicks,
			&row.PriceAmount,
			&row.PriceCurrency,
			&row.ExpiresInDays,
			&row.PurchaseLimit,
			&row.Description,
			&row.MetaData,
			&row.CreatedAt,
			&row.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to select account app products: %w", err)
		}

		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to select account app products: %w", err)
	}

	return result, nil
}


//
// Count account app products.
//
// Version:
//   - 2026-05-02: Added.
//
func (s *ProductStore) Count(option *ProductSelectOption) (int64, error) {
	if s == nil {
		return 0, fmt.Errorf("failed to count account app products: missing required parameter: product_store=null")
	}
	if s.db == nil {
		return 0, fmt.Errorf("failed to count account app products: missing required parameter: db=null")
	}
	if s.tableName == "" {
		return 0, fmt.Errorf("failed to count account app products: missing required parameter: table_name=%q", "empty")
	}
	if err := option.Validate(); err != nil {
		return 0, fmt.Errorf("failed to count account app products: %w", err)
	}

    query, args := option.BuildQuery("SELECT COUNT(*) FROM " + s.tableName)

	var result int64
	if err := s.db.QueryRow(query, args...).Scan(&result); err != nil {
		return 0, fmt.Errorf("failed to count account app products: %w", err)
	}

	return result, nil
}


//
// Update account app product.
//
// Version:
//   - 2026-05-02: Added.
//
func (s *ProductStore) Update(option *ProductUpdateOption) error {
	if s == nil {
		return fmt.Errorf("failed to update account app product: missing required parameter: product_store=null")
	}
	if s.db == nil {
		return fmt.Errorf("failed to update account app product: missing required parameter: db=null")
	}
	if s.tableName == "" {
		return fmt.Errorf("failed to update account app product: missing required parameter: table_name=%q", "empty")
	}
	if err := option.Validate(); err != nil {
		return fmt.Errorf("failed to update account app product: %w", err)
	}

	assignments := make([]string, 0, 11)
	args := make([]any, 0, 12)

	if option.Name != nil {
		assignments = append(assignments, ColName+" = ?")
		args = append(args, *option.Name)
	}

	if option.Status != nil {
		assignments = append(assignments, ColStatus+" = ?")
		args = append(args, *option.Status)
	}

	if option.Type != nil {
		assignments = append(assignments, ColType+" = ?")
		args = append(args, *option.Type)
	}

	if option.CreditTicks != nil {
		assignments = append(assignments, ColCreditTicks+" = ?")
		args = append(args, *option.CreditTicks)
	}

	if option.BonusTicks != nil {
		assignments = append(assignments, ColBonusTicks+" = ?")
		args = append(args, *option.BonusTicks)
	}

	if option.PriceAmount != nil {
		assignments = append(assignments, ColPriceAmount+" = ?")
		args = append(args, *option.PriceAmount)
	}

	if option.PriceCurrency != nil {
		assignments = append(assignments, ColPriceCurrency+" = ?")
		args = append(args, *option.PriceCurrency)
	}

	if option.ExpiresInDays != nil {
		assignments = append(assignments, ColExpiresInDays+" = ?")
		args = append(args, *option.ExpiresInDays)
	}

	if option.PurchaseLimit != nil {
		assignments = append(assignments, ColPurchaseLimit+" = ?")
		args = append(args, *option.PurchaseLimit)
	}

	if option.Description != nil {
		assignments = append(assignments, ColDescription+" = ?")
		args = append(args, *option.Description)
	}

	if option.MetaData != nil {
		assignments = append(assignments, ColMetaData+" = ?")
		args = append(args, *option.MetaData)
	}

	if len(assignments) == 0 {
		return fmt.Errorf("failed to update account app product: invalid parameter: assignments=empty")
	}

	args = append(args, option.ID)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?;", s.tableName, strings.Join(assignments, ", "), ColID)

	if _, err := s.db.Exec(query, args...); err != nil {
		return fmt.Errorf("failed to update account app product: %w", err)
	}

	return nil
}


//
// Delete account app product by ID.
//
// Version:
//   - 2026-05-02: Added.
//
func (s *ProductStore) DeleteByID(id uint64) error {
	if s == nil {
		return fmt.Errorf("failed to delete account app product by id: missing required parameter: product_store=null")
	}
	if s.db == nil {
		return fmt.Errorf("failed to delete account app product by id: missing required parameter: db=null")
	}
	if s.tableName == "" {
		return fmt.Errorf("failed to delete account app product by id: missing required parameter: table_name=%q", "empty")
	}
	if id == 0 {
		return fmt.Errorf("failed to delete account app product by id: invalid parameter: id=0")
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", s.tableName, ColID)

	if _, err := s.db.Exec(query, id); err != nil {
		return fmt.Errorf("failed to delete account app product by id: %w", err)
	}

	return nil
}


//
// Build query.
//
// Version:
//   - 2026-05-12: Added.
//
func (o *ProductSelectOption) BuildQuery(selectFromClause string) (string, []any) {
    // Guard.
    if o == nil {
        return selectFromClause, nil
    }
        
    var query strings.Builder
    query.WriteString(selectFromClause)
    
	conditions := make([]string, 0, 16)
	args := make([]any, 0, 18)

	if o.ID != nil {
		conditions = append(conditions, ColID + " = ?")
		args = append(args, *o.ID)
	}
	if o.Name != nil {
		conditions = append(conditions, ColName + " = ?")
		args = append(args, *o.Name)
	}
	if o.Status != nil {
		conditions = append(conditions, ColStatus + " = ?")
		args = append(args, *o.Status)
	}
	if o.Type != nil {
		conditions = append(conditions, ColType + " = ?")
		args = append(args, *o.Type)
	}
	if o.CreditTicksGTE != nil {
		conditions = append(conditions, ColCreditTicks + " >= ?")
		args = append(args, *o.CreditTicksGTE)
	}
	if o.CreditTicksLTE != nil {
		conditions = append(conditions, ColCreditTicks + " <= ?")
		args = append(args, *o.CreditTicksLTE)
	}
	if o.BonusTicksGTE != nil {
		conditions = append(conditions, ColBonusTicks + " >= ?")
		args = append(args, *o.BonusTicksGTE)
	}
	if o.BonusTicksLTE != nil {
		conditions = append(conditions, ColBonusTicks + " <= ?")
		args = append(args, *o.BonusTicksLTE)
	}
	if o.PriceAmountGTE != nil {
		conditions = append(conditions, ColPriceAmount + " >= ?")
		args = append(args, *o.PriceAmountGTE)
	}
	if o.PriceAmountLTE != nil {
		conditions = append(conditions, ColPriceAmount + " <= ?")
		args = append(args, *o.PriceAmountLTE)
	}
	if o.PriceCurrency != nil {
		conditions = append(conditions, ColPriceCurrency + " = ?")
		args = append(args, *o.PriceCurrency)
	}
	if o.ExpiresInDays != nil {
		conditions = append(conditions, ColExpiresInDays + " = ?")
		args = append(args, *o.ExpiresInDays)
	}
	if o.PurchaseLimit != nil {
		conditions = append(conditions, ColPurchaseLimit + " = ?")
		args = append(args, *o.PurchaseLimit)
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
// Validate account app product select option.
//
// Version:
//   - 2026-05-02: Added.
//
func (o *ProductSelectOption) Validate() error {
	if o == nil {
		return fmt.Errorf("missing required parameter: product_select_option=null")
	}

	if o.ID != nil {
		p := Product{
			ID: *o.ID,
		}
		if err := p.ValidateID(); err != nil {
			return err
		}
	}

	if o.Status != nil {
		p := Product{
			Status: *o.Status,
		}
		if err := p.ValidateStatus(); err != nil {
			return err
		}
	}

	if o.Type != nil {
		p := Product{
			Type: *o.Type,
		}
		if err := p.ValidateType(); err != nil {
			return err
		}
	}

	if o.Name != nil {
		p := Product{
			Name: *o.Name,
		}
		if err := p.ValidateName(); err != nil {
			return err
		}
	}

	if o.CreditTicksGTE != nil && o.CreditTicksLTE != nil && *o.CreditTicksGTE > *o.CreditTicksLTE {
		return fmt.Errorf("invalid parameter: credit_ticks_gte=%d credit_ticks_lte=%d", *o.CreditTicksGTE, *o.CreditTicksLTE)
	}

	if o.BonusTicksGTE != nil && o.BonusTicksLTE != nil && *o.BonusTicksGTE > *o.BonusTicksLTE {
		return fmt.Errorf("invalid parameter: bonus_ticks_gte=%d bonus_ticks_lte=%d", *o.BonusTicksGTE, *o.BonusTicksLTE)
	}

	if o.PriceAmountGTE != nil && o.PriceAmountLTE != nil && *o.PriceAmountGTE > *o.PriceAmountLTE {
		return fmt.Errorf("invalid parameter: price_amount_gte=%d price_amount_lte=%d", *o.PriceAmountGTE, *o.PriceAmountLTE)
	}

	if o.PriceCurrency != nil {
		p := Product{
			PriceCurrency: *o.PriceCurrency,
		}
		if err := p.ValidatePriceCurrency(); err != nil {
			return err
		}
	}

	if o.CreatedAtGTE != nil && o.CreatedAtLTE != nil && o.CreatedAtGTE.After(*o.CreatedAtLTE) {
		return fmt.Errorf("invalid parameter: created_at_gte=%s created_at_lte=%s", *o.CreatedAtGTE, *o.CreatedAtLTE)
	}

	if o.OrderBy != "" {
		switch o.OrderBy {
		case ColID,
			ColStatus,
			ColType,
			ColName,
			ColCreditTicks,
			ColBonusTicks,
			ColPriceAmount,
			ColPriceCurrency,
			ColExpiresInDays,
			ColPurchaseLimit,
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
// Validate account app product update option.
//
// Version:
//   - 2026-05-02: Added.
//
func (o *ProductUpdateOption) Validate() error {
	if o == nil {
		return fmt.Errorf("missing required parameter: product_update_option=null")
	}

	if o.ID == 0 {
		return fmt.Errorf("invalid parameter: id=0")
	}

	if o.Status != nil {
		p := Product{
			Status: *o.Status,
		}
		if err := p.ValidateStatus(); err != nil {
			return err
		}
	}

	if o.Type != nil {
		p := Product{
			Type: *o.Type,
		}
		if err := p.ValidateType(); err != nil {
			return err
		}
	}

	if o.Name != nil {
		p := Product{
			Name: *o.Name,
		}
		if err := p.ValidateName(); err != nil {
			return err
		}
	}

	if o.Description != nil {
		p := Product{
			Description: o.Description,
		}
		if err := p.ValidateDescription(); err != nil {
			return err
		}
	}

	if o.PriceCurrency != nil {
		p := Product{
			PriceCurrency: *o.PriceCurrency,
		}
		if err := p.ValidatePriceCurrency(); err != nil {
			return err
		}
	}

	if o.MetaData != nil {
		p := Product{
			MetaData: o.MetaData,
		}
		if err := p.ValidateMetaData(); err != nil {
			return err
		}
	}

	return nil
}
