package s4

import (
	"testing"
)

func TestSplit(t *testing.T) {
	tables := []struct {
		secret string
		k, n   byte
	}{
		{
			"-----BEGIN RSA PRIVATE KEY-----MIICXAIBAAKBgQDUEpZIGNkep0QvcPbtdAPOviubJEqFd9yK83wVc+nWOZ7rh0bHPv77EPPltURnOnIu1mglj97BDHDOSzzbCtQrjc0gt5KtLgtra7xooKIJmVdVdcmmeLG+PkaJJqkVjbCvDgsbjrdCySKdOMv5s6OZiMdibBuk7OBm12/3BkI4kQIDAQABAoGAGVSvBqWRKYwpJNGFbQ9ydPtaEgnfrNmISkCTDazuvVvck51w4tOveSWpPELOjNX6VYXeor3wiXaG8t0hw7gjZ2iVWXNl0D7nPC82wbzCC2qYnrWWVri/v1Eo40DkYVd8j0TeRjKTT+nKQnKzd6LyISoQkEblEcNLaiSwEk0vQAECQQDYq8zyU/jheSkR2cQ37LZF2VMfXYLWe97VWp47zjPkrF5gGxUTV9JZLRuSk+XY25/nXaruVopYt2lf2ab2XRORAkEA+pEa+TOYzhe3vZlbK37FX6gWipUSyeWNOt5uMeXtcA0eyw7teA2Zm3LRJuA515V6V/h8xjIuDtLeVJZ6AwZVAQJActZDrqBBBcgAs3xW2kk0pjq0KqiWQvWuOihoO0xkBqww7QENn43pZ+hXd825gcGNs8PaZ64obzLuv6WNL4BPcQJBAPR/z4BINt00C5k5IduJcnNrA8Pqv1C/bFZ7Ya/TGyPxyZB/Wn2BV9h1610yA384Xw+ka6zCnmrcnRKRWzHEZgECQGryukIadTG0NO7vEd1hPXTTjLqSANXiSq4TUWd1raJrC3gMlNS2LCWRmKDwVrw0WyNiHh+XZDdVweQOA/Gzm/8=-----END RSA PRIVATE KEY-----",
			3,
			5,
		},
		{
			"PC",
			3,
			5,
		},
	}

	for _, table := range tables {
		bSecret := []byte(table.secret)
		res, err := Split(bSecret, table.k, table.n)
		if err != nil {
			t.Error(err.Error())
		}
		for i, share := range res {
			t.Logf("Share %d X: %d Values %v B64: %s\n", i, share.X, share.Values, share.GetBase64())
		}

		interpol := interpolate(res[:3])

		if table.secret != string(interpol) {
			t.Errorf("Result: %v Want: %v", string(interpol), table.secret)
		}

		// interpol, err = interpolate(res[1:4])
		// if err != nil {
		// 	t.Error(err.Error())
		// }

		// if table.secret != string(interpol) {
		// 	t.Errorf("Result: %v Want: %v", string(interpol), table.secret)
		// }

		// interpol, err = interpolate(res[2:])
		// if err != nil {
		// 	t.Error(err.Error())
		// }

		// if table.secret != string(interpol) {
		// 	t.Errorf("Result: %v Want: %v", string(interpol), table.secret)
		// }
	}

}

// Note: These tests are built with Reed-Solomon

func TestInterpolate(t *testing.T) {
	tables := []struct {
		Shares []Share
		secret string
	}{
		{
			[]Share{
				Share{
					X: 208,
					Values: []byte{
						47,
					},
				},
				Share{
					X: 242,
					Values: []byte{
						79,
					},
				},
				Share{
					X: 17,
					Values: []byte{
						128,
					},
				},
			},
			"g",
		},
	}

	for _, table := range tables {
		s := interpolate(table.Shares)

		if table.secret != string(s) {
			t.Errorf("Interpolation failed. Got: %v Want: %v", string(s), table.secret)
		}
	}
}

