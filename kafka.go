package danmu

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/larspensjo/config"
	"goim/src/github.com/wvanbergen/kafka/consumergroup"
	"log"
	"os"
	"strings"
	"time"
)

var (
	Topic    string
	Group    string
	producer sarama.AsyncProducer
	consumer sarama.Consumer
)

func InitKafka(kafkaAddrs []string) error {
	err := InitKafkaProducer(kafkaAddrs)
	if err != nil {
		return err
	}
	Topic = Conf.GetConfig("kafka", "topic")
	Group = Conf.GetConfig("kafka", "group")
	return nil
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

//


