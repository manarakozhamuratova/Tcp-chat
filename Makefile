build:
	docker build -t forum .
run:
	docker run -d --name forum-app -p8080:8080 forum