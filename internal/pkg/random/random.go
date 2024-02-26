package random

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func NewRandomString(n int) (string, error) {
	const chars = "01234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ret := make([]byte, n)
	for i := range ret {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random string: %w", err)
		}
		ret[i] = chars[num.Int64()]
	}
	return string(ret), nil
}
