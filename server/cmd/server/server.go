package server

import (
	"log"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/graph"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/config"
)

const defaultPort = "8080"

func Start() {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Infrastructure
	config.ConnectDatabase()
	config.ConnectRedis()

	// Fiber app
	app := fiber.New()

	// GraphQL server (gqlgen)
	srv := handler.New(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{
					DB: config.DB, // pass DB if your resolver has it
				},
			},
		),
	)

	// GraphQL config (same as your net/http version)
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New)

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New,
	})

	// Routes
	app.All("/query", adaptor.HTTPHandler(srv))
	app.Get("/", adaptor.HTTPHandler(playground.Handler("GraphQL Playground", "/query")))

	// Port
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Printf("GraphQL playground → http://localhost:%s/", port)
	log.Fatal(app.Listen(":" + port))
}