package worker

import (
	"time"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
	"github.com/sweet-go/file-server/internal/config"
	"github.com/sweet-go/file-server/model"
	workerPkg "github.com/sweet-go/stdlib/worker"
)

type worker struct {
	server workerPkg.Server
}

// Server defines worker server interface
type Server interface {
	Start(errch chan error)
	Stop()
	RegisterScheduledMediaDeleteScheduler() error
	RegisterTaskHandler(deletableMediaRepo model.DeletableMediaRepository)
}

// NewServer creates a new worker server
func NewServer(redisHost string, concurrency int) (Server, error) {
	srv, err := workerPkg.NewServer(
		redisHost,
		asynq.Config{
			Concurrency:         concurrency,
			Queues:              workerPkg.DefaultQueue,
			Logger:              logrus.WithField("source", "file server worker"),
			HealthCheckFunc:     workerPkg.DefaultHealtCheckFn,
			HealthCheckInterval: 5 * time.Minute,
			IsFailure:           workerPkg.DefaultIsFailureCheckerFn,
			StrictPriority:      true,
			RetryDelayFunc:      workerPkg.DefaultRetryDelayFn,
		},
		&asynq.SchedulerOpts{
			LogLevel: config.WorkerLogLevel(),
			Logger:   logrus.New(),
			Location: time.UTC,
		},
	)

	if err != nil {
		logrus.WithError(err).Error("failed to create worker server")
		return nil, err
	}

	return &worker{
		server: srv,
	}, nil
}

var mux = asynq.NewServeMux()

func (w *worker) RegisterTaskHandler(deletableMediaRepo model.DeletableMediaRepository) {
	th := newTaskHandler(deletableMediaRepo)
	mux.HandleFunc(string(model.TaskScheduledDelete), th.handleScheduledMediaDeleteTask)
}

func (w *worker) Start(errch chan error) {
	w.server.Start(mux, errch)
}

func (w *worker) Stop() {
	w.server.Stop()
}

func (w *worker) RegisterScheduledMediaDeleteScheduler() error {
	return w.server.RegisterScheduler(asynq.NewTask(string(model.TaskScheduledDelete), nil), config.ScheduledDeleteSchedulerCronspec())
}
