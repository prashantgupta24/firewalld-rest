## Firewalld-rest

This is a REST service to allow users to dynamically update firewalld rules on a server

## Good links:

- https://theithollow.com/2019/02/04/kubernetes-endpoints/
- https://www.digitalocean.com/community/tutorials/how-to-set-up-a-firewall-using-firewalld-on-centos-7
- https://www.thepolyglotdeveloper.com/2017/03/authenticate-a-golang-api-with-json-web-tokens/
- https://stackoverflow.com/questions/39151420/golang-executing-command-with-spaces-in-one-of-the-parts

## Commands for generating public/private key

```
openssl genrsa -key private-key-sc.pem
openssl req -new -x509 -key private-key-sc.pem -out publickey-sc.cer
```
