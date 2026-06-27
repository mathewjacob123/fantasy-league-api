package services

import (
	"fantasy-league-api/models"
	"testing"
)

func TestCalculatePointsFootball(t *testing.T) {
	tests := []struct {
		name     string
		stats    models.PlayerStats
		expected int
	}{
		{
			name: "goal scorer with 90 mins",
			stats: models.PlayerStats{
				Goals:         2,
				Assists:       1,
				MinutesPlayed: 90,
				YellowCards:   0,
				RedCards:      0,
			},
			expected: 17, // 2*6 + 1*3 + 2
		},
		{
			name: "yellow card deduction",
			stats: models.PlayerStats{
				Goals:         1,
				Assists:       0,
				MinutesPlayed: 90,
				YellowCards:   1,
				RedCards:      0,
			},
			expected: 7, // 1*6 + 2 - 1
		},
		{
			name: "red card deduction",
			stats: models.PlayerStats{
				Goals:         0,
				Assists:       0,
				MinutesPlayed: 60,
				YellowCards:   0,
				RedCards:      1,
			},
			expected: -3, // 0 + 0 - 3
		},
		{
			name: "less than 90 mins no bonus",
			stats: models.PlayerStats{
				Goals:         1,
				Assists:       0,
				MinutesPlayed: 89,
				YellowCards:   0,
				RedCards:      0,
			},
			expected: 6, // 1*6, no 90 min bonus
		},
		{
			name: "exactly 90 mins gets bonus",
			stats: models.PlayerStats{
				Goals:         0,
				Assists:       0,
				MinutesPlayed: 90,
				YellowCards:   0,
				RedCards:      0,
			},
			expected: 2, // just the 90 min bonus
		},
		{
			name: "no contribution zero points",
			stats: models.PlayerStats{
				Goals:         0,
				Assists:       0,
				MinutesPlayed: 45,
				YellowCards:   0,
				RedCards:      0,
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculatePoints("football", tt.stats)
			if result != tt.expected {
				t.Errorf("expected %d points got %d", tt.expected, result)
			}
		})
	}
}

func TestCalculatePointsCricket(t *testing.T) {
	tests := []struct {
		name     string
		stats    models.PlayerStats
		expected int
	}{
		{
			name: "batsman with 50 runs",
			stats: models.PlayerStats{
				Runs:    50,
				Wickets: 0,
				Catches: 0,
			},
			expected: 5, // 50/10 * 1
		},
		{
			name: "bowler with 3 wickets",
			stats: models.PlayerStats{
				Runs:    0,
				Wickets: 3,
				Catches: 0,
			},
			expected: 75, // 3*25
		},
		{
			name: "all rounder",
			stats: models.PlayerStats{
				Runs:    30,
				Wickets: 2,
				Catches: 1,
			},
			expected: 61, // 3 + 50 + 8
		},
		{
			name: "catch taken",
			stats: models.PlayerStats{
				Runs:    0,
				Wickets: 0,
				Catches: 2,
			},
			expected: 16, // 2*8
		},
		{
			name: "no contribution",
			stats: models.PlayerStats{
				Runs:    0,
				Wickets: 0,
				Catches: 0,
			},
			expected: 0,
		},
		{
			name: "runs below 10 give zero points",
			stats: models.PlayerStats{
				Runs:    9,
				Wickets: 0,
				Catches: 0,
			},
			expected: 0, // 9/10 = 0 in integer division
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculatePoints("cricket", tt.stats)
			if result != tt.expected {
				t.Errorf("expected %d points got %d", tt.expected, result)
			}
		})
	}
}

func TestCalculatePointsUnknownSport(t *testing.T) {
	stats := models.PlayerStats{
		Goals: 5,
		Runs:  100,
	}

	result := calculatePoints("basketball", stats)

	if result != 0 {
		t.Errorf("expected 0 points for unknown sport, got %d", result)
	}
}