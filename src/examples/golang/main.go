package main

import (
	"database/sql"
	"fmt"
	"log"
	/* "os"
	"path/filepath"
	"strings"
	"time"
	"io/ioutil" */
	_ "github.com/lib/pq"  // Importação do driver PostgreSQL
)

/* func sayHelloWorld() {
	fmt.Println("Hello, World!!")
}

func listXMLFiles() {
	files, err := ioutil.ReadDir("/xml")
	if err != nil {
		fmt.Printf("Error accessing /xml: %s\n", err)
		return
	}

	fmt.Println("Listing all available XML files!")
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".xml") {
			fmt.Printf("\t> %s\n", f.Name())
		}
	}
}

func main() {
	sayHelloWorld()
	listXMLFiles()
} */


// Configuração da string de conexão
const connStr = "user=is password=is dbname=is sslmode=disable"

// Diretório de arquivos XML
const xmlDir = "/xml"

func connection() {
	// Abrir a conexão com o banco de dados
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Mensagem de log antes do Ping
	fmt.Println("Before Ping")

	// Verificar se a conexão está funcionando
	err = db.Ping()
	if err != nil {
		log.Fatalf("Ping failed: %v", err)
	} else {
		fmt.Println("Connection successful!")
	}
}

func main() {

	
	connection()
}

	/* // Loop infinito para verificar novos arquivos a cada minuto
	for {
		// Verificar novos arquivos XML
		newFiles, err := getUnprocessedFiles(db)
		if err != nil {
			log.Println("Erro ao obter novos arquivos:", err)
		}

		// Processar novos arquivos
		for _, fileName := range newFiles {
			fmt.Printf("Novo arquivo encontrado: %s\n", fileName)

			// Gerar mensagem para o serviço broker (substitua isso com sua lógica real)
			// Aqui você pode adicionar uma tarefa de importação para cada entidade, por exemplo
			generateTaskForBroker(fileName)
		}

		// Aguardar por 1 minuto antes de verificar novamente
		time.Sleep(time.Minute)
	} */


/* // Exemplo de consulta ao banco de dados para obter arquivos não processados
func getUnprocessedFiles(db *sql.DB) ([]string, error) {
	var files []string

	// Consulta para obter arquivos não processados
	rows, err := db.Query("SELECT file_name FROM public.imported_documents WHERE deleted_on IS NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Processar resultados
	for rows.Next() {
		var fileName string
		if err := rows.Scan(&fileName); err != nil {
			return nil, err
		}
		files = append(files, fileName)
	}

	return files, nil
}

// Função para gerar tarefa para o serviço broker (substitua isso com sua lógica real)
func generateTaskForBroker(fileName string) {
	// Aqui você pode implementar a lógica para gerar mensagens para o serviço broker
	// Substitua este exemplo com sua lógica real
	fmt.Printf("Gerar mensagem para o serviço broker: %s\n", fileName)
}

// Função para verificar novos arquivos no diretório XML
func checkForNewFiles() ([]string, error) {
	var newFiles []string

	// Listar todos os arquivos no diretório XML
	files, err := ioutil.ReadDir(xmlDir)
	if err != nil {
		return nil, err
	}

	// Verificar se há novos arquivos
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()

		// Verificar se o arquivo não está na lista de arquivos processados no banco de dados
		if !isProcessed(fileName) {
			newFiles = append(newFiles, filepath.Join(xmlDir, fileName))
		}
	}

	return newFiles, nil
}

// Função para verificar se um arquivo já foi processado com base no banco de dados
func isProcessed(fileName string) bool {
	// Sua lógica para verificar se o arquivo já foi processado
	// Consulte o banco de dados ou outra fonte de informação
	// Retorne verdadeiro se o arquivo já foi processado, falso caso contrário
	return false 
} */
