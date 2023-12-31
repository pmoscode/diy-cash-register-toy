module serial-writer

go 1.18

replace gitlab.com/pmoscode/golang-shared-libs => ../golang-shared-libs

require (
	github.com/tarm/serial v0.0.0-20180830185346-98f6abe2eb07
	gitlab.com/pmoscode/golang-shared-libs v0.0.0-00010101000000-000000000000
)

require (
	github.com/eclipse/paho.mqtt.golang v1.3.5 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
)
