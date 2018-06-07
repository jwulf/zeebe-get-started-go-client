package main

import (
	"errors"
	"fmt"
	"github.com/zeebe-io/zbc-go/zbc"
	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
	"github.com/zeebe-io/zbc-go/zbc/services/zbsubscribe"
	"os"
	"os/signal"
)

const topicName = "default-topic"
const brokerAddr = "0.0.0.0:51015"

var errClientStartFailed = errors.New("cannot start client")

func main() {
	zbClient, err := zbc.NewClient(brokerAddr)
	if err != nil {
		panic(err)
	}

	// deploy workflow
	response, err := zbClient.CreateWorkflowFromFile(topicName, zbcommon.BpmnXml, "order-process.bpmn")
	if err != nil {
		panic(err)
	}

	fmt.Println(response.String())

	// create a new workflow instance
	payload := make(map[string]interface{})
	payload["orderId"] = "31243"

	instance := zbc.NewWorkflowInstance("order-process", -1, payload)
	msg, err := zbClient.CreateWorkflowInstance(topicName, instance)

	if err != nil {
		panic(err)
	}

	fmt.Println(msg.String())

	subscription, err := zbClient.JobSubscription(topicName, "sample-app", "payment-service", 1000, 32, func(client zbsubscribe.ZeebeAPI, event *zbsubscriptions.SubscriptionEvent) {
		fmt.Println(event.String())

		// complete job after processing
		response, _ := client.CompleteJob(event)
		fmt.Println(response)
	})

	if err != nil {
		panic("Unable to open subscription")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		err := subscription.Close()
		if err != nil {
			panic("Failed to close subscription")
		}

		fmt.Println("Closed subscription")
		os.Exit(0)
	}()

	subscription.Start()
}
