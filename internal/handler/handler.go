package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"task_axxon/internal/service"
)

type Handler struct {
	srv service.Service
}

func NewHandler(srv service.Service) *Handler  {
	return &Handler{srv: srv}
}

func (h *Handler) RedirectApi(c echo.Context) error  {
	//struct to unmarshall request body
	var bodyMap map[string]interface{}

	//echo framework function to unmarshall request body
	if bindErr := c.Bind(&bodyMap); bindErr!=nil{
		log.Warnf("invalid struct: %v, in %v", bindErr.Error(), c.Request().RequestURI)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body struct")
	}

	//service method call
	res, err := h.srv.MakeRequest(c.Request().Context(), bodyMap)
	if err != nil{
		log.Errorf("failed response: %v", err)
		return err
	}

	return c.JSON(http.StatusOK, res)
}
