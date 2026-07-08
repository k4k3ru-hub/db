//
// types.go
//
package account

import (
    "database/sql/driver"
    "fmt"
    "strconv"

    "github.com/k4k3ru-hub/db/go/mysql/helper"
)


type Status uint8
const (
    StatusPending Status = iota
    StatusActive
    StatusInactive
    StatusSuspended
    StatusDeleted
)


//
// Convert status to string.
//
// Version:
//   - 2026-05-08: Added.
//
func (s Status) String() string {
    switch s {
    case StatusPending:
        return "pending"
    case StatusActive:
        return "active"
    case StatusInactive:
        return "inactive"
    case StatusSuspended:
        return "suspended"
    case StatusDeleted:
        return "deleted"
    default:
        return ""
    }
}

    
func (s Status) IsValid() bool {
    return s <= StatusDeleted
}   


func (s Status) Value() (driver.Value, error) {
    if !s.IsValid() {
        return nil, fmt.Errorf("invalid parameter: status=%d", s)
    }
    return int64(s), nil
}


func (s *Status) Scan(src any) error {
    // Guard.
    if s == nil {
        return fmt.Errorf("missing required parameter: status=null")
    }

    switch v := src.(type) {
    case int64:
        *s = Status(v)
    case []byte:
        n, err := strconv.ParseUint(string(v), 10, 8)
        if err != nil {
            return err
        }
        *s = Status(n)
    case uint8:
        *s = Status(v)
    case nil:
        return fmt.Errorf("missing required parameter: status=null")
    default:
        return fmt.Errorf("unsupported parameter: type=%T", src)
    }

    if !s.IsValid() {
        return fmt.Errorf("invalid parameter: status=%d", *s)
    }

    return nil
}


type Role uint8
const (
    RoleViewer Role = iota
    RoleEditor
    RoleAdmin
)


func (r Role) IsValid() bool {
    return r <= RoleAdmin
}


func (r Role) Value() (driver.Value, error) {
    if !r.IsValid() {
        return nil, fmt.Errorf("invalid parameter: role=%d", r)
    }
    return int64(r), nil
}


func (r *Role) Scan(src any) error {
    if r == nil {
        return fmt.Errorf("missing required parameter: role=null")
    }

    switch v := src.(type) {
    case int64:
        *r = Role(v)
    case []byte:
        n, err := strconv.ParseUint(string(v), 10, 8)
        if err != nil {
            return err
        }
        *r = Role(n)
    case uint8:
        *r = Role(v)
    case nil:
        return fmt.Errorf("missing required parameter: role=null")
    default:
        return fmt.Errorf("unsupported parameter: type=%T", src)
    }

    if !r.IsValid() {
        return fmt.Errorf("invalid parameter: role=%d", *r)
    }

    return nil
}


//
// OTPStatus.
//
// Version:
//   - 2026-06-25: Added.
//
type OTPStatus uint8

const (
    OTPStatusActive OTPStatus = iota
    OTPStatusVerified
    OTPStatusExpired
)


//
// Convert OTP status to string.
//
// Version:
//   - 2026-06-25: Added.
//
func (s OTPStatus) String() string {
    switch s {
    case OTPStatusActive:
        return "active"
    case OTPStatusVerified:
        return "verified"
    case OTPStatusExpired:
        return "expired"
    default:
        return ""
    }
}


//
// Validate OTP status.
//
// Version:
//   - 2026-06-25: Added.
//
func (s OTPStatus) IsValid() bool {
    switch s {
    case OTPStatusActive,
        OTPStatusVerified,
        OTPStatusExpired:
        return true
    default:
        return false
    }
}


//
// Validate OTP status.
//
// Version:
//   - 2026-06-25: Added.
//
func (s OTPStatus) Validate() error {
    if !s.IsValid() {
        return fmt.Errorf("invalid parameter: otp_status=%d", s)
    }
    return nil
}


//
// Get OTP status as driver.Valuer.
//
// Version:
//   - 2026-06-25: Added.
//
func (s OTPStatus) Value() (driver.Value, error) {
    if err := s.Validate(); err != nil {
        return nil, err
    }

    return int64(s), nil
}


//
// Scan OTP status as sql.Scanner.
//
// Version:
//   - 2026-06-25: Added.
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
    if err := scanned.Validate(); err != nil {
        return fmt.Errorf("failed to scan: %w", err)
    }

    *s = scanned

    return nil
}


