OUTDIR=.

feedme: feedme-webapp feedme-cli

feedme-webapp: ./cmd/webapp/main.go
	go build -o $(OUTDIR)/feedme-webapp ./cmd/webapp/main.go

feedme-cli: ./cmd/cli/main.go
	go build -o $(OUTDIR)/feedme-cli ./cmd/cli/main.go

clean:
	rm -f $(OUTDIR)/feedme-cli
	rm -f $(OUTDIR)/feedme-webapp