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



type WalletStatus uint8

const (
    WalletStatusActive WalletStatus = iota + 1
    WalletStatusDisabled
    WalletStatusArchived
)


//
// Check whether wallet status is valid.
//
// Version:
//   - 2026-05-25: Added.
//
func (s WalletStatus) IsValid() bool {
    switch s {
    case WalletStatusActive, WalletStatusDisabled, WalletStatusArchived:
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
func (s WalletStatus) Validate() error {
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
func (s WalletStatus) Value() (driver.Value, error) {
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
func (s *WalletStatus) Scan(src any) error {
    if s == nil {
        return fmt.Errorf("missing required parameter: wallet_status=null")
    }

    switch v := src.(type) {
    case int64:
        *s = WalletStatus(v)
    case []byte:
        n, err := strconv.ParseUint(string(v), 10, 8)
        if err != nil {
            return err
        }
        *s = WalletStatus(n)
    case uint8:
        *s = WalletStatus(v)
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


type Chain k4k3ruOnchainCore.Chain

const (
    ChainEthereum  Chain = Chain(k4k3ruOnchainCore.ChainEthereum)
    ChainBase      Chain = Chain(k4k3ruOnchainCore.ChainBase)
    ChainBNB       Chain = Chain(k4k3ruOnchainCore.ChainBNB)
    ChainPolygon   Chain = Chain(k4k3ruOnchainCore.ChainPolygon)
    ChainAvalanche Chain = Chain(k4k3ruOnchainCore.ChainAvalanche)
    ChainSolana    Chain = Chain(k4k3ruOnchainCore.ChainSolana)
    ChainSui       Chain = Chain(k4k3ruOnchainCore.ChainSui)
)


//
// Check whether chain is valid.
//
// Version:
//   - 2026-05-25: Added.
//
func (c Chain) IsValid() bool {
    return k4k3ruOnchainCore.Chain(c).IsValid()
}


//
// Validate chain.
//
// Version:
//   - 2026-05-25: Added.
//
func (c Chain) Validate() error {
    return k4k3ruOnchainCore.Chain(c).Validate()
}


//
// Get chain as driver.Valuer.
//
// Version:
//   - 2026-05-25: Added.
//
func (c Chain) Value() (driver.Value, error) {
    if err := c.Validate(); err != nil {
        return nil, err
    }
    return string(c), nil
}


//
// Scan chain as sql.Scanner.
//
// Version:
//   - 2026-05-25: Added.
//
func (c *Chain) Scan(src any) error {
    if c == nil {
        return fmt.Errorf("missing required parameter: chain=null")
    }

    switch v := src.(type) {
    case string:
        *c = Chain(v)
    case []byte:
        *c = Chain(string(v))
    case nil:
        return fmt.Errorf("missing required parameter: chain=null")
    default:
        return fmt.Errorf("unsupported parameter: chain: type=%T", src)
    }

    if err := c.Validate(); err != nil {
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
