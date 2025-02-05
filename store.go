package qstnnr

import (
	"errors"
	"fmt"
	"sync"
)

type Store interface {
	GetQuestions() (map[QuestionID]Question, error)
	GetSolutions() (map[QuestionID]OptionID, error)
	SaveScore(score Score) error
	GetAllScores() ([]Score, error)
}

type memoryStore struct {
	questions map[QuestionID]Question
	solutions map[QuestionID]OptionID
	scores    []Score
	mu        sync.RWMutex
}

type QuestionID int
type OptionID int
type Score = int // Number of correct answers.
type Stat = int  // Percentile calculation

type Question struct {
	ID      QuestionID
	Text    string
	Options map[OptionID]Option
}

type Option struct {
	ID   OptionID
	Text string
}

type answer struct {
	QuestionID QuestionID
	OptionID   OptionID
}

type Solution struct {
	Question        Question
	CorrectOptionID OptionID
}

type submitResult struct {
	Solutions []*Solution
	Stats     int
}

type InitialData struct {
	Questions map[QuestionID]Question
	Solutions map[QuestionID]OptionID
}

type StoreError struct {
	error
}

func NewMemoryStore(data InitialData) (Store, error) {
	if data.Questions == nil || data.Solutions == nil {
		return nil, StoreError{errors.New("questions and solutions maps cannot be nil")}
	}
	return &memoryStore{
		questions: data.Questions,
		solutions: data.Solutions,
		scores:    make([]Score, 0),
		mu:        sync.RWMutex{},
	}, nil
}

func (s *memoryStore) GetQuestions() (map[QuestionID]Question, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.questions, nil
}

func (s *memoryStore) GetSolutions() (map[QuestionID]OptionID, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.solutions, nil
}

func (s *memoryStore) SaveScore(score Score) error {
	if score < 0 {
		return StoreError{fmt.Errorf("score cannot be negative: %d", score)}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.scores = append(s.scores, score)
	return nil
}

func (s *memoryStore) GetAllScores() ([]Score, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	scoresCopy := make([]Score, len(s.scores))
	copy(scoresCopy, s.scores)
	return scoresCopy, nil
}
