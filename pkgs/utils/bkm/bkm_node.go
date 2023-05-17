package bkm

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gogf/gf/encoding/gjson"
)

var toolbarFlag bool

type BrowserType string

const (
	Chrome  BrowserType = "chrome"
	Firefox BrowserType = "firefox"
	BFolder string      = "folder"
	BUrl    string      = "url"
)

type BkmNode struct {
	Children     []*BkmNode
	DateAdded    string
	DateLastUsed string
	DateModified string
	Guid         string
	Id           string
	Name         string
	Type         string
	Url          string
	BType        BrowserType
}

func (that *BkmNode) ParseTree(content interface{}) {
	if that.BType == Chrome {
		j := gjson.New(content)
		that.DateAdded = j.GetString("date_added")
		that.DateLastUsed = j.GetString("date_last_used")
		that.DateModified = j.GetString("date_modified")
		that.Guid = j.GetString("guid")
		that.Id = j.GetString("id")
		that.Name = j.GetString("name")
		that.Type = j.GetString("type")
		that.Url = j.GetString("url")
		if that.Type == BFolder {
			children := j.GetArray("children")
			for _, v := range children {
				child := &BkmNode{BType: Chrome}
				that.Children = append(that.Children, child)
				child.ParseTree(v)
			}
		}
	}
}

func parseFirfoxBkmType(a int64) string {
	switch a {
	case 2:
		return BFolder
	default:
		return BUrl
	}
}

const (
	FirefoxSQLChildren string = `SELECT b.id, b.type, b.title, b.dateAdded, b.lastModified, b.guid, p.url 
From moz_bookmarks AS b LEFT JOIN moz_places AS p ON b.fk=p.id WHERE b.parent=%d`
	FirefoxSQLParent string = `SELECT b.id, b.type, b.title, b.dateAdded, b.lastModified, b.guid, p.url 
	From moz_bookmarks AS b LEFT JOIN moz_places AS p ON b.fk=p.id WHERE b.id=%d`
	CloseJournalMode string = `PRAGMA journal_mode=off`
)

func (that *BkmNode) ParseFirefoxBkm(parentId int64, db *sql.DB) {
	if that.BType != Firefox {
		return
	}
	psql := fmt.Sprintf(FirefoxSQLParent, parentId)
	rows, err := db.Query(psql)
	if err != nil {
		fmt.Println("[Query Failed]", err)
		return
	}
	for rows.Next() {
		var (
			id, bType, dateAdded, lastModified int64
			title, guid                        string
			bUrl                               any
		)
		if err = rows.Scan(&id, &bType, &title, &dateAdded, &lastModified, &guid, &bUrl); err != nil {
			fmt.Println(err)
			continue
		}
		that.Id = fmt.Sprintf("%d", id)
		that.Type = parseFirfoxBkmType(bType)
		that.Name = title
		that.DateAdded = fmt.Sprintf("%d", dateAdded)
		that.DateModified = fmt.Sprintf("%d", lastModified)
		that.Guid = guid
		if _url, ok := bUrl.(string); ok {
			that.Url = _url
		}
		break
	}
	rows.Close()

	sql_ := fmt.Sprintf(FirefoxSQLChildren, parentId)
	rows, err = db.Query(sql_)
	if err != nil {
		fmt.Println("[Query Failed]", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id, bType, dateAdded, lastModified int64
			title, guid                        string
			bUrl                               any
		)
		if err = rows.Scan(&id, &bType, &title, &dateAdded, &lastModified, &guid, &bUrl); err != nil {
			fmt.Println(err)
			continue
		}

		if that.Type == BFolder {
			child := &BkmNode{BType: Firefox}
			child.ParseFirefoxBkm(id, db)
			that.Children = append(that.Children, child)
		} else {
			child := &BkmNode{BType: Firefox}
			child.Id = fmt.Sprintf("%d", id)
			child.Type = parseFirfoxBkmType(bType)
			child.Name = title
			child.DateAdded = fmt.Sprintf("%d", dateAdded)
			child.DateModified = fmt.Sprintf("%d", lastModified)
			child.Guid = guid
			if _url, ok := bUrl.(string); ok {
				child.Url = _url
			}
			that.Children = append(that.Children, child)
		}
	}
}

func (that *BkmNode) Html() string {
	if len(that.DateAdded) > 10 {
		that.DateAdded = that.DateAdded[:10]
	}
	if len(that.DateModified) > 10 {
		that.DateModified = that.DateModified[:10]
	}
	if len(that.DateLastUsed) > 10 {
		that.DateLastUsed = that.DateLastUsed[:10]
	}
	if that.Type == BFolder {
		var h3Str string
		if (that.Name == "toolbar" || that.Name == "收藏夹栏" || that.Name == "收藏夹") && !toolbarFlag {
			h3Str = fmt.Sprintf(H3,
				that.DateAdded,
				that.DateModified,
				TOOLBAR,
				that.Name)
			toolbarFlag = true
		} else {
			h3Str = fmt.Sprintf(H3,
				that.DateAdded,
				that.DateModified,
				"",
				that.Name)
		}

		dlStrList := []string{}
		for _, node := range that.Children {
			dlStrList = append(dlStrList, node.Html())
		}
		dlStr := strings.Join(dlStrList, "\n")
		return fmt.Sprintf(Folder, h3Str, dlStr)
	} else if that.Type == BUrl {
		aStr := fmt.Sprintf(A,
			that.Url,
			that.DateAdded,
			that.Name)
		dtStr := fmt.Sprintf(Dt, aStr)
		return dtStr
	}
	return ""
}

func NewRoot(btype BrowserType) *BkmNode {
	return &BkmNode{BType: btype}
}
