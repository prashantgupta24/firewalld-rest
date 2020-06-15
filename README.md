# Firewalld-rest

This is a REST service to allow users to dynamically update firewalld rules on a server

## Table of Contents

<!-- @import "[TOC]" {cmd="toc" depthFrom=2 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [Table of Contents](#table-of-contents)
- [1. How to install and use](#1-how-to-install-and-use)
  - [Routes](#routes)
    - [Index page](#index-page)
      - [Sample query](#sample-query)
    - [Add new IP](#add-new-ip)
      - [Sample query](#sample-query-1)
    - [Show all IPs](#show-all-ips)
      - [Sample query](#sample-query-2)
    - [Show if IP is present](#show-if-ip-is-present)
      - [Sample query](#sample-query-3)
    - [Delete IP](#delete-ip)
      - [Sample query](#sample-query-4)
  - [IP struct](#ip-struct)
- [2. Helpful tips/links](#2-helpful-tipslinks)
  - [2.1 Kubernetes endpoint](#21-kubernetes-endpoint)
  - [2.2 Firewalld](#22-firewalld)
    - [2.2.1 Commands](#221-commands)
    - [2.2.2 Rich rules](#222-rich-rules)
    - [2.2.3 Documentation](#223-documentation)
  - [2.3 JWT in Go](#23-jwt-in-go)
  - [2.4 Golang Exec](#24-golang-exec)
  - [2.5 Systemd](#25-systemd)
  - [2.6 Commands for generating public/private key](#26-commands-for-generating-publicprivate-key)

<!-- /code_chunk_output -->

## 1. How to install and use

### Routes

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
curl --location --request GET '<SERVER_IP>:8080' \
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
curl --location --request POST '<SERVER_IP>:8080/ip' \
--header 'Authorization: Bearer <signed_jwt>' \
--header 'Content-Type: application/json' \
--data-raw '{"ip":"10.xx.xx.xx","domain":"example.com"}'
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
curl --location --request GET '<SERVER_IP>:8080/ip' \
--header 'Authorization: Bearer <signed_jwt>'
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
curl --location --request GET '<SERVER_IP>:8080/ip/10.xx.xx.xx' \
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
curl --location --request DELETE '<SERVER_IP>:8080/ip/10.xx.xx.xx' \
--header 'Authorization: Bearer <signed_jwt>'
```

### IP struct

```
type IP struct {
	IP     string `json:"ip"`
	Domain string `json:"domain"`
}
```

## 2. Helpful tips/links

### 2.1 Kubernetes endpoint

- https://theithollow.com/2019/02/04/kubernetes-endpoints/

### 2.2 Firewalld

#### 2.2.1 Commands

```
firewall-cmd --get-default-zone
firewall-cmd --get-active-zones

firewall-cmd --list-all-zones | less

firewall-cmd --zone=internal --list-sources
firewall-cmd --zone=internal --list-services
firewall-cmd --zone=internal --list-all

firewall-cmd --zone=public --add-service=ssh --permanent

firewall-cmd --zone=internal --add-source=73.223.28.39/32 --permanent

firewall-cmd --reload
```

#### 2.2.2 Rich rules

`firewall-cmd --permanent --zone=public --list-rich-rules`

`firewall-cmd --permanent --zone=public --add-rich-rule='rule family="ipv4" source address="10.10.99.10/32" port protocol="tcp" port="22" accept'`

`firewall-cmd --permanent --zone=public --add-rich-rule='rule family="ipv4" source address="192.168.100.0/24" invert="True" drop'`

> Reject will reply back with an ICMP packet noting the rejection, while a drop will just silently drop the traffic and do nothing else, so a drop may be preferable in terms of security as a reject response confirms the existence of the system as it is rejecting the request.

#### 2.2.3 Documentation

> --add-source=IP can be used to add an IP address or range of addresses to a zone. This will mean that if any source traffic enters the systems that matches this, the zone that we have set will be applied to that traffic. In this case we set the ‘testing’ zone to be associated with traffic from the 10.10.10.0/24 range.

`[root@centos7 ~]# firewall-cmd --permanent --zone=testing --add-source=10.10.10.0/24`
success

- https://www.digitalocean.com/community/tutorials/how-to-set-up-a-firewall-using-firewalld-on-centos-7

### 2.3 JWT in Go

- https://www.thepolyglotdeveloper.com/2017/03/authenticate-a-golang-api-with-json-web-tokens/

### 2.4 Golang Exec

- https://stackoverflow.com/questions/39151420/golang-executing-command-with-spaces-in-one-of-the-parts

### 2.5 Systemd

- https://medium.com/@benmorel/creating-a-linux-service-with-systemd-611b5c8b91d6
- https://www.digitalocean.com/community/tutorials/understanding-systemd-units-and-unit-files
- [See logs using journalctl](https://www.linode.com/docs/quick-answers/linux/how-to-use-journalctl/)

### 2.6 Commands for generating public/private key

```
openssl genrsa -key private-key-sc.pem
openssl req -new -x509 -key private-key-sc.pem -out public.cert
```
