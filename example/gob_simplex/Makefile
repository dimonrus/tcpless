help:	    	## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

server:		## Start server
	go build -o json && env ENV=local ./json -app=server

client:		## Start client
	go build -o json && env ENV=local ./json -app=client