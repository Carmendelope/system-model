/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package entities

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-user-go"
	"time"
)

// -------
// User
// -------

type UserContactInfo struct {
	FullName 		string `json:"full_name,omitempty" cql:"full_name"`
	Address 		string `json:"address,omitempty" cql:"address"`
	Phone 			map[string]string `json:"phone,omitempty" cql:"phone"`
	AltEmail		string `json:"alt_email,omitempty" cql:"alt_email"`
	CompanyName		string `json:"company_name,omitempty" cql:"company_name"`
	Title 			string `json:"title,omitempty" cql:"title"`

}
func NewContactInfoFromGRPC (full_name string, address string, phone map[string]string,
	alt_email string, company_name string, title string) *UserContactInfo{
	return &UserContactInfo{
		FullName: 	full_name,
		Address: 	address,
		Phone: 		phone,
		AltEmail:	alt_email,
		CompanyName:company_name,
		Title:		title,
	}
}
/*
func (u *UserContactInfo) ToGRPC() *grpc_user_go.ContactInfo{
	if u == nil {
		return nil
	}
	return &grpc_user_go.ContactInfo{
		FullName: u.FullName,
		Address : u.Address,
		Phone: u.Phone,
		AltEmail: u.AltEmail,
		CompanyName: u.CompanyName,
		Title: u.Title,
	}
}
*/
// User model the information available regarding a User of an organization
type User struct {
	OrganizationId       string   `json:"organization_id,omitempty"` //Deprecated, it will be deleted in version 0.5.0
	Email                string   `json:"email,omitempty"`
	Name                 string   `json:"name,omitempty"`
	PhotoUrl             string   `json:"photo_url,omitempty"`
	MemberSince          int64    `json:"member_since,omitempty"`
	ContactInfo			 UserContactInfo `json:"contact_info,omitempty"`
}

func NewUserFromGRPC(addUserRequest *grpc_user_go.AddUserRequest) * User{
	return &User{
		OrganizationId: addUserRequest.OrganizationId,
		Email:          addUserRequest.Email,
		Name:           addUserRequest.Name,
		PhotoUrl:       "",
		MemberSince:    time.Now().Unix(),
		// TODO: add contactInfo
	}
}

func (u * User) ToGRPC() * grpc_user_go.User {
	return &grpc_user_go.User{
		OrganizationId:       u.OrganizationId,
		Email:                u.Email,
		Name:                 u.Name,
		PhotoUrl:             u.PhotoUrl,
		MemberSince:          u.MemberSince,
		// TODO: add contactInfo
	}
}

func (u * User) ApplyUpdate(request * grpc_user_go.UpdateUserRequest) {
	if request.Name != ""{
		u.Name = request.Name
	}
}

func ValidUserID(userID *grpc_user_go.UserId) derrors.Error {

	// OrganizationID is deprecated,
	// for compatibility of the two versions, we delete this check
	//if userID.OrganizationId == "" {
	//	return derrors.NewInvalidArgumentError(emptyOrganizationId)
	//}
	if userID.Email == "" {
		return derrors.NewInvalidArgumentError(emptyEmail)
	}
	return nil
}

func ValidAddUserRequest(addUserRequest *grpc_user_go.AddUserRequest) derrors.Error {
	// OrganizationID is deprecated,
	// for compatibility of the two versions, we delete this check
	//if addUserRequest.OrganizationId == "" {
	//	return derrors.NewInvalidArgumentError(emptyOrganizationId)
	//}
	if addUserRequest.Email == "" {
		return derrors.NewInvalidArgumentError(emptyEmail)
	}
	if addUserRequest.Name == "" {
		return derrors.NewInvalidArgumentError(emptyName)
	}
	return nil
}

func ValidUpdateUserRequest(request *grpc_user_go.UpdateUserRequest) derrors.Error {
	// OrganizationID is deprecated,
	// for compatibility of the two versions, we delete this check
	//if request.OrganizationId == "" {
	//	return derrors.NewInvalidArgumentError(emptyOrganizationId)
	//}
	if request.Email == "" {
		return derrors.NewInvalidArgumentError(emptyEmail)
	}
	return nil
}

func ValidRemoveUserRequest(removeRequest *grpc_user_go.RemoveUserRequest) derrors.Error {
	// OrganizationID is deprecated,
	// for compatibility of the two versions, we delete this check
	//if removeRequest.OrganizationId == "" {
	//	return derrors.NewInvalidArgumentError(emptyOrganizationId)
	//}
	if removeRequest.Email == "" {
		return derrors.NewInvalidArgumentError(emptyEmail)
	}
	return nil
}

// --------------
// Account User
// --------------

// AccountUser message with the information of the status of a user in an account
type AccountUser struct{
	AccountId 	string `json:"account_id,omitempty"`
	Email 		string `json:"email,omitempty"`
	RoleId  	string `json:"role_id,omitempty"`
	Internal 	bool   `json:"internal,omitempty"`
	Status 		int    `json:"status,omitempty"`
}

// ---------------------
// User Account Invite
// ---------------------


type AccountUserInvite struct{
	AccountId 	string `json:"account_id,omitempty"`
	Email 		string `json:"email,omitempty"`
	RoleId  	string `json:"role_id,omitempty"`
	InvitedBy 	string `json:"invited_by,omitempty"`
 	Msg 		string `json:"msg,omitempty"`
	Expires 	int64  `json:"expires,omitempty"`
}

// --------------
// User Project
// --------------
type ProjectUser struct{
	AccountId 	string `json:"account_id,omitempty"`
	ProjectId 	string `json:"project_id,omitempty"`
	Email 		string `json:"email,omitempty"`
	RoleId 		string `json:"role_id,omitempty"`
	Internal 	bool   `json:"internal,omitempty"`
	Status 		int    `json:"status,omitempty"`
}
// ---------------------
// User Project Invite
// ---------------------

type ProjectUserInvite struct{
	AccountId 	string `json:"account_id,omitempty"`
	ProjectId 	string `json:"project_id,omitempty"`
	Email 		string `json:"email,omitempty"`
	RoleId 		string `json:"role_id,omitempty"`
	InvitedBy	string `json:"invited_by,omitempty"`
	Msg 		string `json:"msg,omitempty"`
	Expires 	int64  `json:"expires,omitempty"`
}
