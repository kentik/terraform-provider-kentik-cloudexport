package provider_test

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"testing"

	cloudexportpb "github.com/kentik/api-schema-public/gen/go/kentik/cloud_export/v202101beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const cloudExportNotFound = -1

type fakeCloudExportServer struct {
	cloudexportpb.UnimplementedCloudExportAdminServiceServer
	server *grpc.Server

	url  string
	done chan struct{}
	t    testing.TB

	data []*cloudexportpb.CloudExport
}

func newFakeCloudExportServer(t testing.TB, ces []*cloudexportpb.CloudExport) *fakeCloudExportServer {
	return &fakeCloudExportServer{
		done: make(chan struct{}),
		t:    t,
		data: ces,
	}
}

func (s *fakeCloudExportServer) Start() {
	l, err := net.Listen("tcp", "localhost:0")
	require.NoError(s.t, err)

	s.url = l.Addr().String()
	s.server = grpc.NewServer()
	cloudexportpb.RegisterCloudExportAdminServiceServer(s.server, s)

	go func() {
		err = s.server.Serve(l)
		assert.NoError(s.t, err)
		s.done <- struct{}{}
	}()
}

// Stop blocks until the server is stopped.
func (s *fakeCloudExportServer) Stop() {
	s.server.GracefulStop()
	<-s.done
}

// URL returns the server URL.
func (s *fakeCloudExportServer) URL() string {
	return fmt.Sprintf("http://%v", s.url)
}

func (s *fakeCloudExportServer) ListCloudExport(
	_ context.Context, _ *cloudexportpb.ListCloudExportRequest,
) (*cloudexportpb.ListCloudExportResponse, error) {
	return &cloudexportpb.ListCloudExportResponse{
		Exports:             s.data,
		InvalidExportsCount: 0,
	}, nil
}

func (s *fakeCloudExportServer) GetCloudExport(
	ctx context.Context, req *cloudexportpb.GetCloudExportRequest,
) (*cloudexportpb.GetCloudExportResponse, error) {
	if idx := s.findByID(req.GetId()); idx != cloudExportNotFound {
		return &cloudexportpb.GetCloudExportResponse{Export: s.data[idx]}, nil
	}
	return nil, fmt.Errorf("cloud export with ID %q not found", req.GetId())
}

func (s *fakeCloudExportServer) CreateCloudExport(
	ctx context.Context, req *cloudexportpb.CreateCloudExportRequest,
) (*cloudexportpb.CreateCloudExportResponse, error) {
	newExport := req.GetExport()

	if s.findByName(newExport.Name) != cloudExportNotFound {
		return nil, fmt.Errorf("cloud export %q already exists", newExport.Name)
	}

	newExport.Id = s.allocateNewID()
	newExport.CurrentStatus = &cloudexportpb.Status{
		Status:               "OK",
		ErrorMessage:         "No errors",
		FlowFound:            &wrapperspb.BoolValue{Value: true},
		ApiAccess:            &wrapperspb.BoolValue{Value: true},
		StorageAccountAccess: &wrapperspb.BoolValue{Value: true},
	}

	s.data = append(s.data, newExport)

	return &cloudexportpb.CreateCloudExportResponse{
		Export: newExport,
	}, nil
}

func (s *fakeCloudExportServer) UpdateCloudExport(
	ctx context.Context, req *cloudexportpb.UpdateCloudExportRequest,
) (*cloudexportpb.UpdateCloudExportResponse, error) {
	exportUpdate := req.GetExport()
	if i := s.findByID(exportUpdate.GetId()); i != cloudExportNotFound {
		s.data[i] = exportUpdate
		return &cloudexportpb.UpdateCloudExportResponse{
			Export: exportUpdate,
		}, nil
	}
	return nil, fmt.Errorf("cloud export of id %q doesn't exists", exportUpdate.Id)
}

func (s *fakeCloudExportServer) DeleteCloudExport(
	ctx context.Context, req *cloudexportpb.DeleteCloudExportRequest,
) (*cloudexportpb.DeleteCloudExportResponse, error) {
	return &cloudexportpb.DeleteCloudExportResponse{}, nil
}

