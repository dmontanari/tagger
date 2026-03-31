# Tagger: O Canivete Suíço de Tags do Git

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

. : Versionamento Semântico Pragmático para CI/CD baseado em Git : .

![tagger](images/tagger.png)

O tagger é uma ferramenta CLI escrita em Go desenvolvida para automatizar o Versionamento Semântico (SemVer) baseado puramente em Tags do Git.

Este projeto está em desenvolvimento (WIP). Ainda não está pronto para produção.

### Motivação

Nos ecossistemas modernos de CI/CD, o gerenciamento de versões de artefatos (imagens Docker, Binários, Releases) tornou-se desnecessariamente complexo. As soluções atuais muitas vezes dependem de estados externos (Redis, bancos de dados), arquivos de versão duplicados dentro do repositório, ou pior, extraem metadados de mensagens de commit poluídas com emojis e prefixos redundantes.

O tagger é construído sob a premissa de que o Git é a única fonte da verdade. As tags são marcadores nativos e imutáveis; usá-las para gerenciar versões é a maneira mais limpa e determinística de orquestrar uma pipeline GitOps.

### O que o Tagger Resolve

- Zero Dependências Externas: Sem necessidade de bancos de dados externos ou arquivos `.version` que geram commits circulares na sua pipeline.

- Logs Pragmáticos: Libera seu histórico de commits para o propósito original: documentar a intenção técnica e soluções, e não carregar flags de infraestrutura.

- Consistência SemVer: Garante que os incrementos de versão (Major, Minor, Patch) sigam estritamente a especificação SemVer 2.0.0.

- Performance em CI/CD: Sendo um binário estático em Go, o tagger é ideal para runners leves, não exigindo runtimes pesados (Python, Node, PHP).

### Princípios de Engenharia

- Única Fonte da Verdade: A versão atual é sempre a tag SemVer mais alta encontrada no histórico do Git.

- Imutabilidade: Uma vez enviada para o remote, uma tag é a autoridade definitiva da release.

- Zero Inchaço: Focado em fazer apenas uma coisa — gerenciar tags com o máximo de eficiência.

### Uso - Fluxo de Trabalho GitOps Sugerido

O Tagger foi desenhado para atuar como o gatilho de transição entre o seu processo de Integração Contínua (CI) e o processo de Entrega/Deploy Contínuo (CD). Ele garante que o Git seja a única fonte da verdade, independentemente de onde a sua aplicação será hospedada (Kubernetes, AWS Lambda, Azure Web Apps ou Servidores Bare Metal).

Abaixo está uma sugestão de arquitetura dividindo o ciclo de vida em três fases lógicas:

**Fase 1: Integração Contínua (App Repository)**
1. Um desenvolvedor abre um Pull Request (PR) com a nova funcionalidade.
2. O Code Review é efetuado e o PR é aprovado.
3. O merge na branch principal dispara a pipeline primária de CI.
4. A pipeline executa a bateria de testes unitários, linters e validações de segurança.
5. Com todos os testes passando, a pipeline executa o **Tagger** (`tagger inc . -m`), gerando uma nova tag anotada no repositório.

**Fase 2: Release e Empacotamento (Trigado pela nova Tag)**
1. A criação da nova tag dispara automaticamente uma segunda pipeline focada exclusivamente no processo de Release.
2. Esta pipeline lê a versão da tag recém-criada (ex: `v1.2.0`) e constrói o **artefato da aplicação** atrelado a essa versão (pode ser uma imagem Docker, um binário compilado, um arquivo `.zip` para o AWS Lambda, etc.).
3. O artefato versionado é publicado no seu repositório de destino (Docker Hub, ACR, AWS S3, GitHub Releases, etc.).

