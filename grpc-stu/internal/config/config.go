package config

type (
	Config struct {
		MethodConfig []*MethodConfig `json:"methodConfig,omitempty"`
	}

	MethodConfig struct {
		Name        []*NameConfig `json:"name,omitempty"`
		Timeout     string        `json:"timeout,omitempty"`
		RetryPolicy *RetryPolicy  `json:"retryPolicy,omitempty"`
	}

	NameConfig struct {
		Service string `json:"service,omitempty"`
		Method  string `json:"method"`
	}

	RetryPolicy struct {
		MaxAttempts          int      `json:"maxAttempts,omitempty"`
		InitialBackoff       string   `json:"initialBackoff,omitempty"`
		MaxBackoff           string   `json:"maxBackoff,omitempty"`
		BackoffMultiplier    int      `json:"backoffMultiplier,omitempty"`
		RetryableStatusCodes []string `json:"retryableStatusCodes,omitempty"`
	}
)
