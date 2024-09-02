# GitHub Rate Limit Receiver

The GitHub Rate Limit Receiver is a component of the OpenTelemetry Collector that fetches and emits GitHub API rate limit metrics. This receiver is useful for monitoring the rate limit usage of your GitHub API requests.

## Features

- Fetches GitHub API rate limit metrics.
- Emits metrics in OpenTelemetry format.
- Configurable scrape interval.

## Configuration

To use the GitHub Rate Limit Receiver, add the following configuration to your OpenTelemetry Collector configuration file:

```yaml
receivers:
  githubratelimit:
    token: "<your_github_token>"
    scrape_interval: 60 # in seconds
```

## Example
```yaml
receivers:
  githubratelimit:
    token: "your_github_token"
    scrape_interval: 60

exporters:
  logging:
    loglevel: debug

service:
  pipelines:
    metrics:
      receivers: [githubratelimit]
      exporters: [logging]
```
