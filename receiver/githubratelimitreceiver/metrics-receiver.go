package githubratelimitreceiver

import (
	"context"
	"time"

	"github.com/google/go-github/v64/github"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

type githubRateLimitReceiver struct {
	host    component.Host
	cancel  context.CancelFunc
	metrics consumer.Metrics
	logger  *zap.Logger
	config  *Config
	client  *github.Client
}

func (grlr *githubRateLimitReceiver) Start(ctx context.Context, host component.Host) error {
	grlr.host = host
	grlr.client = github.NewClient(nil).WithAuthToken(grlr.config.Token)
	ctx, grlr.cancel = context.WithCancel(ctx)
	grlr.logger.Info("Starting github rate limit receiver")
	go grlr.scrape(ctx)
	return nil
}

func (grlr *githubRateLimitReceiver) Shutdown(ctx context.Context) error {
	grlr.logger.Info("Shutting down github rate limit receiver")
	if grlr.cancel != nil {
		grlr.cancel()
	}
	return nil
}

func (grlr *githubRateLimitReceiver) scrape(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(grlr.config.ScrapeInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			grlr.fetchAndEmitMetrics(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (grlr *githubRateLimitReceiver) fetchAndEmitMetrics(ctx context.Context) {
	grlr.logger.Debug("Fetching rate limit")
	rateLimits, _, err := grlr.client.RateLimits(ctx)
	if err != nil {
		grlr.logger.Error("Failed to fetch rate limit", zap.Error(err))
		return
	}
	grlr.emitMetrics(rateLimits)
}

func (grlr *githubRateLimitReceiver) emitMetrics(rateLimit *github.RateLimits) {
	metrics := pmetric.NewMetrics()
	rm := metrics.ResourceMetrics().AppendEmpty()
	ilm := rm.ScopeMetrics().AppendEmpty()
	ilm.Scope().SetName("github.com/open-telemetry/opentelemetry-collector-contrib/receiver/githubratelimitreceiver")
	ilm.Scope().SetVersion("v0.1.0")

	if rateLimit != nil {
		grlr.createAndAppendMetric(ilm, "github_rate_limit_remaining", "requests", float64(rateLimit.Core.Remaining))
	}

	if grlr.metrics != nil {
		grlr.metrics.ConsumeMetrics(context.Background(), metrics)
	} else {
		grlr.logger.Error("Metrics consumer is not initialized")
	}
}

func (grlr *githubRateLimitReceiver) createAndAppendMetric(ilm pmetric.ScopeMetrics, name, unit string, value float64) {
	metric := ilm.Metrics().AppendEmpty()
	metric.SetName(name)
	metric.SetUnit(unit)
	gauge := metric.SetEmptyGauge()
	dp := gauge.DataPoints().AppendEmpty()
	dp.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
	dp.SetDoubleValue(value)
	labels := map[string]string{
		"name":   grlr.config.Name,
		"target": grlr.config.Target,
	}
	for k, v := range labels {
		dp.Attributes().PutStr(k, v)
	}
}
