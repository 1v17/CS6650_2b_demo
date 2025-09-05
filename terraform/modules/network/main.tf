# Fetch the default VPC
data "aws_vpc" "default" {
  default = true
}

# List all subnets in that VPC
data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

# Create a security group to allow HTTP to your container port
resource "aws_security_group" "this" {
  name        = "${var.service_name}-sg"
  description = "Allow inbound on ${var.container_port}"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    from_port   = var.container_port
    to_port     = var.container_port
    protocol    = "tcp"
    cidr_blocks = var.cidr_blocks
    description = "Allow HTTP traffic"
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow all outbound"
  }
}

# Additional security group rule to allow ALB traffic
resource "aws_security_group_rule" "alb_ingress" {
  type                     = "ingress"
  from_port               = 8080
  to_port                 = 8080
  protocol                = "tcp"
  cidr_blocks             = ["0.0.0.0/0"]
  security_group_id       = var.alb_security_group_id
  
  depends_on = [var.alb_security_group_id]
}
