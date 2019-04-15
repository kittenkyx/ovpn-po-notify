package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"io/ioutil"
	"time"
	"encoding/json"
	"github.com/gregdel/pushover"
	"github.com/hpcloud/tail"
)

type ConfigData struct {
	User string
	App string
	Location string
}

var (
	// Mon Apr 15 02:43:05 2019 x.x.x.x:port [kyxlaptop] Peer Connection Initiated with [AF_INET]y.y.y.y:port2
	r = regexp.MustCompile(`(?:^\w+\s+\w+\s+\d+\s+\d{2}:\d{2}:\d{2}\s\d+)\s(\d+\.\d+\.\d+.\d+):(\d+)\s\[(\w+)\]\sPeer Connection Initiated with.+`)
	// matches[1] -> ip, matches[2] -> port, matches[3] => device name
	userToken string
	appToken string
	logFile string
)

func main() {
	var config ConfigData

    jsonData, err := ioutil.ReadFile("config.json")
    if err != nil {
        fmt.Println("File reading error", err)
    	return
	}

	err = json.Unmarshal(jsonData, &config)
	if err != nil {
	    log.Println(err)
	}

	userToken= config.User
	appToken = config.App
	logFile = config.Location

	app := pushover.New(appToken)
	recipient := pushover.NewRecipient(userToken)

	tailConfig := tail.Config{
		Follow:    true,
		ReOpen:    true,
		MustExist: true,
		Location:  &tail.SeekInfo{0, os.SEEK_END},
		Poll:      true,
	}

	t, err := tail.TailFile(logFile, tailConfig)
	
	if err != nil {
		log.Fatal(err)
	}

	for line := range t.Lines {
		matches := r.FindStringSubmatch(line.Text)
		if matches != nil {	
			t := time.Now()
			dt := t.Format("2006-01-02 15:04:05 MST")
			msg := pushover.NewMessageWithTitle(dt, fmt.Sprintf("OVPN: %s (%s)", matches[3], matches[1]))
			log.Println(msg)
			response, err := app.SendMessage(msg, recipient)
			if err != nil {
				log.Panic(err)
			}
    		log.Println(response)
		}
	}
}