# 1. Define the AWS Provider
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region     = "us-east-1"
  access_key = "mock_key"
  secret_key = "mock_secret"

  # skip calling the real AWS APIs for validation
  skip_credentials_validation = true
  skip_requesting_account_id  = true
  skip_metadata_api_check     = true
}

# 2. Create the VPC
resource "aws_vpc" "sre_vpc" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true

  tags = {
    Name = "sre-terraform-vpc"
  }
}

# 3. Create an Internet Gateway (so the VPC can connect to the internet)
resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.sre_vpc.id

  tags = {
    Name = "sre-igw"
  }
}

# 4. Create a Public Subnet
resource "aws_subnet" "public_subnet" {
  vpc_id                  = aws_vpc.sre_vpc.id
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = true
  availability_zone       = "us-east-1a"

  tags = {
    Name = "sre-public-subnet"
  }
}