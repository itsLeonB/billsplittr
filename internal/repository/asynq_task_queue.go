package repository

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/rotisserie/eris"
)

type asynqTaskQueue struct {
	logger ezutil.Logger
	client *asynq.Client
}

func NewTaskQueue(
	logger ezutil.Logger,
	client *asynq.Client,
) TaskQueue {
	return &asynqTaskQueue{
		logger,
		client,
	}
}

func (tq *asynqTaskQueue) Enqueue(ctx context.Context, task entity.Task) error {
	payload, err := json.Marshal(task)
	if err != nil {
		return eris.Wrap(err, "error marshaling task")
	}

	asynqTask := asynq.NewTask(task.Type, payload)

	info, err := tq.client.EnqueueContext(ctx, asynqTask, asynq.Queue(task.Type))
	if err != nil {
		return eris.Wrap(err, "error enqueuing task")
	}

	tq.logger.Infof("enqueued task: BatchID=%s, Queue=%s", info.ID, info.Queue)

	return nil
}

func (tq *asynqTaskQueue) Ping() error {
	if err := tq.client.Ping(); err != nil {
		return eris.Wrap(err, "asynq client is not ready")
	}
	return nil
}
