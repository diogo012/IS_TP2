package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	_ "github.com/lib/pq" // Importação do driver PostgreSQL
	"github.com/streadway/amqp"
	//"github.com/streadway/amqp"
)

// Configuração da string de conexão
const connStr = "user=is password=is dbname=is sslmode=disable host=db-xml port=5432"

// Diretório de arquivos XML
const xmlDir = "/xml"

func connection() {
	// Abrir a conexão com o banco de dados
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Verificar se a conexão está funcionando
	err = db.Ping()
	if err != nil {
		log.Fatalf("Ping failed: %v", err)
	} else {
		fmt.Println("Connection successful!")
	}

	// Loop infinito para verificar novas entidades a cada minuto
	for {
		// Verificar novas entidades XML
		newEntities, err := getUnprocessedEntities(db)
		if err != nil {
			log.Println("Erro ao obter novas entidades:", err)
		}

		// Processar novas entidades
		for _, EntitiesName := range newEntities {
			fmt.Printf("Nova entidade encontrada: %s\n", EntitiesName)

			// Verificar se a entidade já foi processado
			if isProcessed(EntitiesName) {
				fmt.Printf("Entidade já foi encontrado anteriormente: %s\n", EntitiesName)
			} else {
				// Gerar mensagem para o serviço broker (substitua isso com sua lógica real)
				// Aqui você pode adicionar uma tarefa de importação para cada entidade, por exemplo
				generateTaskForBroker(EntitiesName)
			}
		}

		// Aguardar por 1 minuto antes de verificar novamente
		time.Sleep(10 * time.Minute)
	}
}

// Exemplo de consulta ao banco de dados para obter entidades não processados
func getUnprocessedEntities(db *sql.DB) ([]string, error) {
	var files []string

	// Consulta para obter entidades não processados
	rows, err := db.Query("SELECT unnest(xpath('/Jobs/JobPortals/JobPortal', xml)) AS job_attributes FROM public.imported_documents;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Processar resultados
	for rows.Next() {
		var EntitiesName string
		if err := rows.Scan(&EntitiesName); err != nil {
			return nil, err
		}
		files = append(files, EntitiesName)
	}

	return files, nil
}

// Função para gerar tarefa para o serviço broker (substitua isso com sua lógica real)
func generateTaskForBroker(EntitiesName string) {
	// Conectar ao servidor RabbitMQ
	conn, err := amqp.Dial("amqp://is:is@rabbitmq:5672/is")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Criar um canal
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declarar uma fila
	q, err := ch.QueueDeclare(
		"queue_migrator", // Nome da fila
		false,            // Durable
		false,            // Delete when unused
		false,            // Exclusive
		false,            // No-wait
		nil,              // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Publicar uma mensagem na fila
	err = ch.Publish(
		"",     // Exchange
		q.Name, // Key
		false,  // Mandatory
		false,  // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(EntitiesName),
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}

	fmt.Printf("Mensagem enviada para o serviço broker \n")
}

// Função para verificar novos arquivos no diretório XML
func checkFornewEntities() ([]string, error) {
	var newEntities []string

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
		EntitiesName := file.Name()

		// Verificar se o arquivo não está na lista de arquivos processados no banco de dados
		if !isProcessed(EntitiesName) {
			newEntities = append(newEntities, filepath.Join(xmlDir, EntitiesName))
		}
	}

	return newEntities, nil
}

// Função para verificar se um arquivo já foi processado com base no banco de dados
func isProcessed(EntitiesName string) bool {
	// Sua lógica para verificar se o arquivo já foi processado
	// Consulte o banco de dados ou outra fonte de informação
	// Retorne verdadeiro se o arquivo já foi processado, falso caso contrário
	return false
}

func main() {
	connection()
}
