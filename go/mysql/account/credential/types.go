//
// dto.go
//
package credential

import (
    "database/sql/driver"
    "fmt"
    "strconv"
)


type Status uint8
const (
    StatusUnverified Status = iota
    StatusActive
    StatusInactive
    StatusSuspended
    StatusDeleted
)

    
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


