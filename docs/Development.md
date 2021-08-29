# Development references

- [Testing Code that depends on google.golang.org/api](https://github.com/googleapis/google-api-go-client/blob/master/testing.md)
- [Semantic Commit Messages](https://gist.github.com/joshbuchea/6f47e86d2510bce28f8e7f42ae84c716)

## Semantic Commit Messages

See how a minor change to your commit message style can make you a better programmer.

Format: `<type>(<scope>): <subject>`

`<scope>` is optional

### Example

```text
feat: add hat wobble
^--^  ^------------^
|     |
|     +-> Summary in present tense.
|
+-------> Type: chore, docs, feat, fix, refactor, style, or test.
```

More Examples:

- `feat`: new feature for the user, not a new feature for build script
- `fix`: bug fix for the user, not a fix to a build script
- `docs`: changes to the documentation
- `style`: formatting, missing semi colons, etc; no production code change
- `refactor`: refactoring production code, eg. renaming a variable
- `perf`: A code change that improves performance
- `test`: adding missing tests, refactoring tests; no production code change
- `build`: Changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm)
- `ci`: Changes to our CI configuration files and scripts (example scopes: Travis, Circle, BrowserStack, SauceLabs)
- `chore`: updating grunt tasks etc; no production code change
- `license`: Edits regarding licensing; no production code change
- `revert`: Reverts a previous commit
- `bump`: Increase the version of something e.g. dependency
- `make`: Change the build process, or tooling, or infra
- `localize`: Translations update

References:

- [https://www.conventionalcommits.org/](https://www.conventionalcommits.org/)
- [https://seesparkbox.com/foundry/semantic_commit_messages](https://seesparkbox.com/foundry/semantic_commit_messages)
- [http://karma-runner.github.io/1.0/dev/git-commit-msg.html](http://karma-runner.github.io/1.0/dev/git-commit-msg.html)

## Development

### Mocks

SCIM Service Mocks

```bash
mockgen -package=mocks -destination=internal/mocks/scim_mocks.go -source=internal/core/scim.go
```

Identity Provider Service Mocks

```bash
mockgen -package=mocks -destination=internal/mocks/provider_mocks.go -source=internal/core/provider.go
```

Repository Mocks

```bash
mockgen -package=mocks -destination=internal/mocks/repository_mocks.go -source=internal/core/repository.go
```

Directory Service Mocks

```bash
mockgen -package=mocks -destination=internal/mocks/directory_mocks.go -source=internal/google/directory.go DirectoryService
```

## Practices

- [Accept interfaces, return structs](https://bryanftan.medium.com/accept-interfaces-return-structs-in-go-d4cab29a301b)
- [Always use interfaces](https://medium.com/@bryanftan/always-use-interfaces-in-go-d8f9f8f8f9c0)
