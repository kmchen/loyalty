package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"test-fullstack-loyalty/backend/consts"
	"test-fullstack-loyalty/backend/model"
	"test-fullstack-loyalty/backend/riderOps"
	"test-fullstack-loyalty/backend/store"

	"github.com/garyburd/redigo/redis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/streadway/amqp"
)

var amaqAddr = flag.String("amaqAddr", "amqp://guest:guest@localhost:5672/", "amaq address")
var wsPort = flag.String("wsPort", ":3000", "wesocket port")
var redisPort = flag.String("redisPort", ":6379", "redis port")
var riderChanSize = flag.Int("riderChanSize", 10, "rider channel size")
var serverAddr = flag.String("serverAddr", "localhost:6061", "server address")

var ROUTING_KEY []string = []string{
	"rider.signup",
	"ride.create",
	"ride.completed",
	"rider.phone_update",
}

func listenFromProducer(msgChan <-chan amqp.Delivery, ops *riderOps.RiderOperation) {
	for msg := range msgChan {
		//log.Printf(" [x] %s", d.Body)
		rideData := model.RideData{}
		err := json.Unmarshal(msg.Body, &rideData)
		if err != nil {
			log.Println("Fail to unmarshal data: ", err)
		}
		switch rideData.Type {
		case consts.RIDER_SIGNED_UP:
			ops.RiderSignUp(rideData.Payload)
		case consts.RIDE_CREATED:
			ops.CreateRide(rideData.Payload)
		case consts.RIDE_COMPLETED:
			ops.CompleteRide(rideData.Payload)
		}
	}
}

func main() {

	flag.Parse()

	var err error

	// riderPusher is a chan for rider status live update
	riderPusher := make(chan model.Rider, *riderChanSize)

	// store serves for redis I/O
	store, err := store.NewRedisStore(
		func() (redis.Conn, error) { return redis.Dial("tcp", *redisPort) })
	if err != nil {
		log.Fatalf("Fail to create redis store")
	}

	// ops is a wrapper for rider operation
	ops := riderOps.NewRiderOps(store, riderPusher)

	// proxy is a websocket wrapper
	proxy := newProxy(ops, riderPusher)
	go proxy.Run()

	// Prepare connection to message queue
	// msgChan serves to receive message from message queue
	queueConn, err := NewQueueConn(*amaqAddr)
	if err != nil {
		log.Fatalf("Fail to connect to %s ", *amaqAddr)
	}
	ch, err := QueueChan(queueConn)
	if err != nil {
		log.Fatalf("Fail to create channel to %s ", *amaqAddr)
	}
	msgChan, err := Bind(ch, ROUTING_KEY, consts.EXCHANGE)
	if err != nil {
		log.Fatalf("Fail to connect to rabbit mq @ %s", *amaqAddr)
	}

	// Spawn a goroutine to listen message from producer
	go listenFromProducer(msgChan, ops)

	// Enable metrics monitoring
	http.Handle("/metrics", prometheus.Handler())

	log.Fatal(http.ListenAndServe(*serverAddr, nil))
}
