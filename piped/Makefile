NAME = piped
REPO = hyperboloide

all:
	go build

lint:
	gometalinter ./...

container:
	GOOS=linux GOARCH=amd64 go build -o $(NAME)
	upx -7 $(NAME)
	docker build -t $(NAME) .
	docker tag $(NAME):latest $(REPO)/$(NAME):latest

push: container
	docker push $(REPO)/$(NAME):latest

clean:
	rm -fr $(NAME) $(NAME).upx

.PHONY: all lint container push clean
