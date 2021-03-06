package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type mergeCompany struct {
	similarity     float64
	idArq          string
	nmRazaoSocial  string
	nmCidade       string
	nmEstado       string
	nmRazaoSocial2 string
	cdCnpj         string
}

var db *sql.DB

func initConection() {
	password := os.Getenv("POSTGRES_USER")
	user := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbname := os.Getenv("POSTGRES_DB")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	fmt.Println(psqlInfo)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		time.Sleep(1 * time.Second)
		initConection()
	}

	if err = db.Ping(); err != nil {
		time.Sleep(1 * time.Second)
		initConection()
	}
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Docker")
	})
	initConection()

	for i := 1; i < len(os.Args); i++ {
		readCsv(os.Args[i])
	}

	similarityTable()

	fmt.Println("Listening on :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func similarityTable() {
	fmt.Println("Espere similaridade entre os dois campos textuais")
	command := `CREATE EXTENSION pg_trgm;
	SELECT similarity(n1.nm_razao_social, n2.nm_razao_social) AS sim,n1.id_arq, n1.nm_razao_social, n1.nm_cidade, n1.nm_estado, n2.nm_razao_social, n2.cd_cnpj
	FROM   empresainfo1 n1
	JOIN   empresainfo2 n2 ON n1.nm_razao_social <> n2.nm_razao_social
				   AND n1.nm_razao_social % n2.nm_razao_social AND similarity(n1.nm_razao_social, n2.nm_razao_social) > 0.5
	ORDER  BY sim DESC;
	`
	rows, err := db.Query(command)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var (
			similarity     float64
			idArq          string
			nmRazaoSocial  string
			nmCidade       string
			nmEstado       string
			nmRazaoSocial2 string
			cdCnpj         string
		)
		if err := rows.Scan(&similarity, &idArq, &nmRazaoSocial, &nmCidade, &nmEstado, &nmRazaoSocial2, &cdCnpj); err != nil {
			log.Fatal(err)
		}
		company := mergeCompany{similarity, idArq, nmRazaoSocial, nmCidade, nmEstado, nmRazaoSocial2, cdCnpj}
		insertMergeComapy(company)
	}
	fmt.Println("Terminado os insert entre as similaridades das duas tabelas company")
}

func readCsv(fileName string) {
	file, _ := os.Open(fileName)
	r := csv.NewReader(bufio.NewReader(file))

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		var lastInsertId int
		if len(record) == 2 {
			err = db.QueryRow("INSERT INTO empresainfo2(cd_cnpj,nm_razao_social) VALUES($1,$2) returning em_id;", record[0], record[1]).Scan(&lastInsertId)
			if err != nil {
				fmt.Println("Error insert element ", err)
			}
			fmt.Println("Sucess insert ", lastInsertId)
		} else if len(record) == 4 {
			err = db.QueryRow("INSERT INTO empresainfo1(id_arq, nm_razao_social,nm_cidade, nm_estado) VALUES($1,$2,$3,$4) returning em_id;", record[0], record[1], record[2], record[3]).Scan(&lastInsertId)
			if err != nil {
				fmt.Println("Error insert element ", err)
			}
			fmt.Println("Sucess insert ", lastInsertId)
		}
	}
}

func insertMergeComapy(company mergeCompany) {
	var lastInsertId int
	err := db.QueryRow(`INSERT INTO merge_table_company(id_arq, nm_razao_social,
										nm_cidade, nm_estado, nm_razao_social2, cd_cnpj) VALUES
										($1, $2, $3, $4, $5, $6) returning id;`, company.idArq,
		company.nmRazaoSocial, company.nmCidade, company.nmEstado,
		company.nmRazaoSocial2, company.cdCnpj).Scan(&lastInsertId)

	if err != nil {
		fmt.Println("Error insert element ", err)
	}
	fmt.Println("Sucess insert ", lastInsertId)
}
