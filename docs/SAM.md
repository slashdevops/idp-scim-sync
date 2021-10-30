# SAM

## Preparation

first export the environment variables

```bash
export AWS_PROFILE=<profile name here>
export AWS_REGION=<region here>
export AWS_IDPSCIM_BUNCKET_NAME=<bucket name here>
export IDPSCIM_VERSION=<version here>
```

## Deployment

```bash
aws cloudformation validate-template --template-body file://template.yaml 1>/dev/null --profile $AWS_PROFILE
sam validate --profile $AWS_PROFILE
sam build --base-dir cmd/idpscim/ --profile $AWS_PROFILE

sam package \
  --template-file template.yaml \
  --output-template-file packaged.yaml \
  --s3-bucket $AWS_IDPSCIM_BUNCKET_NAME \
  --profile $AWS_PROFILE

sam publish \
  --semantic-version $IDPSCIM_VERSION \
  --template packaged.yaml \
  --region $AWS_REGION \
  --profile $AWS_PROFILE
```
