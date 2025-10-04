# Get ALB DNS name from Terraform outputs
$albDnsName = terraform output -raw load_balancer_dns_name

Write-Host "Getting ALB DNS name for API access..."
Write-Host "ALB DNS Name: $albDnsName"

if ([string]::IsNullOrEmpty($albDnsName) -or $albDnsName -eq "None") {
    Write-Host "No ALB DNS name found in Terraform outputs"
    exit 1
}

Write-Host "Application URL: http://${albDnsName}:8080"