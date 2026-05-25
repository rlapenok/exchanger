package http

import (
	"net/http"

	tmpl "github.com/a-h/templ"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	ucaction "github.com/rlapenok/exchanger/internal/uc/action"
	"github.com/rlapenok/exchanger/web/templ"
)
type Handlers struct {
	Auth         *AuthHandler
	Account      *AccountHandler
	Currency     *CurrencyHandler
	ExchangeRate *ExchangeRateHandler
	Exchange     *ExchangeHandler
}

// RouterDeps contains dependencies required by the router.
type RouterDeps struct {
	Handlers    Handlers
	SessionStore sessions.Store
	SessionName  string
	RecordAction *ucaction.RecordUseCase
}

// NewRouter creates a new router
func NewRouter(deps RouterDeps) *gin.Engine {
	// create new gin engine
	router := gin.New()
	router.Use(RequestLogger(), Recoverer())
	router.Use(sessions.Sessions(deps.SessionName, deps.SessionStore))

	// create new frontend routes
	frontendRoutes(router)

	// create new version v1 routes
	v1 := versionV1Routes(router)
	v1.Use(ActionJournal(deps.RecordAction))

	// register public auth routes
	authRoutes(v1, deps.Handlers.Auth)

	accountRoutes(v1, deps.Handlers.Account)
	currencyRoutes(v1, deps.Handlers.Currency)
	exchangeRateRoutes(v1, deps.Handlers.ExchangeRate)
	exchangeRoutes(v1, deps.Handlers.Exchange)

	return router
}

// frontendRoutes creates a new frontend routes
func frontendRoutes(router *gin.Engine) {

	isAuthenticated := func(c *gin.Context) bool {
		session := sessions.Default(c)
		return session.Get("name") != nil && session.Get("role") != nil
	}
	router.GET("/", func(c *gin.Context) {
		if isAuthenticated(c) {
			c.Redirect(http.StatusFound, "/dashboard")
			return
		}

		c.Redirect(http.StatusFound, "/login")
	})
	router.GET("/login", func(c *gin.Context) {
		if isAuthenticated(c) {
			c.Redirect(http.StatusFound, "/dashboard")
			return
		}
		rederTemplate(c, http.StatusOK, templ.LoginPage())
	})
	router.GET("/dashboard", func(c *gin.Context) {
		if !isAuthenticated(c) {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		rederTemplate(c, http.StatusOK, templ.DashboardPage())
	})
}

// versionV1Routes creates a new version v1 routes
func versionV1Routes(router *gin.Engine) *gin.RouterGroup {
	return router.Group("/v1")
}

// authRoutes registers the public auth routes
func authRoutes(router *gin.RouterGroup, authHandler *AuthHandler) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
	}
}

// accountRoutes registers the account routes
func accountRoutes(router *gin.RouterGroup, accountHandler *AccountHandler) {
	account := router.Group("/account")
	account.Use(Authenticate())
	{
		account.GET("/me", accountHandler.Account)
		account.GET("/actions", accountHandler.SessionActions)
	}
}

// currencyRoutes registers the currency routes
func currencyRoutes(router *gin.RouterGroup, currencyHandler *CurrencyHandler) {
	currency := router.Group("/currency")
	currency.Use(Authenticate(), Authorization())
	{
		currency.GET("", currencyHandler.GetAllCurrencies)
		currency.POST("", currencyHandler.CreateCurrency)
		currency.GET("/:code", currencyHandler.GetCurrencyByCode)
		currency.PUT("/:code", currencyHandler.UpdateCurrency)
		currency.DELETE("/:code", currencyHandler.DeleteCurrency)
	}
}

// exchangeRateRoutes registers exchange rate routes.
func exchangeRateRoutes(router *gin.RouterGroup, handler *ExchangeRateHandler) {
	group := router.Group("/exchange-rate")
	group.Use(Authenticate())
	{
		group.GET("", handler.ListLive)
		group.GET("/rates", handler.ListRatesByDate)
		group.GET("/report", RequireAdmin(), handler.ListReport)
		group.GET("/:id/history", RequireAdmin(), handler.ListHistory)
	}

	admin := router.Group("/exchange-rate")
	admin.Use(Authenticate(), Authorization())
	{
		admin.POST("", handler.Create)
		admin.PUT("/:id", handler.Update)
	}
}

// exchangeRoutes registers currency exchange routes.
func exchangeRoutes(router *gin.RouterGroup, handler *ExchangeHandler) {
	group := router.Group("/exchange")
	group.Use(Authenticate(), RequireOperator())
	{
		group.POST("", handler.Execute)
	}
}

// rederTemplate renders a template
func rederTemplate(c *gin.Context, status int, component tmpl.Component) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Status(status)

	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		c.Error(err)
	}
}
