package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sparkedhost/go-mysqldump"
)

func main() {
	// Definir flags de línea de comandos
	var (
		host     = flag.String("host", "localhost", "MySQL host")
		port     = flag.Int("port", 3306, "MySQL port")
		username = flag.String("user", "", "MySQL username (required)")
		password = flag.String("password", "", "MySQL password")
		database = flag.String("database", "", "Database name (required)")
		output   = flag.String("output", "", "Output file (optional, defaults to stdout)")
		dir      = flag.String("dir", ".", "Directory for dump file (when using Register)")
		filename = flag.String("filename", "dump", "Filename for dump (when using Register)")
	)

	flag.Parse()

	// Validar parámetros requeridos
	if *username == "" || *database == "" {
		fmt.Println("Error: username and database are required")
		flag.Usage()
		os.Exit(1)
	}

	// Crear string de conexión MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", *username, *password, *host, *port, *database)

	// Conectar a la base de datos
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Probar conexión
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	fmt.Printf("Connected to MySQL database '%s' successfully!\n", *database)

	if *output != "" {
		// Usar Register para crear archivo de dump
		fmt.Printf("Creating dump file: %s/%s.sql\n", *dir, *filename)

		dumper, err := mysqldump.Register(db, *dir, *filename)
		if err != nil {
			log.Fatalf("Error registering dumper: %v", err)
		}
		defer dumper.Close()

		err = dumper.Dump()
		if err != nil {
			log.Fatalf("Error during dump: %v", err)
		}

		fmt.Printf("Dump completed successfully! File saved as: %s/%s.sql\n", *dir, *filename)
	} else {
		// Usar Dump directo a stdout
		fmt.Println("Dumping to stdout...")
		fmt.Println("-- MySQL Dump Output:")

		err = mysqldump.Dump(db, os.Stdout)
		if err != nil {
			log.Fatalf("Error during dump: %v", err)
		}
		fmt.Println("\n-- Dump completed successfully!")
	}
}
