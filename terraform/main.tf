# Wire together focused modules: network, ecr, logging, iam, ecs.

module "network" {
  source                = "./modules/network"
  service_name          = var.service_name
  container_port        = var.container_port
  cidr_blocks           = var.cidr_blocks
  alb_security_group_id = module.alb.alb_security_group_id
}

module "ecr" {
  source          = "./modules/ecr"
  repository_name = var.ecr_repository_name
}

module "logging" {
  source            = "./modules/logging"
  service_name      = var.service_name
  retention_in_days = var.log_retention_days
}

module "alb" {
  source         = "./modules/alb"
  vpc_id         = module.network.vpc_id
  public_subnets = module.network.public_subnet_ids
  service_name   = var.service_name
  cidr_blocks    = var.cidr_blocks
}

# IAM roles module
module "iam" {
  source        = "./modules/iam"
  service_name  = var.service_name
  database_type = var.database_type
}

# Conditionally create MySQL RDS instance
module "rds" {
  count                 = var.database_type == "mysql" ? 1 : 0
  source                = "./modules/rds"
  service_name          = var.service_name
  vpc_id                = module.network.vpc_id
  subnet_ids            = module.network.subnet_ids
  ecs_security_group_id = module.network.security_group_id
  db_name               = "ecommerce"
  db_username           = "admin"
}

# Conditionally create DynamoDB table
module "dynamodb" {
  count                         = var.database_type == "dynamodb" ? 1 : 0
  source                        = "./modules/dynamodb"
  service_name                  = var.service_name
  environment                   = "dev"
  enable_point_in_time_recovery = true
}

module "ecs" {
  source                    = "./modules/ecs"
  service_name              = var.service_name
  image                     = "${module.ecr.repository_url}:latest"
  container_port            = var.container_port
  subnet_ids                = module.network.subnet_ids
  security_group_ids        = [module.network.security_group_id]
  execution_role_arn        = module.iam.ecs_task_execution_role_arn
  task_role_arn             = module.iam.ecs_task_role_arn
  log_group_name            = module.logging.log_group_name
  ecs_count                 = var.ecs_count
  region                    = var.aws_region
  cpu                       = var.cpu
  memory                    = var.memory
  target_group_arn          = module.alb.target_group_arn
  enable_auto_scaling       = var.enable_auto_scaling
  environment_variables = var.database_type == "mysql" ? [
    {
      name  = "DATABASE_TYPE"
      value = "mysql"
    },
    {
      name  = "DB_HOST"
      value = module.rds[0].rds_endpoint
    },
    {
      name  = "DB_PORT"
      value = tostring(module.rds[0].rds_port)
    },
    {
      name  = "DB_USER"
      value = module.rds[0].db_username
    },
    {
      name  = "DB_PASSWORD"
      value = module.rds[0].db_password
    },
    {
      name  = "DB_NAME"
      value = module.rds[0].db_name
    }
  ] : [
    {
      name  = "DATABASE_TYPE"
      value = "dynamodb"
    },
    {
      name  = "DYNAMODB_TABLE_NAME"
      value = module.dynamodb[0].table_name
    },
    {
      name  = "AWS_REGION"
      value = var.aws_region
    }
  ]
}


// Build & push the Go app image into ECR
resource "docker_image" "app" {
  # Use the URL from the ecr module, and tag it "latest"
  name = "${module.ecr.repository_url}:latest"

  build {
    # relative path from terraform/ → src/
    context = "../src"
    # Dockerfile defaults to "Dockerfile" in that context
  }
}

resource "docker_registry_image" "app" {
  # this will push :latest → ECR
  name = docker_image.app.name
}
