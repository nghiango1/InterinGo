mkdir ~/.ssl

openssl req -x509 -nodes -days 365 \
 -newkey rsa:2048 \
 -keyout ~/.ssl/localhost.key \
 -out ~/.ssl/localhost.crt \
 -config localhost.cnf \
 -extensions ext
