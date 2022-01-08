# AWS SAM

This document is a reference to the AWS Serverless Application Model (SAM) and how to use it to develop `serverless` applications.

## Preparation

first export the environment variables

```bash
export AWS_PROFILE=<profile name here>
export AWS_REGION=<region here>
export SAM_APP_BUCKET=<bucket name here>
export SAM_APP_VERSION=<version here>
```

## Deployment

Validate, Build and Publish

```bash
aws cloudformation validate-template --template-body file://template.yaml 1>/dev/null --profile $AWS_PROFILE
sam validate --profile $AWS_PROFILE

sam build --base-dir cmd/idpscim/ --profile $AWS_PROFILE

sam package \
  --output-template-file packaged.yaml \
  --s3-bucket $SAM_APP_BUCKET \
  --profile $AWS_PROFILE

sam publish \
  --semantic-version $SAM_APP_VERSION \
  --template packaged.yaml \
  --region $AWS_REGION \
  --profile $AWS_PROFILE
```
