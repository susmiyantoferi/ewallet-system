package dto

import "math"

type WebResponse struct {
	Success  bool          `json:"success"`
	Message  string        `json:"message"`
	MetaData *PageMetaData `json:"meta_data,omitempty"`
	Data     any           `json:"data"`
}

func (w *WebResponse) WithMetadata(metadata *PageMetaData) *WebResponse {
	w.MetaData = metadata
	return w
}

func SuccessResponse(data any, message ...string) *WebResponse {
	resMessage := "Success"
	if len(message) > 0 {
		resMessage = message[0]
	}
	return &WebResponse{
		Success: true,
		Message: resMessage,
		Data:    data,
	}
}

func ErrorResponse(message ...string) *WebResponse {
	resMessage := "Error"
	if len(message) > 0 {
		resMessage = message[0]
	}
	return &WebResponse{
		Success: false,
		Message: resMessage,
		Data:    nil,
	}
}

type PageMetaData struct {
	Page      int   `json:"page"`
	Limit     int   `json:"per_page"`
	TotalItem int64 `json:"total_item"`
	TotalPage int64 `json:"total_page"`
}

func (p *PageMetaData) WithCountPage() *PageMetaData {
	p.TotalPage = int64(math.Ceil(float64(p.TotalItem) / float64(p.Limit)))
	return p
}
