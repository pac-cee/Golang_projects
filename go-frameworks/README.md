1.GIN basics



.go
package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func main() {
    r := gin.Default()

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "pong",
        })
    })

    // REST API example
    r.POST("/users", createUser)
    r.GET("/users/:id", getUser)

    r.Run(":8080")
}

type User struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

func createUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, user)
}

func getUser(c *gin.Context) {
    id := c.Param("id")
    c.JSON(http.StatusOK, gin.H{"id": id})
}


2.Fiber framework 


package main

import (
    "github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New()

    // Middleware
    app.Use(func(c *fiber.Ctx) error {
        c.Set("X-Custom-Header", "Hello, Fiber!")
        return c.Next()
    })

    // Routes
    app.Get("/", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "message": "Hello, World!",
        })
    })

    // Route groups
    api := app.Group("/api")
    v1 := api.Group("/v1")

    v1.Post("/users", func(c *fiber.Ctx) error {
        var user User
        if err := c.BodyParser(&user); err != nil {
            return err
        }
        return c.Status(201).JSON(user)
    })

    app.Listen(":3000")
}



3.Buffalo framework


package actions

import (
    "github.com/gobuffalo/buffalo"
    "github.com/gobuffalo/buffalo/middleware"
)

func App() *buffalo.App {
    app := buffalo.New(buffalo.Options{
        Env: "development",
    })

    app.Use(middleware.PopTransaction)
    app.Use(middleware.CSRF)

    app.GET("/", HomeHandler)

    api := app.Group("/api")
    api.GET("/users", UsersHandler)
    api.POST("/users", CreateUserHandler)

    return app
}

func HomeHandler(c buffalo.Context) error {
    return c.Render(200, r.JSON(map[string]string{
        "message": "Welcome to Buffalo!",
    }))
}


4.Echo framework


package main

import (
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "net/http"
)

func main() {
    e := echo.New()

    // Middleware
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    // Routes
    e.GET("/", func(c echo.Context) error {
        return c.JSON(http.StatusOK, map[string]string{
            "message": "Hello, Echo!",
        })
    })

    // Route groups
    api := e.Group("/api")
    api.POST("/users", createUser)
    api.GET("/users/:id", getUser)

    e.Logger.Fatal(e.Start(":1323"))
}


5.commo features comparison



// Gin Middleware
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
            return
        }
        c.Next()
    }
}

// Fiber Middleware
func AuthMiddleware(c *fiber.Ctx) error {
    token := c.Get("Authorization")
    if token == "" {
        return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
    }
    return c.Next()
}

// Echo Middleware
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        token := c.Request().Header.Get("Authorization")
        if token == "" {
            return echo.NewHTTPError(401, "unauthorized")
        }
        return next(c)
    }
}


framework setup command 

gin setup

mkdir gin-project
cd gin-project
go mod init gin-project
go get -u github.com/gin-gonic/gin

fiber setup 

mkdir fiber-project
cd fiber-project
go mod init fiber-project
go get github.com/gofiber/fiber/v2

buffalo setup

go get github.com/gobuffalo/buffalo/buffalo
buffalo new my-buffalo-project
cd my-buffalo-project
buffalo dev

echo setup

mkdir echo-project
cd echo-project
go mod init echo-project
go get github.com/labstack/echo/v4
