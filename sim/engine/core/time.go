package core

import (
	"fmt"
	"go-experiments/sim/config"
)

type Time struct {
	SimTime float32
	DayTime float32
	Days    int
}

func NewTime() Time {
	return Time{
		SimTime: 0,
		DayTime: 0,
		Days:    0}
}

// Updates the time. Returns true if we passed a day marker.
func (t *Time) Update(step float32) bool {
	t.SimTime += step
	t.DayTime += step

	if t.DayTime > config.Config.Sim.SecondsPerDay {
		t.DayTime -= config.Config.Sim.SecondsPerDay
		t.Days++

		fmt.Printf(">>> Advanced to day %v <<<\n", t.Days)
		return true
	}

	return false
}
