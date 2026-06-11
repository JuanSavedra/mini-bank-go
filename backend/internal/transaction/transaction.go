package transaction

import "errors"

const DomainName = "transaction"

type Type int

const (
	Debit Type = iota
	Credit
	Pix
)

func (t Type) String() string {
	switch t {
	case Debit:
		return "DEBIT"
	case Credit:
		return "CREDIT"
	case Pix:
		return "PIX"
	default:
		return "UNKNOWN"
	}
}

var (
	ErrInvalidAmount       = errors.New("transaction amount must be greater than zero")
	ErrInsufficientBalance = errors.New("insufficient balance for the operation")
	ErrInvalidType         = errors.New("invalid transaction type")
)

type Transaction struct {
	Type        Type  `json:"type"`
	AmountCents int64 `json:"amount_cents"`
}

func NewTransaction(t Type, amountCents int64) (Transaction, error) {
	if err := ValidateType(t); err != nil {
		return Transaction{}, err
	}
	if err := ValidateAmount(amountCents); err != nil {
		return Transaction{}, err
	}
	return Transaction{Type: t, AmountCents: amountCents}, nil
}

func ValidateAmount(amountCents int64) error {
	if amountCents <= 0 {
		return ErrInvalidAmount
	}
	return nil
}

func ValidateType(t Type) error {
	switch t {
	case Debit, Credit, Pix:
		return nil
	default:
		return ErrInvalidType
	}
}

func ValidateWithdrawal(balanceCents, amountCents int64) error {
	if err := ValidateAmount(amountCents); err != nil {
		return err
	}
	if amountCents > balanceCents {
		return ErrInsufficientBalance
	}
	return nil
}
