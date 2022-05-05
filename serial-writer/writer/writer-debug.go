package writer

import "log"

type Debug struct{}

func (d Debug) Connect() {
	log.Println("## DEBUG ## --> Connection successful")
}

func (d *Debug) Write(message string) {
	log.Println("## DEBUG ## --> ", message)
}
