# What's New

This document tracks notable changes, new features, and bug fixes across releases.

## Unreleased

### Go 1.27, dependency refresh, three crash fixes, and new architecture documentation (targeting `v0.46.0`, minor)

> [!IMPORTANT]
> This release fixes three crash paths that could **permanently stall the sync**, and changes the
> element ordering inside the state file. Expect one larger-than-usual set of updates on the first run
> after deploying, then a return to normal. No configuration change is required.

#### 🐛 Crash fixes (the important part)

Three code paths panicked on data that Google Workspace and AWS IAM Identity Center legitimately
return. Because the state file is only written after a *successful* run, a panic meant the next
scheduled invocation re-read the same record and panicked again — so a single malformed directory
entry could block all synchronization indefinitely, with no self-recovery.

* **A SCIM user with no email address crashed membership resolution.**
  `internal/scim.GetGroupsMembers` indexed `user.Emails[0]` without a guard. `buildUser` only keeps
  addresses flagged primary, so an AWS SCIM user with `"emails": []` — or with addresses but none
  marked primary — produced a nil slice and `index out of range [0] with length 0`. AWS does return
  that shape for users created without an email. Now resolved through a helper that prefers the
  primary address and falls back to the first, so no address is ever silently dropped.

* **A Google user without a family name crashed hashing.**
  `idp.buildUser` and `scim.buildUser` return `nil` to reject a record that AWS SCIM would refuse
  (missing name, given name, family name, primary email, or SCIM id), logging the reason as they do.
  Three callers stored that `nil` instead of honouring it, and it then reached `SetHashCode`, whose
  sort comparator dereferences every element. All three now skip and warn, matching the behaviour
  `GetGroups` already had for duplicate group names — one unusable record can no longer stop every
  other user from syncing.

  As defence in depth, `SetHashCode` now also drops `nil` elements from its hash input. A nil-safe
  comparator alone would not have been enough: `gob` panics outright on a nil pointer, so the failure
  would simply have moved into `MarshalBinary`.

* **A persistent 409 Conflict could exhaust the stack.**
  `pkg/aws.CreateOrGetUser` and `CreateOrGetGroup` cleared `externalId` and called themselves with no
  attempt limit. If AWS kept answering 409 while the follow-up lookup kept returning nothing, the
  recursion was unbounded. Both now retry exactly once — which is what their comments always
  described — and return the new `ErrConflictUnresolved` afterwards.

#### ⬆️ Go 1.27

The `go` directive moves from `1.26.5` to `1.27.0`. No workflow changes were needed: every CI job
already resolves the toolchain from `go.mod`.

Applied the Go 1.27 modernizers via `go fix ./...`: dropped a redundant loop-variable
re-declaration, and adopted the new **struct-literal field selectors** language feature in two
composite literals.

Every Go 1.27 behavioural change was checked against this codebase, and none affects it:

| Change | Assessment |
| --- | --- |
| `encoding/json` now backed by v2 | The v1 API still sorts map keys (verified empirically). The state file contains no map fields, and hashing uses `gob`, not JSON. |
| Timer channels always unbuffered | No `time.After`, `time.NewTimer`, `time.NewTicker` or `time.Tick` anywhere. |
| `http.Response.Body` drains on `Close` | Beneficial — improves connection reuse in `pkg/aws`. Nothing aborts a body early. |
| Stricter `//go:linkname` validation | No `linkname` and no `unsafe` import. |
| HTTP/2 server priorities | Server-side only; this is a client. |
| `SSL_CERT_FILE` / `SSL_CERT_DIR` on macOS | Affects local development only; noted in [Development.md](Development.md). |

Free wins with no code change: Go 1.27's size-specialized allocator cuts the cost of objects under
80 bytes by up to 30 %, which suits a codebase that allocates very large numbers of small `Member`,
`Group` and `Email` values.

The `omitzero` modernizer was **deliberately not applied**: `omitempty` and `omitzero` differ for
`false`, `0` and empty-but-non-nil slices, and `internal/model` uses `omitempty` on around 40 fields
that are serialized into the state file.

#### 📦 Dependencies

All 8 direct dependencies with updates available moved forward; no majors, no module-path changes.
`aws-lambda-go` v1.55.0, `aws-sdk-go-v2` v1.45.1 (with `config` v1.33.1, `credentials` v1.20.1, `s3`
v1.109.1, `secretsmanager` v1.46.1), `testify` v1.12.1, `google.golang.org/api` v0.295.0. Verified
with `go mod verify`, `govulncheck` (0 vulnerabilities called), regenerated mocks, and the full
`-race` suite.

#### ⚡ Concurrency and cancellation

Three fan-outs used the same hand-rolled `sync.WaitGroup` + channel-semaphore + buffered-error-channel
shape. It had no cancellation — one failure left every sibling running, spending upstream API quota on
results that were then discarded — and it silently dropped every error after the first.

