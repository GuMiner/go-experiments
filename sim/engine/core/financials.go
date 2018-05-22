package core

import (
	"fmt"
	"go-experiments/sim/config"
)

type Financials struct {
	Savings float32
}

func NewFinancials() Financials {
	return Financials{
		Savings: config.Config.Sim.StartingSavings}
}

func (f *Financials) BuyItem(name string, cost float32) {
	f.Savings -= cost
	fmt.Printf("> Purchased a %v for %.0f. Savings: %.0f\n", name, cost, f.Savings)
}

// Updates financials. Returns false if the player is in too much debt to escape
func (f *Financials) Update() bool {
	// TODO: Implement
	return f.Savings > -config.Config.Sim.MaxDebt
}
