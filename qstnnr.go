package qstnnr

import (
	"math"
)

// QService defines the questionnaire operations.
type QService interface {
	GetQuestions() (map[QuestionID]Question, error)
	SubmitAnswers(answers map[QuestionID]OptionID) (*SubmitResult, error)
	GetSolutions() (map[QuestionID]OptionID, error)
}

// QstnnrService implements QService using a persistent store.
type QstnnrService struct {
	store Store
}

// SubmitResult contains a map of questions and their correct options,
// and the user's percentile ranking.
type SubmitResult struct {
	Solutions map[QuestionID]OptionID
	Stat      Stat
	Correct   int
}

// SubmitResult contains quiz submission results and ranking.
type ServiceError struct {
	error
}

// NewQstnnrService creates a new questionnaire service.
func NewQstnnrService(store Store) QService {
	return &QstnnrService{store: store}
}

// GetQuestions returns all available questions.
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

// SubmitAnswers processes a questionnaire submission and returns results.
func (qs *QstnnrService) SubmitAnswers(answers map[QuestionID]OptionID) (*SubmitResult, error) {
	if len(answers) == 0 {
		return nil, ServiceError{wrapError(nil, ErrorCodeInvalidInput, "no answers provided")}
	}

	qsts, err := qs.store.GetQuestions()
	if err != nil {
		if _, ok := err.(StoreError); !ok {
			return nil, err
		}
		return nil, ServiceError{wrapError(err, ErrorCodeInternal, "failed to get questions")}
	}

	if len(answers) != len(qsts) {
		msg := "number of answers (%d) must match number of questions (%d)"
		return nil, ServiceError{wrapError(nil, ErrorCodeInvalidInput, msg, len(answers), len(qsts))}
	}

	for qID := range answers {
		if _, ok := qsts[qID]; !ok {
			return nil, ServiceError{wrapError(nil, ErrorCodeInvalidInput, "couldn't find question with id: %d", qID)}
		}
	}

	solutions, err := qs.store.GetSolutions()
	if err != nil {
		if _, ok := err.(StoreError); !ok {
			return nil, err
		}
		return nil, StoreError{wrapError(err, ErrorCodeInternal, "failed to get solutions")}
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
		return nil, ServiceError{wrapError(err, ErrorCodeInternal, "calculating stats")}
	}

	if err := qs.store.SaveScore(correct); err != nil {
		if _, ok := err.(StoreError); !ok {
			return nil, err
		}
		return nil, ServiceError{wrapError(err, ErrorCodeInternal, "failed to save score: %d", correct)}
	}

	return &SubmitResult{Solutions: solutions, Stat: stat, Correct: correct}, nil
}

// stats calculates the percentile ranking for a score.
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

// GetSolutions returns the correct answers for all questions.
func (qs *QstnnrService) GetSolutions() (map[QuestionID]OptionID, error) {
	solutions, err := qs.store.GetSolutions()
	if err != nil {
		if _, ok := err.(StoreError); !ok {
			return nil, err
		}
		return nil, ServiceError{wrapError(err, ErrorCodeInternal, "failed to get solutions")}
	}
	return solutions, nil
}
