package provider

import (
	"github.com/itsLeonB/billsplittr/internal/message"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/meq"
	"github.com/rotisserie/eris"
)

type Queues struct {
	OrphanedBillCleanup      meq.TaskQueue[message.OrphanedBillCleanup]
	ExpenseBillTextExtracted meq.TaskQueue[message.ExpenseBillTextExtracted]
}

func ProvideQueues(logger ezutil.Logger, db meq.DB) (*Queues, error) {
	if db == nil {
		return nil, eris.New("db cannot be nil")
	}
	return &Queues{
		OrphanedBillCleanup:      meq.NewTaskQueue[message.OrphanedBillCleanup](logger, db),
		ExpenseBillTextExtracted: meq.NewTaskQueue[message.ExpenseBillTextExtracted](logger, db),
	}, nil
}
