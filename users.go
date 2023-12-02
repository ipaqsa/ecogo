package goeco

import (
	"context"
	"fmt"
	"net/http"
	"path"
)

var (
	listUsersPath = path.Join(apiV1Path, "users/%s")

	usersBasePathRD = path.Join(apiV1Path, "/users")

	listRolesPath = "/v1/users/roles"
	//rolePath = path.Join(listRolesPath, "%s")
)

var _ UserServiceI = &UserService{}
var _ RoleServiceI = &RoleService{}

type UserServiceI interface {
	List(ctx context.Context, clusterID uint, namespace string) (Users, error)

	RequestConfig(ctx context.Context, clusterID uint, opt *UserOpt) (*Config, error)
	Delete(ctx context.Context, clusterID uint, opt *UserOpt) error
}
type UserService struct {
	client *Client
}

type Users struct {
	Users []*UserOpt `json:"users"`
}
type UserOpt struct {
	Role       string `json:"role"`
	Name       string `json:"name"`
	User       string `json:"user"`
	Namespace  string `json:"namespace"`
	Created    string `json:"created"`
	Expired    string `json:"expired"`
	SecondsExp uint64 `json:"secondsExp"`
}

type Config struct {
	Content string `json:"content"`
}

type RoleServiceI interface {
	List(ctx context.Context) ([]string, error)

	//Request(ctx context.Context, role string)
}
type RoleService struct {
	client *Client
}

func (s *UserService) List(ctx context.Context, clusterID uint, namespace string) (Users, error) {
	requestPath := fmt.Sprintf(listUsersPath, s.client.projectID, s.client.regionID, clusterID, namespace)
	req, err := s.client.NewRequest(ctx, http.MethodGet, requestPath, nil)
	if err != nil {
		return Users{}, err
	}

	var users Users
	if _, err = s.client.Do(ctx, req, &users); err != nil {
		return Users{}, err
	}
	return users, err
}

func (s *UserService) RequestConfig(ctx context.Context, clusterID uint, opt *UserOpt) (*Config, error) {
	requestPath := fmt.Sprintf(usersBasePathRD, s.client.projectID, s.client.regionID, clusterID)
	req, err := s.client.NewRequest(ctx, http.MethodPost, requestPath, opt)
	if err != nil {
		return nil, err
	}

	var config Config
	if _, err = s.client.Do(ctx, req, &config); err != nil {
		return nil, err
	}
	return &config, err
}
func (s *UserService) Delete(ctx context.Context, clusterID uint, opt *UserOpt) error {
	requestPath := fmt.Sprintf(usersBasePathRD, s.client.projectID, s.client.regionID, clusterID)
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

func (s *RoleService) List(ctx context.Context) ([]string, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, listRolesPath, nil)
	if err != nil {
		return nil, err
	}

	var roles []string
	if _, err = s.client.Do(ctx, req, &roles); err != nil {
		return nil, err
	}
	return roles, err
}
