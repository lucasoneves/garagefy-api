# Relatório de Itens Pendentes — MVP Garagefy API

**Data:** 22/06/2026 (atualizado)  
**Projeto:** Garagefy API — Backend em Go (Gin + GORM + PostgreSQL)  
**Commits:** 9 (29/05 a 17/06/2026)

---

## 1. Resumo do Projeto

API de gerenciamento de manutenção veicular. Permite cadastro de veículos, registro de abastecimentos, histórico de serviços e logbook (observações/lembretes/tarefas). Possui autenticação JWT e recuperação de senha via e-mail (Mailtrap).

**23 endpoints implementados** (5 públicos + 18 protegidos).

---

## 2. O que já está implementado e funcional

### Autenticação
- [x] Registro de usuário (`POST /api/auth/register`)
- [x] Login com JWT (`POST /api/auth/login`)
- [x] Recuperação de senha via e-mail (`POST /api/auth/forgot-password`)
- [x] Redefinição de senha com token (`POST /api/auth/reset-password`)
- [x] Formulário HTML de redefinição (`GET /api/auth/reset-password`)
- [x] Middleware JWT com validação Bearer token

### Veículos (CRUD completo)
- [x] Criar veículo (`POST /api/vehicles`)
- [x] Listar veículos do usuário (`GET /api/vehicles`)
- [x] Obter veículo por ID (`GET /api/vehicles/:id`)
- [x] Atualizar veículo (`PUT /api/vehicles/:id`)
- [x] Remover veículo (soft delete) (`DELETE /api/vehicles/:id`)

### Logbook (CRUD completo)
- [x] Criar entrada (`POST /api/vehicles/:id/logbook`)
- [x] Listar entradas (`GET /api/vehicles/:id/logbook`)
- [x] Obter entrada por ID (`GET /api/vehicles/:id/logbook/:logbookId`)
- [x] Atualizar entrada (`PUT /api/vehicles/:id/logbook/:logbookId`)
- [x] Remover entrada (`DELETE /api/vehicles/:id/logbook/:logbookId`)

### Serviços (CRUD completo)
- [x] Criar serviço (`POST /api/services`)
- [x] Listar serviços por veículo (`GET /api/services?vehicle_id=UUID`)
- [x] Obter serviço por ID (`GET /api/services/:id`)
- [x] Atualizar serviço (`PUT /api/services/:id`)
- [x] Remover serviço (`DELETE /api/services/:id`)

### Abastecimentos (CRUD completo)
- [x] Criar abastecimento com cálculo automático de km/l (`POST /api/fuels`)
- [x] Listar abastecimentos por veículo (`GET /api/fuels?vehicle_id=UUID`)
- [x] Obter abastecimento por ID (`GET /api/fuels/:id`)
- [x] Atualizar abastecimento com recálculo de consumo (`PUT /api/fuels/:id`)
- [x] Remover abastecimento (`DELETE /api/fuels/:id`)
- [x] Sincronização automática do odômetro do veículo

### Infraestrutura
- [x] Conexão PostgreSQL via GORM
- [x] Migrations automáticas (AutoMigrate)
- [x] Docker + docker-compose (PostgreSQL + API)
- [x] CORS configurado
- [x] Envio de e-mail via SMTP (Mailtrap)
- [x] Carregamento de variáveis de ambiente (.env)
- [x] Upload de arquivos (`POST /api/upload`)
- [x] Servir arquivos estáticos (`GET /uploads/*`)

---

## 3. Itens já corrigidos

