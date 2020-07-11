package aws

import (
	"fmt"
	"time"
)

// ProvisionerCommon - interface for all provisioners.
type ProvisionerCommon interface {
	Deploy() error
	Destroy() error
	GetKubeConfig() (string, error)
	WaitWithTimeout(timeout time.Duration) error
}

// NewProvisioner create new provisioner instance.
func NewProvisioner(conf providerConfSpec) (ProvisionerCommon, error) {

	provisionerType, ok := conf.Provisioner["type"].(string)
	if !ok {
		return nil, fmt.Errorf("can't determinate provisioner type")
	}
	switch provisionerType {
	case "minikube":
		return nil, nil
	case "eks":
		var pv ProvisionerCommon
		pv, err := NewProvisionerEks(conf)
		return pv, err
	default:

	}
	return nil, fmt.Errorf("unknown provisioner type '%s'", provisionerType)
}
