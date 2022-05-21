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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const cloudExportNotFound = -1

type testAPIServer struct {
	cloudexportpb.UnimplementedCloudExportAdminServiceServer
	server *grpc.Server

	url  string
	done chan struct{}
	t    testing.TB

	data []*cloudexportpb.CloudExport
}

func newTestAPIServer(t testing.TB, ces []*cloudexportpb.CloudExport) *testAPIServer {
	return &testAPIServer{
		done: make(chan struct{}),
		t:    t,
		data: ces,
	}
}

func (s *testAPIServer) Start() {
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
func (s *testAPIServer) Stop() {
	s.server.GracefulStop()
	<-s.done
}

// URL returns the server URL.
func (s *testAPIServer) URL() string {
	return fmt.Sprintf("http://%v", s.url)
}

func (s *testAPIServer) ListCloudExport(
	_ context.Context, _ *cloudexportpb.ListCloudExportRequest,
) (*cloudexportpb.ListCloudExportResponse, error) {
	return &cloudexportpb.ListCloudExportResponse{
		Exports:             s.data,
		InvalidExportsCount: 0,
	}, nil
}

func (s *testAPIServer) GetCloudExport(
	ctx context.Context, req *cloudexportpb.GetCloudExportRequest,
) (*cloudexportpb.GetCloudExportResponse, error) {
	if idx := s.findByID(req.GetId()); idx != cloudExportNotFound {
		return &cloudexportpb.GetCloudExportResponse{Export: s.data[idx]}, nil
	}
	return nil, status.Errorf(codes.NotFound, "cloud export with ID %q not found", req.GetId())
}

func (s *testAPIServer) CreateCloudExport(
	ctx context.Context, req *cloudexportpb.CreateCloudExportRequest,
) (*cloudexportpb.CreateCloudExportResponse, error) {
	newExport := req.GetExport()

	if s.findByName(newExport.Name) != cloudExportNotFound {
		return nil, status.Errorf(codes.AlreadyExists, "cloud export %q already exists", newExport.Name)
	}

	newID, err := s.allocateNewID()
	if err != nil {
		return nil, err
	}
	newExport.Id = newID
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

func (s *testAPIServer) UpdateCloudExport(
	ctx context.Context, req *cloudexportpb.UpdateCloudExportRequest,
) (*cloudexportpb.UpdateCloudExportResponse, error) {
	exportUpdate := req.GetExport()
	if i := s.findByID(exportUpdate.GetId()); i != cloudExportNotFound {
		s.data[i] = exportUpdate
		return &cloudexportpb.UpdateCloudExportResponse{
			Export: exportUpdate,
		}, nil
	}
	return nil, status.Errorf(codes.NotFound, "cloud export of id %q doesn't exists", exportUpdate.Id)
}

func (s *testAPIServer) DeleteCloudExport(
	ctx context.Context, req *cloudexportpb.DeleteCloudExportRequest,
) (*cloudexportpb.DeleteCloudExportResponse, error) {
	if i := s.findByID(req.GetId()); i != cloudExportNotFound {
		s.data = append(s.data[:i], s.data[i+1:]...)
		return &cloudexportpb.DeleteCloudExportResponse{}, nil
	}
	return nil, status.Errorf(codes.NotFound, "cloud export of id %q doesn't exists", req.GetId())
}

func (s *testAPIServer) allocateNewID() (string, error) {
	var id int

	for _, item := range s.data {
		itemID, err := strconv.Atoi(item.Id)
		if err != nil {
			return "", status.Errorf(codes.Internal, "cannot allocate ID to new cloud export, "+
				"%q string is not a valid integer", item.Id)
		}
		if itemID > id {
			id = itemID
		}
	}
	return strconv.FormatInt(int64(id)+1, 10), nil
}

func (s *testAPIServer) findByName(name string) int {
	for i, ce := range s.data {
		if ce.Name == name {
			return i
		}
	}
	return cloudExportNotFound
}

func (s *testAPIServer) findByID(id string) int {
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
