package writer

type Writer interface {
	Write(message string)
	Connect()
	Disconnect()
}
