NAME			= rodney

GO_RUN		= go run
GO_BUILD	= env GOOS=linux GOARCH=arm GOARM=6 go build

SRCS			= main.go

RM				= rm -rf
SCP				= scp
SSH				=	ssh

PI_USER		= pi
PI_HOST		= 192.168.0.36

.phony: build transfer run fclean

build:
	$(GO_BUILD) -o $(NAME) $(SRCS)

transfer: build
	$(SCP) $(NAME) $(PI_USER)@$(PI_HOST):$(NAME)

run: transfer
	$(SSH) $(PI_USER)@$(PI_HOST) "bash -c ./$(NAME)"

fclean:
	$(RM) $(NAME)