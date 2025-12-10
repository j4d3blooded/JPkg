package jpkg

import (
	"fmt"
	"io"

	jpkg_bin "github.com/j4d3blooded/jpkg/bin"
)

func ReadJPkg(r io.ReadSeeker, encKey []byte) (*JPkg, error) {

	header, err := parseHeader(r)
	if err != nil {
		return nil, fmt.Errorf("error reading header: %w", err)
	}

	manifest, err := parseManifest(r)
	if err != nil {
		return nil, fmt.Errorf("error reading manifest: %w", err)
	}

	return &JPkg{Header: *header, Manifest: *manifest}, nil
}

func parseHeader(r io.ReadSeeker) (*JPkgHeader, error) {
	header, err := jpkg_bin.BinaryRead[JPkgHeader](r)

	if err != nil {
		return nil, fmt.Errorf("error reading jpkg header: %w", err)
	}

	if header.MagicNumber != MAGIC_NUMBER {
		return nil, fmt.Errorf("incorrect magic number")
	}

	if header.Version != 0 {
		return nil, fmt.Errorf("unsupported version: %v", header.Version)
	}

	return header, nil
}

func parseManifest(r io.ReadSeeker) (*JPkgManifest, error) {
	header, err := jpkg_bin.BinaryRead[JPkgManifest](r)

	if err != nil {
		return nil, fmt.Errorf("error reading jpkg manifest: %w", err)
	}

	return header, nil
}
