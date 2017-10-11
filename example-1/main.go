package main

import (
	"github.com/zeebe-io/zbc-go/zbc"
	"fmt"
	"errors"
)

const BrokerAddr = "0.0.0.0:51015"

var errClientStartFailed = errors.New("cannot start client")

func main() {
	zbClient, err := zbc.NewClient(BrokerAddr)
	if err != nil {
		panic(errClientStartFailed)
	}

	fmt.Println(zbClient.Cluster.TopicLeaders)
}