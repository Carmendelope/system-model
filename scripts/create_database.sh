#!/bin/bash
echo "Creating database..."
echo "docker exec -i scylla cqlsh < ./database.cql"
docker exec -i scylla cqlsh < ./database.cql
echo "Done!"
