package server

import (
	"time"

	"github.com/spf13/cobra"
)

func NewServerRestartCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "restart",
		Short: "Restart the qstnnr server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := NewServerStopCommand().RunE(cmd, args); err != nil {

				return err
			}

			// Small delay to ensure proper shutdown
			time.Sleep(time.Second)

			return NewServerStartCommand().RunE(cmd, args)
		},
	}
}
