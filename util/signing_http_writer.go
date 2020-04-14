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
	"sync"
)

var dummyKey *rsa.PrivateKey
var dummyKeyOnce sync.Once

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
		privkey: getKey(),
	}

	return sw
}

// Return the RSA key
//
// XXX: Just returns a random new key for now. Not actually useful
// for signing real responses.
func getKey() *rsa.PrivateKey {
	dummyKeyOnce.Do(func() {
		var err error

		if dummyKey, err = rsa.GenerateKey(rand.Reader, 2048); err != nil {
			log.Fatal(err)
		}
	})

	return dummyKey
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

	signer := pw.newSigner()

	if err := signer.SignResponse(privkey, hostname, w, body); err != nil {
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

func (pw *SigningHTTPWriter) newSigner() httpsig.Signer {
	// as contains the algorithms we consider using
	as := []httpsig.Algorithm{
		httpsig.RSA_SHA512, httpsig.RSA_SHA256,
	}

	// the digest algorithm to use
	da := httpsig.DigestSha256

	// hs contains the headers that our signature will cover
	var hs []string
	//hs = append(hs, httpsig.RequestTarget)
	hs = append(hs, "date")
	if pw.body.Len() > 0 {
		hs = append(hs, "digest")
	}

	// create the signer; the api object we will use to sign
	// our response
	signer, _, err := httpsig.NewSigner(as, da, hs, httpsig.Signature)
	if err != nil {
		log.Fatal(err)
	}

	return signer
}

// Add all headers in src to dst. Existing headers in dst are preserved.
func copyHeaders(dst, src http.Header) {
	for key, values := range src {
		dst[key] = append(src[key], values...)
	}
}
