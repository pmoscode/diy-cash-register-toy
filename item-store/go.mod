module item-store

go 1.18

replace gitlab.com/pmoscode/golang-shared-libs => ../golang-shared-libs

require (
	gitlab.com/pmoscode/golang-shared-libs v0.0.0-00010101000000-000000000000
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/eclipse/paho.mqtt.golang v1.4.3 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
)
