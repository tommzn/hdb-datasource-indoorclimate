package main

import (
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
)

type IndoorClimateDate struct {
	DeviceId       string `json:"device_id"`
	Characteristic string `json:"characteristic"`
	TimeStamp      int64  `json:"timestamp"`
	Value          string `json:"value"`
}

type IotMessageHandler struct {
	logger log.Logger
	conf   config.Config
}
