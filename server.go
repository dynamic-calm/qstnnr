package qstnnr

import (
	"context"

	"github.com/mateopresacastro/qstnnr/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// To assert implementation
var _ api.QuestionnaireServer = (*server)(nil)

type server struct {
	api.QuestionnaireServer
	service QService
}

func NewServer(data InitialData) (*grpc.Server, error) {
	store, err := NewMemoryStore(data)
	if err != nil {
		return nil, err
	}
	service := NewQstnnrService(store)
	server := &server{service: service}
	grpcsrv := grpc.NewServer()
	api.RegisterQuestionnaireServer(grpcsrv, server)
	return grpcsrv, nil
}

func (s *server) GetQuestions(ctx context.Context, _ *emptypb.Empty) (*api.GetQuestionsResponse, error) {
	qsts, err := s.service.GetQuestions()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get questions")
	}
	var questions []*api.Question
	for qID, q := range qsts {
		var options []*api.Option
		var question = &api.Question{Id: int32(qID), Text: q.Text, Options: options}
		for oID, o := range q.Options {
			question.Options = append(question.Options, &api.Option{Id: int32(oID), Text: o.Text})
		}
		questions = append(questions, question)
	}

	return &api.GetQuestionsResponse{Questions: questions}, nil
}

func (s *server) SubmitAnswers(ctx context.Context, req *api.SubmitAnswersRequest) (*api.SubmitAnswersResponse, error) {
	answers := make(map[QuestionID]OptionID)
	for _, a := range req.Answers {
		answers[QuestionID(a.QustionId)] = OptionID(a.OptionId)
	}
	result, err := s.service.SubmitAnswers(answers)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to process submission") // TODO: error boundaries
	}

	qsts, err := s.service.GetQuestions()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get questions")
	}

	var solutions []*api.Solution
	for qID, oID := range result.Solutions {
		q := &api.Question{Id: int32(qID), Text: qsts[qID].Text}
		s := &api.Solution{Question: q, CorrectOptionId: int32(oID), CorrectOptionText: q.Text}
		solutions = append(solutions, s)
	}

	return &api.SubmitAnswersResponse{Solutions: solutions, Stats: int64(result.Stat)}, nil
}

func (s *server) GetSolutions(ctx context.Context, req *emptypb.Empty) (*api.GetSolutionsResponse, error) {
	solutions, err := s.service.GetSolutions()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get solutions")
	}
	processed, err := s.processSolutions(solutions)
	return &api.GetSolutionsResponse{Solutions: processed}, nil
}

func (s *server) processSolutions(ss map[QuestionID]OptionID) ([]*api.Solution, error) {
	qsts, err := s.service.GetQuestions()
	if err != nil {
		return nil, err
	}
	var processed []*api.Solution
	for qID, oID := range ss {
		q := &api.Question{Id: int32(qID), Text: qsts[qID].Text}
		s := &api.Solution{
			Question: q, CorrectOptionId: int32(oID),
			CorrectOptionText: qsts[qID].Options[oID].Text,
		}
		processed = append(processed, s)
	}
	return processed, nil
}
