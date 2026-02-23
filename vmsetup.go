package vmsetup

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jessevdk/go-flags"
)

type Opts struct {
	MemoryGB int  `short:"m" long:"memory" description:"Memory in GB" default:"2"`
	CPUs     int  `short:"c" long:"cpus" description:"Number of vCPUs" default:"2"`
	Spice    bool `long:"spice" description:"Enable SPICE graphics"`
	Args     struct {
		Name string `positional-arg-name:"NAME" required:"yes"`
	} `positional-args:"yes"`
}

func makeSeedISO(
	hostname string,
	outfile string,
	sshKeys []string,
) error {
	ud := &CloudConfig{
		ManageEtc:     true,
		PackageUpdate: false,
		LockPasswd:    true,

		Users: []User{
			{
				Name:              "jsu",
				Sudo:              "ALL=(ALL) NOPASSWD:ALL",
				Groups:            []string{"sudo"},
				Shell:             "/bin/bash",
				SSHAuthorizedKeys: sshKeys,
			},
		},

		RunCmd: []string{},
	}

	md := &MetaData{
		InstanceID:    hostname,
		LocalHostname: hostname,
	}

	nc := &NetworkConfig{
		Version: 2,
		Ethernets: map[string]Ethernet{
			"enp1s0": {
				DHCP4: true,
			},
		},
	}

	udOut, err := ud.MarshalYAML()
	if err != nil {
		return fmt.Errorf("marshal ud: %w", err)
	}

	mdOut, err := md.MarshalYAML()
	if err != nil {
		return fmt.Errorf("marshal md: %w", err)
	}

	ncOut, err := nc.MarshalYAML()
	if err != nil {
		return fmt.Errorf("marshal nc: %w", err)
	}

	err = writeCloudInitSeedISO(
		outfile,
		udOut,
		mdOut,
		ncOut,
	)
	if err != nil {
		return fmt.Errorf("write cloud-init seed: %w", err)
	}

	return nil
}

const (
	img          = "https://cloud-images.ubuntu.com/noble/20260108/noble-server-cloudimg-amd64.img"
	expectedHash = "00786c0936a7dd91a6b07941ca60bb56652975e0e72f9dacf73c887ada420966"
)

func run(
	vmName string,
	memGB, cpus int,
	spice bool,
	sshKeys []string,
) error {
	base := filepath.Base(img)

	backingFile := "/var/lib/libvirt/images/" + base
	_, err := os.Stat(backingFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%s does not exist, downloading\n", backingFile)
			if err := download(context.Background(), img, backingFile); err != nil {
				return fmt.Errorf("download: %w", err)
			}
		}
	}

	hash, err := HashFileSHA256(backingFile)
	if err != nil {
		return fmt.Errorf("hash: %w", err)
	}

	if hash != expectedHash {
		return fmt.Errorf("expected hash %s, got %s", expectedHash, hash)
	}

	overlay := fmt.Sprintf("/var/lib/libvirt/images/%s.qcow2", vmName)
	seedISO := fmt.Sprintf("/var/lib/libvirt/images/%s-seed.iso", vmName)

	_, err = os.Stat(overlay)
	if err == nil {
		return fmt.Errorf("overlay %s already exists", overlay)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("check overlay exists: %w", err)
	}

	cmd := exec.Command(
		"qemu-img",
		"create",
		"-f",
		"qcow2",
		"-b",
		backingFile,
		"-F",
		"qcow2",
		overlay,
		"20G",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("qemu-img: %w: %s", err, string(out))
	}

	if err := makeSeedISO(vmName, seedISO, sshKeys); err != nil {
		return fmt.Errorf("make seed iso: %w", err)
	}

	if err := chown("libvirt-qemu", "libvirt-qemu", backingFile, overlay, seedISO); err != nil {
		return fmt.Errorf("chown: %w", err)
	}

	memMB := memGB * 1024
	args := []string{
		"virt-install",
		"--name", vmName,
		"--memory", fmt.Sprint(memMB),
		"--vcpus", fmt.Sprint(cpus),
		"--disk", fmt.Sprintf("path=%s,format=qcow2", overlay),
		"--disk", fmt.Sprintf("path=%s,device=cdrom", seedISO),
		"--os-variant", "ubuntu24.04",
		"--network", "bridge=br0",
		"--import",
	}

	if spice {
		args = append(args,
			"--graphics", "spice,listen=127.0.0.1",
			"--video", "qxl",
		)
	} else {
		args = append(args, "--graphics", "none")
	}

	fmt.Println(strings.Join(args, " "))
	return nil
}

func Run(sshKeys []string) error {
	var opts Opts

	p := flags.NewParser(&opts, flags.Default)

	if _, err := p.Parse(); err != nil {
		return fmt.Errorf("parse options: %w", err)
	}

	if err := run(opts.Args.Name, opts.MemoryGB, opts.CPUs, opts.Spice, sshKeys); err != nil {
		return fmt.Errorf("run: %w", err)
	}

	return nil
}
