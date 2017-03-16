package logrunner

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func roundTime(t time.Time, interval time.Duration) time.Time {

	roundedMinutes := int(interval / (1 * time.Minute))
	if roundedMinutes < 60 {
		//Round to minutes e.g. 5min, 30 mins etc.
		minuteInHour := t.Minute()
		minute := (minuteInHour / roundedMinutes) * roundedMinutes
		tr := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), int(minute), 0, 0, time.Local)
		return tr
	}
	//Round to some multiple of an hour e.g. 1hr, 2hr, 24hr
	roundedHours := roundedMinutes / 60
	hour := (t.Hour() / roundedHours) * roundedHours
	tr := time.Date(t.Year(), t.Month(), t.Day(), hour, 0, 0, 0, time.Local)
	return tr
}

func updateTicker(interval time.Duration) *time.Ticker {

	nextTick := roundTime(time.Now(), interval).Add(interval)
	diff := nextTick.Sub(time.Now())
	fmt.Printf("NextTick at: %s Diff: %s\n", nextTick.String(), diff.String())
	return time.NewTicker(diff)
}
func logName(filename string, t time.Time, interval time.Duration) string {
	directory := filepath.Dir(filename)
	base := filepath.Base(filename)
	ext := filepath.Ext(base)
	if ext != "" {
		base = strings.TrimSuffix(base, ext)
	}
	//                  YYYY-MM-DD-HH-MM
	suffix := t.Format("2006-01-02-15-04")
	fmt.Printf("Directory: %s Base: %s\n", directory, base)
	fmt.Printf("Suffix: %s\n", suffix)
	newFilename := path.Join(directory, base+"-"+suffix+ext)
	return newFilename
}

func setLogOutput(filename string, t time.Time, interval time.Duration) *os.File {
	nextTime := roundTime(t, interval)
	newFilename := logName(filename, nextTime, interval)
	fmt.Printf("newFilename: %s\n", newFilename)
	f, err := os.OpenFile(newFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file - %s err: %s", newFilename, err)
	}
	//defer f.Close()
	log.SetOutput(f)
	return f
}
func LogRunner(filename string, interval time.Duration) {

	ticker := updateTicker(interval)
	f := setLogOutput(filename, time.Now().UTC(), interval)
	for {
		<-ticker.C
		fmt.Println(time.Now(), "- just ticked - Roll over the log files")
		ticker = updateTicker(interval)
		f.Close()
		f = setLogOutput(filename, time.Now().UTC(), interval)
	}
}
