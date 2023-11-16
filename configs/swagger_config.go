package configs

import (
	swagger_docs_v1 "general/fiber-swagger/docs/v1"
	"os"

	"github.com/gofiber/swagger"
)

func init() {
	swagger_docs_v1.SwaggerInfo.Title = "TODO API"
	swagger_docs_v1.SwaggerInfo.Description = "TODO API Fiber Swagger docs"
	swagger_docs_v1.SwaggerInfo.Version = "1.0"
	swagger_docs_v1.SwaggerInfo.Host = os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT")
	swagger_docs_v1.SwaggerInfo.BasePath = "/api"
	swagger_docs_v1.SwaggerInfo.Schemes = []string{"http", "https"}
}

func SwaggerConfig() swagger.Config {
	return swagger.Config{
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "list",

		// Prefill OAuth ClientId on Authorize popup
		OAuth: &swagger.OAuthConfig{
			AppName:  "OAuth Provider",
			ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		},
		// Ability to change OAuth2 redirect uri location
		OAuth2RedirectUrl: "http://" + os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT") + "/swagger/oauth2-redirect.html",
	}
}
