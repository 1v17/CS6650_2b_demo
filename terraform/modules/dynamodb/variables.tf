variable "service_name" {
  type        = string
  description = "Name of the service"
}

variable "environment" {
  type        = string
  description = "Environment name"
  default     = "dev"
}

variable "enable_point_in_time_recovery" {
  type        = bool
  description = "Enable point-in-time recovery for DynamoDB table"
  default     = false
}
