// Copyright (c) 2014 The SkyDNS Authors. All rights reserved.
// Use of this source code is governed by The MIT License (MIT) that can be
// found in the LICENSE file.

package server

import "log"

type BackendFailure int

const (
	etcdFailure BackendFailure = iota + 1
	otherFailure
)

// Printf calls log.Printf with the parameters given.
func Logf(format string, a ...interface{}) {
	log.Printf("skydns: "+format, a...)
}

// Fatalf calls log.Fatalf with the parameters given.
func Fatalf(format string, a ...interface{}) {
	log.Fatalf("skydns: "+format, a...)
}

// log and Inc the promBackendFailureCount{type="typ"}.
func logBackendFailure(typ BackendFailure, format string, a ...interface{}) {
	Logf(format, a...)
	switch typ {
	case otherFailure:
		promBackendFailureCount.WithLabelValues("other").Inc()
	case etcdFailure:
		promBackendFailureCount.WithLabelValues("etcd").Inc()
	}
}
