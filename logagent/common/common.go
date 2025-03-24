package common

import "github.com/sourcegraph/conc"

var (
	SG = &conc.WaitGroup{}
)

type (
	CollectEntry struct {
		Path  string `json:"path"`
		Topic string `json:"topic"`
	}
)
