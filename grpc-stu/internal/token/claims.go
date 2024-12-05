package token

type Token struct {
	Sub string `mapstructure:"sub" json:"sub,omitempty"`
	Iss string `mapstructure:"iss" json:"iss,omitempty"`
	Aud string `mapstructure:"aud" json:"aud,omitempty"`
}
