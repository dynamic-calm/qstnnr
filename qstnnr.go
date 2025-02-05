package qstnnr

import (
	"errors"
	"fmt"
	"math"
)

type QService interface {
	GetQuestions() (map[QuestionID]Question, error)
	SubmitAnswers(answers map[QuestionID]OptionID) (*SubmitResult, error)
	GetSolutions() (map[QuestionID]OptionID, error)
}

type QstnnrService struct {
	store Store
}

type SubmitResult struct {
	Solutions map[QuestionID]OptionID
	Stat      Stat
}
type ServiceError struct {
	error
}

func NewQstnnrService(store Store) QService {
	return &QstnnrService{store: store}
}

func (qs *QstnnrService) GetQuestions() (map[QuestionID]Question, error) {
	questions, err := qs.store.GetQuestions()
	if err != nil {
		if _, ok := err.(StoreError); !ok {
			// If this error is not a StoreError we know it's a bug and not a known edge case.
			return nil, err
		}
		return nil, ServiceError{err}
	}
	return questions, nil
}

func (qs *QstnnrService) SubmitAnswers(answers map[QuestionID]OptionID) (*SubmitResult, error) {
	if len(answers) == 0 {
		return nil, ServiceError{errors.New("no answers provided")}
	}

	qsts, err := qs.store.GetQuestions()
	if err != nil {
		if _, ok := err.(StoreError); !ok {
			return nil, err
		}
		return nil, ServiceError{err}
	}

	if len(answers) != len(qsts) {
		return nil, ServiceError{errors.New("number of answers must match number of questions")}
	}

	for qID := range answers {
		if _, ok := qsts[qID]; !ok {
			return nil, ServiceError{fmt.Errorf("couldn't find question with id: %d", qID)}
		}
	}

	solutions, err := qs.store.GetSolutions()
	if err != nil {
		if _, ok := err.(StoreError); !ok {
			return nil, err
		}
		return nil, StoreError{err}
	}

	correct := 0
	for k, v := range solutions {
		if answers[k] == v {
			correct++
		}
	}

	stat, err := qs.stats(correct)
	if err != nil {
		if _, ok := err.(StoreError); !ok {
			return nil, err
		}
		return nil, ServiceError{fmt.Errorf("calculating stats: %w", err)}
	}

	if err := qs.store.SaveScore(correct); err != nil {
		if _, ok := err.(StoreError); !ok {
			return nil, err
		}
		return nil, ServiceError{err}
	}

	return &SubmitResult{Solutions: solutions, Stat: stat}, nil
}

func (qs *QstnnrService) stats(score Score) (Stat, error) {
	scores, err := qs.store.GetAllScores()
	if err != nil {
		return 0, err
	}

	if len(scores) == 0 {
		return 100, nil // First quiz taker.
	}

	betterThan := 0
	for _, s := range scores {
		if score > s {
			betterThan++
		}
	}

	percentage := float64(betterThan) / float64(len(scores)) * 100
	return Stat(math.Round(percentage)), nil
}

func (qs *QstnnrService) GetSolutions() (map[QuestionID]OptionID, error) {
	solutions, err := qs.store.GetSolutions()
	if err != nil {
		if _, ok := err.(StoreError); !ok {
			return nil, err
		}
		return nil, ServiceError{err}
	}
	return solutions, nil
}
