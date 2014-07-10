package henchman

import (
	"code.google.com/p/go.crypto/ssh"
	"io/ioutil"
)

const (
	ECHO          = 53
	TTY_OP_ISPEED = 128
	TTY_OP_OSPEED = 129
)

func loadPEM(file string) (ssh.Signer, error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	key, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func ClientKeyAuth(keyFile string) (ssh.AuthMethod, error) {
	key, err := loadPEM(keyFile)
	return ssh.PublicKeys(key), err
}

func PasswordAuth(pass string) (ssh.AuthMethod, error) {
	return ssh.Password(pass), nil
}
