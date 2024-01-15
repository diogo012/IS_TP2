package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

const connStr = "user=is password=is dbname=is sslmode=disable host=db-xml port=5432"

type Company struct {
	ID         string `xml:"id,attr"`
	Name       string `xml:"company,attr"`
	Size       string `xml:"companySize,attr"`
	CountryRef string `xml:"country_ref,attr"`
	Benefits   string `xml:"Benefits"`
}

type Country struct {
	ID        string `xml:"id,attr"`
	Name      string `xml:"country,attr"`
	Location  string `xml:"location,attr"`
	Latitude  string `xml:"latitude,attr"`
	Longitude string `xml:"longitude,attr"`
}

type Job struct {
	ID             string `xml:"id,attr"`
	JobTitle       string `xml:"jobTitle,attr"`
	Experience     string `xml:"experience,attr"`
	WorkType       string `xml:"workType,attr"`
	Qualifications string `xml:"qualifications,attr"`
	Preference     string `xml:"preference,attr"`
	JobPostingDate string `xml:"jobPostingDate,attr"`
	PersonRef      string `xml:"person_ref,attr"`
	RoleRef        string `xml:"role_ref,attr"`
	Description    string `xml:"Description"`
	Skills         string `xml:"Skills"`
}

type JobPortal struct {
	//XMLName   xml.Name `xml:"JobPortal"`
	ID        string `xml:"id,attr"`
	JobPortal string `xml:"jobPortal,attr"`
	Jobs      []Job  `xml:"Jobs>Job"`
}

func parseXML(xmlData string) (interface{}, error) {
	decoder := xml.NewDecoder(strings.NewReader(xmlData))

	var entity interface{}

	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		switch se := token.(type) {
		case xml.StartElement:
			switch se.Name.Local {
			case "JobPortal":
				var jobPortal JobPortal
				if err := decoder.DecodeElement(&jobPortal, &se); err != nil {
					return nil, err
				}
				entity = jobPortal
			case "Company":
				var company Company
				if err := decoder.DecodeElement(&company, &se); err != nil {
					return nil, err
				}
				entity = company
			case "Country":
				var country Country
				if err := decoder.DecodeElement(&country, &se); err != nil {
					return nil, err
				}
				entity = country
			}
		}
	}

	return entity, nil
}

func getUnprocessedEntities(db *sql.DB) error {
	rows, err := db.Query("SELECT unnest(xpath('/Jobs/JobPortals/JobPortal|Jobs/Companies/Company|Jobs/Countries/Country', xml)) AS entity_attributes FROM public.imported_documents;")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var xmlContent string
		if err := rows.Scan(&xmlContent); err != nil {
			log.Println("Erro ao obter dados da entidade:", err)
			continue
		}

		fmt.Printf("Nova entidade encontrada: %s\n", xmlContent)

		entity, err := parseXML(xmlContent)
		if err != nil {
			log.Printf("Erro ao parsear XML: %v", err)
			continue
		}

		// Agora você tem a estrutura preenchida com os dados do XML
		// Faça o que precisar com a estrutura, por exemplo, enviar para o broker
		err = sendToBroker(fmt.Sprintf("%T", entity), entity)
		if err != nil {
			log.Printf("Erro ao enviar entidade para o broker: %v", err)
		}

		// Determine o tipo da entidade e aja de acordo
		switch e := entity.(type) {
		case JobPortal:
			// Faça algo específico para JobPortal
			fmt.Printf("Entidade JobPortal encontrada. ID: %s\n", e.ID)
		case Company:
			// Faça algo específico para Company
			fmt.Printf("Entidade Company encontrada. Nome: %s\n", e.Name)
		case Country:
			// Faça algo específico para Country
			fmt.Printf("Entidade Country encontrada. País: %s\n", e.Name)
		default:
			log.Printf("Tipo de entidade não reconhecido: %v", e)
		}
	}

	return nil
}

func sendToBroker(entityType string, entity interface{}) error {
	conn, err := amqp.Dial("amqp://is:is@rabbitmq:5672/is")
	if err != nil {
		return fmt.Errorf("Falha ao conectar ao RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("Falha ao abrir um canal: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"queue_migrator",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("Falha ao declarar uma fila: %v", err)
	}

	// Converter a entidade para JSON
	entityJSON, err := json.Marshal(entity)
	if err != nil {
		return fmt.Errorf("Falha ao converter entidade para JSON: %v", err)
	}

	messageBody := fmt.Sprintf("Entity type: %s, Entity: %s", entityType, entityJSON)

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(messageBody),
		})
	if err != nil {
		return fmt.Errorf("Falha ao publicar uma mensagem: %v", err)
	}

	fmt.Printf("Mensagem enviada para o serviço broker\n")
	return nil
}

func main() {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Ping failed: %v", err)
	} else {
		fmt.Println("Connection successful!")
	}

	for {
		err := getUnprocessedEntities(db)
		if err != nil {
			log.Println("Erro ao obter novas entidades:", err)
		}

		time.Sleep(10 * time.Minute)
	}
}
