package cert

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

const (
	exportFolderPath   string = "./cert"
	certFileName       string = "cert.crt"
	publicKeyFileName  string = "public_key.pem"
	privateKeyFileName string = "private_key.pem"
)

func init() {
	// создаём шаблон сертификата
	cert := &x509.Certificate{
		// указываем уникальный номер сертификата
		SerialNumber: big.NewInt(1658),
		// заполняем базовую информацию о владельце сертификата
		Subject: pkix.Name{
			Organization: []string{"Yandex.Praktikum"},
			Country:      []string{"RU"},
		},
		// разрешаем использование сертификата для 127.0.0.1 и ::1
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		// сертификат верен, начиная со времени создания
		NotBefore: time.Now(),
		// время жизни сертификата — 10 лет
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		// устанавливаем использование ключа для цифровой подписи,
		// а также клиентской и серверной авторизации
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	// создаём новый приватный RSA-ключ длиной 4096 бит
	// обратите внимание, что для генерации ключа и сертификата
	// используется rand.Reader в качестве источника случайных данных
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	check(err)

	publicKey := &privateKey.PublicKey
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	check(err)

	// создаём сертификат x.509
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, publicKey, privateKey)
	check(err)

	// кодируем сертификат и ключ в формате PEM, который
	// используется для хранения и обмена криптографическими ключами
	var certPEM bytes.Buffer
	pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	var privateKeyPEM bytes.Buffer
	pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	var publicKeyPEM bytes.Buffer
	pem.Encode(&publicKeyPEM, &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	// сохраняем сертификат, приватный и публичный ключи
	certFile, err := os.OpenFile(fmt.Sprintf("%s/%s", exportFolderPath, certFileName), os.O_RDWR|os.O_CREATE, 0644)
	check(err)
	_, err = certFile.Write(certBytes)
	check(err)
	check(certFile.Close())

	pubKeyFile, err := os.OpenFile(fmt.Sprintf("%s/%s", exportFolderPath, publicKeyFileName), os.O_RDWR|os.O_CREATE, 0644)
	check(err)
	_, err = pubKeyFile.Write(publicKeyPEM.Bytes())
	check(err)
	check(pubKeyFile.Close())

	privKeyFile, err := os.OpenFile(fmt.Sprintf("%s/%s", exportFolderPath, privateKeyFileName), os.O_RDWR|os.O_CREATE, 0644)
	check(err)
	_, err = privKeyFile.Write(privateKeyPEM.Bytes())
	check(err)
	check(privKeyFile.Close())
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
