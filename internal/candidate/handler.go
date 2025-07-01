package candidate

import (
	"database/sql"
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

// @Summary Create a new candidate
// @Description Create a new candidate with the provided name and party
// @Tags candidate
// @Accept json
// @Produce json
// @Param candidate body CreateCandidateInput true "Candidate Data"
// @Success 201 {object} map[string]interface{} "Candidate created successfully"
// @Failure 400 {object} map[string]string "Bad request - cannot parse JSON or missing required fields"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /candidates [post]
func (h *Handler) CreateCandidate(c *fiber.Ctx) error {
	input := new(CreateCandidateInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	candidate, err := h.service.CreateCandidate(input)
	if err != nil {
		if err.Error() == "name is required" || err.Error() == "party is required" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create candidate",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(candidate)
}

// @Summary Get all candidates
// @Description Get all candidates with optional filtering and sorting
// @Tags candidate
// @Accept json
// @Produce json
// @Param name query string false "Filter candidates by name"
// @Param party query string false "Filter candidates by party"
// @Param sort_by query string false "Sort by field (name, party, vote_count)"
// @Param order query string false "Sort order (asc, desc)"
// @Success 200 {array} map[string]interface{} "List of candidates"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /candidates [get]
func (h *Handler) GetAllCandidates(c *fiber.Ctx) error {
	name := c.Query("name")
	party := c.Query("party")

	sortBy := c.Query("sort_by")
	order := c.Query("order")

	candidates, err := h.service.GetAllCandidates(name, party, sortBy, order)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get candidates",
		})
	}
	return c.JSON(candidates)
}

// @Summary Get candidate by ID
// @Description Get a specific candidate by their ID
// @Tags candidate
// @Accept json
// @Produce json
// @Param id path int true "Candidate ID"
// @Success 200 {object} map[string]interface{} "Candidate details"
// @Failure 400 {object} map[string]string "Bad request - invalid candidate ID"
// @Failure 404 {object} map[string]string "Not found - candidate not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /candidates/{id} [get]
func (h *Handler) GetCandidateByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid candidate ID",
		})
	}
	candidate, err := h.service.GetCandidateByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get candidate",
		})
	}
	if candidate == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Candidate not found",
		})
	}

	return c.JSON(candidate)
}

// @Summary Update candidate
// @Description Update an existing candidate's information
// @Tags candidate
// @Accept json
// @Produce json
// @Param id path int true "Candidate ID"
// @Param candidate body UpdateCandidateInput true "Updated candidate data"
// @Success 200 {object} map[string]interface{} "Updated candidate details"
// @Failure 400 {object} map[string]string "Bad request - invalid candidate ID or cannot parse JSON"
// @Failure 404 {object} map[string]string "Not found - candidate not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /candidates/{id} [put]
func (h *Handler) UpdateCandidate(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid candidate ID",
		})
	}

	input := new(UpdateCandidateInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	candidate, err := h.service.UpdateCandidate(id, input)
	if err != nil {
		if err.Error() == "name is required" || err.Error() == "party is required" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Candidate not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update candidate",
		})
	}

	if candidate == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Candidate not found",
		})
	}

	return c.JSON(candidate)
}

// @Summary Delete candidate
// @Description Delete a candidate by their ID
// @Tags candidate
// @Accept json
// @Produce json
// @Param id path int true "Candidate ID"
// @Success 200 {object} map[string]string "Candidate deleted successfully"
// @Failure 400 {object} map[string]string "Bad request - invalid candidate ID"
// @Failure 404 {object} map[string]string "Not found - candidate not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /candidates/{id} [delete]
func (h *Handler) DeleteCandidate(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid candidate ID",
		})
	}

	err = h.service.DeleteCandidate(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Candidate not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete candidate",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Candidate deleted successfully",
	})
}
