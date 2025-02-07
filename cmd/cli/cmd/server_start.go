package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func (c *CLI) newServerStartCommand() *cobra.Command {
	var verbose bool
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the qstnnr server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if isServerRunning() {
				return fmt.Errorf("server is already running")
			}

			ex, err := os.Executable()
			if err != nil {
				return fmt.Errorf("failed to get executable path: %v", err)
			}
			binDir := filepath.Dir(ex)
			serverBin := filepath.Join(binDir, "server")

			serverCmd := exec.Command(serverBin)
			env := os.Environ()
			if !verbose {
				env = append(env, "LOG_LEVEL=error")
			}
			serverCmd.Env = env

			if verbose {
				serverCmd.Stdout = os.Stdout
				serverCmd.Stderr = os.Stderr
			}

			if err := serverCmd.Start(); err != nil {
				return fmt.Errorf("failed to start server: %v", err)
			}

			pidFile := filepath.Join(os.TempDir(), "qstnnr-server.pid")
			if err := os.WriteFile(pidFile, []byte(fmt.Sprint(serverCmd.Process.Pid)), 0644); err != nil {
				return fmt.Errorf("failed to write PID file: %v", err)
			}

			port := os.Getenv("PORT")
			if port == "" {
				port = c.port
			}
			fmt.Printf("Server started on port %s (PID: %d)\n", port, serverCmd.Process.Pid)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show server logs")
	return cmd
}
