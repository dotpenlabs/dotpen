package main

import (
	"embed"
	"log"
	"net/mail"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/jsvm"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"github.com/pocketbase/pocketbase/tools/types"

	"dotpen.co/server/hooks/modules"
)

//go:embed emails
var emails embed.FS

func main() {
	os.Args = append(os.Args, "--dir", "data")
	app := pocketbase.New()

	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Dir:          "./migrations",
		TemplateLang: "go",
		Automigrate:  isGoRun,
	})

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/api/crawl", func(e *core.RequestEvent) error {
			modules.UseCrawl(e.Response, e.Request, app)
			return nil
		}).Bind(apis.RequireAuth())

		se.Router.GET("/api/proxy", func(e *core.RequestEvent) error {
			modules.UseProxy(e.Response, e.Request)
			return nil
		}).Bind(apis.RequireAuth())

		se.Router.GET("/api/waitlist-members", func(e *core.RequestEvent) error {
			var result struct {
				Count int `db:"c"`
			}

			err := e.App.DB().
				NewQuery("SELECT COUNT(*) AS c FROM waitlist").
				One(&result)
			if err != nil {
				return err
			}

			e.Response.WriteHeader(200)
			e.Response.Write([]byte(strconv.Itoa(result.Count)))

			return nil
		})

		se.Router.GET("/{path...}", apis.Static(os.DirFS("./public"), false))
		se.Router.GET("/_/images/{path...}", func(e *core.RequestEvent) error {
			staticFS := os.DirFS("./public/assets")
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

	app.OnRecordUpdate("waitlist").BindFunc(func(e *core.RecordEvent) error {
		if e.Record.GetBool("accept") {
			emailBytes, err := emails.ReadFile("emails/invited.html")
			if err != nil {
				return err
			}

			message := &mailer.Message{
				From: mail.Address{
					Address: e.App.Settings().Meta.SenderAddress,
					Name:    e.App.Settings().Meta.SenderName,
				},
				To:      []mail.Address{{Address: e.Record.GetString("email")}},
				Subject: "🎉 You've been invited to join Dotpen!",
				HTML:    string(emailBytes),
			}

			e.App.NewMailClient().Send(message)
		}
		return e.Next()
	})

	app.OnRecordUpdateExecute("bookmarks").BindFunc(func(e *core.RecordEvent) error {
		app.Logger().Debug("RecordUpdate: bookmarks", "action", "update")

		if e.Record.GetBool("deleted") {
			app.Logger().Info("RecordUpdate: bookmarks", "action", "delete")

			e.Record.Set("label", "")
			e.Record.Set("link", "")
			e.Record.Set("favicon", "")
			e.Record.Set("cover", "")

			app.Save(e.Record)
		}
		return e.Next()
	})

	app.Cron().MustAdd("Remove deleted bookmarks", "0 0 * * *", func() {
		app.Logger().Debug("Cron: Remove deleted bookmarks")

		records, err := app.FindAllRecords("bookmarks",
			dbx.NewExp("deleted = true"),
		)

		if err != nil {
			app.Logger().Error("Cron: Remove deleted bookmarks", "error", err.Error())
			return
		}

		app.Logger().Debug("Cron: Found deleted bookmarks", "count", len(records))

		cutoff, err := types.ParseDateTime(time.Now().Add(24 * time.Hour))
		if err != nil {
			app.Logger().Error("Cron: Remove deleted bookmarks", "error", err.Error())
			return
		}

		for _, r := range records {
			if r.GetDateTime("updated").Before(cutoff) || r.GetDateTime("created").Before(cutoff) {
				app.Logger().Debug("Cron: Deleting bookmark", "id", r.Id)
				err := app.Delete(r)
				if err != nil {
					app.Logger().Error("Cron: Deleting bookmark", "id", r.Id, "error", err.Error())
				}
			}
		}
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
