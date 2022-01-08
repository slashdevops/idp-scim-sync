# AWS SAM

This document is a reference to the [AWS Serverless Application Model (SAM)](https://aws.amazon.com/serverless/sam/) and how to use it to deploy a `serverless` applications.

## Preparation

first export the environment variables

```bash
export AWS_PROFILE=<profile name here>
export AWS_REGION=<region here>

# only needed to Package and Deploy
export SAM_APP_BUCKET=<bucket name here>
export SAM_APP_VERSION=<version here>
```

## Deployment

Validate and Build

```bash
aws cloudformation validate-template --template-body file://template.yaml 1>/dev/null --profile $AWS_PROFILE
sam validate --profile $AWS_PROFILE

sam build --base-dir cmd/idpscim/ --profile $AWS_PROFILE
```

Deploy (Local way)

```bash
sam deploy --guided --profile $AWS_PROFILE
```

Package and Publish (Only project Owners)

```bash
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
