//
// types.go
//
package payment

import (
    "database/sql/driver"
    "fmt"
    "strings"
)


type WebhookSignatureAlgorithm string

const (
    WebhookSignatureAlgorithmHMACSHA256 WebhookSignatureAlgorithm = "hmac_sha256"
)


//
// Check whether webhook signature algorithm is valid.
//
// Version:
//   - 2026-06-01: Added.
//
func (a WebhookSignatureAlgorithm) IsValid() bool {
    switch a {
    case WebhookSignatureAlgorithmHMACSHA256:
        return true
    default:
        return false
    }
}


//
// Validate webhook signature algorithm.
//
// Version:
//   - 2026-06-01: Added.
//
func (a WebhookSignatureAlgorithm) Validate() error {
    s := strings.TrimSpace(string(a))
    if s == "" {
        return fmt.Errorf("missing required parameter: webhook_signature_algorithm=%q", "empty")
    }
    if len(s) > 32 {
        return fmt.Errorf("invalid parameter: webhook_signature_algorithm=%q", "too long")
    }
    if !a.IsValid() {
        return fmt.Errorf("invalid parameter: webhook_signature_algorithm=%q", string(a))
    }
    return nil
}


//
// Get webhook signature algorithm as driver.Valuer.
//
// Version:
//   - 2026-06-01: Added.
//
func (a WebhookSignatureAlgorithm) Value() (driver.Value, error) {
    if err := a.Validate(); err != nil {
        return nil, err
    }
    return string(a), nil
}


//
// Scan webhook signature algorithm as sql.Scanner.
//
// Version:
//   - 2026-06-01: Added.
//
func (a *WebhookSignatureAlgorithm) Scan(src any) error {
    if a == nil {
        return fmt.Errorf("missing required parameter: webhook_signature_algorithm=null")
    }

    switch v := src.(type) {
    case string:
        *a = WebhookSignatureAlgorithm(v)
    case []byte:
        *a = WebhookSignatureAlgorithm(string(v))
    case nil:
        return fmt.Errorf("missing required parameter: webhook_signature_algorithm=null")
    default:
        return fmt.Errorf("unsupported parameter: webhook_signature_algorithm: data_type=%T", src)
    }

    if err := a.Validate(); err != nil {
        return err
    }

    return nil
}
