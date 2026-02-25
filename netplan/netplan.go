package netplan

import (
	"fmt"

	"github.com/goccy/go-yaml"
)

type InFile struct {
	Network InNetwork `yaml:"network"`
}

type InNetwork struct {
	Version   int                   `yaml:"version"`
	Ethernets map[string]InEthernet `yaml:"ethernets"`
}

type InEthernet struct {
	DHCP4 *bool `yaml:"dhcp4"`
}

type OutFile struct {
	Network OutNetwork `yaml:"network"`
}

type OutNetwork struct {
	Version   int                    `yaml:"version"`
	Ethernets map[string]OutEthernet `yaml:"ethernets"`
	Bridges   map[string]OutBridge   `yaml:"bridges"`
}

type OutEthernet struct{}

type OutBridge struct {
	Interfaces []string `yaml:"interfaces"`
	DHCP4      bool     `yaml:"dhcp4"`
}

func ToBridgeStrict(inYAML []byte) ([]byte, error) {
	var in InFile
	if err := yaml.UnmarshalWithOptions(
		inYAML,
		&in,
		yaml.DisallowUnknownField(),
	); err != nil {
		return nil, fmt.Errorf("decode (strict): %w", err)
	}

	if in.Network.Version != 2 {
		return nil, fmt.Errorf("network.version must be 2 (got %d)", in.Network.Version)
	}
	if in.Network.Ethernets == nil {
		return nil, fmt.Errorf("network.ethernets is required")
	}
	if len(in.Network.Ethernets) != 1 {
		return nil, fmt.Errorf("network.ethernets must contain exactly one interface (got %d)", len(in.Network.Ethernets))
	}
	eno, ok := in.Network.Ethernets["eno1"]
	if !ok {
		return nil, fmt.Errorf(`network.ethernets must contain only "eno1"`)
	}
	if eno.DHCP4 == nil {
		return nil, fmt.Errorf(`network.ethernets.eno1.dhcp4 is required`)
	}

	out := OutFile{
		Network: OutNetwork{
			Version: in.Network.Version,
			Ethernets: map[string]OutEthernet{
				"eno1": {},
			},
			Bridges: map[string]OutBridge{
				"br0": {
					Interfaces: []string{"eno1"},
					DHCP4:      *eno.DHCP4,
				},
			},
		},
	}

	b, err := yaml.Marshal(out)
	if err != nil {
		return nil, fmt.Errorf("encode: %w", err)
	}
	return b, nil
}
