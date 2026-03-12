# Makefile для удобного запуска тестов

# Переменные
GO_TEST = go test
PKGS = ./...

# 1. Запуск всех тестов
.PHONY: test
test:
	$(GO_TEST) $(PKGS)

# 2. Подробный вывод тестов
.PHONY: test-verbose
test-verbose:
	$(GO_TEST) -v $(PKGS)

# 3. Покрытие кода с выводом процента
.PHONY: test-cover
test-cover:
	$(GO_TEST) -cover -coverprofile=coverage.out $(PKGS)
	@echo "Coverage report saved to coverage.out"

# 4. Генерация HTML отчета о покрытии
.PHONY: test-cover-html
test-cover-html:
	$(GO_TEST) -coverprofile=coverage.out $(PKGS)
	go tool cover -html=coverage.out -o coverage.html
	@echo "HTML coverage report generated at coverage.html"
