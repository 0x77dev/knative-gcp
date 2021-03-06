/*
Copyright 2020 Google LLC
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	nethttp "net/http"
	"time"
)

// healthChecker carries timestamps of the latest handled events from the
// forwarder and receiver, as well as the longest tolerated staleness time.
// The healthChecker is expected to declare the probe helper as unhealthy if
// the probe helper is unable to handle either sort of event.
type healthChecker struct {
	lastProbeEventTimestamp    eventTimestamp
	lastReceiverEventTimestamp eventTimestamp
	maxStaleDuration           time.Duration
}

// stalenessHandlerFunc returns the HTTP handler for probe helper health checks.
func (c *healthChecker) stalenessHandlerFunc() nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, req *nethttp.Request) {
		if req.URL.Path != "/healthz" {
			w.WriteHeader(nethttp.StatusNotFound)
			return
		}
		now := time.Now()
		if now.Sub(c.lastProbeEventTimestamp.getTime()) > c.maxStaleDuration ||
			now.Sub(c.lastReceiverEventTimestamp.getTime()) > c.maxStaleDuration {
			w.WriteHeader(nethttp.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(nethttp.StatusOK)
	}
}
