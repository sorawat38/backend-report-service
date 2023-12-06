package service

import "time"

type Report interface {
	GenerateMontlyReport(date time.Time) error
}
