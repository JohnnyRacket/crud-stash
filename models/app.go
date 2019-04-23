package models

import (
	"time"
)

//App reperesents an app the users can register in the system
type App struct {
	Name        string
	Owner       string
	CreatedDate time.Time
}
