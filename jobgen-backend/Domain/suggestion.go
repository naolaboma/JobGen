package domain

type Suggestion struct{
	ID        string    `json:"id" `
	Type 	  string 	`json:"type"`
	Content   string 	`json:"content"`
	Applied   string 	`json:"applied"`
}
