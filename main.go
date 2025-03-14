package main

import (
	"golangapi/config"
	"golangapi/models"
	"golangapi/controllers"
	"log"
	"golangapi/middleware"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error Loading ENV")
	}

	r := gin.Default()
	db := config.ConnectDatabase()

	db.AutoMigrate(&models.User{}, &models.Post{}, &models.Tag{}, &models.PostTag{}, &models.Product{}, &models.Inventory{}, &models.Order{})

	authController := controllers.NewAuthController(db)
	userController := controllers.NewUserController(db)
	postController := controllers.NewPostController(db)
	productController := controllers.NewProductController(db)
	inventoryController := controllers.NewInventoryController(db)
	orderController := controllers.NewOrderController(db)
	sysController := controllers.NewSysController(db)
	// productImageController := controllers.NewProductImageController(db)

	api := r.Group("/api")
	{
		auth := api.Group("/auth") 
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		 // Protected routes
			protected := api.Group("/")
			protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/users", userController.GetUsers)
			protected.POST("/users", userController.CreateUser)

		// 	 // Tag routes
			protected.POST("/tags", postController.CreateTag)

		// 	 // Without DB routes
			protected.POST("/send", controllers.CreateUserWithoutDB)
			protected.GET("/get", controllers.GetUserWithoutDB)

		// 	 // Post Routes
			protected.POST("/post", postController.CreatePost)
			protected.GET("/post", postController.GetPosts)
			protected.GET("/posts/:id", postController.GetPost)
			protected.PUT("/posts/:id", postController.UpdatePost)
			protected.DELETE("/posts/:id", postController.DeletePost)
		
		//	//  Product
			protected.GET("/products", productController.GetProducts)
			protected.GET("/products/:id", productController.GetProductByID)
			protected.GET("/products/category/:category", productController.GetProductsByCategory)
			protected.POST("/products", productController.CreateProduct)
			protected.PUT("/products/:id", productController.UpdateProduct)
			protected.DELETE("/products/:id", productController.DeleteProduct)
			protected.POST("/products/:id/upload", productController.UploadProductImage)
			protected.GET("/products/:id/download", productController.DownloadProductImage)

		//  //  Inventory
			protected.GET("/inventory", inventoryController.GetInventory)
			protected.GET("/inventory/:product_id", inventoryController.GetStock)
			protected.PUT("/inventory/:product_id", inventoryController.UpdateStock)

		//  //  Order
			protected.GET("/orders", orderController.GetOrders)
			protected.GET("/orders/:order_id", orderController.GetOrderByID)
			protected.POST("/orders", orderController.CreateOrder)

		//  //  SYS Routes
			protected.POST("/directory", sysController.CreateDirectory)
			protected.POST("/file", sysController.CreateFile)
			protected.POST("/file/read", sysController.ReadFile)
			protected.PUT("/file/rename", sysController.RenameFile)
			protected.POST("/file/upload", sysController.UploadFile)
			protected.GET("/file/download", sysController.DownloadFile)

		//	//	Upload Product
			// protected.POST("/products/upload-image", productImageController.UploadImage)
			// protected.GET("/products/download-image/:id", productImageController.DownloadImage)
		}
	}

	r.Run(":8080")
}