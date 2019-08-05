/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

/*

 Commands to execute before launching the tests
------------------------------------------------
docker run --name scylla -p 9042:9042 -d scylladb/scylla
docker exec -it scylla cqlsh

create KEYSPACE nalej WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
create type IF NOT EXISTS nalej.account_billing_info (account_id text, full_name text, company_name text, address text, additional_info text);
create table IF NOT EXISTS nalej.Account (account_id text, name text, created bigint, billing_info FROZEN<account_billing_info>, state int, state_info text, primary key (account_id) );
create index IF NOT EXISTS accountName on nalej.Account(name);

Environment variables
IT_SCYLLA_HOST=127.0.0.1
RUN_INTEGRATION_TEST=true
IT_NALEJ_KEYSPACE=nalej
IT_SCYLLA_PORT=9042

*/
package account

import (
	"github.com/nalej/system-model/internal/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

var _ = ginkgo.Describe("Scylla account provider", func(){

	if ! utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	var scyllaHost = os.Getenv("IT_SCYLLA_HOST")
	if scyllaHost == "" {
		ginkgo.Fail("missing environment variables")
	}
	var nalejKeySpace = os.Getenv("IT_NALEJ_KEYSPACE")
	if scyllaHost == "" {
		ginkgo.Fail("missing environment variables")
	}
	scyllaPort, _ := strconv.Atoi(os.Getenv("IT_SCYLLA_PORT"))
	if scyllaPort <= 0 {
		ginkgo.Fail("missing environment variables")
	}

	// create a provider and connect it
	provider := NewScyllaAccountProvider(scyllaHost, scyllaPort, nalejKeySpace)

	ginkgo.AfterSuite(func() {
		provider.Disconnect()
	})

	RunTest(provider)

})