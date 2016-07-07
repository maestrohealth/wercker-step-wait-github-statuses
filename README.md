# wercker-step-wait-github-statuses

This step will wait for successful GitHub Status API reports against the
commit that triggered the current build, for one or more configured statuses.
The step will complete successfully when all of the configured statuses have
reported successful in GitHub, and will fail if any of them report failure.

For statuses that have not reported yet, the step will continue to wait until 
the configured timeout has elapsed.

This step is intended to be used as a workaround for joining two or more
Wercker Workflows pipelines back into a single pipeline, until Workflows 
supports that natively. It may also be useful for waiting on external processes
that report back to GitHub.

# Options

- `status_contexts` Comma-separated list of status context names to wait for.
  Example: "wercker/build,wercker/pipeline-1" will wait for the main Wercker
  build pipeline and the "pipeline-1" pipeline to complete.
- `timeout` Timeout in minutes; defaults to 25.

# Example

```yaml
build:
    steps:
      - maestrohealthcaretechnologies/wait-github-statuses
          status_contexts: wercker/pipeline-1,wercker/pipeline-2
          timeout: 10
```

# License

TBD

# Changelog

## 0.1.0

- Initial release

