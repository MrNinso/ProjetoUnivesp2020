build:
	mkdir build
	go build -tags=jsoniter main.go -o build/main & #Build Main bin
	pids[0]=$!
	go build -tags=jsoniter BDmain.go -o build/DBmain & # Build DataBaseTool
	pids[1]=$!
	cp buildCert.sh build/
	bash build/buildCert.sh # create SSL self-signed certificate
	for pid in ${pids[*]}; do wait $pid done # wait for gobuild


vagrant-up: build

