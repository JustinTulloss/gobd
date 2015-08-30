SOURCES=cmd/main.go ./gobd.go

gobd: $(SOURCES)
	go build cmd/main.go -o gobd

arm-gobd: $(SOURCES)
	GOARM=7 GOOS=linux GOARCH=arm go build -o arm-gobd cmd/main.go

clean:
	rm -f gobd arm-gobd

run-pi: arm-gobd
	scp arm-gobd pi@pirunner.local:
	ssh pi@pirunner.local ./arm-gobd /dev/rfcomm0
