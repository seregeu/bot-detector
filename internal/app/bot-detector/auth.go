package botdetector

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (i *Implementation) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var r = &RegisterRequest{}
		if err := c.ShouldBindJSON(r); err != nil {
			newErrorResponse(c, http.StatusUnauthorized, err)
			return
		}

		user, err := i.repo.GetUser(r.Username)
		if err != nil {
			log.Printf("ERROR: GetUser username=%s, err=%v", r.Username, err)
			newErrorResponse(c, http.StatusBadRequest, err)
			return
		}

		fmt.Printf("%+v\n", user)
		if user.Password != r.Password {
			log.Printf("ERROR: incorrect username or password")
			newErrorResponse(c, http.StatusUnauthorized, errors.New("incorrect username or password"))
			return
		}

		token, err := i.auth.GenerateToken(user.ID)
		if err != nil {
			log.Printf("ERROR: GenerateToken username=%s, err=%v", r.Username, err)
			newErrorResponse(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, RegisterResponse{Token: token})
	}
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
