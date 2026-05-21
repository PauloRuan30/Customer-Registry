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