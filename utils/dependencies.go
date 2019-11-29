package utils

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type Dependencies struct {
	Db     *gorm.DB
	Logger *logrus.Logger
}