func (s *fakeCloudExportServer) allocateNewID() string {
	var id int

	for _, item := range s.data {
		itemID, err := strconv.Atoi(item.Id)
		if err != nil {
			itemID = 1000000 // str conversion error, assume some high integer for the id
		}
		if itemID > id {
			id = itemID
		}
	}
	return strconv.FormatInt(int64(id)+1, 10)
}

func (s *fakeCloudExportServer) findByName(name string) int {
	for i, ce := range s.data {
		if ce.Name == name {
			return i
		}
	}
	return cloudExportNotFound
}

func (s *fakeCloudExportServer) findByID(id string) int {
	for i, ce := range s.data {
		if ce.Id == id {
			return i
		}
	}
	return cloudExportNotFound
}

func makeInitialCloudExports() []*cloudexportpb.CloudExport {
	return []*cloudexportpb.CloudExport{
		{
			Id:          "1",
			Type:        cloudexportpb.CloudExportType_CLOUD_EXPORT_TYPE_KENTIK_MANAGED,
			Enabled:     true,
			Name:        "test_terraform_aws_export",
			Description: "terraform aws cloud export",
			PlanId:      "11467",
			Bgp: &cloudexportpb.BgpProperties{
				ApplyBgp:       true,
				UseBgpDeviceId: "dummy-device-id",
				DeviceBgpType:  "dummy-device-bgp-type",
			},
			CurrentStatus: &cloudexportpb.Status{
				Status:               "OK",
				ErrorMessage:         "No errors",
				FlowFound:            &wrapperspb.BoolValue{Value: true},
				ApiAccess:            &wrapperspb.BoolValue{Value: true},
				StorageAccountAccess: &wrapperspb.BoolValue{Value: true},
			},
			CloudProvider: "aws",
			Properties: &cloudexportpb.CloudExport_Aws{
				Aws: &cloudexportpb.AwsProperties{
					Bucket:          "terraform-aws-bucket",
					IamRoleArn:      "arn:aws:iam::003740049406:role/trafficTerraformIngestRole",
					Region:          "us-east-2",
					DeleteAfterRead: false,
					MultipleBuckets: false,
				},
			},
		},
		{
			Id:          "2",
			Type:        cloudexportpb.CloudExportType_CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED,
			Enabled:     true,
			Name:        "test_terraform_gce_export",
			Description: "terraform gce cloud export",
			PlanId:      "21600",
			CurrentStatus: &cloudexportpb.Status{
				Status:       "NOK",
				ErrorMessage: "Timeout",
			},
			CloudProvider: "gce",
			Properties: &cloudexportpb.CloudExport_Gce{
				Gce: &cloudexportpb.GceProperties{
					Project:      "project gce",
					Subscription: "subscription gce",
				},
			},
		},
		{
			Id:          "3",
			Type:        cloudexportpb.CloudExportType_CLOUD_EXPORT_TYPE_KENTIK_MANAGED,
			Enabled:     false,
			Name:        "test_terraform_ibm_export",
			Description: "terraform ibm cloud export",
			PlanId:      "11467",
			CurrentStatus: &cloudexportpb.Status{
				Status:       "OK",
				ErrorMessage: "No errors",
			},
			CloudProvider: "ibm",
			Properties: &cloudexportpb.CloudExport_Ibm{
				Ibm: &cloudexportpb.IbmProperties{
					Bucket: "terraform-ibm-bucket",
				},
			},
		},
		{
			Id:          "4",
			Type:        cloudexportpb.CloudExportType_CLOUD_EXPORT_TYPE_KENTIK_MANAGED,
			Enabled:     true,
			Name:        "test_terraform_azure_export",
			Description: "terraform azure cloud export",
			PlanId:      "11467",
			CurrentStatus: &cloudexportpb.Status{
				Status:       "OK",
				ErrorMessage: "No errors",
			},
			CloudProvider: "azure",
			Properties: &cloudexportpb.CloudExport_Azure{
				Azure: &cloudexportpb.AzureProperties{
					Location:                 "centralus",
					ResourceGroup:            "traffic-generator",
					StorageAccount:           "kentikstorage",
					SubscriptionId:           "784bd5ec-122b-41b7-9719-22f23d5b49c8",
					SecurityPrincipalEnabled: true,
				},
			},
		},
	}
}
