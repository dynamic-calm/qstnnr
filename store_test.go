package qstnnr_test

import (
	"testing"

	"github.com/mateopresacastro/qstnnr"
)

func TestStore(t *testing.T) {
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

	t.Run("should fail with nil data", func(t *testing.T) {
		_, err := qstnnr.NewMemoryStore(qstnnr.InitialData{})
		if err == nil {
			t.Fatal("expected error with nil data")
		}
	})

	store, err := qstnnr.NewMemoryStore(qstnnr.InitialData{
		Questions: questions,
		Solutions: solutions,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("should get questions", func(t *testing.T) {
		qs, err := store.GetQuestions()
		if err != nil {
			t.Fatal(err)
		}
		if len(qs) != len(questions) {
			t.Fatalf("expected %d questions, got %d", len(questions), len(qs))
		}
		if qs[1].Text != questions[1].Text {
			t.Fatalf("expected question text %q, got %q", questions[1].Text, qs[1].Text)
		}
	})

	t.Run("should get solutions", func(t *testing.T) {
		sols, err := store.GetSolutions()
		if err != nil {
			t.Fatal(err)
		}
		if len(sols) != len(solutions) {
			t.Fatalf("expected %d solutions, got %d", len(solutions), len(sols))
		}
		if sols[1] != solutions[1] {
			t.Fatalf("expected solution %d, got %d", solutions[1], sols[1])
		}
	})

	t.Run("should save and get scores", func(t *testing.T) {
		// Valid scores
		scores := []qstnnr.Score{2, 3, 1}
		for _, score := range scores {
			if err := store.SaveScore(score); err != nil {
				t.Fatalf("failed to save score %d: %v", score, err)
			}
		}

		// Get scores
		savedScores, err := store.GetAllScores()
		if err != nil {
			t.Fatal(err)
		}
		if len(savedScores) != len(scores) {
			t.Fatalf("expected %d scores, got %d", len(scores), len(savedScores))
		}

		for i, score := range scores {
			if savedScores[i] != score {
				t.Fatalf("expected score %d at position %d, got %d", score, i, savedScores[i])
			}
		}
	})

	t.Run("should not save negative scores", func(t *testing.T) {
		err := store.SaveScore(-1)
		if err == nil {
			t.Fatal("expected error when saving negative score")
		}
	})

	t.Run("error should be of correct type", func(t *testing.T) {
		err := store.SaveScore(-1)
		if _, ok := err.(qstnnr.StoreError); !ok {
			t.Fatal("error is not of correct type")
		}
	})

	t.Run("should return copy of scores", func(t *testing.T) {
		scores1, err := store.GetAllScores()
		if err != nil {
			t.Fatal(err)
		}

		// Modify the returned scores
		scores1[0] = 999

		// Get scores again
		scores2, err := store.GetAllScores()
		if err != nil {
			t.Fatal(err)
		}

		// Verify the modification didn't affect the store
		if scores2[0] == 999 {
			t.Fatal("modification of returned scores affected the store")
		}
	})
}
