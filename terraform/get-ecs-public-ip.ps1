$clusterName = terraform output -raw ecs_cluster_name
$serviceName = terraform output -raw ecs_service_name

$taskArn = aws ecs list-tasks `
    --cluster $clusterName `
    --service-name $serviceName `
    --query 'taskArns[0]' --output text

$networkInterfaceId = aws ecs describe-tasks `
    --cluster $clusterName `
    --tasks $taskArn `
    --query "tasks[0].attachments[0].details[?name=='networkInterfaceId'].value" `
    --output text

aws ec2 describe-network-interfaces `
    --network-interface-ids $networkInterfaceId `
    --query 'NetworkInterfaces[0].Association.PublicIp' `
    --output text