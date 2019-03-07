package gf256

import (
	"testing"
)

// Note: These tests are built with the Rijndael polynomial

func TestAddGF256(t *testing.T) {
	tables := []struct {
		a byte
		b byte
		r byte
	}{
		{0x53, 0xca, 0x99},
		{0xca, 0x53, 0x99},
	}

	f := NewField(RijndaelPolynomial, RijndaelGenerator)

	for _, table := range tables {
		res := f.Add(table.a, table.b)
		if res != table.r {
			t.Errorf("Result of (%d+%d) wrong, got: %d, want: %d.", table.a, table.b, res, table.r)
		}
	}
}

func TestSubGF256(t *testing.T) {
	tables := []struct {
		a byte
		b byte
		r byte
	}{
		{0x53, 0xca, 0x99},
		{0xca, 0x53, 0x99},
	}

	f := NewField(RijndaelPolynomial, RijndaelGenerator)

	for _, table := range tables {
		res := f.Add(table.a, table.b)
		if res != table.r {
			t.Errorf("Result of (%d-%d) wrong, got: %d, want: %d.", table.a, table.b, res, table.r)
		}
	}
}

func TestMulGF256(t *testing.T) {
	tables := []struct {
		a byte
		b byte
		r byte
	}{
		{0, 0, 0},
		{0x1, 0x5a, 0x5a},
		{0x5a, 0x51, 0xfc},
	}

	f := NewField(RijndaelPolynomial, RijndaelGenerator)

	for _, table := range tables {
		res := f.Mul(table.a, table.b)
		if res != table.r {
			t.Errorf("Result of (%d*%d) wrong, got: %d, want: %d.", table.a, table.b, res, table.r)
		}
	}
}

func TestInvGF256(t *testing.T) {
	tables := []struct {
		a byte
		b byte
	}{
		{0x53, 0xca},
	}

	f := NewField(RijndaelPolynomial, RijndaelGenerator)

	for _, table := range tables {
		res := f.Inv(table.a)
		if res != table.b {
			t.Errorf("Result of (%d)^-1 wrong, got: %d, want: %d.", table.a, res, table.b)
		}
	}
}

func TestPowerGF256(t *testing.T) {
	tables := []struct {
		x        byte
		exponent int
		r        byte
	}{
		{0x53, 2, 0xb5},
		{0x53, 3, 0xc3},
	}

	f := NewField(RijndaelPolynomial, RijndaelGenerator)

	for _, table := range tables {
		poweredX := byte(0x1)
		poweredX = f.Mul(poweredX, table.x)
		for i := 1; i < table.exponent; i++ {
			poweredX = f.Mul(poweredX, table.x)
		}
		if poweredX != table.r {
			t.Errorf("Result of %d^%d wrong, got: %d, want: %d.", table.x, table.exponent, poweredX, table.r)
		}
	}
}
