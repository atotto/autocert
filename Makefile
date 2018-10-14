
build:
	docker build . -t atotto/autocert:latest

deploy: build
	docker push atotto/autocert:latest
