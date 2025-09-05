package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/itsLeonB/billsplittr-protos/gen/go/expensebill/v1"
	"github.com/itsLeonB/billsplittr-protos/gen/go/expenseitem/v1"
	"github.com/itsLeonB/billsplittr-protos/gen/go/groupexpense/v1"
	"github.com/itsLeonB/billsplittr-protos/gen/go/otherfee/v1"
	"github.com/itsLeonB/billsplittr/internal/provider"
	"github.com/rotisserie/eris"
	"google.golang.org/grpc"
)

type Servers struct {
	GroupExpense groupexpense.GroupExpenseServiceServer
	ExpenseItem  expenseitem.ExpenseItemServiceServer
	OtherFee     otherfee.OtherFeeServiceServer
	ExpenseBill  expensebill.ExpenseBillServiceServer
}

func ProvideServers(services *provider.Services) *Servers {
	if services == nil {
		panic("services is nil")
	}

	validate := validator.New()

	return &Servers{
		GroupExpense: newGroupExpenseServer(validate, services.GroupExpense),
		ExpenseItem:  newExpenseItemServer(validate, services.ExpenseItem),
		OtherFee:     newOtherFeeServer(validate, services.OtherFee),
		ExpenseBill:  newExpenseBillServer(validate, services.ExpenseBill),
	}
}

func (s *Servers) Register(grpcServer *grpc.Server) error {
	if s.GroupExpense == nil {
		return eris.New("group expense server is nil")
	}
	if s.ExpenseItem == nil {
		return eris.New("expense item server is nil")
	}
	if s.OtherFee == nil {
		return eris.New("other fee server is nil")
	}
	if s.ExpenseBill == nil {
		return eris.New("expense bill server is nil")
	}

	groupexpense.RegisterGroupExpenseServiceServer(grpcServer, s.GroupExpense)
	expenseitem.RegisterExpenseItemServiceServer(grpcServer, s.ExpenseItem)
	otherfee.RegisterOtherFeeServiceServer(grpcServer, s.OtherFee)
	expensebill.RegisterExpenseBillServiceServer(grpcServer, s.ExpenseBill)

	return nil
}
