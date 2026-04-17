package main

import (
	"AppFitness/handlers"
	"AppFitness/middleware"
	"AppFitness/repositories"
	"AppFitness/services"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
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

	// 2. Dependencias

	// --- Repositorios ---
	userRepo := repositories.NewUserRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)
	exerciseRepo := repositories.NewExcerciseRepository(db)
	routineRepo := repositories.NewRoutineRepository(db)
	workoutRepo := repositories.NewWorkoutRepository(db)

	// --- Servicios ---
	authService := services.NewAuthService(userRepo, sessionRepo)
	userService := services.NewUserService(userRepo)
	exerciseService := services.NewExcerciseService(exerciseRepo)
	routineService := services.NewRoutineService(routineRepo, exerciseRepo)
	workoutService := services.NewWorkoutService(workoutRepo, routineRepo, userRepo)
	adminService := services.NewAdminService(userRepo, exerciseRepo, routineRepo, sessionRepo)

	// --- Handlers ---
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	exerciseHandler := handlers.NewExerciseHandler(exerciseService)
	routineHandler := handlers.NewRoutineHandler(routineService)
	workoutHandler := handlers.NewWorkoutHadler(workoutService)
	adminHandler := handlers.NewAdminHandler(adminService)

	router := gin.Default()

	// Configurar archivos státic y templates
	router.Static("/statics", "./statics")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register1.html", nil)
	})
	router.GET("/register2", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register2.html", nil)
	})
	router.GET("/dashboard-user", func(c *gin.Context) {
		c.HTML(http.StatusOK, "user-dashboard.html", nil)
	})
	router.GET("/dashboard-admin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin-dashboard.html", nil)
	})

	// Rutas de páginas de Usuario
	router.GET("/user-exercise", func(c *gin.Context) {
		c.HTML(http.StatusOK, "user-exercise.html", nil)
	})
	router.GET("/user-routines", func(c *gin.Context) {
		c.HTML(http.StatusOK, "user-routines.html", nil)
	})
	router.GET("/user-progress", func(c *gin.Context) {
		c.HTML(http.StatusOK, "user-progress.html", nil)
	})
	router.GET("/user-record", func(c *gin.Context) {
		c.HTML(http.StatusOK, "user-record.html", nil)
	})
	router.GET("/user-routines-new", func(c *gin.Context) {
		c.HTML(http.StatusOK, "user-routines-new.html", nil)
	})
	router.GET("/user-routine-view.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "user-routine-view.html", nil)
	})
	router.GET("/user-routine-edit.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "user-routine-edit.html", nil)
	})

	// Rutas de Perfil
	router.GET("/profile", func(c *gin.Context) {
		c.HTML(http.StatusOK, "profile.html", nil)
	})
	router.GET("/profile-edit.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "profile-edit.html", nil)
	})
	router.GET("/profile-edit-password.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "profile-edit-password.html", nil)
	})

	// Rutas de páginas de Admin
	router.GET("/admin-exercises", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin-exercises.html", nil)
	})
	router.GET("/admin-users", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin-users.html", nil)
	})
	router.GET("/admin-stats", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin-stats.html", nil)
	})
	router.GET("/admin-system-logs", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin-system-logs.html", nil)
	})
	// Rutas para crear/editar ejercicios de admin
	router.GET("/admin-excercise-new", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin-excercise-new.html", nil)
	})
	router.GET("/admin-excercise-edit", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin-excercise-edit.html", nil)
	})
	router.GET("/admin-exercise-ranking", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin-exercise-ranking.html", nil)
	})

	// Rutas Públicas (Autenticación y Registro)
	router.POST("/register", userHandler.PostUser)
	router.POST("/login", authHandler.PostLogin)
	router.POST("/logout", authHandler.PostLogout)
	router.POST("/refresh", authHandler.PostRefresh)

	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())

	//  Rutas de Perfil de Usuario
	userRoutes := api.Group("/users")
	{
		userRoutes.GET("/:id", userHandler.GetUserByID)
		userRoutes.PUT("/:id", userHandler.PutUser)
		userRoutes.POST("/:id/password", userHandler.PasswordModify)
	}
	exerciseRoutes := api.Group("/exercises")
	{
		exerciseRoutes.GET("/", exerciseHandler.GetExcercises)
		exerciseRoutes.GET("/filter", exerciseHandler.GetByFilters) // Búsqueda y filtros
		exerciseRoutes.GET("/:id", exerciseHandler.GetExcerciseByID)

		adminExercise := exerciseRoutes.Group("/")
		adminExercise.Use(middleware.CheckAdmin())
		{
			adminExercise.POST("/", exerciseHandler.PostExcercise)  // Alta
			adminExercise.PUT("/:id", exerciseHandler.PutExcercise) // Edición
			adminExercise.DELETE("/:id", exerciseHandler.DeleteExcercise)
		}
	}

	// Rutas de Rutinas
	routineRoutes := api.Group("/routines")
	routineRoutes.Use(middleware.CheckUser())
	{
		routineRoutes.POST("/", routineHandler.PostRoutine)
		routineRoutes.GET("/", routineHandler.GetRoutines)
		routineRoutes.GET("/:id", routineHandler.GetRoutineByID)
		routineRoutes.PUT("/:id", routineHandler.PutRoutine)
		routineRoutes.DELETE("/:id", routineHandler.DeleteRoutine)

		// Manejo de ejercicios dentro de una rutina
		routineRoutes.POST("/:id/exercises", routineHandler.AddExcerciseToRoutine)
		routineRoutes.PUT("/:id/exercises/:exercise_id", routineHandler.UpdateExerciseInRoutine) // handler espera un DTO en el body, así que no usamos params
		routineRoutes.DELETE("/exercises", routineHandler.RemoveExerciseFromRoutine)
	}

	// Rutas de Seguimiento (Workouts)
	workoutRoutes := api.Group("/workouts")
	workoutRoutes.Use(middleware.CheckUser())
	{
		workoutRoutes.GET("/", workoutHandler.GetWorkouts)

		workoutRoutes.POST("/:id_routine", workoutHandler.PostWorkout)

		workoutRoutes.GET("/stats", workoutHandler.GetWorkoutStats)

		workoutRoutes.GET("/:id", workoutHandler.GetWorkoutByID) // Ver un workout específico

		workoutRoutes.DELETE("/:id", workoutHandler.DeleteWorkout)
	}

	// --- Rutas del Panel de Administración ---
	adminRoutes := api.Group("/admin")
	adminRoutes.Use(middleware.CheckAdmin()) // Protegido solo para Admins
	{
		adminRoutes.GET("/users", userHandler.GetUsers) // Gestión de usuarios
		adminRoutes.GET("/stats/users", adminHandler.GetLogs)
		adminRoutes.GET("/stats/exercises", adminHandler.GetGlobalStats)
	}

	// 5. Iniciar Servidor
	log.Println("Servidor escuchando en http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error al iniciar el servidor Gin: %v", err)
	}
}