//
// OTPChannel.
//
// Version:
//   - 2026-06-25: Added.
//
type OTPChannel uint8

const (
    OTPChannelEmail OTPChannel = iota + 1
    OTPChannelSMS
)


//
// Convert OTP channel to string.
//
// Version:
//   - 2026-06-25: Added.
//
func (c OTPChannel) String() string {
    switch c {
    case OTPChannelEmail:
        return "email"
    case OTPChannelSMS:
        return "sms"
    default:
        return ""
    }
}


//
// Validate OTP channel.
//
// Version:
//   - 2026-06-25: Added.
//
func (c OTPChannel) IsValid() bool {
    switch c {
    case OTPChannelEmail,
         OTPChannelSMS:
        return true
    default:
        return false
    }
}


//
// Validate OTP channel.
//
// Version:
//   - 2026-06-25: Added.
//
func (c OTPChannel) Validate() error {
    if !c.IsValid() {
        return fmt.Errorf("invalid parameter: otp_channel=%d", c)
    }
    return nil
}


//
// Get OTP channel as driver.Valuer.
//
// Version:
//   - 2026-06-25: Added.
//
func (c OTPChannel) Value() (driver.Value, error) {
    if err := c.Validate(); err != nil {
        return nil, err
    }

    return int64(c), nil
}


//
// Scan OTP channel as sql.Scanner.
//
// Version:
//   - 2026-06-25: Added.
//
func (c *OTPChannel) Scan(value interface{}) error {
    if c == nil {
        return fmt.Errorf("failed to scan: missing required parameter: otp_channel=null")
    }

    v, err := helper.ScanUint8("otp_channel", value)
    if err != nil {
        return fmt.Errorf("failed to scan: %w", err)
    }

    scanned := OTPChannel(v)
    if err := scanned.Validate(); err != nil {
        return fmt.Errorf("failed to scan: %w", err)
    }

    *c = scanned

    return nil
}


//
// OTPPurpose.
//
// Version:
//   - 2026-06-25: Added.
//
type OTPPurpose uint16

const (
    // Account: Email
    OTPPurposeAccountEmailCreateCredential OTPPurpose = 1001
    OTPPurposeAccountEmailUpdateCredential OTPPurpose = 1002
    OTPPurposeAccountEmailLogin            OTPPurpose = 1003

    // Account: API
    OTPPurposeAccountAPICreateCredential OTPPurpose = 2001
    OTPPurposeAccountAPIUpdateCredential OTPPurpose = 2002
)


//
// Convert OTP purpose to string.
//
// Version:
//   - 2026-06-25: Added.
//
func (p OTPPurpose) String() string {
    switch p {
    case OTPPurposeAccountEmailCreateCredential:
        return "account-email-create-credential"
    case OTPPurposeAccountAPICreateCredential:
        return "account-api-create-credential"
    case OTPPurposeAccountAPIUpdateCredential:
        return "account-api-update-credential"
    default:
        return ""
    }
}


//
// Validate OTP purpose.
//
// Version:
//   - 2026-06-25: Added.
//
func (p OTPPurpose) IsValid() bool {
    switch p {
    case OTPPurposeAccountEmailCreateCredential,
         OTPPurposeAccountAPICreateCredential,
         OTPPurposeAccountAPIUpdateCredential:
        return true
    default:
        return false
    }
}


//
// Validate OTP purpose.
//
// Version:
//   - 2026-06-25: Added.
//
func (p OTPPurpose) Validate() error {
    if !p.IsValid() {
        return fmt.Errorf("invalid parameter: otp_purpose=%d", p)
    }
    return nil
}


//
// Get OTP purpose as driver.Valuer.
//
// Version:
//   - 2026-06-25: Added.
//
func (p OTPPurpose) Value() (driver.Value, error) {
    if err := p.Validate(); err != nil {
        return nil, err
    }

    return int64(p), nil
}


//
// Scan OTP purpose as sql.Scanner.
//
// Version:
//   - 2026-06-25: Added.
//
func (p *OTPPurpose) Scan(value interface{}) error {
    if p == nil {
        return fmt.Errorf("failed to scan: missing required parameter: otp_purpose=null")
    }

    v, err := helper.ScanUint16("otp_purpose", value)
    if err != nil {
        return fmt.Errorf("failed to scan: %w", err)
    }

    scanned := OTPPurpose(v)
    if err := scanned.Validate(); err != nil {
        return fmt.Errorf("failed to scan: %w", err)
    }

    *p = scanned

    return nil
}
