//
// types.go
//
package email

import (
    "database/sql/driver"
    "fmt"

    "github.com/k4k3ru-hub/db/go/mysql/helper"
)


//
// CredentialStatus.
//
// Version:
//   - 2026-05-03: Added.
//
type CredentialStatus uint8

const (
    CredentialStatusPending CredentialStatus = iota
    CredentialStatusActive
    CredentialStatusInactive
    CredentialStatusSuspended
    CredentialStatusDeleted
)


//
// Convert credential status to string.
//
// Version:
//   - 2026-05-08: Added.
//
func (s CredentialStatus) String() string {
    switch s {
    case CredentialStatusPending:
        return "pending"
    case CredentialStatusActive:
        return "active"
    case CredentialStatusInactive:
        return "inactive"
    case CredentialStatusSuspended:
        return "suspended"
    case CredentialStatusDeleted:
        return "deleted"
    default:
        return ""
    }
}


//
// Validate credential status.
//
// Version:
//   - 2026-05-03: Added.
//
func (s CredentialStatus) IsValid() bool {
    return s <= CredentialStatusDeleted
}


//
// Get credential status as driver.Valuer.
//
// Version:
//   - 2026-05-03: Added.
//
func (s CredentialStatus) Value() (driver.Value, error) {
    if !s.IsValid() {
        return nil, fmt.Errorf("invalid parameter: credential_status=%d", s)
    }

    return int64(s), nil
}


//
// Scan credential status as sql.Scanner.
//
// Version:
//   - 2026-05-03: Added.
//
func (s *CredentialStatus) Scan(value any) error {
    if s == nil {
        return fmt.Errorf("failed to scan: missing required parameter: credential_status=null")
    }

    v, err := helper.ScanUint8("credential_status", value)
    if err != nil {
        return fmt.Errorf("failed to scan: %w", err)
    }

    scanned := CredentialStatus(v)
    if !scanned.IsValid() {
        return fmt.Errorf("failed to scan: invalid parameter: credential_purpose=%d", scanned)
    }

    *s = scanned

    return nil
}
