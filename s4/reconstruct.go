package s4

import (
	"errors"

	"github.com/ceriath/rsa-shamir-secret-sharing/gf256"
)

// Reconstruct takes a set of shares and tries to recover the secret from those shares
func Reconstruct(shares []Share) ([]byte, error) {
	if err := validateInput(shares); err != nil {
		return nil, err
	}
	return interpolate(shares), nil
}

func validateInput(shares []Share) error {
	if len(shares) < 1 {
		return errors.New("No shares provided")
	}

	req := shares[0].RequiredShares
	indexes := make([]int, 0)
	for _, share := range shares {

		// check if all shares have the same constant values
		if req != share.RequiredShares {
			return errors.New("At least one share is corrupted.")
		}

		//check for duplicate shares
		for _, i := range indexes {
			if int(share.Index) == i {
				return errors.New("A share has been provided twice or is corrupted.")
			}
		}
		indexes = append(indexes, int(share.Index))
	}

	return nil
}

func interpolate(shares []Share) []byte {

	result := make([]byte, 0)

	// calculate in gf256
	field := gf256.NewField(gf256.RijndaelPolynomial, gf256.RijndaelGenerator)

	// do an interpolation for each byte seperately
	for byteCounter := 0; byteCounter < len(shares[0].Values); byteCounter++ {
		singleResult := byte(0)

		// recover the byte from given shares
		for i := 0; i < len(shares); i++ {

			// get the y value
			term := shares[i].Values[byteCounter]

			// for each other x value (foreign shares), compute the partial term for x=0, since f(0) is our secret
			for j := 0; j < len(shares); j++ {
				if j != i {

					// 0 - xj
					numerator := shares[j].X
					// xi - xj
					denominator := field.Add(shares[i].X, shares[j].X)

					// (0-xj) / (xi-xj)
					fraction := field.Mul(numerator, field.Inv(denominator))

					// multiply the partial term with the full term for the current share. multiplies with the coefficient the first time. continue with next foreign share
					term = field.Mul(term, fraction)
				}
			}

			// add the result for the current share to the resulting lagrange-term and continue with next share
			singleResult = field.Add(singleResult, term)
		}

		// add the byte to the result and continue with next byte
		result = append(result, singleResult)
	}

	return result

}
