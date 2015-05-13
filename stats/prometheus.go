// Copyright (c) 2014 The SkyDNS Authors. All rights reserved.
// Use of this source code is governed by The MIT License (MIT) that can be
// found in the LICENSE file.

package server

import (
	"github.com/miekg/skydns/server"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	server.StatsForwardCount = newCounter(prometheus.NewCounter(prometheus.CounterOpts{
		Name: "dns_forward_count",
		Help: "Counter of DNS requests forwarded",
	}))
	server.StatsStubForwardCount = newCounter(prometheus.NewCounter(prometheus.CounterOpts{
		Name: "dns_stub_forward_count",
		Help: "Counter of DNS requests forwarded to stubs",
	}))
	// dns_stub_forward_error_count
	// dns_stub_update_count

	server.StatsLookupCount = newCounter(prometheus.NewCounter(prometheus.CounterOpts{
		Name: "dns_lookup_count",
		Help: "Counter of DNS lookups performed",
	}))
	server.StatsRequestCount = newCounter(prometheus.NewCounter(prometheus.CounterOpts{
		Name: "dns_request_count",
		Help: "Counter of DNS requests made",
	}))
	server.StatsDnssecOkCount = newCounter(prometheus.NewCounter(prometheus.CounterOpts{
		Name: "dns_dnssec_ok_count",
		Help: "Counter of DNSSEC requests",
	}))
	server.StatsDnssecCacheMiss = newCounter(prometheus.NewCounter(prometheus.CounterOpts{
		Name: "dns_dnssec_cache_miss_count",
		Help: "Counter of DNSSEC requests that missed the cache",
	}))
	server.StatsNameErrorCount = newCounter(prometheus.NewCounter(prometheus.CounterOpts{
		Name: "dns_name_error_count",
		Help: "Counter of DNS requests resulting in a name error",
	}))
	server.StatsNoDataCount = newCounter(prometheus.NewCounter(prometheus.CounterOpts{
		Name: "dns_no_data_count",
		Help: "Counter of DNS requests that contained no data",
	}))
}
