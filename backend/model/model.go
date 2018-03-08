package model

import "test-fullstack-loyalty/backend/consts"

type Rider struct {
	Id           int64   `json:"id"`
	Name         string  `json:"name"`
	Loyalty      float32 `json:"loyalty"`
	Phone_number string  `json:"phone_number"`
	NumRides     int     `json:"numRides"`
	Grade        string  `json:"grade"`
}

type Request struct {
	Type   consts.RequestType `json:"type"`
	UserId string             `json:"userId"`
}

type Response struct {
	Type  string `json:"type"`
	Rider Rider  `json:"rider"`
}

type Ride struct {
	Id      int64   `json:"id"`
	Amount  float32 `json:"amount"`
	RiderId int64   `json:"rider_id"`
}

type Payload struct {
	Id           int64   `json:"id, omitempty"`
	Amount       float32 `json:"amount, omitempty"`
	RiderId      int64   `json:"rider_id, omitempty"`
	Phone_number string  `json:"Phone_number, omitempty"`
	Name         string  `json:"name, omitempty"`
}

type RideData struct {
	Type    consts.RiderType `json:"type"`
	Payload Payload          `json:"payload"`
}
