package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		jsonData := `{
			"createRule": "@request.auth.id = user.id",
			"deleteRule": "@request.auth.id = user.id",
			"fields": [
				{
					"autogeneratePattern": "[a-z0-9]{15}",
					"hidden": false,
					"id": "text3208210256",
					"max": 15,
					"min": 15,
					"name": "id",
					"pattern": "^[a-z0-9]+$",
					"presentable": false,
					"primaryKey": true,
					"required": true,
					"system": true,
					"type": "text"
				},
				{
					"autogeneratePattern": "[a-z0-9]{8}",
					"hidden": false,
					"id": "text2560465762",
					"max": 0,
					"min": 0,
					"name": "slug",
					"pattern": "^[a-z0-9]+$",
					"presentable": false,
					"primaryKey": false,
					"required": true,
					"system": false,
					"type": "text"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "text1579384326",
					"max": 0,
					"min": 0,
					"name": "name",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": false,
					"system": false,
					"type": "text"
				},
				{
					"cascadeDelete": true,
					"collectionId": "_pb_users_auth_",
					"hidden": false,
					"id": "relation2375276105",
					"maxSelect": 1,
					"minSelect": 0,
					"name": "user",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "relation"
				},
				{
					"cascadeDelete": false,
					"collectionId": "pbc_601157786",
					"hidden": false,
					"id": "relation4232930610",
					"maxSelect": 1,
					"minSelect": 0,
					"name": "collection",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "relation"
				},
				{
					"hidden": false,
					"id": "json2918445923",
					"maxSize": 0,
					"name": "data",
					"presentable": false,
					"required": true,
					"system": false,
					"type": "json"
				},
				{
					"hidden": false,
					"id": "autodate2990389176",
					"name": "created",
					"onCreate": true,
					"onUpdate": false,
					"presentable": false,
					"system": false,
					"type": "autodate"
				},
				{
					"hidden": false,
					"id": "autodate3332085495",
					"name": "updated",
					"onCreate": true,
					"onUpdate": true,
					"presentable": false,
					"system": false,
					"type": "autodate"
				}
			],
			"id": "pbc_700096677",
			"indexes": [],
			"listRule": null,
			"name": "snapshots",
			"system": false,
			"type": "base",
			"updateRule": "@request.auth.id = user.id",
			"viewRule": ""
		}`

		collection := &core.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_700096677")
		if err != nil {
			return err
		}

		return app.Delete(collection)
	})
}
