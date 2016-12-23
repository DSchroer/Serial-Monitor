package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type logValue struct {
	timeStamp time.Time
	value     int
}

type byDate []logValue

func (a byDate) Len() int           { return len(a) }
func (a byDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byDate) Less(i, j int) bool { return !a[i].timeStamp.Before(a[j].timeStamp) }

func openOrCreateFile(file string) *os.File {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			newFile, err := os.Create(file)
			if err != nil {
				log.Print(err)
				return nil
			}
			return newFile
		}
		log.Print(err)
		return nil
	}

	inFile, err := os.OpenFile(file, os.O_APPEND|os.O_RDWR, os.ModeAppend)
	if err != nil {
		log.Print(err)
		return nil
	}
	return inFile
}

func getTimestamp() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func getValues(mon monitor) []logValue {
	fileName := mon.name + ".txt"

	file := openOrCreateFile(fileName)
	defer file.Close()

	var logValues []logValue
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var text []string
		text = strings.Split(scanner.Text(), ":")

		if len(text) != 2 {
			log.Print("Data was in invalid format.")
			continue
		}

		unixStamp, err := strconv.ParseInt(text[0], 10, 64)
		if err != nil {
			log.Print(err)
			continue
		}
		timeStamp := time.Unix(unixStamp, 0)

		value, err := strconv.Atoi(text[1])
		if err != nil {
			log.Print(err)
			continue
		}

		nextValue := logValue{timeStamp: timeStamp, value: value}

		logValues = append(logValues, nextValue)
	}

	return logValues
}

func logDates(mon monitor) {

	val, err := mon.dailyUsage()
	if err != nil {
		log.Print(err)
		return
	}

	var newData []logValue
	for index, value := range val[:] {
		date := getTimestamp().AddDate(0, 0, -index)
		newData = append(newData, logValue{timeStamp: date, value: value})
	}

	oldFileName := mon.name + ".txt"
	newFileName := mon.name + ".txt.tmp"

	inFile := openOrCreateFile(oldFileName)

	err = os.Remove(newFileName)
	if err != nil {
		log.Print(err)
	}

	outFile := openOrCreateFile(newFileName)

	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		var text []string
		text = strings.Split(scanner.Text(), ":")

		if len(text) != 2 {
			log.Print("Data was in invalid format.")
			continue
		}

		unixStamp, err := strconv.ParseInt(text[0], 10, 64)
		if err != nil {
			log.Print(err)
			continue
		}
		timeStamp := time.Unix(unixStamp, 0)

		value, err := strconv.Atoi(text[1])
		if err != nil {
			log.Print(err)
			continue
		}

		nextValue := logValue{timeStamp: timeStamp, value: value}

		if nextValue.timeStamp.Before(time.Now().AddDate(0, -3, 0)) {
			continue
		}

		tmpNewData := newData
		for index, value := range newData[:] {
			if nextValue.timeStamp.Equal(value.timeStamp) {
				nextValue.value = value.value
				tmpNewData = append(newData[:index], newData[index+1:]...)
			}
		}
		newData = tmpNewData

		_, err = outFile.WriteString(fmt.Sprintf("%d:%d\n", nextValue.timeStamp.Unix(), nextValue.value))
		if err != nil {
			log.Print(err)
			continue
		}
	}

	for _, value := range newData[:] {
		_, err = outFile.WriteString(fmt.Sprintf("%d:%d\n", value.timeStamp.Unix(), value.value))
		if err != nil {
			log.Print(err)
			continue
		}
	}

	inFile.Close()
	outFile.Close()

	err = os.Remove(oldFileName)
	if err != nil {
		log.Print(err)
		return
	}

	err = os.Rename(newFileName, oldFileName)
	if err != nil {
		log.Print(err)
		return
	}
}
