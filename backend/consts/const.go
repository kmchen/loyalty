package consts

type RequestType string
type RiderType string
type Grade int

const (
	RIDER_SIGNED_UP            RiderType = "rider_signed_up"
	RIDER_UPDATED_PHONE_NUMBER RiderType = "rider_updated_Phone_number"
	RIDE_CREATED               RiderType = "ride_created"
	RIDE_COMPLETED             RiderType = "ride_completed"

	RIDER_RECORD_REQUEST  RequestType = "RIDER_RECORD_REQUEST"
	RIDER_RECORD_RESPONSE RequestType = "RIDER_RECORD_RESPONSE"

	EXCHANGE string = "events"

	BRONZE Grade = iota
	SILVER
	GOLD
	PLATINUM
)
