package jobs

import (
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/humbertovnavarro/krat/pkg/database"
)

const (
	SatusDone            = 0
	StatusRunning        = 1
	StatusStoppedExpired = 2
	StatusStoppedTimeout = 3
	StatusStoppedError   = 4
)

var GarbageCollectionTime time.Duration = time.Minute * 5
var jobsById = make(map[string]*job)

type JobCallback = func(j job)

type job struct {
	Scheduler *gocron.Scheduler
	JobStartConfig
	Process  *os.Process
	Status   int
	OnDone   JobCallback
	OnStop   JobCallback
	OnStart  JobCallback
	OnReject JobCallback
	Output   []byte
}

type JobStartConfig struct {
	Executable string
	Args       string
	UUID       string
	Expires    *time.Time
	StopsAt    *time.Time
	StartsAt   *time.Time
	OnDone     JobCallback
	OnStop     JobCallback
	OnStart    JobCallback
	OnReject   JobCallback
}

// Create a new job using JobStartConfig
func New(config JobStartConfig) *job {
	return &job{
		Scheduler: gocron.NewScheduler(time.UTC),
		OnDone:    config.OnDone,
		OnStop:    config.OnStop,
		OnStart:   config.OnStart,
	}
}

// Execute the binary specified by Executable with Args as arguments. Calls OnStop and sets the exit status.
func (j *job) Execute() {
	command := exec.Command(j.Executable, strings.Split(j.Args, " ")...)
	output, err := command.Output()
	j.Output = output
	if err != nil {
		j.Status = StatusStoppedError
		j.OnStop(*j)
	}
	j.Status = StatusRunning
	j.Process = command.Process
}

// Schedules the job to be ran
func (j *job) Schedule() {
	if time.Now().UnixMilli() < j.Expires.UnixMilli() {
		j.Status = StatusStoppedExpired
		j.Stop()
	}
	if j.StopsAt != nil {
		j.Scheduler.At(j.Expires).Do(func() {
			j.Status = StatusStoppedTimeout
			j.Stop()
		})
	}
	if j.StartsAt != nil {
		j.Scheduler.At(j.StartsAt).Do(func() {
			j.Execute()
		})
	}
}

func (j *job) Stop() {
	if j.Process != nil {
		j.Process.Kill()
	}
	j.OnStop(*j)
	j.Scheduler.Stop()
	if GarbageCollectionTime.Milliseconds() != 0 {
		j.Scheduler.At(time.Now().UnixMilli() + GarbageCollectionTime.Milliseconds()).Do(func() {
			j.destroy()
		})
		return
	}
	j.destroy()
}

func (j *job) destroy() {
	delete(jobsById, j.UUID)
	j = nil
}

func find(id string) *job {
	return jobsById[id]
}

func Find(id string) *job {
	if find(id) == nil {
		return findInDB(id)
	}
	return find(id)
}

func findInDB(id string) *job {
	return nil
}

func (j *job) serialiaze() {
	database.DB.Exec(createJobTableStatement, j.UUID, j.Status, j.Executable, j.Args, j.Expires.UnixMicro(), j.StopsAt.UnixMicro(), j.Output)
}
