package s4

import (
	"crypto/rand"
	"errors"

	"github.com/ceriath/rsa-shamir-secret-sharing/gf256"
)

var usedXValues []byte

// Split splits a secret into n shares where the threshold k applies
func Split(secret []byte, k, n byte) ([]Share, error) {

	usedXValues = make([]byte, 0)

	if k > n {
		return nil, errors.New("the threshold can not be greater than the shares to generate")
	}

	shares := make([]Share, n)

	// make n shares with random x values != 0
	for i := byte(0); i < n; i++ {
		shares[i] = Share{
			X:              getXValue(),
			Values:         make([]byte, 0),
			Index:          i,
			RequiredShares: k,
		}
	}

	// split the secret into single bytes and create a new polynomial for each byte
	for _, singleSecretByte := range secret {
		poly := generatePolynomial(singleSecretByte, k)

		// calculate f(x) for each share and append the resulting y to the share
		for i := 0; i < len(shares); i++ {
			shareValue := calculateForPolynomial(poly, shares[i].X)
			shares[i].Values = append(shares[i].Values, shareValue)
		}
	}

	return shares, nil
}

func calculateForPolynomial(poly []byte, x byte) byte {
	result := byte(0)
	field := gf256.NewField(gf256.RijndaelPolynomial, gf256.RijndaelGenerator)

	// calculate each coefficient * x ^ exponent in the given polynomial and add it to the result
	for exponent, coefficient := range poly {
		if exponent == 0 {
			result = field.Add(result, coefficient)
			continue
		}

		// calculate x ^ exponent by multiplying x exponent times with itself
		poweredX := x
		for i := 1; i < exponent; i++ {
			poweredX = field.Mul(poweredX, x)
		}

		// multiply the resulting x ^ exponent with the coefficient
		totalX := field.Mul(coefficient, poweredX)

		// add it to the result of the polynomial and continue with next coefficient/exponent pair
		result = field.Add(result, totalX)
	}

	return result
}

func generatePolynomial(secret byte, k byte) []byte {
	coefficients := make([]byte, 1)
	coefficients[0] = secret

	// generate a random coefficient for each non-zero exponent (1 to k-1)
	for i := byte(1); i < k; i++ {
		coefficients = append(coefficients, getRandomByte())
	}

	return coefficients
}

func getXValue() byte {
	x := getRandomByte()

	// make sure the value is not used yet and to never ever get a zero since that would reveal the secret
	for contains(usedXValues, x) || x == 0x0 {
		x = getRandomByte()
	}

	usedXValues = append(usedXValues, x)
	return x
}

func contains(slice []byte, x byte) bool {
	for _, v := range slice {
		if v == x {
			return true
		}
	}
	return false
}

func getRandomByte() byte {
	b := make([]byte, 1)

	_, err := rand.Read(b)
	if err != nil {
		panic("error getting a secure random")
	}

	return b[0]
}
