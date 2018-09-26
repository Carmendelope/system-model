/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package server

import (
	"github.com/rs/zerolog/log"
	"github.com/nalej/grpc-utils/pkg/tools"
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
	return &Providers{}
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
	//p := s.GetProviders()
	return nil
}

