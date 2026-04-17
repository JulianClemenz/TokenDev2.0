package main

import (
	"AppFitness/handlers"
	"AppFitness/middleware"
	"AppFitness/repositories"
	"AppFitness/services"
	"fmt"
	"log"
	"net/http"
	"os" // ✅ CORRECCIÓN: Para leer variables de entorno

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv" // ✅ CORRECCIÓN: Librería para .env
)

func main() {
	// ✅ CORRECCIÓN: Cargamos el archivo .env al inicio
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: No se encontró archivo .env, usando variables de sistema")
	}

	fmt.Println("Iniciando AppFitness...")

	// 1. Conexión a la Base de Datos
	db := repositories.NewMongoDB()

	defer func() {
		log.Println("Cerrando conexion con MongoDB...")
		if err := db.Disconnect(); err != nil {
			log.Fatalf("Error al desconectar MongoDB: %v", err)
		}
	}()
	log.Println("Conectado a MongoDB exitosamente.")

	// 2. Dependencias (Se mantienen tus variables originales)
	userRepo := repositories.NewUserRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)
	exerciseRepo := repositories.NewExcerciseRepository(db)
	routineRepo := repositories.NewRoutineRepository(db)
	workoutRepo := repositories.NewWorkoutRepository(db)

	authService := services.NewAuthService(userRepo, sessionRepo)
	userService := services.NewUserService(userRepo)
	exerciseService := services.NewExcerciseService(exerciseRepo)
	routineService := services.NewRoutineService(routineRepo, exerciseRepo)
	workoutService := services.NewWorkoutService(workoutRepo, routineRepo, userRepo)
	adminService := services.NewAdminService(userRepo, exerciseRepo, routineRepo, sessionRepo)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	exerciseHandler := handlers.NewExerciseHandler(exerciseService)
	routineHandler := handlers.NewRoutineHandler(routineService)
	workoutHandler := handlers.NewWorkoutHadler(workoutService)
	adminHandler := handlers.NewAdminHandler(adminService)

	router := gin.Default()

	// Configurar archivos estáticos y templates
	router.Static("/statics", "./statics")
	router.LoadHTMLGlob("templates/*")

	// --- RUTAS HTML (Sin cambios) ---
	router.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", nil) })
	router.GET("/login", func(c *gin.Context) { c.HTML(http.StatusOK, "login.html", nil) })
	router.GET("/register", func(c *gin.Context) { c.HTML(http.StatusOK, "register1.html", nil) })
	router.GET("/register2", func(c *gin.Context) { c.HTML(http.StatusOK, "register2.html", nil) })
	router.GET("/dashboard-user", func(c *gin.Context) { c.HTML(http.StatusOK, "user-dashboard.html", nil) })
	router.GET("/dashboard-admin", func(c *gin.Context) { c.HTML(http.StatusOK, "admin-dashboard.html", nil) })
	router.GET("/user-exercise", func(c *gin.Context) { c.HTML(http.StatusOK, "user-exercise.html", nil) })
	router.GET("/user-routines", func(c *gin.Context) { c.HTML(http.StatusOK, "user-routines.html", nil) })
	router.GET("/user-progress", func(c *gin.Context) { c.HTML(http.StatusOK, "user-progress.html", nil) })
	router.GET("/user-record", func(c *gin.Context) { c.HTML(http.StatusOK, "user-record.html", nil) })
	router.GET("/user-routines-new", func(c *gin.Context) { c.HTML(http.StatusOK, "user-routines-new.html", nil) })
	router.GET("/user-routine-view.html", func(c *gin.Context) { c.HTML(http.StatusOK, "user-routine-view.html", nil) })
	router.GET("/user-routine-edit.html", func(c *gin.Context) { c.HTML(http.StatusOK, "user-routine-edit.html", nil) })
	router.GET("/profile", func(c *gin.Context) { c.HTML(http.StatusOK, "profile.html", nil) })
	router.GET("/profile-edit.html", func(c *gin.Context) { c.HTML(http.StatusOK, "profile-edit.html", nil) })
	router.GET("/profile-edit-password.html", func(c *gin.Context) { c.HTML(http.StatusOK, "profile-edit-password.html", nil) })
	router.GET("/admin-exercises", func(c *gin.Context) { c.HTML(http.StatusOK, "admin-exercises.html", nil) })
	router.GET("/admin-users", func(c *gin.Context) { c.HTML(http.StatusOK, "admin-users.html", nil) })
	router.GET("/admin-stats", func(c *gin.Context) { c.HTML(http.StatusOK, "admin-stats.html", nil) })
	router.GET("/admin-system-logs", func(c *gin.Context) { c.HTML(http.StatusOK, "admin-system-logs.html", nil) })
	router.GET("/admin-excercise-new", func(c *gin.Context) { c.HTML(http.StatusOK, "admin-excercise-new.html", nil) })
	router.GET("/admin-excercise-edit", func(c *gin.Context) { c.HTML(http.StatusOK, "admin-excercise-edit.html", nil) })
	router.GET("/admin-exercise-ranking", func(c *gin.Context) { c.HTML(http.StatusOK, "admin-exercise-ranking.html", nil) })

	// --- RUTAS API PÚBLICAS ---
	router.POST("/register", userHandler.PostUser)
	router.POST("/login", authHandler.PostLogin)
	router.POST("/logout", authHandler.PostLogout)
	router.POST("/refresh", authHandler.PostRefresh)

	// --- GRUPO API PROTEGIDO ---
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		userRoutes := api.Group("/users")
		{
			userRoutes.GET("/:id", userHandler.GetUserByID)
			userRoutes.PUT("/:id", userHandler.PutUser)
			userRoutes.POST("/:id/password", userHandler.PasswordModify)
		}

		exerciseRoutes := api.Group("/exercises")
		{
			exerciseRoutes.GET("/", exerciseHandler.GetExcercises)
			exerciseRoutes.GET("/filter", exerciseHandler.GetByFilters)
			exerciseRoutes.GET("/:id", exerciseHandler.GetExcerciseByID)

			adminExercise := exerciseRoutes.Group("/")
			adminExercise.Use(middleware.CheckAdmin()) // ✅ Usando tu middleware original
			{
				adminExercise.POST("/", exerciseHandler.PostExcercise)
				adminExercise.PUT("/:id", exerciseHandler.PutExcercise)
				adminExercise.DELETE("/:id", exerciseHandler.DeleteExcercise)
			}
		}

		routineRoutes := api.Group("/routines")
		routineRoutes.Use(middleware.CheckUser()) // ✅ Usando tu middleware original
		{
			routineRoutes.POST("/", routineHandler.PostRoutine)
			routineRoutes.GET("/", routineHandler.GetRoutines)
			routineRoutes.GET("/:id", routineHandler.GetRoutineByID)
			routineRoutes.PUT("/:id", routineHandler.PutRoutine)
			routineRoutes.DELETE("/:id", routineHandler.DeleteRoutine)
			routineRoutes.POST("/:id/exercises", routineHandler.AddExcerciseToRoutine)
			routineRoutes.PUT("/:id/exercises/:exercise_id", routineHandler.UpdateExerciseInRoutine)
			routineRoutes.DELETE("/exercises", routineHandler.RemoveExerciseFromRoutine)
		}

		workoutRoutes := api.Group("/workouts")
		workoutRoutes.Use(middleware.CheckUser())
		{
			workoutRoutes.GET("/", workoutHandler.GetWorkouts)
			workoutRoutes.POST("/:id_routine", workoutHandler.PostWorkout)
			workoutRoutes.GET("/stats", workoutHandler.GetWorkoutStats)
			workoutRoutes.GET("/:id", workoutHandler.GetWorkoutByID)
			workoutRoutes.DELETE("/:id", workoutHandler.DeleteWorkout)
		}

		adminRoutes := api.Group("/admin")
		adminRoutes.Use(middleware.CheckAdmin())
		{
			adminRoutes.GET("/users", userHandler.GetUsers)
			adminRoutes.GET("/stats/users", adminHandler.GetLogs)
			adminRoutes.GET("/stats/exercises", adminHandler.GetGlobalStats)
		}
	}

	// 5. Iniciar Servidor
	// ✅ CORRECCIÓN: Puerto dinámico desde el .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor escuchando en http://localhost:%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor Gin: %v", err)
	}
}
