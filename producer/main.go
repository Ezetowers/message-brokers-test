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
	v.BindEnv("nsqd", "server")
	v.BindEnv("loop", "interval")

	return v, nil
}

func main() {
	v, err := InitConfig()
	if err != nil {
		log.Fatalf("%s", err)
	}

	// Instantiate a consumer that will subscribe to the provided channel.
	producer, err := nsq.NewProducer(
		v.GetString("nsqd_server"),
		nsq.NewConfig(),
	)

	if err != nil {
		log.Panic("Could not create producer")
	}

	topicName := v.GetString("topic")
	messageBody := []byte(fmt.Sprintf("[PRODUCER %s] Sending message", v.GetString("id")))
	// channelName := v.GetString("channel")

	for {
		if err := producer.Publish(topicName, messageBody); err != nil {
			log.Fatal(err)
		}
		log.Printf("%s", messageBody)
		time.Sleep(v.GetDuration("loop_interval"))
	}

	// Gracefully stop the producer when appropriate (e.g. before shutting down the service)
	producer.Stop()
}
