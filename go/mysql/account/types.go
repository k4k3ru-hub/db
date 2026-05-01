//
// dto.go
//
package account

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
