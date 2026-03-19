# Roteiro do Projeto: Rar-Cracker

Este documento descreve a especificação do projeto **rar-cracker**, 
uma ferramenta de linha de comando para recuperação de senhas de arquivos RAR5.  


## 1. Visão Geral do Projeto

O objetivo é criar uma ferramenta que, dado um arquivo `.rar`, 
uma parte conhecida da senha (máscara) e um conjunto de caracteres, 
tente descobrir a senha completa.

### Recursos principais
- **CPU-Bound:** todo o processamento ocorre em memória, sem processos externos.  
- **Máscara Dinâmica:** combinações de prefixo e sufixo em torno da parte conhecida.  
- **Concorrência:** uso de goroutines para explorar todos os núcleos do processador.  
- **Resiliência:** possibilidade de pausar e retomar o progresso.  
- **Simplicidade:** dependência mínima, apenas biblioteca padrão do Go.  


## 2. Arquitetura e Componentes

- **cmd/** → ponto de entrada (`main.go`), que apenas orquestra a execução.  
- **internal/** → implementações concretas (parsing, workers, geração de senhas, persistência e UI).  
- **domain/** → contratos, modelos e entidades que representam as regras de negócio centrais.  


## 3. Etapas de Desenvolvimento

### Etapa 1 — Configuração  
Parsing dos argumentos de linha de comando.  
- Parâmetros obrigatórios: `-file`,  `-charset`.  
- Parâmetros opcionais: `-knownPart`, `-maxLength`, `-workers`, `-stateFile`.  

### Etapa 2 — Análise do Arquivo RAR  
Leitura do cabeçalho do RAR5 e extração das informações criptográficas necessárias para validação das senhas.  

### Etapa 3 — Worker Criptográfico  
Executa o teste de cada senha candidata.  
- Deriva chave com PBKDF2.  
- Descriptografa bloco de verificação.  
- Compara resultado com hash esperado.  

### Etapa 4 — Gerador de Senhas  
Produz senhas válidas a partir da parte conhecida e do charset.  
- Envia para os workers por canais.  
- Suporta retomar a execução usando estado salvo.  

### Etapa 5 — Orquestração  
Responsabilidade do `main.go`.  
- Carrega configuração e estado.  
- Inicializa workers, gerador e UI.  
- Gerencia ciclo de vida e encerramento limpo das goroutines.  


## 4. Funcionalidades Avançadas

### Persistência de Estado  
- Progresso salvo em JSON.  
- Permite retomar a execução do ponto exato de parada.  

### Interface de Terminal  
- Barra de progresso, taxa de senhas testadas, tempo decorrido e estimativa de conclusão (ETA).  
- Atualizações em tempo real, sem poluir o terminal.  

### Eficiência de Execução  
- As informações do RAR são lidas apenas uma vez na inicialização.  
- Compartilhadas em memória com todos os workers para evitar reprocessamento.  


## 5. Estrutura de Arquivos

* rar-cracker/
* ├── cmd/
* │   └── rar-cracker/
* │   │   └── main.go
* │
* ├── internal/ ................................................ # Implementações específicas (infra)
* │   ├── config/ ............................................... # Parsing de flags e configuração
* │   │   ├── config.go
* │   │   └── config_validations.go
* │   │
* │   ├── generator/ ....................................... # Geração de senhas
* │   │   └── generator.go
* │   │
* │   ├── rar/ ..................................................... # Leitura do cabeçalho do arquivo RAR
* │   │   └── parser.go
* │   │
* │   ├── state/ .................................................. # Persistência de estado (pausar/retomar)
* │   │   └── state.go
* │   │
* │   ├── ui/ ........................................................ # UI do terminal (barra de progresso)
* │   │   └── ui.go
* │   │
* │   └── worker/ ............................................. # Worker para teste de senhas (CPU-bound)
* │       └── worker.go
* │
* ├── domain/ ................................................ # Regras de negócio (entidades e contratos)
* │   ├── rar_metadata.go ............................. # Estruturas de dados e entidades
* │   ├── worker_service.go .......................... # Contrato do serviço de workers
* │   ├── generator_service.go .................... # Contrato do gerador de senhas
* │   └── state_model.go
* │
* ├── go.mod
* ├── go.sum
* └── README.md