package retryer

import (
	"errors"
	"net"
	"syscall"
	"time"

	"go.uber.org/zap"
)

var (
	logRetryConnection = "cannot connect to the server. next try in %v seconds"
)

type ConnRetryer interface {
	DoWithRetry(retryableFunc func() error) error
}

type ConnRetryerFn interface {
	DoWithRetryFn() error
}

type retryer struct {
	attempts        int
	triesStartAfter time.Duration
	triesGrowsDelta time.Duration
	log             *zap.Logger
}

type retryerFn struct {
	attempts        int
	triesStartAfter time.Duration
	triesGrowsDelta time.Duration
	log             *zap.Logger
	fn              func() error
}

// ConnRetryer retries function call if function returns connection refused error
// retry time grows as: triesStartAfter = triesStartAfter + triesGrowsDelta
func NewConnRetryer(attempts int, triesStartAfter, triesGrowsDelta time.Duration, logger *zap.Logger) *retryer {
	return &retryer{
		attempts:        attempts,
		triesStartAfter: triesStartAfter,
		triesGrowsDelta: triesGrowsDelta,
		log:             logger,
	}
}

// ConnRetryerFn retries function call if function returns connection refused error
// retry time grows as: triesStartAfter = triesStartAfter + triesGrowsDelta
func NewConnRetryerFn(attempts int, triesStartAfter, triesGrowsDelta time.Duration, logger *zap.Logger, fn func() error) *retryerFn {
	return &retryerFn{
		attempts:        attempts,
		triesStartAfter: triesStartAfter,
		triesGrowsDelta: triesGrowsDelta,
		log:             logger,
		fn:              fn,
	}
}

func (r *retryer) DoWithRetry(retryableFunc func() error) error {
	err := retryableFunc()
	if err != nil {
		var netOpErr *net.OpError
		if errors.As(err, &netOpErr) && errors.Is(netOpErr, syscall.ECONNREFUSED) {
			a := r.attempts
			ti := r.triesStartAfter
			for a > 0 && err != nil {
				r.log.Sugar().Infof(logRetryConnection, ti.Seconds())
				time.Sleep(ti)
				err = retryableFunc()
				a--
				ti = ti + r.triesGrowsDelta
			}
		}
	}
	return err
}

func (r *retryerFn) DoWithRetryFn() error {
	err := r.fn()
	if err != nil {
		var netOpErr *net.OpError
		if errors.As(err, &netOpErr) && errors.Is(netOpErr, syscall.ECONNREFUSED) {
			a := r.attempts
			ti := r.triesStartAfter
			for a > 0 && err != nil {
				r.log.Sugar().Infof(logRetryConnection, ti.Seconds())
				time.Sleep(ti)
				err = r.fn()
				a--
				ti = ti + r.triesGrowsDelta
			}
		}
	}
	return err
}
