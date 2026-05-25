package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	redisstore "github.com/gin-contrib/sessions/redis"
	"github.com/rlapenok/exchanger/internal/config"
	"github.com/rlapenok/exchanger/internal/infra/postgres"
	"github.com/rlapenok/exchanger/internal/logger"
	httptransport "github.com/rlapenok/exchanger/internal/transport/http"
	"github.com/rlapenok/exchanger/internal/uc/auth"
	ucaction "github.com/rlapenok/exchanger/internal/uc/action"
	"github.com/rlapenok/exchanger/internal/uc/currency"
	"github.com/rlapenok/exchanger/internal/uc/exchange"
	"github.com/rlapenok/exchanger/internal/uc/exchangerate"
	"github.com/rlapenok/exchanger/internal/uc/user"
)

const (
	sessionName   = "exchanger_session"
	sessionSecret = "exchanger-local-session-secret-32b"
	sessionMaxAge = 86400 * 7
)

func Run() error {
	// load config
	config, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// redefine logger
	logger.RedefineLogger(config.Logger)

	// create new postgres
	db, err := postgres.NewPostgres(context.Background(), config.DB)
	if err != nil {
		return fmt.Errorf("failed to create postgres: %w", err)
	}
	defer db.Close()

	// wire dependencies
	userRepo := postgres.NewUserRepo(db)
	actionRepo := postgres.NewActionRepo(db)
	loginUseCase := auth.NewLoginUseCase(userRepo)
	accountUseCase := user.NewAccountUseCase(userRepo)
	recordActionUseCase := ucaction.NewRecordUseCase(actionRepo)
	listSessionActionsUseCase := ucaction.NewListSessionUseCase(actionRepo)
	authHandler := httptransport.NewAuthHandler(loginUseCase)
	accountHandler := httptransport.NewAccountHandler(accountUseCase, listSessionActionsUseCase)
	currencyRepo := postgres.NewCurrencyRepo(db)
	getAllCurrenciesUseCase := currency.NewGetAllCurrenciesUseCase(currencyRepo)
	getByCodeUseCase := currency.NewGetByCodeUseCase(currencyRepo)
	createCurrencyUseCase := currency.NewCreateUseCase(currencyRepo)
	updateCurrencyUseCase := currency.NewUpdateUseCase(currencyRepo)
	deleteCurrencyUseCase := currency.NewDeleteUseCase(currencyRepo)
	currencyHandler := httptransport.NewCurrencyHandler(
		getAllCurrenciesUseCase,
		getByCodeUseCase,
		createCurrencyUseCase,
		updateCurrencyUseCase,
		deleteCurrencyUseCase,
	)
	exchangeRateRepo := postgres.NewExchangeRateRepo(db)
	listLiveExchangeRatesUseCase := exchangerate.NewListLiveUseCase(exchangeRateRepo)
	createExchangeRateUseCase := exchangerate.NewCreateUseCase(exchangeRateRepo)
	updateExchangeRateUseCase := exchangerate.NewUpdateUseCase(exchangeRateRepo)
	listExchangeRateHistoryUseCase := exchangerate.NewListHistoryUseCase(exchangeRateRepo)
	listExchangeRateReportUseCase := exchangerate.NewListReportUseCase(exchangeRateRepo)
	exchangeRepo := postgres.NewExchangeRepo(db)
	executeExchangeUseCase := exchange.NewExecuteUseCase(exchangeRateRepo, exchangeRepo)
	listExchangeReportUseCase := exchange.NewListReportUseCase(exchangeRepo)
	exchangeRateHandler := httptransport.NewExchangeRateHandler(
		listLiveExchangeRatesUseCase,
		createExchangeRateUseCase,
		updateExchangeRateUseCase,
		listExchangeRateHistoryUseCase,
		listExchangeRateReportUseCase,
		listExchangeReportUseCase,
	)
	exchangeHandler := httptransport.NewExchangeHandler(executeExchangeUseCase)

	sessionStore, err := newSessionStore(config)
	if err != nil {
		return fmt.Errorf("failed to create session store: %w", err)
	}

	// create new router
	router := httptransport.NewRouter(httptransport.RouterDeps{
		Handlers: httptransport.Handlers{
			Auth:         authHandler,
			Account:      accountHandler,
			Currency:     currencyHandler,
			ExchangeRate: exchangeRateHandler,
			Exchange:     exchangeHandler,
		},
		SessionStore: sessionStore,
		SessionName:  sessionName,
		RecordAction: recordActionUseCase,
	})

	return router.Run(":8080")
}

func newSessionStore(config *config.Config) (sessions.Store, error) {
	address := fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port)
	store, err := redisstore.NewStoreWithDB(
		10,
		"tcp",
		address,
		"",
		config.Redis.Password,
		config.Redis.DB,
		[]byte(sessionSecret),
	)
	if err != nil {
		return nil, err
	}

	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   sessionMaxAge,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	return store, nil
}
