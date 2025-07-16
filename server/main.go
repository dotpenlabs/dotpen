package main

import (
	"log"
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/jsvm"

	"dotpen.co/server/pkg/modules"
)

func main() {
	app := pocketbase.New()

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// se.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), false))
		se.Router.GET("/api/crawl", func(e *core.RequestEvent) error {
			modules.UseCrawl(e.Response, e.Request, app)
			return nil
		}).Bind(apis.RequireAuth())
		
		se.Router.GET("/api/proxy", func(e *core.RequestEvent) error {
			modules.UseProxy(e.Response, e.Request)
			return nil
		}).Bind(apis.RequireAuth())

		se.Router.GET("/_/images/{path...}", func(e *core.RequestEvent) error {
			staticFS := os.DirFS("./pb_public/assets")
			return apis.Static(staticFS, false)(e)
		})

		jsvm.MustRegister(app, jsvm.Config{
			HooksWatch:    true,
			HooksPoolSize: 15,
		})

		return se.Next()
	})

	app.OnRecordCreateRequest("users").BindFunc(func(e *core.RecordRequestEvent) error {
		if e.HasSuperuserAuth() {
			return e.Next()
		}

		wli, err := app.FindFirstRecordByFilter("waitlist", "email = '"+e.Record.Email()+"'")
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				return apis.NewApiError(400, "user is not accepted on the waitlist", "WaitlistNoAccept")
			}
			return err
		}

		if !wli.GetBool("accept") {
			return apis.NewApiError(400, "user is not accepted on the waitlist", "WaitlistNoAccept")
		}

		return e.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
