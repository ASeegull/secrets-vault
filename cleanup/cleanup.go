package cleanup

import (
	"time"

	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

type Manager struct {
	sched *cron.Cron
	db    querier
}

type querier interface {
	ClearExpired() (string, error)
}

func New(db querier, schedule string) *Manager {
	m := &Manager{
		db:    db,
		sched: cron.New(),
	}

	m.sched.AddFunc(schedule, func() {
		res, err := m.db.ClearExpired()
		if err != nil {
			log.Warn(err)
		}
		log.WithFields(log.Fields{"date": time.Now, "job": "cleanup expired secrets"}).Info(res)
	})

	return m
}
