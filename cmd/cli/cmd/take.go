// cmd/cli/cmd/take.go
package cmd

import (
	"context"
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/mateopresacastro/qstnnr"
	"github.com/mateopresacastro/qstnnr/api"
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
		return err
	}

	answers := make(map[qstnnr.QuestionID]qstnnr.OptionID)
	for i, q := range questions.Questions {
		fmt.Printf("Question %d of %d\n", i+1, len(questions.Questions))
		// Create options slice for the select prompt
		options := make([]string, len(q.Options))
		for j, opt := range q.Options {
			options[j] = opt.Text
		}

		// Create the selection prompt
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
		answers[qstnnr.QuestionID(q.Id)] = qstnnr.OptionID(q.Options[index].Id)
	}
	fmt.Println(answers)

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
	fmt.Printf("That's better than %d%% of participants!\n", submitRes.BetterThan)

	return nil
}