func TestPerformance(t *testing.T) {
	tables := []struct {
		secret string
		k, n   byte
	}{
		{
			"-----BEGIN RSA PRIVATE KEY-----MIICXAIBAAKBgQDUEpZIGNkep0QvcPbtdAPOviubJEqFd9yK83wVc+nWOZ7rh0bHPv77EPPltURnOnIu1mglj97BDHDOSzzbCtQrjc0gt5KtLgtra7xooKIJmVdVdcmmeLG+PkaJJqkVjbCvDgsbjrdCySKdOMv5s6OZiMdibBuk7OBm12/3BkI4kQIDAQABAoGAGVSvBqWRKYwpJNGFbQ9ydPtaEgnfrNmISkCTDazuvVvck51w4tOveSWpPELOjNX6VYXeor3wiXaG8t0hw7gjZ2iVWXNl0D7nPC82wbzCC2qYnrWWVri/v1Eo40DkYVd8j0TeRjKTT+nKQnKzd6LyISoQkEblEcNLaiSwEk0vQAECQQDYq8zyU/jheSkR2cQ37LZF2VMfXYLWe97VWp47zjPkrF5gGxUTV9JZLRuSk+XY25/nXaruVopYt2lf2ab2XRORAkEA+pEa+TOYzhe3vZlbK37FX6gWipUSyeWNOt5uMeXtcA0eyw7teA2Zm3LRJuA515V6V/h8xjIuDtLeVJZ6AwZVAQJActZDrqBBBcgAs3xW2kk0pjq0KqiWQvWuOihoO0xkBqww7QENn43pZ+hXd825gcGNs8PaZ64obzLuv6WNL4BPcQJBAPR/z4BINt00C5k5IduJcnNrA8Pqv1C/bFZ7Ya/TGyPxyZB/Wn2BV9h1610yA384Xw+ka6zCnmrcnRKRWzHEZgECQGryukIadTG0NO7vEd1hPXTTjLqSANXiSq4TUWd1raJrC3gMlNS2LCWRmKDwVrw0WyNiHh+XZDdVweQOA/Gzm/8=-----END RSA PRIVATE KEY-----",
			3,
			5,
		},
	}

	for i := 0; i < 100; i++ {
		for _, table := range tables {
			bSecret := []byte(table.secret)
			res, err := Split(bSecret, table.k, table.n)
			if err != nil {
				t.Error(err.Error())
			}
			// for i, share := range res {
			// 	t.Logf("Share %d X: %d Values %v B64: %s\n", i, share.X, share.Values, share.GetBase64())
			// }

			interpol, err := Reconstruct(res[:3])
			if err != nil {
				t.Error(err.Error())
			}

			if table.secret != string(interpol) {
				t.Errorf(":3 Result: %v Want: %v", string(interpol), table.secret)
				for i, share := range res {
					t.Logf("Share %d X: %d Values %v B64: %s\n", i, share.X, share.Values, share.GetBase64())
				}
			}

			interpol, err = Reconstruct(res[1:4])
			if err != nil {
				t.Error(err.Error())
			}

			if table.secret != string(interpol) {
				t.Errorf("1:4 Result: %v Want: %v", string(interpol), table.secret)
				for i, share := range res {
					t.Logf("Share %d X: %d Values %v B64: %s\n", i, share.X, share.Values, share.GetBase64())
				}
			}

			interpol, err = Reconstruct(res[2:])
			if err != nil {
				t.Error(err.Error())
			}

			if table.secret != string(interpol) {
				t.Errorf("2: Result: %v Want: %v", string(interpol), table.secret)
				for i, share := range res {
					t.Logf("Share %d X: %d Values %v B64: %s\n", i, share.X, share.Values, share.GetBase64())
				}
			}
		}
	}
}

// func TestEntropy(t *testing.T) {
// 	dups := 0
// 	for i := 0; i < 1000; i++ {
// 		usedX := make([]byte, 0)
// 		for j := 0; j < 4; j++ {
// 			b := make([]byte, 1)

// 			_, err := rand.Read(b)
// 			if err != nil {
// 				panic("error getting a secure random")
// 			}

// 			for _, v := range usedX {
// 				if v == b[0] {
// 					dups++
// 				}
// 			}
// 			usedX = append(usedX, b[0])
// 		}
// 	}
// 	t.Logf("found duplicates over 1000: %d", dups)
// }