| # | Item | Correção |
|---|------|----------|
| P1 | **JWT secret hardcoded** | `services/auth_service.go` — substituído por `getJWTSecret()` que lê da env var `JWT_SECRET` com fallback dev. |
| P2 | **Validação de entrada rudimentar** | Criado `utils/validation.go` com `FormatValidationError()` — mensagens em português (ex: `"O campo Brand é obrigatório"`). Aplicado em todos os 5 controllers. |
| P3 | **Service.VehicleID tipo `string`** | `models/services.go` — alterado para `uuid.UUID`, consistente com os demais modelos. `controllers/service.go` — `vehicle.ID.String()` simplificado para `vehicle.ID`. |
| P4 | **Service sem `UpdatedAt`/`DeletedAt`** | `models/services.go` — adicionados os campos auditivos. |
| P5 | **Service sem `BeforeCreate` hook** | `models/services.go` — adicionado hook gerando `uuid.New()`. |
| P6 | **User com UUID redundante** | `models/user.go` — removido `default:gen_random_uuid()` da tag (hook já gerencia). |
| P7 | **Sistema de upload removido** | `main.go` — re-adicionado `r.Static("/uploads", "./uploads")`. Criado `controllers/upload.go` com `POST /api/upload`. |

---

## 4. Itens Pendentes para o MVP

### 4.1. Funcionalidades Faltantes — Médio

| # | Item | Local | Detalhes |
|---|------|-------|----------|
| P8 | **Sem endpoints de dashboard/relatórios** | — | Não há endpoints agregados para exibir métricas e resumos. |
| P9 | **Sem notificações/lembretes automáticos** | — | O logbook tem categoria "Reminder" e "To-do", mas não há disparo de notificações quando um lembrete está próximo. |
| P10 | **Sem paginação** | Todos os endpoints de listagem | `GetVehicles`, `GetLogbookEntries`, `GetServicesByVehicle` e `GetFuelLogsByVehicle` retornam todos os registros sem suporte a `page`/`limit`. |

---

#### P8 — Endpoints de Dashboard e Relatórios

**Objetivo:** Fornecer endpoints agregados para alimentar uma tela de dashboard/resumo no frontend.

**Endpoints propostos:**

```
GET  /api/dashboard/summary              # Resumo geral do usuário
GET  /api/dashboard/vehicle/:id/summary  # Resumo de um veículo específico
GET  /api/dashboard/vehicle/:id/costs    # Gastos por período (query: ?start=&end=)
GET  /api/dashboard/vehicle/:id/consumption  # Histórico de consumo (km/l ao longo do tempo)
```

**Exemplo de resposta de `/api/dashboard/summary`:**

```json
{
  "total_vehicles": 2,
  "total_services": 12,
  "total_fuel_logs": 45,
  "total_spent": 8500.00,
  "total_km": 15200,
  "avg_consumption": 12.5,
  "upcoming_services": [
    {
      "vehicle_id": "uuid",
      "vehicle_name": "Fiat Uno 2014",
      "last_service_odo": 85000,
      "next_service_odo": 86000,
      "current_odo": 85750
    }
  ],
  "recent_activity": [
    {
      "type": "service",
      "description": "Troca de óleo",
      "vehicle": "Fiat Uno",
      "date": "2026-06-20T10:00:00Z"
    }
  ]
}
```

**Exemplo de resposta de `/api/dashboard/vehicle/:id/summary`:**

```json
{
  "vehicle": {
    "id": "uuid",
    "brand": "Fiat",
    "model": "Uno",
    "year": 2014,
    "plate": "ABC-1234",
    "current_odo": 85750
  },
  "total_services": 8,
  "total_fuel_logs": 30,
  "total_spent": 5200.00,
  "total_km": 12000,
  "avg_consumption": 13.2,
  "last_service": {
    "title": "Troca de óleo",
    "date": "2026-06-01",
    "cost": 180.00,
    "odo": 85600
  },
  "last_fuel": {
    "date": "2026-06-20",
    "km_liter": 13.5,
    "odometer": 85750
  },
  "costs_by_year": [
    {"year": 2026, "total": 3200.00}
  ],
  "next_service_due_km": 86000,
  "km_until_next_service": 250
}
```

