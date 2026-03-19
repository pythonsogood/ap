package Tickets

import "errors"

type TicketStore struct {
	items map[string]Ticket
}

func (ts *TicketStore) Create(id string, title string, description string, priority int) (*Ticket, error) {
	if len(id) <= 0 {
		return nil, errors.New("Id must be non-empty")
	}

	if priority != 1 && priority != 2 && priority != 3 {
		return nil, errors.New("Priority must be 1..3")
	}

	ticket, ok := ts.items[id]

	if ok {
		return nil, errors.New("ID must be unique")
	}

	ticket = Ticket{ID: id, Title: title, Description: description, Priority: priority, Status: "OPEN"}

	return &ticket, nil
}

func (ts *TicketStore) Assign(ticketID string, assigneeID string) error {
	if len(assigneeID) <= 0 {
		return errors.New("assigneeID must be non-empty")
	}

	ticket, ok := ts.items[ticketID]

	if !ok {
		return errors.New("Ticket must exist")
	}

	if ticket.Status != "OPEN" {
		return errors.New("Status must be OPEN")
	}

	ticket.AssigneeID = assigneeID
	ts.items[ticketID] = ticket

	return nil
}

func (ts *TicketStore) Resolve(ticketID string) error {
	ticket, ok := ts.items[ticketID]

	if !ok {
		return errors.New("Ticket must exist")
	}

	ticket.Status = "DONE"
	ts.items[ticketID] = ticket

	return nil
}

func (ts *TicketStore) ListAll() []Ticket {
	var tickets []Ticket

	for _, ticket := range ts.items {
		tickets = append(tickets, ticket)
	}

	return tickets
}

func (ts *TicketStore) ListByStatus(status string) []Ticket {
	var tickets []Ticket

	for _, ticket := range ts.items {
		if ticket.Status == status {
			tickets = append(tickets, ticket)
		}
	}

	return tickets
}

func (ts *TicketStore) ListUnassigned() []Ticket {
	var tickets []Ticket

	for _, ticket := range ts.items {
		if len(ticket.AssigneeID) <= 0 {
			tickets = append(tickets, ticket)
		}
	}

	return tickets
}

func NewTicketStore() *TicketStore {
	return &TicketStore{items: make(map[string]Ticket)}
}
