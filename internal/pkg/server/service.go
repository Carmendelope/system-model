/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package server

import (
	"fmt"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-role-go"
	"github.com/nalej/grpc-user-go"
	"github.com/nalej/system-model/internal/pkg/server/asset"
	"github.com/nalej/system-model/internal/pkg/server/cluster"
	"github.com/nalej/system-model/internal/pkg/server/device"
	"github.com/nalej/system-model/internal/pkg/server/eic"
	"github.com/nalej/system-model/internal/pkg/server/node"
	"github.com/nalej/system-model/internal/pkg/server/role"
	"github.com/nalej/system-model/internal/pkg/server/user"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/tools"
	appProvider "github.com/nalej/system-model/internal/pkg/provider/application"
	clusterProvider "github.com/nalej/system-model/internal/pkg/provider/cluster"
	devProvider "github.com/nalej/system-model/internal/pkg/provider/device"
	nodeProvider "github.com/nalej/system-model/internal/pkg/provider/node"
	orgProvider "github.com/nalej/system-model/internal/pkg/provider/organization"
	rProvider "github.com/nalej/system-model/internal/pkg/provider/role"
	uProvider "github.com/nalej/system-model/internal/pkg/provider/user"
	aProvider "github.com/nalej/system-model/internal/pkg/provider/asset"
	eicProvider "github.com/nalej/system-model/internal/pkg/provider/eic"
	"github.com/nalej/system-model/internal/pkg/server/application"
	"github.com/nalej/system-model/internal/pkg/server/organization"
)

// Service structure containing the configuration and gRPC server.
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

// Providers structure with all the providers in the system.
type Providers struct {
	organizationProvider orgProvider.Provider
	clusterProvider clusterProvider.Provider
	nodeProvider nodeProvider.Provider
	applicationProvider appProvider.Provider
	roleProvider rProvider.Provider
	userProvider uProvider.Provider
	deviceProvider devProvider.Provider
	assetProvider aProvider.Provider
	controllerProvider eicProvider.Provider
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
		clusterProvider:      clusterProvider.NewMockupClusterProvider(),
		nodeProvider:         nodeProvider.NewMockupNodeProvider(),
		applicationProvider:  appProvider.NewMockupOrganizationProvider(),
		roleProvider:         rProvider.NewMockupRoleProvider(),
		userProvider:         uProvider.NewMockupUserProvider(),
		deviceProvider:       devProvider.NewMockupDeviceProvider(),
		assetProvider:        aProvider.NewMockupAssetProvider(),
		controllerProvider:   eicProvider.NewMockupEICProvider(),
	}
}

// CreateDBScyllaProviders returns a set of in-memory providers.
func (s *Service) CreateDBScyllaProviders() * Providers {
	return &Providers{
		organizationProvider: orgProvider.NewScyllaOrganizationProvider(
			s.Configuration.ScyllaDBAddress, s.Configuration.ScyllaDBPort, s.Configuration.KeySpace),
		clusterProvider: clusterProvider.NewScyllaClusterProvider(
			s.Configuration.ScyllaDBAddress, s.Configuration.ScyllaDBPort, s.Configuration.KeySpace),
		nodeProvider: nodeProvider.NewScyllaNodeProvider(
			s.Configuration.ScyllaDBAddress, s.Configuration.ScyllaDBPort, s.Configuration.KeySpace),
		applicationProvider: appProvider.NewScyllaApplicationProvider(
			s.Configuration.ScyllaDBAddress, s.Configuration.ScyllaDBPort, s.Configuration.KeySpace),
		roleProvider: rProvider.NewSScyllaRoleProvider(
			s.Configuration.ScyllaDBAddress, s.Configuration.ScyllaDBPort, s.Configuration.KeySpace),
		userProvider: uProvider.NewScyllaUserProvider(
			s.Configuration.ScyllaDBAddress, s.Configuration.ScyllaDBPort, s.Configuration.KeySpace),
		deviceProvider: devProvider.NewScyllaDeviceProvider(
			s.Configuration.ScyllaDBAddress, s.Configuration.ScyllaDBPort, s.Configuration.KeySpace),
		assetProvider:      aProvider.NewScyllaAssetProvider(
			s.Configuration.ScyllaDBAddress, s.Configuration.ScyllaDBPort, s.Configuration.KeySpace),
		controllerProvider: eicProvider.NewScyllaControllerProvider(
			s.Configuration.ScyllaDBAddress, s.Configuration.ScyllaDBPort, s.Configuration.KeySpace),
	}
}

// GetProviders builds the providers according to the selected backend.
func (s *Service) GetProviders() * Providers {
	if s.Configuration.UseInMemoryProviders {
		return s.CreateInMemoryProviders()
	} else if s.Configuration.UseDBScyllaProviders {
		return s.CreateDBScyllaProviders()
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
	// nodes
	nodeManager := node.NewManager(p.organizationProvider, p.clusterProvider, p.nodeProvider)
	nodeHandler := node.NewHandler(nodeManager)
	// applications
	appManager := application.NewManager(p.organizationProvider, p.applicationProvider, p.deviceProvider, s.Configuration.PublicHostDomain)
	applicationHandler := application.NewHandler(appManager)
	// roles
	roleManager := role.NewManager(p.organizationProvider, p.roleProvider)
	roleHandler := role.NewHandler(roleManager)
	// users
	userManager := user.NewManager(p.organizationProvider, p.userProvider)
	userHandler := user.NewHandler(userManager)
	//device
	deviceManager := device.NewManager(p.deviceProvider, p.organizationProvider)
	deviceHandler := device.NewHandler(deviceManager)

	assetManager := asset.NewManager(p.organizationProvider, p.assetProvider)
	assetHandler := asset.NewHandler(assetManager)

	controllerManager := eic.NewManager(p.controllerProvider, p.organizationProvider)
	controllerHandler := eic.NewHandler(controllerManager)

	grpcServer := grpc.NewServer()
	grpc_organization_go.RegisterOrganizationsServer(grpcServer, organizationHandler)
	grpc_infrastructure_go.RegisterClustersServer(grpcServer, clusterHandler)
	grpc_infrastructure_go.RegisterNodesServer(grpcServer, nodeHandler)
	grpc_application_go.RegisterApplicationsServer(grpcServer, applicationHandler)
	grpc_role_go.RegisterRolesServer(grpcServer, roleHandler)
	grpc_user_go.RegisterUsersServer(grpcServer, userHandler)
	grpc_device_go.RegisterDevicesServer(grpcServer, deviceHandler)
	grpc_inventory_go.RegisterAssetsServer(grpcServer, assetHandler)
	grpc_inventory_go.RegisterControllersServer(grpcServer, controllerHandler)

	if s.Configuration.Debug{
		log.Info().Msg("Enabling gRPC server reflection")
		// Register reflection service on gRPC server.
		reflection.Register(grpcServer)
	}
	
	log.Info().Int("port", s.Configuration.Port).Msg("Launching gRPC server")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Errs("failed to serve: %v", []error{err})
	}
	return nil
}

