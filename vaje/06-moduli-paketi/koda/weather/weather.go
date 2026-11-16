// Package weather provides a simulated weather station that generates random weather data.
//
// The package implements a concurrent weather station that runs three independent sensors
// (temperature, humidity, and pressure) as goroutines. Each sensor generates random
// measurements at configurable intervals and sends them through a shared channel.
//
// Example usage:
//
//	station := weather.NewStation(1*time.Second, 5*time.Second)
//	defer station.Stop()
//
//	data, err := station.GetData()
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Type: %s, Value: %.2f\n", data.MType, data.Value)
//
// The weather station generates the following measurements:
//   - Temperature: -5-30°C
//   - Humidity: 40-80%
//   - Pressure: 950-1010 mbar
package weather

import (
	"errors"
	"math/rand"
	"time"
)

// ErrTimeout is returned by GetData when no measurement arrives within the timeout.
var ErrTimeout = errors.New("weather station is not responding")

// WeatherData represents a single weather measurement.
type WeatherData struct {
	MType string
	Value float32
}

// Station represents a weather station with multiple sensors.
type Station struct {
	data            chan WeatherData
	stopCh          chan struct{}
	pollingInterval time.Duration
	timeout         time.Duration
}

// NewStation creates a new weather station with the specified polling interval and timeout.
func NewStation(pollingInterval time.Duration, timeout time.Duration) *Station {
	s := &Station{
		data:   make(chan WeatherData, 3),
		stopCh: make(chan struct{}),
	}
	s.pollingInterval = pollingInterval
	s.timeout = timeout
	go s.runSensor("Temperature", -5, 35)
	go s.runSensor("Humidity", 40, 40)
	go s.runSensor("Pressure", 950, 60)

	return s
}

// Stop stops the weather station and all its sensors.
func (s *Station) Stop() {
	close(s.stopCh)
}

// runSensor simulates a sensor of the given type that generates random
// values in the interval [min, min+span).
func (s *Station) runSensor(mType string, min, span float32) {
	for {
		m := WeatherData{mType, min + rand.Float32()*span}
		select {
		case s.data <- m:
		case <-s.stopCh:
			return
		}
		time.Sleep(s.pollingInterval)
	}
}

// GetData retrieves a weather measurement from the station, respecting the timeout.
func (s *Station) GetData() (WeatherData, error) {
	select {
	case m := <-s.data:
		return m, nil
	case <-time.After(s.timeout):
		return WeatherData{}, ErrTimeout
	}
}
