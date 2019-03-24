package main

/*
	Autor: Adriano Souza
	Data: 24/03/2019
	Versão: 1.0
	Descrição: Essa aplicação tem como objetivo ler um arquivo a ser sempre destinado no diretório raíz da mesma.
	Ao ler o arquivo será feito o cleanup dos campos e a validação dos campos referentes a CPF e CNPJ(nesse caso gerando arquivos de retorno).
	Após a validação de todos registros do arquivo CSV (com tabulação definida como espaço), os mesmos serão criados na tabela
	baseteste do banco de dados PROJETO no postgresql. Como se trata de um arquivo de teste, foi usada uma abordagem simples
	e direta, por isso que a tabela a ser usada sempre é dropada e recriada no início da execução da aplicação.

*/

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"errors"
	"strconv"
	"strings"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"regexp"
)

// vriáveis utilizadas no processo de validação de CPF e CNPJ, mas também utilizadas na criação do registro na tabela
var cpf string
var cnpj1 string
var cnpj2 string
var Cpfaux string
var Privateaux string
var Incompletoaux string
var DataUltimaCompraaux string
var TicketMedioaux string
var TicketUltimaCompraaux string
var LojaMaisFrequenteaux string
var LojaUltimaCompraaux string

// Variável responsável por conter o nome absoluto do arquivo a ser lido. 
// Para essa aplicação não foi necessário colocar o caminho pois o arquivo se encontrará sempre no
// diretório raíz da aplicação
var caminho = "BaseTeste.txt"

var caminhoValidacaoCPF = "ValidacaoCPF.txt"
var caminhoValidacaoCNPJ = "ValidacaoCNPJ.txt"

// Variável responsável por armazenar o comando que criará a tabela no postgresql
var schema = `DROP TABLE baseteste; 
				CREATE TABLE baseteste (
				    cpf text,
				    private text,
				    incompleto text,
				    data_ultima_compra text,
				    ticket_medio text,
				    ticket_ultima_compra text,
				    loja_mais_frequente text,
				    loja_ultima_compra text
);`

// Record mapeando os campos da tabela basedteste a ser utilizado no processo de criação de cada registro
type Record struct {
	Cpf                string `db:"cpf" json:"cpf"`
	Private            string `db:"private" json:"private"`
	Incompleto         string `db:"incompleto" json:"incompleto"`
	DataUltimaCompra   string `db:"data_ultima_compra" json:"data_ultima_compra"`
	TicketMedio        string `db:"ticket_medio" json:"ticket_medio"`
	TicketUltimaCompra string `db:"ticket_ultima_compra" json:"ticket_ultima_compra"`
	LojaMaisFrequente  string `db:"loja_mais_frequente" json:"loja_mais_frequente"`
	LojaUltimaCompra   string `db:"loja_ultima_compra" json:"loja_ultima_compra"`
}

// constante com todos os dados necessários para a conexão ao banco de dados PROJETO no postgresql
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "PROJETO"
)



func main() {

	var conteudoCPF []string
	var conteudoCNPJ []string
	var err error

	// String de conexão com banco de dados postgresql bem como a própria abertura da conexão
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		log.Fatalln(err)
	}

	// Trecho de código que lê o conteudo do arquivo(variável caminho) e retorna um slice the string com todas as linhas do arquivo
	// o Split realizado fo baseado em espaçoes, mantendo o padrão do arquivo. Nota que em nenum momento o conteúdo do arquivo foi alterado
	// mas sim somente seu nome.
	var headerFlag int
	csvFile, _ := os.Open(caminho)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = ' '
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1

	db.MustExec(schema)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		if headerFlag == 0 {
			headerFlag = 1
			reader.FieldsPerRecord = 8
		} else {
			
			// Validação dos campos de CPF e CNPJ. armazenando em arquivos (variáveis caminhoValidacaoCPF e caminhoValidacaoCNPJ).
			cpf, err =  validaCPF(row[0])
			conteudoCPF = append(conteudoCPF,cpf)
			cnpj1, err =  validaCNPJ(row[6])
			conteudoCNPJ = append(conteudoCNPJ,cnpj1)
			cnpj2, err =  validaCNPJ(row[7])
			conteudoCNPJ = append(conteudoCNPJ,cnpj2)

			//Clean up das variáveis
			Cpfaux, err = limpaString(row[0])
			Privateaux, err = limpaString(row[1])
			Incompletoaux, err = limpaString(row[2])
			DataUltimaCompraaux, err = limpaString(row[3])
			TicketMedioaux, err = limpaString(row[4])
			TicketUltimaCompraaux, err = limpaString(row[5])
			LojaMaisFrequenteaux, err = limpaString(row[6])
			LojaUltimaCompraaux, err = limpaString(row[7])
			
			var record = Record{
				Cpf:                Cpfaux,
				Private:            Privateaux,
				Incompleto:         Incompletoaux,
				DataUltimaCompra:   DataUltimaCompraaux,
				TicketMedio:        TicketMedioaux,
				TicketUltimaCompra: TicketUltimaCompraaux,
				LojaMaisFrequente:  LojaMaisFrequenteaux,
				LojaUltimaCompra:   LojaUltimaCompraaux,
			}
			tx := db.MustBegin()
			tx.NamedExec(`INSERT INTO baseteste
                    VALUES (:cpf, :private, :incompleto, :data_ultima_compra,
                      :ticket_medio, :ticket_ultima_compra, :loja_mais_frequente,
                      :loja_ultima_compra)`, &record)
			tx.Commit()
		}
	}
	// Crias os arquivos com as validações de CPF e CNPJ
	err = escreverTextoArquivo(conteudoCPF,caminhoValidacaoCPF)
	err = escreverTextoArquivo(conteudoCNPJ,caminhoValidacaoCNPJ)
}

