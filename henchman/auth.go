package henchman

import (
	"code.google.com/p/go.crypto/ssh"
	"io"
	"io/ioutil"
	"strings"
)

func strip(v string) string {
	return strings.TrimSpace(strings.Trim(v, "\n"))
}

type password string

func (p password) Password(pass string) (string, error) {
	return string(p), nil
}

type keychain struct {
	keys []ssh.Signer
}

func (k *keychain) Key(i int) (ssh.PublicKey, error) {
	if i < 0 || i >= len(k.keys) {
		return nil, nil
	}
	return k.keys[i].PublicKey(), nil
}

func (k *keychain) Sign(i int, rand io.Reader, data []byte) (sig []byte, err error) {
	return k.keys[i].Sign(rand, data)
}

func (k *keychain) add(key ssh.Signer) {
	k.keys = append(k.keys, key)
}

func (k *keychain) loadPEM(file string) error {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	key, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return err
	}
	k.add(key)
	return nil
}

func ClientKeyAuth(keyFile string) (ssh.ClientAuth, error) {
	k := new(keychain)
	err := k.loadPEM(keyFile)
	return ssh.ClientAuthKeyring(k), err
}

func PasswordAuth(pass string) (ssh.ClientAuth, error) {
	return ssh.ClientAuthPassword(password(pass)), nil
}
