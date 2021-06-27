package http

const (
	ResponseCodeOk      = 200
	ResponseCodeCreated = 201
	BadRequest          = 400
	InternalServerError = 500
)


type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error"`
}

