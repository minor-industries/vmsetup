package vmsetup

import (
	"fmt"
	"os"

	"github.com/diskfs/go-diskfs"
	"github.com/diskfs/go-diskfs/disk"
	"github.com/diskfs/go-diskfs/filesystem"
	"github.com/diskfs/go-diskfs/filesystem/iso9660"
)

func writeCloudInitSeedISO(
	isoPath string,
	userData string,
	metaData string,
	networkConfig string,
) error {
	const diskSize int64 = 8 * 1024 * 1024
	const isoBlockSize int64 = 2048

	d, err := diskfs.Create(isoPath, diskSize, diskfs.SectorSizeDefault)
	if err != nil {
		return fmt.Errorf("create disk image: %w", err)
	}

	d.LogicalBlocksize = isoBlockSize

	fs, err := d.CreateFilesystem(disk.FilesystemSpec{
		Partition:   0,
		FSType:      filesystem.TypeISO9660,
		VolumeLabel: "cidata",
	})
	if err != nil {
		return fmt.Errorf("create iso9660 filesystem: %w", err)
	}

	if err := writeFile(fs, "/user-data", userData); err != nil {
		return err
	}

	if err := writeFile(fs, "/meta-data", metaData); err != nil {
		return err
	}

	if err := writeFile(fs, "/network-config", networkConfig); err != nil {
		return err
	}

	iso, ok := fs.(*iso9660.FileSystem)
	if !ok {
		return fmt.Errorf("not an iso9660 filesystem (got %T)", fs)
	}

	if err := iso.Finalize(iso9660.FinalizeOptions{
		RockRidge:        true,
		VolumeIdentifier: "cidata",
	}); err != nil {
		return fmt.Errorf("finalize iso: %w", err)
	}

	return nil
}

func writeFile(fs filesystem.FileSystem, path, contents string) error {
	f, err := fs.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	if _, err := f.Write([]byte(contents)); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
