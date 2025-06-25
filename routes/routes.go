package routes

import (
	"patient-portal/controllers"
	"patient-portal/middlewares"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(
	authController *controllers.AuthController,
	patientController *controllers.PatientController,
) *gin.Engine {
	router := gin.Default()

	// ADD THIS LINE - CORS Middleware
	router.Use(middlewares.CORSMiddleware())

	// Public routes
	public := router.Group("/")
	{
		public.POST("/login", authController.Login)
	}

	// Protected routes (require authentication)
	api := router.Group("/api")
	api.Use(middlewares.AuthMiddleware())
	{
		// Patient routes
		patients := api.Group("/patients")
		{
			patients.GET("", patientController.GetAllPatients)
			patients.GET("/:id", patientController.GetPatientByID)
			patients.PATCH("/:id", patientController.UpdatePatient) // Shared update endpoint

			// Receptionist-only routes
			receptionistRoutes := patients.Group("")
			receptionistRoutes.Use(middlewares.RoleMiddleware("receptionist"))
			{
				receptionistRoutes.POST("", patientController.CreatePatient)
				receptionistRoutes.DELETE("/:id", patientController.DeletePatient)
			}

			// Doctor-only features
			doctorRoutes := patients.Group("/medical-notes")
			doctorRoutes.Use(middlewares.RoleMiddleware("doctor"))
			{
				doctorRoutes.PUT("/:id", patientController.UpdateMedicalNotes)
			}
		}
	}

	// Swagger documentation route
	url := ginSwagger.URL("/swagger/doc.json") // Point to the correct JSON file
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
