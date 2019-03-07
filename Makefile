    GOCMD=go
    GOBUILD=$(GOCMD) build
    GOCLEAN=$(GOCMD) clean
    GOTEST=$(GOCMD) test
    GOGET=$(GOCMD) get
    BINARY_UNIX=rsasss
    BINARY_WIN=$(BINARY_NAME).exe
    
    all: test build install
    build: 
            $(GOBUILD) -o $(BINARY_UNIX) -v
    install:
			build
			cp $(BINARY_UNIX) $(GOHOME)/bin
    test: 
            $(GOTEST) ./...
    clean: 
            $(GOCLEAN)
			rm -f $(GOHOME)/bin/$(BINARY_UNIX)
            rm -f $(BINARY_UNIX)
            rm -f $(BINARY_WIN)
    run:
            $(GOBUILD) -o $(BINARY_UNIX) -v ./...
            ./$(BINARY_UNIX)
    deps:
            $(GOGET) github.com/spf13/cobra

    # cross compiling
    build-linux:
			CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v			
    build-win:
			CGO_ENABLED=0 GOOS=windows GOARCH=adm64 $(GOBUILD) -o $(BINARY_WIN) -v