// Função para validar o CPF referente ao campo Cpf row[0] do record
func validaCPF(cpf string) (string,error) {
	cpf = strings.Replace(cpf, ".", "", -1)
	cpf = strings.Replace(cpf, "-", "", -1)
	if len(cpf) != 11 {
		return "CPF inválido: " + cpf, errors.New("CPF inválido: " + cpf)
	}
	var eq bool
	var dig string
	for _, val := range cpf {
		if len(dig) == 0 {
			dig = string(val)
		}
		if string(val) == dig {
			eq = true
			continue
		}
		eq = false
		break
	}
	if eq {
		return "CPF inválido: " + cpf,  errors.New("CPF inválido: " + cpf)
	}

	i := 10
	sum := 0
	for index := 0; index < len(cpf)-2; index++ {
		pos, _ := strconv.Atoi(string(cpf[index]))
		sum += pos * i
		i--
	}

	prod := sum * 10
	mod := prod % 11
	if mod == 10 {
		mod = 0
	}
	digit1, _ := strconv.Atoi(string(cpf[9]))
	if mod != digit1 {
		return "CPF inválido: " + cpf, errors.New("CPF inválido: " + cpf)
	}
	i = 11
	sum = 0
	for index := 0; index < len(cpf)-1; index++ {
		pos, _ := strconv.Atoi(string(cpf[index]))
		sum += pos * i
		i--
	}
	prod = sum * 10
	mod = prod % 11
	if mod == 10 {
		mod = 0
	}
	digit2, _ := strconv.Atoi(string(cpf[10]))
	if mod != digit2 {
		return "CPF inválido: " + cpf, errors.New("CPF inválido: " + cpf)
	}

	return "CPF válido: " + cpf, nil
}

// Função para validar o CNPJ referente aos campos LojaMaisFrequente row[6] e LojaUltimaCompra row[7] do record
func validaCNPJ(cnpj string) (string, error) {
	cnpj = strings.Replace(cnpj, ".", "", -1)
	cnpj = strings.Replace(cnpj, "-", "", -1)
	cnpj = strings.Replace(cnpj, "/", "", -1)
	if len(cnpj) != 14 {
		return "CNPJ inválido: " + cnpj, errors.New("CNPJ inválido: " + cnpj)
	}

	algs := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	var algProdCpfDig1 = make([]int, 12, 12)
	for key, val := range algs {
		intParsed, _ := strconv.Atoi(string(cnpj[key]))
		sumTmp := val * intParsed
		algProdCpfDig1[key] = sumTmp
	}
	sum := 0
	for _, val := range algProdCpfDig1 {
		sum += val
	}
	digit1 := sum % 11
	if digit1 < 2 {
		digit1 = 0
	} else {
		digit1 = 11 - digit1
	}
	char12, _ := strconv.Atoi(string(cnpj[12]))
	if char12 != digit1 {
		return "CNPJ inválido: " + cnpj, errors.New("CNPJ inválido: " + cnpj)
	}
	algs = append([]int{6}, algs...)

	var algProdCpfDig2 = make([]int, 13, 13)
	for key, val := range algs {
		intParsed, _ := strconv.Atoi(string(cnpj[key]))

		sumTmp := val * intParsed
		algProdCpfDig2[key] = sumTmp
	}
	sum = 0
	for _, val := range algProdCpfDig2 {
		sum += val
	}

	digit2 := sum % 11
	if digit2 < 2 {
		digit2 = 0
	} else {
		digit2 = 11 - digit2
	}
	char13, _ := strconv.Atoi(string(cnpj[13]))
	if char13 != digit2 {
		return "CNPJ inválido: " + cnpj, errors.New("CNPJ inválido: " + cnpj)
	}

	return "CNPJ válido: " + cnpj, nil
}

// Funcao que escreve um texto no arquivo e retorna um erro caso tenha algum problema
func escreverTextoArquivo(linhas []string, caminhoDoArquivo string) error {
	// Cria o arquivo de texto
	arquivo, err := os.Create(caminhoDoArquivo)
	// Caso tenha encontrado algum erro retornar ele
	if err != nil {
		return err
	}
	// Garante que o arquivo sera fechado apos o uso
	defer arquivo.Close()

	// Cria um escritor responsavel por escrever cada linha do slice no arquivo de texto
	escritor := bufio.NewWriter(arquivo)
	for _, linha := range linhas {
		fmt.Fprintln(escritor, linha)
	}

	// Caso a funcao flush retorne um erro ele sera retornado aqui tambem
	return escritor.Flush()
}

func limpaString(texto string) (string,error) {
	
	reg, err := regexp.Compile("[^a-zA-Z0-9.,-/]+")

    if err != nil {
        return "Erro limpando string: " + texto, err
    }

    processedString := reg.ReplaceAllString(texto, "")

	return processedString, nil
}
