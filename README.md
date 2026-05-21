
# Customer Registry API

Uma API REST de alta performance desenvolvida em Go para gerir registos de clientes fictícios. Este projeto foi construído rigorosamente como resposta a um desafio técnico, enfatizando a clareza da arquitetura em camadas, a validação estrita de dados de negócio, a conteinerização completa e a observabilidade avançada.

## 🚀 Tecnologias e Ferramentas

- **Linguagem:** Go 1.26
- **Roteamento HTTP:** `go-chi/chi/v5` (Altamente performático e compatível com `net/http`)
- **Base de Dados:** PostgreSQL 15 (`pgx/v5` nativo, sem sobrecarga de ORMs)
- **Migrações:** `golang-migrate` integrado ao ciclo de vida do Docker
- **Infraestrutura:** Docker & Docker Compose
- **Documentação da API:** Swagger / OpenAPI 2.0 (`swaggo/http-swagger`)
- **Observabilidade:** Logs Estruturados em JSON (`log/slog`) e APM Tracer (Datadog Agent)
- **CI/CD:** Pipeline automatizada por ramos via GitHub Actions

## 🏗 Arquitetura do Sistema

O projeto adota uma arquitetura desacoplada em camadas claras, isolando responsabilidades de forma a simplificar a manutenção e a escrita de testes:
- **Handler (Transporte):** Processa requisições HTTP, valida a integridade do JSON de entrada, mapeia códigos de estado HTTP adequados baseados em erros de domínio e expõe o Swagger.
- **Service (Domínio/Negócio):** Concentra as regras de negócio cruciais, validações de integridade (scores, documentos com prefixo exigido, estados operacionais) e fluxos lógicos.
- **Repository (Persistência):** Camada puramente focada na execução de queries SQL otimizadas e interações com o pool de conexões do PostgreSQL.

## 📦 Pré-requisitos

