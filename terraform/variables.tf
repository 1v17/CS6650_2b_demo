# Region to deploy into
variable "aws_region" {
  type    = string
  default = "us-west-2"
}

# ECR & ECS settings
variable "ecr_repository_name" {
  type    = string
  default = "ecr_service"
}

variable "service_name" {
  type    = string
  default = "CS6650L2"
}

variable "container_port" {
  type    = number
  default = 8080
}

variable "ecs_count" {
  type    = number
  default = 1
}

# How long to keep logs
variable "log_retention_days" {
  type    = number
  default = 7
}

# ECS Task CPU
variable "cpu" {
  type        = string
  description = "CPU units for the ECS task"
}

# ECS Task Memory
variable "memory" {
  type        = string
  description = "Memory for the ECS task"
}

# CIDR blocks for security group access
variable "cidr_blocks" {
  type        = list(string)
  description = "CIDR blocks allowed to access the service"
}

# Auto Scaling Configuration
variable "min_capacity" {
  type        = number
  default     = 2
  description = "Minimum number of ECS tasks"
}

variable "max_capacity" {
  type        = number
  default     = 8
  description = "Maximum number of ECS tasks"
}

variable "response_time_threshold" {
  type        = number
  default     = 2
  description = "Response time threshold in seconds for monitoring alarm"
}

variable "error_rate_threshold" {
  type        = number
  default     = 10
  description = "Error rate threshold (count) for monitoring alarm"
}

variable "enable_alb_monitoring" {
  type        = bool
  default     = true
  description = "Enable ALB monitoring alarms (response time and error rate)"
}
