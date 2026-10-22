package measure

import (
	"context"
	"encoding/json"
	"errors"
	"fastcom_cli/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type SpeedMeasure struct {
	apiUrl       string
	token        string
	targetsCount uint8
	tick         time.Duration
	//
	httpClient  *http.Client
	buf         *utils.FixedByteBuf
	sizeCounter atomic.Uint32
	wg          sync.WaitGroup
	adtIndex    int
	targets     []string
	jobs        chan string
	//
	allMeasuresCount int
	allMeasures      float64
	currentSpeed     float64
	avgSpeed         float64
}

func NewSpeedMeasure(
	apiUrl string,
	token string,
	targetsCount uint8,
	bufSize int8,
	tick time.Duration,
) *SpeedMeasure {
	return &SpeedMeasure{
		apiUrl:       apiUrl,
		token:        token,
		targetsCount: targetsCount,
		tick:         tick,
		httpClient: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        int(targetsCount),
				MaxIdleConnsPerHost: int(targetsCount),
				IdleConnTimeout:     90 * time.Second,
			},
		},
		//
		buf:      utils.NewFixedByteBuf(bufSize),
		adtIndex: -1,
	}
}

func (sm *SpeedMeasure) Run(ctx context.Context) {
	targets, err := sm.findTargets(ctx)
	if err != nil {
		log.Panic(err)
	}
	sm.targets = targets

	sm.runWorkers(ctx)

	sm.runTickMeasure(ctx)

	sm.wg.Wait()
	close(sm.jobs)
}

func (sm *SpeedMeasure) GetCurrentSpeed() float64 {
	return sm.currentSpeed
}

func (sm *SpeedMeasure) GetAVGSpeed() float64 {
	return sm.avgSpeed
}

func (sm *SpeedMeasure) findTargets(ctx context.Context) ([]string, error) {
	type re_s struct {
		Targets []struct {
			Name string `json:"name"`
		} `json:"targets"`
	}

	tUrl, err := url.Parse(sm.apiUrl)
	if err != nil {
		return nil, err
	}

	q := tUrl.Query()
	q.Add("https", "true")
	q.Add("token", sm.token)
	q.Add("urlCount", fmt.Sprintf("%d", sm.targetsCount))

	tUrl.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", tUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	response, err := sm.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode > 399 {
		return nil, errors.New(string(body))
	}

	var dest re_s

	if err := json.Unmarshal(body, &dest); err != nil {
		return nil, err
	}

	res := make([]string, 0, len(dest.Targets))

	for _, row := range dest.Targets {
		u, _ := url.Parse(row.Name)
		u.Path += "/range/0-26214400"
		res = append(res, u.String())
	}

	return res, nil
}

func (sm *SpeedMeasure) runWorkers(ctx context.Context) {
	var i uint8

	sm.jobs = make(chan string, sm.targetsCount)

	for ; i < sm.targetsCount; i++ {
		sm.wg.Add(1)
		go sm.downloadWorker(ctx)
	}

	for _, tar := range sm.targets {
		sm.jobs <- tar
	}
}

func (sm *SpeedMeasure) downloadWorker(ctx context.Context) {
	defer sm.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case j, ok := <-sm.jobs:
			if !ok {
				continue
			}

			req, err := http.NewRequestWithContext(ctx, "GET", j, nil)
			if err != nil {
				continue
			}

			response, err := sm.httpClient.Do(req)
			if err != nil {
				continue
			}

			_, err = io.Copy(io.Discard, &utils.Reader{
				Target: j,
				R:      response.Body,
				OnEof:  sm.onEof,
				OnRead: sm.onRead,
			})
			response.Body.Close()
			if err != nil {
				continue
			}
		}
	}
}

func (sm *SpeedMeasure) runTickMeasure(ctx context.Context) {
	sm.wg.Add(1)
	go func() {
		defer sm.wg.Done()
		ticker := time.NewTicker(sm.tick)
		defer ticker.Stop()

		last := time.Now()

		for {
			select {
			case <-ctx.Done():
				return

			case <-ticker.C:
				size := sm.sizeCounter.Swap(0)
				now := time.Now()
				elapsed := now.Sub(last).Seconds()
				last = now

				sm.buf.Add((float64(size) * 8 / elapsed) / 1_000_000)

				var sum float64
				for _, row := range sm.buf.Arr {
					sum += row
				}

				sm.currentSpeed = sum / float64(sm.buf.ASize)

				sm.allMeasuresCount += 1
				sm.allMeasures += sm.currentSpeed
				sm.avgSpeed = sm.allMeasures / float64(sm.allMeasuresCount)
			}

		}
	}()
}

func (sm *SpeedMeasure) onEof(url string) {
	sm.jobs <- url
}

func (sm *SpeedMeasure) onRead(n int) {
	sm.sizeCounter.Add(uint32(n))
}
