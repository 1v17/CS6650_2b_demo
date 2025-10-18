output "ecs_cluster_name" {
  description = "Name of the created ECS cluster"
  value       = module.ecs.cluster_name
}

output "ecs_service_name" {
  description = "Name of the running ECS service"
  value       = module.ecs.service_name
}

output "load_balancer_dns_name" {
  description = "DNS name of the Application Load Balancer"
  value       = module.alb.load_balancer_dns_name
}

# Auto Scaling Outputs
output "autoscaling_target_resource_id" {
  description = "Auto scaling target resource ID"
  value       = module.ecs.autoscaling_target_resource_id
}

output "scale_out_policy_arn" {
  description = "Scale out policy ARN"
  value       = module.ecs.scale_out_policy_arn
}

output "scale_in_policy_arn" {
  description = "Scale in policy ARN"
  value       = module.ecs.scale_in_policy_arn
}

output "high_cpu_alarm_name" {
  description = "High CPU alarm name"
  value       = module.ecs.high_cpu_alarm_name
}

output "low_cpu_alarm_name" {
  description = "Low CPU alarm name"
  value       = module.ecs.low_cpu_alarm_name
}

output "high_response_time_alarm_name" {
  description = "High response time alarm name (if enabled)"
  value       = module.ecs.high_response_time_alarm_name
}

output "error_rate_alarm_name" {
  description = "Error rate alarm name (if enabled)"
  value       = module.ecs.error_rate_alarm_name
}

# RDS Outputs
output "rds_endpoint" {
  description = "RDS instance endpoint address"
  value       = module.rds.rds_endpoint
}

output "rds_port" {
  description = "RDS instance port"
  value       = module.rds.rds_port
}

output "db_name" {
  description = "Database name"
  value       = module.rds.db_name
}

output "db_username" {
  description = "Database master username"
  value       = module.rds.db_username
}

output "rds_security_group_id" {
  description = "Security group ID for RDS instance"
  value       = module.rds.rds_security_group_id
}