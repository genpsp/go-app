#!/bin/bash

# Repository
rm -rf domain/repository/mock_repositories/*repository.go
for repository in $(ls -F domain/repository/ | grep -v "/" | grep -v "_test.go" | grep -xv "db.go"); do
    mockgen -source domain/repository/${repository} -package mock_repositories -destination domain/repository/mock_repositories/${repository}
done

# Service
rm -rf services/src/services/mock/*.go
for service in $(ls -F services/src/services/ | grep -v "/" | grep -v "_test.go"); do
    mockgen -source services/src/services/${service} -package mock_services -destination services/src/services/mock/${service}
done

# Handler
rm -rf services/src/handler/mock/*.go
for handler in $(ls -F services/src/handler/ | grep -v "/"  | grep -v "_test.go" | grep -xv "handler.go"); do
    mockgen -source services/src/handler/${handler} -package mock_handler -destination services/src/handler/mock/${handler}
done
