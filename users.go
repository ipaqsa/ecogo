package goeco

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

var (
	userPath  = path.Join(apiV1Path, "/kube/user")
	adminPath = path.Join(apiV1Path, "/kube/admin")
)

var _ UserServiceI = &UserService{}

type UserServiceI interface {
	RequestAdminConfig(ctx context.Context, clusterID uint, ttl string) (*UserConfig, error)
	RequestUserConfig(ctx context.Context, clusterID uint, opt *UserOpt) (*UserConfig, error)
	DeleteUser(ctx context.Context, clusterID uint, opt *UserOpt) error
}
type UserService struct {
	client *Client
}

type UserConfig struct {
	Content string `json:"content"`
}

type UserOpt struct {
	Role       string `json:"role"`
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	SecondsExp uint64 `json:"secondsExp"`
}

func (s *UserService) RequestAdminConfig(ctx context.Context, clusterID uint, ttl string) (*UserConfig, error) {
	requestPath := fmt.Sprintf(adminPath, s.client.projectID, s.client.regionID, clusterID)
	req, err := s.client.NewRequest(ctx, http.MethodGet, requestPath, nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = url.Values{
		"ttl": {ttl},
	}.Encode()
	var config UserConfig
	if _, err = s.client.Do(ctx, req, &config); err != nil {
		return nil, err
	}
	return &config, err
}
func (s *UserService) RequestUserConfig(ctx context.Context, clusterID uint, opt *UserOpt) (*UserConfig, error) {
	requestPath := fmt.Sprintf(userPath, s.client.projectID, s.client.regionID, clusterID)
	req, err := s.client.NewRequest(ctx, http.MethodPost, requestPath, opt)
	if err != nil {
		return nil, err
	}

	var config UserConfig
	if _, err = s.client.Do(ctx, req, &config); err != nil {
		return nil, err
	}
	return &config, err
}
func (s *UserService) DeleteUser(ctx context.Context, clusterID uint, opt *UserOpt) error {
	requestPath := fmt.Sprintf(userPath, s.client.projectID, s.client.regionID, clusterID)
	req, err := s.client.NewRequest(ctx, http.MethodDelete, requestPath, opt)
	if err != nil {
		return err
	}

	var resp string
	if _, err = s.client.Do(ctx, req, &resp); err != nil {
		return err
	}
	return nil
}
