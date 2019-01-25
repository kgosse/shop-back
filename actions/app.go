package actions

import (
	"log"

	"github.com/casbin/casbin"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	forcessl "github.com/gobuffalo/mw-forcessl"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/unrolled/secure"

	rbac "github.com/kgosse/buffalo-mw-rbac"

	"github.com/gobuffalo/buffalo-pop/pop/popmw"
	contenttype "github.com/gobuffalo/mw-contenttype"
	tokenauth "github.com/gobuffalo/mw-tokenauth"
	"github.com/gobuffalo/x/sessions"
	"github.com/kgosse/shop-back/models"
	"github.com/rs/cors"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionStore: sessions.Null{},
			PreWares: []buffalo.PreWare{
				cors.Default().Handler,
			},
			SessionName: "_shop_session",
		})

		// setup casbin auth rules
		authEnforcer, err := casbin.NewEnforcerSafe(envy.Get("RBAC_AUT_MODEL_PATH", "auth_model.conf"), envy.Get("RBAC_POLICY_PATH", "policy.csv"))
		if err != nil {
			log.Fatal(err)
		}

		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Set the request content type to JSON
		app.Use(contenttype.Set("application/json"))

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.Connection)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))

		tauth := tokenauth.New(tokenauth.Options{})
		// app.Use(tauth)

		roleFunc := func(c buffalo.Context) (string, error) {
			if c.Value("claims") == nil {
				return "anonymous", nil
			}
			claims := c.Value("claims").(jwt.MapClaims)
			r := claims["role"].(string)
			if r != "" {
				return r, nil
			}
			return "anonymous", nil
		}

		mwCashbin := rbac.New(authEnforcer, roleFunc)

		app.GET("/", HomeHandler)

		api := app.Group("/v1/")
		api.Use(tauth)
		api.Use(mwCashbin)

		adminAPI := app.Group("/v1/admin")
		adminAPI.Use(tauth)
		adminAPI.Use(mwCashbin)

		pr := ProductsResource{}
		ur := UsersResource{}

		api.POST("auth/login", ur.Login)
		api.Middleware.Skip(tauth, ur.Login)

		adminAPI.POST("auth/login", ur.LoginAdmin)
		adminAPI.Middleware.Skip(tauth, ur.LoginAdmin)

		api.Resource("/wishlists", WishlistsResource{})

		wur := api.Resource("/users", ur)
		wur.Middleware.Skip(tauth, ur.Create, ur.List)

		// wr := api.Resource("/products", pr)
		api.Resource("/products", pr)
		// wr.Middleware.Skip(tauth, pr.List)
		app.Resource("/roles", RolesResource{})
	}

	return app
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}
