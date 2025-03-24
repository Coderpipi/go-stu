package main

import (
	"context"

	"url-shorter/cmd"
)

func main() {
	ctx := context.Background()
	if err := cmd.Web(ctx); err != nil {
		panic(err)
	}
}
