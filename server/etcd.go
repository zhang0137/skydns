// Copyright (c) 2014 The SkyDNS Authors. All rights reserved.
// Use of this source code is governed by The MIT License (MIT) that can be
// found in the LICENSE file.

package server

import (
	"github.com/coreos/go-etcd/etcd"
	"github.com/miekg/skydns/msg"
)

// Ectd implementation

type etc etcd.Client

func (client *etc) Records(path string, recursive bool) []msg.Service {
	response, err := get(client, path, recursive)
	if err != nil {
		return nil, err
	}
	// process them to msg.Service and use those in server.go
}

// get is a wrapper for client.Get that uses SingleInflight to suppress multiple outstanding queries.
func get(client *etcd.Client, path string, recursive bool) (*etcd.Response, error) {
	resp, err, _ := etcdInflight.Do(path, func() (*etcd.Response, error) {
		r, e := client.Get(path, false, recursive)
		if e != nil {
			return nil, e
		}
		return r, e
	})
	if err != nil {
		return resp, err
	}
	// shared?
	return resp, err
}
