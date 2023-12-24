package cert

//func ProvideTLSCertificateOld(config Config) (TLSCertificate, error) {
//	certificate := tlsCertificate{}
//
//	// Try to load the Certificate from the file
//	err := certificate.loadCertificateFromFile(config.CertificateFileLocation, config.PrivateKeyFileLocation)
//	if err == nil {
//		logrus.Info("Loaded Certificate from file: ", config.CertificateFileLocation, ":", config.PrivateKeyFileLocation)
//		return &certificate, nil
//	}
//
//	logrus.Warn("Failed to load certificate from file with this error: ", err)
//
//	// Create the new certificate instead.
//	err = certificate.createCertificate(config.ValidDuration, config.CertificateSubject, config.ParentCertificate)
//	if err != nil {
//		logrus.Error("Failed to create new certificate with this error: ", err)
//		return nil, err
//	}
//
//	// Save the created certificate to file
//	err = certificate.SaveCertificateToFile(config.CertificateFileLocation, config.PrivateKeyFileLocation)
//	if err != nil {
//		logrus.Error("Failed to save the created certificate with this error", err)
//		return nil, err
//	}
//	return &certificate, nil
//}
