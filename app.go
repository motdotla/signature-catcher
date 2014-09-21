package main

import (
	"github.com/go-martini/martini"
	"github.com/handshakejs/handshakejserrors"
	"github.com/joho/godotenv"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/motdotla/signaturelogic"
	"net/http"
	"os"
)

const (
	LOGIC_ERROR_CODE_UNKNOWN = "unknown"
)

var (
	ORCHESTRATE_API_KEY string
)

func CrossDomain() martini.Handler {
	return func(res http.ResponseWriter) {
		res.Header().Add("Access-Control-Allow-Origin", "*")
		res.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	}
}

type DocumentsProcessedPayload struct {
	Documents []Document `json:"documents"`
	Meta      *struct {
		Postscript string `json:"postscript,omitempty"`
	} `json:"meta,omitempty"`
}

type Document struct {
	Pages  []Page `json:"pages"`
	Status string `json:"status"`
	Url    string `json:"url"`
}

type Page struct {
	Number int    `json:"number"`
	Url    string `json:"url"`
}

func main() {
	loadEnvs()

	signaturelogic.Setup(ORCHESTRATE_API_KEY)

	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(CrossDomain())

	m.Any("/webhook/v0/documents/processed.json", binding.Bind(DocumentsProcessedPayload{}), DocumentsProcessed)

	m.Run()
}

func ErrorPayload(logic_error *handshakejserrors.LogicError) map[string]interface{} {
	error_object := map[string]interface{}{"code": logic_error.Code, "field": logic_error.Field, "message": logic_error.Message}
	errors := []interface{}{}
	errors = append(errors, error_object)
	payload := map[string]interface{}{"errors": errors}

	return payload
}

func DocumentsProcessed(documents_processed_payload DocumentsProcessedPayload, req *http.Request, r render.Render) {
	if len(documents_processed_payload.Documents) <= 0 {
		logic_error := &handshakejserrors.LogicError{"incorrect_payload", "", "the payload was in an unexpected format"}
		payload := ErrorPayload(logic_error)
		r.JSON(400, payload)
	} else {
		id := documents_processed_payload.Meta.Postscript
		pages := documents_processed_payload.Documents[0].Pages

		params := map[string]interface{}{"id": id, "pages": pages, "status": "processed"}
		_, logic_error := signaturelogic.DocumentsUpdate(params)
		if logic_error != nil {
			payload := ErrorPayload(logic_error)
			statuscode := determineStatusCodeFromLogicError(logic_error)
			r.JSON(statuscode, payload)
		} else {
			payload := map[string]interface{}{"success": true}
			r.JSON(200, payload)
		}
	}
}

func determineStatusCodeFromLogicError(logic_error *handshakejserrors.LogicError) int {
	code := 400
	if logic_error.Code == LOGIC_ERROR_CODE_UNKNOWN {
		code = 500
	}

	return code
}

func loadEnvs() {
	godotenv.Load()

	ORCHESTRATE_API_KEY = os.Getenv("ORCHESTRATE_API_KEY")
}
