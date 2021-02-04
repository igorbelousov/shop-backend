package main

// openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
// openssl rsa -pubout -in private.pem -out public.pem

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/igorbelousov/shop-backend/foundation/database"
	"github.com/igorbelousov/shop-backend/internal/data/schema"
)

func main() {
	// tokengen()
	// genkey()
	migrate()
}

func tokengen() {
	privatePEM, err := ioutil.ReadFile("private.pem")
	if err != nil {
		log.Fatal(err, 1)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		log.Fatal(err, 2)
	}

	claims := struct {
		jwt.StandardClaims
		Roles []string
	}{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "Shop backend ",
			Subject:   "123456789",
			ExpiresAt: time.Now().Add(8760 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Roles: []string{"USER"},
	}

	method := jwt.GetSigningMethod("RS256")
	tkn := jwt.NewWithClaims(method, claims)
	tkn.Header["kid"] = "54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"
	str, err := tkn.SignedString(privateKey)
	if err != nil {
		log.Fatal(err, 3)
	}

	fmt.Printf("-----BEGIN TOKEN-----\n%s\n-----END TOKEN-----\n", str)

}

func genkey() {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	privateFile, err := os.Create("private.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer privateFile.Close()

	privateBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	if err := pem.Encode(privateFile, &privateBlock); err != nil {
		log.Fatal(err)
	}

	//--------------------------------------------------------------------

	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Fatal(err)
	}

	publicFile, err := os.Create("public.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer publicFile.Close()

	publicBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	if err := pem.Encode(publicFile, &publicBlock); err != nil {
		log.Fatal(err)
	}

	fmt.Println("KEY GEN DONE")
}

func migrate() {

	dbConfig := database.Config{
		User:       "postgres",
		Password:   "postgres",
		Host:       "0.0.0.0",
		Name:       "postgres",
		DisableTLS: true,
	}

	db, err := database.Open(dbConfig)
	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	if err := schema.Migrate(db); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("migrate comlete")

	if err := schema.Seed(db); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("seed data comlete")
}
