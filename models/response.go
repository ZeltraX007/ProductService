package models

type Result struct {
	ResponseCode        string      `json:"response_code" validate:"required"`
	ResponseStatus      string      `json:"response_status"`
	ResponseDescription string      `json:"response_description"`
	ResponseBody        interface{} `json:"response_body"`
}

type PaginatedResponse struct {
	ResponseCode        string                    `json:"response_code" validate:"required"`
	ResponseStatus      string                    `json:"response_status"`
	ResponseDescription string                    `json:"response_description"`
	ResponseBody        PaginationProductResponse `json:"response_body"`
}

type PaginationProductResponse struct {
	PageNo     int         `json:"page_no"`
	PageSize   int         `json:"page_size"`
	TotalCount int         `json:"total_count"`
	TotalPages int         `json:"total_pages"`
	Offset     int         `json:"offset"`
	Products   interface{} `json:"products"`
}
