package genetics

import (
	"fmt"
	"go-experiments/voxelli/config"
	"math"
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

func NewPopulation(agentCount int, agentCreator func(int) *Agent) *Population {
	population := Population{
		generationCount: 0,
		agents:          make([]*Agent, agentCount)}

	for i := 0; i < agentCount; i++ {
		population.agents[i] = agentCreator(i)

		// We only save the best run, mutation and all should trickle down fairly fast.
		if i == 0 {
			population.agents[i].LoadNet()
		}
	}

	population.prepareNewGeneration()

	return &population
}

func recombine(agents Agents) {
	bestAgentCount := int(float32(len(agents)) * config.Config.Simulation.Evolver.SelectionPercent)

	agentIdx := bestAgentCount
	for true {
		for i := 0; i < bestAgentCount-1; i++ {
			for j := i + 1; j < bestAgentCount; j++ {
				if agentIdx == len(agents) {
					return
				}

				agents[agentIdx].CrossBreed(agents[i], agents[j], config.Config.Simulation.Evolver.CrossoverProbability)
				agentIdx++
			}
		}
	}
}

func mutate(agents Agents) {
	for i, agent := range agents {
		if i > len(agents)/4 {
			agent.net.ProbablyRandomize(
				config.Config.Simulation.Evolver.MutationProbability,
				config.Config.Simulation.Evolver.MutationAmount)
		}
	}
}

func allAgentsDead(agents Agents) bool {
	for _, agent := range agents {
		if agent.isAlive {
			return false
		}
	}

	return true
}

func allAgentsStopped(agents Agents) bool {
	for _, agent := range agents {
		if math.Abs(float64(agent.car.Velocity)) >
			float64(config.Config.Simulation.Agent.WiggleSpeed*100) {
			return false
		}
	}

	return true
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

	for _, agent := range p.agents {
		agent.Reset()
	}

	fmt.Printf("==Generation: %v ==\n", p.generationCount)
	lastFrameDivisor = 0
}

var lastFrameDivisor int = 0

func (p *Population) Update(frameTime float32, agentUpdater func(*Agent)) {
	for _, agent := range p.agents {
		agentUpdater(agent)
	}

	if config.Config.Simulation.Evolver.Mode == "train" {
		p.currentGenerationLifetime += frameTime

		// speedCheckTime: Time after which we can check to make sure all agents are not stopped, in seconds.
		if p.currentGenerationLifetime > config.Config.Simulation.Evolver.MaxGenerationLifetime ||
			allAgentsDead(p.agents) ||
			(p.currentGenerationLifetime > config.Config.Simulation.Evolver.SpeedCheckTime && allAgentsStopped(p.agents)) {

			// Create a new generation by sorting, creating (in-place) new agents, and mutating them
			sort.Sort(sort.Reverse(p.agents))

			// Save the best agent
			fmt.Printf("High Score: %.2f\n", p.agents[0].GetFinalScore())
			p.agents[0].SaveNet()

			recombine(p.agents)
			mutate(p.agents)

			p.prepareNewGeneration()
		} else {
			if int(p.currentGenerationLifetime/5) != lastFrameDivisor {
				fmt.Printf("  %v seconds (%v agents left)\n", int(p.currentGenerationLifetime), agentAliveCount(p.agents))
				lastFrameDivisor = int(p.currentGenerationLifetime / 5)
			}
		}
	}
}

func (p *Population) Render(agentRenderer func(*Agent, float32)) {
	maxScore := float32(0)
	for _, agent := range p.agents {
		if agent.car.Score > maxScore {
			maxScore = agent.car.Score
		}
	}

	for _, agent := range p.agents {
		agentRenderer(agent, maxScore)
	}
}
