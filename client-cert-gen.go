package main

import (
	"crypto/ecdsa"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"runtime"
	"strings"
	"time"
)

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

func main() {
	var mac = flag.String("mac", "", "MAC address of DUT")
	var certValidDuration = flag.Int("dur", 3650, "Unit: days. The duration of SSL certificate validity")
	flag.Parse()
	flagset := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { flagset[f.Name] = true })

	if !flagset["mac"] {
		fmt.Println("ERR: --mac not set")
		os.Exit(2)
	}

	fmt.Println("mac =", *mac)
	var macPath string
	macPath = strings.Replace(*mac, ":", "-", -1)
	if runtime.GOOS == "windows" {
		macPath = macPath + "\\"
	} else {
		macPath = macPath + "/"
	}
	fmt.Printf(time.Now().Format("[2006-01-02 15:04:05] ")+"Creating %v folder ...\n", macPath)
	os.MkdirAll(macPath, 0755)

	// load CA key pair
	//      public key

	caPublicKeyFile, err := ioutil.ReadFile("rootCA.crt")
	if err != nil {
		panic(err)
	}
	pemBlock, _ := pem.Decode(caPublicKeyFile)
	if pemBlock == nil {
		panic("pem.Decode failed")
	}
	// fmt.Println("pemBlock: ", pemBlock)
	caCRT, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		panic(err)
	}
	// fmt.Println("caCRT: ", caCRT)
	//      private key
	caPrivateKeyFile, err := ioutil.ReadFile("rootCA.key")
	if err != nil {
		panic(err)
	}
	pemBlock, _ = pem.Decode(caPrivateKeyFile)
	if pemBlock == nil {
		panic("pem.Decode failed")
	}
	// fmt.Println("pemBlock: ", pemBlock)
	// fmt.Println("isEncryptedPEMBlock: ", x509.IsEncryptedPEMBlock(pemBlock))
	// der, err := x509.DecryptPEMBlock(pemBlock, []byte("password"))
	// if err != nil {
	// 	panic(err)
	// }
	caPrivateKey, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if err != nil {
		panic(err)
	}

	// load client certificate request
	// clientCSRFile, err := ioutil.ReadFile("client.csr")
	// if err != nil {
	// 	panic(err)
	// }
	// pemBlock, _ = pem.Decode(clientCSRFile)
	// if pemBlock == nil {
	// 	panic("pem.Decode failed")
	// }
	// clientCSR, err := x509.ParseCertificateRequest(pemBlock.Bytes)
	// if err != nil {
	// 	panic(err)
	// }
	// if err = clientCSR.CheckSignature(); err != nil {
	// 	panic(err)
	// }
	// generate client private key
	priv, err := rsa.GenerateKey(rand.Reader, 2048)

	// create client certificate template

	clientCRTTemplate := x509.Certificate{
		// Signature:          clientCSR.Signature,
		// SignatureAlgorithm: clientCSR.SignatureAlgorithm,

		// PublicKeyAlgorithm: clientCSR.PublicKeyAlgorithm,
		// PublicKey:          clientCSR.PublicKey,

		SerialNumber: big.NewInt(2),
		// Issuer:       caCRT.Subject,
		// Subject:      clientCSR.Subject,
		Subject: pkix.Name{
			CommonName:   "*.engeniuscloud.com",
			Organization: []string{"EnGenius Technologies, Inc."},
			SerialNumber: *mac,
			ExtraNames: []pkix.AttributeTypeAndValue{
				{
					Type:  []int{0, 9, 2342, 19200300, 100, 1, 1},
					Value: *mac,
				},
			},
		},
		NotBefore: time.Now(),
		// NotAfter:    time.Now().Add(24 * time.Hour * 365 * 10),
		NotAfter:    time.Now().Add(24 * time.Hour * time.Duration(*certValidDuration)),
		KeyUsage:    x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	// create client certificate from template and CA public key
	// clientCRTRaw, err := x509.CreateCertificate(rand.Reader, &clientCRTTemplate, caCRT, clientCSR.PublicKey, caPrivateKey)
	clientCRTRaw, err := x509.CreateCertificate(rand.Reader, &clientCRTTemplate, caCRT, publicKey(priv), caPrivateKey)
	if err != nil {
		panic(err)
	}
	// save the certificate
	clientCRTFilePath := macPath + "client.crt"
	clientCRTMd5Path := clientCRTFilePath + ".md5"
	fmt.Printf(time.Now().Format("[2006-01-02 15:04:05] ")+"Creating %v ...\n", clientCRTFilePath)
	clientCRTFile, err := os.Create(clientCRTFilePath)
	if err != nil {
		panic(err)
	}
	pem.Encode(clientCRTFile, &pem.Block{Type: "CERTIFICATE", Bytes: clientCRTRaw})
	clientCRTFile.Close()
	fmt.Printf(time.Now().Format("[2006-01-02 15:04:05] ")+"Creating %v ...\n", clientCRTMd5Path)
	clientCrt, err := ioutil.ReadFile(clientCRTFilePath)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(clientCRTMd5Path, []byte(fmt.Sprintf("%x", md5.Sum(clientCrt))), 0644)
	if err != nil {
		panic(err)
	}

	clientPrivKeyFilePath := macPath + "client.key"
	clientPrivKeyMd5Path := clientPrivKeyFilePath + ".md5"
	fmt.Printf(time.Now().Format("[2006-01-02 15:04:05] ")+"Creating %v ...\n", clientPrivKeyFilePath)
	clientPrivKeyFile, err := os.Create(clientPrivKeyFilePath)
	if err != nil {
		panic(err)
	}
	pem.Encode(clientPrivKeyFile, pemBlockForKey(priv))
	clientPrivKeyFile.Close()
	fmt.Printf(time.Now().Format("[2006-01-02 15:04:05] ")+"Creating %v ...\n", clientPrivKeyMd5Path)
	clientPrivKey, err := ioutil.ReadFile(clientPrivKeyFilePath)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(clientPrivKeyMd5Path, []byte(fmt.Sprintf("%x", md5.Sum(clientPrivKey))), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Printf(time.Now().Format("[2006-01-02 15:04:05] ") + "Done.\n")

}
