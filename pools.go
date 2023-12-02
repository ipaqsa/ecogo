package goeco

import (
	"context"
	"fmt"
	"net/http"
	"path"
)

var (
	poolBasePath   = path.Join(apiV1Path, "pool")
	poolBasePathUD = path.Join(poolBasePath, "%s")
)

var _ PoolServiceI = &PoolService{}

type PoolServiceI interface {
	Create(ctx context.Context, clusterID uint, opt *PoolOpt) error
	Update(ctx context.Context, clusterID, poolID uint, opt *PoolOpt) error
	Delete(ctx context.Context, clusterID, poolID uint, opt *PoolOpt) error
}
type PoolService struct {
	client *Client
}

type Pool struct {
	ID           uint   `json:"id"`
	NodeCount    int    `json:"nodeCount"`
	MaxNodeCount int    `json:"maxNodeCount"`
	MinNodeCount int    `json:"minNodeCount"`
	VolumeSize   int    `json:"volumeSize"`
	SetTaint     bool   `json:"set-k8s-taint"`
	Role         string `json:"k8s-role"`
	Name         string `json:"name"`
	Flavor       string `json:"flavor"`
	State        string `json:"state"`
	Status       string `json:"status"`
	VolumeType   string `json:"volumeType"`
}
type PoolOpt struct {
	NodeCount    int    `json:"nodeCount"`
	MaxNodeCount int    `json:"maxNodeCount"`
	MinNodeCount int    `json:"minNodeCount"`
	VolumeSize   int    `json:"volumeSize"`
	Repair       bool   `json:"repair"`
	SetTaint     bool   `json:"set-taint"`
	Role         string `json:"k8s-role"`
	VolumeType   string `json:"volumeType"`
	Name         string `json:"name"`
	Flavor       string `json:"flavor"`
}

func (s *PoolService) Create(ctx context.Context, clusterID uint, opt *PoolOpt) error {
	requestPath := fmt.Sprintf(poolBasePath, s.client.projectID, s.client.regionID, clusterID)
	req, err := s.client.NewRequest(ctx, http.MethodPost, requestPath, opt)
	if err != nil {
		return err
	}

	var resp string
	if _, err = s.client.Do(ctx, req, &resp); err != nil {
		return err
	}
	return nil
}
func (s *PoolService) Update(ctx context.Context, clusterID, poolID uint, opt *PoolOpt) error {
	requestPath := fmt.Sprintf(poolBasePathUD, s.client.projectID, s.client.regionID, clusterID, poolID)
	req, err := s.client.NewRequest(ctx, http.MethodPatch, requestPath, opt)
	if err != nil {
		return err
	}

	var resp string
	if _, err = s.client.Do(ctx, req, &resp); err != nil {
		return err
	}
	return nil
}
func (s *PoolService) Delete(ctx context.Context, clusterID, poolID uint, opt *PoolOpt) error {
	requestPath := fmt.Sprintf(poolBasePathUD, s.client.projectID, s.client.regionID, clusterID, poolID)
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
