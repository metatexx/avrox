package avrox

import "time"

// AvroTime truncates a go time.Time to the value that gets stored the avro logicalTime (which has a granularity of milliseconds while go has nanoseconds)
// It also makes sure that the time is expressed in UTC()
func AvroTime(t time.Time) time.Time {
	return t.Truncate(time.Millisecond).UTC()
}

// AvroDate truncates a go time.Time to the value that gets stored the avro logicalDate
// It also makes sure that the time is expressed in UTC()
func AvroDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
