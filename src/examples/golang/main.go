package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // Importação do driver PostgreSQL
	"github.com/streadway/amqp"
	//"github.com/streadway/amqp"
)

// Configuração da string de conexão
const connStr = "user=is password=is dbname=is sslmode=disable host=db-xml port=5432"

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
			if isProcessed(db, EntitiesName) {
				fmt.Printf("Entidade já foi encontrado anteriormente: %s\n", EntitiesName)
			} else {
				// Gerar mensagem para o serviço broker (substitua isso com sua lógica real)
				// Aqui você pode adicionar uma tarefa de importação para cada entidade, por exemplo
				generateTaskForBroker(EntitiesName)
			}
		}

		// Aguardar por 1 minuto antes de verificar novamente
		time.Sleep(time.Minute)
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

// Função para verificar se uma entidade já foi processada com base no banco de dados
func isProcessed(db *sql.DB, EntitiesName string) bool {
	// Consulta SQL para verificar se a entidade já existe na tabela
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM public.imported_documents WHERE unnest(xpath('/Jobs/JobPortals/JobPortal', xml)) = $1)", EntitiesName).Scan(&exists)
	if err != nil {
		log.Printf("Erro ao verificar se a entidade já foi processada: %v", err)
		return false
	}

	//A entidade existe na tabela (já foi processada)
	return exists
}

func main() {
	connection()
}
