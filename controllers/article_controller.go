package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nemdull/myapi/apperrors"
	"github.com/nemdull/myapi/controllers/services"
	"github.com/nemdull/myapi/models"
)

type ArticleController struct {
	// ここに必要なフィールドを追加
	service services.ArticleServicer
}

func NewArticleController(s services.ArticleServicer) *ArticleController {
	return &ArticleController{
		service: s,
	}
}

func (c *ArticleController) HelloHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello, world!\n")
}

func (c *ArticleController) PostArticleHandler(w http.ResponseWriter, req *http.Request) {
	var reqArticle models.Article
	if err := json.NewDecoder(req.Body).Decode(&reqArticle); err != nil {
		err = apperrors.ReqBodyDecodeFailed.Wrap(err, "failed to decode request body")
		apperrors.ErrorHandler(w, req, err)
		return
	}

	article, err := c.service.PostArticleService(reqArticle)
	if err != nil {
		apperrors.ErrorHandler(w, req, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(article)
}
func (c *ArticleController) ArticleListHandler(w http.ResponseWriter, req *http.Request) {
	queryMap := req.URL.Query()

	// クエリパラメータpageを取得
	var page int
	if p, ok := queryMap["page"]; ok && len(p) > 0 {
		var err error
		page, err = strconv.Atoi(p[0])
		if err != nil {
			err = apperrors.BadParams.Wrap(err, "invalid page parameter")
			http.Error(w, "Invalid query parameter", http.StatusBadRequest)
			return
		}
	}

	articles, err := c.service.GetArticleListService(page)
	if err != nil {
		http.Error(w, "fail to get article list\n", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}
func (c *ArticleController) ArticleDetailHandler(w http.ResponseWriter, req *http.Request) {
	articleID, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}
	article, err := c.service.GetArticleService(articleID)
	if err != nil {
		http.Error(w, "fail to get article\n", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}
func (c *ArticleController) PostNiceHandler(w http.ResponseWriter, req *http.Request) {
	var reqNice models.Article
	if err := json.NewDecoder(req.Body).Decode(&reqNice); err != nil {
		http.Error(w, "fail to decode json\n", http.StatusBadRequest)
		return
	}

	nice, err := c.service.PostNiceService(reqNice)
	if err != nil {
		http.Error(w, "fail to post nice\n", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nice)
}
