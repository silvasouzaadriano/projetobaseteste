# Projeto baseteste
Projeto BaseTeste desenvolvido e testado em ambiente Windows 7. Essa aplicação tem como objetivo ler um arquivo a ser sempre destinado no diretório raíz da mesma. Ao ler o arquivo será feito o cleanup dos campos e a validação dos campos referentes a CPF e CNPJ(nesse caso gerando arquivos de retorno).	Após a validação de todos registros do arquivo CSV (com tabulação definida como espaço), os mesmos serão criados na tabela baseteste do banco de dados PROJETO no postgresql. Como se trata de um arquivo de teste, foi usada uma abordagem simples
e direta, por isso que a tabela a ser usada sempre é dropada e recriada no início da execução da aplicação.

# Passos para utilização

 1- Instalar go1.12.1.windows-amd64 ou a versão mais atual
 
 2- Instalar postgresql-11.2-1-windows-x64 ou versão mais atual
 
 3- Através de linha de comando innstalar o pacotes GO: "github.com/jmoiron/sqlx" (arquivo sqlx-master) e _ "github.com/lib/pq" (arquivo pq-1.0.0). Exemplo: go get github.com/jmoiron/sqlx e go get github.com/lib/pq
 
 4- Criar banco de dados postgresql denominado: PROJETO 
  
 5- Manter no diretório raíz da aplicação os arquivos BaseTeste.txt (dados) e importBaseTeste.go
 
 6- Utilizando o prompt de comando do Go, executar a linha de comando: go run importBaseTeste.go
 
 7- No diretório raíz serão criados dois arquivos referentes a validação: ValidacaoCPF.txt e ValidacaoCNPJ.txt
 
 8- No banco de dados PROJETO, tabela baseteste, serão criados os registros já sanitizado
