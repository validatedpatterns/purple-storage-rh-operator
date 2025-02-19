#!/bin/bash
set -e

command -v aws >/dev/null 2>&1 || { echo >&2 "I require aws cli but it's not installed.  Aborting."; exit 1; }

if [ $# -lt 1 ]; then
    echo "Usage: $0 sg-XXXXXXXXXXXXXXXXX"
    echo "Please pass the aws sg group id of the nodes you want the openings"
    exit 1
fi 
AWS_REGION="${AWS_REGION:-$(aws configure get region)}"

aws ec2 --region $AWS_REGION authorize-security-group-ingress --group-id $1 --protocol tcp --port 12345 --source-group $1 --no-cli-pager
aws ec2 --region $AWS_REGION authorize-security-group-ingress --group-id $1  --protocol tcp --port 1191 --source-group $1 --no-cli-pager
aws ec2 --region $AWS_REGION authorize-security-group-ingress --group-id $1  --protocol tcp --port 60000-61000 --source-group $1 --no-cli-pager
