resource "aws_instance" "vm" {
  ami                    = "ami-00d8fc944fb171e29" # ubuntu
  instance_type          = var.instance_type
  vpc_security_group_ids = [var.sg_id]

  tags = {
    Name = "10ms-demo-ec2"
  }
}
