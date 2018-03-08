package riderOps

import (
	"fmt"
	"strconv"
	"testing"

	"test-fullstack-loyalty/backend/model"
	"test-fullstack-loyalty/backend/store"

	"github.com/garyburd/redigo/redis"
)

var testStore, _ = store.NewRedisStore(
	func() (redis.Conn, error) { return redis.Dial("tcp", ":6379") })

var ops = NewRiderOps(testStore, make(chan model.Rider, 10))

func TestRiderOperation(t *testing.T) {
	var riderId int64 = 12345
	var riderIdStr = strconv.FormatInt(riderId, 10)
	var riderName = "testName"
	var rider *model.Rider
	var err error

	// Test sign up a new rider
	newRider := model.Payload{Id: riderId, Name: riderName}
	if _, err := ops.RiderSignUp(newRider); err != nil {
		t.Errorf("[Fail] Fail to riderSignUp %v", err)
	}
	if rider, err = ops.GetRider(riderIdStr); err != nil {
		t.Errorf("[Fail] Fail to getRider %v with key %s", err, riderIdStr)
	}
	if rider.Id != riderId || rider.Name != riderName || rider.Grade != "BRONZE" {
		t.Errorf("[Fail] Wrong data retrieved with key %s, %+v", riderIdStr, rider)
	}

	var rideComplete = model.Payload{
		Id:      7878,
		Amount:  2,
		RiderId: riderId,
	}
	// Test loyalty logic
	for i := 1; i <= 100; i++ {
		if _, err = ops.CompleteRide(rideComplete); err != nil {
			t.Errorf("[Fail] Fail to getRider %v with key %s", err, riderIdStr)
		}
		if rider, err = ops.GetRider(riderIdStr); err != nil {
			t.Errorf("[Fail] Fail to getRider %v with key %s", err, riderIdStr)
		}
		<-ops.riderPusher

		if i == 1 {
			if rider.Loyalty != 2 || rider.NumRides != 1 || rider.Grade != "BRONZE" {
				t.Errorf("[Fail] invalid loyalty computation %v or num of rides %v or grade %v", rider.Loyalty, rider.NumRides, rider.Grade)
			}
		}
		if i == 20 {
			if rider.Loyalty != 44 || rider.NumRides != 20 || rider.Grade != "SILVER" {
				t.Errorf("[Fail] invalid loyalty computation %v or num of rides %v or grade %v", rider.Loyalty, rider.NumRides, rider.Grade)
			}
		}
		if i == 50 {
			if rider.Loyalty != 228 || rider.NumRides != 50 || rider.Grade != "GOLD" {
				t.Errorf("[Fail] invalid loyalty computation %v or num of rides %v or grade %v", rider.Loyalty, rider.NumRides, rider.Grade)
			}
		}
		if i == 100 {
			if rider.Loyalty != 738 || rider.NumRides != 100 || rider.Grade != "PLATINUM" {
				t.Errorf("[Fail] invalid loyalty computation %v or num of rides %v or grade %v", rider.Loyalty, rider.NumRides, rider.Grade)
			}
		}
	}

	tearDown(riderIdStr)

}

func tearDown(riderIdStr string) {
	formatKey := fmt.Sprintf("%s:*:%s", riderIdStr, "user")
	testStore.Del(formatKey)

	formatKey = fmt.Sprintf("%s:*:%s", riderIdStr, "complete")
	testStore.Del(formatKey)
}
