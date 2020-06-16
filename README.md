# Firewalld-rest

A REST service to allow users to dynamically update firewalld rules on a server.

## Table of Contents

<!-- @import "[TOC]" {cmd="toc" depthFrom=2 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [Table of Contents](#table-of-contents)
- [The idea](#the-idea)
- [2. How to install and use](#2-how-to-install-and-use)
  - [2.1 Remove SSH from public zone](#21-remove-ssh-from-public-zone)
  - [2.2 Copy build file over to machine.](#22-copy-build-file-over-to-machine)
  - [2.3 Configure k8s service and ingress.](#23-configure-k8s-service-and-ingress)
  - [2.4 Configure linux systemd service](#24-configure-linux-systemd-service)
  - [2.5 Start and enable systemd service.](#25-start-and-enable-systemd-service)
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
    - [2.2.1 Useful commands](#221-useful-commands)
    - [2.2.2 Rich rules](#222-rich-rules)
    - [2.2.3 Misc tips](#223-misc-tips)
  - [2.3 JWT in Go](#23-jwt-in-go)
  - [2.4 Golang Exec](#24-golang-exec)
  - [2.5 Systemd](#25-systemd)
  - [2.6 Commands for generating public/private key](#26-commands-for-generating-publicprivate-key)

<!-- /code_chunk_output -->

## The idea

The simple idea behind this repo is to have a system running `Firewalld` that does not permit SSH access to anyone by default. The only way to access the system is by communicating with a REST server running on the system, by sending a valid request containing your public IP address.

The REST server validates your request (it checks for a signed JWT, covered later), and if the request is valid, it will add your IP to the firewalld rule for the public zone for SSH service, which gives you SSH access for the machine.

Once you are done using the machine, you can remove your IP using the same REST server, and the server shuts itself off from SSH access from everyone again.

## 2. How to install and use

### 2.1 Remove SSH from public zone

The first step is to remove SSH access from the public zone, which will cease SSH access from everywhere.

```
firewall-cmd --zone=public --remove-service=ssh --permanent
```

This removes ssh access for everyone. This is where the application comes into play, and we enable access based on IP.

**Confirm**:

```
firewall-cmd --zone=public --list-all
```

### 2.2 Copy build file over to machine.

### 2.3 Configure k8s service and ingress.

See the sample `ingress.yaml` and the `svc.yaml` inside the `k8s` folder to get an idea.

### 2.4 Configure linux systemd service

### 2.5 Start and enable systemd service.

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
curl --location --request GET '<SERVER_IP>:8080/m1' \
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

## 2. Helpful tips/links

### 2.1 Kubernetes endpoint

- https://theithollow.com/2019/02/04/kubernetes-endpoints/

### 2.2 Firewalld

- https://www.digitalocean.com/community/tutorials/how-to-set-up-a-firewall-using-firewalld-on-centos-7

#### 2.2.1 Useful commands

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

#### 2.2.2 Rich rules

`firewall-cmd --permanent --zone=public --list-rich-rules`

`firewall-cmd --permanent --zone=public --add-rich-rule='rule family="ipv4" source address="10.10.99.10/32" port protocol="tcp" port="22" accept'`

`firewall-cmd --permanent --zone=public --add-rich-rule='rule family="ipv4" source address="192.168.100.0/24" invert="True" drop'`

> Reject will reply back with an ICMP packet noting the rejection, while a drop will just silently drop the traffic and do nothing else, so a drop may be preferable in terms of security as a reject response confirms the existence of the system as it is rejecting the request.

#### 2.2.3 Misc tips

> --add-source=IP can be used to add an IP address or range of addresses to a zone. This will mean that if any source traffic enters the systems that matches this, the zone that we have set will be applied to that traffic. In this case we set the ‘testing’ zone to be associated with traffic from the 10.10.10.0/24 range.

```
[root@centos7 ~]# firewall-cmd --permanent --zone=testing --add-source=10.10.10.0/24
success
```

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
