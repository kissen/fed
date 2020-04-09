package util

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"github.com/go-fed/httpsig"
	"gitlab.cs.fau.de/kissen/fed/config"
	"gitlab.cs.fau.de/kissen/fed/errors"
	"io"
	"log"
	"net/http"
)

// Implements http.ResponseWriter
type SigningHTTPWriter struct {
	// The body of the response. We cache it so we can calculate
	// a digest before sending it out.
	body bytes.Buffer

	// The headers set ton the respone.
	header http.Header

	// The return status as set by WriteHeaders. Initialized to
	// http.StatusOK in NewSigningWriter.
	status int

	// The signer used to sign the request.
	signer httpsig.Signer

	// The private key used to sign the request.
	privkey *rsa.PrivateKey
}

// Create a new placeholder HTTP response writer that implements
// the http.ResponseWriter interface. You can use it to record
// interactions to an http.ResponseWriter and then "replay" them
// by applying them to another (real) response writer.
func NewSigningWriter() *SigningHTTPWriter {
	sw := &SigningHTTPWriter{
		header:  make(http.Header),
		status:  http.StatusOK,
		signer:  newSigner(),
		privkey: newKey(),
	}

	return sw
}

func newSigner() httpsig.Signer {
	as := []httpsig.Algorithm{
		httpsig.RSA_SHA512, httpsig.RSA_SHA256,
	}

	hs := []string{
		//httpsig.RequestTarget,
		"date", "digest",
	}

	signer, _, err := httpsig.NewSigner(
		as, httpsig.DigestSha256, hs, httpsig.Signature,
	)

	if err != nil {
		log.Fatal(err)
	}

	return signer
}

func newKey() *rsa.PrivateKey {
	privkey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	return privkey
}

func (pw *SigningHTTPWriter) Header() http.Header {
	return pw.header
}

func (pw *SigningHTTPWriter) Write(bs []byte) (int, error) {
	return pw.body.Write(bs)
}

func (pw *SigningHTTPWriter) WriteHeader(status int) {
	pw.status = status
}

// Return a copy of the response body.
func (pw *SigningHTTPWriter) Body() []byte {
	bs := pw.body.Bytes()
	return append([]byte(nil), bs...)
}

// Apply all operations that were done on this placeholder to
// response writer w.
func (pw *SigningHTTPWriter) ApplyTo(w http.ResponseWriter) error {
	// set up headers

	SetDateHeader(pw.Header())
	delete(pw.Header(), "Digest") // will be re-set by SignResponse
	copyHeaders(w.Header(), pw.Header())

	// sign response

	body := pw.Body()
	privkey := pw.privkey
	hostname := config.Get().Hostname

	if err := pw.signer.SignResponse(privkey, hostname, w, body); err != nil {
		return errors.Wrap(err, "could not sign response")
	}

	// write out to w

	copyHeaders(w.Header(), pw.Header())
	w.WriteHeader(pw.status)

	if _, err := io.Copy(w, &pw.body); err != nil {
		return errors.Wrap(err, "copying buffered body failed")
	}

	return nil
}

// Add all headers in src to dst. Existing headers in dst are preserved.
func copyHeaders(dst, src http.Header) {
	for key, values := range src {
		dst[key] = append(src[key], values...)
	}
}
