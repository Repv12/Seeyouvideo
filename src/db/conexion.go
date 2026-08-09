package db

import (
	"database/sql"
	"fmt"

	// El guion bajo "_" significa: importar el paquete solo para que
	// registre el driver de MySQL, aunque no use su nombre directamente.
	_ "github.com/go-sql-driver/mysql"
)

// Config agrupa los datos de conexión a MySQL.
// Encapsular esto en un struct evita pasar 5 parámetros sueltos a Conectar().
type Config struct {
	Usuario   string
	Password  string
	Host      string
	Puerto    string
	BaseDatos string
}

// Conectar abre la conexión a MySQL y crea las tablas si no existen.
// Devuelve un puntero a sql.DB (la conexión) o un error si algo falla.
func Conectar(cfg Config) (*sql.DB, error) {
	// DSN = Data Source String: el formato que espera el driver de MySQL
	// para saber a qué servidor, usuario y base conectarse.
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.Usuario, cfg.Password, cfg.Host, cfg.Puerto, cfg.BaseDatos)

	conexion, err := sql.Open("mysql", dsn)
	if err != nil {
		// %w envuelve el error original, así no perdemos el detalle real
		// de qué falló (útil para debug).
		return nil, fmt.Errorf("error abriendo conexión: %w", err)
	}

	// sql.Open no verifica la conexión de inmediato — Ping() sí la prueba
	// realmente contra el servidor.
	if err := conexion.Ping(); err != nil {
		return nil, fmt.Errorf("no se pudo conectar a MySQL: %w", err)
	}

	if err := crearTablas(conexion); err != nil {
		return nil, fmt.Errorf("error creando tablas: %w", err)
	}

	return conexion, nil
}

// crearTablas define el esquema de la base de datos. Se ejecuta cada vez
// que arranca el programa, pero "if not exists" evita duplicar tablas
// si ya existen.
func crearTablas(conexion *sql.DB) error {
	// separamos en dos Exec porque el driver no deja mandar
	// varias sentencias SQL juntas de una
	_, err := conexion.Exec(`
	CREATE TABLE IF NOT EXISTS usuarios (
		id INT AUTO_INCREMENT PRIMARY KEY,
		nombre VARCHAR(100) NOT NULL,
		email VARCHAR(150) NOT NULL UNIQUE,
		password_hash VARCHAR(255) NOT NULL,
		plan ENUM('basico', 'estandar', 'premium') NOT NULL DEFAULT 'basico',
		creado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		return err
	}

	_, err = conexion.Exec(`
	CREATE TABLE IF NOT EXISTS videos (
		id INT AUTO_INCREMENT PRIMARY KEY,
		titulo VARCHAR(200) NOT NULL,
		categoria VARCHAR(100) NOT NULL,
		url VARCHAR(500) NOT NULL,
		id_usuario INT NOT NULL,
		creado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (id_usuario) REFERENCES usuarios(id)
	);`)
	return err
}
