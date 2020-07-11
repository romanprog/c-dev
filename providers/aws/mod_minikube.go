package aws

import (
	"fmt"

	"github.com/apex/log"
	"github.com/romanprog/c-dev/executor"
)

// Variables set for minikube module tfvars.
type minikubeVarsSpec struct {
	HostedZone      string `json:"hosted_zone"`
	Region          string `json:"region"`
	ClusterName     string `json:"cluster_name"`
	AwsInstanceType string `json:"aws_instance_type"`
}

// Minikube type for minikube module instance.
type Minikube struct {
	config      minikubeVarsSpec
	backendConf executor.BackendSpec
	terraform   *executor.TerraformRunner
}

const minikubeModulePath = "terraform/aws/vpc"
const minikubeModuleBackendKey = "states/terraform-k8s.state"

// NewMinikube create new minikube instance.
func NewMinikube(providerConf providerConfSpec) (*Minikube, error) {
	var bk Minikube
	bk.backendConf = executor.BackendSpec{
		Bucket: providerConf.ClusterName,
		Key:    minikubeModuleBackendKey,
		Region: providerConf.Region,
	}
	instanceType, ok := providerConf.Provisioner["instanceType"].(string)
	if !ok {
		return nil, fmt.Errorf("can't determinate instance type for minikube")
	}
	bk.config = minikubeVarsSpec{
		HostedZone:      fmt.Sprintf("%s.%s", providerConf.ClusterName, providerConf.Domain),
		Region:          providerConf.Region,
		ClusterName:     providerConf.ClusterName,
		AwsInstanceType: instanceType,
	}
	var err error
	bk.terraform, err = executor.NewTerraformRunner(minikubeModulePath)
	if err != nil {
		return nil, err
	}
	return &bk, nil
}

// Deploy - create vpc.
func (s *Minikube) Deploy() error {
	// sss
	log.Debug("Terraform init/plan.")
	err := s.terraform.Clear()
	if err != nil {
		return err
	}
	// Init terraform without backend speck.
	s.terraform.Init(s.backendConf)
	if err != nil {
		return err
	}
	// Plan.
	s.terraform.Plan(s.config, "-compact-warnings", "-out=tfplan")
	if err != nil {
		return err
	}
	// Apply. Create DNS.
	err = s.terraform.ApplyPlan("tfplan", "-compact-warnings")
	if err != nil {
		return err
	}
	return nil
}

// Destroy - remove vpc.
func (s *Minikube) Destroy() error {
	// Init terraform without backend speck.
	err := s.terraform.Init(s.backendConf)
	if err != nil {
		return err
	}
	// Plan.
	return s.terraform.Destroy(s.config, "-compact-warnings")
}

// Check - if s3 bucket exists.
func (s *Minikube) Check() (bool, error) {
	return true, nil
}
