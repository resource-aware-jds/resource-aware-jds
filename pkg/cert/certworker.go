package cert

type WorkerNodeCACertificate TLSCertificate

type WorkerNodeCACertificateConfig struct {
	CACertificateFilePath string
}

func ProvideWorkerNodeCACertificate(config WorkerNodeCACertificateConfig) (WorkerNodeCACertificate, error) {
	certificateChain, err := LoadCertificatesFromFile(config.CACertificateFilePath)
	if err != nil {
		return nil, err
	}

	return ProvideTLSCertificate(certificateChain, nil, true)
}

type WorkerNodeTransportCertificate struct {
	CertificateFileLocation string
	PrivateKeyFileLocation  string
}

func ProvideWorkerNodeTransportCertificate(config WorkerNodeTransportCertificate) (TransportCertificate, error) {
	privateKeyData, err := LoadKeyFromFile(config.PrivateKeyFileLocation)
	if err != nil {
		// TODO: Call register worker node with CP
		return nil, err
	}

	certificateChain, err := LoadCertificatesFromFile(config.CertificateFileLocation)
	if err != nil {
		return nil, err
	}

	return ProvideTLSCertificate(certificateChain, privateKeyData, false)
}
