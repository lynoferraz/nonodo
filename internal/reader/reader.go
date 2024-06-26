// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// This package is responsible for serving the GraphQL reader API.
package reader

//go:generate go run github.com/99designs/gqlgen generate

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	nonodomodel "github.com/calindra/nonodo/internal/model"
	"github.com/calindra/nonodo/internal/reader/graph"
	"github.com/calindra/nonodo/internal/reader/model"
	"github.com/labstack/echo/v4"
)

// Register the GraphQL reader API to echo.
func Register(e *echo.Echo, nonodomodel *nonodomodel.NonodoModel) {
	resolver := Resolver{model.NewModelWrapper(nonodomodel)}
	config := graph.Config{Resolvers: &resolver}
	schema := graph.NewExecutableSchema(config)
	graphqlHandler := handler.NewDefaultServer(schema)
	playgroundHandler := playground.Handler("GraphQL", "/graphql")
	e.POST("/graphql", func(c echo.Context) error {
		graphqlHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})
	e.GET("/graphql", func(c echo.Context) error {
		playgroundHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})
}
