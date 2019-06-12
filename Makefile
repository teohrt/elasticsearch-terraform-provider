BINARY_NAME=terraform-provider-esdynamiconfig

build:
	go build
	mv $(BINARY_NAME) ~/.terraform.d/plugins/$(BINARY_NAME)