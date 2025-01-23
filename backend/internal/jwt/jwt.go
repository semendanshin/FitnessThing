package jwt

import (
	"context"
	"crypto/rsa"
	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/logger"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Credentials interface {
	NewWithClaims(claims jwt.Claims) (string, error)
	GetKey(*jwt.Token) (interface{}, error)
}

type PemCredentials struct {
	PrivateKey *rsa.PublicKey
	PublicKey  *rsa.PublicKey
}

func NewPemCredentials(privateKey *rsa.PublicKey, publicKey *rsa.PublicKey) *PemCredentials {
	return &PemCredentials{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
}

func (c *PemCredentials) NewWithClaims(claims jwt.Claims) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(c.PrivateKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (c *PemCredentials) GetKey(token *jwt.Token) (interface{}, error) {
	return c.PublicKey, nil
}

type SecretCredentials struct {
	Secret string
}

func NewSecretCredentials(secret string) *SecretCredentials {
	return &SecretCredentials{
		Secret: secret,
	}
}

func (c *SecretCredentials) NewWithClaims(claims jwt.Claims) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(c.Secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (c *SecretCredentials) GetKey(token *jwt.Token) (interface{}, error) {
	return []byte(c.Secret), nil
}

type Options struct {
	Credentials Credentials

	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type OptionsFunc func(*Options)

func WithCredentials(credentials Credentials) OptionsFunc {
	return func(o *Options) {
		o.Credentials = credentials
	}
}

func WithAccessTTL(ttl time.Duration) OptionsFunc {
	return func(o *Options) {
		o.AccessTTL = ttl
	}
}

func WithRefreshTTL(ttl time.Duration) OptionsFunc {
	return func(o *Options) {
		o.RefreshTTL = ttl
	}
}

type Provider struct {
	opts Options
}

func NewProvider(opts ...OptionsFunc) *Provider {
	p := &Provider{
		opts: Options{
			AccessTTL:  time.Hour,
			RefreshTTL: time.Hour * 24 * 7,
		},
	}

	for _, o := range opts {
		o(&p.opts)
	}

	if p.opts.Credentials == nil {
		logger.Panic("credentials are required")
	}

	return p
}

func (p *Provider) GeneratePair(ctx context.Context, userID, pairID domain.ID, atTime time.Time) (domain.Tokens, error) {
	accessClaims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": atTime.Add(p.opts.AccessTTL).Unix(),
		"iat": atTime.Unix(),
		"jti": pairID.String(),
	}

	refreshClaims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": atTime.Add(p.opts.RefreshTTL).Unix(),
		"iat": atTime.Unix(),
		"jti": pairID.String(),
	}

	accessToken, err := p.opts.Credentials.NewWithClaims(accessClaims)
	if err != nil {
		return domain.Tokens{}, err
	}

	refreshToken, err := p.opts.Credentials.NewWithClaims(refreshClaims)
	if err != nil {
		return domain.Tokens{}, err
	}

	return domain.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (p *Provider) VerifyPair(ctx context.Context, userID domain.ID, tokens domain.Tokens, atTime time.Time) error {
	accessClaims := jwt.MapClaims{}
	accessTokenParsed, err := jwt.ParseWithClaims(tokens.AccessToken, &accessClaims, p.opts.Credentials.GetKey, jwt.WithoutClaimsValidation())
	if err != nil {
		logger.Errorf("failed to parse access token: %v", err)
		return domain.ErrUnauthorized
	}

	refreshClaims := jwt.MapClaims{}
	refreshTokenParsed, err := jwt.ParseWithClaims(tokens.RefreshToken, &refreshClaims, p.opts.Credentials.GetKey)
	if err != nil {
		logger.Errorf("failed to parse refresh token: %v", err)
		return domain.ErrUnauthorized
	}

	if !accessTokenParsed.Valid || !refreshTokenParsed.Valid {
		logger.Errorf("access token or refresh token is invalid")
		return domain.ErrUnauthorized
	}

	accessSubject, err := accessClaims.GetSubject()
	if err != nil {
		logger.Errorf("access token subject not found: %v", err)
		return domain.ErrUnauthorized
	}

	refreshSubject, err := refreshClaims.GetSubject()
	if err != nil {
		logger.Errorf("refresh token subject not found: %v", err)
		return domain.ErrUnauthorized
	}

	if accessSubject != userID.String() || refreshSubject != userID.String() {
		logger.Errorf("access token and refresh token subject mismatch: %v != %v", accessSubject, refreshSubject)
		return domain.ErrUnauthorized
	}

	accessJTI, ok := accessClaims["jti"]
	if !ok {
		logger.Errorf("access token jti not found")
		return err
	}

	refreshJTI, ok := refreshClaims["jti"]
	if !ok {
		logger.Errorf("refresh token jti not found")
		return domain.ErrUnauthorized
	}

	if accessJTI != refreshJTI {
		logger.Errorf("access token and refresh token jti mismatch")
		return domain.ErrUnauthorized
	}

	return nil
}

func (p *Provider) ParseToken(ctx context.Context, token string) (domain.ID, error) {
	claims := jwt.MapClaims{}
	tokenParsed, err := jwt.ParseWithClaims(token, &claims, p.opts.Credentials.GetKey)
	if err != nil {
		logger.Errorf("failed to parse token: %v", err)
		return domain.ID{}, domain.ErrUnauthorized
	}

	if !tokenParsed.Valid {
		logger.Errorf("token is invalid")
		return domain.ID{}, domain.ErrUnauthorized
	}

	subject, err := claims.GetSubject()
	if err != nil {
		logger.Errorf("token subject not found: %v", err)
		return domain.ID{}, domain.ErrUnauthorized
	}

	return domain.ParseID(subject)
}
