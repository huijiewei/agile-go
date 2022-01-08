package util

import "time"

func getMonthName() string {
	return time.Now().Format("200601")
}
