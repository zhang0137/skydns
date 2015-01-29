// Copyright (c) 2014 The SkyDNS Authors. All rights reserved.
// Use of this source code is governed by The MIT License (MIT) that can be
// found in the LICENSE file.

package server

import "github.com/miekg/skydns/msg"

type Getter interface {
	Records(path string, recursive bool) ([]msg.Service, error)
	// Config
}
