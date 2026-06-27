package handlers

import (
	"fantasy-league-api/models"
	"fantasy-league-api/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MatchHandler struct {
	scoringService *services.ScoringService
}

func NewMatchHandler(scoringService *services.ScoringService) *MatchHandler {
	return &MatchHandler{scoringService: scoringService}
}

// CreateMatch godoc
// @Summary      Create a match (admin)
// @Description  Create a new real-world match — admin only
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.CreateMatchRequest true "Create match request"
// @Success      201 {object} models.Match
// @Failure      400 {object} map[string]string
// @Router       /api/admin/matches [post]
func (h *MatchHandler) CreateMatch(c *gin.Context) {
	var req models.CreateMatchRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	match, err := h.scoringService.CreateMatch(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, match)
}

// SubmitStats godoc
// @Summary      Submit match stats (admin)
// @Description  Submit player stats for a match and trigger scoring — admin only
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path int                      true "Match ID"
// @Param        request body models.SubmitStatsRequest true "Stats request"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /api/admin/matches/{id}/stats [post]
func (h *MatchHandler) SubmitStats(c *gin.Context) {
	matchID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
		return
	}

	var req models.SubmitStatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.scoringService.SubmitMatchStats(matchID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stats submitted and points calculated successfully"})
}

// GetMatch godoc
// @Summary      Get a match
// @Description  Get match details by ID
// @Tags         matches
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Match ID"
// @Success      200 {object} models.Match
// @Failure      404 {object} map[string]string
// @Router       /api/matches/{id} [get]
func (h *MatchHandler) GetMatch(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
		return
	}

	match, err := h.scoringService.GetMatchByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, match)
}

// GetMatchScores godoc
// @Summary      Get match scores
// @Description  Get fantasy points scored by each team in a match
// @Tags         matches
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Match ID"
// @Success      200 {object} map[string]interface{}
// @Router       /api/matches/{id}/scores [get]
func (h *MatchHandler) GetMatchScores(c *gin.Context) {
	matchID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
		return
	}

	scores, err := h.scoringService.GetMatchScores(matchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": scores})
}