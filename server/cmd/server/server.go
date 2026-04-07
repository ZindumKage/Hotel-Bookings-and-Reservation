package server

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/google/uuid"

	"github.com/gofiber/contrib/otelfiber"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/joho/godotenv"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/graph"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/graph/directives"

	application "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/application/booking"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/config"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/database"
)

const defaultPort = "8080"

func Start() {

	// Load env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Infrastructure
	database.ConnectDatabase()
	config.ConnectRedis()
	config.InitServices()


	// Inject Fiber ctx. into Graphql Context
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		ctx := context.WithValue(c.UserContext(), "fiberCtx", c)
		c.SetUserContext(ctx)
		return c.Next()
	})

	app.Use(otelfiber.Middleware())

	// Rate Limiter
	app.Use(limiter.New(limiter.Config{
		Max: 50,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx)error  {
			return c.Status(429).JSON(fiber.Map{
				"error": "Too many requests",
			})
		},
	}))

	app.Use(func(c *fiber.Ctx) error {

	deviceID := c.Cookies("device_id")

	if deviceID == "" {

		deviceID = uuid.New().String()

		c.Cookie(&fiber.Cookie{
			Name:     "device_id",
			Value:    deviceID,
			HTTPOnly: true,
			Secure:   os.Getenv("ENV") == "production",
			Path:     "/",
			MaxAge:   60 * 60 * 24 * 365,
		})
	}

	return c.Next()
})

	// GraphQL server
	srv := handler.New(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{
					UserService:     config.UserService,
					RoomService:     config.RoomService,
					BookingService:  config.BookingService,
					AuditLogService: config.AuditLogService,
				},
				Directives: graph.DirectiveRoot{
					HasRole: directives.HasRole,
				},
			},
		),
	)



	// Transports
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	// Cache
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	// Extensions
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})



	// Booking expiry worker

	expiryWorker := application.NewExpiryWorker(config.BookingService, 1 * time.Minute)
	ctx := context.Background()
	go expiryWorker.Start(ctx)

	// Routes
	app.All("/query", adaptor.HTTPHandler(srv))
	app.Get("/", adaptor.HTTPHandler(playground.Handler("GraphQL Playground", "/query")))

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Printf("GraphQL Playground → http://localhost:%s/", port)
	log.Fatal(app.Listen(":" + port))
}
