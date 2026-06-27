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

// CreatePlayer godoc
// @Summary      Create a player (admin)
// @Description  Create a new real-world player — admin only
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.CreatePlayerRequest true "Create player request"
// @Success      201 {object} models.Player
// @Failure      400 {object} map[string]string
// @Router       /api/admin/players [post]
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

// GetPlayers godoc
// @Summary      List players
// @Description  Get all players with optional filters
// @Tags         players
// @Produce      json
// @Security     BearerAuth
// @Param        sport    query string false "Filter by sport (football/cricket)"
// @Param        position query string false "Filter by position"
// @Param        page     query int    false "Page number" default(1)
// @Param        page_size query int   false "Page size"  default(20)
// @Success      200 {object} map[string]interface{}
// @Router       /api/players [get]
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

// GetPlayerByID godoc
// @Summary      Get a player
// @Description  Get player details by ID
// @Tags         players
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Player ID"
// @Success      200 {object} models.Player
// @Failure      404 {object} map[string]string
// @Router       /api/players/{id} [get]
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

// AddPlayerToTeam godoc
// @Summary      Add player to team
// @Description  Add a real-world player to your fantasy team
// @Tags         players
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path int                          true "Team ID"
// @Param        request body models.AddPlayerToTeamRequest true "Add player request"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /api/teams/{id}/players [post]
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

// RemovePlayerFromTeam godoc
// @Summary      Remove player from team
// @Description  Remove a player from your fantasy team
// @Tags         players
// @Produce      json
// @Security     BearerAuth
// @Param        id       path int true "Team ID"
// @Param        playerId path int true "Player ID"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /api/teams/{id}/players/{playerId} [delete]
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
// GetTeamPlayers godoc
// @Summary      Get team roster
// @Description  Get all players in a fantasy team
// @Tags         players
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Team ID"
// @Success      200 {object} map[string]interface{}
// @Router       /api/teams/{id}/players [get]
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