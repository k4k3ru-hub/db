//
// types.go
//
package app

import (
    "database/sql/driver"
    "fmt"
)


//
// UsageLedgerStatus.
//
// Version:
//   - 2026-05-01: Added.
//
type UsageLedgerStatus uint8

const (
    UsageLedgerStatusActive     UsageLedgerStatus = iota
    UsageLedgerStatusSuspended
)

//
// Is valid.
//
func (s UsageLedgerStatus) IsValid() bool {
    switch s {
    case UsageLedgerStatusActive,
        UsageLedgerStatusSuspended:
        return true
    default:
        return false
    }
}


//
// Value.
//
func (s UsageLedgerStatus) Value() (driver.Value, error) {
    if !s.IsValid() {
        return nil, fmt.Errorf("invalid usage ledger status: %d", s)
    }
    return int64(s), nil
}


//
// Scan.
//
func (s *UsageLedgerStatus) Scan(value any) error {
    if s == nil {
        return fmt.Errorf("failed to scan usage ledger status: receiver=null")
    }

    v, err := scanUint8(value)
    if err != nil {
        return fmt.Errorf("failed to scan usage ledger status: %w", err)
    }

    status := UsageLedgerStatus(v)
    if !status.IsValid() {
        return fmt.Errorf("failed to scan usage ledger status: invalid value=%d", v)
    }

    *s = status
    return nil
}


type UsageLedgerEntryType uint8

const (
    UsageTypeUnknown UsageLedgerEntryType = iota
    UsageTypeDeposit
    UsageTypeBonusGrant
    UsageTypeRefund
    UsageTypeAdjustment
    UsageTypeExpiration
    UsageTypeHTTPRequest
    UsageTypeWSConnection
    UsageTypeWSSubscription
    UsageTypeFIXSession
    UsageTypeFIXMessage
)

func (t UsageLedgerEntryType) IsValid() bool {
    return t <= UsageTypeFIXMessage
}


//
// Value.
//
func (t UsageLedgerEntryType) Value() (driver.Value, error) {
    if !t.IsValid() {
        return nil, fmt.Errorf("invalid usage ledger entry type: %d", t)
    }
    return int64(t), nil
}

//
// Scan.
//
func (t *UsageLedgerEntryType) Scan(value any) error {
    if t == nil {
        return fmt.Errorf("failed to scan usage ledger entry type: receiver=null")
    }

    v, err := scanUint8(value)
    if err != nil {
        return fmt.Errorf("failed to scan usage ledger entry type: %w", err)
    }

    typ := UsageLedgerEntryType(v)
    if !typ.IsValid() {
        return fmt.Errorf("failed to scan usage ledger entry type: invalid value=%d", v)
    }

    *t = typ
    return nil
}



type ProductStatus uint8

const (
    ProductStatusInactive ProductStatus = iota
    ProductStatusActive
)

func (s ProductStatus) IsValid() bool {
    switch s {
    case ProductStatusInactive,
        ProductStatusActive:
        return true
    default:
        return false
    }
}


//
// Value.
//
func (s ProductStatus) Value() (driver.Value, error) {
    if !s.IsValid() {
        return nil, fmt.Errorf("invalid product status: %d", s)
    }
    return int64(s), nil
}

//
// Scan.
//
func (s *ProductStatus) Scan(value any) error {
    if s == nil {
        return fmt.Errorf("failed to scan product status: receiver=null")
    }

    v, err := scanUint8(value)
    if err != nil {
        return fmt.Errorf("failed to scan product status: %w", err)
    }

    status := ProductStatus(v)
    if !status.IsValid() {
        return fmt.Errorf("failed to scan product status: invalid value=%d", v)
    }

    *s = status
    return nil
}


type ProductType uint8

const (
    ProductTypeGeneral ProductType = iota
    ProductTypeCampaign
    ProductTypeTrial
)

func (t ProductType) IsValid() bool {
    switch t {
    case ProductTypeGeneral,
        ProductTypeCampaign,
        ProductTypeTrial:
        return true
    default:
        return false
    }
}


//
// Value.
//
func (t ProductType) Value() (driver.Value, error) {
    if !t.IsValid() {
        return nil, fmt.Errorf("invalid product type: %d", t)
    }
    return int64(t), nil
}

//
// Scan.
//
func (t *ProductType) Scan(value any) error {
    if t == nil {
        return fmt.Errorf("failed to scan product type: receiver=null")
    }

    v, err := scanUint8(value)
    if err != nil {
        return fmt.Errorf("failed to scan product type: %w", err)
    }

    typ := ProductType(v)
    if !typ.IsValid() {
        return fmt.Errorf("failed to scan product type: invalid value=%d", v)
    }

    *t = typ
    return nil
}


type PriceCurrency string

const (
    PriceCurrencyUSD PriceCurrency = "usd"
    PriceCurrencyEUR PriceCurrency = "eur"
    PriceCurrencyJPY PriceCurrency = "jpy"
)

func (c PriceCurrency) IsValid() bool {
    switch c {
    case PriceCurrencyUSD,
        PriceCurrencyEUR,
        PriceCurrencyJPY:
        return true
    default:
        return false
    }
}

func (c PriceCurrency) Value() (driver.Value, error) {
    if !c.IsValid() {
        return nil, fmt.Errorf("invalid price currency: %s", c)
    }
    return string(c), nil
}

func (c *PriceCurrency) Scan(value any) error {
    if c == nil {
        return fmt.Errorf("failed to scan price currency: receiver=null")
    }

    switch v := value.(type) {
    case string:
        currency := PriceCurrency(v)
        if !currency.IsValid() {
            return fmt.Errorf("failed to scan price currency: invalid value=%s", v)
        }
        *c = currency
        return nil

    case []byte:
        currency := PriceCurrency(string(v))
        if !currency.IsValid() {
            return fmt.Errorf("failed to scan price currency: invalid value=%s", string(v))
        }
        *c = currency
        return nil

    default:
        return fmt.Errorf("failed to scan price currency: unsupported scan type=%T", value)
    }
}


//
// Scan uint8.
//
// Version:
//   - 2026-05-03: Added.
//
func scanUint8(value any) (uint8, error) {
    switch v := value.(type) {
    case int64:
        if v < 0 || v > 255 {
            return 0, fmt.Errorf("invalid value=%d", v)
        }
        return uint8(v), nil
    case uint8:
        return v, nil
    case []byte:
        var n uint64
        if _, err := fmt.Sscanf(string(v), "%d", &n); err != nil {
            return 0, fmt.Errorf("invalid value=%q", string(v))
        }
        if n > 255 {
            return 0, fmt.Errorf("invalid value=%d", n)
        }
        return uint8(n), nil
    case string:
        var n uint64
        if _, err := fmt.Sscanf(v, "%d", &n); err != nil {
            return 0, fmt.Errorf("invalid value=%q", v)
        }
        if n > 255 {
            return 0, fmt.Errorf("invalid value=%d", n)
        }
        return uint8(n), nil
    default:
        return 0, fmt.Errorf("unsupported scan type=%T", value)
    }
}
