package server

import (
	"auth/pkg/config"
	"github.com/gofiber/fiber/v3"
	_ "github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"sync"
)

const apiStr = "/api/v1"

// Router - сущность маршрутизатора, содержит приватные поля для работы исключительно внутри пакета
type Router struct {
	r   *fiber.App
	mu  sync.Mutex
	cfg *config.Config
}

// New - создает новый роутер для маршрутизации
func New(cfg *config.Config) *fiber.App {
	router := &Router{
		r:   fiber.New(),
		mu:  sync.Mutex{},
		cfg: cfg,
	}

	router.r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods: []string{"GET", "HEAD", "PUT", "POST", "DELETE"},
	}))

	return router.loadEndpoints()
}

func (route *Router) loadEndpoints() *fiber.App {
	api := route.r.Group(apiStr)

	//Эндпоинты tasks
	tasks := api.Group("/auth")
	tasks.Get("/token", route.createToken)
	tasks.Post("/refresh", route.refreshToken)

	return route.r
}
