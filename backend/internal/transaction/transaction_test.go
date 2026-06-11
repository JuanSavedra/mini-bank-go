package transaction

import (
	"errors"
	"testing"
)

func TestTypeValues(t *testing.T) {
	cases := []struct {
		name     string
		typ      Type
		intValue int
		text     string
	}{
		{"debit", Debit, 0, "DEBIT"},
		{"credit", Credit, 1, "CREDIT"},
		{"pix", Pix, 2, "PIX"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if int(c.typ) != c.intValue {
				t.Errorf("iota value: want %d, got %d", c.intValue, int(c.typ))
			}
			if c.typ.String() != c.text {
				t.Errorf("String(): want %q, got %q", c.text, c.typ.String())
			}
		})
	}
}

func TestValidateAmount(t *testing.T) {
	cases := []struct {
		name    string
		amount  int64
		wantErr error
	}{
		{"positive", 1000, nil},
		{"zero", 0, ErrInvalidAmount},
		{"negative", -500, ErrInvalidAmount},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := ValidateAmount(c.amount)
			if !errors.Is(err, c.wantErr) {
				t.Errorf("want %v, got %v", c.wantErr, err)
			}
		})
	}
}

func TestValidateWithdrawal(t *testing.T) {
	cases := []struct {
		name    string
		balance int64
		amount  int64
		wantErr error
	}{
		{"sufficient balance", 10000, 3000, nil},
		{"exact balance", 3000, 3000, nil},
		{"insufficient balance", 1000, 3000, ErrInsufficientBalance},
		{"zero amount", 10000, 0, ErrInvalidAmount},
		{"negative amount", 10000, -1, ErrInvalidAmount},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := ValidateWithdrawal(c.balance, c.amount)
			if !errors.Is(err, c.wantErr) {
				t.Errorf("want %v, got %v", c.wantErr, err)
			}
		})
	}
}

func TestNewTransaction(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		tx, err := NewTransaction(Pix, 2500)
		if err != nil {
			t.Fatalf("did not expect error, got %v", err)
		}
		if tx.Type != Pix || tx.AmountCents != 2500 {
			t.Errorf("incorrect transaction: %+v", tx)
		}
	})

	t.Run("invalid amount", func(t *testing.T) {
		tx, err := NewTransaction(Debit, 0)
		if !errors.Is(err, ErrInvalidAmount) {
			t.Fatalf("want ErrInvalidAmount, got %v", err)
		}
		if tx != (Transaction{}) {
			t.Errorf("expected zero transaction on error, got %+v", tx)
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		_, err := NewTransaction(Type(99), 1000)
		if !errors.Is(err, ErrInvalidType) {
			t.Fatalf("want ErrInvalidType, got %v", err)
		}
	})
}
