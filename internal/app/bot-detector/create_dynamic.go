package botdetector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/AlyonaAg/bot-detector/internal/model"
	"github.com/gin-gonic/gin"
)

func (i *Implementation) CreateDynamic() gin.HandlerFunc {
	return func(c *gin.Context) {
		var r = &CreateDynamicRequest{
			MaxDeviceOffs:      -1,
			MinDeviceOffs:      -1,
			MaxDevAcceleration: -1,
			MinDevAcceleration: -1,
			MinLight:           -1,
			MaxLight:           -1,
			HitY:               -1,
			HitX:               -1,
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

		id, err := i.repo.CreateDynamic(model.Dynamic{
			UserID:             idToken,
			MaxDeviceOffs:      r.MaxDeviceOffs,
			MinDeviceOffs:      r.MinDeviceOffs,
			MaxDevAcceleration: r.MaxDevAcceleration,
			MinDevAcceleration: r.MinDevAcceleration,
			MinLight:           r.MinLight,
			MaxLight:           r.MaxLight,
			HitY:               r.HitY,
			HitX:               r.HitX,
		})
		if err != nil {
			log.Printf("ERROR: CreateDynamic err=%v", err)
			newErrorResponse(c, http.StatusBadRequest, err)
			return
		}

		dynamicRequestsCounter += 1
		if dynamicRequestsCounter == CONTER_MODULO {
			dynamicData, err := i.repo.GetLastCountDynamic(idToken, dynamicRequestsCounter)
			if err != nil {

			}
			dataToSend, err := json.Marshal(dynamicData)
			http.Post("http://localhost:5000/dynamic_analyze", "application/json", bytes.NewBuffer(dataToSend))
			dynamicRequestsCounter = 0
		}

		c.JSON(http.StatusOK, CreateDynamicResponse{ID: fmt.Sprintf("%d", id)})
	}
}

type CreateDynamicRequest struct {
	MaxDeviceOffs      float64 `json:"max_device_offs"`
	MinDeviceOffs      float64 `json:"min_device_offs"`
	MaxDevAcceleration float64 `json:"max_dev_acceleration"`
	MinDevAcceleration float64 `json:"min_dev_acceleration"`
	MinLight           float64 `json:"min_light"`
	MaxLight           float64 `json:"max_light"`
	HitY               float64 `json:"hit_y"`
	HitX               float64 `json:"hit_x"`
}

type CreateDynamicResponse struct {
	ID string `json:"id"`
}
