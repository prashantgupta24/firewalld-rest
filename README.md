[![Go Report Card](https://goreportcard.com/badge/github.com/prashantgupta24/firewalld-rest)](https://goreportcard.com/report/github.com/prashantgupta24/firewalld-rest) [![codecov](https://codecov.io/gh/prashantgupta24/firewalld-rest/branch/master/graph/badge.svg)](https://codecov.io/gh/prashantgupta24/firewalld-rest)

# Firewalld-rest

A REST service to allow users to dynamically update firewalld rules on a server.

## What it does

The simple idea behind this repo is to have a secure system, a system running `Firewalld` that does not permit SSH access to any IP address by default so there are no brute-force attacks. The only way to access the system is by communicating with a REST server running on the system through a valid request containing your public IP address.

The REST server validates your request (it checks for a signed JWT, covered later), and if the request is valid, it will add your IP to the `firewalld` rule for the `public` zone for SSH, which gives your IP SSH access to the machine.

Once you are done using the machine, you can remove your IP using the same REST server, and the server shuts SSH access off again.

## Table of Contents

<!-- @import "[TOC]" {cmd="toc" depthFrom=2 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [What it does](#what-it-does)
- [Table of Contents](#table-of-contents)
- [Pre-requisites](#pre-requisites)
- [About the application](#about-the-application)
  - [Authorization](#authorization)
  - [DB](#db)
  - [Routing](#routing)
  - [Tests](#tests)
- [How to install and use on server](#how-to-install-and-use-on-server)
  - [Local changes required](#local-changes-required)
  - [Build the application](#build-the-application)
  - [Copy binary file over to server](#copy-binary-file-over-to-server)
  - [Remove SSH from public firewalld zone](#remove-ssh-from-public-firewalld-zone)
  - [Expose the REST server](#expose-the-rest-server)
  - [Configure linux systemd service](#configure-linux-systemd-service)
  - [Start and enable systemd service.](#start-and-enable-systemd-service)
  - [Interacting with the REST server](#interacting-with-the-rest-server)
    - [Index page](#index-page)
      - [Sample query](#sample-query)
    - [Show all IPs](#show-all-ips)
      - [Sample query](#sample-query-1)
    - [Add new IP](#add-new-ip)
      - [Sample query](#sample-query-2)
    - [Show if IP is present](#show-if-ip-is-present)
      - [Sample query](#sample-query-3)
    - [Delete IP](#delete-ip)
      - [Sample query](#sample-query-4)
  - [IP struct](#ip-struct)
- [Helpful tips/links](#helpful-tipslinks)
  - [Creating custom kubernetes endpoint](#creating-custom-kubernetes-endpoint)
  - [Firewalld rules](#firewalld-rules)
    - [Useful commands](#useful-commands)
    - [Rich rules](#rich-rules)
    - [Misc tips](#misc-tips)
  - [Using JWT in Go](#using-jwt-in-go)
  - [Using golang Exec](#using-golang-exec)
  - [Systemd services](#systemd-services)
  - [Commands for generating public/private key](#commands-for-generating-publicprivate-key)

<!-- /code_chunk_output -->

## Pre-requisites

This repo assumes you have:

1. A linux server with `firewalld` installed.
1. `root` access to the machine. (without `root` access, the application will not be able to run the `firewall-cmd` commands needed to add the rule for SSH access)
1. Kubernetes running on the system (so that the REST server can be exposed outside)

## About the application

### Authorization

### DB

### Routing

### Tests

## How to install and use on server

### Local changes required

1. Make sure you have updated the [publicCert.go](https://github.com/prashantgupta24/firewalld-rest/blob/master/route/publicCert.go) with your own public cert for which you have the private key. See the section on [generating your own public/private key](#commands-for-generating-publicprivate-key). Once you have your own public and private key pair, then you can go to jwt.io and generate a valid signed JWT using `RS256 algorithm` (the payload doesn't matter). You will be using that JWT to make calls to the REST server.

1. Make sure you update the path to where you want to keep your binary on the server in the [Linux systemd service](#configure-linux-systemd-service). In the definition I have here, it assumes you have kept it in `/root/rest/firewalld-rest`. If **not**, make sure to change the service definition (covered again later).

### Build the application

Run the command:

```
make build-linux DB_PATH=/dir/to/db/
```

It will create a binary under the build directory, called `firewalld-rest`. The `DB_PATH=/dir/to/keep/db` statement sets the path where the `.db` file will be saved. It should be saved in a protected location such that it is not accidentally deleted on server restart or by any other user. A good place for it would be in the same directory where you will copy the binary over to (in the next step). That way you will not forget where it is.

If `DB_PATH` variable is not set, the db file will be created by default under `/`. (_This happens because the binary is run by systemd. If we manually ran the binary file on the server, the db file would be created in the same directory._)

Once the binary is built, it should contain everything required to run the application on a linux based server.

### Copy binary file over to server

```
scp build/firewalld-rest root@<server>:/root/rest
```

_Note_: if you want to change the directory where you want to keep the binary, then make sure you edit the `firewalld-rest.service` file, as the `linux systemd service` definition example in this repo expects the location of the binary to be `/root/rest`.

### Remove SSH from public firewalld zone

This is to remove SSH access from the public zone, which will cease SSH access from everywhere.

```
firewall-cmd --zone=public --remove-service=ssh --permanent
```

then reload (since we are using `--permanent`):

```
firewall-cmd --reload
```

This removes ssh access for everyone. This is where the application comes into play, and we enable access based on IP.

**Confirmirmation for the step**:

```
firewall-cmd --zone=public --list-all
```

_Notice the `ssh` service will not be listed in public zone anymore._

### Expose the REST server

The REST server can be exposed in a number of different ways, I have 2 examples on how it can be exposed:

1. Using a `NodePort` kubernetes service ([link](https://github.com/prashantgupta24/firewalld-rest/blob/master/k8s/svc-nodeport.yaml))
2. Using `ingress` along with a kubernetes service ([link](https://github.com/prashantgupta24/firewalld-rest/blob/master/k8s/ingress.yaml))

### Configure linux systemd service

See [this](https://github.com/prashantgupta24/firewalld-rest/blob/master/firewalld-rest.service) for an example of a linux systemd service.

**Note**: This service assumes your binary is at `/root/rest/firewalld-rest`. You can change that in the file above.

### Start and enable systemd service.

**Start**

```
systemctl start firewalld-rest
```

**Enable**

```
systemctl enable firewalld-rest
```

**Logs**

You can see the logs for the service using:

```
journalctl -f
```

### Interacting with the REST server

#### Index page

```
route{
    "Index Page",
    "GET",
    "/",
}
```

##### Sample query

```
curl --location --request GET '<SERVER_IP>:8080/m1' \
--header 'Authorization: Bearer <signed_jwt>'
```

#### Show all IPs

```
route{
    "Show all IPs present",
    "GET",
    "/ip",
}
```

##### Sample query

```
curl --location --request GET '<SERVER_IP>:8080/m1/ip' \
--header 'Authorization: Bearer <signed_jwt>'
```

#### Add new IP

```
route{
    "Add New IP",
    "POST",
    "/ip",
}
```

##### Sample query

```
curl --location --request POST '<SERVER_IP>:8080/m1/ip' \
--header 'Authorization: Bearer <signed_jwt>' \
--header 'Content-Type: application/json' \
--data-raw '{"ip":"10.xx.xx.xx","domain":"example.com"}'
```

#### Show if IP is present

```
route{
    "Show if particular IP is present",
    "GET",
    "/ip/{ip}",
}
```

##### Sample query

```
curl --location --request GET '<SERVER_IP>:8080/m1/ip/10.xx.xx.xx' \
--header 'Authorization: Bearer <signed_jwt>'
```

#### Delete IP

```
route{
    "Delete IP",
    "DELETE",
    "/ip/{ip}",
}
```

##### Sample query

```
curl --location --request DELETE '<SERVER_IP>:8080/m1/ip/10.xx.xx.xx' \
--header 'Authorization: Bearer <signed_jwt>'
```

### IP struct

```
type IP struct {
	IP     string `json:"ip"`
	Domain string `json:"domain"`
}
```

## Helpful tips/links

### Creating custom kubernetes endpoint

- https://theithollow.com/2019/02/04/kubernetes-endpoints/

### Firewalld rules

- https://www.digitalocean.com/community/tutorials/how-to-set-up-a-firewall-using-firewalld-on-centos-7

#### Useful commands

```
firewall-cmd --get-default-zone
firewall-cmd --get-active-zones

firewall-cmd --list-all-zones | less

firewall-cmd --zone=public --list-sources
firewall-cmd --zone=public --list-services
firewall-cmd --zone=public --list-all

firewall-cmd --zone=public --add-service=ssh --permanent

firewall-cmd --zone=internal --add-source=70.xx.xx.xxx/32 --permanent

firewall-cmd --reload
```

#### Rich rules

`firewall-cmd --permanent --zone=public --list-rich-rules`

`firewall-cmd --permanent --zone=public --add-rich-rule='rule family="ipv4" source address="10.10.99.10/32" port protocol="tcp" port="22" accept'`

`firewall-cmd --permanent --zone=public --add-rich-rule='rule family="ipv4" source address="192.168.100.0/24" invert="True" drop'`

> Reject will reply back with an ICMP packet noting the rejection, while a drop will just silently drop the traffic and do nothing else, so a drop may be preferable in terms of security as a reject response confirms the existence of the system as it is rejecting the request.

#### Misc tips

> --add-source=IP can be used to add an IP address or range of addresses to a zone. This will mean that if any source traffic enters the systems that matches this, the zone that we have set will be applied to that traffic. In this case we set the ‘testing’ zone to be associated with traffic from the 10.10.10.0/24 range.

```
[root@centos7 ~]# firewall-cmd --permanent --zone=testing --add-source=10.10.10.0/24
success
```

### Using JWT in Go

- https://www.thepolyglotdeveloper.com/2017/03/authenticate-a-golang-api-with-json-web-tokens/

### Using golang Exec

- https://stackoverflow.com/questions/39151420/golang-executing-command-with-spaces-in-one-of-the-parts

### Systemd services

- https://medium.com/@benmorel/creating-a-linux-service-with-systemd-611b5c8b91d6
- https://www.digitalocean.com/community/tutorials/understanding-systemd-units-and-unit-files
- [Logs using journalctl](https://www.linode.com/docs/quick-answers/linux/how-to-use-journalctl/)

### Commands for generating public/private key

```
openssl genrsa -key private-key-sc.pem
openssl req -new -x509 -key private-key-sc.pem -out public.cert
```
