# IAM roles for ECS tasks

# ECS Task Execution Role - Used by ECS to pull images and send logs
resource "aws_iam_role" "ecs_task_execution_role" {
  name = "${var.service_name}-ecs-task-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })

  tags = {
    ISBStudent = "true"
    Name       = "${var.service_name}-ecs-task-execution-role"
  }
}

# Attach AWS managed policy for ECS task execution
resource "aws_iam_role_policy_attachment" "ecs_task_execution_role_policy" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# ECS Task Role - Used by the application container
resource "aws_iam_role" "ecs_task_role" {
  name = "${var.service_name}-ecs-task-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })

  tags = {
    ISBStudent = "true"
    Name       = "${var.service_name}-ecs-task-role"
  }
}

# Add DynamoDB full access to task role when using DynamoDB
resource "aws_iam_role_policy_attachment" "ecs_task_dynamodb_policy" {
  count      = var.database_type == "dynamodb" ? 1 : 0
  role       = aws_iam_role.ecs_task_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess"
}
