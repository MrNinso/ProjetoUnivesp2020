PROJECT_ROOT := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
PROJECT_BUILD := $(PROJECT_ROOT)/build

build:
	mkdir $(PROJECT_BUILD)
	go build -tags=jsoniter -o "$(PROJECT_BUILD)/main" $(PROJECT_ROOT)/main.go & #Build Main bin
	pidMain=$!
	go build -tags=jsoniter -o  $(PROJECT_BUILD)/DBmain $(PROJECT_ROOT)/DBmain.go & # Build DataBaseTool
	pidDBmain=$!
	cp $(PROJECT_ROOT)/lab.service $(PROJECT_BUILD)
	cp $(PROJECT_ROOT)/buildCert.sh $(PROJECT_BUILD)
	cd $(PROJECT_BUILD); bash $(PROJECT_BUILD)/buildCert.sh # create SSL self-signed certificate)
	mkdir -p $(PROJECT_BUILD)/public/site
	cp -r $(PROJECT_ROOT)/public/res $(PROJECT_BUILD)/public/res
	cp -r $(PROJECT_ROOT)/public/site/build $(PROJECT_BUILD)/public/site/build
	wait $(pidMain)
	wait $(pidDBmain)

vagrant-build: build
	cp $(PROJECT_ROOT)/Vagrantfile $(PROJECT_BUILD)
	vagrant up

clean:
	rm -rf $(PROJECT_BUILD)
	rm -rf $(PROJECT_ROOT)/.vagrant