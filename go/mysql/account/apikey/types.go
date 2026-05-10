//
// types.go
//
package apikey

import (
    "database/sql/driver"
    "fmt"

    "github.com/k4k3ru-hub/db/go/mysql/helper"
)

type CredentialStatus uint8

const (
    StatusPending CredentialStatus = iota
    StatusActive
    StatusExpired
    StatusSuspended
)

//
// Convert credential status to string.
//
// Version:
//   - 2026-05-09: Added.
//
func (s CredentialStatus) String() string {
    switch s {
    case StatusPending:
        return "pending"
    case StatusActive:
        return "active"
    case StatusExpired:
        return "expired"
    case StatusSuspended:
        return "suspended"
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
    return s <= StatusSuspended
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




