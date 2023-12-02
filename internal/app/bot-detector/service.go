package botdetector

import (
	"log"

	"github.com/AlyonaAg/bot-detector/internal/model"
	"github.com/gin-gonic/gin"
)

type repo interface {
	CreateUser(user model.User) (int64, error)
	CreateStatic(static model.Static) (int64, error)
	CreateDynamic(static model.Dynamic) (int64, error)
	GetUser(username string) (*model.User, error)
	GetLastCountDynamic(userId int64, count int) ([]*model.Dynamic, error)
}

type authSrv interface {
	GenerateToken(userID int64) (string, error)
	ParseToken(tokenString string) (int64, error)
}

type Implementation struct {
	router *gin.Engine
	repo   repo
	auth   authSrv
}

func NewCheckerServer(repo repo, auth authSrv) *Implementation {
	return &Implementation{
		router: gin.Default(),
		repo:   repo,
		auth:   auth,
	}
}

func (i *Implementation) Start() error {
	i.configureRouter()

	log.Print("starting server.")

	return i.router.Run(":4567")
}

func (i *Implementation) configureRouter() {
	i.router.POST("/api/account/register", i.Register())
	i.router.POST("/mob-api/static-params", i.CreateStatic())
	i.router.POST("/mob-api/dynamic-params", i.CreateDynamic())
	i.router.POST("/api/token/auth", i.Auth())
}
