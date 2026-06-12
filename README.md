# MiniBank

Core bancário didático porém realista, construído em Go + Angular sobre AWS. O objetivo é entregar features completas de produto — ponta a ponta, com segurança e observabilidade — ao longo de 9 fases incrementais.

## O que é

MiniBank é um sistema bancário enxuto e **correto**: cada movimentação financeira é registrada como par débito/crédito num ledger double-entry (os lançamentos sempre somam zero), o saldo é **derivado** desses lançamentos — nunca um campo mutável solto — e valores são sempre em centavos (`int64`), nunca `float`. Dados sensíveis (CPF, PII) são tratados com cifragem em repouso via KMS.

O projeto cobre desde cadastro do cliente até backoffice administrativo, com RBAC, MFA, trilha de auditoria por correlation ID e notificações em tempo real via WebSocket.

## Stack

| Camada | Tecnologias |
|--------|-------------|
| **Backend** | Go 1.25 · `net/http` (Gin/Echo) · `pgx` + `sqlc` · `golang-migrate` · JWT + bcrypt + TOTP · `log/slog` · `testcontainers-go` |
| **Frontend** | Angular 21 standalone · Signals · Reactive Forms · RxJS · `HttpClient` com interceptor JWT |
| **Infra** | AWS: ECS Fargate · ALB · RDS (Postgres) · ElastiCache (Redis) · S3 + CloudFront · Lambda · Secrets Manager · ECR |
| **IaC / CI** | Terraform · GitHub Actions |
| **Observabilidade** | `slog` JSON + correlation ID · OpenTelemetry (traces/métricas) · CloudWatch |

## Estrutura atual

```
minibank/
├── docker-compose.yml      # API Go + Postgres (ambiente local)
├── Makefile                # atalhos de format/lint (back + front)
├── .env.example            # variáveis de ambiente do Postgres
├── .vscode/                # format-on-save e extensões recomendadas
├── backend/                # módulo Go (module backend)
│   ├── cmd/api/            # entrada da API (net/http): / e /health
│   ├── internal/           # domínios: account, auth, transaction
│   ├── Dockerfile          # build multi-stage (imagem final enxuta)
│   ├── .golangci.yml       # config do golangci-lint
│   └── go.mod
└── frontend/               # Angular 21 standalone + SSR
    ├── src/app/core/       # models e services transversais (ApiService)
    ├── src/app/shared/     # models de domínio (Account, Transaction)
    ├── src/app/features/   # telas por feature (a partir da F3)
    ├── src/environments/   # config dev/prod (baseUrl da API)
    └── eslint.config.js
```

## Roadmap (F0 → F8)

| Fase | Tema | Status |
|------|------|--------|
| **F0** | Fundamentos, setup e toolchain (Go + Angular + Git + Docker local) | ✅ concluída |
| **F1** | Core Domain Backend: contas, ledger double-entry e saldo | 🔄 em andamento |
| **F2** | API REST, autenticação (JWT) e autorização (RBAC + MFA) | ⏳ |
| **F3** | Frontend Angular: telas, forms, services e estado reativo | ⏳ |
| **F4** | Integração full-stack: transferência + Pix-like + tempo real | ⏳ |
| **F5** | Infra AWS, containerização e CI/CD | ⏳ |
| **F6** | Segurança / DevSecOps: hardening, OWASP, threat modeling | ⏳ |
| **F7** | Observabilidade, resiliência e produção | ⏳ |
| **F8** | Backoffice, RBAC e auditoria | ⏳ |

> Detalhes de cada fase, Definition of Done e matriz de trilhas: [mentorship-juansavedra](../mentorship-juansavedra/README.md).

## Módulos de domínio (visão futura)

| Módulo | Responsabilidade |
|--------|-----------------|
| `customer` | Cadastro de clientes e KYC básico; CPF cifrado via KMS |
| `account` | Contas bancárias (1:N por cliente); ciclo de vida e limites |
| `ledger` | Razão contábil double-entry; invariantes por transação ACID |
| `transfer` | Transferências internas e Pix-like; idempotência via `Idempotency-Key` |
| `statement` | Extrato paginado; comprovante PDF no S3; push via WebSocket |
| `auth` | JWT access+refresh; bcrypt; TOTP; RBAC (CLIENTE / OPERADOR / ADMIN) |
| `notification` | Push em tempo real ao Angular; jobs agendados via Lambda |
| `backoffice` | Área administrativa com maker-checker e trilha de auditoria |

## Princípios não-negociáveis

- **Consistência de saldo** — derivado do ledger, atualizado na mesma transação ACID, protegido por `CHECK` e lock pessimista/otimista.
- **Idempotência** — `Idempotency-Key` persistida; retry não gera lançamento duplicado.
- **Concorrência segura** — suíte rodada com `-race`; lock no banco é a fonte da verdade.
- **Segurança shift-left** — CIA Triad, Least Privilege, STRIDE no fluxo de transferência, segredos fora do código.
- **Observabilidade nativa** — logs JSON com correlation ID + OpenTelemetry desde o primeiro endpoint.

## Como rodar (F0)

### Ambiente local com Docker (um comando)

Sobe a API Go + Postgres com reprodutibilidade:

```bash
docker compose up --build
```

A API responde em `http://localhost:8080/health`. O Postgres fica em `localhost:5432` e persiste os dados no volume `pgdata` entre restarts.

Credenciais vêm de variáveis de ambiente, com defaults no `compose`. Para sobrescrever, copie o exemplo e edite:

```bash
cp .env.example .env
```

Derrubar (mantendo os dados): `docker compose down`. Apagar também o volume: `docker compose down -v`.

### Rodar os apps direto (sem container)

```bash
# Backend
cd backend
go run ./cmd/api

# Frontend
cd frontend
npm install
npm start
```

### Qualidade de código (format e lint)

Ferramentas necessárias (instalar uma vez):

```bash
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
# angular-eslint já está no frontend (devDependencies)
```

Comandos via `Makefile` na raiz (rodam backend + frontend):

| Comando | O que faz |
|---------|-----------|
| `make fmt` | Formata: `gofmt` + `goimports` (Go) e `prettier` (front) |
| `make fmt-check` | Falha se houver arquivo Go fora do padrão |
| `make vet` | `go vet ./...` no backend |
| `make lint` | `golangci-lint run` (back) + `ng lint` (front) |
| `make check` | Gate completo: `fmt-check` + `vet` + `lint` |

O VS Code aplica **format-on-save** (config em `.vscode/settings.json`); extensões recomendadas em `.vscode/extensions.json`.
