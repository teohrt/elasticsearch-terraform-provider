BINARY_NAME=terraform-provider-esdynamiconfig

build:
	go build
	mv $(BINARY_NAME) ~/.terraform.d/plugins/$(BINARY_NAME)

init: build
	cd infrastructure/terraform; terraform init

plan:
	cd infrastructure/terraform; terraform plan