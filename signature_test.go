package http_signature_auth

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"testing"

	. "github.com/onsi/gomega"
)

func TestVerifySignatureWithExistingMaterial(t *testing.T) {
	RegisterTestingT(t)

	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	Expect(err).ToNot(HaveOccurred())
	rsaKeyID := "testKeyIDrsa"
	keys := NewKeysDatabase()
	keys.AddKey(KeyID(rsaKeyID), &rsaKey.PublicKey)
	
	ecdsaKey, err:= ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	Expect(err).ToNot(HaveOccurred())
	ecdsaKeyID := "testKeyIDEcdsa"
	keys.AddKey(KeyID(ecdsaKeyID), &ecdsaKey.PublicKey)

	ed25519PubKey, ed25519PrivKey, err := ed25519.GenerateKey(rand.Reader)
	Expect(err).ToNot(HaveOccurred())
	ed25519KeyID := "testKeyIDed25519"
	keys.AddKey(KeyID(ed25519KeyID), ed25519PubKey)

	material := &TLSExporterMaterial{}
	copy(material.signatureInput[:], []byte("testSignatureInput"))
	copy(material.verification[:], []byte("testVerification"))


	// rsa
	signature, err := NewSignatureWithMaterial(material, KeyID(rsaKeyID), rsaKey, &rsaKey.PublicKey, tls.PSSWithSHA256)
	Expect(err).ToNot(HaveOccurred())
	ok, err := ValidateSignatureWithMaterial(keys, signature, material)
	Expect(err).ToNot(HaveOccurred())
	Expect(ok).To(BeTrue())

	// ecdsa
	signature, err = NewSignatureWithMaterial(material, KeyID(ecdsaKeyID), ecdsaKey, &ecdsaKey.PublicKey, tls.ECDSAWithP256AndSHA256)
	Expect(err).ToNot(HaveOccurred())
	ok, err = ValidateSignatureWithMaterial(keys, signature, material)
	Expect(err).ToNot(HaveOccurred())
	Expect(ok).To(BeTrue())

	// ed25519
	signature, err = NewSignatureWithMaterial(material, KeyID(ed25519KeyID), ed25519PrivKey, ed25519PubKey, tls.Ed25519)
	Expect(err).ToNot(HaveOccurred())
	ok, err = ValidateSignatureWithMaterial(keys, signature, material)
	Expect(err).ToNot(HaveOccurred())
	Expect(ok).To(BeTrue())
}

func TestParseSignatureAuthorizationPayload(t *testing.T) {
	RegisterTestingT(t)
	stringVal :="Signature k=YmFzZW1lbnQ, a=VGhpcyBpcyBhIHB1YmxpYyBrZXkgaW4gdXNlIGhlcmU, "+
	"s=2055, v=dmVyaWZpY2F0aW9uXzE2Qg, p=SW5zZXJ0IHNpZ25hdHVyZSBvZiBub25jZSBoZXJlIHdo"+
	"aWNoIHRha2VzIDUxMiBiaXRzIGZvciBFZDI1NTE5IQ"

	decodedKeyId := "basement"
	decodedPubKey := "This is a public key in use here"

	decodedVerification := "verification_16B"
	decodedSignature := "Insert signature of nonce here which takes 512 bits for Ed25519!"

	signature, err := ParseSignatureAuthorizationPayload(stringVal)
	Expect(err).ToNot(HaveOccurred())
	Expect(signature.keyID).To(Equal(KeyID(decodedKeyId)))
	Expect(signature.pubkey).To(BeEquivalentTo(decodedPubKey))
	Expect(signature.exporterVerification).To(BeEquivalentTo(decodedVerification))
	Expect(signature.proof).To(BeEquivalentTo(decodedSignature))
}