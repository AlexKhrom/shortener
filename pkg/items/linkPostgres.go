package items

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Link struct {
	URL string
}

type LinkRepoInterface interface {
	NewLink(origLink string) string
	GetOrigLink(unifiedLink string) string
}

type LinkRepoPostgres struct {
	db  *sql.DB
	Mut sync.Mutex
}

func NewLinkRepo(db *sql.DB, storageType string) LinkRepoInterface {
	switch storageType {
	case "inMemory":
		repo := new(LinkRepoMemory)
		repo.origLinks = make(map[string]string)
		repo.newLinks = make(map[string]string)
		return repo
	case "postgres":
		repo := new(LinkRepoPostgres)
		repo.db = db
		return repo
	}
	return nil
}

func (r *LinkRepoPostgres) GetUnifiedLink(origLink string) string {
	rows, err := r.db.Query(`SELECT "unified_link" FROM "links" WHERE original_link=$1`,
		origLink,
	)
	if err != nil {
		fmt.Println("err creat get unified link sql1 = ", err)
		return ""
	}
	defer rows.Close()
	for rows.Next() {
		var link Link
		err = rows.Scan(
			&link.URL,
		)
		if err != nil {
			fmt.Println("err  creat unified link sql scan = ", err)
			return ""
		}
		return link.URL
	}
	return ""
}
func (r *LinkRepoPostgres) GetOriginalLink(unified string) bool {
	rows, err := r.db.Query(`SELECT "original_link" FROM "links" WHERE unified_link=$1`,
		unified,
	)
	if err != nil {
		fmt.Println("err creat new link sql2 = ", err)
		return false
	}
	defer rows.Close()
	for rows.Next() {
		return false
	}
	return true
}

func (r *LinkRepoPostgres) NewLink(origLink string) string {
	unifiedLink := r.GetUnifiedLink(origLink)
	if unifiedLink != "" {
		return unifiedLink
	}
	key := origLink
	newLink := ""
	for {
		newLink = CreatNewLinkString(key)
		ok := r.GetOriginalLink(newLink)
		if ok {
			r.db.QueryRow(
				`INSERT INTO links("original_link","unified_link") VALUES ($1,$2)`,
				origLink,
				newLink,
			)
			return newLink
		}

	}

	return ""
}

func (r *LinkRepoPostgres) GetOrigLink(unifiedLink string) string {
	rows, err := r.db.Query(`SELECT "original_link" FROM "links" WHERE unified_link=$1`,
		unifiedLink,
	)
	if err != nil {
		fmt.Println("err get orders sql = ", err)
		return ""
	}
	defer rows.Close()
	for rows.Next() {
		var link Link
		err = rows.Scan(
			&link.URL,
		)
		if err != nil {
			fmt.Println("err get orders sql scan = ", err)
			return ""
		}
		return link.URL

	}
	return ""
}

func CreatNewLinkString(origLink string) string {
	h := sha1.New()
	h.Write([]byte(origLink))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	return sha1_hash[:10]
}

func JSONError(w http.ResponseWriter, status int, msg string) {
	resp, err := json.Marshal(map[string]interface{}{
		"error": msg,
	})
	w.WriteHeader(status)
	if err != nil {
		fmt.Println("error in JSONError ")
		return
	}
	_, err2 := w.Write(resp)
	if err2 != nil {
		fmt.Println("some bad in JSONError write response")
	}
}
