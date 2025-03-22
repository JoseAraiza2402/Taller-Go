package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Usuario struct {
	ID     int    `json:"id"`
	Nombre string `json:"name"`
	Email  string `json:"email"`
}

var usuarios []Usuario

func main() {
	usuarios = append(usuarios, Usuario{
		ID:     1,
		Nombre: "Alfredo",
		Email:  "Alfredo@mail.com",
	})

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		fmt.Printf("Solicitud recibida: %s %s\n", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})

	r.GET("/", func(c *gin.Context) {
		content, err := os.ReadFile("./public/index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "error leyendo el html")
			return
		}
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, string(content))
	})

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	v1 := r.Group("/v1/users")
	{
		v1.GET("", func(c *gin.Context) {
			c.Header("test-header", "header")
			c.JSON(http.StatusOK, usuarios)
		})

		v1.POST("", func(c *gin.Context) {
			var user Usuario
			if err := c.ShouldBindJSON(&user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error parseando el JSON: %v", err)})
				return
			}
			user.ID = len(usuarios) + 1
			usuarios = append(usuarios, user)
			c.Header("test-header", "header")
			c.JSON(http.StatusOK, user)
		})

		v1.PUT("/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
				return
			}

			var updatedUser Usuario
			if err := c.ShouldBindJSON(&updatedUser); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error parseando el JSON: %v", err)})
				return
			}
			fmt.Printf("Body recibido en PUT: %+v\n", updatedUser) // Depuración
			updatedUser.ID = id

			for i, user := range usuarios {
				if user.ID == id {
					usuarios[i] = updatedUser
					c.Header("test-header", "header")
					c.JSON(http.StatusOK, updatedUser)
					return
				}
			}
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		})

		v1.DELETE("/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
				return
			}

			for i, user := range usuarios {
				if user.ID == id {
					usuarios = append(usuarios[:i], usuarios[i+1:]...)
					c.Status(http.StatusNoContent)
					return
				}
			}
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		})
	}

	fmt.Println("Servidor escuchando en el puerto 3000")
	r.Run(":3000")
}
