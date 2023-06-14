// Copyright (c) 2023 Yandex LLC. All rights reserved.
// Author: Andrey Khaliullin <avhaliullin@yandex-team.ru>

package iam_auth

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	ycsdk "github.com/yandex-cloud/go-sdk"
)

var _ Service = &metadata{}

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
