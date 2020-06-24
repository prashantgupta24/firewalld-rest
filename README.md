[![Go Report Card](https://goreportcard.com/badge/github.com/prashantgupta24/firewalld-rest)](https://goreportcard.com/report/github.com/prashantgupta24/firewalld-rest) [![codecov](https://codecov.io/gh/prashantgupta24/firewalld-rest/branch/master/graph/badge.svg)](https://codecov.io/gh/prashantgupta24/firewalld-rest) [![version][version-badge]][releases]

# Firewalld-rest

A REST application to dynamically update firewalld rules on a linux server

## Purpose

The simple idea behind this is to have a completely isolated system, a system running Firewalld that does not permit SSH access to any IP address by default so there are no brute-force attacks. The only way to access the system is by communicating with a REST application running on the server through a valid request containing your public IP address.

The REST application validates your request (it checks for a valid JWT, covered later), and if the request is valid, it will add your IP to the firewalld rule for the public zone for SSH, which gives **only your IP** SSH access to the machine.

Once you are done using the machine, you can remove your IP interacting with the same REST application, and it changes rules in firewalld, shutting off SSH access and isolating the system again.

## Table of Contents

<!-- @import "[TOC]" {cmd="toc" depthFrom=2 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [Purpose](#purpose)
- [Table of Contents](#table-of-contents)
- [1. Pre-requisites](#1-pre-requisites)
- [2. About the application](#2-about-the-application)
  - [2.1 Authorization](#21-authorization)
  - [2.2 DB](#22-db)
  - [2.3 Tests](#23-tests)
- [3. How to install and use on server](#3-how-to-install-and-use-on-server)
  - [3.1 Generate JWT](#31-generate-jwt)
  - [3.2 Build the application](#32-build-the-application)
  - [3.3 Copy binary file over to server](#33-copy-binary-file-over-to-server)
  - [3.4 Remove SSH from public firewalld zone](#34-remove-ssh-from-public-firewalld-zone)
  - [3.5 Expose the REST application](#35-expose-the-rest-application)
    - [3.5.1 Single node cluster](#351-single-node-cluster)
    - [3.5.2 Multi-node cluster](#352-multi-node-cluster)
  - [3.6 Configure linux systemd service](#36-configure-linux-systemd-service)
  - [3.7 Start and enable systemd service.](#37-start-and-enable-systemd-service)
  - [3.8 IP JSON](#38-ip-json)
  - [3.9 Interacting with the REST application](#39-interacting-with-the-rest-application)
    - [3.9.1 Index page](#391-index-page)
      - [Sample query](#sample-query)
    - [3.9.2 Show all IPs](#392-show-all-ips)
      - [Sample query](#sample-query-1)
    - [3.9.3 Add new IP](#393-add-new-ip)
      - [Sample query](#sample-query-2)
    - [3.9.4 Show if IP is present](#394-show-if-ip-is-present)
      - [Sample query](#sample-query-3)
    - [3.9.5 Delete IP](#395-delete-ip)
      - [Sample query](#sample-query-4)
- [4. Helpful tips/links](#4-helpful-tipslinks)
- [5. Commands for generating public/private key](#5-commands-for-generating-publicprivate-key)

<!-- /code_chunk_output -->

## 1. Pre-requisites

This repo assumes you have:

1. A linux server with `firewalld` installed.
1. `root` access to the server. (without `root` access, the application will not be able to run the `firewall-cmd` commands needed to add the rule for SSH access)
1. Some way of exposing the application externally (there are examples in this repo on how to use Kubernetes to expose the service)

## 2. About the application

### 2.1 Authorization

The application uses `RS256` type algorithm to verify the incoming requests.

> RS256 (RSA Signature with SHA-256) is an asymmetric algorithm, and it uses a public/private key pair: the identity provider has a private (secret) key used to generate the signature, and the consumer of the JWT gets a public key to validate the signature.

The public certificate is in this file [publicCert.go](https://github.com/prashantgupta24/firewalld-rest/blob/master/route/publicCert.go), which is something that will have to be changed before you can use it. (more information on how to create a new one later).

### 2.2 DB

The application uses a file DB for now. The architecture allows easy integration of any other type of DB. The interface in `db.go` is what is required to be fulfilled to introduce a new type of DB.

### 2.3 Tests

The test can be run using `make test`. The emphasis has been given to testing the handler functions and making sure that IPs get added and removed successfully from the DB. I still have to figure out how to actually automate the tests for the firewalld rules (contributions are welcome!)

## 3. How to install and use on server

### 3.1 Generate JWT

Update the file [publicCert.go](https://github.com/prashantgupta24/firewalld-rest/blob/master/route/publicCert.go) with your own `public cert` for which you have the private key.

If you want to create a new set, see the section on [generating your own public/private key](#5-commands-for-generating-publicprivate-key). Once you have your own public and private key pair, then after updating the file above, you can go to `jwt.io` and generate a valid JWT using `RS256 algorithm` (the payload doesn't matter). You will be using that JWT to make calls to the REST application, so keep the JWT safe.

### 3.2 Build the application

Run the command:

```
make build-linux DB_PATH=/dir/to/db/
```

It will create a binary under the build directory, called `firewalld-rest`. The `DB_PATH=/dir/to/keep/db` statement sets the path where the `.db` file will be saved **on the server**. It should be saved in a protected location such that it is not accidentally deleted on server restart or by any other user. A good place for it could be the same directory where you will copy the binary over to (in the next step). That way you will not forget where it is.

If `DB_PATH` variable is not set, the db file will be created by default under `/`. (_This happens because the binary is run by systemd. If we manually ran the binary file on the server, the db file would be created in the same directory._)

Once the binary is built, it should contain everything required to run the application on a linux based server.

### 3.3 Copy binary file over to server

```
scp build/firewalld-rest root@<server>:/root/rest
```

_Note_: if you want to change the directory where you want to keep the binary, then make sure you edit the `firewalld-rest.service` file, as the `linux systemd service` definition example in this repo expects the location of the binary to be `/root/rest`.

### 3.4 Remove SSH from public firewalld zone

This is to remove SSH access from the public zone, which will cease SSH access from everywhere.

SSH into the server, and run the following command:

```
firewall-cmd --zone=public --remove-service=ssh --permanent
```

then reload (since we are using `--permanent`):

```
firewall-cmd --reload
```

This removes ssh access for everyone. This is where the application will come into play, and we enable access based on IP.

**Confirmation for the step**:

```
firewall-cmd --zone=public --list-all
```

_Notice the `ssh` service will not be listed in public zone anymore._

Also try SSH access into the server from another terminal. It should reject the attempt.

### 3.5 Expose the REST application

The REST application can be exposed in a number of different ways, I have 2 examples on how it can be exposed:

1. Using a `NodePort` kubernetes service ([link](https://github.com/prashantgupta24/firewalld-rest/blob/master/k8s/svc-nodeport.yaml))
2. Using `ingress` along with a kubernetes service ([link](https://github.com/prashantgupta24/firewalld-rest/blob/master/k8s/ingress.yaml))

#### 3.5.1 Single node cluster

For a single-node cluster, see the kubernetes service example [here](https://github.com/prashantgupta24/firewalld-rest/blob/master/k8s/svc-nodeport.yaml). The important thing to note is that we manually add the `Endpoints` resource for the service, which points to our node's private IP address and port `8080`.

Once deployed, your service might look like this:

```
kubernetes get svc

external-rest | NodePort | 10.xx.xx.xx | 169.xx.xx.xx | 8080:31519/TCP
```

Now, you can interact with the application on:

> 169.xx.xx.xx:31519/m1/

_Note: Since there's only 1 node in the cluster, you will only ever use `/m1`. For more than 1 node, see the next section._

#### 3.5.2 Multi-node cluster

For a multi-node cluster, an ingress resource would be highly beneficial.

The **first** step would be to create the kubernetes service in each individual node, using the example [here](https://github.com/prashantgupta24/firewalld-rest/blob/master/k8s/svc.yaml). The important thing to note is that we manually add the `Endpoints` resource for the service, which points to our node's private IP address and port `8080`.

The **second** step is the [ingress](https://github.com/prashantgupta24/firewalld-rest/blob/master/k8s/ingress.yaml) resource. It redirects different routes to different nodes in the cluster. For example, in the ingress file above,

A request to `/m1` will be redirected to the `first` node, a request to `/m2` will be redirected to the `second` node, and so on. This will let you control each node's individual SSH access through a single endpoint.

### 3.6 Configure linux systemd service

See [this](https://github.com/prashantgupta24/firewalld-rest/blob/master/firewalld-rest.service) for an example of a linux systemd service.

The `.service` file should be placed under `etc/systemd/system` directory.

**Note**: This service assumes your binary is at `/root/rest/firewalld-rest`. You can change that in the file above.

### 3.7 Start and enable systemd service.

**Start**

```
systemctl start firewalld-rest
```

**Logs**

You can see the logs for the service using:

```
journalctl -r
```

**Enable**

```
systemctl enable firewalld-rest
```

### 3.8 IP JSON

This is how the IP JSON looks like, so that you know how you have to pass your IP and domain to the application:

```
type IP struct {
	IP     string `json:"ip"`
	Domain string `json:"domain"`
}
```

### 3.9 Interacting with the REST application

#### 3.9.1 Index page

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
--header 'Authorization: Bearer <jwt>'
```

#### 3.9.2 Show all IPs

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
--header 'Authorization: Bearer <jwt>'
```

#### 3.9.3 Add new IP

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
--header 'Authorization: Bearer <jwt>' \
--header 'Content-Type: application/json' \
--data-raw '{"ip":"10.xx.xx.xx","domain":"example.com"}'
```

#### 3.9.4 Show if IP is present

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
--header 'Authorization: Bearer <jwt>'
```

#### 3.9.5 Delete IP

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
--header 'Authorization: Bearer <jwt>'
```

## 4. Helpful tips/links

- ### 4.1 Creating custom kubernetes endpoint

  - https://theithollow.com/2019/02/04/kubernetes-endpoints/

- ### 4.2 Firewalld rules

  - https://www.digitalocean.com/community/tutorials/how-to-set-up-a-firewall-using-firewalld-on-centos-7

  #### 4.2.1 Useful commands

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

  #### 4.2.2 Rich rules

  `firewall-cmd --permanent --zone=public --list-rich-rules`

  `firewall-cmd --permanent --zone=public --add-rich-rule='rule family="ipv4" source address="10.10.99.10/32" port protocol="tcp" port="22" accept'`

  `firewall-cmd --permanent --zone=public --add-rich-rule='rule family="ipv4" source address="192.168.100.0/24" invert="True" drop'`

  > Reject will reply back with an ICMP packet noting the rejection, while a drop will just silently drop the traffic and do nothing else, so a drop may be preferable in terms of security as a reject response confirms the existence of the system as it is rejecting the request.

  #### 4.2.3 Misc tips

  > --add-source=IP can be used to add an IP address or range of addresses to a zone. This will mean that if any source traffic enters the systems that matches this, the zone that we have set will be applied to that traffic. In this case we set the ‘testing’ zone to be associated with traffic from the 10.10.10.0/24 range.

  ```
  [root@centos7 ~]# firewall-cmd --permanent --zone=testing --add-source=10.10.10.0/24
  success
  ```

- ### 4.3 Using JWT in Go

  - https://www.thepolyglotdeveloper.com/2017/03/authenticate-a-golang-api-with-json-web-tokens/

- ### 4.4 Using golang Exec

  - https://stackoverflow.com/questions/39151420/golang-executing-command-with-spaces-in-one-of-the-parts

- ### 4.5 Systemd services

  - https://medium.com/@benmorel/creating-a-linux-service-with-systemd-611b5c8b91d6
  - https://www.digitalocean.com/community/tutorials/understanding-systemd-units-and-unit-files
  - [Logs using journalctl](https://www.linode.com/docs/quick-answers/linux/how-to-use-journalctl/)

- ### 4.6 Using LDFlags in golang

  - https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications

## 5. Commands for generating public/private key

```
openssl genrsa -key private-key-sc.pem
openssl req -new -x509 -key private-key-sc.pem -out public.cert
```

[version-badge]: https://img.shields.io/github/v/release/prashantgupta24/firewalld-rest
[releases]: https://github.com/prashantgupta24/firewalld-rest/releases
