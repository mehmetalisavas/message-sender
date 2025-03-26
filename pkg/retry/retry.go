package retry

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

var DefaultConfig = Config{
	MaxRetries:     5,
	InitialBackoff: 1 * time.Second,
	MaxBackoff:     20 * time.Second,
	BackoffFactor:  2,
}

type Config struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	BackoffFactor  int
}

func Retry(ctx context.Context, fn func() (*http.Response, error), config Config) (*http.Response, error) {
	var finalErr error
	backoff := config.InitialBackoff

	for i := 0; i < config.MaxRetries; i++ {
		resp, err := fn()
		if err == nil {
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return resp, nil
			}
			finalErr = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		} else {
			finalErr = err
		}

		select {
		case <-time.After(backoff):
			backoff *= time.Duration(config.BackoffFactor)
			if backoff > config.MaxBackoff {
				backoff = config.MaxBackoff
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return nil, finalErr
}
