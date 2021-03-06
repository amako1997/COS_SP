package peeronionprotocol

import (
	onionrepository "onionRouting/go-torPeer/repositories/onion"
	cryptoserviceinterface "onionRouting/go-torPeer/services/crypto/crypto-service-interface"
	"onionRouting/go-torPeer/types"
	"os"

	logger "github.com/apsdehal/go-logger"
	"github.com/pkg/errors"
)

type PeerOnionService struct {
	onionRepo     onionrepository.OnionRepository
	cryptoService cryptoserviceinterface.CryptoService
	log           *logger.Logger
}

func NewOnionService(onionRepo onionrepository.OnionRepository, cryptoService cryptoserviceinterface.CryptoService) PeerOnionService {
	log, _ := logger.New("PeerOnionService", 1, os.Stdout)
	return PeerOnionService{
		onionRepo:     onionRepo,
		cryptoService: cryptoService,
		log:           log,
	}
}
func (this *PeerOnionService) SaveCircuit(circuit types.P2PBuildCircuitRequest) error {

	linkParamaeters := types.CircuitLinkParameters{
		Next:     circuit.Next,
		Previous: circuit.Previous,
	}
	this.log.Noticef("saved circuit under key %v", string(circuit.ID))
	if err := this.onionRepo.SaveCircuitLink(circuit.ID, linkParamaeters); err != nil {
		return err
	}

	return nil
}
func (this *PeerOnionService) GetSavedCircuit(cId []byte) (CircuitLinkGetDTO, error) {

	link, err := this.onionRepo.GetCircuitLinkParamaters(cId, this.log)
	if err != nil {
		return CircuitLinkGetDTO{}, err
	}
	savedLink := CircuitLinkGetDTO{
		ID:           cId,
		Next:         link.Next,
		Previous:     link.Previous,
		SharedSecret: link.SharedSecret,
	}
	return savedLink, nil
}
func (this *PeerOnionService) PeelOnionLayer(CircuitPayload types.CircuitPayload) ([]byte, string, error) {

	link, err := this.onionRepo.GetCircuitLinkParamaters(CircuitPayload.ID, this.log)
	if err != nil {
		return nil, "", err
	}
	peeledData, err := this.cryptoService.Decrypt(CircuitPayload.Payload, link.SharedSecret)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed  to peel encryption leayer")
	}
	this.log.Noticef("PeelOnionLayer decrypted layer is %v \n", string(peeledData))

	return peeledData, link.Next, nil

}
func (this *PeerOnionService) DecryptData(data []byte, key []byte) ([]byte, error) {
	decrypted, err := this.cryptoService.Decrypt(data, key)
	if err != nil {
		return nil, errors.Wrap(err, "failed during DecryptData")
	}
	return decrypted, nil
}
func (this *PeerOnionService) Forward(data []byte, circuitID []byte, nxt string, forwardType string, sendingCircuit []byte) (bool, []byte, error) {

	if nxt == "" {
		return false, nil, nil
	}
	next := "http://" + nxt
	this.log.Debugf("Forward dialing next %v \n", next)
	body, err := this.onionRepo.DialNext(circuitID, next, data, this.log, forwardType, sendingCircuit)
	if err != nil {
		return false, nil, err
	}
	return true, body, nil
}
func (this *PeerOnionService) BackTrack(chainID []byte) ([]byte, types.CircuitLinkParameters, error) {

	link, err := this.onionRepo.GetCircuitLinkParamaters(chainID, this.log)
	if err != nil {
		return nil, types.CircuitLinkParameters{}, errors.Wrap(err, "failed to backtrack, could not get link paramaters ")
	}
	return chainID, link, nil

}
func (this *PeerOnionService) AddOnionLayer(data []byte, link types.CircuitLinkParameters) ([]byte, string, error) {

	encrypted, err := this.cryptoService.Encrypt(data, link.SharedSecret)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to encrypt data while adding onion layer ")
	}
	this.log.Noticef("AddOnionLayer encrypted layer is %v \n", string(encrypted))

	return encrypted, link.Previous, nil
}
