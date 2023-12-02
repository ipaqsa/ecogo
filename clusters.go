package goeco

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"time"
)

var (
	listClustersPath  = "/v1/cluster/%s/%s/list"
	createClusterPath = "/v1/cluster/%s/%s"

	versionPath = "/version"

	apiV1Path = "/v1/cluster/%s/%s/%d"

	repairClusterPath  = path.Join(apiV1Path, "repair")
	makeClusterHAPath  = path.Join(apiV1Path, "ha")
	healthyClusterPath = path.Join(apiV1Path, "health")
	//certClusterPath    = path.Join("cert")
)

var _ ClusterServiceI = &ClusterService{}

type ClusterServiceI interface {
	ServerVersion(ctx context.Context) (string, error)

	List(ctx context.Context) ([]ClusterShortInfo, error)

	Create(ctx context.Context, opt *ClusterOpt) (uint, error)
	Request(ctx context.Context, clusterID uint) (*Cluster, error)
	Delete(ctx context.Context, clusterID uint) error

	Repair(ctx context.Context, clusterID uint) error
	MakeHA(ctx context.Context, clusterID uint) error
	Healthy(ctx context.Context, clusterID uint) error

	//ClusterCert(ctx context.Context, clusterID uint)
}
type ClusterService struct {
	client *Client
}

type ClusterShortInfo struct {
	ID         uint      `json:"id" yaml:"id"`
	Processing bool      `json:"processing" yaml:"processing"`
	Name       string    `json:"name" yaml:"name"`
	State      string    `json:"state" yaml:"state"`
	Status     string    `json:"status" yaml:"status"`
	Existed    string    `json:"existed" yaml:"existed"`
	Created    time.Time `json:"created" yaml:"created"`
}
type Cluster struct {
	ID           uint      `json:"id" yaml:"id"`
	ProjectID    int       `json:"projectID" yaml:"projectID"`
	RegionID     int       `json:"regionID" yaml:"regionID"`
	Processing   bool      `json:"processing" yaml:"processing"`
	HA           bool      `json:"ha" yaml:"ha"`
	InternalLB   bool      `json:"internalLB" yaml:"internalLB"`
	Endpoint     string    `json:"endpoint" yaml:"endpoint"`
	KubeAuthType string    `json:"kubeAuthType" yaml:"kubeAuthType"`
	Version      string    `json:"version" yaml:"version"`
	Name         string    `json:"name" yaml:"name"`
	NetworkID    string    `json:"networkID" yaml:"networkID"`
	SubnetID     string    `json:"subnetID" yaml:"subnetID"`
	State        string    `json:"state" yaml:"state"`
	Status       string    `json:"status" yaml:"status"`
	Existed      string    `json:"existed" yaml:"existed"`
	Created      time.Time `json:"created" yaml:"created"`
	MastersPool  Pool      `json:"mastersPool" yaml:"mastersPool"`
	WorkersPools []Pool    `json:"workersPools" yaml:"workersPools"`
}
type ClusterOpt struct {
	HA          bool      `json:"ha" yaml:"ha"`
	InternalLB  bool      `json:"internalLB" yaml:"internalLB"`
	AuthType    string    `json:"authType" yaml:"authType"`
	Name        string    `json:"name" yaml:"name"`
	Version     string    `json:"version" yaml:"version"`
	APILbFlavor string    `json:"apiLbFlavor" yaml:"APILbFlavor"`
	NetworkID   string    `json:"networkID" yaml:"networkID"`
	SubnetID    string    `json:"subnetID" yaml:"subnetID"`
	MasterOpt   PoolOpt   `json:"masterOpt" yaml:"masterOpt"`
	WorkerOpts  []PoolOpt `json:"workerOpts" yaml:"workerOpts"`
}

func (s *ClusterService) ServerVersion(ctx context.Context) (string, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, versionPath, nil)
	if err != nil {
		return "", err
	}
	var version string
	if _, err = s.client.Do(ctx, req, &version); err != nil {
		return "", err
	}
	return version, nil
}

func (s *ClusterService) List(ctx context.Context) ([]ClusterShortInfo, error) {
	requestPath := fmt.Sprintf(listClustersPath, s.client.projectID, s.client.regionID)
	req, err := s.client.NewRequest(ctx, http.MethodGet, requestPath, nil)
	if err != nil {
		return nil, err
	}
	var clusters []ClusterShortInfo
	if _, err = s.client.Do(ctx, req, &clusters); err != nil {
		return nil, err
	}
	return clusters, nil
}

func (s *ClusterService) Create(ctx context.Context, opt *ClusterOpt) (uint, error) {
	requestPath := fmt.Sprintf(createClusterPath, s.client.projectID, s.client.regionID)
	req, err := s.client.NewRequest(ctx, http.MethodPost, requestPath, opt)
	if err != nil {
		return 0, err
	}

	var id uint
	if _, err = s.client.Do(ctx, req, &id); err != nil {
		return 0, err
	}
	return id, err
}
func (s *ClusterService) Request(ctx context.Context, clusterID uint) (*Cluster, error) {
	requestPath := fmt.Sprintf(apiV1Path, s.client.projectID, s.client.regionID, clusterID)
	req, err := s.client.NewRequest(ctx, http.MethodGet, requestPath, nil)
	if err != nil {
		return nil, err
	}

	var cluster Cluster
	if _, err = s.client.Do(ctx, req, &cluster); err != nil {
		return nil, err
	}
	return &cluster, err
}
func (s *ClusterService) Delete(ctx context.Context, clusterID uint) error {
	requestPath := fmt.Sprintf(apiV1Path, s.client.projectID, s.client.regionID, clusterID)
	req, err := s.client.NewRequest(ctx, http.MethodDelete, requestPath, nil)
	if err != nil {
		return err
	}

	var resp string
	if _, err = s.client.Do(ctx, req, &resp); err != nil {
		return err
	}
	return err
}

func (s *ClusterService) Repair(ctx context.Context, clusterID uint) error {
	requestPath := fmt.Sprintf(repairClusterPath, s.client.projectID, s.client.regionID, clusterID)
	req, err := s.client.NewRequest(ctx, http.MethodPost, requestPath, nil)
	if err != nil {
		return err
	}

	var resp string
	if _, err = s.client.Do(ctx, req, &resp); err != nil {
		return err
	}
	return err
}
func (s *ClusterService) MakeHA(ctx context.Context, clusterID uint) error {
	requestPath := fmt.Sprintf(makeClusterHAPath, s.client.projectID, s.client.regionID, clusterID)
	req, err := s.client.NewRequest(ctx, http.MethodPost, requestPath, nil)
	if err != nil {
		return err
	}

	var resp string
	if _, err = s.client.Do(ctx, req, &resp); err != nil {
		return err
	}
	return err
}

func (s *ClusterService) Healthy(ctx context.Context, clusterID uint) error {
	requestPath := fmt.Sprintf(healthyClusterPath, s.client.projectID, s.client.regionID, clusterID)
	req, err := s.client.NewRequest(ctx, http.MethodGet, requestPath, nil)
	if err != nil {
		return err
	}

	var resp string
	if _, err = s.client.Do(ctx, req, &resp); err != nil {
		return err
	}
	return err
}
