package ssl

import (
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
)

type CaIface interface {
	GetCred() credentials.TransportCredentials
}

type serverTwoWayVerify struct {
	ServerPemPath string
	ServerKeyPath string
	CaPemPath     string
}

func NewServerTwoWayVerify(serverPemPath, serverKeyPath, caPemPath string) CaIface {
	return &serverTwoWayVerify{
		ServerPemPath: serverPemPath,
		ServerKeyPath: serverKeyPath,
		CaPemPath:     caPemPath,
	}
}

func (c *serverTwoWayVerify) GetCred() credentials.TransportCredentials {
	cert, err := tls.LoadX509KeyPair(c.ServerPemPath, c.ServerKeyPath)
	if err != nil {
		log.Fatal(err)
	}
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(c.CaPemPath)
	if err != nil {
		log.Fatal(err)
	}
	certPool.AppendCertsFromPEM(ca)
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})
	return creds
}

type serverOneWayVerify struct {
	CrtPath       string
	ServerKeyPath string
}

func NewServerOneWayVerify(crtPath, serverKeyPath string) CaIface {
	return &serverOneWayVerify{
		CrtPath:       crtPath,
		ServerKeyPath: serverKeyPath,
	}
}

func (c *serverOneWayVerify) GetCred() credentials.TransportCredentials {
	cred, err := credentials.NewServerTLSFromFile(c.CrtPath, c.ServerKeyPath)
	if err != nil {
		log.Fatal(err)
	}
	return cred
}
