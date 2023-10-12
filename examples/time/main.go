package main

import (
	"time"

	"github.com/sascha-andres/reuse"
)

func main() {
	jan := time.Date(2023, 01, 01, 01, 01, 01, 01,
		time.FixedZone("UTC+1", 3600))
	println("-- 2023-01-01 01:01:01")
	println(reuse.Date(jan))
	println(reuse.DateTime(jan))
	println(reuse.DateTimeWithTimeZone(jan))
	println("-- now")
	println(reuse.CurrentDate())
	println(reuse.CurrentDateTime())
	println(reuse.CurrentDateTimeWithTimeZone())
	println("-- now UTC")
	println(reuse.UTCCurrentDate())
	println(reuse.UTCCurrentDateTime())
	println(reuse.UTCCurrentDateTimeWithTimeZone())
}
