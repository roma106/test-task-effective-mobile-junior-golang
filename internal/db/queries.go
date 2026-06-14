package db

import (
	"fmt"
	"log/slog"
	"subs_service/internal/config"
	"subs_service/internal/entities"
	"subs_service/internal/utils"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectToDB(cfg *config.Config) (*sqlx.DB, error) {
	slog.Info(fmt.Sprintf("Connecting to database %s...", cfg.PostgresDB))
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB))
	if err != nil {
		slog.Error("Failed to connect to database: " + err.Error())
		return nil, err
	}
	slog.Info("Database " + cfg.PostgresDB + " connected!")
	return db, nil
}

func CreateSubscription(db *sqlx.DB, subs entities.Subscription) error {
	_, err := db.NamedExec(`INSERT INTO "subscriptions" (name, price, user_id, start_date, end_date) VALUES (:name, :price, :user_id, :start_date, :end_date)`, subs)
	if err != nil {
		slog.Error("Failed to insert subscription: ", "Error", err.Error())
		return err
	}
	return nil
}

func GetSubscriptions(db *sqlx.DB) ([]entities.Subscription, error) {
	subs := []entities.Subscription{}
	err := db.Select(&subs, `SELECT * FROM "subscriptions"`)
	if err != nil {
		slog.Error("Failed to select subscriptions: ", "Error", err.Error())
		return []entities.Subscription{}, err
	}
	return subs, nil
}
func GetSubscriptionByID(db *sqlx.DB, id string) (entities.Subscription, error) {
	sub := entities.Subscription{}

	err := db.Get(&sub, `SELECT * FROM "subscriptions" WHERE id = $1`, id)
	if err != nil {
		slog.Error("Failed to select subscription", "error", err.Error())
		return entities.Subscription{}, err
	}

	return sub, nil
}

func UpdateSubscription(db *sqlx.DB, newSub entities.Subscription) error {
	_, err := db.NamedExec(`UPDATE "subscriptions" SET 
			name = :name,
			price = :price,
			user_id = :user_id,
			start_date = :start_date,
			end_date = :end_date
		WHERE id = :id`, newSub)

	if err != nil {
		slog.Error("Failed to update subscription", "error", err.Error())
		return err
	}

	return nil
}

func DeleteSubscription(db *sqlx.DB, id string) error {
	_, err := db.Exec(`DELETE FROM "subscriptions" WHERE id = $1`, id)
	if err != nil {
		slog.Error("Failed to delete subscription", "error", err.Error())
		return err
	}

	return nil
}

func SumPriceWithFilters(db *sqlx.DB, filterName string, filterUser string, periodStart, periodEnd *time.Time) (int, error) {
	subs := []entities.Subscription{}
	query := `SELECT * FROM "subscriptions" WHERE 1=1`
	args := []interface{}{}

	if filterName != "" {
		query += " AND name = ?"
		args = append(args, filterName)
	}

	if filterUser != "" {
		query += " AND user_id = ?"
		args = append(args, filterUser)
	}

	query = db.Rebind(query)

	err := db.Select(&subs, query, args...)
	if err != nil {
		slog.Error("Failed to select subscriptions: ", "Error", err.Error())
		return 0, err
	}
	sum := 0

	for _, sub := range subs {
		nachalo := sub.StartDate
		if periodStart != nil && nachalo.Before(*periodStart) {
			nachalo = *periodStart
		}

		konec := time.Now()
		if sub.EndDate != nil {
			konec = *sub.EndDate
		}

		if periodEnd != nil && konec.After(*periodEnd) {
			konec = *periodEnd
		}

		if konec.Before(nachalo) {
			continue
		}

		sum += sub.Price * utils.CountMonths(nachalo, konec)
	}
	return sum, nil
}
