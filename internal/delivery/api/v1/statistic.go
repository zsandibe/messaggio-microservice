package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zsandibe/messaggio-microservice/internal/domain"
)

// GetStatById godoc
// @Summary Get stat by id
// @Description Getting stat info by   id
// @Tags stats
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} entity.Stats
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /stats/{id} [get]
func (h *Handler) getStatById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, fmt.Errorf("invalid id param: %v", err))
		return
	}

	stat, err := h.service.Statistic.GetStatById(c, id)
	if err != nil {
		if errors.Is(err, domain.ErrStatisticNotFound) {
			newErrorResponse(c, http.StatusText(http.StatusNotFound), http.StatusNotFound, err)
			return
		}
		newErrorResponse(c, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, fmt.Errorf("something was wrong: %s", err))
		return
	}

	c.JSON(http.StatusOK, stat)
}

// GetStatsList godoc
// @Summary Get stats list by filter
// @Description Getting stats info by filter
// @Tags stats
// @Produce json
// @Param  id path string true "id"
// @Success 200 {object} []entity.Message
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /stats [get]
func (h *Handler) getStatsList(c *gin.Context) {
	stats, err := h.service.Statistic.GetStatsList(c)
	if err != nil {
		newErrorResponse(c, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, fmt.Errorf("something was wrong: %s", err))
		return
	}
	c.JSON(http.StatusOK, stats)
}
