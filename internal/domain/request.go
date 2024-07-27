package domain

type CreateMessageRequest struct {
	Content string `json:"content"`
}

type UpdateMessageRequest struct {
	Id int `json:"id"`
}

type MessagesListParams struct {
	Content string `json:"content"`
	Status  bool   `json:"status"`
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
}
