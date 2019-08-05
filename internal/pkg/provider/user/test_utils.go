package user

import (
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
	"time"
)

func CreateUser (email string) *entities.User {
	return &entities.User{
		Email: email,
		Name: "User test",
		PhotoUrl: "url",
		MemberSince: time.Now().Unix(),
		ContactInfo: entities.UserContactInfo{
			FullName: "full name",
			Address: "address",
			Phone: map[string]string{"home":"00.000.00.00", "mobile":"000.00.00.00"},
			AltEmail:"alt_email@nalej.com",
			CompanyName:"company",
			Title:"title",
		},
	}
}

func CreateAddUser(provider Provider) (*entities.User, derrors.Error) {

	toAdd := CreateUser(fmt.Sprintf("%s@nalej.com", entities.GenerateUUID()),)
	err := provider.Add(*toAdd)
	if err != nil {
		return nil, err
	}
	return toAdd, err
}

func AddUser(provider Provider, user *entities.User) (derrors.Error) {

	err := provider.Add(*user)
	if err != nil {
		return  err
	}
	return  err
}

func CreateAccountUser(email string) *entities.AccountUser {
	return &entities.AccountUser{
		AccountId: entities.GenerateUUID(),
		Email: email,
		RoleId: entities.GenerateUUID(),
		Internal: true,
		Status: 1,
	}
}

func CreateAccountUserInvite(email string) *entities.AccountUserInvite {
	return &entities.AccountUserInvite{
		AccountId: entities.GenerateUUID(),
		Email: email,
		RoleId: entities.GenerateUUID(),
		InvitedBy: "alt_email@nalej.com",
		Msg: "You are invited to operate in this account",
		Expires:time.Now().Unix(),
	}
}

func CreateProjectUser(accountID string, projectID string, email string) *entities.ProjectUser {
	return &entities.ProjectUser{
		AccountId: accountID,
		ProjectId: projectID,
		Email: email,
		RoleId: entities.GenerateUUID(),
		Internal: true,
		Status: 1,
	}
}

func generateMail() string {
	return fmt.Sprintf("%s.nalej.com", entities.GenerateUUID())
}

func CreateProjectUserInvite(accountID string, projectID string, email string) *entities.ProjectUserInvite {
	return &entities.ProjectUserInvite{
		AccountId: accountID,
		ProjectId: projectID,
		Email:     email,
		RoleId:    entities.GenerateUUID(),
		InvitedBy: "alt_email@nalej.com",
		Msg:       "You are invited to see this project",
		Expires:   time.Now().Unix(),
	}
}