BINARY_NAME=terraform-provider-esdynamiconfig
DIRCMD=cd infrastructure/terraform;
TERRAFORM=$(DIRCMD) terraform 

build:
	go build
	mv $(BINARY_NAME) ~/.terraform.d/plugins/$(BINARY_NAME)

init: build
	$(TERRAFORM) init

plan:
	$(TERRAFORM) plan
	
apply:
	$(TERRAFORM) apply

destroy:
	$(TERRAFORM) destroy

clean:
	$(DIRCMD) rm -f terraform.txt terraform.tfstate terraform.tfstate.backup crash.log