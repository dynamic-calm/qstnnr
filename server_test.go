package qstnnr_test

import (
	"context"
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

	server, err := qstnnr.NewServer()
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
		if len(resp.Questions) != 1 {
			t.Errorf("expected 1 question, got %d", len(resp.Questions))
		}
		if resp.Questions[0].Id != 1 {
			t.Errorf("expected question ID 1, got %d", resp.Questions[0].Id)
		}
		if resp.Questions[0].Text != "What is 2+2?" {
			t.Errorf("expected question 'What is 2+2?', got '%s'", resp.Questions[0].Text)
		}
	})

	t.Run("Should submit answers", func(t *testing.T) {
		resp, err := client.SubmitAnswers(ctx, &api.SubmitAnswersRequest{
			Answers: []*api.Answer{
				{
					QustionId: 1,
					OptionId:  2,
				},
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Solutions) != 1 {
			t.Errorf("expected 1 solution, got %d", len(resp.Solutions))
		}
		if resp.Stats != 1 {
			t.Errorf("expected stats 1, got %d", resp.Stats)
		}
		if resp.Solutions[0].CorrectOptionId != 2 {
			t.Errorf("expected correct option ID 2, got %d", resp.Solutions[0].CorrectOptionId)
		}
	})

	t.Run("Should get solutions", func(t *testing.T) {
		resp, err := client.GetSolutions(ctx, &emptypb.Empty{})
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Solutions) != 1 {
			t.Errorf("expected 1 solution, got %d", len(resp.Solutions))
		}
		if resp.Solutions[0].CorrectOptionId != 2 {
			t.Errorf("expected correct option ID 2, got %d", resp.Solutions[0].CorrectOptionId)
		}
		if resp.Solutions[0].CorrectOptionText != "4" {
			t.Errorf("expected correct option '4', got '%s'", resp.Solutions[0].CorrectOptionText)
		}
	})
}
