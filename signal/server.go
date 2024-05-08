package signal

import (
	"errors"
	"fmt"
	"sync"
)

type Server struct {
	users map[string]*UserData
	mutex sync.RWMutex
}

type UserData struct {
	IdentityKey  []byte
	PreKeys      []*PreKey
	SignedPreKey *PreKey
}

func NewServer() *Server {
	return &Server{
		users: make(map[string]*UserData),
	}
}

// RegisterUser adds a new user to the server along with their identity key and prekeys.
func (s *Server) RegisterUser(userID string, identityKey []byte, preKeys []*PreKey, signedPreKey *PreKey) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.users[userID]; exists {
		return errors.New("user already registered")
	}

	s.users[userID] = &UserData{
		IdentityKey:  identityKey,
		PreKeys:      preKeys,
		SignedPreKey: signedPreKey,
	}
	return nil
}

// FetchPreKeys returns a prekey for a specified user to start a new session.
func (s *Server) FetchPreKeys(userID string) (*PreKey, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	userData, exists := s.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}

	if len(userData.PreKeys) == 0 {
		return nil, errors.New("no available prekeys")
	}

	// Return the first prekey and remove it from the list
	preKey := userData.PreKeys[0]
	userData.PreKeys = userData.PreKeys[1:]
	return preKey, nil
}

// FetchSignedPreKey returns a signed prekey for a specified user.
func (s *Server) FetchSignedPreKey(userID string) (*PreKey, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	userData, exists := s.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}

	if userData.SignedPreKey == nil {
		return nil, errors.New("no signed prekey available")
	}

	return userData.SignedPreKey, nil
}

// ForwardMessage forwards an encrypted message from the sender to the recipient.
func (s *Server) ForwardMessage(senderID, recipientID, message string) error {
	s.mutex.RLock()
	_, exists := s.users[recipientID]
	s.mutex.RUnlock()
	if !exists {
		return errors.New("recipient not found")
	}

	// In a real application, this would place the message into a queue or similar structure
	// for delivery to the recipient. Here we'll just simulate that process.
	fmt.Printf("Message from %s to %s: %s\n", senderID, recipientID, message)
	return nil
}
