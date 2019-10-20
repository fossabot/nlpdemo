package language

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
	"rayyildiz.dev/nlp/internal/infra"
)

func RegisterHandler(r chi.Router, log *zap.Logger) {

	r.Post("/detect", handleLangDetection(log))
}

func handleLangDetection(log *zap.Logger) http.HandlerFunc {
	type req struct {
		Text string `json:"text"`
	}

	type resp struct {
		Language   string  `json:"language"`
		Confidence float64 `json:"confidence"`
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		var r req
		err := render.DecodeJSON(request.Body, &r)
		if err != nil {
			infra.CaptureException(request, err)
			log.Error("could not decode request", zap.Error(err))
			http.Error(writer, err.Error(), http.StatusNotAcceptable)
			return
		}

		lang, conf, err := detectLanguage(r.Text)
		if err != nil {
			infra.CaptureException(request, err)
			log.Error("could not decode request", zap.Error(err))
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		render.JSON(writer, request, resp{
			Language:   lang,
			Confidence: conf,
		})
	}
}
