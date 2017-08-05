package models

import (
	"fmt"
	"strings"
)

const (
	PACHeaderFilename  = "_metadata.json" // Filename for dumped header data.
	PACEntryMetaSuffix = ".json"          // Suffix for dumped entry metadata.
)

// Returns the dumped filename for the i'th file in an archive, with the given base filename.
func ToPACEntryFilename(i uint32, filename string) string {
	return fmt.Sprintf("%03d_%s", i, filename)
}

// Returns the base filename from the given sequenced filename (from ToPACEntryFilename).
func FromPACEntryFilename(filename string) string {
	return strings.SplitN(filename, "_", 2)[1]
}
