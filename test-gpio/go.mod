module test-gpio

go 1.23.5

replace gitlab.com/pmoscode/golang-shared-libs => ../golang-shared-libs

require (
	gitlab.com/pmoscode/golang-shared-libs v0.0.0-00010101000000-000000000000
	periph.io/x/conn/v3 v3.7.2
	periph.io/x/host/v3 v3.8.5
)
