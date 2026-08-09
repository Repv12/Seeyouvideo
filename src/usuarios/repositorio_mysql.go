package usuarios

import (
	"database/sql"
	"errors"
	"fmt"
)

// RepositorioMySQL implementa la interfaz RepositorioUsuarios usando una
// base de datos MySQL real. Es la ÚNICA parte del sistema que sabe
// escribir SQL — el resto del código (ServicioAuth, main.go) solo conoce
// la interfaz, no estos detalles.
type RepositorioMySQL struct {
	conexion *sql.DB
}

// NuevoRepositorioMySQL construye el repositorio a partir de una conexión
// ya abierta (la que crea db.Conectar()).
func NuevoRepositorioMySQL(conexion *sql.DB) *RepositorioMySQL {
	return &RepositorioMySQL{conexion: conexion}
}

// Guardar inserta un nuevo usuario en la tabla. Los "?" son placeholders:
// evitan inyección SQL, porque los valores se pasan por separado, no
// concatenados directamente en el texto de la query.
func (r *RepositorioMySQL) Guardar(u *Usuario) error {
	query := `INSERT INTO usuarios (nombre, email, password_hash, plan) VALUES (?, ?, ?, ?)`
	resultado, err := r.conexion.Exec(query, u.Nombre(), u.Email(), u.PasswordHash(), u.Plan())
	if err != nil {
		return fmt.Errorf("error guardando usuario: %w", err)
	}

	// MySQL genera el ID automáticamente (AUTO_INCREMENT). Lo recuperamos
	// y lo guardamos en el objeto Usuario para tenerlo disponible después.
	id, err := resultado.LastInsertId()
	if err == nil {
		u.SetID(int(id))
	}
	return nil
}

// BuscarPorEmail se usa en el login: busca un usuario por su email y
// reconstruye el struct Usuario a partir de la fila de la base de datos.
func (r *RepositorioMySQL) BuscarPorEmail(email string) (*Usuario, error) {
	query := `SELECT id, nombre, email, password_hash, plan FROM usuarios WHERE email = ?`
	fila := r.conexion.QueryRow(query, email)

	var u Usuario
	// Scan copia cada columna de la fila al campo correspondiente del struct.
	err := fila.Scan(&u.id, &u.nombre, &u.email, &u.passwordHash, &u.plan)
	if errors.Is(err, sql.ErrNoRows) {
		// Por seguridad, no decimos "ese email no existe" —
		// se usa el mismo error genérico que una contraseña incorrecta.
		return nil, ErrCredenciales
	}
	if err != nil {
		return nil, fmt.Errorf("error buscando usuario: %w", err)
	}
	return &u, nil
}

// ExisteEmail se usa antes de registrar, para no permitir emails duplicados.
func (r *RepositorioMySQL) ExisteEmail(email string) (bool, error) {
	query := `SELECT COUNT(*) FROM usuarios WHERE email = ?`
	var contador int
	err := r.conexion.QueryRow(query, email).Scan(&contador)
	if err != nil {
		return false, fmt.Errorf("error verificando email: %w", err)
	}
	return contador > 0, nil
}
