package utils

import (
	"log/slog"
	"subs_service/internal/entities"
	"time"

	"github.com/google/uuid"
)

func ParseDate(s string) (time.Time, error) {
	date, err := time.Parse("01-2006", s)
	if err != nil {
		slog.Error("Wrong start or end date. Pattern: 01-2006. ", "Errror", err.Error())
		return time.Now(), err
	}
	return date, nil
}

func ParseFrontendSub(fs entities.FrontendSubscription) (entities.Subscription, error) {
	startDate, err := ParseDate(fs.StartDate)
	if err != nil {
		return entities.Subscription{}, err
	}

	var endDate *time.Time
	if fs.EndDate != "" {
		endDateTime, err := ParseDate(fs.EndDate)
		if err != nil {
			return entities.Subscription{}, err
		}

		endDate = &endDateTime
	}

	userID, err := uuid.Parse(fs.UserID)
	if err != nil {
		slog.Error("Incorrenct user id", "error", err.Error())
		return entities.Subscription{}, err
	}

	return entities.Subscription{
		ID:        fs.ID,
		Name:      fs.Name,
		Price:     fs.Price,
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}, nil
}

func CountMonths(nachalo time.Time, konec time.Time) int {
	y1, m1, _ := nachalo.Date()
	y2, m2, _ := konec.Date()

	return (y2-y1)*12 + int(m2-m1) + 1
}
