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

// AddMessage godoc
// @Summary Create a new message
// @Description Creates a new message by taking a content
// @Tags message
// @Accept  json
// @Produce  json
// @Param   input  body      domain.CreateMessageRequest  false  "Message Creation Data"
// @Success 201  {object} entity.Message
// @Failure 400  {object}  errorResponse
// @Failure 500 {object} errorResponse
// @Router /messages [post]
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

// DeleteMessageById godoc
// @Summary Delete a message
// @Description Delete a message by Id
// @Tags message
// @Accept  json
// @Produce  json
// @Param   id path string true "id"
// @Success 200 {string} string "Successfully deleted"
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /messages/{id} [delete]
func (h *Handler) deleteMessageById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, fmt.Errorf("invalid id param: %v", err))
		return
	}

	if err := h.service.Message.DeleteMessageById(c, id); err != nil {
		if errors.Is(err, domain.ErrMessageNotFound) {
			newErrorResponse(c, http.StatusText(http.StatusNotFound), http.StatusNotFound, fmt.Errorf("error deleting message: %v", err))
			return
		}
		newErrorResponse(c, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, fmt.Errorf("something was wrong: %s", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted"})
}

// GetMessageById godoc
// @Summary Get message by id
// @Description Getting message info by   id
// @Tags message
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} entity.Message
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /messages/{id} [get]
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

// GetMessagesList godoc
// @Summary Get messages list by filter
// @Description Getting messages info by filter
// @Tags message
// @Accept json
// @Produce json
// @Param content query string false "Message`s` content"
// @Param status query bool false "Message`s` status"
// @Param limit query int false "Message`s limit"
// @Param offset query int false "Message`s offset"
// @Success 200 {object} []entity.Message
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /messages [get]
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
