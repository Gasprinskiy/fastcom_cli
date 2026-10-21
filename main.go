package main

import (
	"context"
	"errors"
	"fastcom_cli/core/measure"
	"fastcom_cli/core/view"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	API_URL     = "https://api.fast.com/netflix/speedtest/v2"
	TOKEN       = "YXNkZmFzZGxmbnNkYWZoYXNkZmhrYWxm"
	URL_COUNT   = 10
	BUFFER_SIZE = 10
	DURATION    = 30 * time.Second
	TICK        = 200 * time.Millisecond
)

var DURATION_FLAG *time.Duration

func init() {
	DURATION_FLAG = flag.Duration("d", DURATION, "")
	flag.Parse()
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ctxT, cancel := context.WithTimeout(ctx, *DURATION_FLAG)
	defer cancel()

	m := measure.NewSpeedMeasure(
		API_URL,
		TOKEN,
		URL_COUNT,
		BUFFER_SIZE,
		TICK,
	)
	vm := view.NewViewModel(
		m.GetCurrentSpeed,
		m.GetAVGSpeed,
		TICK+10,
	)

	p := tea.NewProgram(vm, tea.WithContext(ctxT))

	go m.Run(ctxT)
	if _, err := p.Run(); err != nil && !errors.Is(err, tea.ErrProgramKilled) {
		log.Fatal(err)
	}
}
