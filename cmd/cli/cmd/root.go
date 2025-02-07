package cmd

import (
	"fmt"
	"os"

	"github.com/mateopresacastro/qstnnr/api"
	"github.com/mateopresacastro/qstnnr/cmd/cli/cmd/server"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CLI struct {
	conn    *grpc.ClientConn
	client  api.QuestionnaireClient
	rootCmd *cobra.Command
	port    string
}

var cli *CLI

func init() {
	cli = &CLI{
		port: "5974",
		rootCmd: &cobra.Command{
			Use:   "qstnnr",
			Short: "A simple quiz CLI",
			Long:  `A CLI application for taking quizzes and viewing results.`,
			PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
				if cmd.Parent() != nil && cmd.Parent().Name() == "server" {
					return nil
				}
				return cli.connect()
			},
			PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
				if cmd.Parent() != nil && cmd.Parent().Name() == "server" {
					return nil
				}
				return cli.close()
			},
		},
	}
	cli.rootCmd.PersistentFlags().StringVarP(&cli.port, "port", "p", cli.port, "Port to connect to")
	cli.addCommands()
}

func (c *CLI) connect() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = c.port
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("localhost:%s", port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
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
	c.rootCmd.AddCommand(server.NewServerCommand())
}
