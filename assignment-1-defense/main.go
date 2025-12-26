package main

import (
	"fmt"
	"strings"

	Agents "github.com/SultanYakupov/CheckingTask2/agents"
	Tickets "github.com/SultanYakupov/CheckingTask2/tickets"
)

func main() {
	agents := make(map[string]Agents.Agent)
	store := Tickets.NewTicketStore()

cli:
	for {
		fmt.Println(`
Ticket store menu
[1] Create Ticket
[2] Add Agent (Human or Bot)
[3] Assign Ticket to Agent
[4] Resolve Ticket
[5] List All Tickets
[6] List All OPEN Tickets
[7] List All DONE Tickets
[8] List Unassigned Tickets
[9] Exit`)

		fmt.Println("\nAvailable agents:")
		for _, agent := range agents {
			fmt.Println(Agents.FormatAgent(agent))
		}
		fmt.Println()

		var choice int

		fmt.Print(">>> ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			var ticketID string

			fmt.Print("Ticket ID: ")
			fmt.Scanln(&ticketID)

			var ticketTitle string

			fmt.Print("Ticket Title: ")
			fmt.Scanln(&ticketTitle)

			var ticketDescription string

			fmt.Print("Ticket Description: ")
			fmt.Scanln(&ticketDescription)

			var ticketPriority int

			fmt.Print("Ticket Priority (1/2/3): ")
			fmt.Scanln(&ticketPriority)

			_, err := store.Create(ticketID, ticketTitle, ticketDescription, ticketPriority)

			if err != nil {
				fmt.Println(err)
			}

		case 2:
			var agentID string

			fmt.Print("Agent ID: ")
			fmt.Scanln(&agentID)

			agent, ok := agents[agentID]

			if ok {
				fmt.Println("Agent already exists")
				break
			}

			var agentName string

			fmt.Print("Agent Name: ")
			fmt.Scanln(&agentName)

			var agentBot string

			fmt.Print("Agent Bot (y/N): ")
			fmt.Scanln(&agentBot)

			agentBot = strings.ToLower(agentBot)

			if agentBot == "y" || agentBot == "yes" {
				var agentVersion string

				fmt.Print("Agent Version: ")
				fmt.Scanln(&agentVersion)

				agent = Agents.BotAgent{ID: agentID, Name: agentName, Version: agentVersion}
			} else {
				agent = Agents.HumanAgent{ID: agentID, Name: agentName}
			}

			agents[agentID] = agent

		case 3:
			var ticketID string

			fmt.Print("Ticket ID: ")
			fmt.Scanln(&ticketID)

			var agentID string

			fmt.Print("Agent ID: ")
			fmt.Scanln(&agentID)

			err := store.Assign(ticketID, agentID)

			if err != nil {
				fmt.Println(err)
			}

		case 4:
			var ticketID string

			fmt.Print("Ticket ID: ")
			fmt.Scanln(&ticketID)

			err := store.Resolve(ticketID)

			if err != nil {
				fmt.Println(err)
			}

		case 5:
			for _, ticket := range store.ListAll() {
				assignee := ticket.AssigneeID

				if len(assignee) <= 0 {
					assignee = "(no agent)"
				}

				fmt.Printf("[%s -> %s | %d | %s] %s. %s\n", ticket.ID, assignee, ticket.Priority, ticket.Status, ticket.Title, ticket.Description)
			}

		case 6:
			for _, ticket := range store.ListByStatus("OPEN") {
				assignee := ticket.AssigneeID

				if len(assignee) <= 0 {
					assignee = "(no agent)"
				}

				fmt.Printf("[%s -> %s | %d] %s. %s\n", ticket.ID, assignee, ticket.Priority, ticket.Title, ticket.Description)
			}

		case 7:
			for _, ticket := range store.ListByStatus("DONE") {
				assignee := ticket.AssigneeID

				if len(assignee) <= 0 {
					assignee = "(no agent)"
				}

				fmt.Printf("[%s -> %s | %d] %s. %s\n", ticket.ID, assignee, ticket.Priority, ticket.Title, ticket.Description)
			}

		case 8:
			for _, ticket := range store.ListUnassigned() {
				fmt.Printf("[%s | %d | %s] %s. %s\n", ticket.ID, ticket.Priority, ticket.Status, ticket.Title, ticket.Description)
			}

		case 9:
			break cli
		}
	}
}
