package worker

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/itsLeonB/billsplittr/internal/message"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/meq/task"
	"github.com/rotisserie/eris"
)

func expenseBillTextExtractedHandler(expenseSvc service.GroupExpenseService) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		var taskMsg task.Task[message.ExpenseBillTextExtracted]

		if err := json.Unmarshal(t.Payload(), &taskMsg); err != nil {
			return eris.Wrapf(err, "error unmarshaling payload to: %T", taskMsg)
		}

		return expenseSvc.ParseFromBillText(ctx, taskMsg.Message)
	})
}
