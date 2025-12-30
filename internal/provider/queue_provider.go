package provider

import (
	"github.com/itsLeonB/billsplittr/internal/message"
	"github.com/itsLeonB/billsplittr/internal/pkg/logger"
	"github.com/itsLeonB/meq"
)

type Queues struct {
	OrphanedBillCleanup      meq.TaskQueue[message.OrphanedBillCleanup]
	ExpenseBillTextExtracted meq.TaskQueue[message.ExpenseBillTextExtracted]
}

func ProvideQueues(db meq.DB) *Queues {
	return &Queues{
		OrphanedBillCleanup:      meq.NewTaskQueue[message.OrphanedBillCleanup](logger.Global, db),
		ExpenseBillTextExtracted: meq.NewTaskQueue[message.ExpenseBillTextExtracted](logger.Global, db),
	}
}
