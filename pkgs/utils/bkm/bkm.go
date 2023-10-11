package bkm

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gogf/gf/encoding/gjson"
)

type BookmarkTree struct {
	Root             *BkmNode
	OriginalFilePath string
	ToSavePath       string
}

func NewBkmTree(btype BrowserType, oPath, toSavePath string) *BookmarkTree {
	t := &BookmarkTree{
		Root:             &BkmNode{BType: btype},
		OriginalFilePath: oPath,
		ToSavePath:       toSavePath,
	}
	t.initTree()
	return t
}

func (that *BookmarkTree) initTree() {
	if that.Root.BType == Chrome {
		b, _ := os.ReadFile(that.OriginalFilePath)
		j := gjson.New(b)
		that.Root.ParseTree(j.GetString("roots.bookmark_bar"))
	} else if that.Root.BType == Firefox {
		if db, err := sql.Open("sqlite", that.OriginalFilePath); err == nil {
			sql_ := `SELECT id FROM moz_bookmarks WHERE title="toolbar"`
			rows, err := db.Query(sql_)
			if err != nil {
				fmt.Println(err)
				return
			}
			var id int64
			for rows.Next() {
				if err = rows.Scan(&id); err != nil {
					fmt.Println(err)
					continue
				}
			}
			rows.Close()
			that.Root.ParseFirefoxBkm(id, db)
		} else {
			fmt.Println("[Open firefox bookmark database failed] ", err)
		}
	}
}

func (that *BookmarkTree) SaveHtml() {
	if that.ToSavePath == "" {
		that.ToSavePath = "bookmark.html"
	}
	htmStr := fmt.Sprintf("%s\n%s", Header, that.Root.Html())
	os.WriteFile(that.ToSavePath, []byte(htmStr), 0666)
}
