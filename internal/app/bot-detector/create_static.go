package botdetector

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AlyonaAg/bot-detector/internal/model"
	"github.com/gin-gonic/gin"
)

func (i *Implementation) CreateStatic() gin.HandlerFunc {
	return func(c *gin.Context) {
		var r = &CreateStaticRequest{
			BattaryCharge:  -1,
			BattaryStatus:  -1,
			DataTransStand: -1,
			SimPresence:    -1,
		}
		if err := c.ShouldBindJSON(r); err != nil {
			newErrorResponse(c, http.StatusUnauthorized, err)
			return
		}

		token := c.GetHeader("Authorization")
		fmt.Println("token: ", token[4:])
		idToken, err := i.auth.ParseToken(token[4:])
		if err != nil {
			log.Printf("ERROR: ParseToken err=%v", err)
			newErrorResponse(c, http.StatusInternalServerError, err)
			return
		}

		id, err := i.repo.CreateStatic(model.Static{
			UserID:         idToken,
			BattaryCharge:  r.BattaryCharge,
			BattaryStatus:  r.BattaryStatus,
			DataTransStand: r.DataTransStand,
			SimPresence:    r.SimPresence,
		})
		if err != nil {
			log.Printf("ERROR: CreateStatic err=%v", err)
			newErrorResponse(c, http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, CreateStaticResponse{ID: fmt.Sprintf("%d", id)})
	}
}

type CreateStaticRequest struct {
	BattaryCharge  float64 `json:"battery_charge"`
	BattaryStatus  float64 `json:"battery_status"`
	DataTransStand float64 `json:"data_trans_stand"`
	SimPresence    float64 `json:"sim_presence"`
}

type CreateStaticResponse struct {
	ID string `json:"id"`
}
