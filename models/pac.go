package models

import (
	"bufio"
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

// A PAC structure represents the full contents of a .pac archive.
type PAC struct {
	Header  PACHeader  `yaml:"header"`
	Entries []PACEntry `yaml:"entries"`
}

// Reads a full .pac archive from a stream.
func ReadPAC(r io.ReadSeeker) (out PAC, err error) {
	out.Header, err = ReadPACHeader(r)
	if err != nil {
		return out, errors.Wrap(err, "Header")
	}

	out.Entries = make([]PACEntry, out.Header.Count)
	for i := uint32(0); i < out.Header.Count; i++ {
		if _, err := r.Seek(int64(out.Header.IndexPtr+(i*out.Header.IndexStride)), io.SeekStart); err != nil {
			return out, err
		}
		entry, err := ReadPACEntry(r)
		if err != nil {
			return out, errors.Wrapf(err, "Entry %d", i)
		}
		out.Entries[i] = entry
	}

	return out, err
}

type PACHeaderMagic [4]byte

func (m PACHeaderMagic) String() string {
	return string(m[:3])
}

func (m PACHeaderMagic) MarshalYAML() (interface{}, error) {
	return m.String(), nil
}

func (m *PACHeaderMagic) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s []byte
	if err := unmarshal(&s); err != nil {
		return err
	}
	if len(s) != 3 {
		return errors.Errorf("magic number must be a 3-byte string: %s", s)
	}
	m[0] = s[0]
	m[1] = s[1]
	m[2] = s[2]
	m[3] = 0
	return nil
}

// A PACHeader represents the file header of a PAC.
type PACHeader struct {
	// Magic number. Always "add\0".
	Magic PACHeaderMagic `yaml:"magic" json:"-"`
	// Archive version. Always 4 for Nanoha GoD, may differ for other games on the same engine.
	Version uint32 `yaml:"version" json:"-"`

	// Pointer to the first PACEntry header. Always 32.
	IndexPtr uint32 `yaml:"index_ptr" json:"index_ptr"`
	// Length of a single PACEntry header. Always 32. Unused; the game has this hardcoded.
	IndexStride uint32 `yaml:"index_stride" json:"index_stride"`
	// Number of entries in the archive.
	Count uint32 `yaml:"count" json:"-"`

	// Total file size of the archive. Not actually used, the game checks the file's length.
	FileLen uint32 `yaml:"file_len" json:"-"`
}

// Reads a PACHeader from a stream. It must be positioned at the start of the file.
func ReadPACHeader(r io.ReadSeeker) (out PACHeader, err error) {
	if err := binary.Read(r, binary.LittleEndian, &out); err != nil {
		return out, errors.Wrap(err, "Version")
	}

	if out.Magic[0] != 'a' || out.Magic[1] != 'd' || out.Magic[2] != 'd' || out.Magic[3] != 0 {
		return out, errors.New("Magic: Mismatch; file is not a valid .pac")
	}
	if out.Version != 4 {
		return out, errors.Errorf("Version: Mismatch; expected 4, got %d", out.Version)
	}

	return out, nil
}

type PACEntryFields struct {
	DataPtr     uint32 `yaml:"data_ptr" json:"-"`      // Pointer to the actual data.
	DataLen     uint32 `yaml:"data_len" json:"-"`      // Length of the payload.
	Unused      uint32 `yaml:"unused" json:"-"`        // Seemingly unused, always 0.
	Unknown     uint32 `yaml:"unknown" json:"unknown"` // Unknown purpose.
	FilenamePtr uint32 `yaml:"filename_ptr" json:"-"`  // Pointer to the filename.
}

// A PACEntry represents an entry in a PAC.
type PACEntry struct {
	PACEntryFields `yaml:",inline"`
	Ptr            struct {
		Filename string `yaml:"filename"` // filename_ptr -> null terminator
		Data     []byte `yaml:"-"`        // data_ptr -> data_len
	} `yaml:"ptr" json:"-"`
}

// Reads a PACEntry from a stream. It must already be positioned at the start of the entry (see
// PACHeader's IndexPtr and IndexStride), and will leave the cursor in an undefined position.
func ReadPACEntry(r io.ReadSeeker) (out PACEntry, err error) {
	if err := binary.Read(r, binary.LittleEndian, &out.PACEntryFields); err != nil {
		return out, errors.Wrap(err, "Fields")
	}

	// Follow the filename pointer and read the actual filename.
	if _, err := r.Seek(int64(out.FilenamePtr), io.SeekStart); err != nil {
		return out, errors.Wrap(err, "Ptr.Filename")
	}
	filename, err := bufio.NewReader(r).ReadString(0)
	if err != nil {
		return out, errors.Wrap(err, "Ptr.Filename")
	}
	out.Ptr.Filename = filename[:len(filename)-1]

	// Follow the data pointer and read the payload.
	if _, err := r.Seek(int64(out.DataPtr), io.SeekStart); err != nil {
		return out, errors.Wrap(err, "Ptr.Data")
	}
	out.Ptr.Data = make([]byte, out.DataLen)
	if _, err := io.ReadFull(r, out.Ptr.Data); err != nil {
		return out, errors.Wrap(err, "Ptr.Data")
	}

	return out, nil
}
