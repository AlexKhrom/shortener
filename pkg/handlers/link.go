package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"ozon/pkg/items"
)

type LinkHandler struct {
	Repo items.LinkRepoInterface
}

func NewLinkHandler(db *sql.DB, storageType string) *LinkHandler {
	hand := new(LinkHandler)
	hand.Repo = items.NewLinkRepo(db, storageType)
	return hand
}

func (h *LinkHandler) NewLink(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hi new link")
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println("err with body")
		items.JSONError(w, 500, "some bad in read data")
		return
	}
	req := &items.Link{}

	err1 := json.Unmarshal(body, req)
	if err1 != nil {
		fmt.Println("err unmarshal body")
		items.JSONError(w, 500, "some bad in unmarshal data")
		return
	}
	newLink := h.Repo.NewLink(req.URL)

	response, err := json.Marshal(items.Link{URL: newLink})
	if err != nil {
		fmt.Println("err json marshal")
		items.JSONError(w, 500, "some bad in marshal response")
		return
	}

	w.Write(response)
}

func (h *LinkHandler) GetLink(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hi get link")

	link := r.URL.Query().Get("url")

	if link == "" {
		fmt.Println("err bad link")
		items.JSONError(w, 500, "bad link")
		return
	}

	origLink := h.Repo.GetOrigLink(link)

	response, err := json.Marshal(items.Link{URL: origLink})
	if err != nil {
		fmt.Println("err json marshal")
		items.JSONError(w, 500, "some bad in marshal response")
		return
	}
	w.Write(response)
}
