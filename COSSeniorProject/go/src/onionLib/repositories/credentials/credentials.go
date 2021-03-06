package peercredentialsrepository

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	storageserviceinterface "onionLib/services/storage/storage-interface"
	"onionLib/types"

	"github.com/dgraph-io/badger"

	"github.com/pkg/errors"
)

type PeerCredentials struct {
	db storageserviceinterface.StorageService
}

func NewPeerCredentialsRepository(db storageserviceinterface.StorageService) PeerCredentials {

	return PeerCredentials{
		db: db,
	}
}
func (this *PeerCredentials) GetPeerCredentials(peerID string, PublicKey []byte) (types.PeerCredentials, error) {

	fmt.Println("trying to get peer credentials via state key :" + peerID)
	credentialBytes, err := this.db.Get(peerID)
	if err != nil {
		return types.PeerCredentials{}, errors.Wrap(err, "failed get Peer credentials bytes from the database ")
	}
	peerCredentials := types.PeerCredentials{}
	if err := json.Unmarshal(credentialBytes, &peerCredentials); err != nil {

		return types.PeerCredentials{}, errors.Wrap(err, "failed to unmarshal credentialBytes")
	}
	return peerCredentials, nil
}
func (this *PeerCredentials) SavePeerCredentials(peerID string, credentials types.PeerCredentials) error {

	fmt.Println("peer credentials saved under key:" + peerID)

	savedCredentialBytes, err := this.db.Get(peerID)
	if err != nil && err != badger.ErrKeyNotFound {

		return errors.Wrap(err, "failed to get savedCredentialBytes ")
	}
	if savedCredentialBytes == nil {
		credentialBytes, err := json.Marshal(credentials)
		if err != nil {
			return errors.Wrap(err, "failed to marshal credentials")
		}
		if err := this.db.Put(peerID, credentialBytes); err != nil {
			return errors.Wrap(err, "failed to save peer credentials ")
		}
		return nil
	}
	if savedCredentialBytes != nil {

		peerCredentials := types.PeerCredentials{}
		if err := json.Unmarshal(savedCredentialBytes, &peerCredentials); err != nil {

			return errors.Wrap(err, "failed to unmarshal credentialBytes")
		}
		if credentials.PublicKey != nil {
			peerCredentials.PublicKey = credentials.PublicKey
		}
		peerCredentials.SharedSecret = credentials.SharedSecret
		newPeerCredentialBytes, err := json.Marshal(peerCredentials)
		if err != nil {
			return errors.Wrap(err, "failed to to marshal peerCredentials")
		}
		if err = this.db.Put(peerID, newPeerCredentialBytes); err != nil {
			return errors.Wrap(err, "failed to save newPeerCredentialBytes in the database ")
		}
	}

	return nil
}
func (this *PeerCredentials) createHash(data []byte) ([]byte, error) {

	hasher := sha256.New()
	_, err := hasher.Write(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make hash ")
	}
	hashedData := hasher.Sum(nil)
	return hashedData, nil
}
