package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/nsqio/go-nsq"
    "github.com/spf13/viper"
)

type myMessageHandler struct {
    ID string
}

// HandleMessage implements the Handler interface.
func (h *myMessageHandler) HandleMessage(m *nsq.Message) error {
    if len(m.Body) == 0 {
        // Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
        // In this case, a message with an empty body is simply ignored/discarded.
        return nil
    }

    // do whatever actual message processing is desired
    h.ProcessMessage(m.Body)
    return nil
}

func (h *myMessageHandler) ProcessMessage(msg []byte) {
    log.Printf("[CONSUMER %s] Received message %s", h.ID, msg)
}

// InitConfig Function that uses viper library to parse env variables. If
// some of the variables cannot be parsed, an error is returned
func InitConfig() (*viper.Viper, error) {
    v := viper.New()
    v.AutomaticEnv()

    // Add env variables supported
    v.BindEnv("id")
    v.BindEnv("topic")
    v.BindEnv("channel")
    v.BindEnv("lookupd", "servers")

    return v, nil
}

func main() {
    v, err := InitConfig()
    if err != nil {
        log.Fatalf("%s", err)
    }

    // Instantiate a consumer that will subscribe to the provided channel.
    consumer, err := nsq.NewConsumer(
        v.GetString("topic"),
        v.GetString("channel"),
        nsq.NewConfig(),
    )

    if err != nil {
        log.Panic("Could not create consumer")
    }

    // Set the Handler for messages received by this Consumer. Can be called multiple times.
    // See also AddConcurrentHandlers.
    consumer.AddHandler(&myMessageHandler{
        ID: v.GetString("id"),
    })

    // Use nsqlookupd to discover nsqd instances.
    // See also ConnectToNSQD, ConnectToNSQDs, ConnectToNSQLookupds.
    lookupdServers := v.GetStringSlice("lookupd_servers")
    if err := consumer.ConnectToNSQLookupds(lookupdServers); err != nil {
        log.Fatal(err)
    }

    // wait for signal to exit
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan

    // Gracefully stop the consumer.
    consumer.Stop()
}
