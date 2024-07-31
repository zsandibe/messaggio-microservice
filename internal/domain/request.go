package domain

type CreateMessageRequest struct {
	Content string `json:"content"`
}

type UpdateMessageRequest struct {
	Id int `json:"id"`
}

type MessagesListParams struct {
	Content string `form:"content"`
	Status  bool   `form:"status"`
	Limit   int    `form:"limit"`
	Offset  int    `form:"offset"`
}
