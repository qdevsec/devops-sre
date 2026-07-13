compare building a Virtual private cloud:
- single public subnet
- internet gateway

Use `Terraform` and `Cloudformation (AWS)`


# Terraform:
- uses HCL (hashicorp configuration language)
- relies on providers to translate code to API calls
- cloud-agnostic
- Imperative / Declarative Hybrid
    - Declarative:
        - Terraform manages the state, dependencies, lifecycle and destruction automatically
    - Imperative:
        - Step by step commands telling system exactly how to perform an action
- Terraform keeps tracks of what it builds in `terraform.tfstate`

commands

cd /to/project/where/main.tf is located

$ terraform init

$ terraform plan -out=tfplan  # lock plan into file instead of printing to console

$ terraform apply "tfplan"

clean up after

$ terraform destroy

# Cloudformation (AWS)
- cloud native tool
- JSON / YAML
- Strictly declarative
- With AWS CDK (Imperative Generator) can write imperative code that gets synthesized to a declarative CloudFormation YAML/JSON template
- CloudFormation manages the state implicitly behind the scenes withing the AWS service itself, referred to as a stack


Commands

syntax validation

$ aws cloudformation validate-template --template-body file://template.yaml

Deploy infrastructure

$ aws cloudformation deploy \
    --template-file template.yaml \
    --stack-name sre-network-stack

Clean up

$ aws cloudformation delete-stack --stack-name sre-network-stack


# Sample Output - terraform plan

```
Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # aws_internet_gateway.igw will be created
  + resource "aws_internet_gateway" "igw" {
      + arn      = (known after apply)
      + id       = (known after apply)
      + owner_id = (known after apply)
      + tags     = {
          + "Name" = "sre-igw"
        }
      + tags_all = {
          + "Name" = "sre-igw"
        }
      + vpc_id   = (known after apply)
    }

  # aws_subnet.public_subnet will be created
  + resource "aws_subnet" "public_subnet" {
      + arn                                            = (known after apply)
      + assign_ipv6_address_on_creation                = false
      + availability_zone                              = "us-east-1a"
      + availability_zone_id                           = (known after apply)
      + cidr_block                                     = "10.0.1.0/24"
      + enable_dns64                                   = false
      + enable_resource_name_dns_a_record_on_launch    = false
      + enable_resource_name_dns_aaaa_record_on_launch = false
      + id                                             = (known after apply)
      + ipv6_cidr_block_association_id                 = (known after apply)
      + ipv6_native                                    = false
      + map_public_ip_on_launch                        = true
      + owner_id                                       = (known after apply)
      + private_dns_hostname_type_on_launch            = (known after apply)
      + tags                                           = {
          + "Name" = "sre-public-subnet"
        }
      + tags_all                                       = {
          + "Name" = "sre-public-subnet"
        }
      + vpc_id                                         = (known after apply)
    }

  # aws_vpc.sre_vpc will be created
  + resource "aws_vpc" "sre_vpc" {
      + arn                                  = (known after apply)
      + cidr_block                           = "10.0.0.0/16"
      + default_network_acl_id               = (known after apply)
      + default_route_table_id               = (known after apply)
      + default_security_group_id            = (known after apply)
      + dhcp_options_id                      = (known after apply)
      + enable_dns_hostnames                 = true
      + enable_dns_support                   = true
      + enable_network_address_usage_metrics = (known after apply)
      + id                                   = (known after apply)
      + instance_tenancy                     = "default"
      + ipv6_association_id                  = (known after apply)
      + ipv6_cidr_block                      = (known after apply)
      + ipv6_cidr_block_network_border_group = (known after apply)
      + main_route_table_id                  = (known after apply)
      + owner_id                             = (known after apply)
      + tags                                 = {
          + "Name" = "sre-terraform-vpc"
        }
      + tags_all                             = {
          + "Name" = "sre-terraform-vpc"
        }
    }

Plan: 3 to add, 0 to change, 0 to destroy.
```
