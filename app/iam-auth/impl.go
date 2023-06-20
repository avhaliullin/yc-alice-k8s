// Copyright (c) 2023 Yandex LLC. All rights reserved.
// Author: Andrey Khaliullin <avhaliullin@yandex-team.ru>

package iam_auth

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/avhaliullin/yandex-alice-k8s-skill/app/config"
	ycsdk "github.com/yandex-cloud/go-sdk"
	"github.com/yandex-cloud/go-sdk/iamkey"
)

type Deps interface {
	GetConfig() *config.Config
}

func New(deps Deps) (Service, error) {
	conf := deps.GetConfig()
	if len(conf.SAKey) > 0 {
		return NewStaticKeyBytes([]byte(conf.SAKey))
	} else {
		return NewMetadata()
	}
}

var _ Service = &metadata{}

type sdkAuth struct {
	sdk *ycsdk.SDK
}

func NewStaticKeyBytes(keyJSON []byte) (Service, error) {
	key, err := iamkey.ReadFromJSONBytes(keyJSON)
	if err != nil {
		return nil, err
	}
	return newStaticKey(key)
}

func NewStaticKeyPath(path string) (Service, error) {
	key, err := iamkey.ReadFromJSONFile(path)
	if err != nil {
		return nil, err
	}
	return newStaticKey(key)
}
func newStaticKey(key *iamkey.Key) (Service, error) {
	creds, err := ycsdk.ServiceAccountKey(key)
	if err != nil {
		return nil, err
	}
	sdk, err := ycsdk.Build(context.Background(), ycsdk.Config{
		Credentials: creds,
	})
	if err != nil {
		return nil, err
	}
	return &sdkAuth{sdk: sdk}, nil
}

func (s *sdkAuth) GetToken(ctx context.Context) (string, error) {
	resp, err := s.sdk.CreateIAMToken(ctx)
	if err != nil {
		return "", err
	}
	return resp.GetIamToken(), nil
}

type metadata struct {
	creds ycsdk.NonExchangeableCredentials
}

func NewMetadata() (Service, error) {
	return &metadata{creds: ycsdk.InstanceServiceAccount()}, nil
}

func (s *metadata) GetToken(ctx context.Context) (string, error) {
	resp, err := s.creds.IAMToken(ctx)
	if err != nil {
		return "", err
	}
	return resp.GetIamToken(), nil
}

var _ Service = &console{}

type console struct {
	profile string
}

func NewConsole(profile string) Service {
	return &console{profile: profile}
}

func (c *console) GetToken(ctx context.Context) (string, error) {
	cmd := exec.Command("yc", "--profile", c.profile, "iam", "create-token")
	respBytes, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("reading yc out: %w", err)
	}
	return strings.Trim(string(respBytes), " \n\t\r"), nil
}
