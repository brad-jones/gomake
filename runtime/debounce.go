package runtime

import "time"

// Debounce will only execute a function once in the given time frame.
// credit: https://drailing.net/2018/01/debounce-function-for-golang/
func Debounce(interval time.Duration, input chan string, cb func(arg string)) {
	var item string
	timer := time.NewTimer(interval)
	for {
		select {
		case item = <-input:
			timer.Reset(interval)
		case <-timer.C:
			if item != "" {
				cb(item)
			}
		}
	}
}
