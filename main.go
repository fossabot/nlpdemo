package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"rayyildiz.dev/nlp/pkg/language"
	"rayyildiz.dev/nlp/pkg/sentiment"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"rayyildiz.dev/nlp/internal/infra"
)

func init() {
	godotenv.Load()
}

func main() {
	var spec infra.Specification
	err := envconfig.Process("nlp", &spec)
	if err != nil {
		fmt.Printf("could not load config, %v", err)
		os.Exit(1)
	}

	log := infra.NewLogger(spec.Debug)
	defer sentry.Flush(time.Second * 5)

	daataset := readDataset("./data/yelp_labelled.txt", log)
	// db, err := infra.NewDatabase(spec.PostgresConnection)
	// if err != nil {
	// 	log.Fatal("could not initialize database", zap.Error(err))
	// }
	// defer db.Close()

	r := infra.NewRouter()
	port := fmt.Sprintf("%d", spec.Port)
	if port == "" {
		port = "4000"
	}

	pwd, err := os.Getwd()
	if err != nil {
		sentry.CaptureException(err)
		log.Error("could not get the working dir", zap.Error(err))
		os.Exit(1)
	}
	log.Info("working dir", zap.String("pwd", pwd))

	filesDir := filepath.Join(pwd, "static")
	staticFileServer(r, "/", http.Dir(filesDir))

	// Handlers
	r.Route("/api", func(r chi.Router) {
		r.Route("/sentiment", func(r chi.Router) {
			sentiment.RegisterHandler(r, log, daataset)
		})

		r.Route("/language", func(r chi.Router) {
			language.RegisterHandler(r, log)
		})
	})

	log.Info("server is starting", zap.String("port", port))
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Error("could not start server", zap.String("port", port), zap.Error(err))
	}
}

func staticFileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "max-age=2592000")
		fs.ServeHTTP(w, r)
	})
}

func readDataset(file string, log *zap.Logger) map[string]string {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	dataset := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		l := scanner.Text()
		data := strings.Split(l, "\t")
		if len(data) != 2 {
			continue
		}
		sentence := data[0]
		if data[1] == "0" {
			dataset[sentence] = "negative"
		} else if data[1] == "1" {
			dataset[sentence] = "positive"
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("could not read dataset", zap.Error(err))
	}
	return dataset
}
