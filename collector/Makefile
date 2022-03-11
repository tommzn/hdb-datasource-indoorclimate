build:
	echo "Build IndoorClimate Collector Binary"
	go build -o sensordatacollector

stop_deamon:
	echo "Stopping IndoorClimate Collector"
	sudo systemctl stop  sensordatacollector.service 

start_deamon:
	echo "Starting IndoorClimate Collector"
	sudo systemctl start  sensordatacollector.service 

install:
	echo "Copy Binary"
	cp sensordatacollector ~/go/bin/sensordatacollector 

deploy: build stop_deamon install start_deamon
