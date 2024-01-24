package retryer

import (
	"errors"
	"fmt"
	"net"
	"syscall"
	"time"

	"github.com/Chystik/runtime-metrics/pkg/logger"
)

var (
	logRetryConnection = "cannot connect to the server. next try in %v seconds"
)

type retryer struct {
	attempts        int
	triesStartAfter time.Duration
	triesGrowsDelta time.Duration
	log             logger.Logger
}

type retryerFn struct {
	attempts        int
	triesStartAfter time.Duration
	triesGrowsDelta time.Duration
	log             logger.Logger
	fn              func() error
}

// ConnRetryer retries function call if function returns connection refused error
// retry time grows as: triesStartAfter = triesStartAfter + triesGrowsDelta
func NewConnRetryer(attempts int, triesStartAfter, triesGrowsDelta time.Duration, logger logger.Logger) *retryer {
	return &retryer{
		attempts:        attempts,
		triesStartAfter: triesStartAfter,
		triesGrowsDelta: triesGrowsDelta,
		log:             logger,
	}
}

// ConnRetryerFn retries function call if function returns connection refused error
// retry time grows as: triesStartAfter = triesStartAfter + triesGrowsDelta
func NewConnRetryerFn(attempts int, triesStartAfter, triesGrowsDelta time.Duration, logger logger.Logger, fn func() error) *retryerFn {
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
				r.log.Info(fmt.Sprintf(logRetryConnection, ti.Seconds()))
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
				r.log.Info(fmt.Sprintf(logRetryConnection, ti.Seconds()))
				time.Sleep(ti)
				err = r.fn()
				a--
				ti = ti + r.triesGrowsDelta
			}
		}
	}
	return err
}
