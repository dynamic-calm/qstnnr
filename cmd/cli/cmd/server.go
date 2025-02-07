package cmd

import (
	"os"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"
)

func (c *CLI) newServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Manage the qstnnr server",
		Long:  `Commands to start, stop, restart and check status of the qstnnr server`,
	}

	cmd.AddCommand(c.newServerStartCommand())
	cmd.AddCommand(c.newServerStopCommand())
	cmd.AddCommand(c.newServerStatusCommand())
	cmd.AddCommand(c.newServerRestartCommand())

	return cmd
}

func isServerRunning() bool {
	pidFile := filepath.Join(os.TempDir(), "qstnnr-server.pid")
	pidBytes, err := os.ReadFile(pidFile)
	if err != nil {
		return false
	}

	pid, err := strconv.Atoi(string(pidBytes))
	if err != nil {
		return false
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// Send signal 0 to check if process exists
	err = process.Signal(syscall.Signal(0))
	return err == nil
}
