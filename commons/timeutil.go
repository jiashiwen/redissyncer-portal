package commons

import "time"

func GetCurrentUnixMillisecond() int64 {
	return time.Now().UnixNano() / 1e6
}