**Fase 3: GitOps e Atualização de Estado (Deploy Repository)**
1. A mesma pipeline de Release clona um repositório separado que guarda as configurações de infraestrutura e deploy.
2. Um script atualiza a versão da aplicação nos arquivos de configuração (ex: `terraform.tfvars`, `serverless.yml`, `ansible/group_vars` ou manifestos K8s) para a nova versão `v1.2.0`.
3. A pipeline faz o commit e o push dessa alteração de volta no repositório de infraestrutura.
4. O motor de Entrega Contínua que monitora este repositório (ArgoCD, Flux) detecta a mudança de estado e aplica a nova versão no ambiente de destino.

### Comandos

#### Ajuda

```bash
$ tagger help


	. : Git tag Swiss Army Knife : .

Usage:
  tagger [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  inc         Create new tag incrementing version number.
  last        Return last tag in repository path.
  list        List all tags in repository path.
  version     Show version.

Flags:
  -d, --dry-run         Dry run
  -h, --help            help for tagger
  -r, --remote string   Remote name to use (default "origin")
  -V, --verbose         Verbose mode

Use "tagger [command] --help" for more information about a command.
```

#### Listar todas as tags

```bash
$ tagger help list
list [repository path] List all tags in repository path.

Usage:
  tagger list [repository path] [flags]

Flags:
  -h, --help   help for list

Global Flags:
  -d, --dry-run         Dry run
  -r, --remote string   Remote name to use (default "origin")
  -V, --verbose         Verbose mode

$ tagger list /path/to/repo
Tags:
2026-03-28 12:26  v0.0.1

$ tagger list /path/to/remote --remote repo2
Tags:
2026-03-30 20:24  v0.0.7
2026-03-30 20:22  v0.0.6
2026-03-30 20:07  v0.0.5
2026-03-30 20:03  v0.0.4
2026-03-30 20:03  v0.0.3
2026-03-29 00:50  v0.0.2
2026-03-28 12:26  v0.0.1
```

#### Incrementar versão (Major, Minor ou Patch)

O tagger identifica a versão mais recente, aplica o incremento lógico e, opcionalmente, faz o push para o remote.

Nota: Incrementar um nível de versão mais alto zera os inferiores (ex: um incremento Major na v2.1.35 resulta na v3.0.0).

```bash
inc [repository path] [flags] Create new tag incrementing version number.
	Tags must follow the pattern vM.m.p.
	Incrementing a higher version level resets lower ones (e.g., a Major bump on v2.1.35 results in v3.0.0).

Usage:
  tagger inc [repository path] [--dry-run|-d] [flags]

Flags:
  -a, --author string    Author for commit (default "Tagger")
  -e, --email string     Email for commit (default "tagger@bot")
  -h, --help             help for inc
  -M, --major            Increment major version
      --message string   Commit message (default "Version generated by tagger")
  -m, --minor            Increment minor version
  -p, --patch            Increment patch version

Global Flags:
  -d, --dry-run         Dry run
  -r, --remote string   Remote name to use (default "origin")
  -V, --verbose         Verbose mode

$ tagger inc -M . --dry-run --verbose
v0.0.1 -> v1.0.0
$ tagger inc -m . --dry-run --verbose
v0.0.1 -> v0.1.0
$ tagger inc -p . --dry-run --verbose
v0.0.1 -> v0.0.2
$ tagger inc -p . --dry-run
v0.0.2
```

#### Obter a última tag

```bash
last [repository path] Return last tag in repository path.

Usage:
  tagger last [repository path] [flags]

Flags:
  -f, --full   Full output - date and tag
  -h, --help   help for last

Global Flags:
  -d, --dry-run         Dry run
  -r, --remote string   Remote name to use (default "origin")
  -V, --verbose         Verbose mode

$ tagger last /path/to/repo
v0.0.7

$ tagger last /path/to/repo -f
2026-03-30 20:24  v0.0.7
```

### Compilando a partir do código-fonte

Construindo a partir da fonte:

    Clone este repositório: git clone https://github.com/dmontanari/tagger.git

    Compile o binário: make build

O binário estará no diretório raiz.

### Instalação

Em breve

### Licença

Distribuído sob a licença MIT. Veja o arquivo LICENSE para mais informações.

© 2026 Daniel Montanari. Todos os direitos reservados.
