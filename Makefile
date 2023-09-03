# Makefile for Kafka Docker Compose project

# Define the Docker Compose file
DOCKER_COMPOSE_FILE := docker-compose.yml

# Targets
.PHONY: start
start:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

.PHONY: stop
stop:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

.PHONY: restart
restart: stop start

.PHONY: clean
clean: stop
	docker-compose -f $(DOCKER_COMPOSE_FILE) rm -f

# to start kafka run in "kafka origin directory 	cmd"/./bin/Windows/zookeeper-server-start.sh config/zookeeper.properties

