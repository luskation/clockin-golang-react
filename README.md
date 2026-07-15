# ⏱️ Cronos

Sistema full-stack de controle de ponto desenvolvido para a **Comp Júnior** — múltiplas empresas, usuários com papéis distintos (admin/colaborador) e registro de entrada/saída com um clique.

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-1.12-00ADD8?style=flat)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=flat&logo=postgresql&logoColor=white)
![React](https://img.shields.io/badge/React-19-61DAFB?style=flat&logo=react&logoColor=black)
![TypeScript](https://img.shields.io/badge/TypeScript-6-3178C6?style=flat&logo=typescript&logoColor=white)
![Vite](https://img.shields.io/badge/Vite-8-646CFF?style=flat&logo=vite&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=flat&logo=docker&logoColor=white)
![Tests](https://img.shields.io/badge/tests-39%20passing-2ea44f?style=flat)

---

## 📋 Índice

- [Sobre o projeto](#-sobre-o-projeto)
- [Arquitetura](#-arquitetura)
- [Stack](#-stack)
- [Funcionalidades](#-funcionalidades)
- [Estrutura de pastas](#-estrutura-de-pastas)
- [Como rodar](#-como-rodar)
- [API](#-api)
- [Decisões técnicas](#-decisões-técnicas)
- [Qualidade e testes](#-qualidade-e-testes)
- [Problemas enfrentados](#-problemas-enfrentados)
- [Gestão do projeto](#-gestão-do-projeto)
- [Roadmap](#-roadmap)
- [Autor](#-autor)

---

## 🧭 Sobre o projeto

O **Cronos** resolve um problema real de empresas júnior/pequenas equipes: controlar entrada e saída de funcionários sem planilha. Um admin cadastra empresas e usuários; cada colaborador bate o ponto com um clique e o sistema decide sozinho se é entrada ou saída, com base no último registro.

Construído do zero — modelagem, API, autenticação, frontend, testes e containerização — como projeto full-stack de portfólio.

## 🏗️ Arquitetura

```
┌──────────────┐      HTTPS/JSON       ┌──────────────┐      SQL       ┌──────────────┐
│   Frontend   │ ────────────────────► │   Backend    │ ─────────────► │  PostgreSQL  │
│ React + Vite │ ◄──────────────────── │  Go + Gin    │ ◄───────────── │              │
│  (nginx)     │        JWT             │ handler/     │                │              │
└──────────────┘                       │ service/     │                └──────────────┘
                                        │ repository   │
                                        └──────────────┘
```

Backend em camadas: `handler` (HTTP) → `service` (regra de negócio) → `repository` (acesso a dados via SQL puro, sem ORM). Os services dependem de **interfaces** de repositório, não do struct concreto — o que permite testar regra de negócio com repositórios fake, sem precisar de banco.

## 🧰 Stack

| | Backend | Frontend |
|---|---|---|
| Linguagem | Go 1.26 | TypeScript |
| Framework | Gin | React 19 + Vite |
| Dados | PostgreSQL + pgx (sem ORM) | — |
| Auth | JWT (golang-jwt) + bcrypt | Axios com interceptors |
| Roteamento | — | React Router |
| Estilo | — | CSS Modules (design system próprio) |
| Testes | testify | Vitest + Testing Library |
| Infra | Docker multi-stage | Docker multi-stage + nginx |

**Por que essas escolhas:**
- **SQL puro (pgx) em vez de ORM** — em um domínio pequeno (3 entidades), o ganho de produtividade de um ORM não compensa a perda de controle sobre as queries.
- **CSS Modules sem biblioteca de UI** — escopo de estilo por componente, com um design system próprio (`tokens.css`) em vez de depender da estética padrão de uma lib de componentes.
- **JWT stateless** — sem sessão em banco; o token carrega `user_id` e `role`, verificados em middleware.

## ✨ Funcionalidades

**Admin**
- CRUD completo de empresas (com validação de CNPJ mascarado)
- CRUD completo de usuários (papel e empresa vinculada), sem poder se autoexcluir
- Listagens paginadas em toda a aplicação

**Colaborador**
- Bater ponto com um clique — o sistema alterna entrada/saída automaticamente
- Histórico paginado dos próprios registros
- Corrigir/excluir o próprio registro (admin pode em qualquer um)

**Em toda a aplicação**
- Login com JWT, rotas protegidas por autenticação e por papel
- Validação de formulário com erro por campo (não alerta genérico)
- Loading em toda ação assíncrona, toast de sucesso/erro
- Interceptor Axios centraliza expiração de sessão (401 → logout automático)
- Responsivo (mobile/tablet/desktop)

## 📁 Estrutura de pastas

```
.
├── ponto/                    # backend (Go)
│   ├── cmd/                  # entrypoint
│   ├── internal/
│   │   ├── handler/          # HTTP (Gin)
│   │   ├── service/          # regra de negócio (+ testes unitários)
│   │   ├── repository/       # acesso a dados (pgx)
│   │   ├── middleware/       # auth JWT, checagem de papel
│   │   ├── domain/           # entidades
│   │   └── apperr/           # erros tipados
│   ├── db/migrations/        # schema versionado (golang-migrate)
│   ├── requests.http         # collection de teste manual (REST Client)
│   └── Dockerfile
├── ponto-frontend/           # frontend (React + Vite)
│   └── src/
│       ├── pages/            # uma página por rota (+ pages/admin)
│       ├── components/       # Modal, Toast, Pagination, Header, rotas protegidas...
│       ├── services/         # cliente HTTP (interceptors) + auth
│       ├── hooks/, utils/    # useLiveClock, validação, formatação de data
│       ├── styles/           # tokens de design e reset global
│       └── Dockerfile, nginx.conf
├── docs/
│   └── ponto-api.postman_collection.json   # collection Postman/Insomnia
└── docker-compose.yml        # sobe banco + backend + frontend juntos
```

## 🚀 Como rodar

### Opção 1 — full stack com Docker (recomendado)

Na raiz do repositório:

```bash
cp .env.example .env                     # DB_USER / DB_PASSWORD / DB_NAME
cp ponto/.env.example ponto/.env         # mesmas credenciais + DB_HOST=localhost, DB_PORT, JWT_SECRET
cp ponto-frontend/.env.example ponto-frontend/.env   # VITE_API_URL
docker compose up -d --build
```

- Frontend: `localhost:5173`
- Backend: `localhost:8080`
- Postgres: `localhost:5432`

As migrations ainda precisam ser aplicadas manualmente (ver abaixo).

### Opção 2 — backend local (Go fora do container)

```bash
cd ponto
cp .env.example .env
docker compose up -d      # sobe só o Postgres
migrate -path db/migrations -database "postgres://usuario:senha@localhost:5432/banco?sslmode=disable" up
go run ./cmd/main.go
```

O `migrate` é o CLI do [golang-migrate](https://github.com/golang-migrate/migrate) (`go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest`).

### Opção 3 — frontend local

```bash
cd ponto-frontend
npm install
cp .env.example .env      # VITE_API_URL apontando para o backend
npm run dev
```

## 🔌 API

Base: `/api/v1` · Auth via header `Authorization: Bearer <token>`.

| Método | Rota | Auth | Descrição |
|--------|------|------|-----------|
| GET | `/health` | não | Healthcheck |
| POST | `/auth/login` | não | Login (retorna token + usuário) |
| POST | `/auth/forgot-password` | não | Gera token de reset (15 min) |
| POST | `/auth/reset-password` | não | Redefine senha com o token |
| POST | `/companies` | admin | Cria empresa |
| GET | `/companies` | admin | Lista empresas (paginado) |
| GET | `/companies/:id` | admin | Busca empresa |
| PUT | `/companies/:id` | admin | Atualiza empresa |
| DELETE | `/companies/:id` | admin | Remove empresa |
| POST | `/users` | admin | Cria usuário |
| GET | `/users` | admin | Lista usuários (paginado) |
| GET | `/users/:id` | autenticado | Busca usuário |
| PUT | `/users/:id` | autenticado | Atualiza nome/e-mail/papel |
| DELETE | `/users/:id` | admin | Remove usuário (não pode excluir a si mesmo) |
| POST | `/time-entries` | autenticado | Bate ponto (alterna entrada/saída automaticamente) |
| GET | `/time-entries/me` | autenticado | Meu histórico (paginado) |
| GET | `/time-entries` | admin | Histórico geral (paginado) |
| PUT | `/time-entries/:id` | dono ou admin | Corrige um registro |
| DELETE | `/time-entries/:id` | dono ou admin | Remove um registro |

Collections prontas para testar:
- [`docs/ponto-api.postman_collection.json`](./docs/ponto-api.postman_collection.json) — importável no Postman ou no Insomnia (Import/Export → Postman v2.1). Já vem com scripts que capturam token/IDs das respostas automaticamente.
- [`ponto/requests.http`](./ponto/requests.http) — para a extensão REST Client do VS Code.

## 🧪 Qualidade e testes

```bash
# backend
cd ponto && go test ./... -v -cover

# frontend
cd ponto-frontend && npm test
```

- **Backend**: 22 testes unitários dos services (regra de negócio), com interfaces de repositório e fakes — sem precisar de banco real para testar.
- **Frontend**: 17 testes (validação, formatação, componentes e páginas) com Vitest + Testing Library, cobertura exportada em lcov.
- Erros internos (500) sempre logados com a causa raiz antes de responder — nenhuma falha silenciosa.

## 🩹 Problemas enfrentados

- O backend derrubava com `log.Fatalf` se não encontrasse um `.env` — quebrava a inicialização dentro do Docker, onde não existe esse arquivo e as variáveis já vêm do `docker-compose`. Corrigido para apenas logar e seguir.
- Erros internos (500) eram descartados silenciosamente, sem log — qualquer falha de banco em produção seria impossível de investigar depois. Corrigido carregando a causa raiz no erro e logando antes de responder.
- Os services dependiam do struct concreto do repositório, impedindo testar regra de negócio sem subir um banco real. Refatorado para interfaces enxutas por service.
- As telas de gestão do frontend eram apenas placeholders herdados de uma semana anterior — reescritas do zero mantendo a identidade visual já estabelecida no login em vez de introduzir um estilo novo.

## 🗂️ Gestão do projeto

Board Kanban: https://github.com/users/luskation/projects/2

## 🛣️ Roadmap

- [ ] Envio real de e-mail na recuperação de senha (hoje o token volta na resposta, só para teste manual)
- [ ] Refresh token (hoje o JWT expira em 24h sem renovação)
- [ ] Rate limiting no login
- [ ] Deploy público (Vercel + Railway/Render)
- [ ] Vídeo de demonstração

## 👤 Autor

**Lucas Rodrigues** — [github.com/luskation](https://github.com/luskation)
