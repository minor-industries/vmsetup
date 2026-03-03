package vmsetup

import (
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/minor-industries/vmsetup/netplan"
)

type HostOpts struct {
	Args struct {
		In string `positional-arg-name:"IN" required:"yes"`
	} `positional-args:"yes"`
}

func (h *HostOpts) Execute(args []string) error {
	in, err := os.ReadFile(h.Args.In)
	if err != nil {
		return err
	}

	out, err := netplan.ToBridgeStrict(in)
	if err != nil {
		return err
	}

	_, err = os.Stdout.Write(out)
	return err
}

type Opts struct {
	MemoryGB int  `short:"m" long:"memory" default:"2"`
	CPUs     int  `short:"c" long:"cpus" default:"2"`
	DiskGB   int  `short:"d" long:"disk" default:"20"`
	Spice    bool `long:"spice"`

	Username string   `long:"username"`
	SSHKeys  []string `long:"ssh-key"`

	CloudImageURL   string `long:"cloud-image-url"`
	CloudConfigHash string `long:"cloud-config-hash"`

	Args struct {
		Name string `positional-arg-name:"NAME" required:"yes"`
	} `positional-args:"yes"`
}

func (v *Opts) Execute(args []string) error {
	return run(v)
}

type Root struct {
	Host HostOpts `command:"host" description:"rewrite netplan ethernet config to bridge"`
	VM   Opts     `command:"vm" description:"create vm"`
}

type Config struct {
	SshKeys  []string
	Username string

	CloudImageURL  string
	CloudImageHash string
}

func Run(cfg *Config) error {
	root := &Root{}
	root.VM.Username = cfg.Username
	root.VM.SSHKeys = cfg.SshKeys
	root.VM.CloudImageURL = cfg.CloudImageURL
	root.VM.CloudConfigHash = cfg.CloudImageHash

	p := flags.NewParser(root, flags.Default)

	_, err := p.Parse()
	return err
}