`pkg/google.ListGroupMembersBatch`, `internal/idp.GetUsersByGroupsMembers` and
`internal/setup.Secrets` now use `errgroup.WithContext` with explicit limits, matching the idiom
`internal/scim` already used. The first failure cancels the rest.

`internal/setup.Secrets` additionally takes a `context.Context` instead of calling
`context.Background()` internally, so the four Secrets Manager reads can actually be abandoned when a
Lambda approaches its deadline rather than running until the runtime freezes the environment.

#### 🔀 State-file element ordering (one-off diff)

`GetUsersByGroupsMembers` assembled its result by ranging a map, so the order changed on every run.
The hash was unaffected — container hashes sort a copy before hashing — but the state file preserves
slice order, so **every** sync wrote a large spurious diff against the previous S3 object. The result
is now sorted by primary email.

Relatedly, the six `sort.Slice` calls in `internal/model` became `slices.SortStableFunc`. `sort.Slice`
is unstable, and four of these sort by a key that is not guaranteed unique (`HashCode`, `IPID`), which
left the order of tied elements — and therefore the hash of a set containing ties — unspecified.

Both changes ship together so there is only one reordering diff, not two.

A new [`golden_hash_test.go`](../internal/model/golden_hash_test.go) pins ten hash values against a
fully-populated fixture and round-trips it through both `gob` and JSON. **All ten are unchanged by
this release** — the hashing scheme itself did not move, only the order elements appear in the file.

#### 🔧 Correctness and style

* `err.Error() != "EOF"` string comparisons in the state-file decode path (8 sites) became
  `errors.Is(err, io.EOF)`. `gob` can return both `io.EOF` and `io.ErrUnexpectedEOF`; the string
  comparison caught only the first, so a genuinely truncated state file could be reported as a decode
  failure while a wrapped EOF was mistaken for one.
* Errors wrapped with `%v` instead of `%w` in `pkg/google` and `pkg/aws/secretsmanager` (6 sites) —
  which severed the chain that `internal/core` relies on via `errors.AsType`.
* **Log levels `fatal` and `panic` now work.** `config.Validate` has always accepted them — `logrus`
  leftovers predating the move to `log/slog`, which has no equivalent — but the logger handled
  neither, so they fell through to `info` while emitting a misleading "unknown log level" warning.
  They now map to `slog.LevelError`. The `--log-level` help text also advertised a `trace` level that
  was never valid.
* `buildCreateUserRequest` and `buildPutUserRequest` were ~150 lines of near-identical mapping
  differing only in `ID`, so adding a synced attribute meant editing both and forgetting one would
  make creates and updates disagree. Now one shared mapping.
* `cmd/idpscimcli` carried a 35-line copy of the logger setup; it now shares `setup.Logger`.
* 22 sentinel errors declared with `fmt.Errorf` and no format verbs became `errors.New`; naked
  returns removed from `internal/core/reconciling.go` and `internal/model/operations.go`; a parameter
  shadowing the `max` builtin renamed; `pkg/aws`'s path constants now actually used at all 18 call
  sites, via `path.Join` (which also means an id containing a slash can no longer forge extra path
  segments).

#### 🧪 CI and tooling

* **Codecov has been receiving nothing.** `build.yml` uploaded `./coverage.out`, but the Makefile
  writes coverage to `build/coverage.txt`. That path never existed in CI, so the badge has been
  reporting stale data. Fixed and verified from a clean tree.
* Added a committed [`.golangci.yml`](../.golangci.yml) and a **Lint job** to CI. The contributor
  checklist already required `golangci-lint`, but with no config each runner used its own defaults, so
  results were not reproducible. **The lint baseline is now 0 issues**, down from 18, against a
  stricter configuration.
* Closed 21 missing doc comments on exported symbols, as the Google style guide requires.
* `internal/setup` had no tests at all despite wiring everything together; it now has four.
* Fixed two tests that called `os.Remove(os.TempDir())` — attempting to delete the shared system temp
  directory, failing silently only because `/tmp` is never empty — and three subtests in
  `pkg/aws/config_test.go` that passed only because `os.Setenv` leaked credentials between them. All
  now use `t.Setenv` / `t.TempDir()` and pass individually as well as in file order.

#### 📚 Documentation

Three new guides, with 10 mermaid diagrams so architecture is reviewable in a diff rather than as
binary exports:

* **[Architecture.md](Architecture.md)** — system context, package topology, the sync algorithm, the
  reconciliation model, the hashing scheme, state-file lifecycle, how SCIM membership is inferred,
  concurrency limits, and failure modes.
* **[Implementation-Guide.md](Implementation-Guide.md)** — the invariants that are easy to break, a
  step-by-step for adding a synced user attribute, testing conventions, and state-file compatibility
  rules.
* **[User-Manual.md](User-Manual.md)** — end-to-end deployment and operations, with troubleshooting.

