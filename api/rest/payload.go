package rest

// Request Payload used for creation of exection plan and bash script content
type RequestPayload struct {
	Tasks []Task `json:"tasks"`
}

// This structure is used inside the payload of the requests
// Used solely in the Web part of the service 
type Task struct {
	Name string `json:"name"`
	Command string `json:"command"`
	DependableTasks []string `json:"requires,omitempty"`
}

// Response body used for returning the execution plan
type ExecutionPlanResponse struct {
	Status int `json:"status"`
	Tasks []Task `json:"tasks"`
}


// Response body used for returning error messages to the clint
type ErrorResponse struct {
	ErrorCode int    `json:"errorcode"` //status code of the request - 4xx or 5xx
	ErrorMsg  string `json:"message"`   //desription of the error
}
