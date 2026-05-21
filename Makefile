.PHONY: up down test build run swag logs clean

# Gera a documentação Swagger e sobe os contentores em background
up: swag
	docker compose up --build -d

# Para e remove os contentores, redes e volumes associados
down:
	docker compose down -v

# Roda todos os testes unitários e de integração
test:
	go test -v ./...

# Gera os ficheiros do Swagger (requer o swag instalado)
swag:
	go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/api/main.go

# Acompanha os logs da aplicação e do banco de dados em tempo real
logs:
	docker compose logs -f

# Limpeza de módulos e compilações locais
clean:
	go mod tidy
	rm -rf docs/