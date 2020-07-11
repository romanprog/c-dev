package aws

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/apex/log"
	"github.com/romanprog/c-dev/executor"
)

// ProvisionerMinikube class.
type ProvisionerMinikube struct {
	providerConf   providerConfSpec
	kubeConfig     string
	minikubeModule *Minikube
}

// NewProvisionerMinikube create new instance of EKS provisioner.
func NewProvisionerMinikube(providerConf providerConfSpec) (*ProvisionerMinikube, error) {
	var provisioner ProvisionerMinikube
	provisioner.providerConf = providerConf
	minikubeMod, err := NewMinikube(providerConf)
	if err != nil {
		return nil, err
	}
	// TODO check config.
	provisioner.minikubeModule = minikubeMod
	return &provisioner, nil
}

// Deploy EKS provisioner modules, save kubernetes config to kubeConfig string.
// Upload kube config to s3.
func (p *ProvisionerMinikube) Deploy(timeout time.Duration) error {
	err := p.minikubeModule.Deploy()
	if err != nil {
		return err
	}
	// kube config file path.
	kubeConfigFileName := "kubeconfig_" + p.providerConf.ClusterName
	kubeConfigFile := filepath.Join(p.minikubeModule.ModulePath(), kubeConfigFileName)

	// Init bash runner in module director and export path to kubeConfig.
	varString := fmt.Sprintf("KUBECONFIG=%s", kubeConfigFile)
	bash, err := executor.NewBashRunner(p.minikubeModule.ModulePath(), varString)
	if err != nil {
		return err
	}
	// Ticker for pause and timeout.
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
			// Download kube config (try)
			downloadCommand := fmt.Sprintf("aws s3 cp 's3://%s/%s' '%s'", p.providerConf.ClusterName, kubeConfigFileName, kubeConfigFile)
			stdout, stderr, err := bash.RunMutely(downloadCommand)
			if err != nil {
				log.Info("Minikube cluster is not ready yet. Will retry after 5 seconds...")
				continue
			}
			kubeconfig, err := ioutil.ReadFile(kubeConfigFile)
			if err != nil {
				return err
			}
			p.kubeConfig = string(kubeconfig)
			//log.Debugf("Kubeconfig: %v", kubeconfig)
			stdout, stderr, err = bash.RunMutely("kubectl version --request-timeout=5s")
			if err == nil {
				// Connected! k8s is ready.
				return nil
			}
			log.Info("Minikube cluster is not ready yet. Will retry after 5 seconds...")
			log.Debugf("Error check kubectl version: %s %s", stdout, stderr)
			// Connection error. Wait for next tick (try).
		}
	}
}

// Destroy minikube provisioner objects.
func (p *ProvisionerMinikube) Destroy() error {
	err := p.minikubeModule.Destroy()
	if err != nil {
		return err
	}
	p.kubeConfig = ""
	return nil
}

// GetKubeConfig return 'kubeConfig' or error if config is empty.
func (p *ProvisionerMinikube) GetKubeConfig() (string, error) {
	if p.kubeConfig == "" {
		return "", fmt.Errorf("minikube kube config is empty")
	}
	return p.kubeConfig, nil
}
