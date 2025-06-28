#!/bin/bash

# verify-lambda Lambda Deployment Script

set -e

echo "Building verify-lambda Lambda function..."

# Build for Linux/AMD64 (Lambda runtime)
GOOS=linux GOARCH=amd64 go build -o verify-lambda .

# Create deployment package
zip verify-lambda.zip verify-lambda

echo "✅ Build complete: verify-lambda.zip"
echo "📦 Ready for AWS Lambda deployment"

# Optionally deploy with AWS CLI (uncomment if needed)
# aws lambda update-function-code \
#   --function-name verify-lambda \
#   --zip-file fileb://verify-lambda.zip