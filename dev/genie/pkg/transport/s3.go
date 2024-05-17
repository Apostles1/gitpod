// Copyright (c) 2024 Gitpod GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License.AGPL.txt in the project root for license information.

package transport

import (
	"context"
	"fmt"
)

type S3Config struct {
}

var _ Transport = &S3Transport{}

type S3Transport struct {
	Config *S3Config
}

func NewS3Transport(cfg *S3Config) (*S3Transport, error) {
	return &S3Transport{
		Config: cfg,
	}, nil
}

func (t *S3Transport) CreateSession(ctx context.Context, sessionId string) error {
	return fmt.Errorf("not implemented")
}

func (t *S3Transport) HasSession(ctx context.Context, sessionId string) bool {
	return false
}

func (t *S3Transport) WatchSessions(ctx context.Context) (<-chan string, error) {
	return nil, fmt.Errorf("not implemented")
}

func (t *S3Transport) WatchRequests(ctx context.Context, sessionId string) (<-chan *Message, error) {
	return nil, fmt.Errorf("not implemented")
}

func (t *S3Transport) GetLastRequestID(ctx context.Context, sessionId string) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

func (t *S3Transport) SendUnary(ctx context.Context, sessionId string, msg *Message) (*Message, error) {
	return nil, fmt.Errorf("not implemented")
}

func (t *S3Transport) SendResponse(ctx context.Context, sessionId string, msg *Message) error {
	return fmt.Errorf("not implemented")
}

func (t *S3Transport) SendStream(ctx context.Context, sessionId string, msg *Message) (<-chan *Message, error) {
	return nil, fmt.Errorf("not implemented")
}
