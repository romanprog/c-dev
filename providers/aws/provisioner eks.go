package aws

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/apex/log"
	"github.com/romanprog/c-dev/executor"
)

// ProvisionerEks class.
type ProvisionerEks struct {
	providerConf providerConfSpec
	kubeConfig   string
	eksModule    *Eks
}

// NewProvisionerEks create new instance of EKS provisioner.
func NewProvisionerEks(providerConf providerConfSpec) (*ProvisionerEks, error) {
	var provisioner ProvisionerEks
	provisioner.providerConf = providerConf
	eksMod, err := NewEks(providerConf)
	if err != nil {
		return nil, err
	}
	// TODO check config.
	provisioner.eksModule = eksMod
	return &provisioner, nil
}

// Deploy EKS provisioner modules, save kubernetes config to kubeConfig string.
// Upload kube config to s3.
func (p *ProvisionerEks) Deploy(timeout time.Duration) error {
	err := p.eksModule.Deploy()
	if err != nil {
		return err
	}
	// kube config file path.
	kubeConfigFileName := "kubeconfig_" + p.providerConf.ClusterName
	kubeConfigFile := filepath.Join(p.eksModule.ModulePath(), kubeConfigFileName)
	// Read kube confin from file to string.
	conf, err := ioutil.ReadFile(kubeConfigFile)
	if err != nil {
		return err
	}
	p.kubeConfig = string(conf)

	// Upload kube config to s3 bucket.

	// Init bash runner in module director and export path to kubeConfig.
	varString := fmt.Sprintf("KUBECONFIG=%s", kubeConfigFile)
	bash, err := executor.NewBashRunner(p.eksModule.ModulePath(), varString)
	if err != nil {
		return err
	}
	// Ticker for pause.
	tm := time.After(timeout)
	var tick = make(<-chan time.Time)
	tick = time.Tick(5 * time.Second)
	for {
		select {
		case <-tm:
			// Timeout
			return fmt.Errorf("k8s access timeout")
		// Wait for tick.
		case <-tick:
			// Test k8s access.
			stdout, stderr, err := bash.RunMutely("kubectl version --request-timeout=5s")
			if err == nil {
				// Connected! k8s is ready.
				// Upload kubeConfig to s3
				// https://golang.org/pkg/fmt/ (see "Explicit argument indexes")
				uploadCommand := fmt.Sprintf("aws s3 cp '%[1]s' 's3://%[2]s/%[1]s", kubeConfigFileName, p.providerConf.ClusterName)
				return bash.Run(uploadCommand)
			}
			log.Debugf("%s %s", stdout, stderr)
			// Connection error. Wait for next tick (try).
		}
	}
}

// Destroy EKS provisioner objects.
func (p *ProvisionerEks) Destroy() error {
	err := p.eksModule.Destroy()
	if err != nil {
		return err
	}
	p.kubeConfig = ""
	return nil
}

// GetKubeConfig return 'kubeConfig' or error if config is empty.
func (p *ProvisionerEks) GetKubeConfig() (string, error) {
	if p.kubeConfig == "" {
		return "", fmt.Errorf("eks kube config is empty")
	}
	return p.kubeConfig, nil
}
