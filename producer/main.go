package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/spf13/viper"
)

// InitConfig Function that uses viper library to parse env variables. If
// some of the variables cannot be parsed, an error is returned
func InitConfig() (*viper.Viper, error) {
	v := viper.New()
	v.AutomaticEnv()

	// Add env variables supported
	v.BindEnv("id")
	v.BindEnv("topic")
	v.BindEnv("channel")
	v.BindEnv("nsqd", "servers")
	v.BindEnv("loop", "interval")

	return v, nil
}

func main() {
	v, err := InitConfig()
	if err != nil {
		log.Fatalf("%s", err)
	}

	// Instantiate a consumer that will subscribe to the provided channel.

	nsqConfig := nsq.NewConfig()
	listOfProducers := make([]*nsq.Producer, 0)
	for _, producerAddr := range v.GetStringSlice("nsqd_servers") {
		producer, err := nsq.NewProducer(
			producerAddr,
			nsqConfig,
		)
		if err != nil {
			log.Fatal("MAMAAAA!!")
		}
		listOfProducers = append(listOfProducers, producer)
	}

	if err != nil {
		log.Panic("Could not create producer")
	}

	topicName := v.GetString("topic")
	messageBody := []byte(fmt.Sprintf("[PRODUCER %s] Sending message", v.GetString("id")))
	// channelName := v.GetString("channel")

	for {
		messageSent := false
		for _, producer := range listOfProducers {
			log.Printf("[PRODUCER %s] Sending message to nsqd server %s", v.GetString("id"), producer)
			if err := producer.Publish(topicName, messageBody); err != nil {
				continue
			}
			messageSent = true
			break
		}

		if !messageSent {
			log.Fatal("Esto no es fault tolerant vieja!!")
		}

		log.Printf("%s", messageBody)
		time.Sleep(v.GetDuration("loop_interval"))
	}
}
