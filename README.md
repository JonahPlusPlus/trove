# trove
 Backend Practice in Go

## Dependencies
* Mostly POSIX compliant OS (I use Windows with Debian on WSL)
* Go
* Docker

## Building
* Run `make certificate` to generate certificate
* Install `localhost.crt` to your list of locally trusted roots
* Boot up Kafka image with `make boot`
* Run `make` to run server

### Note
Some browsers may give a security error since the certificate is self-signed.
Using properly signed certificates when deploying will solve this.
