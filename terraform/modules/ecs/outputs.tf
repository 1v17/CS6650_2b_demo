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
  value       = aws_appautoscaling_target.ecs_target.resource_id
}

output "scale_out_policy_arn" {
  description = "Scale out policy ARN"
  value       = aws_appautoscaling_policy.scale_out.arn
}

output "scale_in_policy_arn" {
  description = "Scale in policy ARN"
  value       = aws_appautoscaling_policy.scale_in.arn
}

output "high_cpu_alarm_name" {
  description = "High CPU alarm name"
  value       = aws_cloudwatch_metric_alarm.high_cpu.alarm_name
}

output "low_cpu_alarm_name" {
  description = "Low CPU alarm name"
  value       = aws_cloudwatch_metric_alarm.low_cpu.alarm_name
}

output "high_response_time_alarm_name" {
  description = "High response time alarm name (if enabled)"
  value       = var.enable_alb_monitoring && var.target_group_arn != null && var.load_balancer_arn_suffix != null && var.target_group_arn != "" && var.load_balancer_arn_suffix != "" ? aws_cloudwatch_metric_alarm.high_response_time[0].alarm_name : null
}

output "error_rate_alarm_name" {
  description = "Error rate alarm name (if enabled)"
  value       = var.enable_alb_monitoring && var.target_group_arn != null && var.load_balancer_arn_suffix != null && var.target_group_arn != "" && var.load_balancer_arn_suffix != "" ? aws_cloudwatch_metric_alarm.error_rate[0].alarm_name : null
}
