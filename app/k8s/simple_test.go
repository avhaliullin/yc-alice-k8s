// Copyright (c) 2023 Yandex LLC. All rights reserved.
// Author: Andrey Khaliullin <avhaliullin@yandex-team.ru>

package k8s

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	"github.com/avhaliullin/yandex-alice-k8s-skill/app/config"
	iam_auth "github.com/avhaliullin/yandex-alice-k8s-skill/app/iam-auth"
	"github.com/stretchr/testify/require"
)

func TestListNS(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	iamAuth := iam_auth.NewConsole("prod-yndxandr")
	app := &testApp{
		iamAuth: iamAuth,
		config: &config.Config{
			K8sHost: "https://158.160.0.81",
			K8sCA:   base64Decode(t, "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM1ekNDQWMrZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRJek1EVXlNakUzTURrME5sb1hEVE16TURVeE9URTNNRGswTmxvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTHBSCmQzYk5JeHhBR1dHemVmRGZ4cmh4M0dLY2dSL2NZQnlJY2VoZUxDOHhLR0Z0ME1oak8zazhLMmJMSlRGZWRsRXAKekdnRlJPYXIySFhLci9xZGNTOTVqY3FIanc5TUh2N3RaS0xVbUcyMjFSdTByVW5Kb0NHOVUwbFBxdS9pS1lQYQprdUZleTJhWWNHMGhSWG1pSzFjWWR3cTJsN012T3RTTDV0V0ZocVdGKzFMaGY5ZzNnZzZneWIvdHhtYTBSTGtjCmE2WFppTnFpVXU5Y0sySWtsbzlsVlNsaHVDbVRoZ2s4Y0o1Z1BQeDEwVzk3YTRaTVJjdkVxM0RNK3MrTExVR3MKRkVmOXVPeTF2UVc0M2U3Qm1pTm5TYnptVkd1ZC81djdoaHdwNFkvY0s0b2Z2dDdBbzZ3M1N1aWVuaEwyNlZzMQplQkdheWZuV3FNMUUzb3laaUtzQ0F3RUFBYU5DTUVBd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0hRWURWUjBPQkJZRUZDTmMvTkpwNk1vVEVKUmNwUHY2a3R2QmRVeldNQTBHQ1NxR1NJYjMKRFFFQkN3VUFBNElCQVFDM1YwUHBjVEZ0VEhhRitMMXNWSmFaY0JoUkFpTTMrKzlLZCtOWnNuMkZGRnJoNmhQRgpCeFlYR2pwSmZjclFCaDNaalNzWTR4Y2pJR1M5ajVjQkVvRXBDVjg0RjhxcUxzZ2dKa0YvcjJiQmpma0hYN0FXCjRYaTRyMjdQR0VLZGdnRlBmKzl4MFZEeG4zY093THUvdzdlT0Y3cG9nU2w0eHMwYThNMWNvSmY0SFRyZ1lJM0EKNDdGYlVzNDk0VFZTYTNhVVB2SStpWFFNcklUT2ZMQVBqU2RGUEpQQWF6R3VFdVpyQ1dRRkNzZjNVQ1ROMmlpVQpGYlg3OEJnVzhMK1VTZmkzMXlCeHhSeUxFSm5sU1J1ZnZSeC9WOXZJWHowak5RREdkTFRPYy9yWURkYUsyaWxkClk3VEpxOUcvbCsvT2o2ZUJ3Zmp2ZXVxanNpRHQ2akZrNlRzSAotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="),
		},
	}
	service, err := NewService(app)
	require.NoError(t, err)
	res, err := service.ListNamespaces(ctx, &ListNamespacesReq{})
	require.NoError(t, err)
	t.Logf("Found namespaces: %v", res)
}

var _ Deps = &testApp{}

type testApp struct {
	iamAuth iam_auth.Service
	config  *config.Config
}

func (t *testApp) GetIAMAuth() iam_auth.Service {
	return t.iamAuth
}

func (t *testApp) GetConfig() *config.Config {
	return t.config
}

func base64Decode(t *testing.T, s string) []byte {
	res, err := base64.StdEncoding.DecodeString(s)
	require.NoError(t, err)
	return res
}
