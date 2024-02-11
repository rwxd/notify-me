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
	"fmt"
	"os"
	"strings"

	uptimekuma "github.com/rwxd/notify-me/services/uptimeKuma"
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
			os.Exit(1)
		}

		instance, _ := cmd.Flags().GetString("instance")
		token, _ := cmd.Flags().GetString("token")
		down, _ := cmd.Flags().GetBool("down")
		up, _ := cmd.Flags().GetBool("up")

		if !strings.Contains(instance, "http://") || !strings.Contains(instance, "https://") {
			instance = "https://" + instance
		}

		if err := uptimekuma.SendMonitorStatus(instance, token, up || !down); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Sent status to uptime-kuma")

	},
}

func ensureConfigCorrect(cmd *cobra.Command) error {
	if cmd.Flag("down").Changed && cmd.Flag("up").Changed {
		return fmt.Errorf("You can't set both down and up")
	} else if !cmd.Flag("down").Changed && !cmd.Flag("up").Changed {
		return fmt.Errorf("You must set either down or up")
	}

	if !cmd.Flag("token").Changed {
		return fmt.Errorf("You must set the token")
	}

	if !cmd.Flag("instance").Changed {
		return fmt.Errorf("You must set the instance")
	}

	return nil
}

func init() {
	rootCmd.AddCommand(uptimeKumaCmd)
	uptimeKumaCmd.Flags().StringP("instance", "i", "", "The instance to send the notification to")
	uptimeKumaCmd.Flags().StringP("token", "t", "", "Token for the push monitor")
	uptimeKumaCmd.Flags().Bool("down", false, "Set the monitor to down")
	uptimeKumaCmd.Flags().Bool("up", false, "Set the monitor to up")
}
