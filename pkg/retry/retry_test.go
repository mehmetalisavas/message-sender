// File: pkg/retry/retry_test.go
package retry

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestRetry_SuccessOnFirstAttempt(t *testing.T) {
	ctx := context.Background()
	fn := func() (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK}, nil
	}

	resp, err := Retry(ctx, fn, DefaultConfig)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}
}

func TestRetry_RetriesAndSucceeds(t *testing.T) {
	ctx := context.Background()
	attempts := 0
	fn := func() (*http.Response, error) {
		attempts++
		if attempts < 3 {
			return nil, errors.New("temporary error")
		}
		return &http.Response{StatusCode: http.StatusOK}, nil
	}

	resp, err := Retry(ctx, fn, DefaultConfig)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}

func TestRetry_ExceedsMaxRetries(t *testing.T) {
	ctx := context.Background()
	fn := func() (*http.Response, error) {
		return nil, errors.New("permanent error")
	}

	resp, err := Retry(ctx, fn, Config{
		MaxRetries:     2,
		InitialBackoff: 1 * time.Second,
		MaxBackoff:     5 * time.Second,
		BackoffFactor:  2,
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if resp != nil {
		t.Fatalf("expected nil response, got %v", resp)
	}
}

func TestRetry_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	fn := func() (*http.Response, error) {
		time.Sleep(100 * time.Millisecond)
		return nil, errors.New("timeout error")
	}

	resp, err := Retry(ctx, fn, DefaultConfig)
	if err == nil || err != context.DeadlineExceeded {
		t.Fatalf("expected context deadline exceeded error, got %v", err)
	}
	if resp != nil {
		t.Fatalf("expected nil response, got %v", resp)
	}
}
func TestRetry_HandlesNon2xxStatusCodes(t *testing.T) {
	ctx := context.Background()
	attempts := 0
	fn := func() (*http.Response, error) {
		attempts++
		if attempts < 3 {
			return &http.Response{StatusCode: http.StatusInternalServerError}, nil
		}
		return &http.Response{StatusCode: http.StatusOK}, nil
	}

	resp, err := Retry(ctx, fn, DefaultConfig)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}
