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
    cidr_blocks = ["0.0.0.0/0"]
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

# Additional security group rule to allow ALB to reach ECS tasks
resource "aws_security_group_rule" "alb_to_ecs" {
  type                     = "ingress"
  from_port               = var.container_port
  to_port                 = var.container_port
  protocol                = "tcp"
  source_security_group_id = var.alb_security_group_id
  security_group_id       = aws_security_group.this.id
}
