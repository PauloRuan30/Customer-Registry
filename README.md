# Customer Registry API

Uma API REST desenvolvida em Go para gerenciar registros fictícios de clientes. Desenvolvida como parte de um desafio técnico para demonstrar os fundamentos de arquitetura de software, persistência relacional e conteinerização.

## 🛠 Decisões Técnicas

- **Arquitetura (Handler -> Service -> Repository)**: Esta separação garante que a lógica de negócio (Service) não conheça detalhes do HTTP (Handler) ou do Banco de Dados (Repository). Isso facilita os testes unitários.
- **Roteamento (`go-chi/chi`)**: O Chi foi escolhido por ser 100% compatível com a biblioteca padrão `net/http` do Go, além de ser leve e performático.
- **Banco de Dados (`pgx/v5`)**: Utilizado no lugar de ORMs pesados (como GORM) para demonstrar proficiência com SQL puro e obter máxima performance na comunicação com o PostgreSQL.
- **Migrations (`golang-migrate`)**: As migrations são executadas via container no Docker Compose, garantindo que o banco suba sempre estruturado e pronto para uso, sem intervenção manual.

## 🚀 Pré-requisitos

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- Make (opcional, mas recomendado)

## 📦 Como executar o projeto

Basta rodar o comando abaixo na raiz do projeto. Ele fará o build da aplicação Go, subirá o banco PostgreSQL e executará as migrations automaticamente:

```bash
make up
```

ou:
 
```bash
docker-compose up --build -d
```