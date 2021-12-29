# Contributing

`idp-scim-sync` is [Apache 2.0 licensed](https://github.com/slashdevops/idp-scim-sync/blob/main/LICENSE) and
accepts contributions via GitHub pull requests. This document outlines
some of the conventions on to make it easier to get your contribution
accepted.

We gratefully welcome improvements to issues and documentation as well as to
code.

## Certificate of Origin

By contributing to this project you agree to the [Developer Certificate of
Origin (DCO)]([DCO](https://en.wikipedia.org/wiki/Developer_Certificate_of_Origin#:~:text=The%20Developer%20Certificate%20of%20Origin,Contributor%20License%20Agreement%20(CLA).)). This document was created by the Linux Kernel community and is a
simple statement that you, as a contributor, have the legal right to make the
contribution.

- All contributors to `idp-scim-sync` must sign off each contribution (usually this means each commit).The signature must contain your real name (sorry, no pseudonyms or anonymous contributions).  If your `user.name` and `user.email` are configured in your Git config, you can sign your commit automatically with `git commit -s`.
- Each commit sign off will be reviewed by the idp-scim-sync maintainer (by taking a look at their email address and Github profile) before merging the contribution.

References:

- [https://en.wikipedia.org/wiki/Developer_Certificate_of_Origin#:~:text=The%20Developer%20Certificate%20of%20Origin,Contributor%20License%20Agreement%20(CLA)](https://en.wikipedia.org/wiki/Developer_Certificate_of_Origin#:~:text=The%20Developer%20Certificate%20of%20Origin,Contributor%20License%20Agreement%20(CLA))
- [https://developercertificate.org/](https://developercertificate.org/)
- [https://probot.github.io/apps/dco/](https://probot.github.io/apps/dco/)
- [https://writing.kemitchell.com/2021/07/02/DCO-Not-CLA.html](https://writing.kemitchell.com/2021/07/02/DCO-Not-CLA.html)

## Understanding idp-scim-sync

If you are entirely new to idp-scim-sync, you might want to take a look at:

- [What is the AWS Single Sign-On SCIM implementation?](https://docs.aws.amazon.com/singlesignon/latest/developerguide/what-is-scim.html)
- [AWS Single Sign-On](https://aws.amazon.com/es/single-sign-on/)
- [Workspace Admin SDK -> Directory API](https://developers.google.com/admin-sdk/directory)

### Understanding the code

To get started with developing `idp-scim-sync`, you might want to check the following:

- [https://go.dev/](https://go.dev/)
- [https://github.com/spf13/cobra](https://github.com/spf13/cobra)
- [https://github.com/spf13/viper](https://github.com/spf13/viper)
- [AWS SSO SCIM -> Supported API operations](https://docs.aws.amazon.com/singlesignon/latest/developerguide/supported-apis.html)
- [Workspace Admin SDK -> Directory API -> Go quickstart](https://developers.google.com/admin-sdk/directory/v1/quickstart/go)
- [Code Reviews](https://github.com/golang/go/wiki/CodeReviewComments)
- [Testing Code that depends on google.golang.org/api](https://github.com/googleapis/google-api-go-client/blob/master/testing.md)
- [Unit Testing with the AWS SDK for Go V2](https://aws.github.io/aws-sdk-go-v2/docs/unit-testing/)
- [OpenSSF Best Practices Badge Program](https://bestpractices.coreinfrastructure.org/en)
- [CodeQL documentation](https://codeql.github.com/docs/)

### Mocks

These are generated using [gomock](https://github.com/golang/mock) project.

For better integration I use [go:generate](https://pkg.go.dev/cmd/go/internal/generate) to run `gomock` command inside files when mocks are needed

### Practices

- [Accept interfaces, return structs](https://bryanftan.medium.com/accept-interfaces-return-structs-in-go-d4cab29a301b)
- [CodeReviewComments#interfaces](https://github.com/golang/go/wiki/CodeReviewComments#interfaces)

## How to run the test suite

Prerequisites:

- make >= 3
- go >= 1.17

```bash
make test
```

Clean up the test output:

```bash
make clean
```

## Acceptance policy

These things will make a PR more likely to be accepted:

- a well-described requirement
- tests for new code
- tests for old code!
- new code and tests follow the conventions in old code and tests
- a good commit message (see below)
- all code must abide [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- names should abide [What's in a name](https://talks.golang.org/2014/names.slide#1)
- code must build on both Linux and Darwin, via plain `go build` or using `make build-dist`
- code should have appropriate test coverage and tests should be written to work with `go test ./...` or `make test`

### Format of the Commit Message

Try to follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) or [Semantic Commit Messages](https://gist.github.com/joshbuchea/6f47e86d2510bce28f8e7f42ae84c716)
