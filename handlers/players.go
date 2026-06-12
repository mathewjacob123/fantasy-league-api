package handlers

import (
	"fantasy-league-api/models"
	"fantasy-league-api/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PlayerHandler struct {
	playerService *services.PlayerService
}

func NewPlayerHandler(playerService *services.PlayerService) *PlayerHandler {
	return &PlayerHandler{playerService: playerService}
}

func (h *PlayerHandler) CreatePlayer(c *gin.Context) {
	var req models.CreatePlayerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player, err := h.playerService.CreatePlayer(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, player)
}

func (h *PlayerHandler) GetPlayers(c *gin.Context) {
	page, _     := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	sport       := c.DefaultQuery("sport", "")
	position    := c.DefaultQuery("position", "")

	players, err := h.playerService.GetPlayers(page, pageSize, sport, position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      players,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *PlayerHandler) GetPlayerByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player id"})
		return
	}

	player, err := h.playerService.GetPlayerByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, player)
}

func (h *PlayerHandler) AddPlayerToTeam(c *gin.Context) {
	teamID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team id"})
		return
	}

	var req models.AddPlayerToTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := int(c.GetFloat64("user_id"))

	err = h.playerService.AddPlayerToTeam(teamID, userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "player added to team successfully"})
}

func (h *PlayerHandler) RemovePlayerFromTeam(c *gin.Context) {
	teamID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team id"})
		return
	}

	playerID, err := strconv.Atoi(c.Param("playerId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player id"})
		return
	}

	userID := int(c.GetFloat64("user_id"))

	err = h.playerService.RemovePlayerFromTeam(teamID, playerID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "player removed from team successfully"})
}

func (h *PlayerHandler) GetTeamPlayers(c *gin.Context) {
	teamID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team id"})
		return
	}

	players, err := h.playerService.GetTeamPlayers(teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": players})
}