package auth

import (
	"context"
	"coolcar/shared/auth/token"
	"coolcar/shared/id"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	ImpersonateAccountHeader = "impersonate-account-id"
	authorizationHeader      = "authorization"
	bearerPrefix             = "Bearer "
)

type tokenVerifier interface {
	Verify(token string) (string, error)
}
type interceptor struct {
	verifier tokenVerifier
}

type accountIDKey struct{}

func Interceptor(publicKeyFile string) (grpc.UnaryServerInterceptor, error) {
	pubKeyFile, err := os.Open(publicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("cannot open public key file: %v", err)
	}
	pubKeyBytes, err := ioutil.ReadAll(pubKeyFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read public key: %v", err)
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("cannot parse public key: %v", err)
	}
	i := &interceptor{
		verifier: &token.JWTTokenVerifier{
			PublicKey: pubKey,
		},
	}
	return i.HandleReq, nil
}

func (i *interceptor) HandleReq(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	aid := impersonationFromContext(ctx)
	if aid != "" {
		return handler(ContextWithAccountID(ctx, id.AccountID(aid)), req)
	}
	tkn, err := tokenFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "")
	}
	aid, err = i.verifier.Verify(tkn)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "token not valid: %v", err)
	}

	return handler(ContextWithAccountID(ctx, id.AccountID(aid)), req)
}

func impersonationFromContext(ctx context.Context) string {
	m, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	imp := m[ImpersonateAccountHeader]
	if len(imp) == 0 {
		return ""
	}
	return imp[0]
}

func ContextWithAccountID(c context.Context, aid id.AccountID) context.Context {
	return context.WithValue(c, accountIDKey{}, aid)
}

func tokenFromContext(ctx context.Context) (string, error) {
	unauthenticated := codes.Unauthenticated
	m, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(unauthenticated, "")
	}
	tkn := ""
	for _, v := range m[authorizationHeader] {
		if strings.HasPrefix(v, bearerPrefix) {
			tkn = v[len(bearerPrefix):]
		}
	}
	if tkn == "" {
		return "", status.Error(unauthenticated, "")
	}
	return tkn, nil
}
func AccountIDFromContext(c context.Context) (id.AccountID, error) {
	v := c.Value(accountIDKey{})
	aid, ok := v.(id.AccountID)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "")
	}
	return aid, nil
}
