# trove
 Backend Practice in Go

## Dependencies
* Mostly POSIX compliant OS (I use Windows with Debian on WSL)
* Go
* Docker
* Sass

## Building
* Run `make install` to install quicktemplate (assuming your "go/bin" is already in PATH, if not, do that) and generate localhost certificate
* Install `localhost.crt` to your list of locally trusted roots
* Boot up Kafka image with `make boot`
* Run `make` to build and run server

### Note
Some browsers may give a security error since the certificate is self-signed.
Using properly signed certificates when deploying will solve this.
