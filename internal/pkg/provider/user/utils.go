package user

import (
	"github.com/nalej/derrors"
)

func clearScyllaUsers(address string, keyspace string) derrors.Error{

	sp := NewScyllaUserProvider(address, keyspace)

	err := sp.Connect()
	if err == nil {

	}

	sp.ClearTable()

	return nil

}
