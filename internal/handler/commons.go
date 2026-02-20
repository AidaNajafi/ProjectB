package handler

import "time"

func ParseStringDepartureDateToRFC(s string) (time.Time, error){
	return time.Parse(time.RFC3339, s)
}