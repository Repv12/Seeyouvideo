package videos

import (
	"database/sql"
	"fmt"
)

// RepositorioVideosMySQL es el que realmente habla con la base de datos.
type RepositorioVideosMySQL struct {
	conexion *sql.DB
}

func NuevoRepositorioVideosMySQL(conexion *sql.DB) *RepositorioVideosMySQL {
	return &RepositorioVideosMySQL{conexion: conexion}
}

// Guardar mete un video nuevo en la tabla.
func (r *RepositorioVideosMySQL) Guardar(v *Video) error {
	query := `INSERT INTO videos (titulo, categoria, url, id_usuario) VALUES (?, ?, ?, ?)`
	resultado, err := r.conexion.Exec(query, v.Titulo(), v.Categoria(), v.URL(), v.IDUsuario())
	if err != nil {
		return fmt.Errorf("error guardando video: %w", err)
	}

	id, err := resultado.LastInsertId()
	if err == nil {
		v.SetID(int(id))
	}
	return nil
}

// ListarTodos trae todos los videos pa mostrarlos en el catálogo.
func (r *RepositorioVideosMySQL) ListarTodos() ([]*Video, error) {
	query := `SELECT id, titulo, categoria, url, id_usuario FROM videos ORDER BY creado_en DESC`
	filas, err := r.conexion.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error listando videos: %w", err)
	}
	defer filas.Close()

	// vamos armando la lista recorriendo fila por fila
	var lista []*Video
	for filas.Next() {
		var v Video
		if err := filas.Scan(&v.id, &v.titulo, &v.categoria, &v.url, &v.idUsuario); err != nil {
			return nil, fmt.Errorf("error leyendo video: %w", err)
		}
		lista = append(lista, &v)
	}
	return lista, nil
}
