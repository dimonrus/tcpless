help:	    	## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

server:		## Start server
	go build -o simplex && env ENV=local ./simplex -app=server

client:		## Start client
	go build -o simplex && env ENV=local ./simplex -app=client