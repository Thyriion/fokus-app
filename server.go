package main

import (
	"log"
	"net/http"
	"os"

	"fokus-app/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("🚀 Server läuft unter http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
