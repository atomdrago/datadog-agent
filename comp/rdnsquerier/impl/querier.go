// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.

package rdnsquerierimpl

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/DataDog/datadog-agent/comp/core/log"
)

type querier interface {
	start()
	stop()
	getHostnameAsync(addr string, updateHostname func(string, error)) error
}

// Standard querier implementation
type querierImpl struct {
	config            *rdnsQuerierConfig
	logger            log.Component
	internalTelemetry *rdnsQuerierTelemetry

	rateLimiter rateLimiter
	resolver    resolver

	// Context is used by the rate limiter and also for shutting down worker goroutines via its Done() channel.
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	queryChan chan *query
}

type query struct {
	addr           string
	updateHostname func(string, error)
}

func newQuerier(config *rdnsQuerierConfig, logger log.Component, internalTelemetry *rdnsQuerierTelemetry) querier {
	return &querierImpl{
		config:            config,
		logger:            logger,
		internalTelemetry: internalTelemetry,
		rateLimiter:       newRateLimiter(config),
		resolver:          newResolver(config),
	}
}

func (q *querierImpl) start() {
	q.ctx, q.cancel = context.WithCancel(context.Background())

	q.queryChan = make(chan *query, q.config.chanSize)

	for range q.config.workers {
		q.wg.Add(1)
		go q.worker()
	}
	q.logger.Infof("Reverse DNS Enrichment started %d workers", q.config.workers)
}

func (q *querierImpl) stop() {
	q.cancel()
	q.wg.Wait()

	q.logger.Infof("Reverse DNS Enrichment stopped workers")
}

func (q *querierImpl) worker() {
	defer q.wg.Done()
	for {
		select {
		case query := <-q.queryChan:
			q.getHostname(query)
		case <-q.ctx.Done():
			return
		}
	}
}

func (q *querierImpl) getHostname(query *query) { //JMWTUE if error send error to callback function
	err := q.rateLimiter.wait(q.ctx)
	if err != nil {
		q.internalTelemetry.droppedRateLimiter.Inc()
		q.logger.Debugf("Reverse DNS Enrichment rateLimiter.wait() returned error: %v - dropping query for IP address %s", err, query.addr)
		query.updateHostname("", err) //JMWTUE
		return
	}

	hostname, err := q.resolver.lookup(query.addr)
	if err != nil {
		//JMW add test for these errors - and others
		if dnsErr, ok := err.(*net.DNSError); ok {
			if dnsErr.IsNotFound {
				q.internalTelemetry.lookupErrNotFound.Inc()
				q.logger.Debugf("Reverse DNS Enrichment net.LookupAddr returned not found error '%v' for IP address %v", err, query.addr)
				// no match was found for the requested IP address, so call updateHostname() to make the caller aware of that fact
				query.updateHostname("", nil) //JMW return "", nil, or return w/ err and let caller check error and handle IsNotFound?
				//JMW test this path w/ cache layer
				return
			}
			if dnsErr.IsTimeout {
				q.internalTelemetry.lookupErrTimeout.Inc()
				q.logger.Debugf("Reverse DNS Enrichment net.LookupAddr returned timeout error '%v' for IP address %v", err, query.addr)
				return
			}
			if dnsErr.IsTemporary {
				q.internalTelemetry.lookupErrTemporary.Inc()
				q.logger.Debugf("Reverse DNS Enrichment net.LookupAddr returned temporary error '%v' for IP address %v", err, query.addr)
				return
			}
		}
		q.internalTelemetry.lookupErrOther.Inc()
		q.logger.Debugf("Reverse DNS Enrichment net.LookupAddr returned error '%v' for IP address %v", err, query.addr)
		query.updateHostname(hostname, err)
		return
	}

	q.internalTelemetry.successful.Inc()
	q.logger.Tracef("JMW Reverse DNS Enrichment q.resolver.lookup successfully returned hostname for IP address %v - calling updateHostname(%s, nil)", query.addr, hostname, query.addr)
	query.updateHostname(hostname, nil)
}

func (q *querierImpl) getHostnameAsync(addr string, updateHostname func(string, error)) error {
	select {
	case q.queryChan <- &query{addr: addr, updateHostname: updateHostname}:
		q.internalTelemetry.chanAdded.Inc()
		return nil
	default:
		q.internalTelemetry.droppedChanFull.Inc()
		return fmt.Errorf("channel is full, dropping query for IP address %s", addr)
	}
}
