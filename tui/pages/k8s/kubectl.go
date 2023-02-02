package k8s

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/FrangipaneTeam/bean/tui/pages/errorpanel"
)

// Kubectl runs a kubectl command.
func Kubectl(ctx context.Context, k8sCmd *Cmd) tea.Cmd {
	return func() tea.Msg {
		k8sCmd.Done = false
		if k8sCmd.Verb == "" || len(k8sCmd.Files) == 0 {
			return errorpanel.ErrorMsg{
				Reason: "no verb or files provided",
				Cause: fmt.Errorf(
					"verb : %s - files : %d", k8sCmd.Verb, len(k8sCmd.Files),
				),
			}
		}

		var args []string
		switch k8sCmd.Verb {
		case "managed":
			args = []string{"get", "managed"}
		case "apply":
			args = []string{k8sCmd.Verb, "-f", k8sCmd.JoinedFiles()}
		case "delete":
			args = []string{k8sCmd.Verb, "--wait=false", "-f", k8sCmd.JoinedFiles()}
		}

		cmdChan := make(chan interface{})
		defer close(cmdChan)

		go func() {
			var cmd *exec.Cmd
			if k8sCmd.Debug {
				cmd = exec.CommandContext(ctx, "sleep", "10")
			} else {
				cmd = exec.CommandContext(ctx, "kubectl", args...)
			}
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err := cmd.Run()
			if err != nil {
				if ctx.Err() == nil {
					cmdChan <- errors.New(stderr.String())
				}
				return
			}
			if ctx.Err() == nil {
				cmdChan <- stdout.String()
			}
		}()

		select {
		case cmdResult := <-cmdChan:
			switch result := cmdResult.(type) {
			case error:
				if result != nil {
					return errorpanel.ErrorMsg{
						Reason:   fmt.Sprintf("command kubectl %s failed", strings.Join(args, " ")),
						Cause:    result,
						CmdID:    k8sCmd.ID,
						FromPage: k8sCmd.FromPage,
					}
				}
			case string:
				k8sCmd.Result = result
			}

		case <-ctx.Done():
			return errorpanel.ErrorMsg{
				Reason:   "context done",
				Cause:    errors.New("cancel kubectl command"),
				FromPage: k8sCmd.FromPage,
			}
		}

		return k8sCmd
	}
}
