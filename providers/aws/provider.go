// Package aws - aws provider.
// Common functions of aws provider.
package aws

import (
	"time"

	"github.com/apex/log"

	"gopkg.in/yaml.v2"
)

// Ident - string key for identify provider in the providers map.
const Ident = "aws"

// aws provider sub-config.
type providerConfSpec struct {
	Region            string                 `yaml:"region"`
	Vpc               string                 `yaml:"vpc"`
	Domain            string                 `yaml:"domain"`
	Provisioner       map[string]interface{} `yaml:"provisioner"`
	ProviderType      string                 `yaml:"type"`
	AvailabilityZones []string               `yaml:"availability_zones"`
	ClusterName       string                 `yaml:"cluster_name"`
}

// ModuleCommon - interface for terraform modules instance.
type ModuleCommon interface {
	// Apply module with his defined configuration.
	Deploy() error
	// Destroy infrastructure, created by module.
	Destroy() error
	// Some modules checks.
	Check() (bool, error)
}

// Provider - main provider object.
type Provider struct {
	Config providerConfSpec
}

var test Provider

// Init provider, check config.
func (p *Provider) Init(yamlSpec []byte, clusterName string) error {
	var spec providerConfSpec
	err := yaml.UnmarshalStrict(yamlSpec, &spec)
	if err != nil {
		return err
	}
	p.Config = spec
	p.Config.ClusterName = clusterName
	return nil
}

// Deploy function.
func (p *Provider) Deploy() error {
	// Create bucket.
	backend, err := NewS3Backend(p.Config)
	if err != nil {
		return err
	}

	exists, err := backend.Check()
	if err != nil {
		return err
	}

	if !exists {
		err = backend.Deploy()
		if err != nil {
			return err
		}
	}

	// Deploy DNS.
	route53, err := NewRoute53(p.Config)
	if err != nil {
		return err
	}
	err = route53.Deploy()
	if err != nil {
		return err
	}

	// Deploy VPC.
	vpc, err := NewVpc(p.Config)

	if err != nil {
		return err
	}
	err = vpc.Deploy()
	if err != nil {
		return err
	}
	// Deploy EKS
	provisioner, err := NewProvisioner(p.Config)
	if err != nil {

		return err
	}
	provisioner.Deploy()
	if err != nil {
		return err
	}
	kubeConfig, err := provisioner.GetKubeConfig()
	log.Debugf("%s", kubeConfig)
	err = provisioner.WaitWithTimeout(time.Second * 120)
	if err != nil {
		return err
	}
	return nil
}

// Destroy function.
func (p *Provider) Destroy() error {

	log.Debug("Check if backend bucket exists")

	backend, err := NewS3Backend(p.Config)
	if err != nil {
		return err
	}
	exists, err := backend.Check()
	if err != nil {
		return err
	}

	if !exists {
		log.Infof("Backend bucket '%s' is not found. Nothing to destroy.", p.Config.ClusterName)
		return nil
	}
	// Deploy kubernetes.
	provisioner, err := NewProvisioner(p.Config)
	if err != nil {

		return err
	}
	provisioner.Destroy()
	if err != nil {
		return err
	}

	// Destroy VPC.
	vpc, err := NewVpc(p.Config)
	if err != nil {
		return err
	}
	err = vpc.Destroy()
	if err != nil {
		return err
	}
	// Destroy DNS.
	route53, err := NewRoute53(p.Config)
	if err != nil {
		return err
	}
	err = route53.Destroy()
	if err != nil {
		return err
	}

	// Remove bucket.
	err = backend.Destroy()
	if err != nil {
		return err
	}

	return nil
}
