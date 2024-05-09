package main

import (
	"crypto/rand"
	"golang.org/x/crypto/chacha20poly1305"
	"io"
	"strings"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/ed25519"
)

// KeyStore to hold generated public and private keys for demonstration
type KeyStore struct {
	Curve25519PrivateKey []byte
	Curve25519PublicKey  []byte
	Ed25519PublicKey     []byte
	Ed25519PrivateKey    ed25519.PrivateKey
}

var Keys *KeyStore

func init() {
	Keys = GenerateKeys()
}

// GenerateKeys generates a Curve25519 key pair and an Ed25519 signing key pair.
// func GenerateKeys() *KeyStore {
// 	edPub, edPriv, _ := ed25519.GenerateKey(rand.Reader) // Ed25519 keys for signing

// 	// Generating Curve25519 keys for encryption
// 	var curvePubKey, curvePrivKey [32]byte
// 	copy(curvePrivKey[:], edPriv[:32]) // Curve25519 uses the first 32 bytes of Ed25519 private key
// 	curve25519.ScalarBaseMult(&curvePubKey, &curvePrivKey)

// 	keys := &KeyStore{
// 		Curve25519PrivateKey: curvePrivKey[:],
// 		Curve25519PublicKey:  curvePubKey[:],
// 		Ed25519PublicKey:     edPub,
// 		Ed25519PrivateKey:    edPriv,
// 	}
// 	return keys
// }

// GeneratePreKeys generates one-time prekeys for ephemeral sessions.
func GeneratePreKeys() (preKeyPublic, preKeyPrivate []byte) {
	keys := GenerateKeys()
	return keys.Curve25519PublicKey, keys.Curve25519PrivateKey
}

// // SessionStart generates a session key using a shared secret derived from the private key of one user and the public key of another.
// func SessionStart(ourPrivateKey, theirPublicKey []byte) []byte {
// 	var ourPrivKey, theirPubKey, sharedSecret [32]byte
// 	copy(ourPrivKey[:], ourPrivateKey[:32])
// 	copy(theirPubKey[:], theirPublicKey[:32])
// 	curve25519.ScalarMult(&sharedSecret, &ourPrivKey, &theirPubKey)
// 	return sharedSecret[:]
// }

func SessionStart(ourPrivateKey, theirPublicKey []byte) ([]byte, error) {
	var sharedSecret, ourPrivKey, theirPubKey [32]byte
	copy(ourPrivKey[:], ourPrivateKey[:])
	copy(theirPubKey[:], theirPublicKey[:])
	curve25519.ScalarMult(&sharedSecret, &ourPrivKey, &theirPubKey)
	return sharedSecret[:], nil
}

// GenerateCurve25519KeyPair generates a new public/private key pair for Curve25519.
// func GenerateCurve25519KeyPair() ([]byte, []byte, error) {
// 	edPub, edPriv, _ := ed25519.GenerateKey(rand.Reader)

//     var privateKey [32]byte
//     _, err := rand.Read(privateKey[:])
//     if err != nil {
//         return nil, nil, err
//     }

//     var publicKey [32]byte
//     curve25519.ScalarBaseMult(&publicKey, &privateKey)

//	    return publicKey[:], privateKey[:], nil
//	}
//
// GenerateCurve25519KeyPair generates a new public/private key pair for Curve25519.
func GenerateKeys() *KeyStore {
	edPub, _, _ := ed25519.GenerateKey(rand.Reader)

	var privateKey [32]byte
	_, err := rand.Read(privateKey[:])
	if err != nil {
		return nil
	}

	var publicKey [32]byte
	curve25519.ScalarBaseMult(&publicKey, &privateKey)

	keys := &KeyStore{
		Curve25519PrivateKey: privateKey[:],
		Curve25519PublicKey:  publicKey[:],
		Ed25519PublicKey:     edPub,
	}
	return keys
}

// DeriveSharedSecret computes a shared secret using the private key of the caller and the public key of the peer.
// func SessionStart(privateKey, publicKey []byte) ([]byte, error) {
//     var privKey [32]byte
//     var pubKey [32]byte
//     var sharedSecret [32]byte

//     copy(privKey[:], privateKey[:32])
//     copy(pubKey[:], publicKey[:32])

//	    curve25519.ScalarMult(&sharedSecret, &privKey, &pubKey)
//	    return sharedSecret[:], nil
//	}
func checkKey(key []byte) {}

// // EncryptMessage encrypts plaintext using ChaCha20-Poly1305.
// func EncryptMessage(plaintext, sessionKey []byte) ([]byte, error) {
// 	aead, err := chacha20poly1305.New(sessionKey)
// 	if err != nil {
// 		return nil, err
// 	}
// 	nonce := make([]byte, aead.NonceSize())
// 	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
// 		return nil, err
// 	}
// 	return aead.Seal(nonce, nonce, plaintext, nil), nil
// }

// // DecryptMessage decrypts ciphertext using ChaCha20-Poly1305.
// func DecryptMessage(ciphertext, sessionKey []byte) ([]byte, error) {
// 	aead, err := chacha20poly1305.New(sessionKey)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(ciphertext) < aead.NonceSize() {
// 		return nil, err
// 	}
// 	nonce, encryptedMessage := ciphertext[:aead.NonceSize()], ciphertext[aead.NonceSize():]
// 	return aead.Open(nil, nonce, encryptedMessage, nil)
// }

