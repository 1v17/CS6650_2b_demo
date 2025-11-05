output "rds_endpoint" {
  description = "RDS instance endpoint address (hostname only, without port)"
  value       = split(":", aws_db_instance.mysql.endpoint)[0]
}

output "rds_port" {
  description = "RDS instance port"
  value       = aws_db_instance.mysql.port
}

output "db_name" {
  description = "Database name"
  value       = aws_db_instance.mysql.db_name
}

output "db_username" {
  description = "Database master username"
  value       = aws_db_instance.mysql.username
}

output "db_password" {
  description = "Database password (sensitive)"
  value       = random_password.db_password.result
  sensitive   = true
}

output "rds_security_group_id" {
  description = "Security group ID for RDS instance"
  value       = aws_security_group.rds.id
}