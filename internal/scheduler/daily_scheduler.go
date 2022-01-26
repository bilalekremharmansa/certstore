package scheduler

import (
	"errors"
	"time"

	"bilalekrem.com/certstore/internal/logging"
)

type TimeProvider func() time.Time

var DEFAULT_TIME_PROVIDER TimeProvider = func() time.Time {
	return time.Now()
}

// -------

type dailyScheduler struct {
	scheduled bool

	timeProvider TimeProvider
}

func NewDailyScheduler() *dailyScheduler {
	return &dailyScheduler{timeProvider: DEFAULT_TIME_PROVIDER}
}

func (s *dailyScheduler) Schedule(fn func()) error {
	if s.scheduled {
		return errors.New("scheduler is already scheduled..")
	}

	// ----

	providedTime := s.timeProvider()
	_, minutes, seconds := providedTime.Clock()
	minutesForNextHour := 59 - minutes
	secondsForNextHour := 59 - seconds

	timeForNextHourInSeconds := (minutesForNextHour * 60) + secondsForNextHour

	// ----

	go func() {
		logging.GetLogger().Infof("sleeping for %d seconds", timeForNextHourInSeconds)
		time.Sleep(time.Duration(timeForNextHourInSeconds) * time.Second)

		for {
			fn()

			logging.GetLogger().Info("will sleep for a day, until for next iteration")
			time.Sleep(24 * time.Hour)
		}
	}()
	s.scheduled = true

	return nil
}
