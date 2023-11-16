package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"general/fiber-swagger/configs"
	"general/fiber-swagger/models"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	fiberlog "github.com/gofiber/fiber/v2/log"
)

// @Summary		Github Callback
// @Description	Github Callback
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			code	query		string	true	"code"
// @Success		200		{object}	models.TodoId
// @Failure		400		{object}	models.ErrorResponse
// @Failure		401		{object}	models.ErrorResponse
// @Failure		403		{object}	models.ErrorResponse
// @Failure		500		{object}	models.ErrorResponse
// @Router			/v1/auth/github/callback [get]
func GithubCallback(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// get code from query params
	code := c.Query("code")
	// get state from query params
	// state := c.Query("state")
	fiberlog.Debug(code)
	// get access token from github
	request := fiber.
		Post("https://github.com/login/oauth/access_token").
		Debug().
		Set("Accept", "application/json").
		Set("Content-Type", "application/json").
		JSON(fiber.Map{"client_id": os.Getenv("GITHUB_OAUTH_CLIENT_ID"), "client_secret": os.Getenv("GITHUB_OAUTH_CLIENT_SECRET"), "code": code, "redirect_uri": os.Getenv("GITHUB_OAUTH_CALLBACK_URL")})
	stcode, body, errs := request.Bytes()
	// Check for errors
	if errs != nil {
		fiberlog.Error(errs[0].Error())
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   fiber.ErrInternalServerError.Message,
			Details: errs[0].Error(),
		})
	}
	// handle error from github
	json_resp := make(map[string]interface{})
	if stcode == fiber.StatusOK {
		json.Unmarshal(body, &json_resp)
		if json_resp["error"] != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error:   fiber.ErrBadRequest.Message,
				Details: json_resp["error_description"].(string),
			})
		}
	}

	// get access_token from body unmarshal
	auth_resp := new(models.AuthResponse)
	json.Unmarshal(body, &auth_resp)

	// get user id from github
	request = fiber.
		Get("https://api.github.com/user").
		Debug().
		Set("Accept", "application/json").
		Set("Content-Type", "application/json").
		Set("Authorization", "Bearer "+auth_resp.AccessToken)
	statusCode, body, errs := request.Bytes()
	// Check for errors
	if errs != nil {
		fiberlog.Error(errs[0].Error())
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   fiber.ErrInternalServerError.Message,
			Details: errs[0].Error(),
		})
	}
	// get user from body unmarshal
	user := ""
	if statusCode == fiber.StatusOK {
		json.Unmarshal(body, &json_resp)
		user = json_resp["login"].(string)
	}
	fiberlog.Debug(user)

	// put user and access_token_hash to database
	token_record := new(models.AccessTokenRecord)
	token_record.User = user
	h := sha256.New()
	h.Write([]byte(auth_resp.AccessToken))
	token_record.AcessTokenHash = h.Sum(nil)
	token_record.Created_time = time.Now()
	// Save the token to the mongo database
	tokensCollection := configs.GetCollection(configs.DB, "tokens")
	_, err := tokensCollection.InsertOne(ctx, token_record)
	// Check for errors
	if err != nil {
		fiberlog.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   fiber.ErrInternalServerError.Message,
			Details: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(auth_resp)
}
