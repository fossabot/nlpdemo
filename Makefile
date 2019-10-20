init:
	go mod download
	go get github.com/pilu/fresh
	cd  web/app && yarn 

clean:
	rm -rf .git
	touch .env

start_frontend:
	cd  web/app && yarn start 

start_backend:
	fresh

build:
	docker build -t rayyildiz/nlp -f deployments/backend/Dockerfile .
