package handler

import (
	"log"
	"net/http"
	"strconv"

	avitotask "github.com/ant0nix/avitoTask"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateUser(c *gin.Context) {
	var input avitotask.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	
	err := h.services.Start.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "User created successfully",
	})

}

func (h *Handler) ChangeBalance(c *gin.Context) {
	var input avitotask.Balance
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	req, err := h.services.InternalServices.ChangeBalance(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": req,
	})

}

func (h *Handler) ShowBalance(c *gin.Context) {
	var input avitotask.Balance
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	req, err := h.services.InternalServices.ShowBalance(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": req,
	})
}

func (h *Handler) P2p(c *gin.Context) {
	var input avitotask.P2p
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	req, err := h.services.InternalServices.P2p(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": req,
	})
}

func (h *Handler) CreateServices(c *gin.Context) {
	var input avitotask.Service
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err := h.services.Start.CreateServices(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Service created successfully",
	})

}

func (h *Handler) MakeOrder(c *gin.Context) {
	var input avitotask.Order
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id := c.Param("id")
	tmp, err := strconv.Atoi(id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	input.SId = tmp
	input.Amount, err = h.services.Service.GetServicesPrice(input.SId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println(input)
	req, err := h.services.Service.MakeOrder(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": req,
	})

}

func (h *Handler) ListServices(c *gin.Context) {
	req, err := h.services.Start.ShowServices()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": req,
	})

}

//TODO: привязка по id юзера
func (h *Handler) DoOrder(c *gin.Context) {
	id := c.Param("id")
	tmp, err := strconv.Atoi(id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	req, err := h.services.Service.DoOrder(tmp)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": req,
	})

}

/*
- отправка офера в фирму
*/
