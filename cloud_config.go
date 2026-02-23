package vmsetup

import (
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

type CloudConfig struct {
	Hostname      string `yaml:"hostname,omitempty"`
	FQDN          string `yaml:"fqdn,omitempty"`
	ManageEtc     bool   `yaml:"manage_etc_hosts,omitempty"`
	PackageUpdate bool   `yaml:"package_update,omitempty"`

	Users    []User   `yaml:"users,omitempty"`
	Packages []string `yaml:"packages,omitempty"`

	RunCmd     []string `yaml:"runcmd,omitempty"`
	LockPasswd bool     `yaml:"lock_passwd"`
}

type User struct {
	Name              string   `yaml:"name"`
	Sudo              string   `yaml:"sudo,omitempty"`
	Groups            []string `yaml:"groups,omitempty"`
	Shell             string   `yaml:"shell,omitempty"`
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys,omitempty"`
}

func (cc *CloudConfig) MarshalYAML() (string, error) {
	out, err := yaml.Marshal(cc)
	if err != nil {
		return "", errors.Wrap(err, "marshal cloud-config")
	}

	// cloud-init requires this header
	return "#cloud-config\n" + string(out), nil
}
