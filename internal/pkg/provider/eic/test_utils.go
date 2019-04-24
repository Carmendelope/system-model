/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package eic

import (
	"fmt"
	"github.com/nalej/system-model/internal/pkg/entities"
	"math/rand"
	"time"
)

func CreateTestEdgeController() *entities.EdgeController{
	id:= rand.Intn(200)
	labels := make (map[string]string, 0)
	size := rand.Intn(10) +1
	for i:=0; i<size; i++{
		labels[fmt.Sprintf("label-%d", i)] = fmt.Sprintf("value-%d", i)
	}
	return &entities.EdgeController{
		OrganizationId:   fmt.Sprintf("organization_%d", id),
		EdgeControllerId: entities.GenerateUUID(),
		Show:             true,
		Created:          time.Now().Unix(),
		Name:             fmt.Sprintf("name_%d", id),
		Labels:           labels,
	}
}