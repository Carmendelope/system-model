/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package server

import (
	"fmt"
	"github.com/nalej/grpc-account-go"
	"github.com/nalej/grpc-application-network-go"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-project-go"
	"github.com/nalej/grpc-role-go"
	"github.com/nalej/grpc-user-go"
	"github.com/nalej/system-model/internal/pkg/server/application_network"
	"github.com/nalej/system-model/internal/pkg/server/project"

	"github.com/nalej/system-model/internal/pkg/server/account"
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
	acProvider "github.com/nalej/system-model/internal/pkg/provider/account"
	appProvider "github.com/nalej/system-model/internal/pkg/provider/application"
	anProvider "github.com/nalej/system-model/internal/pkg/provider/application_network"
	aProvider "github.com/nalej/system-model/internal/pkg/provider/asset"
	clusterProvider "github.com/nalej/system-model/internal/pkg/provider/cluster"
	devProvider "github.com/nalej/system-model/internal/pkg/provider/device"
	eicProvider "github.com/nalej/system-model/internal/pkg/provider/eic"
	nodeProvider "github.com/nalej/system-model/internal/pkg/provider/node"
	orgProvider "github.com/nalej/system-model/internal/pkg/provider/organization"
	pProvider "github.com/nalej/system-model/internal/pkg/provider/project"
	rProvider "github.com/nalej/system-model/internal/pkg/provider/role"
	uProvider "github.com/nalej/system-model/internal/pkg/provider/user"

	"github.com/nalej/system-model/internal/pkg/server/application"
	"github.com/nalej/system-model/internal/pkg/server/organization"
)

// Service structure containing the configuration and gRPC server.
type Service struct {
	Configuration Config
}

// NewService creates a new system model service.
func NewService(conf Config) *Service {
	return &Service{
		conf,
	}
}

// Providers structure with all the providers in the system.
type Providers struct {
	organizationProvider orgProvider.Provider
	clusterProvider      clusterProvider.Provider
	nodeProvider         nodeProvider.Provider
	applicationProvider  appProvider.Provider
	roleProvider         rProvider.Provider
	userProvider         uProvider.Provider
	deviceProvider       devProvider.Provider
	assetProvider        aProvider.Provider
	controllerProvider   eicProvider.Provider
	accountProvider      acProvider.Provider
	projectProvider      pProvider.Provider
	appNetProvider       anProvider.Provider
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
func (s *Service) CreateInMemoryProviders() *Providers {
	return &Providers{
		organizationProvider: orgProvider.NewMockupOrganizationProvider(),
		clusterProvider:      clusterProvider.NewMockupClusterProvider(),
		nodeProvider:         nodeProvider.NewMockupNodeProvider(),
		applicationProvider:  appProvider.NewMockupApplicationProvider(),
		roleProvider:         rProvider.NewMockupRoleProvider(),
		userProvider:         uProvider.NewMockupUserProvider(),
		deviceProvider:       devProvider.NewMockupDeviceProvider(),
		assetProvider:        aProvider.NewMockupAssetProvider(),
		controllerProvider:   eicProvider.NewMockupEICProvider(),
		accountProvider:      acProvider.NewMockupAccountProvider(),
		projectProvider:      pProvider.NewMockupProjectProvider(),
		appNetProvider:       anProvider.NewMockupApplicationNetworkProvider(),
	}
}

// CreateDBScyllaProviders returns a set of in-memory providers.
func (s *Service) CreateDBScyllaProviders() *Providers {
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
		assetProvider: aProvider.NewScyllaAssetProvider(
			s.Configuration.ScyllaDBAddress, s.Configuration.ScyllaDBPort, s.Configuration.KeySpace),
		controllerProvider: eicProvider.NewScyllaControllerProvider(
			s.Configuration.ScyllaDBAddress, s.Configuration.ScyllaDBPort, s.Configuration.KeySpace),
		accountProvider: acProvider.NewScyllaAccountProvider(
			s.Configuration.ScyllaDBAddress, s.Configuration.ScyllaDBPort, s.Configuration.KeySpace),
		projectProvider: pProvider.NewScyllaProjectProvider(
			s.Configuration.ScyllaDBAddress, s.Configuration.ScyllaDBPort, s.Configuration.KeySpace),
		appNetProvider: anProvider.NewScyllaApplicationNetworkProvider(
			s.Configuration.ScyllaDBAddress, s.Configuration.ScyllaDBPort, s.Configuration.KeySpace),
	}
}

// GetProviders builds the providers according to the selected backend.
func (s *Service) GetProviders() *Providers {
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

	appNetManager := application_network.NewManager(p.organizationProvider, p.applicationProvider, p.appNetProvider)
	appNetHandler := application_network.NewHandler(appNetManager)

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
	//account
	accountManager := account.NewManager(p.accountProvider)
	accountHandler := account.NewHandler(accountManager)
	//project
	projectManager := project.NewManager(p.accountProvider, p.projectProvider)
	projectHandler := project.NewHandler(projectManager)

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
	grpc_account_go.RegisterAccountsServer(grpcServer, accountHandler)
	grpc_project_go.RegisterProjectsServer(grpcServer, projectHandler)
	grpc_application_network_go.RegisterApplicationNetworkServer(grpcServer, appNetHandler)

	if s.Configuration.Debug {
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
