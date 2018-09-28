/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package application

import "github.com/nalej/system-model/internal/pkg/provider/application"

type Manager struct {
	provider application.Provider
}

func NewManager(provider application.Provider) Manager {
	return Manager{provider}
}