**Arquivos a criar/modificar:**
- `controllers/dashboard.go` — novo controller com os handlers
- `main.go` — registrar as novas rotas no grupo protegido

---

#### P9 — Notificações e Lembretes Automáticos

**Objetivo:** Disparar notificações (e-mail) quando lembretes do logbook estiverem próximos ou vencidos, e alertar sobre serviços programados.

**Funcionalidades:**

| Funcionalidade | Descrição |
|----------------|-----------|
| **Lembretes por data** | Entradas do logbook com categoria "Reminder" que possuem data associada (campo a adicionar no model) disparam e-mail no dia do vencimento. |
| **To-dos pendentes** | Notificação semanal com resumo de to-dos não concluídos. |
| **Alerta de serviço** | Quando o odômetro atual do veículo se aproxima do odômetro do último serviço + intervalo estimado, notificar o usuário. |
| **Checklist de abastecimento** | Se o usuário não registra abastecimento há mais de X dias, enviar lembrete. |

**Modelo de dados — campo novo no LogbookEntry:**

```go
ReminderDate *time.Time `json:"reminder_date,omitempty"`  // data agendada para o lembrete
Done         bool       `gorm:"default:false" json:"done"` // para to-dos
```

**Arquivos a criar/modificar:**
- `models/logbook.go` — adicionar campos `ReminderDate` e `Done`
- `services/notification_service.go` — novo serviço com lógica de verificação e disparo
- `services/email_service.go` — estender para suportar templates de notificação
- `main.go` — agendar task periódica (ex: goroutine com `time.Ticker` rodando a cada 1h)

**Fluxo da task agendada:**

```
a cada 1h:
  1. Buscar lembretes com ReminderDate <= now() e Done = false
  2. Para cada lembrete, disparar e-mail para o dono do veículo
  3. Marcar Done = true (ou criar campo NotifiedAt)
```

---

#### P10 — Paginação nos Endpoints de Listagem

**Objetivo:** Adicionar suporte a paginação em todos os endpoints que retornam listas, evitando sobrecarga com muitos registros.

**Endpoints afetados:**

| Endpoint | Parâmetros |
|----------|------------|
| `GET /api/vehicles` | `?page=1&limit=10` |
| `GET /api/vehicles/:id/logbook` | `?page=1&limit=20` |
| `GET /api/services?vehicle_id=UUID` | `?page=1&limit=20` |
| `GET /api/fuels?vehicle_id=UUID` | `?page=1&limit=20` |

**Formato de resposta padronizado:**

```json
{
  "data": [ ... ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 85,
    "total_pages": 5,
    "has_next": true,
    "has_prev": false
  }
}
```

**Implementação:**

Criar helper reutilizável em `utils/pagination.go`:

```go
type PaginationParams struct {
    Page  int `form:"page"`
    Limit int `form:"limit"`
}

type PaginationMeta struct {
    Page       int   `json:"page"`
    Limit      int   `json:"limit"`
    Total      int64 `json:"total"`
    TotalPages int   `json:"total_pages"`
    HasNext    bool  `json:"has_next"`
    HasPrev    bool  `json:"has_prev"`
}

type PaginatedResponse struct {
    Data       interface{}    `json:"data"`
    Pagination PaginationMeta `json:"pagination"`
}

func ParsePagination(c *gin.Context) PaginationParams {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
    if page < 1 { page = 1 }
    if limit < 1 { limit = 20 }
    if limit > 100 { limit = 100 }
    return PaginationParams{Page: page, Limit: limit}
}

func Paginate(db *gorm.DB, page, limit int, data interface{}) (*PaginationMeta, error) {
    var total int64
    db.Count(&total)

    offset := (page - 1) * limit
    if err := db.Offset(offset).Limit(limit).Find(data).Error; err != nil {
        return nil, err
    }

    totalPages := int(math.Ceil(float64(total) / float64(limit)))
    if totalPages < 1 { totalPages = 1 }

    return &PaginationMeta{
        Page:       page,
        Limit:      limit,
        Total:      total,
        TotalPages: totalPages,
        HasNext:    page < totalPages,
        HasPrev:    page > 1,
    }, nil
}
```

