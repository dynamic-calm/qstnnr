package qstnnr_test

import (
	"context"
	"log/slog"
	"net"
	"testing"

	"github.com/mateopresacastro/qstnnr"
	"github.com/mateopresacastro/qstnnr/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestServer(t *testing.T) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	clientOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.NewClient(ln.Addr().String(), clientOpts...)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	questions := map[qstnnr.QuestionID]qstnnr.Question{
		1: {
			ID:   1,
			Text: "What is the capital of France?",
			Options: map[qstnnr.OptionID]qstnnr.Option{
				1: {ID: 1, Text: "London"},
				2: {ID: 2, Text: "Paris"},
				3: {ID: 3, Text: "Berlin"},
				4: {ID: 4, Text: "Madrid"},
			},
		},
		2: {
			ID:   2,
			Text: "Which planet is known as the Red Planet?",
			Options: map[qstnnr.OptionID]qstnnr.Option{
				1: {ID: 1, Text: "Venus"},
				2: {ID: 2, Text: "Mars"},
				3: {ID: 3, Text: "Jupiter"},
				4: {ID: 4, Text: "Saturn"},
			},
		},
		3: {
			ID:   3,
			Text: "What is 2 + 2?",
			Options: map[qstnnr.OptionID]qstnnr.Option{
				1: {ID: 1, Text: "3"},
				2: {ID: 2, Text: "4"},
				3: {ID: 3, Text: "5"},
				4: {ID: 4, Text: "6"},
			},
		},
	}

	solutions := map[qstnnr.QuestionID]qstnnr.OptionID{
		1: 2, // Paris
		2: 2, // Mars
		3: 2, // 4
	}

	data := qstnnr.InitialData{
		Questions: questions,
		Solutions: solutions,
	}

	store, err := qstnnr.NewMemoryStore(data)
	if err != nil {
		t.Fatal(err)
	}

	service := qstnnr.NewQstnnrService(store)

	cfg := &qstnnr.ServerConfig{
		Logger:  slog.Default(),
		Service: service,
	}

	server, err := qstnnr.NewServer(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer server.GracefulStop()

	go func() {
		if err := server.Serve(ln); err != nil {
			t.Error(err)
		}
	}()

	client := api.NewQuestionnaireClient(conn)
	ctx := context.Background()

	t.Run("Should get the questions", func(t *testing.T) {
		resp, err := client.GetQuestions(ctx, &emptypb.Empty{})
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Questions) != 3 {
			t.Errorf("expected 3 question, got %d", len(resp.Questions))
		}
		if resp.Questions[0].Id != 1 {
			t.Errorf("expected question ID 1, got %d", resp.Questions[0].Id)
		}
		if resp.Questions[0].Text != "What is the capital of France?" {
			t.Errorf("expected question 'What is the capital of France?', got '%s'", resp.Questions[0].Text)
		}
	})

	t.Run("Should error if we don't send the correct number of answers", func(t *testing.T) {
		_, err := client.SubmitAnswers(ctx, &api.SubmitAnswersRequest{
			Answers: []*api.Answer{
				{
					QuestionId: 1,
					OptionId:   2,
				},
				{
					QuestionId: 2,
					OptionId:   1,
				},
			},
		})

		if err == nil {
			t.Fatal("expected error since we are sending 2 answers only")
		}
	})

	t.Run("Should submit answers", func(t *testing.T) {
		resp, err := client.SubmitAnswers(ctx, &api.SubmitAnswersRequest{
			Answers: []*api.Answer{
				{
					QuestionId: 1,
					OptionId:   2,
				},
				{
					QuestionId: 2,
					OptionId:   1,
				},
				{
					QuestionId: 3,
					OptionId:   1,
				},
			},
		})

		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Solutions) != 3 {
			t.Errorf("expected 3 solution, got %d", len(resp.Solutions))
		}
		if resp.BetterThan != 100 {
			t.Errorf("expected stats 100, got %d", resp.BetterThan)
		}
		if resp.Solutions[0].CorrectOptionId != 2 {
			t.Errorf("expected correct option ID 2, got %d", resp.Solutions[0].CorrectOptionId)
		}

		for _, sol := range resp.Solutions {
			switch sol.Question.Id {
			case 1:
				if sol.CorrectOptionId != 2 {
					t.Errorf("question 1: expected correct option ID 2, got %d", sol.CorrectOptionId)
				}
			case 2:
				if sol.CorrectOptionText != "Mars" {
					t.Errorf("question 2: expected correct option 'Mars', got '%s'", sol.CorrectOptionText)
				}
			case 3:
				if sol.CorrectOptionText != "4" {
					t.Errorf("question 3: expected correct option '4', got '%s'", sol.CorrectOptionText)
				}
			default:
				t.Errorf("unexpected question ID: %d", sol.Question.Id)
			}
		}
	})

	t.Run("Should get solutions", func(t *testing.T) {
		resp, err := client.GetSolutions(ctx, &emptypb.Empty{})
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Solutions) != 3 {
			t.Errorf("expected 3 solution, got %d", len(resp.Solutions))
		}

		for _, sol := range resp.Solutions {
			switch sol.Question.Id {
			case 1:
				if sol.CorrectOptionId != 2 {
					t.Errorf("question 1: expected correct option ID 2, got %d", sol.CorrectOptionId)
				}
			case 2:
				if sol.CorrectOptionText != "Mars" {
					t.Errorf("question 2: expected correct option 'Mars', got '%s'", sol.CorrectOptionText)
				}
			case 3:
				if sol.CorrectOptionText != "4" {
					t.Errorf("question 3: expected correct option '4', got '%s'", sol.CorrectOptionText)
				}
			default:
				t.Errorf("unexpected question ID: %d", sol.Question.Id)
			}
		}
	})
}
