package apiserver

import (
	"net/http"
	"strconv"

	"github.com/alex11prog/ups-imitator/internal/app/model"
	"github.com/gin-gonic/gin"
)

type mode struct {
	Mode bool `json:"mode" example:"false" binding:"required"`
}

//	@Summary	method returns imitator mode
//	@Tags		Imitator
//	@Produce	json
//	@Success	200	{object}	mode
//	@Router		/imitator/mode [get]
func (s *server) handlerGetMode(c *gin.Context) {
	c.JSON(http.StatusOK, mode{s.imitator.GetMode()})
}

//	@Summary	method updates imitator mode
//	@Tags		Imitator
//	@Accept		json
//	@Param		input	body	mode	true	"mode"
//	@Produce	json
//	@Success	200	{object}	statusBody
//	@Router		/imitator/mode [put]
func (s *server) handlerUpdateMode(c *gin.Context) {
	var input mode
	if err := c.BindJSON(&input); err != nil {
		s.errorResponse(c, http.StatusBadRequest, err)
		return
	}
	s.imitator.SetMode(input.Mode)
	c.JSON(http.StatusOK, statusBody{"OK"})
}

//	@Summary	method returns all ups
//	@Tags		Imitator
//	@Produce	json
//	@Success	200	{array}	model.UpsParams
//	@Router		/imitator/ups [get]
func (s *server) handlerGetAllUps(c *gin.Context) {
	c.JSON(http.StatusOK, s.imitator.GetUpsParams())
}

//	@Summary	method updates ups params
//	@Tags		Imitator
//	@Accept		json
//	@Param		input	body	model.UpsParamsUpdateForm	true	"params"
//	@Param		ups_id	path	int							true	"UPS id"
//	@Produce	json
//	@Success	200	{object}	statusBody
//	@Failure	400	{object}	errorResponse	"invalid payload"
//	@Failure	422	{object}	errorResponse
//	@Router		/imitator/ups/{ups_id} [patch]
func (s *server) handlerUpdateUps(c *gin.Context) {
	ups_id, err := strconv.Atoi(c.Param("ups_id"))
	if err != nil {
		s.errorResponse(c, http.StatusBadRequest, err)
		return
	}
	var input model.UpsParamsUpdateForm
	if err := c.BindJSON(&input); err != nil {
		s.errorResponse(c, http.StatusBadRequest, err)
		return
	}
	if err := s.imitator.UpdateUps(ups_id, &input); err != nil {
		s.errorResponse(c, http.StatusUnprocessableEntity, err)
		return
	}
	c.JSON(http.StatusOK, statusBody{"OK"})
}

//	@Summary	method updates ups battery params
//	@Tags		Imitator
//	@Accept		json
//	@Param		input	body	model.BatteryParamsUpdateForm	true	"params"
//	@Produce	json
//	@Param		ups_id	path		int	true	"UPS id"
//	@Param		bat_id	path		int	true	"Battery id"
//	@Success	200		{object}	statusBody
//	@Failure	400		{object}	errorResponse	"invalid payload"
//	@Failure	422		{object}	errorResponse
//	@Router		/imitator/ups/{ups_id}/battery/{bat_id} [patch]
func (s *server) handlerUpdateUpsBattery(c *gin.Context) {
	ups_id, err := strconv.Atoi(c.Param("ups_id"))
	if err != nil {
		s.errorResponse(c, http.StatusBadRequest, err)
		return
	}
	bat_id, err := strconv.Atoi(c.Param("bat_id"))
	if err != nil {
		s.errorResponse(c, http.StatusBadRequest, err)
		return
	}
	var input model.BatteryParamsUpdateForm
	if err := c.BindJSON(&input); err != nil {
		s.errorResponse(c, http.StatusBadRequest, err)
		return
	}
	if err := s.imitator.UpdateUpsBattery(ups_id, bat_id, &input); err != nil {
		s.errorResponse(c, http.StatusUnprocessableEntity, err)
		return
	}
	c.JSON(http.StatusOK, statusBody{"OK"})
}
