package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/sean-miningah/sil-backend-assessment/internal/adapters/handlers/graphql/generated"
	"github.com/sean-miningah/sil-backend-assessment/internal/core/ports"
)

type Handler struct {
	resolver *Resolver
}

func NewHandler(ps ports.ProductService, os ports.OrderService) *Handler {
	return &Handler{
		resolver: NewResolver(ps, os),
	}
}

func (h *Handler) GraphQL() gin.HandlerFunc {
	schema := generated.NewExecutableSchema(
		generated.Config{
			Resolvers: h.resolver,
		},
	)

	gqlHandler := handler.New(schema)

	return func(c *gin.Context) {
		gqlHandler.ServeHTTP(c.Writer, c.Request)
	}
}

func (h *Handler) Playground() gin.HandlerFunc {
	playgroundHandler := playground.Handler("GraphQL", "/graphql")

	return func(c *gin.Context) {
		playgroundHandler.ServeHTTP(c.Writer, c.Request)
	}
}
