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
    CredentialStatusActive CredentialStatus = iota
    CredentialStatusInactive
    CredentialStatusSuspended
    CredentialStatusDeleted
)


//
// Validate credential status.
//
// Version:
//   - 2026-05-03: Added.
//
func (s CredentialStatus) IsValid() bool {
    switch s {
    case CredentialStatusActive,
        CredentialStatusInactive,
        CredentialStatusSuspended,
        CredentialStatusDeleted:
        return true
    default:
        return false
    }
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



//
// OTPStatus.
//
// Version:
//   - 2026-05-03: Added.
//
type OTPStatus uint8

const (
    OTPStatusActive OTPStatus = iota
    OTPStatusVerified
)


//
// Validate OTP status.
//
// Version:
//   - 2026-05-03: Added.
//
func (s OTPStatus) IsValid() bool {
    switch s {
    case OTPStatusActive,
        OTPStatusVerified:
        return true
    default:
        return false
    }
}


//
// Get OTP status as driver.Valuer.
//
// Version:
//   - 2026-05-03: Added.
//
func (s OTPStatus) Value() (driver.Value, error) {
    if !s.IsValid() {
        return nil, fmt.Errorf("invalid parameter: otp_status=%d", s)
    }

    return int64(s), nil
}


//
// Scan OTP status as sql.Scanner.
//
// Version:
//   - 2026-05-03: Added.
//
func (s *OTPStatus) Scan(value interface{}) error {
    if s == nil {
        return fmt.Errorf("failed to scan: missing required parameter: otp_status=null")
    }

    v, err := helper.ScanUint8("otp_status", value)
    if err != nil {
        return fmt.Errorf("failed to scan: %w", err)
    }

    scanned := OTPStatus(v)
    if !scanned.IsValid() {
        return fmt.Errorf("failed to scan: invalid parameter: otp_status=%d", scanned)
    }

    *s = scanned

    return nil
}


//
// OTPPurpose.
//
// Version:
//   - 2026-05-03: Added.
//
type OTPPurpose uint8

const (
    OTPPurposeEmailCreate OTPPurpose = iota + 1
    OTPPurposeAPIKeyCreate
)


//
// Validate OTP purpose.
//
// Version:
//   - 2026-05-03: Added.
//
func (p OTPPurpose) IsValid() bool {
    switch p {
    case OTPPurposeEmailCreate,
        OTPPurposeAPIKeyCreate:
        return true
    default:
        return false
    }
}


//
// Get OTP purpose as driver.Valuer.
//
// Version:
//   - 2026-05-03: Added.
//
func (p OTPPurpose) Value() (driver.Value, error) {
    if !p.IsValid() {
        return nil, fmt.Errorf("invalid parameter: otp_purpose=%d", p)
    }

    return int64(p), nil
}


//
// Scan OTP purpose as sql.Scanner.
//
// Version:
//   - 2026-05-03: Added.
//
func (p *OTPPurpose) Scan(value interface{}) error {
    if p == nil {
        return fmt.Errorf("failed to scan: missing required parameter: otp_purpose=null")
    }

    v, err := helper.ScanUint8("otp_purpose", value)
    if err != nil {
        return fmt.Errorf("failed to scan: %w", err)
    }

    scanned := OTPPurpose(v)
    if !scanned.IsValid() {
        return fmt.Errorf("failed to scan: invalid parameter: otp_purpose=%d", scanned)
    }

    *p = scanned

    return nil
}


