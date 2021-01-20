package controller

import (
	"encoding/json"
	"fmt"
	"github.com/MathisBurger/OpenInventory-Backend/models"
	OwnSQL "github.com/MathisBurger/OpenInventory-Backend/mysql"
	"github.com/gofiber/fiber/v2"
)

func RemoveTableEntryController(c *fiber.Ctx) error {
	raw := string(c.Body())
	obj := models.RemoveTableEntryRequestModel{}
	err := json.Unmarshal([]byte(raw), &obj)
	if err != nil {
		response, err := models.GetJsonResponse("Invaild JSON body", "alert alert-danger", "error", "None", 200)
		if err != nil {
			panic(err)
		}
		return c.Send(response)
	}
	if OwnSQL.MySQL_loginWithToken(obj.Username, obj.Password, obj.Token) {
		fmt.Println(obj)
		conn := OwnSQL.GetConn()
		stmt, _ := conn.Prepare("DELETE FROM `table_" + obj.TableName + "` WHERE `id`=?")
		aff, _ := stmt.Exec(obj.RowID)
		aff_res, _ := aff.RowsAffected()
		if aff_res == 0 {
			res, _ := models.GetJsonResponse("EntryID not found", "alert alert-warning", "ok", "None", 200)
			return c.Send(res)
		}
		stmt, _ = conn.Prepare("SELECT `entries` FROM `inv_tables` WHERE `name`=?")
		resp, err := stmt.Query(obj.TableName)
		if err != nil {
			panic(err.Error())
		}
		entries := 0
		for resp.Next() {
			var entry OwnSQL.Entries
			err = resp.Scan(&entry.Entries)
			if err != nil {
				panic(err.Error())
			}
			entries = entry.Entries
		}
		entries -= 1
		stmt, _ = conn.Prepare("UPDATE `inv_tables` SET `entries`=? WHERE `name`=?;")
		stmt.Exec(entries, obj.TableName)
		resp.Close()
		stmt.Close()
		conn.Close()
		res, _ := models.GetJsonResponse("Successfully deleted entry", "alert alert-success", "ok", "None", 200)
		return c.Send(res)
	} else {
		resp, _ := models.GetJsonResponse("You do not have the permission perform this", "alert alert-danger", "ok", "None", 200)
		return c.Send(resp)
	}
}
