package qstnnr_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/mateopresacastro/qstnnr"
	"github.com/mateopresacastro/qstnnr/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestE2E(t *testing.T) {
	// Setup a buffer for logs
	var buf bytes.Buffer
	port := "4000"

	getenv := func(key string) string {
		if key == "PORT" {
			return port
		}
		return ""
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- qstnnr.Run(ctx, getenv, &buf)
	}()

	time.Sleep(100 * time.Millisecond)

	conn, err := grpc.NewClient("localhost:"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := api.NewQuestionnaireClient(conn)

	t.Run("complete quiz flow", func(t *testing.T) {
		questions, err := client.GetQuestions(ctx, &emptypb.Empty{})
		if err != nil {
			t.Fatalf("failed to get questions: %v", err)
		}
		if len(questions.Questions) != 10 {
			t.Errorf("expected 10 questions, got %d", len(questions.Questions))
		}

		// All correct
		answers := []*api.Answer{
			{QuestionId: 1, OptionId: 2},  // defer()
			{QuestionId: 2, OptionId: 2},  // var s []int
			{QuestionId: 3, OptionId: 1},  // nil
			{QuestionId: 4, OptionId: 1},  // go
			{QuestionId: 5, OptionId: 1},  // panic
			{QuestionId: 6, OptionId: 3},  // both interface{} and any
			{QuestionId: 7, OptionId: 1},  // discard unwanted value
			{QuestionId: 8, OptionId: 1},  // lowercase letter
			{QuestionId: 9, OptionId: 1},  // value, exists := map[key]
			{QuestionId: 10, OptionId: 1}, // type I interface {}
		}

		result, err := client.SubmitAnswers(ctx, &api.SubmitAnswersRequest{
			Answers: answers,
		})
		if err != nil {
			t.Fatalf("failed to submit answers: %v", err)
		}

		// All correct
		if result.Correct != 10 {
			t.Errorf("expected 10 correct answers, got %d", result.Correct)
		}
		if result.BetterThan != 100 {
			t.Errorf("expected to be better than 100%% of users, got %d%%", result.BetterThan)
		}

		// 3. Submit some wrong answers
		wrongAnswers := []*api.Answer{
			{QuestionId: 1, OptionId: 1},  // wrong
			{QuestionId: 2, OptionId: 1},  // wrong
			{QuestionId: 3, OptionId: 1},  // correct
			{QuestionId: 4, OptionId: 2},  // wrong
			{QuestionId: 5, OptionId: 2},  // wrong
			{QuestionId: 6, OptionId: 1},  // wrong
			{QuestionId: 7, OptionId: 1},  // correct
			{QuestionId: 8, OptionId: 2},  // wrong
			{QuestionId: 9, OptionId: 1},  // correct
			{QuestionId: 10, OptionId: 1}, // correct
		}

		result, err = client.SubmitAnswers(ctx, &api.SubmitAnswersRequest{
			Answers: wrongAnswers,
		})
		if err != nil {
			t.Fatalf("failed to submit answers: %v", err)
		}

		if result.Correct != 4 {
			t.Errorf("expected 4 correct answers, got %d", result.Correct)
		}
		if result.BetterThan != 0 {
			t.Errorf("expected to be better than 0%% of users, got %d%%", result.BetterThan)
		}

		solutions, err := client.GetSolutions(ctx, &emptypb.Empty{})
		if err != nil {
			t.Fatalf("failed to get solutions: %v", err)
		}

		if len(solutions.Solutions) != 10 {
			t.Errorf("expected 10 solutions, got %d", len(solutions.Solutions))
		}

		for _, sol := range solutions.Solutions {
			switch sol.Question.Id {
			case 1:
				if sol.CorrectOptionText != "defer()" {
					t.Errorf("wrong solution for question 1, got %s", sol.CorrectOptionText)
				}
			case 3:
				if sol.CorrectOptionText != "nil" {
					t.Errorf("wrong solution for question 3, got %s", sol.CorrectOptionText)
				}
			case 4:
				if sol.CorrectOptionText != "go" {
					t.Errorf("wrong solution for question 4, got %s", sol.CorrectOptionText)
				}
			}
		}
	})

	t.Run("error cases", func(t *testing.T) {
		invalidAnswers := []*api.Answer{
			{QuestionId: 999, OptionId: 1},
		}
		_, err := client.SubmitAnswers(ctx, &api.SubmitAnswersRequest{
			Answers: invalidAnswers,
		})
		if err == nil {
			t.Error("expected error for invalid question ID")
		}

		incompleteAnswers := []*api.Answer{
			{QuestionId: 1, OptionId: 1},
		}
		_, err = client.SubmitAnswers(ctx, &api.SubmitAnswersRequest{
			Answers: incompleteAnswers,
		})
		if err == nil {
			t.Error("expected error for incomplete answers")
		}
	})

	cancel()
	if err := <-errCh; err != nil {
		t.Fatalf("server error: %v", err)
	}
}
