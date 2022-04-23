package api

import (
	"clips/internal/db"
	"clips/pkg/models"
	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
)

func GetMe(c *fiber.Ctx) error {

	includeClips := c.Query("includeClips","false") == "true"
	user := c.Locals("user").(models.User)

	if !includeClips {
		return c.JSON(&user)
	}

	err := db.Database.Model(&user).
		Where("user_id = ?",user.UserID).
		Relation("Clips").
		Select(&user)

	if err != nil {
		log.Println(err)

		return c.Status(500).JSON(fiber.Map{
			"message":"Unexpected Error",
		})
	}

	return c.JSON(&user)
}

func GetUser(c *fiber.Ctx) error {
	userID,_ := strconv.ParseInt(c.Params("user_id","0"),10,64)

	user := models.User{}

	err := db.Database.Model(&user).
		Where("user_id = ?",userID).
		Select()

	if err != nil {

		if err == pg.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{
				"error":"User not found",
			})
		}

		log.Println(err)

		return c.Status(500).JSON(fiber.Map{
			"error":"Unknown Error",
		})
	}

	return c.JSON(&user)
}