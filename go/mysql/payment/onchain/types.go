//
// types.go
//
package onchain

import (
    "database/sql/driver"
    "fmt"
    "strconv"

    "github.com/k4k3ru-hub/db/go/mysql/helper"
)


type RequestStatus uint8

const (
    RequestStatusPending RequestStatus = iota
    RequestStatusCompleted
    RequestStatusExpired
    RequestStatusCanceled
    RequestStatusFailed
)

func (s RequestStatus) IsValid() bool {
    return s <= RequestStatusFailed
}


func (s RequestStatus) Validate() error {
    if !s.IsValid() {
        return fmt.Errorf("invalid parameter: request_status=%d", s)
    }
    return nil
}

func (s RequestStatus) Value() (driver.Value, error) {
    if err := s.Validate(); err != nil {
        return nil, err
    }
    return int64(s), nil
}

func (s *RequestStatus) Scan(src any) error {
    if s == nil {
        return fmt.Errorf("missing required parameter: request_status=null")
    }

    switch v := src.(type) {
    case int64:
        *s = RequestStatus(v)
    case []byte:
        n, err := strconv.ParseUint(string(v), 10, 8)
        if err != nil {
            return err
        }
        *s = RequestStatus(n)
    case uint8:
        *s = RequestStatus(v)
    case nil:
        return fmt.Errorf("missing required parameter: request_status=null")
    default:
        return fmt.Errorf("unsupported parameter: type=%T", src)
    }

    if err := s.Validate(); err != nil {
        return err
    }

    return nil
}

type TxStatus uint8

const (
    TxStatusDetected TxStatus = iota
    TxStatusConfirmed
    TxStatusRejected
    TxStatusFailed
)

func (s TxStatus) IsValid() bool {
    return s <= TxStatusFailed
}

func (s TxStatus) Validate() error {
    if !s.IsValid() {
        return fmt.Errorf("invalid parameter: tx_status=%d", s)
    }
    return nil
}

func (s TxStatus) Value() (driver.Value, error) {
    if err := s.Validate(); err != nil {
        return nil, err
    }
    return int64(s), nil
}

func (s *TxStatus) Scan(src any) error {
    if s == nil {
        return fmt.Errorf("missing required parameter: tx_status=null")
    }

    switch v := src.(type) {
    case int64:
        *s = TxStatus(v)
    case []byte:
        n, err := strconv.ParseUint(string(v), 10, 8)
        if err != nil {
            return err
        }
        *s = TxStatus(n)
    case uint8:
        *s = TxStatus(v)
    case nil:
        return fmt.Errorf("missing required parameter: tx_status=null")
    default:
        return fmt.Errorf("unsupported parameter: type=%T", src)
    }

    if err := s.Validate(); err != nil {
        return err
    }

    return nil
}


type Chain string

const (
    ChainBitcoin   Chain = "bitcoin"
    ChainEthereum  Chain = "ethereum"
    ChainBNB       Chain = "bnb"
    ChainPolygon   Chain = "polygon"
    ChainArbitrum  Chain = "arbitrum"
    ChainOptimism  Chain = "optimism"
    ChainBase      Chain = "base"
    ChainAvalanche Chain = "avalanche"
    ChainSui       Chain = "sui"
    ChainSolana    Chain = "solana"
)

func (c Chain) IsValid() bool {
    switch c {
    case ChainBitcoin,
        ChainEthereum,
        ChainBNB,
        ChainPolygon,
        ChainArbitrum,
        ChainOptimism,
        ChainBase,
        ChainAvalanche,
        ChainSui,
        ChainSolana:
        return true
    default:
        return false
    }
}

func (c Chain) Validate() error {
    if string(c) == "" {
        return fmt.Errorf("invalid parameter: chain=%q", "empty")
    }
    if !c.IsValid() {
        return fmt.Errorf("invalid parameter: chain=%q", helper.TruncateRunes(string(c), 16))
    }
    return nil
}

func (c Chain) Value() (driver.Value, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return string(c), nil
}

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
		return fmt.Errorf("unsupported parameter: type=%T", src)
	}

	if err := c.Validate(); err != nil {
        return err
	}

	return nil
}


type Network string

const (
    NetworkMainnet     Network = "mainnet"
    NetworkTestnet     Network = "testnet"
    NetworkDevnet      Network = "devnet"
    NetworkMainnetBeta Network = "mainnet-beta"
)

func (n Network) IsValid() bool {
    switch n {
    case NetworkMainnet,
        NetworkTestnet,
        NetworkDevnet,
        NetworkMainnetBeta:
        return true
    default:
        return false
    }
}

func (n Network) Validate() error {
    if string(n) == "" {
        return fmt.Errorf("invalid parameter: network=%q", "empty")
    }
    if !n.IsValid() {
        return fmt.Errorf("invalid parameter: network=%q", helper.TruncateRunes(string(n), 16))
    }
    return nil
}

func (n Network) Value() (driver.Value, error) {
	if err := n.Validate(); err != nil {
        return nil, err
	}
	return string(n), nil
}

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
		return fmt.Errorf("unsupported parameter: type=%T", src)
	}

    if err := n.Validate(); err != nil {
        return err
	}

	return nil
}


type Asset string

const (
    AssetBTC  Asset = "btc"
    AssetETH  Asset = "eth"
    AssetBNB  Asset = "bnb"
    AssetARB  Asset = "arb"
    AssetSUI  Asset = "sui"
    AssetSOL  Asset = "sol"
    AssetUSDC Asset = "usdc"
)

func (a Asset) IsValid() bool {
    switch a {
    case AssetBTC,
        AssetETH,
        AssetBNB,
        AssetARB,
        AssetSUI,
        AssetSOL,
        AssetUSDC:
        return true
    default:
        return false
    }
}

func (a Asset) Validate() error {
    if string(a) == "" {
        return fmt.Errorf("invalid parameter: asset=%q", "empty")
    }
    if !a.IsValid() {
        return fmt.Errorf("invalid parameter: asset=%q", helper.TruncateRunes(string(a), 16))
    }
    return nil
}

func (a Asset) Value() (driver.Value, error) {
	if err := a.Validate(); err != nil {
		return nil, err
	}
	return string(a), nil
}

func (a *Asset) Scan(src any) error {
	if a == nil {
		return fmt.Errorf("missing required parameter: asset=null")
	}

	switch v := src.(type) {
	case string:
		*a = Asset(v)
	case []byte:
		*a = Asset(string(v))
	case nil:
		return fmt.Errorf("missing required parameter: asset=null")
	default:
		return fmt.Errorf("unsupported parameter: type=%T", src)
	}

	if err := a.Validate(); err != nil {
        return err
	}

	return nil
}




