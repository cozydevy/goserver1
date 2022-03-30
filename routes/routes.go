package routes

import (
	"course-go/config"
	"course-go/controllers"
	"course-go/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func Serve(r *gin.Engine) {
	db := config.GetDB()
	v1 := r.Group("/api/v1")
	authenticate := middleware.Authenticate().MiddlewareFunc()
	authorize := middleware.Authorize()

	authGroup := v1.Group("auth")
	authController := controllers.Auth{DB: db}
	{
		authGroup.POST("/sign-up", authController.Signup)
		authGroup.POST("/sign-in", middleware.Authenticate().LoginHandler)
		authGroup.GET("/profile", authenticate, authController.GetProfile)
		authGroup.PATCH("/profile", authenticate, authController.UpdateProfile)
	}

	usersController := controllers.Users{DB: db}
	usersGroup := v1.Group("users")
	usersGroup.GET("", usersController.FindAll)
	usersGroup.POST("", usersController.Create)
	usersGroup.GET("/:id", usersController.FindOne)
	usersGroup.PATCH("/:id", usersController.Update)
	usersGroup.PATCH("/updateuser/:id", usersController.UpdatebyAdmin)

	usersGroup.DELETE("/:id", usersController.Delete)
	usersGroup.PATCH("/:id/promote", usersController.Promote)
	usersGroup.PATCH("/:id/demote", usersController.Demote)
	usersGroup.Use(authenticate, authorize)
	{

	}

	articleController := controllers.Articles{DB: db}
	articlesGroup := v1.Group("/articles")
	articlesGroup.GET("", articleController.FindAll)
	articlesGroup.GET("/:id", articleController.FindOne)
	articlesGroup.Use(authenticate, authorize)
	{
		articlesGroup.PATCH("/:id", articleController.Update)
		articlesGroup.DELETE("/:id", articleController.Delete)
		articlesGroup.POST("", authenticate, articleController.Create)
	}

	CategoryController := controllers.Categories{DB: db}
	categoriesGroup := v1.Group("/categories")
	categoriesGroup.GET("", CategoryController.FindAll)
	categoriesGroup.GET("/:id", CategoryController.FindOne)
	categoriesGroup.Use(authenticate, authorize)
	{
		categoriesGroup.PATCH("/:id", CategoryController.Update)
		categoriesGroup.DELETE("/:id", CategoryController.Delete)
		categoriesGroup.POST("", CategoryController.Create)
	}

	ProgressController := controllers.Progresses{DB: db}
	progressGroup := v1.Group("/progresses")
	// progressGroup.GET("", ProgressController.FindAll)
	progressGroup.GET("/:id", ProgressController.FindOneuser)
	progressGroup.GET("/mycourse/:id", ProgressController.FindMyCourse)
	progressGroup.GET("/finduserenroll/:id", ProgressController.FinduserEnroll)

	// coursesGroup.Use(authenticate, authorize)

	{
		progressGroup.POST("", ProgressController.Create)
	}

	ProgressDetailController := controllers.ProgressDetails{DB: db}
	progressDetailGroup := v1.Group("/progressesdetail")
	// progressGroup.GET("", ProgressController.FindAll)
	progressDetailGroup.GET("/:id", ProgressDetailController.FindOneuser)

	progressDetailGroup.POST("/couters", ProgressDetailController.CountQuestion)
	progressDetailGroup.POST("/couters2", ProgressDetailController.CountQuestionStatic)

	progressDetailGroup.POST("", ProgressDetailController.Create)

	TaskController := controllers.Tasks{DB: db}
	tasksGroup := v1.Group("/tasks")
	tasksGroup.GET("", TaskController.FindAll)
	tasksGroup.GET("/:id", TaskController.FindOne)

	// tasksGroup.Use(authenticate, authorize)
	{
		tasksGroup.PATCH("/:id", TaskController.Update)
		tasksGroup.DELETE("/:id", TaskController.Delete)
		tasksGroup.POST("", TaskController.Create)
	}

	CourseController := controllers.Courses{DB: db}
	coursesGroup := v1.Group("/courses")
	coursesGroup.GET("", CourseController.FindAll)
	coursesGroup.GET("/:id", CourseController.FindOne)
	// coursesGroup.Use(authenticate, authorize)
	{
		coursesGroup.PATCH("/:id", CourseController.Update)
		coursesGroup.DELETE("/:id", CourseController.Delete)
		coursesGroup.POST("", CourseController.Create)
	}

	QuestionController := controllers.Questions{DB: db}
	questionsGroup := v1.Group("/questions")
	questionsGroup.GET("", QuestionController.FindAll)
	questionsGroup.GET("/:id", QuestionController.FindOne)
	questionsGroup.POST("/CheckAns", QuestionController.CheckAns)
	questionsGroup.POST("/questionall", QuestionController.FindAllName)
	questionsGroup.POST("/couters", QuestionController.FindQuestionByCourse)

	// questionsGroup.Use(authenticate, authorize)
	{
		questionsGroup.PATCH("/:id", QuestionController.Update)
		questionsGroup.DELETE("/:id", QuestionController.Delete)
		questionsGroup.POST("", QuestionController.Create)
	}

	dockersGroup := v1.Group("/containers")
	dockerController := controllers.Dockers{}
	{
		dockersGroup.GET("", dockerController.ListAll)
		dockersGroup.GET("/stop/:id", dockerController.StopContainer)
		dockersGroup.GET("/start/:id", dockerController.StartContainer)
	}

	{
		url := ginSwagger.URL("http://localhost:5200/swagger/doc.json") // The url pointing to API definition
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}

}
