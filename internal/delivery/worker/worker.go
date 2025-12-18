package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/itsLeonB/billsplittr/internal/message"
	"github.com/itsLeonB/billsplittr/internal/pkg/config"
	"github.com/itsLeonB/billsplittr/internal/pkg/logger"
	"github.com/itsLeonB/billsplittr/internal/provider"
	"github.com/rotisserie/eris"
)

type Worker struct {
	srv *asynq.Server
	mux *asynq.ServeMux
}

func Setup() (*Worker, error) {
	providers, err := provider.All()
	if err != nil {
		return nil, err
	}

	asynqCfg := asynq.Config{
		Concurrency: 3,
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			if err != nil {
				logger.Global.Errorf("error processing message: %s", eris.ToString(err, true))
			}
		}),
	}

	srv := asynq.NewServer(config.Global.Valkey.ToRedisOpts(), asynqCfg)
	mux := asynq.NewServeMux()

	mux.Handle(message.ExpenseBillUploaded{}.Type(), expenseBillUploadedHandler(providers.Services.ExpenseBill, providers.Queues.ExpenseBillTextExtracted))
	mux.Handle(message.ExpenseBillTextExtracted{}.Type(), expenseBillTextExtractedHandler(providers.Services.GroupExpense))

	if err := srv.Ping(); err != nil {
		return nil, eris.Wrap(err, "error pinging valkey")
	}

	return &Worker{
		srv,
		mux,
	}, nil
}

func (w *Worker) Run() {
	if err := w.srv.Run(w.mux); err != nil {
		logger.Global.Fatalf("error running worker: %v", err)
	}
}
