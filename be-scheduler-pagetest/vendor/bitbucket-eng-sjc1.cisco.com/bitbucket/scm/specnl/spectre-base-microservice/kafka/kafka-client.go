package kafkaclient

import (
	"errors"

	"bitbucket-eng-sjc1.cisco.com/bitbucket/specnl/spectre-base-microservice/logging"

	"github.com/Shopify/sarama"
)

var log = logging.Log.Logger

//function for creating a producer to kafka cluster
func New_kafka_producer(brokers []string) (sarama.SyncProducer, error) {
	if len(brokers) == 0 {
		err := errors.New("Invalid broker information provided")
		log.Error(err)
		return nil, err
	}
	log.Infof("RequestProcessor: new_kafka_producer: start")
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	log.Infof("Producer Config: %v\n", config)
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Errorf("Failed to start Sarama producer: %s", err)
		return nil, err
	}
	log.Infof("RequestProcessor: new_kafka_producer: end")
	return producer, nil
}

//function for creating a consumer to kafka cluster
func New_kafka_consumer(brokers []string) (sarama.Consumer, error) {
	if len(brokers) == 0 {
		err := errors.New("Invalid broker information provided")
		log.Error(err)
		return nil, err
	}
	log.Infof("RequestProcessor: new_kafka_consumer: start")
	config := sarama.NewConfig()
	log.Infof("Consumer Config: %v\n", config)
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Errorf("Failed to start Sarama consumer: %s", err)
		return nil, err
	}
	log.Infof("RequestProcessor: new_kafka_consumer: end")
	return consumer, nil
}
