package main

import (
	"context"
	"fastcom_cli/core/measure"
	"time"
)

const (
	API_URL     = "https://api.fast.com/netflix/speedtest/v2"
	TOKEN       = "YXNkZmFzZGxmbnNkYWZoYXNkZmhrYWxm"
	URL_COUNT   = 5
	BUFFER_SIZE = 10
	TIMEOUT     = 20 * time.Second
	TICK        = 200 * time.Millisecond
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	measure.NewSpeedMeasure(
		API_URL,
		TOKEN,
		URL_COUNT,
		BUFFER_SIZE,
		TICK,
	).Run(ctx)
}
