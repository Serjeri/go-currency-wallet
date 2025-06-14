GIT := git
DOCKER_COMPOSE := docker-compose
SLEEP := sleep
TIMEOUT := 5


.PHONY: all
all: wallet exchanger notification


.PHONY: wallet
clone-currency-wallet:
	@echo "Клонируем go-currency-wallet и запускаем Docker Compose..."
	$(GIT) clone https://github.com/Serjeri/go-currency-wallet.git
	cd go-currency-wallet/deployments && $(DOCKER_COMPOSE) up -d --build
	@echo "Ожидание $(TIMEOUT) секунд перед следующим проектом..."
	$(SLEEP) $(TIMEOUT)


.PHONY: exchanger
clone-exchanger: clone-currency-wallet
	@echo "Клонируем gw-exchanger и запускаем Docker Compose..."
	$(GIT) clone https://github.com/Serjeri/gw-exchanger.git
	cd gw-exchanger/docker && $(DOCKER_COMPOSE) up -d --build
	@echo "Ожидание $(TIMEOUT) секунд перед следующим проектом..."
	$(SLEEP) $(TIMEOUT)


.PHONY: notification
clone-notification: clone-exchanger
	@echo "Клонируем gw-notification и запускаем Docker Compose..."
	$(GIT) clone https://github.com/Serjeri/gw-notification.git
	cd gw-notification/docker && $(DOCKER_COMPOSE) up -d --build
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
