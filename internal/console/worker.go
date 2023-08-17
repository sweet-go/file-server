package console

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/sweet-go/file-server/internal/config"
	"github.com/sweet-go/file-server/internal/db"
	"github.com/sweet-go/file-server/internal/repository"
	"github.com/sweet-go/file-server/internal/worker"
)

var workerCmd = &cobra.Command{
	Use:  "worker",
	Long: "Starts the worker",
	Run:  workerFn,
}

func init() {
	RootCMD.AddCommand(workerCmd)
}

func workerFn(_ *cobra.Command, _ []string) {
	wrk, err := worker.NewServer(config.WorkerBrokerHost(), config.WorkerConcurency())
	if err != nil {
		logrus.WithError(err).Error("failed to create new worker")
	}

	db.InitializePostgresConn()

	deletableMediaRepo := repository.NewDeletableMediaRepository(db.PostgresDB)

	wrk.RegisterTaskHandler(deletableMediaRepo)

	if err := wrk.RegisterScheduledMediaDeleteScheduler(); err != nil {
		logrus.WithError(err).Error("failed to register scheduled media delete scheduler")
		os.Exit(1)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	errch := make(chan error)

	wrk.Start(errch)
	select {
	case sig := <-sigCh:
		logrus.Infof("receiving signal to stop worker server from console: %s. Gracefully shutting down worker", sig.String())
		wrk.Stop()
		os.Exit(0)
	case err := <-errch:
		logrus.WithError(err).Error("receiving quit signal from worker server. Gracefully shutting down worker")
		wrk.Stop()
		os.Exit(1)
	}
}
