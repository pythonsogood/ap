package Tickets

type Ticket struct {
	ID          string
	Title       string
	Description string
	Priority    int    // 1 low, 2 medium, 3 high
	AssigneeID  string // may be null
	Status      string // OPEN or DONE
}
