package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"seeyouvideo/src/usuarios"
)

// repositorioFalsoParaAPI es lo mismo que usamos en las pruebas de
// integración de usuarios, pero repetido aquí porque las pruebas de
// cada paquete solo ven lo que hay en su propio paquete.
type repositorioFalsoParaAPI struct {
	usuarios map[string]*usuarios.Usuario
}

func (r *repositorioFalsoParaAPI) Guardar(u *usuarios.Usuario) error {
	r.usuarios[u.Email()] = u
	return nil
}
func (r *repositorioFalsoParaAPI) BuscarPorEmail(email string) (*usuarios.Usuario, error) {
	u, existe := r.usuarios[email]
	if !existe {
		return nil, usuarios.ErrCredenciales
	}
	return u, nil
}
func (r *repositorioFalsoParaAPI) ExisteEmail(email string) (bool, error) {
	_, existe := r.usuarios[email]
	return existe, nil
}

// TestAceptacion_RegistroCompleto simula lo que haría un usuario
// real: manda una petición HTTP POST a /registro con datos válidos,
// y verifica que la respuesta sea 201 Created con el usuario creado.
func TestAceptacion_RegistroCompleto(t *testing.T) {
	repo := &repositorioFalsoParaAPI{usuarios: make(map[string]*usuarios.Usuario)}
	auth := usuarios.NuevoServicioAuth(repo)
	servidor := NuevoServidor(auth, repo, nil)

	cuerpo := []byte(`{"nombre":"Ricardo","email":"ricardo@test.com","password":"12345678","plan":"basico"}`)
	peticion := httptest.NewRequest(http.MethodPost, "/registro", bytes.NewReader(cuerpo))
	grabador := httptest.NewRecorder()

	servidor.Rutas().ServeHTTP(grabador, peticion)

	if grabador.Code != http.StatusCreated {
		t.Fatalf("se esperaba status 201, se obtuvo %d", grabador.Code)
	}

	var respuesta map[string]interface{}
	json.NewDecoder(grabador.Body).Decode(&respuesta)
	if respuesta["email"] != "ricardo@test.com" {
		t.Errorf("email esperado 'ricardo@test.com', se obtuvo '%v'", respuesta["email"])
	}
}

// TestAceptacion_LoginConCredencialesInvalidas verifica que el
// servicio responde 401 cuando el usuario no existe.
func TestAceptacion_LoginConCredencialesInvalidas(t *testing.T) {
	repo := &repositorioFalsoParaAPI{usuarios: make(map[string]*usuarios.Usuario)}
	auth := usuarios.NuevoServicioAuth(repo)
	servidor := NuevoServidor(auth, repo, nil)

	cuerpo := []byte(`{"email":"noexiste@test.com","password":"12345678"}`)
	peticion := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(cuerpo))
	grabador := httptest.NewRecorder()

	servidor.Rutas().ServeHTTP(grabador, peticion)

	if grabador.Code != http.StatusUnauthorized {
		t.Fatalf("se esperaba status 401, se obtuvo %d", grabador.Code)
	}
}
