package riderOps

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"test-fullstack-loyalty/backend/model"
	"test-fullstack-loyalty/backend/monitoring"
	"test-fullstack-loyalty/backend/store"

	"github.com/garyburd/redigo/redis"
)

type RiderOps interface {
	RiderSignUp(data model.Payload) (string, error)
	CreateRide(data model.Payload) (string, error)
	GetRider(key string) (*model.Rider, error)
}
type RiderOperation struct {
	store       store.Store
	riderPusher chan model.Rider
}

func NewRiderOps(store store.Store, riderPusher chan model.Rider) *RiderOperation {
	if store == nil || riderPusher == nil {
		log.Printf("Fail to crate RiderOps due to nil store")
		return nil
	}
	return &RiderOperation{store, riderPusher}
}

// RiderSignUp create rider with key riderId:timestamp:key
func (ops *RiderOperation) RiderSignUp(data model.Payload) (string, error) {
	// key in format riderId:timestamp:user
	key := genKey(data.Id, "user")

	// Check rider whether rider already signed up
	found, err := ops.store.Exist(key)
	if err != nil {
		log.Printf("Fail to lookup rider %v", data.Id)
		return "", err
	}

	// if rider alrady exists the do nothing
	if !found {
		newRider := &model.Rider{
			Id:    data.Id,
			Name:  data.Name,
			Grade: "BRONZE",
		}
		if err := ops.store.Set(key, newRider); err != nil {
			return "", err
		}
		monitoring.RiderTotal.Inc()
	} else {
		log.Printf("Duplicate rider sign up %v", data.Id)
	}
	return key, nil
}

// GetRider get rider info with riderId
func (ops *RiderOperation) GetRider(key string) (*model.Rider, error) {
	userKey, _ := ops.getRiderKey(key)
	value, _ := ops.store.Get(userKey)
	var rider = new(model.Rider)
	if err := redis.ScanStruct(value, rider); err != nil {
		return nil, err
	}
	return rider, nil
}

// createRide creates a ride
func (ops *RiderOperation) CreateRide(data model.Payload) (string, error) {
	var err error
	// key is in the format of riderId:timestamp:created
	key := genKey(data.RiderId, "created")
	if err = ops.newRide(data, key); err != nil {
		log.Printf("[Error] Fail to create a ride, %v", err)
	} else {
		monitoring.RideCreated.Inc()
	}
	return key, ops.newRide(data, key)
}

// CompleteRide completes a ride, update loyalty and push to frontend
func (ops *RiderOperation) CompleteRide(data model.Payload) (string, error) {
	var err error
	var rider *model.Rider

	// key is in the format of riderId:timestamp:complete
	key := genKey(data.RiderId, "complete")
	if err = ops.newRide(data, key); err != nil {
		return "", err
	}

	var keyStr = strconv.FormatInt(data.RiderId, 10)
	if rider, err = ops.updateLoyalty(keyStr, data.Amount); err != nil {
		log.Printf("[Error] Fail to update loyalty %v", err)
		return "", err
	}
	ops.riderPusher <- *rider
	return key, nil
}

// newRide writes a created/complete ride record in redis
func (ops *RiderOperation) newRide(data model.Payload, key string) error {
	ride := &model.Ride{
		Id:      data.Id,
		Amount:  data.Amount,
		RiderId: data.RiderId,
	}
	if err := ops.store.Set(key, ride); err != nil {
		log.Printf("[Error] Fail to create/complete new ride %v, %v", err, data)
		return err
	}
	return nil
}

// updateLoyalty updates loyalty, grade and num of rides
func (ops *RiderOperation) updateLoyalty(key string, amount float32) (*model.Rider, error) {
	rider, err := ops.GetRider(key)
	if err != nil {
		log.Printf("Fail to get Rider with key %s", key)
		return nil, err
	}

	// Loyalty update logic
	rider.NumRides = rider.NumRides + 1
	if rider.NumRides < 20 {
		rider.Loyalty = rider.Loyalty + amount
	} else if rider.NumRides < 50 {
		rider.Loyalty = rider.Loyalty + amount*3
		rider.Grade = "SILVER"
	} else if rider.NumRides < 100 {
		rider.Loyalty = rider.Loyalty + amount*5
		rider.Grade = "GOLD"
	} else {
		rider.Loyalty = rider.Loyalty + amount*10
		rider.Grade = "PLATINUM"
	}

	// Update loyalty
	var riderIdStr = strconv.FormatInt(rider.Id, 10)
	riderKey, _ := ops.getRiderKey(riderIdStr)
	if err := ops.store.Set(riderKey, rider); err != nil {
		return nil, err
	}
	return rider, nil
}

// genKey generates key as riderId:timestamp:suffix
func genKey(key int64, suffix string) string {
	var keyStr = strconv.FormatInt(key, 10)
	var timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	return fmt.Sprintf("%s:%s:%s", keyStr, timestamp, suffix)
}

// genUserKey retrieves user based on riderId
func (ops *RiderOperation) getRiderKey(riderId string) (string, error) {
	var riderIds []string
	var err error
	formatRiderId := fmt.Sprintf("%s:*:%s", riderId, "user")
	if riderIds, err = ops.store.Keys(formatRiderId); err != nil {
		log.Printf("Fail to get key %s, %v", riderIds, err)
	}
	if len(riderIds) != 1 {
		return "", fmt.Errorf("More/less than one rider found  %v", riderIds)
	}
	return riderIds[0], nil
}
