package cmd

import (
	"context"
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/mateopresacastro/qstnnr/pkg/api"
	"github.com/mateopresacastro/qstnnr/pkg/store"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *CLI) newTakeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "take",
		Short: "Take the quiz",
		Long:  `Start a new quiz session and answer questions`,
		RunE:  c.runTakeQuiz,
	}
}

func (c *CLI) runTakeQuiz(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	questions, err := c.client.GetQuestions(ctx, &emptypb.Empty{})
	if err != nil {
		return fmt.Errorf("Something went wrong when trying to connect to the server. Did you run `qstnnr server start`?: %w", err)
	}

	answers := make(map[store.QuestionID]store.OptionID)
	for i, q := range questions.Questions {
		fmt.Printf("Question %d of %d\n", i+1, len(questions.Questions))
		// Create options slice for the select prompt
		options := make([]string, len(q.Options))
		for j, opt := range q.Options {
			options[j] = opt.Text
		}

		prompt := promptui.Select{
			Label: q.Text,
			Items: options,
			Templates: &promptui.SelectTemplates{
				Label:    "{{ . }}",
				Selected: fmt.Sprintf(`✔ Question %d: {{ . }}`, i+1),
				Active:   "➜ {{ . | cyan }}",
				Inactive: "  {{ . }}",
			},
		}

		index, _, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("prompt failed: %v", err)
		}

		// Store the answer
		answers[store.QuestionID(q.Id)] = store.OptionID(q.Options[index].Id)
	}

	if len(answers) <= 0 {
		return nil
	}

	confirm := promptui.Prompt{
		Label:     "Submit your answers",
		IsConfirm: true,
	}

	result, err := confirm.Run()
	if err != nil {
		return fmt.Errorf("prompt failed: %v", err)
	}

	if result != "y" && result != "Y" && result != "" {
		return nil
	}

	fmt.Println("\nSubmitting answers...")
	var fmtAnswers []*api.Answer
	for qID, oID := range answers {
		fmtAnswers = append(fmtAnswers, &api.Answer{
			QuestionId: int32(qID),
			OptionId:   int32(oID),
		})
	}

	req := &api.SubmitAnswersRequest{Answers: fmtAnswers}
	submitRes, err := c.client.SubmitAnswers(ctx, req)
	if err != nil {
		return err
	}

	fmt.Printf("\nYou got %d correct!\n", submitRes.Correct)
	fmt.Printf("That's better than %d%% of participants! 🌱\n", submitRes.BetterThan)

	reviewPrompt := promptui.Prompt{
		Label:     "Would you like to check the solutions",
		IsConfirm: true,
	}

	reviewResult, err := reviewPrompt.Run()
	if err != nil {
		return nil
	}

	if reviewResult != "y" && reviewResult != "Y" && reviewResult != "" {
		return nil
	}

	for _, solution := range submitRes.Solutions {
		fmt.Printf("\n%s\n", solution.Question.Text)
		userAnswer := answers[store.QuestionID(solution.Question.Id)]

		if userAnswer == store.OptionID(solution.CorrectOptionId) {
			// Correct
			fmt.Printf("\033[32m✓ %s\033[0m\n", solution.CorrectOptionText)
		} else {
			// Incorrect
			originalQ := findQuestion(questions.Questions, solution.Question.Id)
			userAnswerText := findOptionText(originalQ.Options, int32(userAnswer))

			fmt.Printf("\033[32m✓ Correct: %s\033[0m\n", solution.CorrectOptionText)
			fmt.Printf("\033[31m✗ Your answer: %s\033[0m\n", userAnswerText)
		}
	}

	return nil
}

func findQuestion(questions []*api.Question, id int32) *api.Question {
	for _, q := range questions {
		if q.Id == id {
			return q
		}
	}
	return nil
}

func findOptionText(options []*api.Option, id int32) string {
	for _, opt := range options {
		if opt.Id == id {
			return opt.Text
		}
	}
	return "Unknown option"
}
