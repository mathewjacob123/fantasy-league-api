package handlers

import (
	"fantasy-league-api/models"
	"fantasy-league-api/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LeagueHandler struct {
	leagueService *services.LeagueService
}

func NewLeagueHandler(leagueService *services.LeagueService) *LeagueHandler {
	return &LeagueHandler{leagueService: leagueService}
}

// CreateLeague godoc
// @Summary      Create a league
// @Description  Create a new fantasy league
// @Tags         leagues
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.CreateLeagueRequest true "Create league request"
// @Success      201 {object} models.League
// @Failure      400 {object} map[string]string
// @Router       /api/leagues [post]
func (h *LeagueHandler) CreateLeague(c *gin.Context) {
	var req models.CreateLeagueRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ownerID := int(c.GetFloat64("user_id"))

	league, err := h.leagueService.CreateLeague(req, ownerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, league)
}

// GetLeagues godoc
// @Summary      List leagues
// @Description  Get all leagues with pagination
// @Tags         leagues
// @Produce      json
// @Security     BearerAuth
// @Param        page      query int false "Page number" default(1)
// @Param        page_size query int false "Page size"  default(10)
// @Success      200 {object} map[string]interface{}
// @Router       /api/leagues [get]
func (h *LeagueHandler) GetLeagues(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	userID := int(c.GetFloat64("user_id"))

	leagues, err := h.leagueService.GetLeagues(page, pageSize, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      leagues,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetLeagueByID godoc
// @Summary      Get a league
// @Description  Get league details by ID
// @Tags         leagues
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "League ID"
// @Success      200 {object} models.LeagueResponse
// @Failure      404 {object} map[string]string
// @Router       /api/leagues/{id} [get]
func (h *LeagueHandler) GetLeagueByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid league id"})
		return
	}

	userID := int(c.GetFloat64("user_id"))

	league, err := h.leagueService.GetLeagueByID(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, league)
}

// JoinLeague godoc
// @Summary      Join a league
// @Description  Join an existing league
// @Tags         leagues
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "League ID"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /api/leagues/{id}/join [post]
func (h *LeagueHandler) JoinLeague(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid league id"})
		return
	}

	userID := int(c.GetFloat64("user_id"))

err = h.leagueService.JoinLeague(id, userID)
if err != nil {
    // check what kind of error it is
    if err.Error() == "league is no longer accepting members" || 
       err.Error() == "already a member of this league" ||
       err.Error() == "league is full" {
        c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
    return
}
}
func (h *LeagueHandler) GetLeaderboard(c *gin.Context) {
	leagueID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid league id"})
		return
	}

	teams, err := h.leagueService.GetLeaderboard(leagueID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": teams})
}