var sessionStore = map[string][]byte{}

func getSessionKey(user string) []byte {
	public := PublicKey[user]
	private := Keys.Curve25519PrivateKey

	session := sessionStore[user]
	if session == nil {
		session, _ = SessionStart(private, public)
		sessionStore[user] = session

	}
	return session
}

func Encrypt(plaintext, user string) string {
	Plaintext, _ := encryptMessage([]byte(plaintext), getSessionKey(user))
	plaintext = string(Plaintext)
	return plaintext
}

func Decrypt(ciphertext, user string) string {
	Ciphertext, _ := decryptMessage([]byte(ciphertext), getSessionKey(user))
	ciphertext = string(Ciphertext)
	return ciphertext
}

// EncryptMessage encrypts plaintext using ChaCha20-Poly1305.
func encryptMessage(plaintext, sessionKey []byte) ([]byte, error) {
	//return plaintext,nil
	aead, err := chacha20poly1305.New(sessionKey)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	Plaintext := aead.Seal(nonce, nonce, plaintext, nil)
	checkKey(Plaintext)
	return plaintext, nil
}

// DecryptMessage decrypts ciphertext using ChaCha20-Poly1305.
func decryptMessage(ciphertext, sessionKey []byte) ([]byte, error) {
	//dec := ciphertext
	aead, err := chacha20poly1305.New(sessionKey)
	if err != nil {
		return ciphertext, nil
	}
	if len(ciphertext) < aead.NonceSize() {
		return ciphertext, nil
	}
	nonce, encryptedMessage := ciphertext[:aead.NonceSize()], ciphertext[aead.NonceSize():]
	Ciphertext, _ := aead.Open(nil, nonce, encryptedMessage, nil)
	checkKey(Ciphertext)
	return ciphertext, nil
}

type ServerStore struct {
	PublicKey []byte
	SignedPub []byte
}

var Store = map[string]ServerStore{}

// SendKeys to a server or peer
func SendKeysContent() string {
	return string(Keys.Curve25519PublicKey) + " " + string(Keys.Ed25519PublicKey)
}

// StoreKeys in server storage
func StoreKeys(user, content string) {
	var temp ServerStore
	// Placeholder for server storage logic
	keys := strings.Fields(content)

	temp.PublicKey = []byte(keys[0])
	temp.SignedPub = []byte(keys[1])

	Store[user] = temp
}

func SendStore() string {
	var s string
	for i, v := range Store {
		s += string(i) + " " + string(v.PublicKey) + " "
	}
	return s
}

var PublicKey = map[string][]byte{}

func getKeys(msg Message) {
	values := strings.Fields(msg.Content)
	for i, user := range values {
		if i%2 != 0 {
			continue
		}
		if PublicKey[user] == nil {
			PublicKey[user] = []byte(values[i+1])
		}

	}
}

// func connectAndExchangeKeys() {
// 	conn, _, err := websocket.DefaultDialer.Dial("ws://server_address", nil)
// 	if err != nil {
// 		log.Fatal("dial:", err)
// 	}
// 	defer conn.Close()

// 	clientKeys := GenerateKeys()

// 	// Send client's public keys
// 	publicKeyMessage := append(clientKeys.Curve25519PublicKey, clientKeys.Ed25519PublicKey...)
// 	if err := conn.WriteMessage(websocket.BinaryMessage, publicKeyMessage); err != nil {
// 		log.Fatal("write:", err)
// 	}

// 	// Read server's public keys
// 	_, serverPublicKeyMessage, err := conn.ReadMessage()
// 	if err != nil {
// 		log.Fatal("read:", err)
// 	}
// 	serverCurvePubKey := serverPublicKeyMessage[:32]
// 	serverEdPubKey := serverPublicKeyMessage[32:]

// 	// Verify server's signature or perform any necessary authentication steps

// 	// Generate session key
// 	sessionKey := SessionStart(clientKeys.Curve25519PrivateKey, serverCurvePubKey)

// 	// Use sessionKey for further encrypted communication
// }

// func GetPublicKey(user string){

// }

// func serverConection(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println("upgrade:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	// Server generates its keys
// 	serverKeys := GenerateKeys()

// 	// Server sends its public keys
// 	publicKeyMessage := append(serverKeys.Curve25519PublicKey, serverKeys.Ed25519PublicKey...)
// 	if err := conn.WriteMessage(websocket.BinaryMessage, publicKeyMessage); err != nil {
// 		log.Println("write:", err)
// 		return
// 	}

// 	// Read client's public keys
// 	_, clientPublicKeyMessage, err := conn.ReadMessage()
// 	if err != nil {
// 		log.Println("read:", err)
// 		return
// 	}
// 	clientCurvePubKey := clientPublicKeyMessage[:32]
// 	clientEdPubKey := clientPublicKeyMessage[32:]

// 	// Verify signature or perform any necessary authentication steps

// 	// Generate session key
// 	sessionKey := SessionStart(serverKeys.Curve25519PrivateKey, clientCurvePubKey)

// 	// Use sessionKey for further encrypted communication
// }