- [Docker](https://docs.docker.com/get-docker/) instalado localmente.
- [Docker Compose](https://docs.docker.com/compose/install/) instalado localmente.
- Utilitário `make` (Opcional, mas altamente recomendado para atalhos).

---

## 🔧 Configuração e Execução Rápida

### 1. Configurar as Variáveis de Ambiente
Crie uma cópia do ficheiro de exemplo `.env.example` e renomeie-a para `.env`:
```bash
cp .env.example .env

```

### 2. Subir a Infraestrutura e Aplicação

Execute o comando abaixo para compilar o binário em Go de forma otimizada (multi-stage build), inicializar o cluster PostgreSQL, o Datadog Agent e **executar automaticamente as migrações da base de dados**:

```bash
make up

```

*Nota: Caso não possua o utilitário `make`, execute o comando Docker nativo correspondente:*

```bash
docker compose up --build -d

```

### 3. Verificar o Estado dos Serviços

Garanta que todos os contentores necessários encontram-se em estado `healthy` ou com execução concluída com sucesso (no caso das migrações):

```bash
docker compose ps

```

### 4. Documentação Interativa (Swagger UI)

Com a API em execução, aceda à interface gráfica do Swagger diretamente no seu navegador para explorar e testar interativamente todos os contratos da API:

```
http://localhost:8080/swagger/index.html

```

---

## 🔄 Execução de Migrações da Base de Dados

As migrações de esquemas relacionais são geridas automaticamente através da ferramenta `golang-migrate/migrate`.
O contentor temporário `customer_registry_migrate` aguarda que a base de dados PostgreSQL esteja totalmente pronta para aceitar conexões (via *healthcheck*) e injeta os ficheiros de migração localizados no diretório `./migrations`. Não é necessária qualquer ação manual.

---

## 🧪 Testes Automatizados

O projeto possui uma suite de testes unitários focada em validar rigorosamente o comportamento das regras de negócio na camada de `Service` e o tratamento adequado de erros operacionais.

Para rodar os testes da aplicação e validar os cenários de sucesso e exceções, execute:

```bash
make test

```

*Equivalente nativo:* `go test -v ./...`

---

## 🔍 Exemplos Práticos de Chamadas à API (cURL)

Abaixo constam exemplos reais de interação com os 5 endpoints obrigatórios expostos na porta `8080`.

### 1. Criar um Cliente Fictício (`POST /customers`)

```bash
curl -X POST http://localhost:8080/customers \
  -H "Content-Type: application/json" \
  -d '{
    "document": "FAKE-98765",
    "name": "Arnaldo Silva Simulado",
    "score": 850,
    "risk_level": "LOW",
    "income_range": "4000-6000"
  }'

```

### 2. Listar Clientes com Paginação Desejável (`GET /customers`)

```bash
curl -X GET "http://localhost:8080/customers?page=1&limit=10"

```

### 3. Buscar Cliente por Identificador ID (`GET /customers/{id}`)

```bash
curl -X GET http://localhost:8080/customers/INSIRA-O-UUID-AQUI

```

### 4. Buscar Cliente pelo Documento Fictício (`GET /customers/document/{document}`)

```bash
curl -X GET http://localhost:8080/customers/document/FAKE-98765

```

### 5. Atualizar Apenas o Status Interno (`PATCH /customers/{id}/status`)

```bash
curl -X PATCH http://localhost:8080/customers/INSIRA-O-UUID-AQUI/status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "UNDER_REVIEW"
  }'

```

---

## 💾 Inspeção Direta do Banco de Dados

Para auditar e certificar a persistência física relacional dos dados sem intermediários, aceda diretamente à CLI interativa do PostgreSQL de forma isolada dentro do contentor:

```bash
docker exec -it customer_registry_db psql -U admin -d customer_registry

```

Queries rápidas para validação em entrevista:

```sql
-- Verificar dados persistidos e a conformidade dos campos
SELECT id, document, name, score, risk_level, status FROM customers;

-- Confirmar o funcionamento do histórico temporal (updated_at) após um PATCH
SELECT id, status, created_at, updated_at FROM customers;

-- Encerrar a CLI do psql
\q

```

---

## 🧠 Decisões Técnicas Tomadas

1. **Evitar ORMs Pesados (`pgx/v5` + SQL Puro):** Optou-se por utilizar o driver nativo `pgx` e conexões estruturadas por `pgxpool` em detrimento de ORMs tradicionais (como GORM). Esta abordagem garante controlo absoluto sobre as transações da base de dados, elimina processamentos reflexivos ocultos e reduz drasticamente a latência de execução.
2. **Uso do Roteador `go-chi`:** Escolhido por ser extremamente leve, possuir total compatibilidade com as assinaturas nativas da biblioteca padrão `net/http` e fornecer suporte nativo simplificado para middlewares e sub-rotas.
3. **Validação Desacoplada e Falha Rápida (Fail-Fast):** As regras de validação (como o formato estrito de documento `FAKE-\d+`, score de 0 a 1000 e enums permitidos) foram isoladas no domínio da camada de `Service`. O `Handler` traduz os erros de negócio diretamente para códigos HTTP equivalentes (400 para erros de payload, 404 para entidades em falta, 409 para violações de unicidade de documento).
4. **Logs Estruturados e APM Integrado:** A aplicação faz uso de `log/slog` configurado para saída em JSON no `Stdout`. Isto garante que ferramentas de agregação de logs consigam indexar os metadados das operações de forma eficiente. A inclusão do rastreamento distribuído (Datadog Tracer) demonstra maturidade técnica ao rastrear o ciclo de vida completo de cada transação.

---

## 🔮 O Que Seria Melhorado com Mais Tempo

Se este projeto fosse escalado para um ecossistema de produção de larga escala, os seguintes pontos seriam implementados:

1. **Testes de Integração com Banco Real (Testcontainers):** Substituir/estender os mocks atuais da camada de repositório por testes de integração reais que utilizam instâncias efémeras do Docker em tempo de execução para garantir que as queries SQL nunca quebrem em produção.
2. **Metadados Completos na Paginação:** Enriquecer o corpo da resposta da listagem de clientes ou incluir cabeçalhos HTTP customizados (`X-Total-Count`, `X-Page-Size`) para facilitar o consumo por aplicações de frontend.
3. **Tratamento de Graceful Shutdown:** Capturar sinais do sistema operacional (`SIGINT`, `SIGTERM`) no ficheiro `main.go` para finalizar de forma segura o pool de conexões do banco de dados e interromper o servidor HTTP sem derrubar requisições ativas em trânsito.
4. **Pipeline de CI/CD Avançada por Ramos (Branches):** Configuração de uma esteira automatizada (ex: GitHub Actions) acionada em Pull Requests e pushes para os ramos ;`dev`, `qa` e `main`. O fluxo executaria análises estáticas (linters), testes com validação de concorrência (`-race`) e empacotamento em imagens Docker com tags dinâmicas contextualizadas por ambiente (`dev-latest`, `qa-latest`, `prod-latest`).
