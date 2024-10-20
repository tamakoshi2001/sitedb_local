package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/tamakoshi2001/gextension/model"
	"github.com/tamakoshi2001/gextension/service"
)

// A SiteHandler implements handling REST endpoints.
type SiteHandler struct {
	siteService *service.SiteService
}

// NewSiteHandler returns a new SiteHandler.
func NewSiteHandler(siteService *service.SiteService) *SiteHandler {
	return &SiteHandler{
		siteService: siteService,
	}
}

// Create handles the endpoint that creates the Site.
func (h *SiteHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	var req model.CreateSiteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errMsg := fmt.Sprintf("Error: failed to decode request body. Details: %v", err)

		// エラーをログに出力
		log.Println(errMsg)

		// http error
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	res, err := h.siteService.Create(&req)
	if err != nil {
		errMsg := fmt.Sprintf("Error: failed to create Site. Details: %v", err)

		// エラーをログに出力
		log.Println(errMsg)

		// http error
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		errMsg := fmt.Sprintf("Error: failed to encode response body. Details: %v", err)

		// エラーをログに出力
		log.Println(errMsg)

		// http error
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
}

// Read handles the endpoint that read the Site.
func (h *SiteHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	// クエリパラメータからデータを取得
	query := r.URL.Query().Get("query")

	// ReadSiteRequest にキャスト
	req := model.ReadSiteRequest{
		Query: query,
	}

	// クエリパラメータが存在しない場合、エラーレスポンスを返す
	if req.Query == "" {
		http.Error(w, "query parameter is missing", http.StatusBadRequest)
		return
	}

	res, err := h.siteService.Read(&req)
	if err != nil {
		errMsg := fmt.Sprintf("Error: failed to read Site. Details: %v", err)

		// エラーをログに出力
		log.Println(errMsg)

		// http error
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		errMsg := fmt.Sprintf("Error: failed to encode response body. Details: %v", err)

		// エラーをログに出力
		log.Println(errMsg)

		// http error
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
}

// Delete handles the endpoint that delete the Site.
func (h *SiteHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	var req model.DeleteSiteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errMsg := fmt.Sprintf("Error: failed to decode request body. Details: %v", err)

		// エラーをログに出力
		log.Println(errMsg)

		// http error
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	res, err := h.siteService.Delete(&req)
	if err != nil {
		errMsg := fmt.Sprintf("Error: failed to delete Site. Details: %v", err)

		// エラーをログに出力
		log.Println(errMsg)

		// http error
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		errMsg := fmt.Sprintf("Error: failed to encode response body. Details: %v", err)

		// エラーをログに出力
		log.Println(errMsg)

		// http error
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
}
