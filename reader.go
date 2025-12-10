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

	files, err := parseFiles(r, manifest.FileCount)
	if err != nil {
		return nil, fmt.Errorf("error reading file records: %w", err)
	}

	buildFS(files)

	return &JPkg{Header: *header, Manifest: *manifest, FileRecords: files}, nil
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

func parseFiles(r io.ReadSeeker, fileCount uint64) ([]JPkgFileRecordWithOffset, error) {
	files := make([]JPkgFileRecordWithOffset, fileCount)

	for i := range fileCount {
		record, err := jpkg_bin.BinaryRead[JPkgFileRecordWithoutData](r)
		if err != nil {
			return nil, fmt.Errorf("error reading file record %v: %w", i, err)
		}

		offset, err := r.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, fmt.Errorf("error seeking in file: %w", err)
		}

		_, err = r.Seek(int64(record.CompressedDataSize), io.SeekCurrent)
		if err != nil {
			return nil, fmt.Errorf("error seeking in file: %w", err)
		}

		files[i] = JPkgFileRecordWithOffset{
			JPkgFileRecordWithoutData: *record,
			Offset:                    uint64(offset),
		}
	}

	return files, nil
}

func buildFS(files []JPkgFileRecordWithOffset) {
	panic("unimplemented")
}
