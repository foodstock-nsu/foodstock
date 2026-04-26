package http

import (
	jwtinfra "backend/internal/infrastructure/jwt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Router struct {
	tokenGen *jwtinfra.Generator
	Auth     *AuthHandler
	Client   *ClientHandler
	Location *LocationHandler
	Item     *ItemHandler
}

func NewRouter(
	tokenGen *jwtinfra.Generator,
	auth *AuthHandler,
	client *ClientHandler,
	location *LocationHandler,
	item *ItemHandler,
) *Router {
	return &Router{
		tokenGen: tokenGen,
		Auth:     auth,
		Client:   client,
		Location: location,
		Item:     item,
	}
}

func (r *Router) InitRoutes() *echo.Echo {
	router := echo.New()

	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	router.Use(middleware.Recover())
	router.Use(middleware.RequestLogger())

	router.File("/swagger/doc.json", "api/openapi.yaml")

	router.GET("/swagger", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	router.GET("/swagger/*", echoSwagger.WrapHandler)

	v1 := router.Group("/api/v1")
	{
		client := v1.Group("/client")

		clientLocations := client.Group("/locations")
		{
			clientLocations.GET("/:id/catalog", r.Client.GetCatalog)
		}

		admin := v1.Group("/admin")

		auth := admin.Group("/auth")
		{
			auth.POST("", r.Auth.AdminAuth)
		}

		adminLocations := admin.Group("/locations")
		adminLocations.Use(r.withAuth(r.tokenGen))
		{
			adminLocations.POST("", r.Location.Create)
			adminLocations.PUT("/:id", r.Location.Update)
			adminLocations.DELETE("/:id", r.Location.Delete)
			adminLocations.GET("", r.Location.List)
			adminLocations.GET("/:id/qrcode", r.Location.GetQRCode)
		}

		adminItems := admin.Group("/items")
		adminItems.Use(r.withAuth(r.tokenGen))
		{
			adminItems.POST("", r.Item.Create)
			adminItems.PUT("/:id", r.Item.Update)
			adminItems.DELETE("/:id", r.Item.Delete)
			adminItems.GET("", r.Item.List)
		}
	}

	return router
}

// ---------- Middlewares ----------

func (r *Router) withAuth(gen *jwtinfra.Generator) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing auth header")
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			adminID, err := gen.Validate(tokenString)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
			}

			c.Set("admin_id", adminID)

			return next(c)
		}
	}
}
