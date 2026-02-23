package vmsetup

import (
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

type NetworkConfig struct {
	Version   int                 `yaml:"version"`
	Ethernets map[string]Ethernet `yaml:"ethernets"`
}

type Ethernet struct {
	DHCP4 bool `yaml:"dhcp4,omitempty"`
	DHCP6 bool `yaml:"dhcp6,omitempty"`

	Match   *Match `yaml:"match,omitempty"`
	SetName string `yaml:"set-name,omitempty"`
}

type Match struct {
	MACAddress string `yaml:"macaddress,omitempty"`
}

func (nc *NetworkConfig) MarshalYAML() (string, error) {
	out, err := yaml.Marshal(nc)
	if err != nil {
		return "", errors.Wrap(err, "marshal network-config")
	}

	return string(out), nil
}
