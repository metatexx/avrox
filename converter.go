package avrox

import "time"

// AvroTime truncates a go time.Time to the value that gets stored the avro logicalTime
// It also makes sure that the time is expressed in UTC()
func AvroTime(t time.Time) time.Time {
	return t.UTC().Truncate(time.Millisecond)
}

// AvroDate truncates a go time.Time to the value that gets stored the avro logicalDate
// It also makes sure that the time is expressed in UTC()
func AvroDate(t time.Time) time.Time {
	return t.UTC().Truncate(time.Hour * 24)
}
