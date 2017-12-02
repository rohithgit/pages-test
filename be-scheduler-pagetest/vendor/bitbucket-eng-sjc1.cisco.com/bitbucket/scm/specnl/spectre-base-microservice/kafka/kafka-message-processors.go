package kafkaclient

import (
	"errors"
	"fmt"

	"github.com/Shopify/sarama"
	//"github.com/linkedin/goavro"
)

//Struct for sarama producer, consumer, partition and offset
type MessageProcessor struct {
	ResponseMessenger sarama.SyncProducer
	RequestConsumer   sarama.Consumer
	partition         int32
	offset            int64
}

// function for creation of producer and consumer. Returns reference to struct
func NewMessageProcessor(brokers []string) (*MessageProcessor, error) {
	if len(brokers) == 0 {
		err := errors.New("Invalid broker information provided")
		log.Error(err)
		return nil, err
	}
	producer, err := New_kafka_producer(brokers)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	consumer, err := New_kafka_consumer(brokers)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &MessageProcessor{
		ResponseMessenger: producer,
		RequestConsumer:   consumer,
		offset:            0,
		partition:         0,
	}, nil
}

//fuction to consume message from kafka and returns a channel.
func (ref *MessageProcessor) NewConsumeMsg(topic string) (<-chan string, error) {
	consumer, err := ref.RequestConsumer.ConsumePartition(topic, ref.partition, ref.offset)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	c := make(chan string)

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.Error(err)
			case msg := <-consumer.Messages():
				c <- fmt.Sprintf("%s", string(msg.Value))
				log.Infof("Received messages: %s", string(msg.Value))
			}
		}
	}()

	return c, nil
}

//function to publish message to kafka cluster
func (ref *MessageProcessor) NewPublishMsg(topic, value string) error {

	if value == "" {
		err := errors.New("Publish Message is Empty")
		log.Error(err)
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(value),
	}

	partition, offset, err := ref.ResponseMessenger.SendMessage(msg)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Infof("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}

//function to publish message to kafka cluster which is Avro encoded
func (ref *MessageProcessor) NewPublishMsgAvro(topic string, value []byte) error {

	return ref.NewPublishMsg(topic, string(value))
}
