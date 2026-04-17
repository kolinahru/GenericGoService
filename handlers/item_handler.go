package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"go-day3/service"
)

type ItemHandler struct {
	service service.ItemService
}

func NewItemHandler(service service.ItemService) *ItemHandler {
	return &ItemHandler{service: service}
}

type createItemRequest struct {
	Name string `json:"name"`
}

type updateItemRequest struct {
	Name string `json:"name"`
}

func (h *ItemHandler) HandleItems(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getItems(w, r)
	case http.MethodPost:
		h.createItem(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ItemHandler) HandleItemByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getItemByID(w, r, id)
	case http.MethodPut:
		h.updateItem(w, r, id)
	case http.MethodDelete:
		h.deleteItem(w, r, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ItemHandler) getItems(w http.ResponseWriter, _ *http.Request) {
	items, err := h.service.GetItems()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, items)
}

func (h *ItemHandler) getItemByID(w http.ResponseWriter, _ *http.Request, id int) {
	item, err := h.service.GetItemByID(id)
	if err != nil {
		if service.IsNotFoundError(err) {
			http.Error(w, "item not found", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, item)
}

func (h *ItemHandler) createItem(w http.ResponseWriter, r *http.Request) {
	var req createItemRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item, err := h.service.CreateItem(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusCreated, item)
}

func (h *ItemHandler) updateItem(w http.ResponseWriter, r *http.Request, id int) {
	var req updateItemRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item, err := h.service.UpdateItem(id, req.Name)
	if err != nil {
		if service.IsNotFoundError(err) {
			http.Error(w, "item not found", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, item)
}

func (h *ItemHandler) deleteItem(w http.ResponseWriter, _ *http.Request, id int) {
	err := h.service.DeleteItem(id)
	if err != nil {
		if service.IsNotFoundError(err) {
			http.Error(w, "item not found", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseIDFromPath(path string) (int, error) {
	idPart := strings.TrimPrefix(path, "/items/")
	return strconv.Atoi(idPart)
}
