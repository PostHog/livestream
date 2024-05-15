package main

import (
	"sync/atomic"
)

type Subscription struct {
	// Filters
	Token      string
	DistinctId string
	EventType  string

	// Response channel
	EventChan   chan interface{}
	ShouldClose *atomic.Bool
}

type ResponsePostHogEvent struct {
	Timestamp  string `json:"timestamp"`
	DistinctId string `json:"distinct_id"`
	PersonId   string `json:"person_id"`
	Event      string `json:"event"`
}

type ResponseGeoEvent struct {
	Lat   float64 `json:"lat"`
	Lng   float64 `json:"lng"`
	Count int32   `json:"count"`
}

type Filter struct {
	inboundChan chan PostHogEvent
	subChan     chan *Subscription
	subs        []Subscription
}

func NewFilter(inboundChan chan PostHogEvent) *Filter {
	return &Filter{inboundChan: inboundChan}
}

func (c *Filter) Run() {
	select {
	case event := <-c.inboundChan:
		for i := 0; i < len(c.subs); i++ {
			sub := c.subs[i]

			if sub.ShouldClose.Load() {
				// TODO: Figure this out later. Apparently closing from the read side is dangerous
				// because writing to a closed channel = panic.
				continue
			}

			if sub.Token != "" && event.Token != sub.Token {
				continue
			}

			if sub.DistinctId != "" && event.DistinctID != sub.DistinctId {
				continue
			}

			if sub.EventType != "" && event.Event != sub.EventType {
				continue
			}

			sub.EventChan <- event
		}
	case newSub := <-c.subChan:
		c.subs = append(c.subs, *newSub)
	}
}