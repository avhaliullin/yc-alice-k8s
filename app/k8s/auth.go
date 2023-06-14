// Copyright (c) 2023 Yandex LLC. All rights reserved.
// Author: Andrey Khaliullin <avhaliullin@yandex-team.ru>

package k8s

import (
	"fmt"
	"net/http"

	iam_auth "github.com/avhaliullin/yandex-alice-k8s-skill/app/iam-auth"
	"k8s.io/client-go/rest"
)

func authProviderFactory(iam iam_auth.Service) rest.Factory {
	return func(_ string, _ map[string]string, _ rest.AuthProviderConfigPersister) (rest.AuthProvider, error) {
		return &authProvider{iamAuth: iam}, nil
	}
}

var _ rest.AuthProvider = &authProvider{}

type authProvider struct {
	iamAuth iam_auth.Service
}

func (p *authProvider) WrapTransport(r http.RoundTripper) http.RoundTripper {
	return &authRoundTripper{
		iamAuth: p.iamAuth,
		base:    r,
	}
}

func (p *authProvider) Login() error {
	return fmt.Errorf("not implemented")
}

var _ http.RoundTripper = &authRoundTripper{}

type authRoundTripper struct {
	iamAuth iam_auth.Service
	base    http.RoundTripper
}

func (r *authRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Header.Get("Authorization") != "" {
		return r.base.RoundTrip(req)
	}

	token, err := r.iamAuth.GetToken(req.Context())
	if err != nil {
		return nil, fmt.Errorf("getting credentials: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := r.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
