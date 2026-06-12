//
// types.go
//
package onchain

import (
    "database/sql/driver"
    "fmt"
    "strconv"

    k4k3ruOnchainCore "github.com/k4k3ru-hub/onchain/go/core"
)



type DepositAddressStatus uint8

const (
    DepositAddressStatusActive DepositAddressStatus = iota + 1
    DepositAddressStatusDisabled
    DepositAddressStatusArchived
)


//
// Check whether wallet status is valid.
//
// Version:
//   - 2026-05-25: Added.
//
func (s DepositAddressStatus) IsValid() bool {
    switch s {
    case DepositAddressStatusActive, DepositAddressStatusDisabled, DepositAddressStatusArchived:
        return true
    default:
        return false
    }
}


//
// Validate wallet status.
//
// Version:
//   - 2026-05-25: Added.
//
func (s DepositAddressStatus) Validate() error {
    if !s.IsValid() {
        return fmt.Errorf("invalid parameter: wallet_status=%d", s)
    }
    return nil
}


//
// Get wallet status as driver.Valuer.
//
// Version:
//   - 2026-05-25: Added.
//
func (s DepositAddressStatus) Value() (driver.Value, error) {
    if err := s.Validate(); err != nil {
        return nil, err
    }
    return int64(s), nil
}


//
// Scan wallet status as sql.Scanner.
//
// Version:
//   - 2026-05-25: Added.
//
func (s *DepositAddressStatus) Scan(src any) error {
    if s == nil {
        return fmt.Errorf("missing required parameter: wallet_status=null")
    }

    switch v := src.(type) {
    case int64:
        *s = DepositAddressStatus(v)
    case []byte:
        n, err := strconv.ParseUint(string(v), 10, 8)
        if err != nil {
            return err
        }
        *s = DepositAddressStatus(n)
    case uint8:
        *s = DepositAddressStatus(v)
    case nil:
        return fmt.Errorf("missing required parameter: wallet_status=null")
    default:
        return fmt.Errorf("unsupported parameter: wallet_status: type=%T", src)
    }

    if err := s.Validate(); err != nil {
        return err
    }

    return nil
}


type ChainFamily k4k3ruOnchainCore.ChainFamily

const (
    ChainFamilyEVM    ChainFamily = ChainFamily(k4k3ruOnchainCore.ChainFamilyEVM)
    ChainFamilySolana ChainFamily = ChainFamily(k4k3ruOnchainCore.ChainFamilySolana)
    ChainFamilySui    ChainFamily = ChainFamily(k4k3ruOnchainCore.ChainFamilySui)
)


//
// Check whether chain family is valid.
//
// Version:
//   - 2026-05-25: Added.
//
func (f ChainFamily) IsValid() bool {
    return k4k3ruOnchainCore.ChainFamily(f).IsValid()
}


//
// Validate chain family.
//
// Version:
//   - 2026-05-25: Added.
//
func (f ChainFamily) Validate() error {
    return k4k3ruOnchainCore.ChainFamily(f).Validate()
}


//
// Get chain family as driver.Valuer.
//
// Version:
//   - 2026-05-25: Added.
//
func (f ChainFamily) Value() (driver.Value, error) {
    if err := f.Validate(); err != nil {
        return nil, err
    }
    return string(f), nil
}


//
// Scan chain_ family as sql.Scanner.
//
// Version:
//   - 2026-05-25: Added.
//
func (f *ChainFamily) Scan(src any) error {
    if f == nil {
        return fmt.Errorf("missing required parameter: chain_family=null")
    }

    switch v := src.(type) {
    case string:
        *f = ChainFamily(v)
    case []byte:
        *f = ChainFamily(string(v))
    case nil:
        return fmt.Errorf("missing required parameter: chain_family=null")
    default:
        return fmt.Errorf("unsupported parameter: chain_family: type=%T", src)
    }

    if err := f.Validate(); err != nil {
        return err
    }

    return nil
}


type Network k4k3ruOnchainCore.Network

const (
    NetworkMainnet Network = Network(k4k3ruOnchainCore.NetworkMainnet)
    NetworkTestnet Network = Network(k4k3ruOnchainCore.NetworkTestnet)
    NetworkDevnet  Network = Network(k4k3ruOnchainCore.NetworkDevnet)
    NetworkSepolia Network = Network(k4k3ruOnchainCore.NetworkSepolia)
    NetworkHolesky Network = Network(k4k3ruOnchainCore.NetworkHolesky)
)


//
// Check whether network is valid.
//
// Version:
//   - 2026-05-25: Added.
//
func (n Network) IsValid() bool {
    return k4k3ruOnchainCore.Network(n).IsValid()
}


//
// Validate network.
//
// Version:
//   - 2026-05-25: Added.
//
func (n Network) Validate() error {
    return k4k3ruOnchainCore.Network(n).Validate()
}


//
// Get network as driver.Valuer.
//
// Version:
//   - 2026-05-25: Added.
//
func (n Network) Value() (driver.Value, error) {
    if err := n.Validate(); err != nil {
        return nil, err
    }
    return string(n), nil
}


//
// Scan network as sql.Scanner.
//
// Version:
//   - 2026-05-25: Added.
//
func (n *Network) Scan(src any) error {
    if n == nil {
        return fmt.Errorf("missing required parameter: network=null")
    }

    switch v := src.(type) {
    case string:
        *n = Network(v)
    case []byte:
        *n = Network(string(v))
    case nil:
        return fmt.Errorf("missing required parameter: network=null")
    default:
        return fmt.Errorf("unsupported parameter: network: type=%T", src)
    }

    if err := n.Validate(); err != nil {
        return err
    }

    return nil
}
