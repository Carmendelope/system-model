/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package entities

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-user-go"
	"time"
)

// User model the information available regarding a User of an organization
type User struct {
	OrganizationId       string   `json:"organization_id,omitempty"`
	Email                string   `json:"email,omitempty"`
	Name                 string   `json:"name,omitempty"`
	PhotoUrl             string   `json:"photo_url,omitempty"`
	MemberSince          int64    `json:"member_since,omitempty"`
}

func NewUserFromGRPC(addUserRequest *grpc_user_go.AddUserRequest) * User{
	return &User{
		OrganizationId: addUserRequest.OrganizationId,
		Email:          addUserRequest.Email,
		Name:           addUserRequest.Name,
		PhotoUrl:       "",
		MemberSince:    time.Now().Unix(),
	}
}

func (u * User) ToGRPC() * grpc_user_go.User {
	return &grpc_user_go.User{
		OrganizationId:       u.OrganizationId,
		Email:                u.Email,
		Name:                 u.Name,
		PhotoUrl:             u.PhotoUrl,
		MemberSince:          u.MemberSince,
	}
}

func (u * User) ApplyUpdate(request * grpc_user_go.UpdateUserRequest) {
	if request.Name != ""{
		u.Name = request.Name
	}
}

func ValidUserID(userID *grpc_user_go.UserId) derrors.Error {
	if userID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if userID.Email == "" {
		return derrors.NewInvalidArgumentError(emptyEmail)
	}
	return nil
}

func ValidAddUserRequest(addUserRequest *grpc_user_go.AddUserRequest) derrors.Error {
	if addUserRequest.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if addUserRequest.Email == "" {
		return derrors.NewInvalidArgumentError(emptyEmail)
	}
	if addUserRequest.Name == "" {
		return derrors.NewInvalidArgumentError(emptyName)
	}
	return nil
}

func ValidRemoveUserRequest(removeRequest *grpc_user_go.RemoveUserRequest) derrors.Error {
	if removeRequest.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if removeRequest.Email == "" {
		return derrors.NewInvalidArgumentError(emptyEmail)
	}
	return nil
}
