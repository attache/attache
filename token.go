package attache

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

var gvTokenHeader = base64UrlEncode([]byte(`{"typ":"JWT", "alg":"HS256"`))

type TokenHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type TokenClaims map[string]interface{}

type Token struct {
	conf TokenConfig

	Header TokenHeader
	Claims TokenClaims
}

func (t Token) ClearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Name:     t.conf.Cookie,
		MaxAge:   -1,
		Value:    "",
	})
}

func (t Token) SaveCookie(w http.ResponseWriter) error {
	data, err := t.Encode()
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Name:     t.conf.Cookie,
		MaxAge:   t.conf.MaxAge,
		Value:    string(data),
	})

	return nil
}

func (t Token) MustSaveCookie(w http.ResponseWriter) {
	if err := t.SaveCookie(w); err != nil {
		ErrorFatal(err)
	}
}

func (t Token) encodeClaims() ([]byte, error) {
	if t.Claims == nil {
		return []byte("{}"), nil
	}

	data, err := json.Marshal(t.Claims)
	if err != nil {
		return nil, err
	}

	return base64UrlEncode(data), nil
}

func (t Token) encodeHeader() ([]byte, error) {
	data, err := json.Marshal(t.Header)
	if err != nil {
		return nil, err
	}

	return base64UrlEncode(data), nil
}

func (t Token) Validate() error {
	if t.Claims == nil {
		return errors.New("invalid token (no expiration)")
	}

	now := time.Now()
	switch got := t.Claims["exp"].(type) {
	case float64:
		if time.Unix(int64(got), 0).Before(now) {
			return errors.New("invalid token (expired)")
		}
	case int64:
		if time.Unix(got, 0).Before(now) {
			return errors.New("invalid token (expired)")
		}
	default:
		return errors.New("invalid token (no expiration)")
	}

	return nil
}

// Encode encodes and signs the Token in JWT format
func (t Token) Encode() ([]byte, error) {
	if err := t.Validate(); err != nil {
		return nil, err
	}

	head, err := t.encodeHeader()
	if err != nil {
		return nil, err
	}

	claims, err := t.encodeClaims()
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(
		make([]byte, 0, encodedLen(head)+encodedLen(claims)+1+256),
	)
	buf.Write(head)
	buf.WriteByte('.')
	buf.Write(claims)
	signStr := buf.Bytes()
	buf.WriteByte('.')
	buf.Write(base64UrlEncode(signatureFor(signStr, t.conf.Secret)))
	return buf.Bytes(), nil
}

func (t *Token) reset() { t.Claims = TokenClaims{} }

func (t *Token) Decode(data []byte) error {
	t.reset()

	parts := bytes.SplitN(data, []byte{'.'}, 3)
	if len(parts) != 3 {
		return errors.New("malformed token data")
	}

	head64, claims64, sig64 := parts[0], parts[1], parts[2]

	rawSig, err := base64UrlDecode(sig64)
	if err != nil {
		return errors.New("malformed token data")
	}

	rawHead, err := base64UrlDecode(head64)
	if err != nil {
		return errors.New("malformed token data")
	}

	rawClaims, err := base64UrlDecode(claims64)
	if err != nil {
		return errors.New("malformed token data")
	}

	if test := signatureFor(data[:len(head64)+len(claims64)+1], t.conf.Secret); !hmac.Equal(rawSig, test) {
		return errors.New("signature mismatch")
	}

	if err = json.Unmarshal(rawHead, &t.Header); err != nil {
		t.reset()
		return err
	}

	if err = json.Unmarshal(rawClaims, &t.Claims); err != nil {
		return err
	}

	return t.Validate()
}

func encodedLen(b []byte) int { return base64.RawURLEncoding.EncodedLen(len(b)) }
func decodedLen(b []byte) int { return base64.RawURLEncoding.DecodedLen(len(b)) }

func base64UrlEncode(src []byte) []byte {
	out := make([]byte, encodedLen(src))
	base64.RawURLEncoding.Encode(out, src)
	return out
}

func base64UrlDecode(src []byte) ([]byte, error) {
	out := make([]byte, decodedLen(src))
	if _, err := base64.RawURLEncoding.Decode(out, src); err != nil {
		return nil, err
	}
	return out, nil
}

func signatureFor(data, secret []byte) []byte {
	sig := hmac.New(sha256.New, secret)
	sig.Write(data)
	return sig.Sum(nil)
}
