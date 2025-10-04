output "cluster_name" {
  description = "ECS cluster name"
  value       = aws_ecs_cluster.this.name
}

output "service_name" {
  description = "ECS service name"
  value       = aws_ecs_service.this.name
}

output "autoscaling_target_resource_id" {
  description = "Auto scaling target resource ID"
  value       = var.enable_auto_scaling ? aws_appautoscaling_target.ecs_target[0].resource_id : null
}

output "scale_out_policy_arn" {
  description = "Scale out policy ARN"
  value       = var.enable_auto_scaling ? aws_appautoscaling_policy.scale_out[0].arn : null
}

output "scale_in_policy_arn" {
  description = "Scale in policy ARN"
  value       = var.enable_auto_scaling ? aws_appautoscaling_policy.scale_in[0].arn : null
}

output "high_cpu_alarm_name" {
  description = "High CPU alarm name"
  value       = var.enable_auto_scaling ? aws_cloudwatch_metric_alarm.high_cpu[0].alarm_name : null
}

output "low_cpu_alarm_name" {
  description = "Low CPU alarm name"
  value       = var.enable_auto_scaling ? aws_cloudwatch_metric_alarm.low_cpu[0].alarm_name : null
}

output "high_response_time_alarm_name" {
  description = "High response time alarm name (if enabled)"
  value       = var.enable_alb_monitoring && var.target_group_arn != null && var.load_balancer_arn_suffix != null && var.target_group_arn != "" && var.load_balancer_arn_suffix != "" ? aws_cloudwatch_metric_alarm.high_response_time[0].alarm_name : null
}

output "error_rate_alarm_name" {
  description = "Error rate alarm name (if enabled)"
  value       = var.enable_alb_monitoring && var.target_group_arn != null && var.load_balancer_arn_suffix != null && var.target_group_arn != "" && var.load_balancer_arn_suffix != "" ? aws_cloudwatch_metric_alarm.error_rate[0].alarm_name : null
}
