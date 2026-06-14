package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"subs_service/internal/config"
	"subs_service/internal/db"
	"subs_service/internal/entities"
	"subs_service/internal/utils"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "subs_service/docs"

	echoSwagger "github.com/swaggo/echo-swagger"
)

type App struct {
	db   *sqlx.DB
	http *echo.Echo
	cfg  *config.Config
}

func New() (*App, error) {
	app := new(App)

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLoggerWithConfig(config.LoggerConfig()))
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	cfg, err := config.New("conf.env")
	if err != nil {
		return nil, err
	}

	e.File("/", "frontend/index.html")
	e.File("/script.js", "frontend/script.js")
	e.File("/style.css", "frontend/style.css")

	Db, err := db.ConnectToDB(cfg)
	if err != nil {
		return nil, err
	}
	e.POST("/subs", app.createSubscription)
	e.PUT("/subs/:id", app.updateSubscription)
	e.GET("/subs", app.getSubscriptions)
	e.GET("/subs/:id", app.getSubscription)
	e.DELETE("/subs/:id", app.deleteSubscription)
	e.GET("/subs/sum", app.sumSubsPrices)
	app.db = Db
	app.cfg = cfg
	app.http = e
	return app, nil
}

func (app *App) Run() error {
	err := app.http.Start(fmt.Sprintf(":%s", app.cfg.RestServerPort))
	return err
}

// createSubscription godoc
// @Summary Создать подписку
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body entities.FrontendSubscription true "Данные подписки"
// @Success 200 {string} string "successfully created new sub"
// @Failure 400 {string} string "failed to parse request data"
// @Failure 500 {string} string "failed to create sub"
// @Router /subs [post]
func (app *App) createSubscription(c echo.Context) error {
	subFromFront := entities.FrontendSubscription{}
	if err := c.Bind(&subFromFront); err != nil {
		slog.Error("Bad sub from frontend: ", "Error", err.Error())
		return c.JSON(http.StatusBadRequest, "Error: failed to parse sub")
	}
	sub, err := utils.ParseFrontendSub(subFromFront)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "failed to parse request data")
	}
	err = db.CreateSubscription(app.db, sub)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "failed to create sub")
	}
	return c.String(200, "successfully created new sub")
}

// updateSubscription godoc
// @Summary Обновить подписку
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "ID подписки"
// @Param subscription body entities.FrontendSubscription true "Обновленные данные подпискиЫ"
// @Success 200 {string} string "successfully updated sub"
// @Failure 400 {string} string "failed to parse request data"
// @Failure 500 {string} string "failed to update sub"
// @Router /subs/{id} [put]
func (app *App) updateSubscription(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "failed to parse user id")
	}
	subFromFront := entities.FrontendSubscription{}
	if err := c.Bind(&subFromFront); err != nil {
		slog.Error("Bad sub from frontend: ", "Error", err.Error())
		return c.JSON(http.StatusBadRequest, "Error: failed to parse sub")
	}
	sub, err := utils.ParseFrontendSub(subFromFront)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "failed to parse request data")
	}
	sub.ID = id
	err = db.UpdateSubscription(app.db, sub)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "failed to update sub")
	}
	return c.String(200, "successfully updated sub")
}

// getSubscriptions godoc
// @Summary Получить список подписок
// @Tags subscriptions
// @Produce json
// @Success 200 {array} entities.Subscription
// @Failure 500 {string} string "failed to get sub"
// @Router /subs [get]
func (app *App) getSubscriptions(c echo.Context) error {
	subs, err := db.GetSubscriptions(app.db)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "failed to get sub")
	}
	return c.JSON(200, subs)
}

// getSubscription godoc
// @Summary Получить подписку по ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {object} entities.Subscription
// @Failure 500 {string} string "failed to get sub"
// @Router /subs/{id} [get]
func (app *App) getSubscription(c echo.Context) error {
	sub, err := db.GetSubscriptionByID(app.db, c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "failed to get sub")
	}
	return c.JSON(200, sub)
}

// deleteSubscription godoc
// @Summary Удалить подписку
// @Tags subscriptions
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {string} string "successfully deleted sub"
// @Failure 500 {string} string "Error: Failed to delete subscription"
// @Router /subs/{id} [delete]
func (app *App) deleteSubscription(c echo.Context) error {
	id := c.Param("id")
	err := db.DeleteSubscription(app.db, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error: Failed to delete subscription")
	}
	return c.String(200, "successfully deleted sub")
}

// sumSubsPrices godoc
// @Summary Посчитать суммарную стоимость подписок с фильтрами по имени, uuid пользователя, за выбранный период
// @Tags subscriptions
// @Produce json
// @Param name query string false "Название (имя) подписки"
// @Param user_id query string false "ID пользователя UUID"
// @Param period_start query string false "Начало периода в формате MM-YYYY"
// @Param period_end query string false "Конец периода в формате MM-YYYY"
// @Success 200 {integer} int
// @Failure 400 {string} string "failed to parse period"
// @Failure 500 {string} string "Error: Failed to count summ"
// @Router /subs/sum [get]
func (app *App) sumSubsPrices(c echo.Context) error {
	filterName := c.QueryParam("name")
	filterUserID := c.QueryParam("user_id")
	var periodStart *time.Time
	if c.QueryParam("period_start") != "" {
		t, err := utils.ParseDate(c.QueryParam("period_start"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Error: failed to parse start period")
		}
		periodStart = &t
	}
	var periodEnd *time.Time
	if c.QueryParam("period_end") != "" {
		t, err := utils.ParseDate(c.QueryParam("period_end"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Error: failed to parse start period")
		}
		periodEnd = &t
	}
	sum, err := db.SumPriceWithFilters(app.db, filterName, filterUserID, periodStart, periodEnd)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error: Failed to count summ")
	}
	return c.JSON(200, sum)
}
