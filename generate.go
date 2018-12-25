package main

import (
	rand2 "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand2.Int(rand2.Reader, max)
	subject := pkix.Name{
		Organization: []string{"Monstar Lab Inc."},
		OrganizationalUnit:[]string{"Software Agency"},
		CommonName:"Digital services agency",
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: subject,
		NotBefore: time.Now(),
		NotAfter: time.Now().Add(365*24*time.Hour),
		KeyUsage:x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
	}

	pk, _ := rsa.GenerateKey(rand2.Reader, 2048)

	derBytes, _ := x509.CreateCertificate(rand2.Reader, &template,
		&template, &pk.PublicKey, pk)
	certOut, _ := os.Create("cert.pem")
	pem.Encode(certOut, &pem.Block{Type:"CERTIFICATE", Bytes:derBytes})
	certOut.Close()

	keyOut, _ := os.Create("key.pem")
	pem.Encode(keyOut, &pem.Block{Type:"RSA PRIVATE KEY", Bytes:
		x509.MarshalPKCS1PrivateKey(pk)})
	keyOut.Close()

	http.HandleFunc("/index", handler)
	server := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: nil,
	}
	server.ListenAndServeTLS("cert.pem", "key.pem")

}


func handler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello World, %s!", request.URL.Path[1:])
}
