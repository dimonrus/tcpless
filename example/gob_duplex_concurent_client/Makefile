help:	    	## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

server:		## Start server
	go build -o concur && env ENV=local ./concur -app=server

client:		## Start client
	go build -o concur && env ENV=local ./concur -app=client