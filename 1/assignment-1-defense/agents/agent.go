package Agents

import (
	"fmt"
)

type Agent interface {
	GetID() string
	GetName() string
}

type HumanAgent struct {
	ID   string
	Name string
}

type BotAgent struct {
	ID      string
	Name    string
	Version string
}

func FormatAgent(agent Agent) string {
	return fmt.Sprintf("%s | %s", agent.GetID(), agent.GetName())
}

func (a HumanAgent) GetID() string {
	return a.ID
}

func (a BotAgent) GetID() string {
	return a.ID
}

func (a HumanAgent) GetName() string {
	return a.Name
}

func (a BotAgent) GetName() string {
	return a.Name
}
