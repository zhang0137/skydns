// Copyright (c) 2014 The SkyDNS Authors. All rights reserved.
// Use of this source code is governed by The MIT License (MIT) that can be
// found in the LICENSE file.

package server

import (
	"encoding/json"
	"strings"

	"github.com/coreos/go-etcd/etcd"
	"github.com/skynetservices/skydns/msg"
)

// Ectd implementation

type etc struct {
	c *etcd.Client
}

func (client *etc) Records(p string, rec bool) ([]*msg.Service, error) {
	path, star := msg.PathWithWildcard(p)
	r, err := get(client.c, path, rec)
	if err != nil {
		return nil, err
	}
	if !r.Node.Dir {
		serv := new(msg.Service)
		if err := json.Unmarshal([]byte(r.Node.Value), serv); err != nil {
			return nil, err
		}
		ttl := calculateTtl(r.Node, serv)
		serv.Key = r.Node.Key
		serv.Ttl = ttl
		return []*msg.Service{serv}, nil
	}
	sx, err := loopNodes(&r.Node.Nodes, strings.Split(p, "/"), star, nil)
	if err != nil || len(sx) == 0 {
		return nil, err
	}
	for _, serv := range sx {
		ttl := calculateTtl(r.Node, serv)
		serv.Key = r.Node.Key
		serv.Ttl = ttl
	}
	return sx, nil
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

type bareService struct {
	Host     string
	Port     int
	Priority int
	Weight   int
	Text     string
}

// skydns/local/skydns/east/staging/web
// skydns/local/skydns/west/production/web
//
// skydns/local/skydns/*/*/web
// skydns/local/skydns/*/web

// loopNodes recursively loops through the nodes and returns all the values. The nodes' keyname
// will be match against any wildcards when star is true.
func loopNodes(n *etcd.Nodes, nameParts []string, star bool, bx map[bareService]bool) (sx []*msg.Service, err error) {
	if bx == nil {
		bx = make(map[bareService]bool)
	}
Nodes:
	for _, n := range *n {
		if n.Dir {
			nodes, err := loopNodes(&n.Nodes, nameParts, star, bx)
			if err != nil {
				return nil, err
			}
			sx = append(sx, nodes...)
			continue
		}
		if star {
			keyParts := strings.Split(n.Key, "/")
			for i, n := range nameParts {
				if i > len(keyParts)-1 {
					// name is longer than key
					continue Nodes
				}
				if n == "*" {
					continue
				}
				if keyParts[i] != n {
					continue Nodes
				}
			}
		}
		serv := new(msg.Service)
		if err := json.Unmarshal([]byte(n.Value), serv); err != nil {
			return nil, err
		}
		b := bareService{serv.Host, serv.Port, serv.Priority, serv.Weight, serv.Text}
		if _, ok := bx[b]; ok {
			continue
		}
		bx[b] = true
		serv.Ttl = calculateTtl(n, serv)
		serv.Key = n.Key
		sx = append(sx, serv)
	}
	return sx, nil
}

// calculateTtl returns the smaller of the etcd TTL and the service's
// TTL. If neither of these are set (have a zero value), it leave it at zero.
func calculateTtl(node *etcd.Node, serv *msg.Service) uint32 {
	etcdTtl := uint32(node.TTL)

	if etcdTtl == 0 {
		return serv.Ttl
	}
	if serv.Ttl == 0 {
		return etcdTtl
	}
	if etcdTtl < serv.Ttl {
		return etcdTtl
	}
	return serv.Ttl
}
