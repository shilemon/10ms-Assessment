terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.0"
    }
  }
}

provider "aws" {
  region = var.region
}

# ---------------------------------------
# Security Group
# ---------------------------------------
resource "aws_security_group" "sg" {
  name = "demo-sg"

  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# ---------------------------------------
# EC2 Module
# ---------------------------------------
module "ec2_module" {
  source        = "./modules/ec2"
  instance_type = var.instance_type
  sg_id         = aws_security_group.sg.id
}

