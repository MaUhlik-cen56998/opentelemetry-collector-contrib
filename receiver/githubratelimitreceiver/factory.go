package githubratelimitreceiver

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

var (
	typeStr = component.MustNewType("githubratelimit")
)

const (
	defaultEndpoint = "https://api.github.com/rate_limit"
)

func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		createDefaultConfig,
		receiver.WithMetrics(createMetricsReceiver, component.StabilityLevelAlpha),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		Endpoint: defaultEndpoint,
		LogLevel: "error",
	}
}

func createMetricsReceiver(ctx context.Context, params receiver.Settings, cfg component.Config, consumer consumer.Metrics) (receiver.Metrics, error) {
	rCfg := cfg.(*Config)
	if err := rCfg.Validate(); err != nil {
		return nil, err
	}

	return &githubRateLimitReceiver{
		logger:  params.Logger,
		config:  rCfg,
		metrics: consumer,
	}, nil
}
