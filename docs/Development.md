# Development references

## Container Registry

## AWS ECR

```bash
aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws/l2n7y5s7
docker build -t slashdevops/idp-scim-sync .
docker tag slashdevops/idp-scim-sync:latest public.ecr.aws/l2n7y5s7/slashdevops/idp-scim-sync:latest
docker push public.ecr.aws/l2n7y5s7/slashdevops/idp-scim-sync:latest
```

## Github Container Registry

```bash
echo $CR_PAT | docker login ghcr.io -u USERNAME --password-stdin
docker build -t ghcr.io/slashdevops/idp-scim-sync .
docker push ghcr.io/slashdevops/idp-scim-sync:latest
```
