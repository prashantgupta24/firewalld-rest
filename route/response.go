package route

// JSONResponse defines meta and interface struct
type JSONResponse struct {
	// Reserved field to add some meta information to the API response
	Meta interface{} `json:"meta"`
	Data interface{} `json:"data"`
}

// JSONErrorResponse defines error struct
type JSONErrorResponse struct {
	Error *APIError `json:"error"`
}

// APIError defines api error struct
type APIError struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
}
