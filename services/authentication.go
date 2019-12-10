package services

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/sample-full-api/models"
	"github.com/sample-full-api/utils"
	"github.com/sample-full-api/views"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

type AuthenticationService interface {
	GenerateToken(loginRequest views.LoginRequest) (string, *utils.Error)
}

type authenticationService struct {
	db             *gorm.DB
	logger         *log.Logger
	expTimeMinutes int
}

func NewAuthenticationService(deps utils.Dependencies) AuthenticationService {
	expTime, _ := strconv.Atoi(os.Getenv("JWT_EXP_TIME_MINUTES"))

	return &authenticationService{
		db:             deps.Db,
		logger:         deps.Logger,
		expTimeMinutes: expTime,
	}
}

func (a *authenticationService) GenerateToken(loginRequest views.LoginRequest) (string, *utils.Error) {
	// 1. check if user, pass exists in DB (if not --> unauthorized)
	if err := a.validateUser(loginRequest.User, loginRequest.Password); err != nil {
		a.logger.Error(err)
		return "", utils.ErrorUnauthorized(err.Error())
	}

	// 2. generate token with expiration time and secret
	token, err := a.generateToken()
	if err != nil {
		a.logger.Error(err)
		return "", utils.ErrorInternal(err.Error())
	}

	return token, nil
}

func (a *authenticationService) validateUser(user, password string) error {
	if err := a.db.Where(&models.User{User: user, Password: password}).First(&models.User{}).Error; err != nil {
		return err
	}

	return nil
}

func (a *authenticationService) generateToken() (string, error) {
	expirationTime := time.Now().Add(time.Duration(a.expTimeMinutes) * time.Minute)

	claims := jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("environment-var"))
	if err != nil {
		return "", err
	}

	return token, nil
}
