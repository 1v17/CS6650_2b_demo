# Get cluster and service names from Terraform outputs
$clusterName = terraform output -raw ecs_cluster_name
$serviceName = terraform output -raw ecs_service_name

Write-Host "Getting public IP for ECS service..."
Write-Host "Cluster: $clusterName"
Write-Host "Service: $serviceName"

# Get the task ARN
$taskArn = aws ecs list-tasks --cluster $clusterName --service-name $serviceName --query 'taskArns[0]' --output text
# Write-Host "Task ARN: $taskArn"

if ($taskArn -eq "None" -or [string]::IsNullOrEmpty($taskArn)) {
    Write-Host "No running tasks found for service $serviceName"
    exit 1
}

# Get the network interface ID
$networkInterfaceId = aws ecs describe-tasks --cluster $clusterName --tasks $taskArn --query "tasks[0].attachments[0].details[?name=='networkInterfaceId'].value" --output text
# Write-Host "Network Interface ID: $networkInterfaceId"

if ($networkInterfaceId -eq "None" -or [string]::IsNullOrEmpty($networkInterfaceId)) {
    Write-Host "No network interface found for task"
    exit 1
}

# Get the public IP
$publicIp = aws ec2 describe-network-interfaces --network-interface-ids $networkInterfaceId --query 'NetworkInterfaces[0].Association.PublicIp' --output text
Write-Host "Public IP: $publicIp"

if ($publicIp -eq "None" -or [string]::IsNullOrEmpty($publicIp)) {
    Write-Host "No public IP assigned to network interface"
    exit 1
}

Write-Host "Application URL: http://${publicIp}:8080"