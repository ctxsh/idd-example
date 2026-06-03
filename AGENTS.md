# Agent Instructions

## Issue Tracker

Use GitHub Issues as the issue tracker for this repository. Reference issues
with GitHub issue numbers or links in plans, PR descriptions, and review notes.

## Development Commands

Use the Makefile targets for Go development tasks:

- `make deps` installs local development tools into `./bin`.
- `make build` builds the service into `./bin`.
- `make test` runs the Go test suite.
- `make lint` runs the local `golangci-lint` binary.
- `make fmt` formats Go files.
- `make tidy` tidies Go modules.

Do not install repo tooling globally. Keep tool binaries and build outputs under
`./bin`; the Makefile also keeps Go caches under `./.cache`.

## PR Messaging Style

When preparing PR text:

- Start with an overview that briefly describes the feature or issue and
  summarizes the change. For larger changes, call out high-level design
  decisions and link relevant design docs or tracking issues.
- Include customer-facing messaging only when there is a real customer-facing
  change. Write it as a short summary that could be shared with customers, and
  avoid internal implementation details.
- Make validation explicit. Summarize tests added, commands run, integration
  tests, screenshots, or other evidence.
- Include rollout and rollback notes. Use `N/A` only when rollout or rollback
  does not apply.
- Link tracking issues with the appropriate action keyword.

Use this repository's PR template as follows:

- `Change Summary` should cover the overview.
- `Risks` should include validation gaps, rollout concerns, rollback concerns,
  and other known risks.
- `Issues` should contain the tracking issue action line.

## PR Titles And Commit Messages

Use the same subject style for PR titles and commit messages:

```text
<type>(<scope>): <summary>
```

- Use a specific type such as `feat`, `fix`, `docs`, `chore`, `refactor`,
  `test`, or `infra`.
- Use a short lowercase scope for the package, component, or area.
- Write the summary in lowercase, without a trailing period.
- Keep the summary focused on the behavior or repo change, not the tool used to
  make it.
- Put issue references in the PR body, not the title.
- Do not use vague titles like `chore: updates` or `fix: misc fixes`.

Examples:

- `docs(readme): describe deployment api layout`
- `chore(go): initialize module scaffold`
- `infra(github): add git town workflow`

## Pull Requests

When creating or updating a pull request, use
`.github/pull_request_template.md` as the PR description structure.

Fill out every section:

- `Change Summary` should explain the user-visible or reviewer-relevant change.
- `Risks` should call out known risks, tradeoffs, migrations, rollback concerns,
  or validation gaps. Use `None identified` only when that is accurate.
- `Issues` should be updated with `<action> <issue>` entries, such as
  `Closes #123`, `Fixes #123`, `Resolves #123`, or `Related to #123`.
- `Stack` is reserved for Git Town. Leave the `<!-- branch-stack -->` marker in
  place so the Git Town GitHub Action can update the branch stack information.

Keep issue references concrete by using issue numbers or links when available.
Do not remove template sections just because a section is short.
