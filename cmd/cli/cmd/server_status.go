package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
)

func (c *CLI) newServerStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Check the status of the qstnnr server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !isServerRunning() {
				fmt.Println("Server is not running")
				return nil
			}

			pidFile := filepath.Join(os.TempDir(), "qstnnr-server.pid")
			pidBytes, err := os.ReadFile(pidFile)
			if err != nil {
				return fmt.Errorf("failed to read PID file: %v", err)
			}

			pid, err := strconv.Atoi(string(pidBytes))
			if err != nil {
				return fmt.Errorf("invalid PID in file: %v", err)
			}

			fmt.Printf("Server is running (PID: %d)\n", pid)
			return nil
		},
	}
}
