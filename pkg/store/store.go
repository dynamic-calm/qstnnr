package store

import (
	"errors"
	"fmt"
	"sync"
)

// Store defines the interface for persistent storage operations
// of questions, solutions, and scores.
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

// QuestionID uniquely identifies a question in the store.
type QuestionID int

// OptionID uniquely identifies an answer option within a question.
type OptionID int

// Score represents the number of correct answers in a submission.
type Score = int

// Stat represents a percentile score comparing against other submissions.
type Stat = int

// Question represents a multiple choice question with its available options.
type Question struct {
	ID      QuestionID
	Text    string
	Options map[OptionID]Option
}

// Option represents a single answer choice for a question.
type Option struct {
	ID   OptionID
	Text string
}

type answer struct {
	QuestionID QuestionID
	OptionID   OptionID
}

// Solution combines a question with its correct answer.
type Solution struct {
	Question        Question
	CorrectOptionID OptionID
}

type submitResult struct {
	Solutions []*Solution
	Stats     int
}

// InitialData contains the required data to initialize a new store.
type InitialData struct {
	Questions map[QuestionID]Question
	Solutions map[QuestionID]OptionID
}

// StoreError indicates an expected error condition in the store operations,
// as opposed to unexpected errors that would indicate bugs.
type StoreError struct {
	error
}

// NewMemoryStore initiates an implementation of the Store interface
// with the given data.
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

// GetQuestions returns all available questions from the store.
func (s *memoryStore) GetQuestions() (map[QuestionID]Question, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.questions, nil
}

// GetSolutions returns the correct answers for all questions from the store.
func (s *memoryStore) GetSolutions() (map[QuestionID]OptionID, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.solutions, nil
}

// SaveScore stores a new score in the store. Returns an error if the score is negative.
func (s *memoryStore) SaveScore(score Score) error {
	if score < 0 {
		return StoreError{fmt.Errorf("score cannot be negative: %d", score)}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.scores = append(s.scores, score)
	return nil
}

// GetAllScores returns a copy of all stored scores.
func (s *memoryStore) GetAllScores() ([]Score, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	scoresCopy := make([]Score, len(s.scores))
	copy(scoresCopy, s.scores)
	return scoresCopy, nil
}
