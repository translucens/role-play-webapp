package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/gin-gonic/gin"
	"github.com/mittz/role-play-webapp/webapp/database"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var dbHandler database.DatabaseHandler

// getUserID always returns the same id which is for scstore user.
// This should be updated when this app implements an authentication feature.
func getUserID() int {
	return 2
}

func postInitEndpoint(c *gin.Context) {
	if err := dbHandler.InitDatabase(); err != nil {
		c.String(http.StatusInternalServerError, "%v", err)
		return
	}

	c.String(http.StatusAccepted, "Initiatilized data.")
}

func getCheckoutsEndpoint(c *gin.Context) {
	userID := getUserID()
	checkouts, err := dbHandler.GetCheckouts(c.Request.Context(), userID)
	if err != nil {
		c.String(http.StatusInternalServerError, "%v", err)
		return
	}

	c.HTML(http.StatusOK, "checkouts.html", gin.H{
		"title":     "Checkouts",
		"checkouts": checkouts,
	})
}

func postCheckoutEndpoint(c *gin.Context) {
	userID := getUserID()

	productID, err := strconv.Atoi(c.PostForm("product_id"))
	if err != nil {
		c.String(http.StatusInternalServerError, "%v", err)
		return
	}

	productQuantity, err := strconv.Atoi(c.PostForm("product_quantity"))
	if err != nil {
		c.String(http.StatusInternalServerError, "%v", err)
		return
	}

	createdAt, err := dbHandler.CreateCheckout(userID, productID, productQuantity)
	if err != nil {
		c.String(http.StatusInternalServerError, "%v", err)
		return
	}

	product, err := dbHandler.GetProduct(productID)
	if err != nil {
		c.String(http.StatusInternalServerError, "%v", err)
		return
	}

	c.HTML(http.StatusAccepted, "checkout.html", gin.H{
		"title": "Checkout",
		"checkout": database.Checkout{
			Product:         product,
			ProductQuantity: productQuantity,
			CreatedAt:       createdAt,
		},
	})
}

func getProductEndpoint(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		c.String(http.StatusInternalServerError, "%v", err)
		return
	}

	product, err := dbHandler.GetProduct(productID)
	if err != nil {
		c.String(http.StatusInternalServerError, "%v", err)
		return
	}

	c.Header("Cache-Control", "max-age=300")
	c.HTML(http.StatusOK, "product.html", gin.H{
		"title":   "Product",
		"product": product,
	})
}

func getProductsEndpoint(c *gin.Context) {
	products, err := dbHandler.GetProducts()
	if err != nil {
		c.String(http.StatusInternalServerError, "%v", err)
		return
	}

	c.Header("Cache-Control", "max-age=300")
	c.HTML(http.StatusOK, "products.html", gin.H{
		"title":    "Products",
		"products": products,
	})
}

func getHealthCheck(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func SetupRouter(dbh database.DatabaseHandler, assetsDir string, templatesDirMatch string) *gin.Engine {
	dbHandler = dbh

	// Create exporter.
	ctx := context.Background()
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	exporter, err := texporter.New(texporter.WithProjectID(projectID))
	if err != nil {
		log.Fatalf("texporter.NewExporter: %v", err)
	}

	// Create trace provider with the exporter.
	//
	// By default it uses AlwaysSample() which samples all traces.
	// In a production environment or high QPS setup please use
	// probabilistic sampling.
	// Example:
	//   tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.0001)), ...)
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
	defer tp.ForceFlush(ctx) // flushes any pending spans
	otel.SetTracerProvider(tp)

	router := gin.Default()
	router.Use(otelgin.Middleware("scstore"))

	router.Static("/assets", assetsDir)
	router.StaticFile("/favicon.ico", filepath.Join(assetsDir, "favicon.ico"))
	router.LoadHTMLGlob(templatesDirMatch)

	router.GET("/", getProductsEndpoint)

	router.POST("/admin/init", postInitEndpoint)

	router.GET("/product/:product_id", getProductEndpoint)
	router.GET("/products", getProductsEndpoint)
	router.GET("/checkouts", getCheckoutsEndpoint)
	router.POST("/checkout", postCheckoutEndpoint)

	router.GET("/healthcheck", getHealthCheck)

	return router
}
