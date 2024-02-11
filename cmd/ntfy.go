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
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rwxd/notify-me/services/ntfy"
	"github.com/sagikazarmark/slog-shim"
	"github.com/spf13/cobra"
)

var ntfyCmd = &cobra.Command{
	Use:   "ntfy",
	Short: "Send a push notification to a ntfy instance",
	PreRun: func(cmd *cobra.Command, args []string) {
		setupLogging()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := ensureNtfyConfigCorrect(cmd); err != nil {
			fmt.Println(err)
			cmd.Help()
			os.Exit(1)
		} else if err := ensureNtfyCmdConfigCorrect(cmd); err != nil {
			fmt.Println(err)
			cmd.Help()
			os.Exit(1)
		}

		instance, _ := cmd.Flags().GetString("instance")
		tags, _ := cmd.Flags().GetStringSlice("tags")

		if !strings.Contains(instance, "http://") || !strings.Contains(instance, "https://") {
			instance = "https://" + instance
		}

		notification := ntfy.NewNotification(cmd.Flags().Lookup("topic").Value.String(),
			cmd.Flags().Lookup("title").Value.String(),
			cmd.Flags().Lookup("message").Value.String(),
			ntfy.Priority(cmd.Flags().Lookup("priority").Value.String()),
			tags,
			cmd.Flags().Lookup("url").Value.String(),
			cmd.Flags().Lookup("actions").Value.String(),
			cmd.Flags().Lookup("delay").Value.String(),
			cmd.Flags().Lookup("icon").Value.String(),
			cmd.Flags().Lookup("markdown").Changed,
		)
		if err := ntfy.SendNotification(
			notification,
			instance,
			cmd.Flags().Lookup("user").Value.String(),
			cmd.Flags().Lookup("pass").Value.String(),
		); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Notification sent")
	},
}

var ntfyWrapCmd = &cobra.Command{
	Use:   "wrap",
	Short: "Wrap a command and send a push notification to a ntfy instance",
	PreRun: func(cmd *cobra.Command, args []string) {
		setupLogging()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := ensureNtfyConfigCorrect(cmd); err != nil {
			fmt.Println(err)
			cmd.Help()
			os.Exit(1)
		} else if err := ensureNtfyWrapCmdConfigCorrect(cmd, args); err != nil {
			fmt.Println(err)
			cmd.Help()
			os.Exit(1)
		}

		instance, _ := cmd.Flags().GetString("instance")
		tags, _ := cmd.Flags().GetStringSlice("tags")
		onlyFailure, _ := cmd.Flags().GetBool("error")
		onlySuccess, _ := cmd.Flags().GetBool("success")

		if !strings.Contains(instance, "http://") || !strings.Contains(instance, "https://") {
			instance = "https://" + instance
		}

		notification := ntfy.NewNotification(cmd.Flags().Lookup("topic").Value.String(),
			cmd.Flags().Lookup("title").Value.String(),
			cmd.Flags().Lookup("message").Value.String(),
			ntfy.Priority(cmd.Flags().Lookup("priority").Value.String()),
			tags,
			cmd.Flags().Lookup("url").Value.String(),
			cmd.Flags().Lookup("actions").Value.String(),
			cmd.Flags().Lookup("delay").Value.String(),
			cmd.Flags().Lookup("icon").Value.String(),
			cmd.Flags().Lookup("markdown").Changed,
		)

		program := args[0]
		programArgs := args[1:]
		command := exec.Command(program, programArgs...)

		var output bytes.Buffer
		command.Stdout = &output
		command.Stderr = &output

		slog.Debug("Running command", "program", program, "args", programArgs)
		err := command.Run()

		if notification.Message != "" {
			notification.Message += "\n" + output.String()
		} else {
			notification.Message = output.String()
		}

		if err != nil {
			slog.Debug("Command failed", "error", err)

			if notification.Message == "" {
				notification.Message = err.Error()
			} else {
				notification.Message += "\n" + err.Error()
			}

			if onlySuccess {
				slog.Debug("Only sending on success, not sending notification")
				return
			}

			if err := ntfy.SendNotification(
				notification,
				instance,
				cmd.Flags().Lookup("user").Value.String(),
				cmd.Flags().Lookup("pass").Value.String(),
			); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Notification sent")
		} else {
			slog.Debug("Command succeeded")

			if onlyFailure {
				slog.Debug("Only sending on failure, not sending notification")
				return
			}

			if err := ntfy.SendNotification(
				notification,
				instance,
				cmd.Flags().Lookup("user").Value.String(),
				cmd.Flags().Lookup("pass").Value.String(),
			); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Notification sent")
		}

	},
}

func ensureNtfyConfigCorrect(cmd *cobra.Command) error {
	if cmd.Flags().Changed("user") && !cmd.Flags().Changed("pass") {
		return errors.New("password must be provided if username is provided")
	} else if !cmd.Flags().Changed("user") && cmd.Flags().Changed("pass") {
		return errors.New("username must be provided if password is provided")
	}

	if !cmd.Flags().Changed("topic") {
		return errors.New("topic must be provided")
	}

	return nil
}

func ensureNtfyCmdConfigCorrect(cmd *cobra.Command) error {
	if !cmd.Flags().Changed("message") {
		return errors.New("message must be provided")
	}

	return nil
}

func ensureNtfyWrapCmdConfigCorrect(cmd *cobra.Command, args []string) error {
	if cmd.Flags().Changed("erorr") && cmd.Flags().Changed("success") {
		return errors.New("only one of error or success can be provided")
	}

	if len(args) == 0 {
		return errors.New("command must be provided")
	}

	return nil
}

func init() {
	rootCmd.AddCommand(ntfyCmd)
	ntfyCmd.AddCommand(ntfyWrapCmd)

	ntfyCmd.PersistentFlags().StringP("instance", "i", "ntfy.sh", "ntfy instance")
	ntfyCmd.PersistentFlags().StringP("user", "u", "", "Username for the ntfy instance")
	ntfyCmd.PersistentFlags().StringP("pass", "p", "", "Password for the ntfy instance")
	ntfyCmd.PersistentFlags().StringP("topic", "t", "", "Topic to send the message to")
	ntfyCmd.PersistentFlags().StringP("message", "m", "", "Message")
	ntfyCmd.PersistentFlags().StringP("priority", "P", "", "Message Priority (min, low, default, high, max")
	ntfyCmd.PersistentFlags().StringSlice("tags", []string{}, "Tags for the message")
	ntfyCmd.PersistentFlags().StringP("title", "T", "", "Message title")
	ntfyCmd.PersistentFlags().StringP("url", "U", "", "URL to open when the notification is clicked")
	ntfyCmd.PersistentFlags().String("actions", "", "JSON array or short format of user actions")
	ntfyCmd.PersistentFlags().String("delay", "", "Timestamp or duration for delayed delivery")
	ntfyCmd.PersistentFlags().String("icon", "", "URL to use as notification icon")
	ntfyCmd.PersistentFlags().Bool("markdown", false, "Enable Markdown formatting in the notification body")

	ntfyWrapCmd.Flags().Bool("fail", false, "Send a notification only if the command fails")
	ntfyWrapCmd.Flags().Bool("success", false, "Send a notification only if the command succeeds")
}
