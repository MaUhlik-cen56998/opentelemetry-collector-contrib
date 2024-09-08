// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package webhookeventreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/webhookeventreceiver"

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/receiver"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/webhookeventreceiver/internal/metadata"
)

func reqToLog(sc *bufio.Scanner,
	r *http.Request,
	_ *Config,
	settings receiver.Settings) (plog.Logs, int) {
	log := plog.NewLogs()
	resourceLog := log.ResourceLogs().AppendEmpty()
	appendMetadata(resourceLog, "query.", r.URL.Query())
	appendMetadata(resourceLog, "header.", r.Header)
	scopeLog := resourceLog.ScopeLogs().AppendEmpty()

	scopeLog.Scope().SetName(scopeLogName)
	scopeLog.Scope().SetVersion(settings.BuildInfo.Version)
	scopeLog.Scope().Attributes().PutStr("source", settings.ID.String())
	scopeLog.Scope().Attributes().PutStr("receiver", metadata.Type.String())

	for sc.Scan() {
		logRecord := scopeLog.LogRecords().AppendEmpty()
		logRecord.SetObservedTimestamp(pcommon.NewTimestampFromTime(time.Now()))
		line := sc.Text()
		logRecord.Body().SetStr(line)
	}

	return log, scopeLog.LogRecords().Len()
}

func appendMetadata(resourceLog plog.ResourceLogs, prefix string, params interface{}) {
	switch p := params.(type) {
	case url.Values:
		for k := range p {
			if p.Get(k) != "" {
				resourceLog.Resource().Attributes().PutStr(fmt.Sprintf("%s%s", prefix, k), p.Get(k))
			}
		}
	case http.Header:
		for k, v := range p {
			if len(v) > 0 {
				resourceLog.Resource().Attributes().PutStr(fmt.Sprintf("%s%s", prefix, k), v[0])
			}
		}
	default:
		// handle other types if needed
	}
}
