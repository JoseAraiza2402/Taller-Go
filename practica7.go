package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Modelo de datos
type Usuario struct {
	ID     int    `json:"id"`
	Nombre string `json:"name"`
	Email  string `json:"email"`
}

// Servicio de usuarios
type UsuarioService struct {
	usuarios []Usuario
}

func NewUsuarioService() *UsuarioService {
	return &UsuarioService{
		usuarios: []Usuario{
			{ID: 1, Nombre: "Alfredo", Email: "Alfredo@mail.com"},
		},
	}
}

func (s *UsuarioService) GetAll() []Usuario {
	return s.usuarios
}

func (s *UsuarioService) Create(user Usuario) Usuario {
	user.ID = len(s.usuarios) + 1
	s.usuarios = append(s.usuarios, user)
	return user
}

func (s *UsuarioService) Update(id int, updated Usuario) (*Usuario, bool) {
	for i, user := range s.usuarios {
		if user.ID == id {
			updated.ID = id
			s.usuarios[i] = updated
			return &updated, true
		}
	}
	return nil, false
}

func (s *UsuarioService) Delete(id int) bool {
	for i, user := range s.usuarios {
		if user.ID == id {
			s.usuarios = append(s.usuarios[:i], s.usuarios[i+1:]...)
			return true
		}
	}
	return false
}

// Middleware de API Key
func apiKeyMiddleware(c *gin.Context) {
	apiKey := c.GetHeader("X-API-Key")
	if apiKey != "my-secret-key" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API Key inválida"})
		return
	}
	c.Next()
}

func main() {
	r := gin.Default()

	// Logger de solicitudes
	r.Use(func(c *gin.Context) {
		fmt.Printf("Solicitud recibida: %s %s\n", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})

	usuarioService := NewUsuarioService()

	v1 := r.Group("/v1/users")
	v1.Use(apiKeyMiddleware)
	{
		v1.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, usuarioService.GetAll())
		})

		v1.POST("", func(c *gin.Context) {
			var user Usuario
			if err := c.ShouldBindJSON(&user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error parseando el JSON: %v", err)})
				return
			}
			created := usuarioService.Create(user)
			c.JSON(http.StatusOK, created)
		})

		v1.PUT("/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
				return
			}
			var updated Usuario
			if err := c.ShouldBindJSON(&updated); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error parseando el JSON: %v", err)})
				return
			}
			if res, ok := usuarioService.Update(id, updated); ok {
				c.JSON(http.StatusOK, res)
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
			}
		})

		v1.DELETE("/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
				return
			}
			if usuarioService.Delete(id) {
				c.Status(http.StatusNoContent)
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
			}
		})
	}

	fmt.Println("Servidor escuchando en el puerto 3000")
	r.Run(":3000")
}
