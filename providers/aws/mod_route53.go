package aws

import (
	"github.com/apex/log"
	"github.com/romanprog/c-dev/executor"
)

// Variables set for route53 module tfvars.
type route53VarsSpec struct {
	Region         string `json:"region"`
	ClusterName    string `json:"cluster_name"`
	ClusterDomain  string `json:"cluster_domain"`
	ZoneDelegation string `json:"zone_delegation"`
}

// Route53 type for route53 module instance.
type Route53 struct {
	config      route53VarsSpec
	backendConf executor.BackendSpec
	terraform   *executor.TerraformRunner
}

const route53ModulePath = "terraform/aws/route53"
const route53ModuleBackendKey = "states/terraform-dns.state"

// NewRoute53 create new route53 instance.
func NewRoute53(providerConf providerConfSpec) (*Route53, error) {
	var route53 Route53
	route53.backendConf = executor.BackendSpec{
		Bucket: providerConf.ClusterName,
		Key:    route53ModuleBackendKey,
		Region: providerConf.Region,
	}
	zoneDelegation := "false"
	if providerConf.Domain == "cluster.dev" {
		zoneDelegation = "true"
	}
	route53.config = route53VarsSpec{
		Region:         providerConf.Region,
		ClusterName:    providerConf.ClusterName,
		ClusterDomain:  providerConf.Domain,
		ZoneDelegation: zoneDelegation,
	}
	var err error
	route53.terraform, err = executor.NewTerraformRunner(route53ModulePath)
	if err != nil {
		return nil, err
	}
	return &route53, nil
}

// Deploy - create route53.
func (r *Route53) Deploy() error {
	// sss
	log.Debug("Terraform init/plan.")
	err := r.terraform.Clear()
	if err != nil {
		return err
	}
	// Init terraform without backend speck.
	r.terraform.Init(r.backendConf)
	if err != nil {
		return err
	}
	// Plan.
	r.terraform.Plan(r.config, "-compact-warnings", "-out=tfplan")
	if err != nil {
		return err
	}
	// Apply. Create DNS.
	err = r.terraform.ApplyPlan("tfplan", "-compact-warnings")
	if err != nil {
		return err
	}
	return nil
}

// Destroy - remove s3 bucket.
func (r *Route53) Destroy() error {
	// Init terraform without backend speck.
	err := r.terraform.Init(r.backendConf)
	if err != nil {
		return err
	}
	// Plan.
	return r.terraform.Destroy(r.config, "-compact-warnings")
}

// Check - if s3 bucket exists.
func (r *Route53) Check() (bool, error) {
	return true, nil
}
