package main

type Well struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	DepthM       int     `json:"depth_m"`
	PressureKPa  float64 `json:"pressure_kpa"`
	TemperatureC float64 `json:"temperature_c"`
	Status       string  `json:"status"`
}
