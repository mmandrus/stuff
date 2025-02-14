package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// chance that there is an error in the second phase of saving the object
const chanceOfError = 0.1

var (
	errNotFound = errors.New("not found")
	errNotReady = errors.New("not ready")
	errFatal    = errors.New("fatal")
)

const (
	statusPending = iota
	statusCreated
)

// object to be saved and retrieved
type object struct {
	name        string
	description string
	value       string
}

func (o object) String() string {
	return fmt.Sprintf("name: %s\ndescription: %s\nvalue: %s\n", o.name, o.description, o.value)
}

// object split into parts
type objectMeta struct {
	name        string //id
	description string
	status      int
}

type objectValue struct {
	value string
}

// events
type objMetaSavedEvent struct {
	name  string
	value string
}

type valueSavedEvent struct {
	name string
}

// primary struct that does stuff
type worker struct {
	metaMu      sync.RWMutex
	objMetadata map[string]objectMeta

	valueMu  sync.RWMutex
	objValue map[string]objectValue

	objMetadataSavedEvents []objMetaSavedEvent
	objValueSavedEvents    []valueSavedEvent
}

// store the object metadata and put an event in the queue
func (w *worker) createObject(o object) {
	w.metaMu.Lock()
	defer w.metaMu.Unlock()

	w.objMetadata[o.name] = objectMeta{o.name, o.description, statusPending}
	w.objMetadataSavedEvents = append(w.objMetadataSavedEvents, objMetaSavedEvent{o.name, o.value})
}

// return the object, or an error if it is not available yet
func (w *worker) getObject(name string) (*object, error) {
	w.metaMu.RLock()
	defer w.metaMu.RUnlock()

	m, ok := w.objMetadata[name]
	if !ok {
		return nil, errNotFound
	}

	if m.status != statusCreated {
		return nil, errNotReady
	}

	w.valueMu.RLock()
	defer w.valueMu.RUnlock()
	return &object{m.name, m.description, w.objValue[name].value}, nil
}

// process a metadata event from the queue (i.e. save the value and put an event in the queue) and pop if successful.
// then add a value saved event to the queue
// simulates random errors during the value saving to demonstrate eventuality
func (w *worker) processMetaSavedEvents() {
	w.metaMu.RLock()
	defer w.metaMu.RUnlock()

	if len(w.objMetadataSavedEvents) == 0 {
		return
	}

	e := w.objMetadataSavedEvents[0]
	w.valueMu.Lock()
	defer w.valueMu.Unlock()

	if simulateError() != nil {
		return
	}

	w.objValue[e.name] = objectValue{e.value}
	w.objMetadataSavedEvents = w.objMetadataSavedEvents[1:]
	w.objValueSavedEvents = append(w.objValueSavedEvents, valueSavedEvent{e.name})
}

// process a value event from the queue (i.e. update the metadata status to 'ready') and pop if successful
func (w *worker) processValueSavedEvents() {
	w.valueMu.RLock()
	defer w.valueMu.RUnlock()

	if len(w.objValueSavedEvents) == 0 {
		return
	}

	e := w.objValueSavedEvents[0]
	w.metaMu.Lock()
	defer w.metaMu.Unlock()

	o := w.objMetadata[e.name]
	o.status = statusCreated
	w.objMetadata[e.name] = o
	w.objValueSavedEvents = w.objValueSavedEvents[1:]
}

// n percent chance an error will happen
func simulateError() error {
	if rand.Float32() < chanceOfError {
		return errFatal
	}
	return nil
}

func main() {
	w := worker{
		objMetadata:            make(map[string]objectMeta),
		objValue:               make(map[string]objectValue),
		objMetadataSavedEvents: make([]objMetaSavedEvent, 0),
		objValueSavedEvents:    make([]valueSavedEvent, 0),
	}

	ctx := context.Background()
	defer ctx.Done()

	metaTicker := time.NewTicker(10 * time.Millisecond)
	valueTicker := time.NewTicker(10 * time.Millisecond)

	// start pulling from the queues at a frequent interval
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-metaTicker.C:
				w.processMetaSavedEvents()
			}
		}
	}(ctx)

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-valueTicker.C:
				w.processValueSavedEvents()
			}
		}
	}(ctx)

	awaitingCompletion := make(map[string]time.Time, 100)
	completed := make(map[object]time.Duration, 100)

	// create a bunch of objects and take note of when the request was made
	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("name-%d", i)
		awaitingCompletion[name] = time.Now()
		w.createObject(object{
			name:        name,
			description: fmt.Sprintf("description-%d", i),
			value:       fmt.Sprintf("value-%d", i),
		})
	}

	// continually query for each object until we get them all. take note of how long each one took to create.
	for {
		for k, v := range awaitingCompletion {
			o, err := w.getObject(k)
			if errors.Is(err, errNotFound) {
				panic("whoops")
			}
			if errors.Is(err, errNotReady) {
				continue
			}
			completed[*o] = time.Since(v)
			delete(awaitingCompletion, k)
		}

		if len(awaitingCompletion) == 0 {
			break
		}
	}

	var total int64
	for k, v := range completed {
		fmt.Printf("%v: %s\n\n", k, v)
		total += v.Milliseconds()
	}
	fmt.Printf("average time: %dms", total/100)
}
