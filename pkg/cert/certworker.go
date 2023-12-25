package cert

type ClientCATLSCertificate TLSCertificate

type ClientCATLSCertificateConfig struct {
	CACertificateFilePath string
}

func ProvideClientCATLSCertificate(config ClientCATLSCertificateConfig) (TLSCertificate, error) {
	certificateChain, err := LoadCertificatesFromFile(config.CACertificateFilePath)
	if err != nil {
		return nil, err
	}

	return ProvideTLSCertificate(certificateChain, nil, true)
}
