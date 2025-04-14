package service

import (
	"context"
	"finnbank/common/grpc/statement"
	"finnbank/common/utils"
)

type StatementService struct {
	Logger *utils.Logger
	Grpc   statement.StatementServiceServer
	statement.UnimplementedStatementServiceServer
}

// mustEmbedUnimplementedStatementServiceServer implements statement.StatementServiceServer.

// AddStatement implements statement.StatementServiceServer.
func (s *StatementService) AddStatement(context.Context, *statement.AddStatementRequest) (*statement.AddStatementResponse, error) {
	panic("unimplemented")
}

func (s *StatementService) GetStatement(context.Context, *statement.GetStatementRequest) (*statement.GetStatementResponse, error) {
	panic("unimplemented")
}

// mustEmbedUnimplementedStatementServiceServer implements statement.StatementServiceServer.

// func AddStatement(context.Context, *AddStatementRequest) (*AddStatementResponse, error) {

// }
