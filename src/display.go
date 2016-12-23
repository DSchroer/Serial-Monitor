package main

import (
	"fmt"
	"net/http"
	"sort"
)

func forceLog(w http.ResponseWriter, r *http.Request) {
	control.Lock()

	for _, mon := range getMonitors()[:] {
		logDates(mon)
	}

	fmt.Fprintf(w, "<head></head><body><h1>Reading Taken</h1></body>")

	control.Unlock()
}

func home(w http.ResponseWriter, r *http.Request) {
	control.Lock()
	fmt.Fprintf(w, "<head></head><body>")

	for _, mon := range getMonitors()[:] {
		fmt.Fprintf(w, "<div style='float:left;margin:25px;'><h1>%s</h1>", mon.name)
		current, err := mon.currentDraw()
		if err != nil {
			fmt.Fprintf(w, "<p>Read Error (Refresh Page).</p>")
		} else {
			fmt.Fprintf(w, "<p>Current Reading(Watts): %d</p>", current)
		}

		fmt.Fprintf(w, "<table border='1'><tr><th>Date</th><th>Usage (KWHR)</ht></tr>")
		week, _ := mon.dailyUsage()
		for index, day := range week[:] {
			time := getTimestamp().AddDate(0, 0, -index)
			fmt.Fprintf(w, "<td>%d-%d-%d</td><td>%d</td>", time.Year(), time.Month(), time.Day(), day)
			fmt.Fprintf(w, "</tr>")
		}
		fmt.Fprintf(w, "</table><br/></div>")
	}

	fmt.Fprintf(w, "<form action='/report'><input type='submit' value='Display Report'></form>")
	fmt.Fprintf(w, "<form action='/log'><input type='submit' value='Take Reading'></form>")
	fmt.Fprintf(w, "<form action='/'><input type='submit' value='Refresh Page'></form>")

	fmt.Fprintf(w, "<script>"+
		"var time = new Date().getTime();"+
		"function refresh() {"+
		"     if(new Date().getTime() - time >= 60000)"+
		"         window.location.reload(true);"+
		"     else"+
		"         setTimeout(refresh, 10000);"+
		"}"+
		"setTimeout(refresh, 10000);"+
		"</script>")

	fmt.Fprintf(w, "</body>")

	control.Unlock()
}

func report(w http.ResponseWriter, r *http.Request) {
	control.Lock()
	fmt.Fprintf(w, "<head></head><body>")

	for _, mon := range getMonitors()[:] {
		fmt.Fprintf(w, "<div style='float:left;margin:25px;'><h1>%s</h1>", mon.name)
		dataPoints := getValues(mon)

		sort.Sort(byDate(dataPoints))

		var currentTotal, lastTotal, absoluteTotal int
		for _, point := range dataPoints[:] {

			timeStamp := getTimestamp()
			firstCutoff := timeStamp.AddDate(0, 0, -(timeStamp.Day() - 1))
			secondCutoff := firstCutoff.AddDate(0, -1, 0)

			time := point.timeStamp
			if !time.Before(firstCutoff) {
				currentTotal += point.value
			} else if !time.Before(secondCutoff) {
				lastTotal += point.value
			}

			absoluteTotal += point.value
		}
		fmt.Fprintf(w, "<p>Current Month Total (KWHR): %d</p>", currentTotal)
		fmt.Fprintf(w, "<p>Last Month Total (KWHR): %d</p>", lastTotal)
		fmt.Fprintf(w, "<p>All Time Total (KWHR): %d</p>", absoluteTotal)

		fmt.Fprintf(w, "<table border='1'><tr><th>Date</th><th>Usage (KWHR)</ht></tr>")
		for _, point := range dataPoints[:] {
			fmt.Fprintf(w, "<tr>")
			time := point.timeStamp
			fmt.Fprintf(w, "<td>%d-%d-%d</td><td>%d</td>", time.Year(), time.Month(), time.Day(), point.value)
			fmt.Fprintf(w, "</tr>")
		}
		fmt.Fprintf(w, "</table></div>")
	}

	fmt.Fprintf(w, "</body>")
	control.Unlock()
}
