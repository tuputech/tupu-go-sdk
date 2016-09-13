package gencert

import (
    //"crypto"
    "crypto/rand"
    "crypto/rsa"
    //"crypto/sha256"
    "fmt"
    "os"
)

func Do() {

    jimenaPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)

    if err != nil {
        fmt.Println(err.Error)
        os.Exit(1)
    }

    jimenaPublicKey := &jimenaPrivateKey.PublicKey

    alistairPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)

    if err != nil {
        fmt.Println(err.Error)
        os.Exit(1)
    }

    alistairPublicKey := &alistairPrivateKey.PublicKey

    fmt.Println("Private Key : ", jimenaPrivateKey)
    fmt.Println("Public key ", jimenaPublicKey)
    fmt.Println("Private Key : ", alistairPrivateKey)
    fmt.Println("Public key ", alistairPublicKey)

}
