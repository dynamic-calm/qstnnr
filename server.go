package qstnnr

import (
	"context"

	"github.com/mateopresacastro/qstnnr/api"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// To assert implementation
var _ api.QuestionnaireServer = (*server)(nil)

type server struct {
	api.QuestionnaireServer
}

func NewServer() (*grpc.Server, error) {
	grpcsrv := grpc.NewServer()
	server := &server{}
	api.RegisterQuestionnaireServer(grpcsrv, server)
	return grpcsrv, nil
}

func (s *server) GetQuestions(ctx context.Context, req *emptypb.Empty) (*api.GetQuestionsResponse, error) {
	mockQuestion := &api.Question{
		Id:   1,
		Text: "What is 2+2?",
		Options: []*api.Option{
			{Id: 1, Text: "3"},
			{Id: 2, Text: "4"},
		},
	}
	return &api.GetQuestionsResponse{Questions: []*api.Question{mockQuestion}}, nil
}

func (s *server) SubmitAnswers(ctx context.Context, req *api.SubmitAnswersRequest) (*api.SubmitAnswersResponse, error) {
	mockSolution := &api.Solution{
		Question:          &api.Question{Id: 1, Text: "What is 2+2?"},
		CorrectOptionId:   2,
		CorrectOptionText: "4",
	}
	return &api.SubmitAnswersResponse{Solutions: []*api.Solution{mockSolution}, Stats: 1}, nil
}

func (s *server) GetSolutions(ctx context.Context, req *emptypb.Empty) (*api.GetSolutionsResponse, error) {
	mockSolution := &api.Solution{
		Question:          &api.Question{Id: 1, Text: "What is 2+2?"},
		CorrectOptionId:   2,
		CorrectOptionText: "4",
	}
	return &api.GetSolutionsResponse{Solutions: []*api.Solution{mockSolution}}, nil
}
