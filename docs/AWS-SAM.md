# AWS SAM

This repository uses [AWS Serverless Application Model (AWS SAM)](https://docs.aws.amazon.com/en_en/serverless-application-model/latest/developerguide/what-is-sam.html) in two ways:

1. To deploy `idp-scim-sync` privately from source into your own AWS account.
2. To publish public versions to the [AWS Serverless Application Repository](https://serverlessrepo.aws.amazon.com/applications/us-east-1/889836709304/idp-scim-sync).

The same [template.yaml](../template.yaml) drives both workflows.

## How this repository uses SAM

- [template.yaml](../template.yaml) defines the Lambda function, EventBridge schedule, Secrets Manager secrets, S3 state bucket, KMS key, IAM role, and CloudWatch log group.
- The template includes `AWS::ServerlessRepo::Application` metadata so the application can be published to the AWS Serverless Application Repository.
- The Lambda build uses `Metadata: BuildMethod: makefile`, so `sam build` calls the `build-LambdaFunction` target from [Makefile](../Makefile).
- That build target compiles [cmd/idpscim/main.go](../cmd/idpscim/main.go) into a `bootstrap` binary for the `provided.al2023` runtime on `arm64`.
- Maintainer publishing is automated in [.github/workflows/aws-sam.yml](../.github/workflows/aws-sam.yml).

## Prerequisites

Before using SAM in this repository, install and configure the following:

- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)
- [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html)
- [Go](https://go.dev/doc/install)
- `make`
- AWS credentials with permission to create CloudFormation, Lambda, EventBridge, IAM, Secrets Manager, S3, CloudWatch Logs, and KMS resources
- Google Workspace service account credentials and delegated user email
- AWS IAM Identity Center SCIM endpoint and access token

Recommended environment variables:

```bash
export AWS_PROFILE=<profile-name>
export AWS_REGION=us-east-1
```

For publishing public versions, you also need:

```bash
export SAM_APP_BUCKET=<artifact-bucket-name>
export SAM_APP_VERSION=<semantic-version>
```

## Validate And Build

Use SAM and CloudFormation validation before deploying:

```bash
aws cloudformation validate-template \
  --template-body file://template.yaml \
  --profile "$AWS_PROFILE"

sam validate \
  --profile "$AWS_PROFILE" \
  --region "$AWS_REGION"

GIT_VERSION=dev sam build \
  --profile "$AWS_PROFILE" \
  --region "$AWS_REGION"
```

`sam build` creates `.aws-sam/build/LambdaFunction/bootstrap` by invoking the repository `Makefile`.

## Deploy From Source Into Your AWS Account

For private deployments from source, use `sam deploy --guided` the first time. It is the safest option because this template includes secret values such as the Google Workspace credential JSON and the SCIM token.

```bash
sam deploy --guided \
  --stack-name idp-scim-sync \
  --capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM \
  --profile "$AWS_PROFILE" \
  --region "$AWS_REGION"
```

During the guided prompts, provide at least these required values:

| Parameter | Purpose |
| --- | --- |
| `GWSServiceAccountFile` | Full contents of the Google Workspace service account JSON |
| `GWSUserEmail` | Delegated Google Workspace administrator or service user email |
| `SCIMEndpoint` | AWS IAM Identity Center SCIM endpoint URL |
| `SCIMAccessToken` | AWS IAM Identity Center SCIM bearer token |

Common optional parameters you will usually adjust:

| Parameter | Purpose | Default |
| --- | --- | --- |
| `GWSGroupsFilter` | Restrict which Google Workspace groups are synchronized | empty |
| `SyncUserFields` | Comma-separated optional user attributes to sync | empty = all supported fields |
| `ScheduleExpression` | EventBridge execution schedule | `rate(15 minutes)` |
| `BucketNamePrefix` | Prefix for the state bucket name | `idp-scim-sync-state` |
| `BucketKey` | S3 object key for the state file | `data/state.json` |
| `MemorySize` | Lambda memory | `256` |
| `Timeout` | Lambda timeout in seconds | `300` |
| `LogLevel` | Application log level | `info` |
| `LogFormat` | Lambda log format | `json` |

After the first guided deploy, SAM stores the deployment configuration. Subsequent updates are usually just:

```bash
GIT_VERSION=dev sam build \
  --profile "$AWS_PROFILE" \
  --region "$AWS_REGION"

sam deploy \
  --stack-name idp-scim-sync \
  --capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM \
  --profile "$AWS_PROFILE" \
  --region "$AWS_REGION"
```

## Update An Existing Deployment From The AWS Serverless Application Repository

This project is also published to the public AWS Serverless Application Repository:

- [AWS Serverless Application Repository application page](https://serverlessrepo.aws.amazon.com/applications/us-east-1/889836709304/idp-scim-sync)

The important AWS behavior to remember is this:

- When you deploy an application from the AWS Serverless Application Repository, AWS creates a CloudFormation stack whose actual name is prefixed with `serverlessrepo-`.
- When you update that application later, you must reuse the original application name or stack name that you entered when you first deployed it.
- Do not enter the generated `serverlessrepo-...` stack name during the update flow.

Example:

- Original name you entered: `idp-scim-sync`
- Actual CloudFormation stack created by AWS: `serverlessrepo-idp-scim-sync`
- Name to use for future updates: `idp-scim-sync`

### Console Update Flow

For most users, the console update flow is the clearest option:

1. Open the published application page in the AWS Serverless Application Repository.
2. Choose `Deploy` again for the same application.
3. Enter the same application name you used the first time, without the `serverlessrepo-` prefix.
4. Select the newer published application version.
5. Keep the existing parameter values or change them as needed.
6. Acknowledge IAM capabilities if prompted.
7. Deploy the update.

Use this flow when you want to move to a new published version without building from source locally.

### AWS CLI Update Flow

If you automate updates, the AWS Serverless Application Repository docs use a CloudFormation change set flow:

1. Inspect the application and required capabilities.
2. Create a CloudFormation change set for the same original stack name.
3. Execute that change set.

Example skeleton:

```bash
export AWS_REGION=us-east-1
export APPLICATION_ID=arn:aws:serverlessrepo:us-east-1:889836709304:applications/idp-scim-sync

aws serverlessrepo get-application \
  --application-id "$APPLICATION_ID" \
  --region "$AWS_REGION"

CHANGE_SET_ID=$(aws serverlessrepo create-cloud-formation-change-set \
  --application-id "$APPLICATION_ID" \
  --stack-name idp-scim-sync \
  --capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM \
  --region "$AWS_REGION" \
  --query ChangeSetId \
  --output text)

aws cloudformation execute-change-set \
  --change-set-name "$CHANGE_SET_ID" \
  --region "$AWS_REGION"
```

If your deployment requires explicit parameter overrides, include them when creating the change set. Reuse the same stack name you originally deployed with, not the generated `serverlessrepo-...` name.

## Publish A New Public Version

This section is only for project maintainers who publish new versions of the application.

Package the application artifacts:

```bash
sam package \
  --output-template-file packaged.yaml \
  --s3-bucket "$SAM_APP_BUCKET" \
  --profile "$AWS_PROFILE" \
  --region "$AWS_REGION"
```

Publish the packaged template as a new semantic version:

```bash
sam publish \
  --semantic-version "$SAM_APP_VERSION" \
  --template packaged.yaml \
  --region us-east-1 \
  --profile "$AWS_PROFILE"
```

Then make the published application deployable by other AWS accounts:

```bash
AWS_SAM_APP_ARN=$(aws serverlessrepo list-applications \
  --max-items 100 \
  --region us-east-1 \
  --profile "$AWS_PROFILE" \
  | jq -c '.Applications[] | select(.ApplicationId | contains("idp-scim-sync"))' \
  | jq -r '.ApplicationId')

aws serverlessrepo put-application-policy \
  --application-id "$AWS_SAM_APP_ARN" \
  --statements Principals='*',Actions='Deploy' \
  --region us-east-1 \
  --profile "$AWS_PROFILE"
```

The repository CI workflow in [.github/workflows/aws-sam.yml](../.github/workflows/aws-sam.yml) already automates this process for tagged releases.

## Remove A Published Application

If you need to delete the published application from the AWS Serverless Application Repository:

```bash
AWS_SAM_APP_ARN=$(aws serverlessrepo list-applications \
  --max-items 100 \
  --region us-east-1 \
  --profile "$AWS_PROFILE" \
  | jq -c '.Applications[] | select(.ApplicationId | contains("idp-scim-sync"))' \
  | jq -r '.ApplicationId')

aws serverlessrepo delete-application \
  --application-id "$AWS_SAM_APP_ARN" \
  --region us-east-1 \
  --profile "$AWS_PROFILE"
```
