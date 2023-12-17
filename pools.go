package goeco

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"time"
)

var (
	poolPath   = path.Join(apiV1Path, "pool")
	poolPathUD = path.Join(poolPath, "%d")
)

var _ PoolServiceI = &PoolService{}

type PoolServiceI interface {
	Create(ctx context.Context, clusterID uint, opt *PoolOpt) error
	Update(ctx context.Context, clusterID, poolID uint, opt *PoolOpt) error
	Delete(ctx context.Context, clusterID, poolID uint) error
}
type PoolService struct {
	client *Client
}

type Pool struct {
	ID           uint   `json:"id" yaml:"id"`
	NodeCount    int    `json:"nodeCount" yaml:"nodeCount"`
	MaxNodeCount int    `json:"maxNodeCount" yaml:"maxNodeCount"`
	MinNodeCount int    `json:"minNodeCount" yaml:"minNodeCount"`
	VolumeSize   int    `json:"volumeSize" yaml:"volumeSize"`
	SetTaint     bool   `json:"set-k8s-taint" yaml:"setTaint"`
	Role         string `json:"k8s-role" yaml:"role"`
	Name         string `json:"name" yaml:"name"`
	Flavor       string `json:"flavor" yaml:"flavor"`
	State        string `json:"state" yaml:"state"`
	Status       string `json:"status" yaml:"status"`
	VolumeType   string `json:"volumeType" yaml:"volumeType"`
}
type PoolOpt struct {
	NodeCount    int       `json:"nodeCount" yaml:"nodeCount"`
	MaxNodeCount int       `json:"maxNodeCount" yaml:"maxNodeCount"`
	MinNodeCount int       `json:"minNodeCount" yaml:"minNodeCount"`
	VolumeSize   int       `json:"volumeSize" yaml:"volumeSize"`
	Repair       bool      `json:"repair" yaml:"repair"`
	SetTaint     bool      `json:"set-taint" yaml:"setTaint"`
	Created      time.Time `json:"created" yaml:"created"`
	Role         string    `json:"k8s-role" yaml:"role"`
	VolumeType   string    `json:"volumeType" yaml:"volumeType"`
	Name         string    `json:"name" yaml:"name"`
	Flavor       string    `json:"flavor" yaml:"flavor"`
}

func (s *PoolService) Create(ctx context.Context, clusterID uint, opt *PoolOpt) error {
	requestPath := fmt.Sprintf(poolPath, s.client.projectID, s.client.regionID, clusterID)
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
	requestPath := fmt.Sprintf(poolPathUD, s.client.projectID, s.client.regionID, clusterID, poolID)
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
func (s *PoolService) Delete(ctx context.Context, clusterID, poolID uint) error {
	requestPath := fmt.Sprintf(poolPathUD, s.client.projectID, s.client.regionID, clusterID, poolID)
	req, err := s.client.NewRequest(ctx, http.MethodDelete, requestPath, nil)
	if err != nil {
		return err
	}

	var resp string
	if _, err = s.client.Do(ctx, req, &resp); err != nil {
		return err
	}
	return nil
}
