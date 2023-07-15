module input-reader

go 1.18

replace gitlab.com/pmoscode/golang-shared-libs => ../golang-shared-libs

require github.com/gvalkov/golang-evdev v0.0.0-20220815104727-7e27d6ce89b6

require (
	github.com/eclipse/paho.mqtt.golang v1.3.5 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	gitlab.com/pmoscode/golang-shared-libs v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.7.0 // indirect
)
