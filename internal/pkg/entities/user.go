/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package entities

import (
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-user-go"
	"time"
)

// User model the information available regarding a User of an organization
type User struct {
	OrganizationId string `json:"organization_id,omitempty"`
	Email          string `json:"email,omitempty"`
	Name           string `json:"name,omitempty"`
	PhotoUrl       string `json:"photo_url,omitempty"`
	MemberSince    int64  `json:"member_since,omitempty"`
	LastName       string `json:"last_name,omitempty"`
	Title          string `json:"title,omitempty"`
	Phone          string `json:"phone,omitempty"`
	Location       string `json:"location,omitempty"`
}

func NewUserFromGRPC(addUserRequest *grpc_user_go.AddUserRequest) *User {
	return &User{
		OrganizationId: addUserRequest.OrganizationId,
		Email:          addUserRequest.Email,
		Name:           addUserRequest.Name,
		PhotoUrl:       "",
		MemberSince:    time.Now().Unix(),
		Title:          addUserRequest.Title,
	}
}

func (u *User) ToGRPC() *grpc_user_go.User {
	return &grpc_user_go.User{
		OrganizationId: u.OrganizationId,
		Email:          u.Email,
		Name:           u.Name,
		PhotoUrl:       u.PhotoUrl,
		MemberSince:    u.MemberSince,
		ContactInfo: &grpc_user_go.ContactInfo{
			FullName: fmt.Sprintf("%s %s", u.Name, u.LastName),
			Address:  u.Location,
			Title:    u.Title,
		},
	}
}

func (u *User) ApplyUpdate(request *grpc_user_go.UpdateUserRequest) {
	if request.Name != "" {
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

func ValidUpdateUserRequest(request *grpc_user_go.UpdateUserRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.Email == "" {
		return derrors.NewInvalidArgumentError(emptyEmail)
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
