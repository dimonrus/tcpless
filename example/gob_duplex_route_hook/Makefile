help:	    	## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

server:		## Start server
	go build -o hook && env ENV=local ./hook -app=server

client:		## Start client
	go build -o hook && env ENV=local ./hook -app=client