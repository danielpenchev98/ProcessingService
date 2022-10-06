package rest

type RequestPayload struct {
	Tasks []Task `json:"tasks"`
}

type Task struct {
	Name string `json:"name"`
	Command string `json:"command"`
	DependableTasks []string `json:"requires,omitempty"`
}

type ExecutionPlanResponse struct {
	Status int `json:"status"`
	Tasks []Task `json:"tasks"`
}

type ErrorResponse struct {
	ErrorCode int    `json:"errorcode"` //status code of the request - 4xx or 5xx
	ErrorMsg  string `json:"message"`   //desription of the error
}
