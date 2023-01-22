package broker

import (
	"sync"

	"github.com/jerry0420/queue-system/backend/logging"
)

type Event map[string]interface{}

type Broker struct {
	consumers map[string]map[chan Event]bool
	logger    logging.LoggerTool
	sync.Mutex
}

func NewBroker(logger logging.LoggerTool) *Broker {
	return &Broker{
		consumers: make(map[string]map[chan Event]bool),
		logger:    logger,
	}
}

func (broker *Broker) Subscribe(topicName string) chan Event {
	broker.Lock()
	defer broker.Unlock()

	if topicName == "" {
		broker.logger.FATALf("topic name can't be empty!")
	}

	consumerChan := make(chan Event)
	consumerChans, ok := broker.consumers[topicName]
	if !ok {
		broker.consumers[topicName] = map[chan Event]bool{consumerChan: true}
	} else {
		consumerChans[consumerChan] = true
	}
	return consumerChan
}

func (broker *Broker) UnsubscribeConsumer(topicName string, consumerChan chan Event) {
	broker.Lock()
	defer broker.Unlock()
	close(consumerChan)
	delete(broker.consumers[topicName], consumerChan)
}

func (broker *Broker) UnsubscribeTopic(topicName string) {
	broker.Lock()
	defer broker.Unlock()
	for consumerChan, _ := range broker.consumers[topicName] {
		close(consumerChan)	
	}
	delete(broker.consumers, topicName)
}

func (broker *Broker) CloseAll() {
	broker.Lock()
	defer broker.Unlock()
	for _, consumerChans := range broker.consumers {
		for consumerChan, _ := range consumerChans {
			close(consumerChan)
		}
	}
}

func (broker *Broker) Publish(topicName string, event Event) {
	broker.Lock()
	defer broker.Unlock()
	for consumerChan, _ := range broker.consumers[topicName] {
		// broadcasting...
		consumerChan <- event
	}
}