package wts

import (
	"api/internal/frameworks/obj"
	"api/internal/frameworks/utils"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
	"go.uber.org/zap"
)

type Server struct {
	logger    *zap.Logger
	wt        *webtransport.Server
	handler   *Handler
	isRunning bool
	name      string
	cfg       *WebTransportServerConfig
}

var _ obj.Server = (*Server)(nil)

func NewServer(cfg *WebTransportServerConfig, logger *zap.Logger, hub *Hub, handlerRegistries HandlerRegisters) (*Server, error) {
	var (
		tlsConfig *tls.Config
		cert      *x509.Certificate
		err       error
	)
	name := "webtransport"
	logger = logger.Named(name)
	if cfg.UseSelfCert {
		tlsConfig, cert, err = makeSelfTLS(logger)
	} else {
		tlsConfig, cert, err = loadTLS(cfg.CertFile, cfg.KeyFile)
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if len(tlsConfig.Certificates) == 0 {
		return nil, errors.New("lost TLS certificate")
	}

	s := &Server{
		name:   name,
		logger: logger,
		cfg:    cfg,
		wt: &webtransport.Server{
			CheckOrigin: func(r *http.Request) bool {
				logger.Debug("Check Origin",
					zap.String("origin", r.Header.Get("Origin")),
				)
				return true
			},
			H3: &http3.Server{
				TLSConfig: http3.ConfigureTLSConfig(tlsConfig),
				// QUICConfig: &quic.Config{
				// 	Tracer:                           qlog.DefaultConnectionTracer,
				// 	EnableDatagrams:                  true,
				// 	EnableStreamResetPartialDelivery: true,
				// },
			},
		},
	}

	printTLS(logger, tlsConfig.Certificates[0], cert)

	if cfg.Host == "" {
		s.wt.H3.Addr = fmt.Sprintf(":%d", cfg.Port)
	} else {
		s.wt.H3.Addr = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	}

	if cfg.Path == "" {
		cfg.Path = "/"
	}

	var messageHandlerRegistry = newMessageHandlerRegistry()
	for _, handler := range handlerRegistries {
		handler.Register(messageHandlerRegistry)
	}

	s.handler = NewHandler(s.wt, hub, logger, messageHandlerRegistry)
	webtransport.ConfigureHTTP3Server(s.wt.H3)
	h := http.NewServeMux()
	h.HandleFunc("/readyz", s.readyzHandler)
	h.HandleFunc(cfg.Path, s.handler.ServeHTTP)
	s.wt.H3.Handler = h

	logger.Info("server configured", utils.ObjectToZapFields(cfg)...)
	return s, nil
}

func (s *Server) Name() string {
	return s.name
}

func (s *Server) GetCertHash() string {
	cert := s.wt.H3.TLSConfig.Certificates[0]
	certBytes := cert.Certificate[0]
	hash := sha256.Sum256(certBytes)
	return base64.RawStdEncoding.EncodeToString(hash[:])
}

func (s *Server) IsRunning() bool {
	return s.isRunning
}

func (s *Server) Start() {
	s.isRunning = true
	s.logger.Info("server started", zap.String("addr", s.wt.H3.Addr))
	if err := errors.WithStack(s.wt.ListenAndServe()); err != nil && !errors.Is(err, context.Canceled) {
		s.logger.Error("server stopped", zap.Error(err))
	}
	s.isRunning = false
}

func (s *Server) Close(_ context.Context) error {
	defer func() {
		s.logger.Info("server shutting down")
		s.isRunning = false
	}()
	return errors.WithStack(s.wt.Close())
}

func (*Server) readyzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(
		"Content-Type",
		"text/plain",
	)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}

func loadTLS(certFile, keyFile string) (*tls.Config, *x509.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	x509Cert, err := x509.ParseCertificate(
		cert.Certificate[0],
	)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	// fmt.Printf("privatekey %T\n", cert.PrivateKey)
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		// MinVersion:   tls.VersionTLS13,
	}, x509Cert, nil
}

func printTLS(logger *zap.Logger, tls tls.Certificate, cert *x509.Certificate) {

	var ipAddresses []string
	for _, ip := range cert.IPAddresses {
		ipAddresses = append(ipAddresses, ip.String())
	}
	var b = sha256.Sum256(tls.Certificate[0])
	var extKeyUsage []int = make([]int, len(cert.ExtKeyUsage))
	for i, k := range cert.ExtKeyUsage {
		extKeyUsage[i] = int(k)
	}

	logger.Info("TLS certificate loaded",
		zap.Int("cert chain length", len(tls.Certificate)),
		zap.Any("SerialNumber", cert.SerialNumber),
		zap.String("Subject", cert.Subject.String()),
		zap.String("Issuer:", cert.Issuer.String()),
		zap.Strings("DNSNames", cert.DNSNames),
		zap.Strings("IPAddresses", ipAddresses),
		zap.Int("Keyusage", int(cert.KeyUsage)),
		zap.String("NotBefore", cert.NotBefore.String()),
		zap.String("NotAfter", cert.NotAfter.String()),
		zap.Bool("BasicConstraintsValid", cert.BasicConstraintsValid),
		zap.Ints("ExtKeyUsage", extKeyUsage),
		zap.String("sha256 hash", base64.RawStdEncoding.EncodeToString(b[:])))
}

const privateKey = "webtransport-example-cert-key-01"

func makeSelfTLS(logger *zap.Logger) (*tls.Config, *x509.Certificate, error) {
	privateKeyBytes := []byte(privateKey)
	logger.Info("Use Self Signed Certificate",
		zap.String("privateKey", privateKey),
		zap.Binary("privateKeyBytes", privateKeyBytes))
	certKey, err := ecdsa.ParseRawPrivateKey(elliptic.P256(), privateKeyBytes)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	// The W3C serverCertificateHashes API requires the certificate to be valid now,
	// but for less than two weeks. A weekly window keeps the hash stable across restarts.
	now := time.Now().UTC()
	validFrom := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	validFrom = validFrom.AddDate(0, 0, -int((validFrom.Weekday()+6)%7))
	validUntil := validFrom.Add(13 * 24 * time.Hour)

	certTemplate := x509.Certificate{
		SerialNumber:          big.NewInt(validFrom.Unix()),
		Subject:               pkix.Name{CommonName: "localhost"},
		NotBefore:             validFrom,
		NotAfter:              validUntil,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1), net.ParseIP("::1")},
	}
	// Passing nil makes ECDSA signing deterministic according to RFC 6979.
	certDER, err := x509.CreateCertificate(nil, &certTemplate, &certTemplate, &certKey.PublicKey, certKey)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{certDER},
			PrivateKey:  certKey,
		}},
	}, &certTemplate, nil
}
