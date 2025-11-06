# DynamoDB table for shopping carts
resource "aws_dynamodb_table" "carts" {
  name         = "${var.service_name}-carts-${var.environment}"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "cart_id"

  attribute {
    name = "cart_id"
    type = "S"
  }

  attribute {
    name = "customer_id"
    type = "N"
  }

  # Global Secondary Index for querying by customer_id
  global_secondary_index {
    name            = "customer-index"
    hash_key        = "customer_id"
    projection_type = "ALL"
  }

  # Enable TTL on the ttl attribute
  ttl {
    attribute_name = "ttl"
    enabled        = true
  }

  # Point-in-time recovery
  point_in_time_recovery {
    enabled = var.enable_point_in_time_recovery
  }

  tags = {
    Name        = "${var.service_name}-carts-${var.environment}"
    Environment = var.environment
    Service     = var.service_name
  }
}
