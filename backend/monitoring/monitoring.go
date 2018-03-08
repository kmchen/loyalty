package monitoring

import "github.com/prometheus/client_golang/prometheus"

var RideCreated = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "ride_created",
	Help: "Number of ride created",
})

var RideComplete = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "ride_complete",
	Help: "Number of ride completed",
})

var RiderTotal = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "riders_total",
	Help: "Number of riders",
})

var RideFailSignUp = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "ride_signup_fail",
	Help: "Number of signup failure",
})

func init() {
	prometheus.MustRegister(RideCreated)
	prometheus.MustRegister(RideFailSignUp)
	prometheus.MustRegister(RiderTotal)
	prometheus.MustRegister(RideComplete)
}
