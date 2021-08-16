package silence

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type App struct {
	PageSize int
	//RefreshInterval time.Duration
	indexTemplate *template.Template
	itemCache     []NewsItem
	api           Client
}

func NewApp() *App {
	return &App{
		PageSize:      30,
		indexTemplate: generateTemplateFile(),
		itemCache:     nil,
		api:           Client{},
	}
}

type Page struct {
	TopStories []NewsItem
	Page       int
	NextPage   int
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var page = 1
	var err error

	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			log.Println("error encountered")
		}
	}

	buf := bytes.NewBuffer([]byte{})
	data, err := a.getTopStores(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	nextPage := page + 1

	p := Page{
		TopStories: data,
		Page:       page,
		NextPage:   nextPage,
	}

	fmt.Printf("\nPage:  %d, Next Page: %d", page, p.NextPage)

	if err := a.indexTemplate.Execute(buf, p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseBytes, err := ioutil.ReadAll(buf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseBytes)

}
func (a *App) getTopStores(page int) ([]NewsItem, error) {
	ids, err := a.api.RetrieveTopStories()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve ids: %s", err)
	}

	// ids = []long list
	var newsItems []NewsItem
	idsSlice := ids[(page-1)*a.PageSize : page*a.PageSize+1]

	//FirstPage:    1,
	//LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),

	//func (f Filters) offset() int {
	//	return (f.Page - 1) * f.PageSize
	//}

	for _, newsItemID := range idsSlice {

		item, err := a.api.RetrieveNewsItemById(newsItemID)
		if err != nil {
			return nil, err
		}

		if item.URL == "" || item.Title == "" {
			continue
		}

		newsItems = append(newsItems, item)

		if len(newsItems) > a.PageSize {
			break
		}
	}

	return newsItems, nil

}

//go:embed index.tmpl
var tpls embed.FS

func generateTemplateFile() *template.Template {
	t, err := template.ParseFS(tpls, "index.tmpl")
	if err != nil {
		log.Println(err)
	}

	return t

}
