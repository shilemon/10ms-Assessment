1. Wrong Variable Types
Problem:

One or more variables in the root module used incorrect types.
Example issues:

A variable defined as a string but used as a list

A boolean written as a number

An AWS region expected as string but provided incorrectly

Fix Applied:

Corrected the variable definitions

Ensured types matched module inputs

Aligned variable default values with actual usage in resources

2. Missing Variables
Problem:

The module modules/ec2 referenced variables that were never declared at the root module.
This caused:

Terraform plan failure

Module failing to receive required AMI, instance_type, subnet_id, etc.

Fix Applied:

Added missing variables in root variables.tf

Passed the correct values into the ec2 module

Ensured variable names matched exactly between module and root

3. Non-Existent Subnet
Problem:

The EC2 instance attempted to launch into a subnet ID that did not exist.
This caused:

"Subnet ID not found" AWS API errors

Module failing to create EC2 instance

Fix Applied:

Replaced the invalid subnet with a valid one

Added variable validation to prevent incorrect IDs

Ensured the subnet belonged to the correct VPC

4. Invalid AMI
Problem:

The AMI used in the EC2 module was outdated or not available in the specified region.
Terraform failed with:

“InvalidAMIID.NotFound”

Fix Applied:

Selected an updated Amazon Linux 2 or Ubuntu LTS AMI

Ensured AMI was region-specific

Passed the correct AMI ID as module input

5. Module Missing Inputs
Problem:

The ec2 module expected certain variables but they were never passed from the root module.
Examples:

Missing instance_type

Missing security_group_ids

Missing key_name

Fix Applied:

Added the required inputs to the module block

Ensured names matched exactly

Validated all inputs with Terraform validate

6. Incorrect Module Output Referencing
Problem:

Root module attempted to output values from the module that did not exist.
Example:

Referencing module.ec2.public_ip when module output was named differently

Output block referencing outdated variable names

Fix Applied:

Corrected output names

Updated root output references to match module output definitions

Removed outputs pointing to non-existent resources

7. Security Group Exposed ALL Ports
Problem:

The security group allowed:

0.0.0.0/0

all ports open

all protocols allowed
This is a severe security vulnerability.

Fix Applied:

Limited inbound rules to required ports only (SSH / App port)

Limited outbound rules to standard egress

Added variable for allowed IP ranges

Properly hardened the SG

8. Deprecated AWS Provider Version
Problem:

The Terraform AWS provider was using an outdated version.
This caused:

Syntax incompatibility

Resource behavior issues

AWS API mismatches

Fix Applied:

Updated provider block to a stable, supported version

Added required providers block

Ran terraform init -upgrade

9. Root Module Output Referencing Non-Existent Output
Problem:

The root outputs.tf referenced outputs that were removed or renamed in the EC2 module.

Fix Applied:

Reviewed all module outputs

Removed unused outputs

Aligned root outputs with module’s real output names

10. Module Structure Validation
Problem:

Terraform module folder structure had inconsistent organization and missing required files.

Fix Applied:

Validated module folder structure

Ensured main.tf, variables.tf, and outputs.tf existed

Ensured no circular dependencies

Run terraform init, validate, plan, and applied fixes

Final Result

After all fixes:

Terraform initializes and validates successfully

Plan succeeds with no errors

EC2 instance deploys correctly

Security group follows least-privilege standards

Module inputs/outputs fully aligned

AWS provider up-to-date

Code is production-ready and secure
