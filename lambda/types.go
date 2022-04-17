package main

import (
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	indoorclimate "github.com/tommzn/hdb-datasource-indoorclimate"
)

type IndoorClimateData struct {
	DeviceId       string `json:"device_id"`
	Characteristic string `json:"characteristic"`
	TimeStamp      int64  `json:"timestamp"`
	Value          string `json:"value"`
}

type IotMessageHandler struct {
	logger    log.Logger
	conf      config.Config
	publisher []indoorclimate.Publisher
}
