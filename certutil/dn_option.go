package certutil

// DNOption can set the Distinguished Name for cert
// Distinguished name (DN) is a term that describes the identifying information
// in a certificate and is part of the certificate itself.
type DNOption func(*KeyCert)

func WithDNOptionCountry(country ...string) DNOption {
	return func(kc *KeyCert) {
		kc.cert.Subject.Country = country
	}
}

func WithDNOptionProvince(province ...string) DNOption {
	return func(kc *KeyCert) {
		kc.cert.Subject.Province = province
	}
}

func WithDNOptionLocality(locality ...string) DNOption {
	return func(kc *KeyCert) {
		kc.cert.Subject.Locality = locality
	}
}

func WithDNOptionStreetAddress(streetAddress ...string) DNOption {
	return func(kc *KeyCert) {
		kc.cert.Subject.StreetAddress = streetAddress
	}
}

func WithDNOptionPostalCode(postalCode ...string) DNOption {
	return func(kc *KeyCert) {
		kc.cert.Subject.PostalCode = postalCode
	}
}

func WithDNOptionOrganization(organization ...string) DNOption {
	return func(kc *KeyCert) {
		kc.cert.Subject.Organization = organization
	}
}

func WithDNOptionOrganizationalUnit(organizationalUnit ...string) DNOption {
	return func(kc *KeyCert) {
		kc.cert.Subject.OrganizationalUnit = organizationalUnit
	}
}
