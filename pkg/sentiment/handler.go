package sentiment

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
	"rayyildiz.dev/nlp/internal/infra"
)

func RegisterHandler(r chi.Router, log *zap.Logger, dataset map[string]string) {
	nb := newClassifier()
	nb.train(dataset)

	r.Post("/", handleLSentimentDetection(log, nb))

}

func handleLSentimentDetection(log *zap.Logger, nb *classifier) http.HandlerFunc {
	type req struct {
		Text string `json:"text"`
	}

	type resp struct {
		Sentiment  string  `json:"sentiment"`
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

		result := nb.classify(r.Text)
		if result[positive] > result[negative] {
			render.JSON(writer, request, resp{positive, result[positive]})
		} else {
			render.JSON(writer, request, resp{negative, result[negative]})
		}
	}
}
