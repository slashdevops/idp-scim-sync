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

sam build --profile $AWS_PROFILE
```

Deploy (Local way)

```bash
sam deploy --guided  --capabilities CAPABILITY_IAM --capabilities CAPABILITY_NAMED_IAM --profile $AWS_PROFILE
```

## Publish

Package and Publish (Only project Owners)

```bash
sam package \
  --output-template-file packaged.yaml \
  --s3-bucket $SAM_APP_BUCKET \
  --profile $AWS_PROFILE
```

Buy default this is private the first time and depends on `AWS_REGION`

```bash
#export AWS_PUBLIC_REGIONS=($(aws ec2 describe-regions | jq -c -r '.Regions[] | .RegionName' | tr '\n' ' '))

export AWS_PUBLIC_REGIONS=(\
us-east-1 \
)

for AWS_PUBLIC_REGION in $AWS_PUBLIC_REGIONS; do
  echo "Publishing in $AWS_PUBLIC_REGION ..."
  sam publish \
    --semantic-version $SAM_APP_VERSION \
    --template packaged.yaml \
    --region $AWS_PUBLIC_REGION \
    --profile $AWS_PROFILE

  sleep 2

  export AWS_SAM_APP_ARN=$(\
  aws serverlessrepo list-applications \
    --max-items 100 \
    --region $AWS_PUBLIC_REGION \
    --profile $AWS_PROFILE \
    | jq -c '.Applications[] | select(.ApplicationId | contains("idp-scim-sync"))' \
    | jq -r '.ApplicationId' \
  )

  sleep 2

  aws serverlessrepo put-application-policy \
    --application-id $AWS_SAM_APP_ARN \
    --statements Principals='*',Actions='Deploy' \
    --region $AWS_PUBLIC_REGION \
    --profile $AWS_PROFILE

  sleep 1
done
```

Delete Application

```bash
export AWS_PUBLIC_REGIONS=(\
us-east-1 \
)

for AWS_PUBLIC_REGION in $AWS_PUBLIC_REGIONS; do
  echo "Deleting in $AWS_PUBLIC_REGION ..."

  export AWS_SAM_APP_ARN=$(\
  aws serverlessrepo list-applications \
    --max-items 100 \
    --region $AWS_PUBLIC_REGION \
    --profile $AWS_PROFILE \
    | jq -c '.Applications[] | select(.ApplicationId | contains("idp-scim-sync"))' \
    | jq -r '.ApplicationId' \
  )

  sleep 1

  aws serverlessrepo delete-application \
    --application-id $AWS_SAM_APP_ARN \
    --region $AWS_PUBLIC_REGION \
    --profile $AWS_PROFILE
done
```
