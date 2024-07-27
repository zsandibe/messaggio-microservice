package domain

import "fmt"

var (
	ErrMessageIncorrect  = fmt.Errorf("message content is incorrect")
	ErrMessageNotFound   = fmt.Errorf("message not found")
	ErrStatisticNotFound = fmt.Errorf("statistic not found")
	ErrCreatingMessage   = fmt.Errorf("error creating message")
)
