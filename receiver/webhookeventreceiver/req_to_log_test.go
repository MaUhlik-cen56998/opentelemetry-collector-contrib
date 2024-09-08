// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package webhookeventreceiver

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/receivertest"
)

func TestReqToLog(t *testing.T) {
	defaultConfig := createDefaultConfig().(*Config)

	tests := []struct {
		desc string
		sc   *bufio.Scanner
		r    *http.Request
		tt   func(t *testing.T, reqLog plog.Logs, reqLen int, settings receiver.Settings)
	}{
		{
			desc: "Valid query valid event",
			sc: func() *bufio.Scanner {
				reader := io.NopCloser(bytes.NewReader([]byte("this is a: log")))
				return bufio.NewScanner(reader)
			}(),
			r: func() *http.Request {
				req, err := http.NewRequest("GET", "http://localhost:8080?param1=hello&param2=world", nil)
				if err != nil {
					log.Fatal("failed to create request")
				}
				return req
			}(),
			tt: func(t *testing.T, reqLog plog.Logs, reqLen int, _ receiver.Settings) {
				require.Equal(t, 1, reqLen)

				attributes := reqLog.ResourceLogs().At(0).Resource().Attributes()
				require.Equal(t, 2, attributes.Len())

				scopeLogsScope := reqLog.ResourceLogs().At(0).ScopeLogs().At(0).Scope()
				require.Equal(t, 2, scopeLogsScope.Attributes().Len())

				if v, ok := attributes.Get("query.param1"); ok {
					require.Equal(t, "hello", v.AsString())
				} else {
					require.Fail(t, "failed to set attribute from query parameter 1")
				}
				if v, ok := attributes.Get("query.param2"); ok {
					require.Equal(t, "world", v.AsString())
				} else {
					require.Fail(t, "failed to set attribute query parameter 2")
				}
			},
		},
		{
			desc: "Query is empty",
			sc: func() *bufio.Scanner {
				reader := io.NopCloser(bytes.NewReader([]byte("this is a: log")))
				return bufio.NewScanner(reader)
			}(),
			r: func() *http.Request {
				req, err := http.NewRequest("GET", "http://localhost:8080", nil)
				if err != nil {
					log.Fatal("failed to create request")
				}
				return req
			}(),
			tt: func(t *testing.T, reqLog plog.Logs, reqLen int, _ receiver.Settings) {
				require.Equal(t, 1, reqLen)

				attributes := reqLog.ResourceLogs().At(0).Resource().Attributes()
				require.Equal(t, 0, attributes.Len())

				scopeLogsScope := reqLog.ResourceLogs().At(0).ScopeLogs().At(0).Scope()
				require.Equal(t, 2, scopeLogsScope.Attributes().Len())
			},
		},
		{
			desc: "Validate headers",
			sc: func() *bufio.Scanner {
				reader := io.NopCloser(bytes.NewReader([]byte("this is a: log")))
				return bufio.NewScanner(reader)
			}(),
			r: func() *http.Request {
				req, err := http.NewRequest("GET", "http://localhost:8080", nil)
				if err != nil {
					log.Fatal("failed to create request")
				}
				req.Header.Add("header1", "value1")
				req.Header.Add("header2", "value2")
				return req
			}(),
			tt: func(t *testing.T, reqLog plog.Logs, reqLen int, _ receiver.Settings) {

				attributes := reqLog.ResourceLogs().At(0).Resource().Attributes()
				require.Equal(t, 2, attributes.Len())

				if v, ok := attributes.Get("header.Header1"); ok {
					require.Equal(t, "value1", v.AsString())
				} else {
					require.Fail(t, "failed to set attribute from header 1")
				}
				if v, ok := attributes.Get("header.Header2"); ok {
					require.Equal(t, "value2", v.AsString())
				} else {
					require.Fail(t, "failed to set attribute from header 2")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			reqLog, reqLen := reqToLog(test.sc, test.r, defaultConfig, receivertest.NewNopSettings())
			test.tt(t, reqLog, reqLen, receivertest.NewNopSettings())
		})
	}
}
