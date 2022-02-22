#!/bin/bash
cd /home/ec2-user
aws s3 cp s3://somebucket-123321/app ./
chmod +x ./app
echo "export AWS_REGION=us-east-1" >> .bashrc