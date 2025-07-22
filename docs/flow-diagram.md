# Diagramas do Fluxo do Sistema

## 🔄 Fluxo Principal de Processamento de Pagamentos

```mermaid
sequenceDiagram
    participant Cliente
    participant Nginx
    participant Go1 as Aplicação Go (go1)
    participant Go2 as Aplicação Go (go2)
    participant NATS
    participant Subscriber
    participant HealthCheck
    participant Redis
    participant Default as Processador Default
    participant Fallback as Processador Fallback
    participant PostgreSQL

    Cliente->>Nginx: POST /payments
    Nginx->>Go1: Distribui carga (round-robin)
    Note over Go1: Valida requisição
    Go1->>NATS: Publica mensagem (pub.payments)
    Go1->>Cliente: 202 Accepted

    NATS->>Subscriber: Consome mensagem
    Subscriber->>HealthCheck: Verifica processador ativo
    HealthCheck->>Redis: Consulta cache (8s TTL)
    
    alt Cache válido
        Redis->>HealthCheck: Retorna processador ativo
    else Cache expirado
        HealthCheck->>Default: GET /payments/service-health
        alt Default saudável
            Default->>HealthCheck: 200 OK
            HealthCheck->>Redis: Cache processador default
        else Default falhou
            HealthCheck->>Fallback: GET /payments/service-health
            Fallback->>HealthCheck: 200 OK
            HealthCheck->>Redis: Cache processador fallback
        end
    end

    HealthCheck->>Subscriber: Retorna processador ativo
    Subscriber->>Default: POST /payments (processamento)
    Default->>Subscriber: Resposta do processamento
    Subscriber->>PostgreSQL: Insere resultado
    PostgreSQL->>Subscriber: Confirma inserção
```

## 🏥 Estratégia de Fallback

```mermaid
flowchart TD
    A[Health Check Iniciado] --> B{Default Processador<br/>Saudável?}
    B -->|Sim| C[Usar Default]
    B -->|Não| D{Fallback Processador<br/>Saudável?}
    D -->|Sim| E[Usar Fallback]
    D -->|Não| F[Manter Último Conhecido]
    
    C --> G[Atualizar Cache Redis]
    E --> G
    F --> G
    
    G --> H[Log da Decisão]
    H --> I["Aguardar Próximo Check<br/>(8 segundos)"]
    I --> A
    
    style C fill:#90EE90
    style E fill:#FFB6C1
    style F fill:#FFD700
```

## 🏗️ Arquitetura de Componentes

```mermaid
graph TB
    subgraph "Cliente"
        Client[Cliente HTTP]
    end
    
    subgraph "Load Balancer"
        Nginx[Nginx<br/>Porta 9999]
    end
    
    subgraph "Aplicação"
        Go1[Go Instance 1<br/>Porta 8080]
        Go2[Go Instance 2<br/>Porta 8080]
    end
    
    subgraph "Mensageria"
        NATS[NATS Server<br/>Comunicação Assíncrona]
    end
    
    subgraph "Processamento"
        Subscriber[Subscriber<br/>Consome mensagens]
        HealthCheck[Health Check<br/>Monitora processadores]
    end
    
    subgraph "Cache"
        Redis[Redis<br/>Cache de processadores]
    end
    
    subgraph "Processadores Externos"
        Default[Processador Default<br/>URL Principal]
        Fallback[Processador Fallback<br/>URL Alternativa]
    end
    
    subgraph "Persistência"
        PostgreSQL[(PostgreSQL<br/>Dados de pagamentos)]
    end
    
    Client --> Nginx
    Nginx --> Go1
    Nginx --> Go2
    Go1 --> NATS
    Go2 --> NATS
    NATS --> Subscriber
    Subscriber --> HealthCheck
    HealthCheck --> Redis
    HealthCheck --> Default
    HealthCheck --> Fallback
    Subscriber --> Default
    Subscriber --> Fallback
    Subscriber --> PostgreSQL
```

## 📊 Fluxo de Dados Detalhado

```mermaid
graph LR
    subgraph "1. Recebimento"
        A[POST /payments] --> B[Validação]
        B --> C[Publicação NATS]
        C --> D[202 Accepted]
    end
    
    subgraph "2. Processamento Assíncrono"
        E[Consumo NATS] --> F[Health Check]
        F --> G{Processador Ativo?}
        G -->|Default| H[Envia para Default]
        G -->|Fallback| I[Envia para Fallback]
        H --> J[Persiste no PostgreSQL]
        I --> J
    end
    
    subgraph "3. Monitoramento"
        K[Health Check Loop<br/>8 segundos] --> L{Default OK?}
        L -->|Sim| M[Cache Default]
        L -->|Não| N{Fallback OK?}
        N -->|Sim| O[Cache Fallback]
        N -->|Não| P[Manter Último]
    end
    
    A --> E
    K --> F
```

## 🔍 Estados do Sistema

```mermaid
stateDiagram-v2
    [*] --> Inicializando
    Inicializando --> VerificandoDefault
    VerificandoDefault --> DefaultAtivo: Default OK
    VerificandoDefault --> VerificandoFallback: Default Falhou
    VerificandoFallback --> FallbackAtivo: Fallback OK
    VerificandoFallback --> DefaultAtivo: Ambos Falharam
    
    DefaultAtivo --> VerificandoDefault: Health Check (8s)
    FallbackAtivo --> VerificandoDefault: Health Check (8s)
    
    DefaultAtivo --> VerificandoFallback: Default Falhou
    FallbackAtivo --> DefaultAtivo: Default Recuperado
    
    note right of DefaultAtivo
        Processador Default ativo
        Cache Redis atualizado
        Logs de sucesso
    end note
    
    note right of FallbackAtivo
        Processador Fallback ativo
        Cache Redis atualizado
        Logs de warning
    end note
```

## 📈 Métricas e Monitoramento

```mermaid
graph TD
    subgraph "Métricas Coletadas"
        A[Tempo de Resposta<br/>Processadores]
        B[Taxa de Sucesso/Erro]
        C[Latência de Processamento]
        D[Status Health Check]
        E[Cache Hit/Miss]
    end
    
    subgraph "Logs Gerados"
        F[Mudança de Processador]
        G[Health Check Failures]
        H[Erros de Conexão]
        I[Performance Metrics]
    end
    
    subgraph "Endpoints de Monitoramento"
        J[GET /processors/status]
        K[GET /payments-summary]
        L[Logs Docker]
    end
    
    A --> J
    B --> J
    C --> J
    D --> J
    E --> J
    
    F --> L
    G --> L
    H --> L
    I --> L
```

---

## 📝 Como Usar os Diagramas

Estes diagramas podem ser visualizados em:

1. **GitHub**: Os diagramas Mermaid são renderizados automaticamente
2. **Mermaid Live Editor**: https://mermaid.live/
3. **VS Code**: Com extensão Mermaid Preview
4. **Documentação**: Qualquer ferramenta que suporte Mermaid

## 🔧 Personalização

Para modificar os diagramas:

1. Edite o arquivo `docs/flow-diagram.md`
2. Use a sintaxe Mermaid para criar novos diagramas
3. Teste no Mermaid Live Editor antes de commitar
4. Os diagramas serão atualizados automaticamente no GitHub 