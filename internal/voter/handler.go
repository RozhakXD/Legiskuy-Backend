package voter

import (
	"database/sql"
	"strconv"
	"strings"

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

// @Summary Create a new voter
// @Description Create a new voter with the provided name
// @Tags voter
// @Accept json
// @Produce json
// @Param voter body CreateVoterInput true "Voter Data"
// @Success 201 {object} map[string]interface{} "Voter created successfully"
// @Failure 400 {object} map[string]string "Bad request - cannot parse JSON or name is required"
// @Failure 409 {object} map[string]string "Conflict - voter name already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /voters [post]
func (h *Handler) CreateVoter(c *fiber.Ctx) error {
	input := new(CreateVoterInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	voter, err := h.service.CreateVoter(input)
	if err != nil {
		if err.Error() == "name is required" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Voter name already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create voter",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(voter)
}

// @Summary Get all voters
// @Description Get all voters with optional name filtering
// @Tags voter
// @Accept json
// @Produce json
// @Param name query string false "Filter voters by name"
// @Success 200 {array} map[string]interface{} "List of voters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /voters [get]
func (h *Handler) GetAllVoters(c *fiber.Ctx) error {
	name := c.Query("name")
	voters, err := h.service.GetAllVoters(name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get voters",
		})
	}
	return c.JSON(voters)
}

// @Summary Get voter by ID
// @Description Get a specific voter by their ID
// @Tags voter
// @Accept json
// @Produce json
// @Param id path int true "Voter ID"
// @Success 200 {object} map[string]interface{} "Voter details"
// @Failure 400 {object} map[string]string "Bad request - invalid voter ID"
// @Failure 404 {object} map[string]string "Not found - voter not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /voters/{id} [get]
func (h *Handler) GetVoterByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid voter ID",
		})
	}
	voter, err := h.service.GetVoterByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get voter",
		})
	}
	if voter == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Voter not found",
		})
	}
	return c.JSON(voter)
}

// @Summary Update voter
// @Description Update an existing voter's information
// @Tags voter
// @Accept json
// @Produce json
// @Param id path int true "Voter ID"
// @Param voter body UpdateVoterInput true "Updated voter data"
// @Success 200 {object} map[string]interface{} "Updated voter details"
// @Failure 400 {object} map[string]string "Bad request - invalid voter ID or cannot parse JSON"
// @Failure 404 {object} map[string]string "Not found - voter not found"
// @Failure 409 {object} map[string]string "Conflict - voter name already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /voters/{id} [put]
func (h *Handler) UpdateVoter(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid voter ID",
		})
	}
	input := new(UpdateVoterInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	voter, err := h.service.UpdateVoter(id, input)
	if err != nil {
		if err.Error() == "name is required" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Voter name already exists",
			})
		}
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Voter not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update voter",
		})
	}
	return c.JSON(voter)
}

// @Summary Delete voter
// @Description Delete a voter by their ID
// @Tags voter
// @Accept json
// @Produce json
// @Param id path int true "Voter ID"
// @Success 200 {object} map[string]string "Voter deleted successfully"
// @Failure 400 {object} map[string]string "Bad request - invalid voter ID"
// @Failure 404 {object} map[string]string "Not found - voter not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /voters/{id} [delete]
func (h *Handler) DeleteVoter(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid voter ID",
		})
	}

	err = h.service.DeleteVoter(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Voter not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete voter",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Voter deleted successfully",
	})
}
