package danmu

import (
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"log"
	"strings"
	"time"
)

var (
	Topic    string
	Topics   []string
	Group    string
	producer sarama.AsyncProducer
	consumer *cluster.Consumer
)

func InitKafka(kafkaAddrs []string) error {
	var (
		err error
	)
	Topic = Conf.GetConfig("kafka", "topic")
	Topics = strings.Split(Topic, ",")
	Group = Conf.GetConfig("kafka", "group")

	err = InitKafkaProducer(kafkaAddrs)
	if err != nil {
		return err
	}
	err = InitKafkaConsumer(kafkaAddrs)
	if err != nil {
		return err
	}

	return nil
}

func CloseKafka(){
	producer.Close()
	consumer.Close()
}

// producer

func InitKafkaProducer(kafkaAddrs []string) error {
	var (
		err error
	)
	config := sarama.NewConfig()
	//config.Producer.RequiredAcks = sarama.NoResponse
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Timeout = 5 * time.Second
	producer, err = sarama.NewAsyncProducer(kafkaAddrs, config)
	go handleProducerSuccess()
	go handleProducerError()
	return err
}

func handleProducerSuccess() {
	var (
		pm *sarama.ProducerMessage
	)
	for {
		pm = <-producer.Successes()
		if pm != nil {
			log.Printf("producer message success, partition:%d offset:%d key:%v valus:%s", pm.Partition, pm.Offset, pm.Key, pm.Value)
		}
	}
}

func handleProducerError() {
	var (
		err *sarama.ProducerError
	)
	for {
		err = <-producer.Errors()
		if err != nil {
			log.Printf("producer message error, partition:%d offset:%d key:%v valus:%s error(%v)", err.Msg.Partition, err.Msg.Offset, err.Msg.Key, err.Msg.Value, err.Err)
		}
	}
}

//consumer

func InitKafkaConsumer(kafkaAddrs []string) error {
	var (
		err error
	)
	config := cluster.NewConfig()
	config.Group.Return.Notifications = false
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.CommitInterval = 1 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetNewest //初始从最新的offset开始

	consumer, err = cluster.NewConsumer(kafkaAddrs, Group, Topics, config)
	if err != nil {
		return err
	}

	// consume errors
	go handleConsumerError()

	return OK
}

func handleConsumerError(){
	for err := range consumer.Errors() {
		log.Printf("Consumer Error: %s\n", err.Error())
	}
}
