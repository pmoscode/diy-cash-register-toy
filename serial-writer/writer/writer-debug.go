package writer

import "log"

type Debug struct{}

func (d *Debug) Write(message string) {
	log.Println("## DEBUG ## --> ", message)
}
