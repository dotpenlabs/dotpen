package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_2111211931")
		if err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(1, []byte(`{
			"hidden": false,
			"id": "file1704208859",
			"maxSelect": 1,
			"maxSize": 0,
			"mimeTypes": [
				"image/webp",
				"image/jpeg"
			],
			"name": "cover",
			"presentable": false,
			"protected": false,
			"required": true,
			"system": false,
			"thumbs": [
				"50x50",
				"100x100",
				"500x500"
			],
			"type": "file"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_2111211931")
		if err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(1, []byte(`{
			"hidden": false,
			"id": "file1704208859",
			"maxSelect": 1,
			"maxSize": 0,
			"mimeTypes": [
				"image/webp",
				"image/jpeg"
			],
			"name": "icon",
			"presentable": false,
			"protected": false,
			"required": true,
			"system": false,
			"thumbs": [
				"50x50",
				"100x100",
				"500x500"
			],
			"type": "file"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
