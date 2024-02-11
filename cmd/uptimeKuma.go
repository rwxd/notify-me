/*
Copyright © 2024 rwxd

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	uptimekuma "github.com/rwxd/notify-me/services/uptimeKuma"
	"github.com/sagikazarmark/slog-shim"
	"github.com/spf13/cobra"
)

var uptimeKumaCmd = &cobra.Command{
	Use:   "uptime-kuma",
	Short: "Send a push notification to an uptime-kuma instance",
	PreRun: func(cmd *cobra.Command, args []string) {
		setupLogging()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := ensureConfigCorrect(cmd); err != nil {
			fmt.Println(err)
			cmd.Help()
			os.Exit(1)
		} else if err := ensureDefaultCmdConfigCorrect(cmd); err != nil {
			fmt.Println(err)
			cmd.Help()
			os.Exit(1)
		}

		instance, _ := cmd.Flags().GetString("instance")
		token, _ := cmd.Flags().GetString("token")
		message, _ := cmd.Flags().GetString("message")
		ping, _ := cmd.Flags().GetString("ping")
		down, _ := cmd.Flags().GetBool("down")
		up, _ := cmd.Flags().GetBool("up")

		if !strings.Contains(instance, "http://") || !strings.Contains(instance, "https://") {
			instance = "https://" + instance
		}

		if err := uptimekuma.SendMonitorStatus(instance, token, up || !down, message, ping); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Sent status to uptime-kuma")

	},
}

var uptimeKumaWrapCmd = &cobra.Command{
	Use:   "wrap",
	Short: "Wrap a command and send a push notification to an uptime-kuma instance",
	PreRun: func(cmd *cobra.Command, args []string) {
		setupLogging()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := ensureConfigCorrect(cmd); err != nil {
			fmt.Println(err)
			cmd.Help()
			os.Exit(1)
		} else if err := ensureWrapCmdConfigCorrect(cmd, args); err != nil {
			fmt.Println(err)
			cmd.Help()
			os.Exit(1)
		}

		instance, _ := cmd.Flags().GetString("instance")
		token, _ := cmd.Flags().GetString("token")
		message, _ := cmd.Flags().GetString("message")
		onlyError, _ := cmd.Flags().GetBool("error")
		onlySuccess, _ := cmd.Flags().GetBool("success")
		reverse, _ := cmd.Flags().GetBool("reverse")

		if !strings.Contains(instance, "http://") || !strings.Contains(instance, "https://") {
			instance = "https://" + instance
		}

		program := args[0]
		programArgs := args[1:]

		command := exec.Command(program, programArgs...)

		var output bytes.Buffer
		command.Stdout = &output
		command.Stderr = &output

		slog.Debug("Running command", "program", program, "args", programArgs)
		err := command.Run()

		if message == "" {
			message = output.String()
		}

		statusUp := true
		if err != nil {
			slog.Debug("Command failed", "error", err)

			if message == "" {
				message = err.Error()
			} else {
				message = message + "\n" + err.Error()
			}

			if !reverse {
				statusUp = false
			} else {
				slog.Debug("Reverse is set, setting status to up")
			}

			if onlySuccess {
				slog.Debug("Only success is set, not sending status")
				return
			}

			if err := uptimekuma.SendMonitorStatus(instance, token, statusUp, message, ""); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println("Sent status to uptime-kuma")
		} else {
			if reverse {
				slog.Debug("Reverse is set, setting status to down")
				statusUp = false
			}

			if onlyError {
				slog.Debug("Only error is set, not sending status")
				return
			}

			if err := uptimekuma.SendMonitorStatus(instance, token, statusUp, message, ""); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println("Sent status to uptime-kuma")
		}
	},
}

func ensureConfigCorrect(cmd *cobra.Command) error {
	if !cmd.Flag("token").Changed {
		return fmt.Errorf("You must set the token")
	}

	if !cmd.Flag("instance").Changed {
		return fmt.Errorf("You must set the instance")
	}

	return nil
}

func ensureDefaultCmdConfigCorrect(cmd *cobra.Command) error {
	if cmd.Flag("down").Changed && cmd.Flag("up").Changed {
		return fmt.Errorf("You can't set both down and up")
	} else if !cmd.Flag("down").Changed && !cmd.Flag("up").Changed {
		return fmt.Errorf("You must set either down or up")
	}

	return nil
}

func ensureWrapCmdConfigCorrect(cmd *cobra.Command, args []string) error {
	if cmd.Flag("error").Changed && cmd.Flag("success").Changed {
		return fmt.Errorf("You can't set both error and success")
	}

	if len(args) == 0 {
		return fmt.Errorf("You must provide a command to wrap")
	}

	return nil
}

func init() {
	rootCmd.AddCommand(uptimeKumaCmd)
	uptimeKumaCmd.PersistentFlags().StringP("instance", "i", "", "The instance to send the notification to")
	uptimeKumaCmd.PersistentFlags().StringP("token", "t", "", "Token for the push monitor")
	uptimeKumaCmd.Flags().StringP("message", "m", "", "Message to send to the monitor")
	uptimeKumaCmd.Flags().StringP("ping", "p", "", "Measurement number to send to the monitor")
	uptimeKumaCmd.Flags().Bool("down", false, "Set the monitor to down")
	uptimeKumaCmd.Flags().Bool("up", false, "Set the monitor to up")

	uptimeKumaCmd.AddCommand(uptimeKumaWrapCmd)
	uptimeKumaWrapCmd.Flags().Bool("error", false, "Send a notification only if the command fails")
	uptimeKumaWrapCmd.Flags().Bool("success", false, "Send a notification only if the command succeeds")
	uptimeKumaWrapCmd.Flags().Bool("reverse", false, "Send a up notification if the command fails and a down notification if the command succeeds")
	uptimeKumaWrapCmd.Flags().StringP("message", "m", "", "Custom message, defaults to the command output")
}
