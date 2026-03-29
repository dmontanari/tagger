# Tagger: The Git Tag Swiss Army Knife

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)


. : Pragmatic Semantic Versioning for Git-based CI/CD : .

tagger é uma ferramenta de linha de comando escrita em Go para automação de versionamento semântico (SemVer) baseada puramente em Git Tags.

No ecossistema de CI/CD moderno, o gerenciamento de versões de artefatos (Imagens Docker, Binários, Releases) tornou-se desnecessariamente complexo. Muitas soluções atuais dependem de estados externos (Redis, bancos de dados), arquivos de versão duplicados no repositório ou, pior, a extração de metadados de mensagens de commit poluídas com emojis e prefixos redundantes.

O tagger nasce da premissa de que o Git é a única fonte da verdade. Tags são marcadores imutáveis e nativos; utilizá-las para gerenciar a versão é a forma mais limpa e determinística de orquestrar uma pipeline de GitOps.

### O que o Tagger resolve

* Eliminação de Dependências: Não requer bancos de dados externos ou arquivos .version que geram commits circulares na pipeline.

* Logs Pragmáticos: Libera o histórico de commits para sua função original: documentar intenções e soluções técnicas, não para carregar flags de versionamento.

* Consistência SemVer: Garante que o incremento de versões (Major, Minor, Patch) siga rigorosamente o padrão SemVer 2.0.0.

* Performance em CI/CD: Como um binário estático em Go, o tagger é ideal para execução em runners leves, sem necessidade de runtimes (Python, Node, PHP) ou instalação de pacotes pesados.

Princípios de Engenharia

* Single Source of Truth: A versão atual é sempre a maior tag SemVer encontrada no histórico do Git.

* Imutabilidade: Uma vez criada e enviada ao remoto, a tag é a autoridade máxima do release.

* Zero Bloat: Ferramenta focada em fazer uma única coisa: gerenciar tags com a máxima eficiência.

### Como funciona

#### Listar versões atuais

```bash
tagger list [repository_path]
```

#### Incrementar Versão (Major, Minor ou Patch)

O tagger identifica a última versão, aplica o incremento lógico e, opcionalmente, realiza o push para o remoto.

```bash
tagger inc [repository_path] --patch  # Incrementa v1.0.0 para v1.0.1
tagger inc [repository_path] --minor  # Incrementa v1.0.1 para v1.1.0
tagger inc [repository_path] --major  # Incrementa v1.1.0 para v2.0.0
```

Nota: Incrementar um nível de versão zera os demais. Por exemplo, incrementar o major da versão v2.1.35 vai gerar a tag da versão v3.0.0.

## Compilação

Building from source

1. Clone o repositório: git clone https://github.com/dmontanari/tagger.git

2. Compile: make build

O binário estará no diretório src/tagger.


## Instalação

Em breve via Homebrew e binários pré-compilados.

## Licença

Distribuído sob a licença MIT. Veja LICENSE para mais informações.

© 2026 Daniel Montanari. Todos os direitos reservados.


