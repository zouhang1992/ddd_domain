package facade

import (
	"encoding/json"
	"github.com/zouhang1992/ddd_domain/internal/application/command"
	"github.com/zouhang1992/ddd_domain/internal/application/query"
	buscommand "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/command"
	busquery "github.com/zouhang1992/ddd_domain/internal/infrastructure/bus/query"
	"net/http"
)

// CQRSRoomHandler 基于 CQRS 的房间 HTTP 处理器
type CQRSRoomHandler struct {
	commandBus *buscommand.Bus
	queryBus   *busquery.Bus
}

// NewCQRSRoomHandler 创建基于 CQRS 的房间处理器
func NewCQRSRoomHandler(commandBus *buscommand.Bus, queryBus *busquery.Bus) *CQRSRoomHandler {
	return &CQRSRoomHandler{
		commandBus: commandBus,
		queryBus:   queryBus,
	}
}

// RegisterRoutes 注册路由
func (h *CQRSRoomHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /rooms", h.Create)
	mux.HandleFunc("GET /rooms", h.List)
	mux.HandleFunc("GET /rooms/{id}", h.Get)
	mux.HandleFunc("PUT /rooms/{id}", h.Update)
	mux.HandleFunc("DELETE /rooms/{id}", h.Delete)
}

// Create 创建房间
func (h *CQRSRoomHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		LocationID string   `json:"location_id"`
		RoomNumber string   `json:"room_number"`
		Tags       []string `json:"tags"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := command.CreateRoomCommand{
		LocationID: req.LocationID,
		RoomNumber: req.RoomNumber,
		Tags:       req.Tags,
	}

	result, err := h.commandBus.Dispatch(cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	room := result
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(room)
}

// List 列出房间
func (h *CQRSRoomHandler) List(w http.ResponseWriter, r *http.Request) {
	var q query.ListRoomsQuery

	// 获取查询参数
	locationID := r.URL.Query().Get("location_id")
	if locationID != "" {
		q.LocationID = locationID
	}

	result, err := h.queryBus.Dispatch(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	queryResult := result.(*query.RoomsQueryResult)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(queryResult.Items)
}

// Get 获取房间
func (h *CQRSRoomHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	q := query.GetRoomQuery{ID: id}

	result, err := h.queryBus.Dispatch(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	queryResult := result.(*query.RoomQueryResult)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(queryResult.Room)
}

// Update 更新房间
func (h *CQRSRoomHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct {
		LocationID string   `json:"location_id"`
		RoomNumber string   `json:"room_number"`
		Tags       []string `json:"tags"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := command.UpdateRoomCommand{
		ID:         id,
		LocationID: req.LocationID,
		RoomNumber: req.RoomNumber,
		Tags:       req.Tags,
	}

	result, err := h.commandBus.Dispatch(cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	room := result
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(room)
}

// Delete 删除房间
func (h *CQRSRoomHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cmd := command.DeleteRoomCommand{ID: id}

	if _, err := h.commandBus.Dispatch(cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
