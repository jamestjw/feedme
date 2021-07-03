OUTDIR=.

feedme: feedme-webapp feedme-cli

feedme-webapp: ./cmd/webapp/webapp.go
	go build -o $(OUTDIR)/feedme-webapp ./cmd/webapp/webapp.go

feedme-cli: ./cmd/cli/cli.go
	go build -o $(OUTDIR)/feedme-cli ./cmd/cli/cli.go

clean:
	rm -f $(OUTDIR)/feedme-cli
	rm -f $(OUTDIR)/feedme-webapp