package services

// Use this for resolvers business logic

// GenerateStatementForUserByTimestamp // Query (calls to transaction-service)

import (
	"context"
	"finnbank/common/grpc/statement"
	"finnbank/common/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)


type StatementService struct {
	db *pgxpool.Pool
	st statement.StatementServiceClient
	l *utils.Logger
}

func NewStatementService(db *pgxpool.Pool, logger *utils.Logger, st statement.StatementServiceClient) *StatementService {
	return &StatementService{
		db:     db,
		l: logger,
		st: st,
	}
}

func (s *StatementService) GenerateStatementForUser(creditId int, debit int, savings int, startDate string, endDate string, ctx *context.Context) ([]byte, error) {
	
	stRes, err := s.st.GenerateStatement(*ctx, &statement.ClientRequest{
		CreditId:  int32(creditId),
		DebitId:   int32(debit),
		SavingsId: int32(savings),
		StartDate: startDate,
		EndDate:   endDate,
	}) 

	if err != nil {
		s.l.Error("Error generating statement %v", err)
		return nil, err
	}
	return stRes.PdfBuffer, nil
}