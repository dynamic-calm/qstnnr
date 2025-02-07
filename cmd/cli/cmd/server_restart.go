package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

func (c *CLI) newServerRestartCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "restart",
		Short: "Restart the qstnnr server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.newServerStopCommand().RunE(cmd, args); err != nil {
				return err
			}

			// Small delay to ensure proper shutdown
			time.Sleep(time.Second)

			return c.newServerStartCommand().RunE(cmd, args)
		},
	}
}
