package retry

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"time"
)

func Do(fn func() error, opts ...Option) error {
	fn_ := func() (any, error) {
		return nil, fn()
	}

	_, err := DoWithData(fn_, opts...)
	return err
}

func DoWithData[T any](fn func() (T, error), opts ...Option) (T, error) {
	var t T

	cfg := newDefaultConfig()
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return t, err
		}
	}

	var err error
	for attempt := 0; attempt < cfg.maxAttempts; attempt++ {
		t, err = fn()
		if err == nil {
			return t, nil
		}

		if !cfg.retryIf(err) {
			break
		}

		// don't wait if this was the last attempt
		if attempt == cfg.maxAttempts-1 {
			break
		}

		select {
		case <-time.After(cfg.backoff.Delay(attempt)):
			continue
		case <-cfg.context.Done():
			err = cfg.context.Err()
			break
		}
	}

	return t, err
}

type Config struct {
	maxAttempts int
	retryIf     func(error) bool
	backoff     *Backoff
	context     context.Context
}

func newDefaultConfig() *Config {
	return &Config{
		maxAttempts: 5,
		retryIf:     func(error) bool { return true },
		backoff: &Backoff{
			InitialDelay: 1 * time.Second,
			Factor:       2.0,
			Jitter:       0.1,
			MaxDelay:     120 * time.Second,
		},
		context: context.Background(),
	}
}

type Option func(*Config) error

func MaxAttempts(n int) Option {
	return func(c *Config) error {
		if n < 0 {
			return errors.New("Maximum retries must be non-negative")
		}
		c.maxAttempts = n
		return nil
	}
}

func RetryIf(fn func(error) bool) Option {
	return func(c *Config) error {
		c.retryIf = fn
		return nil
	}
}

func WithContext(ctx context.Context) Option {
	return func(c *Config) error {
		if ctx == nil {
			return errors.New("Context must not be nil")
		}
		c.context = ctx
		return nil
	}
}

func WithBackoff(b *Backoff) Option {
	return func(c *Config) error {
		c.backoff = b
		return nil
	}
}

type Backoff struct {
	// How long to wait before first retry
	InitialDelay time.Duration
	// Upper bound on backoff
	MaxDelay time.Duration
	// Factor with which to multiply backoff after a failed retry
	Factor float64
	// By how much to randomize backoff
	Jitter float64
}

func (b *Backoff) Delay(attempt int) time.Duration {
	if b == nil {
		return time.Duration(0)
	}

	if attempt < 0 {
		return b.InitialDelay
	}

	delay := math.Pow(b.Factor, float64(attempt)) * float64(b.InitialDelay)
	if delay > float64(b.MaxDelay) {
		delay = float64(b.MaxDelay)
	}

	// if jitter is 0.1, then multiply delay by random value in [0.9, 1.1)
	r := -1 + 2*rand.Float64() // pseudo-random number in [-1, 1)
	delay = delay * (1 + b.Jitter*r)

	return time.Duration(delay)
}
