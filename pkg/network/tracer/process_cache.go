// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build linux_bpf || (windows && npm)

package tracer

import (
	"fmt"
	"sync"
	"time"

	"github.com/cihub/seelog"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/DataDog/datadog-agent/pkg/network/events"
	"github.com/DataDog/datadog-agent/pkg/telemetry"
	"github.com/DataDog/datadog-agent/pkg/util/log"
)

const (
	maxProcessQueueLen = 100
	// maxProcessListSize is the max size of a processList
	maxProcessListSize     = 3
	processCacheModuleName = "network_tracer__process_cache"
	defaultExpiry          = 2 * time.Minute
)

var processCacheTelemetry = struct {
	cacheEvicts   telemetry.Counter
	cacheLength   *prometheus.Desc
	eventsDropped telemetry.Counter
	eventsSkipped telemetry.Counter
}{
	telemetry.NewCounter(processCacheModuleName, "cache_evicts", []string{}, "Counter measuring the number of evictions in the process cache"),
	prometheus.NewDesc(processCacheModuleName+"__cache_length", "Gauge measuring the current size of the process cache", nil, nil),
	telemetry.NewCounter(processCacheModuleName, "events_dropped", []string{}, "Counter measuring the number of dropped process events"),
	telemetry.NewCounter(processCacheModuleName, "events_skipped", []string{}, "Counter measuring the number of skipped process events"),
}

type processList []*events.Process

// ProcessCache holds data from the process events to be queried later
type ProcessCache struct {
	mu sync.Mutex

	// cache of pid -> list of processes holds a list of processes
	// with the same pid but differing start times up to a max of
	// maxProcessListSize. this is used to determine the closest
	// match to a connection's timestamp
	cacheByPid map[uint32]processList
	// lru cache; keyed by (pid, start time)
	cache *lru.Cache[processCacheKey, *events.Process]

	in      chan *events.Process
	stopped chan struct{}
	stop    sync.Once
}

type processCacheKey struct {
	pid       uint32
	startTime int64
}

// NewProcessCache creates a new cache to hold process events
func NewProcessCache(maxProcs int) (*ProcessCache, error) {
	pc := &ProcessCache{
		cacheByPid: map[uint32]processList{},
		in:         make(chan *events.Process, maxProcessQueueLen),
		stopped:    make(chan struct{}),
	}

	var err error
	pc.cache, err = lru.NewWithEvict(maxProcs, func(_ processCacheKey, p *events.Process) {
		log.TraceFunc(func() string { return fmt.Sprintf("evicting process %+v", p) })
		//nolint:gosimple // TODO(NET) Fix gosimple linter
		pl, _ := pc.cacheByPid[p.Pid]
		if pl = pl.remove(p); len(pl) == 0 {
			delete(pc.cacheByPid, p.Pid)
			return
		}

		pc.cacheByPid[p.Pid] = pl
	})

	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case <-pc.stopped:
				return
			case p := <-pc.in:
				pc.add(p)
			}
		}
	}()

	return pc, nil
}

func (pc *ProcessCache) HandleProcessEvent(entry *events.Process) {

	select {
	case <-pc.stopped:
		return
	default:
	}

	p := pc.processEvent(entry)
	if p == nil {
		processCacheTelemetry.eventsSkipped.Inc()
		return
	}

	select {
	case pc.in <- p:
	default:
		// dropped
		processCacheTelemetry.eventsDropped.Inc()
	}
}

func (pc *ProcessCache) processEvent(entry *events.Process) *events.Process {
	if len(entry.Tags) == 0 && entry.ContainerID == nil {
		return nil
	}

	return entry
}

func (pc *ProcessCache) Trim() {
	if pc == nil {
		return
	}

	pc.mu.Lock()
	defer pc.mu.Unlock()

	now := time.Now().Unix()
	trimmed := 0
	for _, v := range pc.cache.Values() {
		if now > v.Expiry {
			// Remove will call the evict callback which will
			// delete from the cacheByPid map
			log.TraceFunc(func() string {
				return fmt.Sprintf("trimming process %+v", v)
			})
			pc.cache.Remove(processCacheKey{pid: v.Pid, startTime: v.StartTime})
			trimmed++
		}
	}

	if trimmed > 0 {
		log.Debugf("Trimmed %d process cache entries", trimmed)
	}
}

func (pc *ProcessCache) Stop() {
	if pc == nil {
		return
	}

	pc.stop.Do(func() { close(pc.stopped) })
}

func (pc *ProcessCache) add(p *events.Process) {
	if pc == nil {
		return
	}

	pc.mu.Lock()
	defer pc.mu.Unlock()

	if log.ShouldLog(seelog.TraceLvl) {
		log.Tracef("adding process %+v to process cache", p)
	}
	p.Expiry = time.Now().Add(defaultExpiry).Unix()
	if evicted := pc.cache.Add(processCacheKey{pid: p.Pid, startTime: p.StartTime}, p); evicted {
		processCacheTelemetry.cacheEvicts.Inc()
	}

	pl := pc.cacheByPid[p.Pid]
	pc.cacheByPid[p.Pid] = pl.update(p)
}

func (pc *ProcessCache) Get(pid uint32, ts int64) (*events.Process, bool) {
	if pc == nil {
		return nil, false
	}

	pc.mu.Lock()
	defer pc.mu.Unlock()

	log.TraceFunc(func() string { return fmt.Sprintf("looking up pid %d %v", pid, ts) })

	pl := pc.cacheByPid[pid]
	if closest := pl.closest(ts); closest != nil {
		closest.Expiry = time.Now().Add(defaultExpiry).Unix()
		pc.cache.Get(processCacheKey{pid: closest.Pid, startTime: closest.StartTime})
		log.TraceFunc(func() string { return fmt.Sprintf("found entry for pid %d: %+v", pid, closest) })
		return closest, true
	}

	log.TraceFunc(func() string { return fmt.Sprintf("entry not found for process %d", pid) })
	return nil, false
}

func (pc *ProcessCache) Dump() (interface{}, error) {
	res := map[uint32]interface{}{}
	if pc == nil {
		return res, nil
	}

	pc.mu.Lock()
	defer pc.mu.Unlock()

	for pid, pl := range pc.cacheByPid {
		res[pid] = pl
	}

	return res, nil
}

// Describe returns all descriptions of the collector.
func (pc *ProcessCache) Describe(ch chan<- *prometheus.Desc) {
	ch <- processCacheTelemetry.cacheLength
}

// Collect returns the current state of all metrics of the collector.
func (pc *ProcessCache) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(processCacheTelemetry.cacheLength, prometheus.GaugeValue, float64(pc.cache.Len()))
}

func (pl processList) update(p *events.Process) processList {
	for i := range pl {
		if pl[i].StartTime == p.StartTime {
			pl[i] = p
			return pl
		}
	}

	if len(pl) == maxProcessListSize {
		copy(pl, pl[1:])
		pl = pl[:len(pl)-1]
	}

	if pl == nil {
		pl = make(processList, 0, maxProcessListSize)
	}

	return append(pl, p)
}

func (pl processList) remove(p *events.Process) processList {
	for i := range pl {
		if pl[i] == p {
			return append(pl[:i], pl[i+1:]...)
		}
	}

	return pl
}

func abs(i int64) int64 {
	if i < 0 {
		return -i
	}

	return i
}

func (pl processList) closest(ts int64) *events.Process {
	var closest *events.Process
	for i := range pl {
		if ts >= pl[i].StartTime &&
			(closest == nil ||
				abs(closest.StartTime-ts) > abs(pl[i].StartTime-ts)) {
			closest = pl[i]
		}
	}

	return closest
}
