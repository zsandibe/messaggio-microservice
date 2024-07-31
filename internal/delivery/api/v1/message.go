package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zsandibe/messaggio-microservice/internal/domain"
	"github.com/zsandibe/messaggio-microservice/pkg"
)

func (h *Handler) addMessage(c *gin.Context) {
	var inp domain.CreateMessageRequest

	if err := c.BindJSON(&inp); err != nil {
		newErrorResponse(c, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, fmt.Errorf("invalid request body: %v", err))
		return
	}

	if !pkg.ValidateContent(inp.Content) {
		newErrorResponse(c, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, errors.New("not a valid content"))
		return
	}

	message, err := h.service.Message.CreateMessage(c, inp)
	if err != nil {
		newErrorResponse(c, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, fmt.Errorf("something was wrong: %v", err))
		return
	}

	bytesMessage, err := json.Marshal(message)
	if err != nil {
		newErrorResponse(c, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, fmt.Errorf("something was wrong: %v", err))
		return
	}

	bytesInp, err := json.Marshal(inp)
	if err != nil {
		newErrorResponse(c, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, fmt.Errorf("failed to marshal input: %v", err))
		return
	}

	id := strconv.Itoa(message.Id)

	err = h.service.PublishMessage(c, bytesMessage, bytesInp, id)
	if err != nil {
		newErrorResponse(c, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, fmt.Errorf("failed to publish message: %v", err))
		return
	}

	c.JSON(http.StatusCreated, message)
}

func (h *Handler) deleteMessageById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, fmt.Errorf("invalid id param: %v", err))
		return
	}

	if err := h.service.Message.DeleteMessageById(c, id); err != nil {
		if err == domain.ErrMessageNotFound {
			newErrorResponse(c, http.StatusText(http.StatusNotFound), http.StatusNotFound, fmt.Errorf("error deleting message: %v", err))
			return
		}
		newErrorResponse(c, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, fmt.Errorf("something was wrong: %s", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted"})
}

func (h *Handler) getMessageById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, fmt.Errorf("invalid id param: %v", err))
		return
	}

	message, err := h.service.Message.GetMessageById(c, id)
	if err != nil {
		if errors.Is(err, domain.ErrMessageNotFound) {
			newErrorResponse(c, http.StatusText(http.StatusNotFound), http.StatusNotFound, err)
			return
		}
		newErrorResponse(c, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, fmt.Errorf("something was wrong: %s", err))
		return
	}

	c.JSON(http.StatusOK, message)
}

func (h *Handler) getMessagesList(c *gin.Context) {
	var inp domain.MessagesListParams

	if err := c.ShouldBindQuery(&inp); err != nil {
		newErrorResponse(c, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, fmt.Errorf("invalid request body: %v", err))
		return
	}
	messages, err := h.service.GetMessagesList(c, inp)
	if err != nil {
		newErrorResponse(c, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, fmt.Errorf("something was wrong: %s", err))
		return
	}

	if len(messages) == 0 {
		newErrorResponse(c, http.StatusText(http.StatusNotFound), http.StatusNotFound, domain.ErrMessageNotFound)
		return
	}

	c.JSON(http.StatusOK, messages)
}
