package events

import (
	"log"
	"sync"
)

type Event struct {
	EventType string
	Content   map[string]interface{}
}

type EventDispatcher struct {
	Listeners []Listener
	// Add world state
}

func (ed *EventDispatcher) AddEvent(event Event) {
	log.Println("Event Dispathcer AddEvent: ", event.EventType)
	log.Printf("Will run in #%d: snakes\n", len(ed.Listeners))
	var wg sync.WaitGroup
	wg.Add(len(ed.Listeners))
	for _, l := range ed.Listeners {
		go l.HandleEvent(event, &wg)
	}
	wg.Wait()
}

type Listener interface {
	HandleEvent(event Event, wg *sync.WaitGroup)
}
