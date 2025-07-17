package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/jsvm"
	"github.com/pocketbase/pocketbase/tools/types"

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

	app.OnRecordUpdate("bookmarks").BindFunc(func(e *core.RecordEvent) error {
		app.Logger().Debug("RecordUpdate: bookmarks", "action", "update")

		if e.Record.GetBool("deleted") {
			app.Logger().Info("RecordUpdate: bookmarks", "action", "delete")

			e.Record.Set("label", "")
			e.Record.Set("url", "")
			e.Record.Set("favicon", "")
			e.Record.Set("cover", "")
		}
		return e.Next()
	})

	app.Cron().MustAdd("Remove deleted bookmarks", "0 0 * * *", func() {
		app.Logger().Debug("Cron: Remove deleted bookmarks")

		records, err := app.FindAllRecords("bookmarks",
			dbx.NewExp("deleted = 'true'"),
		)

		if err != nil {
			app.Logger().Error("Cron: Remove deleted bookmarks", "error", err.Error())
			return
		}

		cutoff, _ := types.ParseDateTime(time.Now().Add(-31 * 24 * time.Hour))

		for _, r := range records {
			if r.GetDateTime("updated").Before(cutoff) {
				app.Delete(r)
			}
		}
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
