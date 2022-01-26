package scheduler

import (
	"testing"
	"time"

	"bilalekrem.com/certstore/internal/assert"
	"bilalekrem.com/certstore/internal/logging"
)

func TestSchedule(t *testing.T) {
	mockTimeProvider := func() time.Time {
		return time.Date(2022, 01, 01, 01, 59, 55, 0, time.Local)
	}

	scheduler := &dailyScheduler{timeProvider: mockTimeProvider}

	called := false
	fn := func() {
		logging.GetLogger().Error("running function")
		called = true
	}

	err := scheduler.Schedule(fn)
	assert.NotError(t, err, "scheduling failed")

	time.Sleep(6 * time.Second)
	assert.True(t, called)
}

func TestScheduleAlreadyScheduled(t *testing.T) {
	scheduler := NewDailyScheduler()
	err := scheduler.Schedule(func() {})
	assert.NotError(t, err, "scheduling failed")

	err = scheduler.Schedule(func() {})
	assert.ErrorContains(t, err, "already scheduled")
}

func TestScheduleValidateNotRunBeforeTimeIsUp(t *testing.T) {
	// fn will run after 5 seconds
	mockTimeProvider := func() time.Time {
		return time.Date(2022, 01, 01, 01, 59, 55, 0, time.Local)
	}

	scheduler := &dailyScheduler{timeProvider: mockTimeProvider}

	called := false
	fn := func() {
		logging.GetLogger().Error("running function")
		called = true
	}

	err := scheduler.Schedule(fn)
	assert.NotError(t, err, "scheduling failed")

	time.Sleep(3 * time.Second)
	assert.False(t, called) // should not have run at this moment

	time.Sleep(2 * time.Second)
	assert.True(t, called)
}
