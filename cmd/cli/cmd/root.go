// cmd/cli/cmd/root.go
package cmd

import (
	"fmt"
	"os"

	"github.com/mateopresacastro/qstnnr/api"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CLI struct {
	conn    *grpc.ClientConn
	client  api.QuestionnaireClient
	rootCmd *cobra.Command
}

var cli *CLI

func init() {
	cli = &CLI{
		rootCmd: &cobra.Command{
			Use:   "qstnnr",
			Short: "A simple quiz CLI",
			Long:  `A CLI application for taking quizzes and viewing results.`,
			PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
				return cli.connect()
			},
			PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
				return cli.close()
			},
		},
	}

	cli.addCommands()
}

func (c *CLI) connect() error {
	conn, err := grpc.NewClient("localhost:5974", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	c.conn = conn
	c.client = api.NewQuestionnaireClient(conn)
	return nil
}

func (c *CLI) close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func Execute() {
	if err := cli.rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func (c *CLI) addCommands() {
	c.rootCmd.AddCommand(c.newTakeCommand())
}
