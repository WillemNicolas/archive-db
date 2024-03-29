package routes

import (
	"archive-api/utils/services"
	"time"

	_ "archive-api/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

func BuildSearchRoutes(app *fiber.App, pool *pgxpool.Pool) {

	search_routes := app.Group("/search")

	search_routes.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	search_routes.Use(cache.New(cache.Config{
		Expiration:   30 * time.Minute,
		CacheControl: true,
		Next: func(c *fiber.Ctx) bool {
			return c.Query("refresh") == "true"
		},
		KeyGenerator: func(c *fiber.Ctx) string {
			return utils.CopyString(c.Path() + string(c.Request().URI().QueryString()))
		},
	}))

	search_routes.Get("/looking", func(c *fiber.Ctx) error {
		return services.QueryExperiment(c, pool)
	})
	search_routes.Get("/publication", func(c *fiber.Ctx) error {
		return services.SearchExperimentForPublication(c, pool)
	})
	search_routes.Get("/", func(c *fiber.Ctx) error {
		return services.SearchExperimentLike(c, pool)
	})

}
