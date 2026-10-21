package main

import (
	"context"
	"errors"
	"fastcom_cli/core/measure"
	"fastcom_cli/core/view"
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
	URL_COUNT   = 5
	BUFFER_SIZE = 10
	TIMEOUT     = 30 * time.Second
	TICK        = 200 * time.Millisecond
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ctxT, cancel := context.WithTimeout(ctx, TIMEOUT)
	defer cancel()

	m := measure.NewSpeedMeasure(
		API_URL,
		TOKEN,
		URL_COUNT,
		BUFFER_SIZE,
		TICK,
	)

	go m.Run(ctxT)

	vm := view.NewViewModel(m.GetCurrentSpeed, m.GetAVGSpeed, TICK)

	p := tea.NewProgram(vm, tea.WithContext(ctxT))
	if _, err := p.Run(); err != nil && !errors.Is(err, tea.ErrProgramKilled) {
		log.Fatal(err)
	}
}
