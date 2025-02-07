package server

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
)

func NewServerStopCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop the qstnnr server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !isServerRunning() {
				return fmt.Errorf("server is not running")
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

			process, err := os.FindProcess(pid)
			if err != nil {
				return fmt.Errorf("failed to find process: %v", err)
			}

			if err := process.Signal(os.Interrupt); err != nil {
				return fmt.Errorf("failed to stop server: %v", err)
			}

			if err := os.Remove(pidFile); err != nil {
				return fmt.Errorf("failed to remove PID file: %v", err)
			}

			fmt.Println("Server stopped successfully")
			return nil
		},
	}
}
