package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var control = &sync.Mutex{}

func getMonitors() []monitor {
	monList := []monitor{
		monitor{name: "Warehouse", device: "/dev/ttyUSB0"},
		monitor{name: "Mechanical Room", device: "/dev/ttyUSB1"},
	}

	return monList
}

func main() {
	ticker := time.NewTicker(1 * time.Hour)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Logging Data")
				control.Lock()

				for _, mon := range getMonitors()[:] {
					logDates(mon)
				}

				control.Unlock()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	http.HandleFunc("/log", forceLog)
	http.HandleFunc("/report", report)
	http.HandleFunc("/", home)
	http.ListenAndServe(":80", nil)

}
