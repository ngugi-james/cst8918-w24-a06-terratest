package test

import (
	"testing"
	"strings"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// You normally want to run this under a separate "Testing" subscription
// For lab purposes you will use your assigned subscription under the Cloud Dev/Ops program tenant
var subscriptionID string = "b3dac7bd-3af2-4f39-8244-cff239d62f24"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix": "wair0001",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of output variable
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

	// Confirm VM exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))
}

func TestAzureNICAttachedToVM(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: "../",
	}

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	expectedNicName := terraform.Output(t, terraformOptions, "nic_name")

	// Retrieve NICs attached to the VM
	nicIDs := azure.GetVirtualMachineNics(t, vmName, resourceGroupName, subscriptionID)

	// Extract NIC names
	var nicNames []string
	for _, nicID := range nicIDs {
		parts := strings.Split(nicID, "/")
		nicNames = append(nicNames, parts[len(parts)-1])
	}

	// Validate NIC attachment
	assert.Contains(t, nicNames, expectedNicName, "NIC is not attached to the VM")
}

func TestAzureLinuxVMOSVersion(t *testing.T) {  // Function name updated
	terraformOptions := &terraform.Options{
		TerraformDir: "../",
	}

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// Get Terraform output values
	vmVersion := terraform.Output(t, terraformOptions, "vm_version")
	vmSku := terraform.Output(t, terraformOptions, "vm_sku")

	// Expected values for Ubuntu 22.04 LTS
	expectedSku := "22_04-lts-gen2"
	expectedVersion := "latest" // or you can match the latest stable version

	// Validate that the VM is using the expected OS version
	assert.Equal(t, expectedSku, vmSku, "Expected VM SKU to be Ubuntu 22.04 LTS")
	assert.Equal(t, expectedVersion, vmVersion, "Expected VM version to be 'latest'")
}