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

//
// Convert request status to string.
//
// Version:
//   - 2026-05-20: Added.
//
func (s RequestStatus) String() string {
    switch s {
    case RequestStatusPending:
        return "pending"
    case RequestStatusCompleted:
        return "completed"
    case RequestStatusExpired:
        return "expired"
    case RequestStatusCanceled:
        return "canceled"
    case RequestStatusFailed:
        return "failed"
    default:
        return ""
    }
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

//
// Convert tx status to string.
//
// Version:
//   - 2026-05-20: Added.
//  
func (s TxStatus) String() string {
    switch s {
    case TxStatusDetected:
        return "detected"
    case TxStatusConfirmed:
        return "confirmed"
    case TxStatusRejected:
        return "rejected"
    case TxStatusFailed:
        return "failed"
    default:
        return ""
    }
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

func (c Chain) IsValid() bool {
    return k4k3ruOnchainCore.Chain(c).IsValid()
}

func (c Chain) Validate() error {
    return k4k3ruOnchainCore.Chain(c).Validate()
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


type Network k4k3ruOnchainCore.Network

const (
    NetworkMainnet Network = Network(k4k3ruOnchainCore.NetworkMainnet)
    NetworkTestnet Network = Network(k4k3ruOnchainCore.NetworkTestnet)
    NetworkDevnet  Network = Network(k4k3ruOnchainCore.NetworkDevnet)
    NetworkSepolia Network = Network(k4k3ruOnchainCore.NetworkSepolia)
    NetworkHolesky Network = Network(k4k3ruOnchainCore.NetworkHolesky)
)

func (n Network) IsValid() bool {
    return k4k3ruOnchainCore.Network(n).IsValid()
}

func (n Network) Validate() error {
    return k4k3ruOnchainCore.Network(n).Validate()
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


type Token k4k3ruOnchainCore.Token

const (
    TokenAVAX Token = Token(k4k3ruOnchainCore.TokenAVAX)
    TokenBNB  Token = Token(k4k3ruOnchainCore.TokenBNB)
    TokenETH  Token = Token(k4k3ruOnchainCore.TokenETH)
    TokenPOL  Token = Token(k4k3ruOnchainCore.TokenPOL)
    TokenSOL  Token = Token(k4k3ruOnchainCore.TokenSOL)
    TokenSUI  Token = Token(k4k3ruOnchainCore.TokenSUI)
    TokenUSDC Token = Token(k4k3ruOnchainCore.TokenUSDC)
    TokenUSDT Token = Token(k4k3ruOnchainCore.TokenUSDT)
)

func (t Token) IsValid() bool {
    return k4k3ruOnchainCore.Token(t).IsValid()
}

func (t Token) Validate() error {
    return k4k3ruOnchainCore.Token(t).Validate()
}

func (t Token) Value() (driver.Value, error) {
	if err := t.Validate(); err != nil {
		return nil, err
	}
	return string(t), nil
}

func (t *Token) Scan(src any) error {
	if t == nil {
		return fmt.Errorf("missing required parameter: token=null")
	}

	switch v := src.(type) {
	case string:
		*t = Token(v)
	case []byte:
		*t = Token(string(v))
	case nil:
		return fmt.Errorf("missing required parameter: token=null")
	default:
		return fmt.Errorf("unsupported parameter: token: type=%T", src)
	}

	if err := t.Validate(); err != nil {
        return err
	}

	return nil
}




