package handlers

import (
	"fantasy-league-api/models"
	"fantasy-league-api/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	teamService *services.TeamService
}

func NewTeamHandler(teamService *services.TeamService) *TeamHandler {
	return &TeamHandler{teamService: teamService}
}

// CreateTeam godoc
// @Summary      Create a team
// @Description  Create a fantasy team in a league
// @Tags         teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path int                    true "League ID"
// @Param        request body models.CreateTeamRequest true "Create team request"
// @Success      201 {object} models.Team
// @Failure      400 {object} map[string]string
// @Router       /api/leagues/{id}/teams [post]
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	leagueID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid league id"})
		return
	}

	var req models.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := int(c.GetFloat64("user_id"))

	team, err := h.teamService.CreateTeam(req, leagueID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, team)
}

// GetTeamsByLeague godoc
// @Summary      List teams in a league
// @Description  Get all teams in a specific league
// @Tags         teams
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "League ID"
// @Success      200 {object} map[string]interface{}
// @Router       /api/leagues/{id}/teams [get]
func (h *TeamHandler) GetTeamsByLeague(c *gin.Context) {
	leagueID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid league id"})
		return
	}

	teams, err := h.teamService.GetTeamsByLeague(leagueID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": teams})
}

// GetTeamByID godoc
// @Summary      Get a team
// @Description  Get team details by ID
// @Tags         teams
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Team ID"
// @Success      200 {object} models.Team
// @Failure      404 {object} map[string]string
// @Router       /api/teams/{id} [get]
func (h *TeamHandler) GetTeamByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team id"})
		return
	}

	userID := int(c.GetFloat64("user_id"))

	team, err := h.teamService.GetTeamByID(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, team)
}