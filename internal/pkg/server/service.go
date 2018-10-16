/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package server

import (
	"fmt"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/system-model/internal/pkg/server/cluster"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/tools"
	orgProvider "github.com/nalej/system-model/internal/pkg/provider/organization"
	clusterProvider "github.com/nalej/system-model/internal/pkg/provider/cluster"
	appProvider "github.com/nalej/system-model/internal/pkg/provider/application"

	"github.com/nalej/system-model/internal/pkg/server/organization"
	"github.com/nalej/system-model/internal/pkg/server/application"
	"github.com/nalej/grpc-application-go"
)

type Service struct {
	Configuration Config
	Server * tools.GenericGRPCServer
}

// NewService creates a new system model service.
func NewService(conf Config) *Service {
	return &Service{
		conf,
		tools.NewGenericGRPCServer(uint32(conf.Port)),
	}
}

type Providers struct {
	organizationProvider orgProvider.Provider
	clusterProvider clusterProvider.Provider
	applicationProvider appProvider.Provider
}

// Name of the service.
func (s *Service) Name() string {
	return "System Model Service."
}

// Description of the service.
func (s *Service) Description() string {
	return "Api service of the System Model project."
}

// CreateInMemoryProviders returns a set of in-memory providers.
func (s *Service) CreateInMemoryProviders() * Providers {
	return &Providers{
		organizationProvider: orgProvider.NewMockupOrganizationProvider(),
		clusterProvider: clusterProvider.NewMockupClusterProvider(),
		applicationProvider: appProvider.NewMockupOrganizationProvider(),
	}
}

// GetProviders builds the providers according to the selected backend.
func (s *Service) GetProviders() * Providers {
	if s.Configuration.UseInMemoryProviders {
		return s.CreateInMemoryProviders()
	}
	log.Fatal().Msg("unsupported type of provider")
	return nil
}

// Run the service, launch the REST service handler.
func (s *Service) Run() error {
	s.Configuration.Print()
	p := s.GetProviders()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Configuration.Port))
	if err != nil {
		log.Fatal().Errs("failed to listen: %v", []error{err})
	}
	// organizations
	orgManager := organization.NewManager(p.organizationProvider)
	organizationHandler := organization.NewHandler(orgManager)
	// clusters
	clusterManager := cluster.NewManager(p.organizationProvider, p.clusterProvider)
	clusterHandler := cluster.NewHandler(clusterManager)
	// applications
	appManager := application.NewManager(p.organizationProvider, p.applicationProvider)
	applicationHandler := application.NewHandler(appManager)


	grpcServer := grpc.NewServer()
	grpc_organization_go.RegisterOrganizationsServer(grpcServer, organizationHandler)
	grpc_infrastructure_go.RegisterClustersServer(grpcServer, clusterHandler)
	grpc_application_go.RegisterApplicationsServer(grpcServer, applicationHandler)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	log.Info().Int("port", s.Configuration.Port).Msg("Launching gRPC server")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Errs("failed to serve: %v", []error{err})
	}
	return nil
}

