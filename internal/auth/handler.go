package auth

import (
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

// @Summary Register a new user
// @Description Register a new user account with username, password, name, and role
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterInput true "User Registration Data"
// @Success 201 {object} map[string]interface{} "User successfully registered"
// @Failure 400 {object} map[string]string "Bad request - validation errors or username already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /register [post]
func (h *Handler) Register(c *fiber.Ctx) error {
	input := new(RegisterInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	user, err := h.service.Register(input)
	if err != nil {
		if err.Error() == "username already exists" || err.Error() == "name is required" || err.Error() == "username is required" || err.Error() == "password is required" || err.Error() == "role is required" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to register user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// @Summary Login a user
// @Description Login a user and get a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginInput true "Login Credentials"
// @Success 200 {object} map[string]string "Login successful with JWT token"
// @Failure 400 {object} map[string]string "Bad request - cannot parse JSON"
// @Failure 401 {object} map[string]string "Unauthorized - invalid credentials"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /login [post]
func (h *Handler) Login(c *fiber.Ctx) error {
	input := new(LoginInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	token, err := h.service.Login(input)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "invalid credentials" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to login",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token":   token,
		"message": "Login successful",
	})
}
