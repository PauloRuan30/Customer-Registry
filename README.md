# Customer Registry API

Uma API REST robusta desenvolvida em Go para gerir registos de clientes fictícios. Este projeto foi desenvolvido como resposta a um desafio técnico, enfatizando a clareza da arquitetura, validação estrita de dados e observabilidade.

## 🚀 Tecnologias e Ferramentas

- **Linguagem:** Go 1.22
- **Roteamento:** `go-chi/chi`
- **Base de Dados:** PostgreSQL 15 (`pgx/v5`)
- **Infraestrutura:** Docker & Docker Compose
- **Documentação API:** Swagger (OpenAPI 2.0)
- **Observabilidade:** Logs Estruturados em JSON (`log/slog`) e APM Tracer (Datadog)

## 🏗 Arquitetura

O projeto adota uma arquitetura em camadas focada na clara separação de responsabilidades:
- **Handler (Transport):** Lida com o processamento HTTP (respostas, leitura de JSON, Swagger).
- **Service (Domain):** Encapsula as lógicas de negócio cruciais e validações (ex: scores, estados).
- **Repository (Data):** Totalmente focado nas queries e interação transacional com o PostgreSQL.

## 📦 Pré-requisitos

- [Docker](https://docs.docker.com/get-docker/) e [Docker Compose](https://docs.docker.com/compose/install/) instalados na sua máquina local.
- Make (Opcional, atalhos facilitados).

## 🔧 Configuração e Execução Rápida

1. **Configurar as Variáveis de Ambiente:**
   Crie uma cópia do ficheiro `.env.example` e renomeie-a para `.env`.

```bash
   cp .env.example .env
```


2. **Subir os Contentores da Infraestrutura:**
   Execute o comando abaixo para compilar a aplicação Go, inicializar o PostgreSQL e executar as migrações de esquemas de forma totalmente automatizada:

```bash
    make up
```

*Nota: Caso não tenha o utilitário `make` instalado no seu sistema, execute diretamente:*

```bash
    docker compose up --build -d

```

3. **Verificar o Estado dos Serviços:**
Garanta que todos os contentores (API, Base de Dados, Migrações e Datadog Agent) subiram com sucesso:

```bash
    docker compose ps

```


4. **Interagir via Documentação (Swagger UI):**
Com a API ativa, aceda à interface gráfica do Swagger através do seu browser para testar os endpoints em tempo real:

```bash
    http://localhost:8080/swagger/index.html

```



---

## 🧪 Testes Automatizados

O projeto conta com uma suite de testes unitários focada em garantir a robustez das regras de negócio da camada de `Service` (como validações de limites de score, formatos de documento com o prefixo obrigatório `FAKE-` e estados permitidos).

Para executar os testes locais com feedback detalhado no terminal, utilize:

```bash
make test

```

---

## 🔍 Inspeção Direta da Base de Dados

Para validar que os dados enviados através dos endpoints HTTP estão a ser corretamente persistidos no disco do PostgreSQL (e não apenas retidos em memória), pode aceder diretamente ao cliente interativo `psql` dentro do contentor isolado:

```bash
docker exec -it customer_registry_db psql -U admin -d customer_registry

```

Exemplos de comandos úteis para demonstração:

```sql
-- Listar todos os clientes registados e os seus respetivos scores
SELECT id, document, name, score, status FROM customers;

-- Verificar o funcionamento do histórico de modificações (PATCH)
SELECT id, status, updated_at FROM customers;

-- Sair do terminal do PostgreSQL
\q

```