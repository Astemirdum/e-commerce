package openapiv2

import (
	"embed"
	_ "embed"
	"github.com/gin-gonic/gin"
	"io/fs"
	"log"
	"net/http"
)

//go:embed auth/v1/auth.swagger.json
var AuthV1Swagger []byte

//go:embed product/v1/product.swagger.json
var ProductV1Swagger []byte

//go:embed order/v1/order.swagger.json
var OrderV1Swagger []byte

type serviceSpec struct {
	Spec []byte
}

//go:embed swagger-ui
var swui embed.FS

func SwagUIHandler(path string) http.Handler {
	assets, err := fs.Sub(swui, "swagger-ui") // swagger-ui
	if err != nil {
		log.Println("error: ")
	}
	return http.StripPrefix(path, http.FileServer(http.FS(assets)))
}

func SwagHandler() gin.HandlerFunc {
	specs := map[string]*serviceSpec{
		"auth":    {Spec: AuthV1Swagger},
		"product": {Spec: ProductV1Swagger},
		"order":   {Spec: OrderV1Swagger},
	}
	return func(c *gin.Context) {
		name := c.Param("name")
		if name == "" {
			c.AbortWithStatus(http.StatusBadRequest)
		}
		svc := specs[name]

		c.Header("Content-Type", "application/json")
		_, _ = c.Writer.Write(svc.Spec)
	}
}
