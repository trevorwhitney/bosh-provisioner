package disk

import (
	"time"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshretry "github.com/cloudfoundry/bosh-utils/retrystrategy"
	"github.com/pivotal-golang/clock"
)

type sfdiskPartitionStrategy struct {
	retryable   boshretry.Retryable
	timeService clock.Clock
	logger      boshlog.Logger
}

func NewSfdiskPartitionStrategy(
	retryable boshretry.Retryable,
	timeService clock.Clock,
	logger boshlog.Logger,
) boshretry.RetryStrategy {
	return &sfdiskPartitionStrategy{
		retryable:   retryable,
		logger:      logger,
		timeService: timeService,
	}
}

func (s *sfdiskPartitionStrategy) Try() error {
	var err error
	var isRetryable bool

	for i := 0; i < 20; i++ {
		s.logger.Debug("attemptRetryStrategy", "Making attempt #%d", i)

		isRetryable, err = s.retryable.Attempt()
		if err == nil {
			return nil
		}

		if !isRetryable {
			return err
		}

		s.timeService.Sleep(3 * time.Second)
	}

	return err
}
