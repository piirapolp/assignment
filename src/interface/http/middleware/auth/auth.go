package auth

import (
	"assignment/entity"
	"assignment/global"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strings"
	"time"
)

func DBTokenAuth(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.Set("WWW-Authenticate", `Bearer realm="api", error="invalid_request"`)
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		tokenStr := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))

		var token entity.Tokens
		now := time.Now()

		err := db.WithContext(c.Context()).
			Where("session_id = ?", tokenStr).
			Where("expired_at > ?", now).
			Take(&token).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.Set("WWW-Authenticate", `Bearer realm="api", error="invalid_token"`)
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "invalid or expired token",
				})
			}
			// DB error
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		_ = db.Model(&entity.Tokens{}).
			Where("session_id = ?", token.SessionId).
			Update("expired_at", now.Add(720*time.Minute)).Error

		// Put user info into request context
		c.Locals(global.KEY_USER_ID, token.UserId)
		c.Locals("session_id", token.SessionId)
		return c.Next()
	}
}
