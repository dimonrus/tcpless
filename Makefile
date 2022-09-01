help:	    	## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

cert:		## Generate tls certs for server and client. Required OpenSSL 3.0.0
	@read -p "Enter domain: " domain; \
	mkdir resource; \
    rm resource/*; \
    openssl genrsa -out resource/ca.pem && openssl req -new -x509 -nodes -days 825 -key resource/ca.pem -out resource/ca.crt -subj "/C=RU/ST=MSK/L=Earth/O=$$domain/OU=IT/CN=www.$$domain/emailAddress=email@$$domain" -addext "subjectAltName = IP:0.0.0.0"; \
    openssl req -new -nodes -x509 -out resource/server.crt -keyout resource/server.pem -CA resource/ca.crt -CAkey resource/ca.pem -days 825 -subj "/C=RU/ST=MSK/L=Earth/O=$$domain/OU=IT/CN=www.$$domain/emailAddress=email@$$domain" -addext "subjectAltName = IP:0.0.0.0"; \
    openssl req -new -nodes -x509 -out resource/client.crt -keyout resource/client.pem -CA resource/ca.crt -CAkey resource/ca.pem -days 825 -subj "/C=RU/ST=MSK/L=Earth/O=$$domain/OU=IT/CN=www.$$domain/emailAddress=email@$$domain" -addext "subjectAltName = IP:0.0.0.0"; \
    chmod 0700 resource/*

info:		## Read certificate
	@read -p "Enter cert path: " path; \
	openssl x509 -in $$path -text -noout
