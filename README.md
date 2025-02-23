# renoglaab

`renoglaab` is a Go application designed to run in a GitLab CI job and automatically merge Renovate merge requests based on predefined conditions.

## Prerequisites

- Renovate configured to create MRs
- GitLab Personal Access Token (PAT) or Group Access Token with `api` scope.

## Usage

Add the following job to your `.gitlab-ci.yml`:
```yml
renoglaab:
  image: ghcr.io/xmoelletschi/renoglaab:latest
  script:
    - /renoglaab
```
To run this job in a scheduled pipeline, add a schedule to your GitLab project under **CI / CD > Schedules**.

## Configuration

Modify environment variables to configure merging behavior:

| Variable                              | Description                                      | Default                           | Options                           |
|---------------------------------------|--------------------------------------------------|-----------------------------------|-----------------------------------|
| `CONFIG_PATH`                         | Path to the configuration file                   | `$CI_PROJECT_DIR/config.js`       | Any valid file path               |
| `LOG_LEVEL`                           | Logging level (`debug`, `info`, `warn`, `error`) | `info`                            | `debug`, `info`, `warn`, `error`  |
| `GITLAB_API_TOKEN`                    | GitLab API token (required)                      |                                   | Any valid token                   |
| `GITLAB_URL`                          | GitLab instance URL                              | `https://gitlab.com`              | Any valid URL                     |
| `FILTER_BY_AUTHOR_USERNAME`           | Filter MRs by author username                    | `true`                            | `true`, `false`                   |
| `AUTHOR_USERNAME`                     | Author username to filter MRs                    | `renovate-bot`                    | Any valid username                |
| `FILTER_BY_LABELS`                    | Filter MRs by labels                             | `true`                            | `true`, `false`                   |
| `LABELS`                              | Labels to filter MRs                             | `renovate`                        | Any valid label                   |
| `FILTER_BY_BRANCH`                    | Filter MRs by branch regex                       | `true`                            | `true`, `false`                   |
| `ALLOWED_BRANCH_REGEX`                | Regex for allowed branches                       | `renovate/automerge`            | Any valid regex                   |
| `FILTER_BY_SUCCEEDED_PIPELINE`        | Filter MRs by succeeded pipeline                 | `true`                            | `true`, `false`                   |
| `FILTER_BY_PIPELINE_WITHOUT_WARNINGS` | Filter MRs by pipeline without warnings          | `true`                            | `true`, `false`                   |
| `ADD_COMMENT`                         | Add a comment to the MR                          | `true`                            | `true`, `false`                   |
| `COMMENT`                             | Comment to add to the MR                         | `Approving merge request! :ship:` | Any valid comment                 |
| `APPROVE`                             | Approval command                                 | `/approve`                        | Any valid command                 |

## Examples

For a real-world example, visit the [renoglaab GitLab group](https://gitlab.com/renoglaab). [currently in WIP]

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## Support

For support or questions, open a discussion on the [GitHub repository](https://github.com/xMoelletschi/renoglaab/discussions).

## License

This project is licensed under the [MIT License](LICENSE).
