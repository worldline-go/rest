package rest

// Response of Rest-API.
type Response[T any] struct {
	Message *Message `json:"message,omitempty"`
	Meta    *Meta    `json:"meta,omitempty"`
	Payload T        `json:"payload"`
}

// ResponseMessage is only used for message responses like errors or notify.
type ResponseMessage struct {
	Message *Message `json:"message,omitempty"`
}

// Message object returns a response text key for i18n files which is used for frontend text rendering.
// A params object returns dynamic parameters in the same format as the text string is returned.
// They are also used on the frontend side. It can cover errors, success messages and regular notifications.
type Message struct {
	Text   string         `json:"text,omitempty"`
	Params map[string]any `json:"params,omitempty"`
	Err    string         `json:"error,omitempty"`
}

// Meta is a structure that contains pagination metadata for the ListResponseG.
type Meta struct {
	// TotalItemCount is the total number of entities that match the query.
	TotalItemCount uint64 `json:"total_item_count,omitempty"`
	// Limit is the limit used within the request.
	// If not defined in the query parameters, this should be the default value used in the service endpoint.
	Limit int64 `json:"limit,omitempty"`
	// Offset is the offset used within the request.
	Offset int64 `json:"offset,omitempty"`
}
