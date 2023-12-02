package botdetector

import (
	"log"
	"net/http"

	"github.com/AlyonaAg/bot-detector/internal/model"
	"github.com/gin-gonic/gin"
)

func (i *Implementation) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var r = &RegisterRequest{}
		if err := c.ShouldBindJSON(r); err != nil {
			newErrorResponse(c, http.StatusUnauthorized, err)
			return
		}

		id, err := i.repo.CreateUser(model.User{
			Username:  r.Username,
			Password:  r.Password,
			FirstName: r.FirstName,
			LastName:  r.LastName,
			Phone:     r.Phone,
			Email:     r.Email,
		})
		if err != nil {
			log.Printf("ERROR: CreateUser username=%s, err=%v", r.Username, err)
			newErrorResponse(c, http.StatusBadRequest, err)
			return
		}

		token, err := i.auth.GenerateToken(id)
		if err != nil {
			log.Printf("ERROR: GenerateToken username=%s, err=%v", r.Username, err)
			newErrorResponse(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, RegisterResponse{Token: token})
	}
}

type RegisterRequest struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}
