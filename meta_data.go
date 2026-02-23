package vmsetup

import (
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

type MetaData struct {
	InstanceID    string `yaml:"instance-id"`
	LocalHostname string `yaml:"local-hostname,omitempty"`
}

func (cc *MetaData) MarshalYAML() (string, error) {
	// Marshal YAML
	out, err := yaml.Marshal(cc)
	if err != nil {
		return "", errors.Wrap(err, "marshal cloud-config")
	}

	// cloud-init requires this header
	return string(out), nil
}
