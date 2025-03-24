package shortcode

import "math/rand/v2"

type ShortCode struct {
	length int
}

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func NewShortCode(length int) *ShortCode {
	return &ShortCode{length: length}
}
func (s *ShortCode) GenerateShortCode() string {
	n := len(chars)
	result := make([]byte, s.length)

	for i := range s.length {
		result[i] = chars[rand.IntN(n)]
	}

	return string(result)
}
