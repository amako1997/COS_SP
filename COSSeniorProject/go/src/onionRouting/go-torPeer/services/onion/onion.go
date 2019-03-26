package peeronionprotocol

import (
	"fmt"
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
	fmt.Println(linkParamaeters.Next)
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
func (this *PeerOnionService) PeelOnionLayer(CircuitPayload types.CircuitPayload) error {

	link, err := this.onionRepo.GetCircuitLinkParamaters(CircuitPayload.ID, this.log)
	if err != nil {
		return err
	}
	peeledData, err := this.cryptoService.Decrypt(CircuitPayload.Payload, link.SharedSecret)
	if err != nil {
		return errors.Wrap(err, "failed  to peel encryption leayer")
	}
	this.log.Noticef("decrypted layer is %v \n", string(peeledData))
	next := "http://" + link.Next
	this.log.Debugf("dialing next %v \n", next)
	err = this.onionRepo.DialNext(CircuitPayload.ID, next, peeledData, this.log)
	if err != nil {
		return err
	}
	return nil

}