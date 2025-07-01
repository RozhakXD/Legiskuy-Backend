package election

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// @Summary Cast a vote
// @Description Cast a vote for a candidate in the election
// @Tags election
// @Accept json
// @Produce json
// @Param vote body CastVoteInput true "Vote Data"
// @Success 200 {object} map[string]string "Vote cast successfully"
// @Failure 400 {object} map[string]string "Bad request - cannot parse JSON or missing required fields"
// @Failure 403 {object} map[string]string "Forbidden - election is not currently active"
// @Failure 404 {object} map[string]string "Not found - voter or candidate not found"
// @Failure 409 {object} map[string]string "Conflict - voter has already voted"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /vote [post]
func (h *Handler) CastVote(c *fiber.Ctx) error {
	input := new(CastVoteInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	err := h.service.CastVote(input)
	if err != nil {
		switch err.Error() {
		case "voter not found", "candidate not found":
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		case "voter has already voted":
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		case "voter_id and candidate_id are required":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		case "election is not currently active":
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": err.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to cast vote",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Vote cast successfully",
	})
}

// @Summary Set election time
// @Description Set the start and end time for the election
// @Tags election
// @Accept json
// @Produce json
// @Param time body SetTimeInput true "Election Time Data"
// @Success 200 {object} map[string]string "Election time set successfully"
// @Failure 400 {object} map[string]string "Bad request - cannot parse JSON or invalid time format"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /election/time [post]
func (h *Handler) SetElectionTime(c *fiber.Ctx) error {
	input := new(SetTimeInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	err := h.service.SetElectionTime(input)
	if err != nil {
		if err.Error() == "invalid time format, use RFC3339 format (e.g., 2025-06-13T00:00:00Z)" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to set election time",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Election time set successfully",
	})
}

// @Summary Set election threshold
// @Description Set the minimum vote threshold for candidates to be qualified
// @Tags election
// @Accept json
// @Produce json
// @Param threshold body SetThresholdInput true "Threshold Data"
// @Success 200 {object} map[string]string "Threshold set successfully"
// @Failure 400 {object} map[string]string "Bad request - cannot parse JSON or invalid threshold value"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /election/threshold [post]
func (h *Handler) SetThreshold(c *fiber.Ctx) error {
	input := new(SetThresholdInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	err := h.service.SetThreshold(input)
	if err != nil {
		if err.Error() == "threshold is required" || err.Error() == "threshold must be a non-negative number" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to set threshold",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Threshold set successfully",
	})
}

// @Summary Get election results
// @Description Get election results with optional filtering for qualified candidates only
// @Tags election
// @Accept json
// @Produce json
// @Param qualified query bool false "Filter only qualified candidates (default: false)"
// @Success 200 {object} map[string]interface{} "Election results"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /election/results [get]
func (h *Handler) GetResults(c *fiber.Ctx) error {
	qualifiedStr := c.Query("qualified", "false")
	qualifiedOnly, _ := strconv.ParseBool(qualifiedStr)

	results, err := h.service.GetResults(qualifiedOnly)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get results",
		})
	}
	return c.JSON(results)
}
