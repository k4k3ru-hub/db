//
// types.go
//
package onchain

import (
    "database/sql/driver"
    "fmt"
    "strconv"
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

func (s RequestStatus) Value() (driver.Value, error) {
    if !s.IsValid() {
        return nil, fmt.Errorf("invalid parameter: request_status=%d", s)
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

    if !s.IsValid() {
        return fmt.Errorf("invalid parameter: request_status=%d", *s)
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

func (s TxStatus) Value() (driver.Value, error) {
    if !s.IsValid() {
        return nil, fmt.Errorf("invalid parameter: tx_status=%d", s)
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

    if !s.IsValid() {
        return fmt.Errorf("invalid parameter: tx_status=%d", *s)
    }

    return nil
}


type Chain string

const (
    ChainEVM    Chain = "evm"
    ChainSui    Chain = "sui"
    ChainSolana Chain = "solana"
)

func (c Chain) IsValid() bool {
    switch c {
    case ChainEVM,
        ChainSui,
        ChainSolana:
        return true
    default:
        return false
    }
}

func (c Chain) Value() (driver.Value, error) {
	if !c.IsValid() {
		return nil, fmt.Errorf("invalid parameter: chain=%s", c)
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

	if !c.IsValid() {
		return fmt.Errorf("invalid parameter: chain=%s", *c)
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

func (n Network) Value() (driver.Value, error) {
	if !n.IsValid() {
		return nil, fmt.Errorf("invalid parameter: network=%s", n)
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

	if !n.IsValid() {
		return fmt.Errorf("invalid parameter: network=%s", *n)
	}

	return nil
}
