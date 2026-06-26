//
// types.go
//
package api

import (
    "database/sql/driver"
    "encoding/json"
    "fmt"

    "github.com/k4k3ru-hub/db/go/mysql/helper"
)

type CredentialStatus uint8

const (
    CredentialStatusPending CredentialStatus = iota
    CredentialStatusActive
    CredentialStatusExpired
    CredentialStatusSuspended
)

//
// Convert credential status to string.
//
// Version:
//   - 2026-05-09: Added.
//
func (s CredentialStatus) String() string {
    switch s {
    case CredentialStatusPending:
        return "pending"
    case CredentialStatusActive:
        return "active"
    case CredentialStatusExpired:
        return "expired"
    case CredentialStatusSuspended:
        return "suspended"
    default:
        return ""
    }
}


//
// Check whether credential status is valid.
//
// Version:
//   - 2026-05-03: Added.
//
func (s CredentialStatus) IsValid() bool {
    return s <= CredentialStatusSuspended
}


//
// Validate credential status.
//
// Version:
//   - 2026-06-03: Added.
//
func (s CredentialStatus) Validate() error {
    if !s.IsValid() {
        fmt.Errorf("invalid parameter: credential_status=%d", s)
    }
    return nil
}


//
// Get credential status as driver.Valuer.
//
// Version:
//   - 2026-05-03: Added.
//
func (s CredentialStatus) Value() (driver.Value, error) {
    if err := s.Validate(); err != nil {
        return nil, err
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
    if err := scanned.Validate(); err != nil {
        return fmt.Errorf("failed to scan: %w", err)
    }

    *s = scanned

    return nil
}


type CredentialSignatureAlgorithm uint8

const (
    CredentialSignatureAlgorithmHMACSHA256 CredentialSignatureAlgorithm = iota + 1
    CredentialSignatureAlgorithmEd25519
)


//
// Convert credential signature algorithm to string.
//
// Version:
//   - 2026-05-09: Added.
//
func (a CredentialSignatureAlgorithm) String() string {
    switch a {
    case CredentialSignatureAlgorithmHMACSHA256:
        return "hmac-sha256"
    case CredentialSignatureAlgorithmEd25519:
        return "ed25519"
    default:
        return ""
    }
}


//
// Check whether credential algorithm is valid.
//
// Version:
//   - 2026-05-03: Added.
//
func (a CredentialSignatureAlgorithm) IsValid() bool {
    switch a {
    case CredentialSignatureAlgorithmHMACSHA256, CredentialSignatureAlgorithmEd25519:
        return true
    default:
        return false
    }
}


//
// Validate credential algorithm.
//
// Version:
//   - 2026-06-03: Added.
//
func (s CredentialSignatureAlgorithm) Validate() error {
    if !s.IsValid() {
        fmt.Errorf("invalid parameter: credential_algorithm=%d", s)
    }
    return nil
}


//
// Get credential algorithm as driver.Valuer.
//
// Version:
//   - 2026-05-03: Added.
//
func (s CredentialSignatureAlgorithm) Value() (driver.Value, error) {
    if err := s.Validate(); err != nil {
        return nil, err
    }

    return int64(s), nil
}


//
// Scan credential signature algorithm as sql.Scanner.
//
// Version:
//   - 2026-05-03: Added.
//
func (a *CredentialSignatureAlgorithm) Scan(value any) error {
    if a == nil {
        return fmt.Errorf("failed to scan: missing required parameter: credential_signature_algorithm=null")
    }

    v, err := helper.ScanUint8("credential_signature_algorithm", value)
    if err != nil {
        return fmt.Errorf("failed to scan: %w", err)
    }

    scanned := CredentialSignatureAlgorithm(v)
    if err := scanned.Validate(); err != nil {
        return fmt.Errorf("failed to scan: %w", err)
    }

    *a = scanned

    return nil
}


type CredentialScopes []string


//
// Validate credential scopes.
//
func (s CredentialScopes) Validate() error {
    if s == nil {
        return nil
    }

    seen := make(map[string]struct{}, len(s))

    for _, scope := range s {
        if scope == "" {
            return fmt.Errorf("invalid parameter: scope=%q", "empty")
        }
        if len(scope) > 128 {
            return fmt.Errorf("invalid parameter: scope=%q max_length=128", "too long")
        }
        if _, ok := seen[scope]; ok {
            return fmt.Errorf("invalid parameter: scope=%q duplicated", scope)
        }

        seen[scope] = struct{}{}
    }

    b, err := json.Marshal(s)
    if err != nil {
        return fmt.Errorf("invalid parameter: %w", err)
    }
    if len(b) > 4096 {
        return fmt.Errorf("invalid parameter: scopes=%q max_size=4096", "too large")
    }

    return nil
}


//
// Get credential scopes as driver.Valuer.
//
// Version:
//   - 2026-06-03: Added.
//
func (s CredentialScopes) Value() (driver.Value, error) {
    if s == nil {
        return []byte("[]"), nil
    }
    return json.Marshal([]string(s))
}


//
// Scan credential scopes as sql.Scanner.
//
// Version:
//   - 2026-06-03: Added.
//
func (s *CredentialScopes) Scan(src any) error {
    if s == nil {
        return fmt.Errorf("invalid parameter: scopes=null")
    }

    if src == nil {
        *s = CredentialScopes{}
        return nil
    }

    var b []byte
    switch v := src.(type) {
    case []byte:
        b = v
    case string:
        b = []byte(v)
    default:
        return fmt.Errorf("invalid parameter: scopes: data_type=%T", src)
    }

    if len(b) == 0 || string(b) == "null" {
        *s = CredentialScopes{}
        return nil
    }

    var scopes []string
    if err := json.Unmarshal(b, &scopes); err != nil {
        return fmt.Errorf("failed to unmarshal scopes: %w", err)
    }

    if scopes == nil {
        scopes = []string{}
    }

    *s = CredentialScopes(scopes)
    return nil
}





