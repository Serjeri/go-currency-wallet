GIT := git
DOCKER_COMPOSE := docker-compose -p
SLEEP := sleep
TIMEOUT := 5
NETWORK_NAME := network

.PHONY: all
all: network exchanger notification

.PHONY: network
network:
	@echo "Создаем Docker сеть $(NETWORK_NAME)..."
	docker network create $(NETWORK_NAME) || true

# .PHONY: wallet
# wallet:
# 	@echo "Клонируем go-currency-wallet и запускаем Docker Compose..."
# 	$(GIT) clone https://github.com/Serjeri/go-currency-wallet.git
# 	cd go-currency-wallet && go mod tidy
# 	cd go-currency-wallet/deployments && $(DOCKER_COMPOSE) wallet -f first-compose.yml up -d --build
# 	@echo "Ожидание $(TIMEOUT) секунд перед следующим проектом..."
# 	$(SLEEP) $(TIMEOUT)

.PHONY: wallet
wallet:
	go mod tidy
	go-currency-wallet/deployments && $(DOCKER_COMPOSE) wallet -f first-compose.yml up -d --build
	@echo "Ожидание $(TIMEOUT) секунд перед следующим проектом..."
	$(SLEEP) $(TIMEOUT)

.PHONY: exchanger
exchanger: wallet
	@echo "Клонируем gw-exchanger и запускаем Docker Compose..."
	$(GIT) clone https://github.com/Serjeri/gw-exchanger.git
	cd gw-exchanger && go mod tidy
	cd gw-exchanger/docker && $(DOCKER_COMPOSE) exchanger -f second-compose.yml up -d --build
	@echo "Ожидание $(TIMEOUT) секунд перед следующим проектом..."
	$(SLEEP) $(TIMEOUT)


.PHONY: notification
notification: exchanger
	@echo "Клонируем gw-notification и запускаем Docker Compose..."
	$(GIT) clone https://github.com/Serjeri/gw-notification.git
	cd gw-notification && go mod tidy
	cd gw-notification/docker && $(DOCKER_COMPOSE) notification -f third-compose.yml up -d --build
	@echo "Все проекты успешно развернуты в Docker!"


.PHONY: clean
clean:
	@echo "Останавливаем и удаляем контейнеры..."
	-cd go-currency-wallet/docker && $(DOCKER_COMPOSE) down
	-cd gw-exchanger/docker && $(DOCKER_COMPOSE) down
	-cd gw-notification/docker && $(DOCKER_COMPOSE) down
	@echo "Удаляем клонированные репозитории..."
	rm -rf go-currency-wallet gw-exchanger gw-notification

.PHONY: logs-wallet
logs-wallet:
	cd go-currency-wallet/docker && $(DOCKER_COMPOSE) logs -f

.PHONY: logs-exchanger
logs-exchanger:
	cd gw-exchanger/docker && $(DOCKER_COMPOSE) logs -f

.PHONY: logs-notification
logs-notification:
	cd gw-notification/docker && $(DOCKER_COMPOSE) logs -f
