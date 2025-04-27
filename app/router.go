package app

import (
	"ProductService/db/connector"
	"ProductService/services"
	"ProductService/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func runserver() {
	log.Println("Product Service - Backend Service")
	router := mux.NewRouter()
	router.Use(utils.CorsFilter)

	getProdByIdProc := services.NewGetProdById(connector.RedisConnector, connector.PGDBConnector)
	getProductHandler := ProductHandler(getProdByIdProc)
	router.HandleFunc("/products/{id}", getProductHandler.HandleProduct).Methods("GET", "OPTIONS")

	getAllProdProc := services.NewGetAllProd(connector.RedisConnector, connector.PGDBConnector)
	getAllProductHandler := ProductHandler(getAllProdProc)
	router.HandleFunc("/products", getAllProductHandler.HandleProduct).Methods("GET", "OPTIONS")

	createProduct := services.NewCreateProduct(connector.RedisConnector, connector.PGDBConnector)
	createProductHandler := ProductHandler(createProduct)
	router.HandleFunc("/products", createProductHandler.HandleProduct).Methods("POST", "OPTIONS")

	updateProduct := services.NewUpdateProduct(connector.RedisConnector, connector.PGDBConnector)
	updateProductHandler := ProductHandler(updateProduct)
	router.HandleFunc("/products/{id}", updateProductHandler.HandleProduct).Methods("PUT", "OPTIONS")

	deleteProduct := services.NewDeleteProd(connector.RedisConnector, connector.PGDBConnector)
	deleteProductHandler := ProductHandler(deleteProduct)
	router.HandleFunc("/products/{id}", deleteProductHandler.HandleProduct).Methods("DELETE", "OPTIONS")

	PORT := os.Getenv("PORT")

	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: router,
	}

	log.Printf("Started HTTPS Server on port %v", PORT)
	log.Printf("-------------------------")
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Error listening on port: %v, error: %v", PORT, err)
	}

}
