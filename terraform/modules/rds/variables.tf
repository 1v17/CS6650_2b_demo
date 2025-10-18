variable "service_name" {
  type        = string
  description = "Base name for RDS resources"
}

variable "vpc_id" {
  type        = string
  description = "VPC ID where RDS will be deployed"
}

variable "subnet_ids" {
  type        = list(string)
  description = "Subnet IDs for RDS deployment (should be private subnets)"
}

variable "ecs_security_group_id" {
  type        = string
  description = "Security group ID of ECS tasks that need database access"
}

variable "db_name" {
  type        = string
  description = "Name of the initial database"
  default     = "ecommerce"
}

variable "db_username" {
  type        = string
  description = "Master username for the database"
  default     = "admin"
}

variable "instance_class" {
  type        = string
  description = "RDS instance class"
  default     = "db.t3.micro"
}

variable "engine_version" {
  type        = string
  description = "MySQL engine version"
  default     = "8.0"
}

variable "backup_retention_period" {
  type        = number
  description = "Number of days to retain automated backups"
  default     = 7
}

variable "skip_final_snapshot" {
  type        = bool
  description = "Whether to skip final snapshot on deletion"
  default     = true
}

variable "deletion_protection" {
  type        = bool
  description = "Whether to enable deletion protection"
  default     = false
}