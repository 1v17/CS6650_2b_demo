# Generate a random password for the database
resource "random_password" "db_password" {
  length  = 16
  special = true
}

# Create DB subnet group
resource "aws_db_subnet_group" "this" {
  name       = "${lower(var.service_name)}-db-subnet-group"
  subnet_ids = var.subnet_ids

  tags = {
    Name = "${var.service_name}-db-subnet-group"
  }
}

# Security group for RDS
resource "aws_security_group" "rds" {
  name        = "${lower(var.service_name)}-rds-sg"
  description = "Security group for RDS MySQL database"
  vpc_id      = var.vpc_id

  # Allow MySQL access from ECS tasks
  ingress {
    from_port       = 3306
    to_port         = 3306
    protocol        = "tcp"
    security_groups = [var.ecs_security_group_id]
    description     = "MySQL access from ECS tasks"
  }

  # Allow all outbound traffic
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    description = "All outbound traffic"
  }

  tags = {
    Name = "${var.service_name}-rds-sg"
  }
}

# RDS MySQL instance
resource "aws_db_instance" "mysql" {
  identifier = "${lower(var.service_name)}-mysql-db"
  
  # Engine configuration
  engine         = "mysql"
  engine_version = var.engine_version
  instance_class = var.instance_class
  
  # Storage configuration
  allocated_storage     = 20
  max_allocated_storage = 100
  storage_type          = "gp2"
  storage_encrypted     = true
  
  # Database configuration
  db_name  = var.db_name
  username = var.db_username
  password = random_password.db_password.result
  
  # Network configuration
  db_subnet_group_name   = aws_db_subnet_group.this.name
  vpc_security_group_ids = [aws_security_group.rds.id]
  publicly_accessible    = false
  
  # Backup configuration
  backup_retention_period = var.backup_retention_period
  backup_window          = "03:00-04:00"
  maintenance_window     = "sun:04:00-sun:05:00"
  
  # Deletion configuration
  skip_final_snapshot = var.skip_final_snapshot
  deletion_protection = var.deletion_protection
  
  # Monitoring
  monitoring_interval = 0  # Disable enhanced monitoring to avoid role requirements
  
  tags = {
    Name = "${var.service_name}-mysql-db"
  }
}

# Get current AWS account ID
data "aws_caller_identity" "current" {}