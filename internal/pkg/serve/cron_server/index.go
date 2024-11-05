package cron_server

import (
	_dateUtils "github.com/deeptest-com/deeptest-next/pkg/libs/date"
	_logUtils "github.com/deeptest-com/deeptest-next/pkg/libs/log"
	"gorm.io/gorm"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	once sync.Once
	cc   *cron.Cron
)

type CronServer struct {
	syncMap sync.Map

	CheckEmailJob *CheckEmailJob `inject:""`
}

func NewCronServer() (server *CronServer) {
	server = &CronServer{}

	return
}
func (c CronServer) Start() {
	c.CheckEmailJob.AddTask()

	//GetCronInstance().Start()
}

type CheckEmailJob struct {
	DB *gorm.DB `inject:""`
}

func (j CheckEmailJob) Run() {
	_logUtils.Infof("%s - check email", _dateUtils.TimeStr(time.Now()))
}

func (c CheckEmailJob) AddTask() (err error) {
	spec := "*/6 * * * * *"

	_, err = GetCronInstance().AddJob(spec, c)
	if err != nil {
		return err
	}

	return
}

//func DoOnce(job cron.Job, t ...time.Duration) (err error) {
//	once := time.Now().Add(2 * time.Second)
//	if len(t) == 1 {
//		once = time.Now().Add(t[0] * time.Second)
//	}
//
//	onceSpec := fmt.Sprintf("%d %d %d %d %d %d", once.Second(), once.Minute(), once.Hour(), once.Day(), once.Month(), once.Weekday())
//
//	_, err = GetCronInstance().AddJob(onceSpec, job)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

func GetCronInstance() *cron.Cron {
	once.Do(func() {
		cc = cron.New(cron.WithSeconds())
	})
	return cc
}