**Arquivos a criar/modificar:**
- `utils/pagination.go` — novo helper
- `controllers/vehicle.go` — `GetVehicles` usar paginação
- `controllers/logbook.go` — `GetLogbookEntries` usar paginação
- `controllers/service.go` — `GetServicesByVehicle` usar paginação
- `controllers/fuel.go` — `GetFuelLogsByVehicle` usar paginação

### 4.2. Qualidade e Testes — Médio

| # | Item | Local | Detalhes |
|---|------|-------|----------|
| P11 | **Zero testes automatizados** | Projeto inteiro | Nenhum arquivo `*_test.go`. Sem testes unitários, de integração ou de API. |
| P12 | **Sem validação de categoria do logbook** | `controllers/logbook.go:43` | A string da categoria é convertida diretamente para `models.LogbookCategory(input.Category)` sem validar se é um dos valores permitidos (`Observation`, `Reminder`, `To-do`). Um valor inválido causa erro 500. |
| P13 | **Erro no UpdateService: sem tratamento de erro no `Save`** | `controllers/service.go` | `config.DB.Save(&service)` ignora o retorno de erro. Se falhar, retorna 200 com dados desatualizados. |
| P14 | **Erro no DeleteService: sem tratamento de erro no `Delete`** | `controllers/service.go` | Mesmo problema do P13. |
| P15 | **Erro no UpdateLogbookEntry: sem tratamento de erro no `Save`** | `controllers/logbook.go` | Mesmo problema do P13. |

### 4.3. Documentação — Baixo

| # | Item | Local | Detalhes |
|---|------|-------|----------|
| P16 | **README.md ausente** | Raiz do projeto | Não há README com instruções de setup, configuração, uso da API. |
| P17 | **Sem documentação OpenAPI/Swagger** | Projeto inteiro | Não há `swagger.json` ou endpoints de documentação interativa. |
| P18 | **`.env.example` incompleto** | `.env.example` | Faltam as variáveis `GIN_MODE` e `ALLOWED_ORIGINS` (para CORS restrito em produção). `JWT_SECRET` já foi adicionado. |
| P19 | **CORS com wildcard `*`** | `main.go:35` | `Access-Control-Allow-Origin: *` é aceitável para desenvolvimento, mas inseguro para produção. Deve ser configurável via variável de ambiente. |

### 4.4. DevOps — Baixo

| # | Item | Local | Detalhes |
|---|------|-------|----------|
| P20 | **Sem CI/CD** | — | Nenhum workflow do GitHub Actions para testes, lint, build ou deploy. |
| P21 | **Sem configuração de ambiente produção** | — | Docker Compose apenas para dev. Sem Dockerfile otimizado para produção (multistage existe mas não há separação dev/prod). |

---

## 5. Resumo por Prioridade

| Prioridade | Itens | Esforço estimado |
|------------|-------|------------------|
| **Alto** | P11, P12 | 3-6 horas |
| **Médio** | P8, P9, P10, P13, P14, P15, P19 | 8-16 horas |
| **Baixo** | P16, P17, P18, P20, P21 | 4-8 horas |

**Total restante:** ~15-30 horas.

---

## 6. Recomendações Próximos Passos

1. **Curto prazo:** Validar categoria do logbook (P12) e tratar erros ignorados nos controllers (P13, P14, P15).
2. **Médio prazo:** Adicionar paginação (P10) e testes (P11). Implementar dashboard com endpoints agregados (P8).
3. **Documentação:** Criar README (P16) e completar `.env.example` (P18). Tornar CORS configurável por env var (P19).
4. **DevOps:** Configurar CI/CD (P20) e separação dev/prod no Docker (P21).
