package jobs

import (
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
)

const (
	StatusOk           = 0
	StatusExecuteError = 1
	StatusRunning      = 2
)

var jobsById = make(map[string]*job)
var JobScheduler *gocron.Scheduler = gocron.NewScheduler(time.UTC)

type Job struct {
	Requires   []string   `json:"requires"`
	Executable string     `json:"executable"`
	Args       string     `json:"string"`
	UUID       string     `json:"uuid"`
	Status     int        `json:"status"`
	Expires    *time.Time `json:"expires"`
	StopsAt    *time.Time `json:"stopsAt"`
	StartsAt   *time.Time `json:"startsAt"`
}

type job struct {
	Job
	Process *os.Process
	OnDone  func(j job, err error)
	OnStop  func(j job, err error)
	OnStart func(j job, err error)
}

type JobConfig struct {
	OnDone  func(j job, err error)
	OnStop  func(j job, err error)
	OnStart func(j job, err error)
}

func New(j Job, config JobConfig) *job {
	return &job{
		OnDone:  config.OnDone,
		OnStop:  config.OnStop,
		OnStart: config.OnStart,
	}
}

func (j *job) Execute() {
	command := exec.Command(j.Executable, strings.Split(j.Args, " ")...)
	err := command.Start()
	if err != nil {
		j.Status = StatusExecuteError
		j.OnStop(*j, err)
	}
	j.Status = StatusRunning
	j.Process = command.Process
}

func (j *job) Schedule() {
	if j.StopsAt != nil {
		JobScheduler.At(j.Expires).Do(func() {
			j.Stop()
		})
	}
	if j.StartsAt != nil {
		JobScheduler.At(j.StartsAt).Do(func() {
			j.Execute()
		})
	}
}

func (j *job) Stop() {
	j.OnStop(*j, j.Process.Kill())
}

func Find(id string) *job {
	return jobsById[id]
}
