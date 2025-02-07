package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func (c *CLI) newServerStartCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the qstnnr server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if isServerRunning() {
				return fmt.Errorf("server is already running")
			}

			// Get the executable directory
			ex, err := os.Executable()
			if err != nil {
				return fmt.Errorf("failed to get executable path: %v", err)
			}
			binDir := filepath.Dir(ex)
			serverBin := filepath.Join(binDir, "server")

			// Start the server as a background process
			serverCmd := exec.Command(serverBin)
			serverCmd.Stdout = os.Stdout
			serverCmd.Stderr = os.Stderr

			if err := serverCmd.Start(); err != nil {
				return fmt.Errorf("failed to start server: %v", err)
			}

			// Write PID to file
			pidFile := filepath.Join(os.TempDir(), "qstnnr-server.pid")
			if err := os.WriteFile(pidFile, []byte(fmt.Sprint(serverCmd.Process.Pid)), 0644); err != nil {
				return fmt.Errorf("failed to write PID file: %v", err)
			}

			fmt.Println("Server started successfully")
			return nil
		},
	}
}
