module input-reader

go 1.18

replace github.com/pmoscode/golang-mqtt => ../mqtt-client

require (
	github.com/gvalkov/golang-evdev v0.0.0-20191114124502-287e62b94bcb
	github.com/pmoscode/golang-mqtt v0.0.0-00010101000000-000000000000
)

require (
	github.com/eclipse/paho.mqtt.golang v1.3.5 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	golang.org/x/net v0.0.0-20200425230154-ff2c4b7c35a0 // indirect
)