Also: every package comment expanded (`internal/core`'s described `internal/deepcopy`);
[State-File-example.md](State-File-example.md) now explains what is hashed and why, and that deleting
the state file is safe; [Development.md](Development.md) updated for Go 1.27 and the lint gate;
[Configuration.md](Configuration.md) documents the log levels; and the README gains a Quick Start, an
architecture diagram, a Contributing section, and links to all 14 docs instead of 5.

#### 🔎 Panics now name the record that caused them

Go 1.27 prints `runtime/pprof` goroutine labels in traceback headers. Every bounded fan-out now tags
its goroutines, and so does the sync goroutine itself, so a panic in CloudWatch reads:

```
goroutine 42 [running] {sync: "users", user: "zoe@example.com"}:
```

instead of leaving you to bisect the directory. The crash fixes above were exactly this shape — an
`index out of range` inside a worker, with no indication of which record.

Labels: `sync=root` + `codeVersion` on the sync itself, `sync=users` + `user` on the user fan-out,
`sync=group-members` + `group` on the Google members fan-out, and `sync=scim-group-members` + `user`
on the SCIM membership inversion.

Purely additive instrumentation: no control flow, no API, no hashing, no state-file change. Cost is
~29 ns and 104 B per goroutine, which is noise against a network round-trip. Opt out with
`GODEBUG=tracebacklabels=0`.

> [!NOTE]
> Implemented with `pprof.SetGoroutineLabels`, **not** `pprof.Do`. `pprof.Do` defers restoring the
> previous label set, and on a panic those defers run during unwinding before the runtime prints the
> traceback, so the labels are gone exactly when they are needed. `runtime.Stack` shows them either
> way, so the broken version verifies fine by accident. A test guards the distinction.

#### 🧮 Canonical hash ordering (one-off membership reconciliation)

> [!IMPORTANT]
> This changes the `groupsMembers` container hash, which the sync **reads**. The first run after
> deploying performs one membership reconciliation pass; subsequent runs are stable. Group and user
> hashes are unchanged, so groups and users do **not** re-sync.

Every hash in `internal/model` is a `gob` encoding, and the hand-written `MarshalBinary` methods walk
their `Resources` slices in slice order. Arrival order is therefore part of the hash unless it is
normalised first — and it must be, because the identity-provider fan-outs and the SCIM membership
inversion all build their slices from maps, whose iteration order is randomised.

The normalisation was inlined per method, and two places had drifted:

* **`GroupsMembersResult` sorted its outer slice but not the members nested inside each entry.**
  `GroupMembers.MarshalBinary` walks its own members in slice order, so member order *within* a group
  reached the bytes of the enclosing hash — while `GroupMembers`' own hash stayed stable, because that
  method does sort its copy. Since `internal/core` compares `GroupsMembersResult` hashes to decide
  whether membership needs reconciling, an upstream reordering caused a spurious **full membership
  sync**. Two sources present members in an unstable order: the Google Directory API guarantees no
  member ordering, and `internal/scim`'s membership inversion appends members in goroutine completion
  order.

* **`State.SetHashCode` sorted nothing at all.** It rebuilt each `*Result` through a builder and then
  gob-encoded the result; the builder's own sorting only ever applied to its private copy, so it never
  reached the State hash bytes. Two runs over identical data could produce different
  `State.HashCode` values. Nothing reads that field — `internal/core` compares only the three
  container hashes — so this had no runtime effect, but it made the state file's top-level `hashCode`
  useless as an external change indicator.

Both are fixed by giving the hash input a canonical ordering at every level. The five orderings now
live in one place (`internal/model/sort.go`) and are called by every `SetHashCode`, so they cannot
drift apart again.

Care was taken that computing a hash never reorders the caller's live data: copying a `GroupMembers`
struct shares its `Resources` backing array, so each entry gets its own copied member slice before
sorting. A test pins this.

Golden hashes: `groupsMembers` and `state` updated deliberately; `groups`, `users` and every
individual resource hash unchanged.

#### 🧷 Also fixed

* **`idpscimcli aws users list --filter` never worked.** The flag was registered on the `aws users`
  grouping command, which has no run function, so Cobra advertised it under
  `idpscimcli aws users --help` where nothing consumed it and rejected it as
  `unknown flag: --filter` where it was actually needed. `aws groups list` had it right. Moved to the
  `list` subcommand, with a test asserting both leaf commands accept it and neither grouping command
  advertises it.

* **`.gitignore` was hiding the `cmd/` source trees.** The entries `idpscim` and `idpscimcli` were
  intended for stray binaries built into the repository root, but a gitignore pattern without a slash
  matches at *any* depth — so `cmd/idpscim/` and `cmd/idpscimcli/` were both ignored. Existing tracked
  files were unaffected, which is why nobody noticed, but any **new** file added under either
  directory was invisible to `git status` and silently omitted from commits, and every
  `.gitignore`-aware tool (ripgrep, editors, Docker build contexts) skipped both trees entirely. Both
  patterns are now anchored as `/idpscim` and `/idpscimcli`, verified to still ignore root binaries
  while making `cmd/` visible again.

* **Corrected a wrong note about the `unit` build tag.** An earlier commit in this branch claimed
  `-tags=unit` selected nothing. It does: `cmd/idpscimcli/cmd/*_test.go` carry `//go:build unit`, so a
  bare `go test ./...` reports "no test files" for that package while `make test` runs it. The
  Makefile comment now says so. (The original mistake came from a `.gitignore`-aware grep that could
  not see into `cmd/` — the same root cause as the entry above.)

#### 🔜 Known follow-up

`golang.org/x/oauth2/google.CredentialsFromJSONWithParams` is deprecated for a security reason — it
does not validate the credential configuration, and here the service-account JSON arrives from
Secrets Manager unvalidated. It carries a documented `TODO` and is deliberately deferred to its own
change: it is the Lambda's authentication path and needs separate tests plus a manual run against a
real Google Workspace tenant.

### Template housekeeping: remove dead/misleading IAM grants and standardize on `${AWS::Partition}`

Cleans up the Lambda execution role, the KMS key policy, and a few hardcoded partition strings in `template.yaml`. No runtime behavior change — every removal is a permission or grant that was never reached or never matched at runtime.

**Lambda role (`CustomLambdaPolicy`):**

* **Removed `s3:ListBucket` statement** (`S3ListBucketPolicy`). The runtime code only calls `s3:GetObject` and `s3:PutObject` on the single state object ([`internal/repository/s3.go`](../internal/repository/s3.go)) — `ListBucket` / `ListObjectsV2` / `HeadObject` are never issued, so the statement and its `s3:prefix` condition were unreachable.
* **Removed `kms:GenerateDataKeyPair` statement** (`KMSGetDataPolicy`). `GenerateDataKeyPair` is an asymmetric-key API; the state-bucket CMK is the default symmetric KMS key, so this action could never be called against it. The remaining symmetric action `kms:GenerateDataKey` (needed by S3 for SSE-KMS with bucket keys) is already granted by the renamed `KMSStateObjectPolicy` statement.
* **Removed `secretsmanager:GetResourcePolicy` action.** The code only calls `GetSecretValue` ([`pkg/aws/secretsmanager.go`](../pkg/aws/secretsmanager.go)). `GetResourcePolicy` reads the secret's resource-policy JSON — nothing in this Lambda needs it. Renamed the surviving statement to `SecretsManagerGetSecretValuePolicy` (the old name `SSMGetParameterPolicy` was a copy-paste from a different service).

**KMS key policy:**

* **Removed `AllowAWSLambdaToRetrieveKMSKey` statement.** Its principal was `Service: lambda.amazonaws.com`, but at runtime KMS sees the Lambda's **assumed-role** ARN (not the Lambda service) when S3 forwards the `kms:Decrypt` / `kms:GenerateDataKey` calls. The statement therefore granted nothing at runtime — the real grant comes from `AllowIAMThisAccount` (the standard "delegate to IAM" pattern), combined with the role's identity-based policy. Removing the dead statement makes the grant model unambiguous.

**Partition portability:**

* Replaced four hardcoded `arn:aws:…` ARNs with `arn:${AWS::Partition}:…`:
  * `KMSKey` key-policy principal (`AllowIAMThisAccount`)
  * `Bucket` `BucketEncryption.KMSMasterKeyID`
  * `BucketPolicy` `AllowAWSLambdaFunction` principal
* The rest of the template already used `${AWS::Partition}`; these were the last holdouts. The template is now deployable in non-commercial AWS partitions (GovCloud, China) without manual edits.

### CI fix: cosign now signs the published container manifest by tag (closes the v0.45.0 signing failure)

Fixes the `Cosign sign published container manifest (keyless / Sigstore)` step of the release workflow, which failed with **`MANIFEST_UNKNOWN: manifest unknown`** on every release attempt after the multi-arch build was restored (see ["CI fix: restore multi-arch container builds"](#ci-fix-restore-multi-arch-container-builds-in-the-release-workflow)).

**Root cause.** The step resolved the digest to sign by piping `podman manifest inspect` into `jq -r '.digest // .manifests[0].digest'`. Two compounding problems:

1. **A manifest list's own JSON has no top-level `.digest`** (its digest is computed by hashing the JSON, not stored inside it). So the `//` fallback always wins and returns `.manifests[0].digest` — the digest of the **first per-arch image** (arm64), not the manifest list.
2. **Podman re-serializes manifests when pushing** (media-type conversion between Docker `vnd.docker.distribution.manifest.v2+json` and OCI `vnd.oci.image.manifest.v1+json`). The locally computed digest therefore does not match what GHCR stores, so cosign's lookup of `ghcr.io/…@sha256:<local-digest>` returns 404.

Result: cosign was asked to sign a digest that exists nowhere on the registry.

**Fix.** Sign by tag (`cosign sign --recursive ghcr.io/…:TAG`). Cosign internally HEAD-resolves the tag to its authoritative on-registry digest and signs that digest — the signature is still stored *by digest*, so the resulting artifact is identical to what the previous (broken) code intended to produce. The classic "signing-by-tag races with concurrent pushes" caveat does not apply here: this job exclusively owns the `v<x.y.z>` and `latest` tags and has just pushed them sequentially in the previous step.

No code or release-artifact changes.

### CI fix: restore multi-arch container builds in the release workflow

Fixes the `Publish Container Images` job (failing since v0.44.1, surfaced again on the v0.45.0 release as ["Could not resolve digest for ghcr.io/slashdevops/idp-scim-sync:v0.45.0"](https://github.com/slashdevops/idp-scim-sync/actions/runs/26356807875/job/77585211704)).

**Root cause.** When the workflow was migrated from Docker to Podman in `ab22744`, the `docker/setup-qemu-action` step was removed on the assumption that pre-built binaries no longer required cross-compilation. They don't — but every `RUN` line in `Containerfile` (notably `apk add ca-certificates`) still has to execute under the target architecture. On the `amd64` runner, building the `arm64` variant therefore needs QEMU user-mode emulation registered with `binfmt_misc`. Without it, `podman build --platform linux/arm64` died with `exec /bin/sh: Exec format error`, the `arm64` image was never created, `podman manifest add` then failed, no manifest was pushed, and `cosign sign` finally failed with "Could not resolve digest". The two upstream `make` failures were not surfaced as red steps because `make ... | tee $GITHUB_STEP_SUMMARY` runs under bash without `pipefail`, so `tee`'s exit 0 masked the make exit 1.

**Changes in `.github/workflows/container-image.yml`:**

* Install `qemu-user-static` and `binfmt-support` alongside `podman` so the kernel can run arm64 binaries during the build.
* Set `defaults.run.shell: bash -eo pipefail {0}` at the job level so future `cmd | tee` failures fail the step instead of being silently swallowed.

No code changes; release artifacts and signing flow are unchanged.

### SCIM members sync — major security & performance improvement (closes [#520](https://github.com/slashdevops/idp-scim-sync/issues/520))

> [!IMPORTANT]
> This release replaces the brute-force algorithm that reconstructed group memberships on the AWS IAM Identity Center side. The change is **internal-only** (no CLI/config change) but materially improves both the **security posture** and the **runtime cost** of every sync.

#### Security improvements

* **Drastically smaller attack & failure surface per sync.** The number of authenticated SCIM requests issued per sync drops by ~2 orders of magnitude (see *Performance* below). Each request is a credential-bearing call to AWS IAM Identity Center — fewer requests means fewer opportunities for a credential leak, log capture, in-flight tampering, or partial-failure half-state to be observed.
* **Shorter sync window = smaller inconsistency window.** Previously, a sync of a few hundred users could run for many minutes while ~100k requests trickled out behind a hand-rolled 10–150 ms random sleep. During that window, group membership in AWS could be partially reconciled — an externally-observable inconsistent state. The new path completes in a fraction of the time, shrinking that window proportionally.
* **No more time-based throttling band-aid.** The previous `time.Sleep(rand.Intn(...))` jitter existed solely to avoid tripping AWS SCIM throttles under the brute-force call volume. It has been removed: the new call profile is light enough that artificial gapping is no longer required. This eliminates a source of non-determinism in the sync path and removes timing-dependent behavior from a security-sensitive code path.
* **Deterministic pagination.** Cursor-based pagination (`?cursor` + `nextCursor`) walks the full result set deterministically, so memberships can no longer be silently truncated by hitting an undocumented page cap mid-sync.

#### Performance improvements

* **Before:** `internal/scim.GetGroupsMembersBruteForce` issued one `ListGroups` call for *every* (group, user) combination — `O(N_groups × N_users)` requests per sync, throttled by a 10–150 ms random sleep and capped at concurrency 5. For an org with 200 groups and 500 users this is **~100,000 calls per sync run**.
* **Now:** `internal/scim.GetGroupsMembers(ctx, groups, users)` issues one cursor-paginated `?cursor&filter=members.value eq "<user-id>"` request per user (plus one extra request per additional page of memberships, when a single user belongs to more than 100 groups). The result is then inverted into the group → members map the rest of the pipeline expects. For the same 200-group / 500-user org this is **~500 calls per sync — roughly two orders of magnitude fewer requests**.
* **Lower Lambda execution time and cost.** Fewer requests + no random sleeps directly reduces billable Lambda duration on every scheduled invocation.

This is enabled by two AWS IAM Identity Center SCIM features documented at <https://docs.aws.amazon.com/singlesignon/latest/developerguide/listgroups.html>:

* The `members.value eq "<user-id>"` filter, which returns every group containing a given user.
* Cursor-based pagination (`?cursor` + `nextCursor`), which lifts the historical 50-result page cap to 100 results per page and supports walking the full result set deterministically.

**API changes (internal-only — no user-facing CLI/config change):**

* `pkg/aws.SCIMService` gained `ListGroupsWithCursor(ctx, filter, cursor) (*ListGroupsResponse, error)`. `ListGroups` is unchanged and remains the non-paginated single-page call.
* `pkg/aws.ListResponse` gained a `NextCursor string` field.
* `internal/core.SCIMService.GetGroupsMembers` now takes `(ctx, *model.GroupsResult, *model.UsersResult)`. The previous `GetGroupsMembers(ctx, gr)` and `GetGroupsMembersBruteForce(ctx, gr, ur)` methods, plus their AWS-side brute-force scaffolding, have been removed entirely — there is no compatibility shim.
* Memberships pointing at AWS-side groups that are *not* part of the in-scope `gr` (for example AWS-managed groups created outside the sync) are silently ignored, matching prior behavior.

**Tests:** the concurrency cap is now exercised under `testing/synctest` (graduated to the standard library in Go 1.26 — see <https://go.dev/blog/testing-time>) so the test asserts the true peak in-flight count using virtual time, instead of waiting on a wall-clock sleep race.

### IAM least-privilege hardening for the state-file Lambda role

Tightens the Lambda execution role in `template.yaml` so it can only touch the single state object via the single intended path. No behavior change for normal operation; the role is now strictly scoped.

**S3 statement (split in two):**

* `S3ObjectPolicy` — `s3:GetObject*` / `s3:PutObject*` are now scoped to the exact state object (`arn:${Partition}:s3:::<bucket>/<BucketKey>`) instead of `<bucket>/*`.
* `S3ListBucketPolicy` — `s3:ListBucket` stays on the bucket ARN (required by S3) but is now gated by a `s3:prefix` condition matching `BucketKey`.

**KMS statements (`KMSGetDataPolicy` and `KMSDecryptPolicy`):**

Both now carry the AWS-recommended SSE-KMS scoping conditions:

* `kms:ViaService = s3.<region>.amazonaws.com` — the CMK can only be used through requests S3 forwards on the role's behalf, not via direct `kms:Decrypt` / `kms:Encrypt` calls.
* `kms:EncryptionContext:aws:s3:arn = arn:${Partition}:s3:::<bucket>` — the KMS grant only applies when S3 passes that exact **bucket** ARN as encryption context. The bucket ARN (not the object ARN) is the correct value here because the state bucket has `BucketKeyEnabled: true`: S3 generates a short-lived bucket-level data key once and derives per-object keys locally from it, so the encryption context S3 forwards to KMS is the bucket ARN. The role can still only reach the single state object through the `S3ObjectPolicy` ARN scoping. See [AWS docs — S3 Bucket Keys](https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucket-key.html).

These are belt-and-braces additions on top of the existing bucket policy, public-access block, `aws:SecureTransport` deny, and SSE-KMS-with-bucket-key configuration.

References: [AWS docs — kms:ViaService](https://docs.aws.amazon.com/kms/latest/developerguide/policy-conditions.html#conditions-kms-via-service), [AWS docs — SSE-KMS encryption context](https://docs.aws.amazon.com/AmazonS3/latest/userguide/specifying-kms-encryption.html), [AWS docs — S3 Bucket Keys](https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucket-key.html).

### OpenSSF Scorecard Hardening (Phase 4) — Fuzzing

Closes the **Fuzzing** Scorecard check by adding native Go fuzz targets and a CI workflow that exercises them.

**Targets added:**

* `internal/repository.FuzzDiskRepositoryGetState` — fuzzes JSON state-file loading, the primary untrusted-input surface in the project (state files can be persisted to disk or S3 and re-read).
* `internal/model.FuzzGroupsResultUnmarshalBinary` — fuzzes gob deserialization, which loops `Items` times calling `dec.Decode`. An attacker who controls the encoded blob can request an enormous number of items; the fuzzer ensures this fails safely rather than panicking or hanging.

**Tooling:**

* New `make fuzz` target enumerates every `Fuzz*` function via `go test -list` and runs each one for `FUZZ_TIME` (default `60s`, override per-invocation).
* New `.github/workflows/fuzz.yml`:
  * Runs every PR with a **60s budget per target** (fast enough to not block merges).
  * Runs weekly on Wednesdays at 06:00 UTC with a **10-minute budget per target** to actually find bugs.
  * `workflow_dispatch` lets a maintainer pick any duration on demand.

When a fuzz target finds a crashing input, Go saves it under `testdata/fuzz/<FuzzName>/<hash>`. Commit those files: they become permanent regression tests that run as part of `go test ./...`.

### OpenSSF Scorecard Hardening (Phase 3) — Signed releases + SLSA provenance

Closes the **Signed-Releases** Scorecard check (0/10 → 10/10) by adopting two complementary supply-chain primitives:

* **SLSA Level 3 provenance for binary release zips.** A new `provenance` job in `.github/workflows/release.yml` calls the official `slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml@v2.1.0` reusable workflow. It binds every release zip's `sha256` to a Sigstore-signed in-toto attestation (`multiple.intoto.jsonl`) and uploads it to the GitHub release. Anyone can verify the link between artifact ↔ source ↔ build with [`slsa-verifier`](https://github.com/slsa-framework/slsa-verifier).
* **Cosign keyless signing for container images.** After `podman manifest push`, `.github/workflows/container-image.yml` resolves the multi-arch manifest to its content digest and signs by digest with `cosign sign --recursive` (keyless via Sigstore Fulcio + Rekor). Both the tagged image and `latest` are signed.

Also added a `Verifying release artifacts` section to `SECURITY.md` with copy-pasteable verification commands. Fixed a stale `codeql-analysis.yml` reference in the same file.

> [!NOTE]
> The `slsa-github-generator` reusable workflow is intentionally pinned by tag (not by SHA). The SLSA verifier validates the workflow ref against its own allow-list of signed releases — SHA-pinning would break verification. This is the only documented exception to the project's SHA-pinning rule.

### OpenSSF Scorecard Hardening (Phase 2) — Vulnerability remediation

Closed the **Vulnerabilities** Scorecard check (0/10 → 10/10). The previous baseline reported 19 known vulnerabilities, all tracing to two indirect dependencies:

* `golang.org/x/crypto` v0.51.0 → v0.52.0 (closes 12 advisories)
* `golang.org/x/net` v0.54.0 → v0.55.0 (closes 7 advisories, including `GO-2026-5026` — the only one actually reachable from project code via `idna.ToASCII` in `pkg/aws/scim.go`)

Also added a **`govulncheck` step to the `Build` workflow** so future regressions block CI before merge. The step uses `golang/govulncheck-action@v1.0.4` (SHA-pinned).

### OpenSSF Scorecard Hardening (Phase 1)

Started hardening the project against the [OpenSSF Scorecard](https://github.com/ossf/scorecard) checks. This first phase targets the `Token-Permissions` and `Pinned-Dependencies` checks and adds continuous Scorecard analysis to CI. No behavior change — CI security posture only.

**What changed:**

* **Workflow tokens scoped to minimum.** All `.github/workflows/*.yml` now declare `permissions: contents: read` at the top level. Elevated permissions (`packages: write`, `id-token: write`, `contents: write`) are granted only to the specific jobs that need them — release publishing, container image push, AWS SAM OIDC auth.
* **Every GitHub Action pinned by full commit SHA**, with the human-readable tag preserved as a trailing comment. Covers `actions/checkout`, `actions/setup-go`, `actions/upload-artifact`, `actions/download-artifact`, `actions/setup-python`, `codecov/codecov-action`, `softprops/action-gh-release`, `github/codeql-action/*`, `aws-actions/setup-sam`, and `aws-actions/configure-aws-credentials`.
* **Container base image pinned by digest.** `Containerfile` now uses `alpine:3.21@sha256:…` instead of the floating `alpine` tag.
* **Dependabot now watches GitHub Actions and Docker** in addition to Go modules, so the SHA pins above stay current.
* **New `OpenSSF Scorecard` workflow** runs on push to `main` and weekly, publishing SARIF results to GitHub code scanning.
* **Dropped unused `security-events: write` permission** from the build workflow.

### Documentation Refresh

Refreshed the project documentation to align it with the current Go, AWS SAM, and AWS Serverless Application Repository workflows.

This update expands the README documentation map, clarifies how to deploy and update the serverless application, documents the two binaries built from `cmd/`, and modernizes the remaining docs in the `docs/` folder so they are linked and project-valid.

## v0.40.1

### Improved HTTP Retry Library

Replaced the `httpretrier` library with [httpx](https://github.com/slashdevops/httpx), a zero-dependency HTTP client with built-in retry support.

**Why:** The previous library did not properly handle HTTP `429 Too Many Requests` responses, which caused issues with AWS SSO SCIM API throttling under high load.

**What changed:**

* The `httpx` library automatically retries on `429` and `5xx` responses with configurable backoff strategies.
* AWS SCIM API calls now use **jitter backoff** instead of simple exponential backoff, reducing the chance of thundering herd effects during rate limiting.
* Google Workspace API calls use **exponential backoff** for reliable retries.
* The `httpx` library has zero external dependencies and integrates with Go's `slog` logging.

### AWS SCIM Client Improvements (`pkg/aws`)

Several code quality improvements and bug fixes in the AWS SCIM client:

* **Bug fix:** `CreateOrGetUser` used `reflect.DeepEqual` to compare a `*CreateUserRequest` with a `*GetUserResponse` — different types, so the comparison always returned `false`, causing unnecessary PUT updates on every 409 conflict. Replaced with a typed `usersEqual` function that compares only sync-relevant attributes.
* **Removed `pkg/errors` dependency:** Replaced with stdlib `errors` and `fmt` packages. Sentinel errors now use `errors.New` instead of `errors.Errorf`.
* **Go 1.26 `errors.AsType`:** Migrated all `errors.As` calls to the generic `errors.AsType[T]` for compile-time type safety and better performance.
* **Fixed `String()` methods:** `User.String()` and `Group.String()` no longer call `os.Exit(1)` on marshal failure. They return a safe fallback string instead.
* **Eliminated double JSON decode:** `GetUserByUserName` and `GetGroupByDisplayName` no longer marshal a resource to JSON and re-parse it. They use direct type conversion instead.
* **Fixed decode error fallback:** `CreateGroup` and `CreateOrGetGroup` no longer attempt to read an already-consumed response body on decode failure.
* **Removed redundant context set:** `do()` no longer calls `req.WithContext(ctx)` since the request is already created with `http.NewRequestWithContext`.
* **Simplified type conversions:** `CreateOrGetUser` and `CreateOrGetGroup` use type conversions instead of manual field-by-field struct copies.

### Go 1.26 Modernization

Applied Go 1.26 best practices across the codebase:

* **Removed `github.com/pkg/errors` dependency:** Replaced all `errors.Wrap` and `errors.Errorf` with stdlib `fmt.Errorf` (with `%w`) and `errors.New` in `internal/setup`, `internal/repository`, `pkg/aws`, and `cmd/idpscim`.
* **`errors.AsType[T]`:** Migrated `errors.As` calls to the generic `errors.AsType[T]` in `internal/core/sync.go` for type safety and performance.
* **Fixed `os.Exit` in `Hash()`:** `internal/model.Hash()` no longer calls `os.Exit(1)` on nil input or encoding failure. It panics instead (appropriate for programming errors, recoverable, produces stack trace).

### CLI Modernization (`idpscimcli`)

Improved the `idpscimcli` command-line tool for code quality and testability:

* **Fixed logger initialization bug:** Log handler options were applied after the handler was created, so log level and format from config had no effect. Fixed by configuring options before creating the handler.
* **Testable error handling:** `getGWSDirectoryService` now returns errors instead of calling `os.Exit(1)`, allowing proper error propagation and unit testing.
* **`show()` error handling:** The output function now returns errors instead of silently discarding marshal failures.
* **Reduced code duplication:** Extracted `newSCIMHTTPClient()`, `newGWSHTTPClient()`, and `newAWSSCIMService()` helpers, eliminating 6 duplicated builder blocks.
* **Added unit tests:** New tests for `show()` covering JSON, YAML, default format, marshal errors, and empty structs.
* **Fixed typos:** "usrs" → "users", "Servive" → "Service" in command help text.

## v0.44.0

### Configurable User Fields

You can now choose which optional user attributes are synced from Google Workspace to AWS SSO SCIM using the new `sync_user_fields` configuration option.

For example, sync only phone numbers and enterprise data while excluding addresses, locale, or timezone. When not configured, all fields are synced as before (fully backward compatible).

**Available fields:** `phoneNumbers`, `addresses`, `title`, `preferredLanguage`, `locale`, `timezone`, `nickName`, `profileURL`, `userType`, `enterpriseData`.

See [Configurable User Fields](../README.md#configurable-user-fields) for configuration examples and behavior notes.

### Bug Fix: Unnecessary member re-syncs

Fixed a bug where group members were re-synced on every Lambda execution even when nothing changed in Google Workspace.

**Root cause:** `MergeGroupsMembersResult` was not consolidating entries for the same group when merging "created" and "equal" member sets. This produced duplicate group entries in the state file, causing the groups-members hash to never match the IDP data on subsequent syncs.

**Impact:** After upgrading, the first sync will reconcile the state file automatically. Subsequent syncs will correctly skip member reconciliation when no changes are detected.

### Performance Improvements

* **Concurrent user fetching:** `GetUsersByGroupsMembers` now fetches user details from the Google Workspace API concurrently (up to 10 parallel requests) instead of sequentially. For deployments with 100+ users, this reduces the user retrieval phase from minutes to seconds.

* **Optimized member comparison:** Removed a redundant O(m) inner loop in `membersDataSets` that iterated over the entire SCIM member set to find an email already confirmed by a direct map lookup. Benchmarks show ~16-19% improvement for large groups.

* **Goroutine leak safety:** Concurrent operations are verified with `synctest.Test` (Go 1.26) to ensure no goroutine leaks in both success and error paths.
