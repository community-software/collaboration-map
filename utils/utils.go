package utils

import (
	"encoding/json"
	"log"
)

func Log(u interface{}) {
	newMsg, _ := json.Marshal(u)

	log.Printf("\033[1;34m%s\033[0m ", "Message received: ")
	log.Printf("\033[1;34m%s\033[0m ", string(newMsg))
	log.Printf("\033[1;34m%s\033[0m ", "End message")
}
