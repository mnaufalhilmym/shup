package model

import (
	"encoding/json"
	"math"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Response[T any] struct {
	Timestamp  time.Time   `json:"timestamp"`
	RequestID  uuid.UUID   `json:"request_id"`
	Code       int         `json:"code"`
	Status     bool        `json:"status"`
	Message    string      `json:"message"`
	Path       string      `json:"path,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Data       T           `json:"data"`
}

type Pagination struct {
	Page      int `json:"page"`
	Size      int `json:"size"`
	TotalItem int `json:"total_item"`
	TotalPage int `json:"total_page"`
}

func WriteResponse[T any](w http.ResponseWriter, body Response[T]) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(body.Code)
	json.NewEncoder(w).Encode(body)
}

func ResponseOK[T comparable](w http.ResponseWriter, requestID uuid.UUID, data T) {
	WriteResponse(
		w,
		Response[T]{
			Timestamp: time.Now(),
			RequestID: requestID,
			Code:      http.StatusOK,
			Status:    true,
			Message:   "Success",
			Data:      data,
		},
	)
}

func ResponseOKPaginated[T any](w http.ResponseWriter, requestID uuid.UUID, data T, total int, page int, size int) {
	WriteResponse(
		w,
		Response[T]{
			Timestamp: time.Now(),
			RequestID: requestID,
			Code:      http.StatusOK,
			Status:    true,
			Message:   "Success",
			Pagination: &Pagination{
				Page:      page,
				Size:      size,
				TotalItem: total,
				TotalPage: int(math.Ceil(float64(total) / float64(size))),
			},
			Data: data,
		},
	)
}

func ResponseCreated[T any](w http.ResponseWriter, requestID uuid.UUID, data T) {
	WriteResponse(
		w,
		Response[T]{
			Timestamp: time.Now(),
			RequestID: requestID,
			Code:      http.StatusCreated,
			Status:    true,
			Message:   "Success",
			Data:      data,
		},
	)
}

func ResponseError(w http.ResponseWriter, requestID uuid.UUID, err error) {
	appError, ok := err.(*Error)
	if ok {
		WriteResponse(
			w,
			Response[any]{
				Timestamp: time.Now(),
				RequestID: requestID,
				Code:      appError.Code,
				Status:    false,
				Message:   appError.Err.Error(),
				Path:      appError.Path,
			},
		)
		return
	}

	WriteResponse(
		w,
		Response[any]{
			Timestamp: time.Now(),
			RequestID: requestID,
			Code:      http.StatusInternalServerError,
			Status:    false,
			Message:   err.Error(),
		},
	)
}
