package genetics

import (
	"fmt"
	"go-experiments/voxelli/config"
	"math"
	"math/rand"
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

	lastFrameDivisor int
	agentsRetrained  int
}

func NewPopulation(agentCount int, agentCreator func(int) *Agent) *Population {
	population := Population{
		generationCount:  0,
		agents:           make([]*Agent, agentCount),
		lastFrameDivisor: 0,
		agentsRetrained:  0}

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
	p.lastFrameDivisor = 0
}

func (p *Population) Update(frameTime float32, agentUpdater func(*Agent)) {
	for _, agent := range p.agents {
		agentUpdater(agent)
	}

	isReportInterval := func() bool {
		return int(p.currentGenerationLifetime/float32(config.Config.Simulation.Evolver.ReportInterval)) != p.lastFrameDivisor
	}

	resetReportInterval := func() {
		p.lastFrameDivisor = int(p.currentGenerationLifetime / float32(config.Config.Simulation.Evolver.ReportInterval))
	}

	saveBestAgent := func(agent *Agent) {
		// Save the best agent
		fmt.Printf("Best Agent High Score: %.2f\n", agent.GetFinalScore())
		agent.SaveNet()
	}

	if config.Config.Simulation.Evolver.Mode == "train" {
		p.currentGenerationLifetime += frameTime

		// speedCheckTime: Time after which we can check to make sure all agents are not stopped, in seconds.
		if p.currentGenerationLifetime > config.Config.Simulation.Evolver.MaxGenerationLifetime ||
			allAgentsDead(p.agents) ||
			(p.currentGenerationLifetime > config.Config.Simulation.Evolver.SpeedCheckTime && allAgentsStopped(p.agents)) {

			// Create a new generation by sorting, creating (in-place) new agents, and mutating them
			sort.Sort(sort.Reverse(p.agents))

			saveBestAgent(p.agents[0])

			recombine(p.agents)
			mutate(p.agents)

			p.prepareNewGeneration()
		} else {
			if isReportInterval() {
				fmt.Printf("  %v seconds (%v agents left)\n", int(p.currentGenerationLifetime), agentAliveCount(p.agents))
				resetReportInterval()
			}
		}
	} else if config.Config.Simulation.Evolver.Mode == "demoTrain" {
		p.currentGenerationLifetime += frameTime

		// Resort all agents to get the highest storing agents.
		sort.Sort(sort.Reverse(p.agents))

		for i, agent := range p.agents {

			// The agent has lived long enough it shouldn't be stopped, or it is the lowest scoring agent. Respawn it
			// -agent.wallHitTime
			if agent.GetLifetime() > config.Config.Simulation.Evolver.SpeedCheckTime &&
				(math.Abs(float64(agent.car.Velocity)) < float64(config.Config.Simulation.Agent.WiggleSpeed*100) ||
					i == len(p.agents)-1) {

				agent.Reset()

				// Never regenerate the best agent
				if i != 0 {
					// Always pick among the four best agents for regeneration
					agentToUse := rand.Uint32() % 4
					nextAgentToUse := agentToUse + 1 + rand.Uint32()%2

					agent.CrossBreed(p.agents[agentToUse], p.agents[nextAgentToUse], config.Config.Simulation.Evolver.CrossoverProbability)

					agent.net.ProbablyRandomize(
						config.Config.Simulation.Evolver.MutationProbability,
						config.Config.Simulation.Evolver.MutationAmount)
				} else {
					// The best agent died, so save it for future speed in regeneration
					saveBestAgent(agent)
				}

				p.agentsRetrained++
			}
		}

		if isReportInterval() {
			fmt.Printf("  %v seconds (%v agents retrained) High Score: (%v)\n",
				int(p.currentGenerationLifetime), p.agentsRetrained, p.agents[0].GetFinalScore())
			resetReportInterval()
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
