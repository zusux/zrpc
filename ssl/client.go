package ssl

import (
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
)

type ClientTwoWayVerify struct {
	ClientPemPath string
	ClientKeyPath string
	CaPemPath string
	ServerName string
}

func NewClientTwoWayVerify(clientPemPath string, clientKeyPath string, caPemPath string, serverName string) CaIface {
	return &ClientTwoWayVerify{ClientPemPath: clientPemPath, ClientKeyPath: clientKeyPath, CaPemPath: caPemPath, ServerName: serverName}
}
func (c *ClientTwoWayVerify) GetCred()(credentials.TransportCredentials){
	cert,err := tls.LoadX509KeyPair(c.ClientPemPath,c.ClientKeyPath)
	if err != nil{
		log.Fatal(err)
	}
	certPool := x509.NewCertPool()
	ca , err := ioutil.ReadFile(c.CaPemPath)
	if err != nil{
		log.Fatal(err)
	}
	certPool.AppendCertsFromPEM(ca)
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName: c.ServerName,
		RootCAs: certPool,
	})
	return creds
}

type clientOneWayVerify struct {
	ServerName string
	CrtPath string
}
func (c *clientOneWayVerify) GetCred()(credentials.TransportCredentials){
	cert,err := credentials.NewClientTLSFromFile(c.CrtPath,c.ServerName)
	if err != nil{
		log.Fatal(cert)
	}
	return cert
}
func NewClientOneWayVerify(crtPath,serverName string) CaIface{
	return &clientOneWayVerify{
		CrtPath:crtPath,
		ServerName:serverName,
	}
}
