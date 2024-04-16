package main

type IssueResponse struct {
	Issue Issue `json:"issue"`
}

type IssuesResponse struct {
	Issues []Issue `json:"issues"`
}

type IssueProperty struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Status struct {
	IssueProperty
	IsClosed bool `json:"is_closed"`
}

type Parent struct {
	ID int `json:"id"`
}

type Issue struct {
	ID                  int           `json:"id"`
	Project             IssueProperty `json:"project"`
	Tracker             IssueProperty `json:"tracker"`
	Status              Status        `json:"status"`
	Priority            IssueProperty `json:"priority"`
	Author              IssueProperty `json:"author"`
	AssignedTo          IssueProperty `json:"assigned_to"`
	Parent              *Parent       `json:"parent,omitempty"` // Optional field
	Subject             string        `json:"subject"`
	Description         string        `json:"description"`
	StartDate           string        `json:"start_date"`
	DueDate             *string       `json:"due_date,omitempty"` // Optional field
	DoneRatio           int           `json:"done_ratio"`
	IsPrivate           bool          `json:"is_private"`
	EstimatedHours      *float64      `json:"estimated_hours,omitempty"`       // Optional field
	TotalEstimatedHours *float64      `json:"total_estimated_hours,omitempty"` // Optional field
	SpentHours          float64       `json:"spent_hours"`
	TotalSpentHours     float64       `json:"total_spent_hours"`
	CreatedOn           string        `json:"created_on"`
	UpdatedOn           string        `json:"updated_on"`
	ClosedOn            *string       `json:"closed_on,omitempty"` // Optional field
}
