package domain

import (
	"time"
)

type Weather struct{
	City string
	Time time.Time
	Temperature float32
	FeelsLike float32
	Humidity int32
	WindSpeed float32
}
