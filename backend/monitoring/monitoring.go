package monitoring

import "github.com/prometheus/client_golang/prometheus"

var RideCreatedReceived = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "ride_created_received",
	Help: "Number of ride_created received",
})

var RideCreated = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "ride_created",
	Help: "Number of ride_created saved",
})

var RideCompleteReceived = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "ride_complete_received",
	Help: "Number of ride_completed received",
})

var RideComplete = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "ride_complete",
	Help: "Number of ride_completed saved",
})

var RiderTotalReceived = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "riders_total_received",
	Help: "Number of riders received",
})

var RiderTotal = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "riders_total",
	Help: "Number of riders saved",
})

var NumBronzeRider = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "riders_bronze",
	Help: "Number of bronze rider",
})

var NumSilverRider = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "riders_silver",
	Help: "Number of silver rider",
})

var NumGoldRider = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "riders_gold",
	Help: "Number of gold rider",
})

var NumPlatinumRider = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "riders_platinum",
	Help: "Number of platinum rider",
})

var NumLiveUpdate = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "num_live_update",
	Help: "Number of live update",
})

func init() {
	prometheus.MustRegister(RideCreatedReceived)
	prometheus.MustRegister(RideCreated)
	prometheus.MustRegister(RiderTotalReceived)
	prometheus.MustRegister(RiderTotal)
	prometheus.MustRegister(RideCompleteReceived)
	prometheus.MustRegister(RideComplete)
	prometheus.MustRegister(NumBronzeRider)
	prometheus.MustRegister(NumSilverRider)
	prometheus.MustRegister(NumGoldRider)
	prometheus.MustRegister(NumPlatinumRider)
	prometheus.MustRegister(NumLiveUpdate)
}
