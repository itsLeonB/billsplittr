package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/itsLeonB/billsplittr-protos/gen/go/expensebill/v1"
	"github.com/itsLeonB/billsplittr-protos/gen/go/expenseitem/v1"
	"github.com/itsLeonB/billsplittr-protos/gen/go/groupexpense/v1"
	groupexpensev2 "github.com/itsLeonB/billsplittr-protos/gen/go/groupexpense/v2"
	"github.com/itsLeonB/billsplittr-protos/gen/go/otherfee/v1"
	"github.com/itsLeonB/billsplittr/internal/provider"
	"github.com/rotisserie/eris"
	"google.golang.org/grpc"
)

type Servers struct {
	groupExpense   groupexpense.GroupExpenseServiceServer
	expenseItem    expenseitem.ExpenseItemServiceServer
	otherFee       otherfee.OtherFeeServiceServer
	expenseBill    expensebill.ExpenseBillServiceServer
	groupExpenseV2 groupexpensev2.GroupExpenseServiceServer
}

func ProvideServers(services *provider.Services) *Servers {
	if services == nil {
		panic("services is nil")
	}

	validate := validator.New()

	return &Servers{
		groupExpense:   newGroupExpenseServer(validate, services.GroupExpense),
		expenseItem:    newExpenseItemServer(validate, services.ExpenseItem),
		otherFee:       newOtherFeeServer(validate, services.OtherFee),
		expenseBill:    newExpenseBillServer(validate, services.ExpenseBill),
		groupExpenseV2: newGroupExpenseServerV2(validate, services.GroupExpense),
	}
}

func (s *Servers) Register(grpcServer *grpc.Server) error {
	if s.groupExpense == nil {
		return eris.New("group expense server is nil")
	}
	if s.expenseItem == nil {
		return eris.New("expense item server is nil")
	}
	if s.otherFee == nil {
		return eris.New("other fee server is nil")
	}
	if s.expenseBill == nil {
		return eris.New("expense bill server is nil")
	}

	groupexpense.RegisterGroupExpenseServiceServer(grpcServer, s.groupExpense)
	expenseitem.RegisterExpenseItemServiceServer(grpcServer, s.expenseItem)
	otherfee.RegisterOtherFeeServiceServer(grpcServer, s.otherFee)
	expensebill.RegisterExpenseBillServiceServer(grpcServer, s.expenseBill)
	groupexpensev2.RegisterGroupExpenseServiceServer(grpcServer, s.groupExpenseV2)

	return nil
}
