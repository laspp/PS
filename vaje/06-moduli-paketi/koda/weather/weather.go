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
	go s.runTemperatureSensor()
	go s.runHumiditySensor()
	go s.runPressureSensor()

	return s
}

// Stop stops the weather station and all its sensors.
func (s *Station) Stop() {
	close(s.stopCh)
}

// runTemperatureSensor simulates the temperature sensor.
func (s *Station) runTemperatureSensor() {
	for {
		select {
		case <-s.stopCh:
			return
		default:
			value := -5 + rand.Float32()*35
			m := WeatherData{"Temperature", value}
			s.data <- m
			time.Sleep(s.pollingInterval)
		}
	}
}

// runHumiditySensor simulates the humidity sensor.
func (s *Station) runHumiditySensor() {
	for {
		select {
		case <-s.stopCh:
			return
		default:
			m := WeatherData{"Humidity", 40 + rand.Float32()*40}
			s.data <- m
			time.Sleep(s.pollingInterval)
		}
	}
}

// runPressureSensor simulates the pressure sensor.
func (s *Station) runPressureSensor() {
	for {
		select {
		case <-s.stopCh:
			return
		default:
			m := WeatherData{"Pressure", 950 + rand.Float32()*60}
			s.data <- m
			time.Sleep(s.pollingInterval)
		}
	}
}

// GetData retrieves a weather measurement from the station, respecting the timeout.
func (s *Station) GetData() (WeatherData, error) {
	select {
	case m := <-s.data:
		return m, nil
	case <-time.After(s.timeout):
		return WeatherData{}, errors.New("napaka: vremenska postaja se ne odziva")
	}
}
