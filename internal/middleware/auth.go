package middleware

import (
	"school_attendance_backend/internal/utils"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func AuthRequired(jwtUtil *utils.JWTUtil) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Bearer token is required",
			})
		}

		claims, err := jwtUtil.GetClaims(tokenString)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Add user info to context
		ctx.Locals("userID", claims["user_id"])
		ctx.Locals("userRole", claims["role"])

		return ctx.Next()
	}
}

func AdminOnly(ctx fiber.Ctx) error {
	role := ctx.Locals("userRole").(string)
	if role != "admin" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}
	return ctx.Next()
}
