package genetics

import (
	"fmt"
	"sort"
)

type Agents []*Agent

func (a Agents) Len() int           { return len(a) }
func (a Agents) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Agents) Less(i, j int) bool { return a[i].GetFinalScore() < a[j].GetFinalScore() }

type Population struct {
	agents Agents

	generationCount           int
	currentGenerationLifetime float32
}

const maxGenerationLifetime = 30.0 // seconds

const mutationProbability = 0.20
const mutationAmount = 1.5

const selectionPercent = 0.04

const crossoverProbability = 0.20

func NewPopulation(agentCount int, agentCreator func(int) *Agent) *Population {
	population := Population{
		generationCount: 0,
		agents:          make([]*Agent, agentCount)}

	for i := 0; i < agentCount; i++ {
		population.agents[i] = agentCreator(i)
	}

	population.prepareNewGeneration()

	return &population
}

func allAgentsDead(agents Agents) bool {
	for _, agent := range agents {
		if agent.isAlive {
			return false
		}
	}

	return true
}

func recombine(agents Agents) {
	bestAgentCount := int(float32(len(agents)) * selectionPercent)

	agentIdx := bestAgentCount
	for true {
		for i := 0; i < bestAgentCount-1; i++ {
			for j := i + 1; j < bestAgentCount; j++ {
				if agentIdx == len(agents) {
					return
				}

				agents[agentIdx].CrossBreed(agents[i], agents[j], crossoverProbability)
				agentIdx++
			}
		}
	}
}

func mutate(agents Agents) {
	for i, agent := range agents {
		if i > len(agents)/4 {
			agent.net.ProbablyRandomize(mutationProbability, mutationAmount)
		}
	}
}

func agentAliveCount(agents Agents) int {
	alive := 0
	for _, agent := range agents {
		if agent.isAlive {
			alive++
		}
	}

	return alive
}

func (p *Population) prepareNewGeneration() {
	p.generationCount++
	p.currentGenerationLifetime = 0.0

	fmt.Printf("==Generation: %v\n==", p.generationCount)
	lastFrameDivisor = 0
}

var lastFrameDivisor int = 0

func (p *Population) Update(frameTime float32, agentUpdater func(*Agent)) {
	for _, agent := range p.agents {
		agentUpdater(agent)
	}

	p.currentGenerationLifetime += frameTime
	if p.currentGenerationLifetime > maxGenerationLifetime || allAgentsDead(p.agents) {

		// Create a new generation by sorting, creating (in-place) new agents, and mutating them
		sort.Sort(sort.Reverse(p.agents))
		fmt.Printf("High Score: %.2f\n", p.agents[0].GetFinalScore())

		recombine(p.agents)
		mutate(p.agents)

		p.prepareNewGeneration()
	} else {
		if int(p.currentGenerationLifetime/5) != lastFrameDivisor {
			fmt.Printf("  %v seconds into generation %v (%v agents left)\n", int(p.currentGenerationLifetime), p.generationCount, agentAliveCount(p.agents))
			lastFrameDivisor = int(p.currentGenerationLifetime / 5)
		}
	}
}

func (p *Population) Render(agentRenderer func(*Agent)) {
	for _, agent := range p.agents {
		agentRenderer(agent)
	}
}
