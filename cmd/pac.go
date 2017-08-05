// Copyright © 2017 The Clockwork Creators
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"

	"github.com/liclac/clockwork/models"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

// pacCmd represents the pac command
var pacCmd = &cobra.Command{
	Use:   "pac file",
	Short: "Dump the index of a .pac archive.",
	Long: `Dump the index of a .pac archive as YAML to stdout.

Refer to the various subcommands for how to dump and manipulate the contents.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer func() { _ = file.Close() }()

		p, err := models.ReadPAC(file)
		if err != nil {
			return err
		}

		data, err := yaml.Marshal(p)
		if err != nil {
			return err
		}
		fmt.Print(string(data))

		return nil
	},
}

func init() {
	RootCmd.AddCommand(pacCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pacCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pacCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
