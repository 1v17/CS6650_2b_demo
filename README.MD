## Instructions

### Prepare Credentials

Retrieve you temporary credentials from Learner's Lab.
Enter your configuration when prompted:
```
aws configure
```

And set your session token:
```
aws configure set aws_session_token <YOUR-TEMP-SESSION-TOKEBN>
```

### Apply Infrastructure
```
cd terraform
terraform init
terraform apply -auto-approve
```

### Send Requests
Make sure you are in terraform folder

Get public ip address
```
aws ec2 describe-network-interfaces \
--network-interface-ids $(
    aws ecs describe-tasks \
    --cluster $(terraform output -raw ecs_cluster_name) \
    --tasks $(
        aws ecs list-tasks \
        --cluster $(terraform output -raw ecs_cluster_name) \
        --service-name $(terraform output -raw ecs_service_name) \
        --query 'taskArns[0]' --output text
    ) \
    --query "tasks[0].attachments[0].details[?name=='networkInterfaceId'].value" \
    --output text
) \
--query 'NetworkInterfaces[0].Association.PublicIp' \
--output text
```

Alternatively, use Powershell:

```pwsh
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
```
The above commands are written into a [Powershell script](/terraform/get-ecs-public-ip.ps1).

Send some requests:
```
curl http://<PUBLIC-IP-ADDRESS>:8080/albums
```

## Clean Up
```
terraform destroy -auto-approve
```