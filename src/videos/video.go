package videos

import (
	"errors"
	"strings"
)

// Video es el dato de cada video que se sube al catálogo.
// Igual que con Usuario, los campos van en minúscula pa que
// nadie los toque directo desde afuera, solo por los métodos de aquí.
type Video struct {
	id        int
	titulo    string
	categoria string
	url       string
	idUsuario int
}

var (
	ErrTituloVacio    = errors.New("el título no puede estar vacío")
	ErrCategoriaVacia = errors.New("la categoría no puede estar vacía")
	ErrURLInvalida    = errors.New("la url/enlace no puede estar vacío")
)

// NuevoVideo es el constructor, valida antes de armar el objeto.
// idUsuario es quién lo subió (así sabemos de quién es cada video).
func NuevoVideo(titulo, categoria, url string, idUsuario int) (*Video, error) {
	if strings.TrimSpace(titulo) == "" {
		return nil, ErrTituloVacio
	}
	if strings.TrimSpace(categoria) == "" {
		return nil, ErrCategoriaVacia
	}
	if strings.TrimSpace(url) == "" {
		return nil, ErrURLInvalida
	}

	return &Video{
		titulo:    titulo,
		categoria: categoria,
		url:       url,
		idUsuario: idUsuario,
	}, nil
}

// getters, solo lectura como en Usuario
func (v *Video) ID() int           { return v.id }
func (v *Video) Titulo() string    { return v.titulo }
func (v *Video) Categoria() string { return v.categoria }
func (v *Video) URL() string       { return v.url }
func (v *Video) IDUsuario() int    { return v.idUsuario }

// esto lo usa el repositorio nomas, pa setear el id que da mysql
func (v *Video) SetID(id int) {
	v.id = id
}
