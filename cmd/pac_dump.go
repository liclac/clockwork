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
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/liclac/clockwork/models"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// pacDumpCmd represents the pac dump command
var pacDumpCmd = &cobra.Command{
	Use:   "dump file",
	Short: "Dump the contents of a .pac archive to a directory.",
	Long: `Dump the contents of a .pac archive to a directory, specified by --out/-o.

By default, the output will be placed in a directory with the same name as the
archive, in the current dir; eg. /path/to/USRDIR/global.pac -> ./global.pac/*.

Each archive entry will yield two files: 00X_filename and 00X_filename.json,
the latter being a metadata file for information that cannot be encoded in the
file itself. 00X denotes the sequence number, starting at 000.

Use 'pac repack' to reassemble a directory dump into a .pac archive.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := args[0]

		out := cmd.Flag("out").Value.String()
		if out == "" {
			out = filepath.Base(filename)
		}

		fs := afero.NewOsFs()
		file, err := fs.Open(filename)
		if err != nil {
			return err
		}
		defer func() { _ = file.Close() }()

		p, err := models.ReadPAC(file)
		if err != nil {
			return err
		}

		fmt.Printf("- %s\n", out)
		if err := fs.MkdirAll(out, 0755); err != nil {
			return err
		}

		metapath := filepath.Join(out, models.PACHeaderFilename)
		metadata, err := json.MarshalIndent(p.Header, "", "  ")
		if err != nil {
			return err
		}
		if err := afero.WriteFile(fs, metapath, metadata, 0644); err != nil {
			return err
		}

		for i, entry := range p.Entries {
			datapath := filepath.Join(out, models.ToPACEntryFilename(uint32(i), entry.Ptr.Filename))
			metapath := datapath + models.PACEntryMetaSuffix

			fmt.Printf("- %s\n", metapath)
			metadata, err := json.MarshalIndent(entry, "", "  ")
			if err != nil {
				return err
			}
			if err := afero.WriteFile(fs, metapath, metadata, 0644); err != nil {
				return err
			}

			fmt.Printf("- %s\n", datapath)
			if err := afero.WriteFile(fs, datapath, entry.Ptr.Data, 0644); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	pacCmd.AddCommand(pacDumpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pacDumpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pacDumpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	pacDumpCmd.Flags().StringP("out", "o", "", "Output directory; if not set, creates a directory from the archive's base name (dropping the .pac suffix).")
}
