package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
)

func main() {
	app := fiber.New()

	app.Post("/login", func(c *fiber.Ctx) error {
		var body struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).SendString("invalid request")
		}

		// hash password input
		hash := sha1.Sum([]byte(body.Password))
		passInput := hex.EncodeToString(hash[:])

		// ambil password dari redis
		passSaved, err := rdb.HGet(ctx, "user:"+body.Username, "password").Result()
		if err != nil {
			return c.Status(401).SendString("user not found")
		}

		if passInput != passSaved {
			return c.Status(401).SendString("wrong password")
		}

		return c.SendString(fmt.Sprintf("welcome %s", body.Username))
	})

	app.Listen(":3000")
}