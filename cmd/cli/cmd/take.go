// cmd/cli/cmd/take.go
package cmd

import (
	"context"
	"fmt"

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

	answers := make(map[int32]int32)

	for _, q := range questions.Questions {
		fmt.Printf("\nQuestion: %s\n", q.Text)
		for _, opt := range q.Options {
			fmt.Printf("%d. %s\n", opt.Id, opt.Text)
		}

		var answer int32
		fmt.Print("\nYour answer (enter the number): ")
		fmt.Scanf("%d", &answer)
		answers[q.Id] = answer
	}

	return nil
}
