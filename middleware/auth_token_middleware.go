package middleware

import (
	"crypto/sha256"
	"encoding/json"
	"general/fiber-swagger/configs"
	"general/fiber-swagger/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"go.mongodb.org/mongo-driver/bson"

	fiberlog "github.com/gofiber/fiber/v2/log"
)

func AuthTokenMiddleware() func(*fiber.Ctx) error {
	// Create config for Bearer authentication middleware.
	config := keyauth.Config{
		AuthScheme: "Bearer",
		Validator: func(c *fiber.Ctx, key string) (bool, error) {
			user := ""
			// make hash from key
			h := sha256.New()
			h.Write([]byte(key))
			access_token_hash := h.Sum(nil)
			// get hash key from db
			tokensCollection := configs.GetCollection(configs.DB, "tokens")
			cached_token := tokensCollection.FindOne(c.Context(), bson.M{"access_token_hash": access_token_hash})
			token_record := new(models.AccessTokenRecord)
			if cached_token != nil {
				cached_token.Decode(&token_record)
				if token_record.User != "" {
					fiberlog.Debug(token_record.User)
					user = token_record.User
					c.Locals("user", user)
					return true, nil
				}
			}

			// if token is not in db, get user id from github and save to db
			request := fiber.Get("https://api.github.com/user")
			request.Debug()
			request.Set("Accept", "application/json").Set("Content-Type", "application/json").Set("Authorization", "Bearer "+key)
			statusCode, body, errs := request.Bytes()
			// Check for errors
			if errs != nil {
				fiberlog.Error(errs[0].Error())
				return false, errs[0]
			}
			// get user from body unmarshal
			if statusCode == fiber.StatusOK {
				json_resp := make(map[string]interface{})
				json.Unmarshal(body, &json_resp)
				user = json_resp["login"].(string)
				token_record.User = user
				token_record.AcessTokenHash = access_token_hash
				token_record.Created_time = time.Now()
				fiberlog.Debug(user)
				if user != "" {
					_, err := tokensCollection.InsertOne(c.Context(), token_record)
					c.Locals("user", user)
					// Check for errors
					if err != nil {
						fiberlog.Error(err)
						return false, err
					}
					return true, nil
				}
			}

			return false, keyauth.ErrMissingOrMalformedAPIKey
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Return status 403 Forbidden.
			return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{
				Error:   fiber.ErrForbidden.Message,
				Details: "Sorry you don't have permission to access this resource.",
			})
		},
	}

	return keyauth.New(config)
}
