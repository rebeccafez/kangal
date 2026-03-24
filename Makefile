.PHONY: clean

kangal:
	go build -o kangal cmd/kangal/main.go

clean:
	@rm kangal
