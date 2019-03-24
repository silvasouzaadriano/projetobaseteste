# projetobaseteste
Projeto BaseTeste desenvolvido e testado em ambiente Windows 7

# Passos para utilização

1- Instalar go1.12.1.windows-amd64 ou a versão mais atual
2- Instalar postgresql-11.2-1-windows-x64 ou versão mais atual
3- Através de linha de comando innstalar o pacotes GO: "github.com/jmoiron/sqlx" (arquivo sqlx-master) e _ "github.com/lib/pq" (arquivo pq-1.0.0). Exemplo: go get github.com/jmoiron/sqlx e go get github.com/lib/pq
3- Criar banco de dados (postgresql denominado: PROJETO
4- Manter no diretório raíz da aplicação os arquivos BaseTeste.txt (dados) e importBaseTeste.go
5- Utilizando o prompt de comando do Go, executar a linha de comando: go run importBaseTeste.go
6- No diretório raíz serão criados dois arquivos referentes a validação: ValidacaoCPF.txt e ValidacaoCNPJ.txt
7- No banco de dados PROJETO, tabela baseteste, serão criados os registros já sanitizados
