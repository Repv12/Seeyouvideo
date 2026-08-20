package usuarios

import "testing"

// repositorioFalso es una implementación en memoria de
// RepositorioUsuarios, solo para pruebas — así no necesitamos
// MySQL real corriendo para probar la lógica de negocio.
type repositorioFalso struct {
	usuarios map[string]*Usuario
}

func nuevoRepositorioFalso() *repositorioFalso {
	return &repositorioFalso{usuarios: make(map[string]*Usuario)}
}

func (r *repositorioFalso) Guardar(u *Usuario) error {
	r.usuarios[u.Email()] = u
	return nil
}

func (r *repositorioFalso) BuscarPorEmail(email string) (*Usuario, error) {
	u, existe := r.usuarios[email]
	if !existe {
		return nil, ErrCredenciales
	}
	return u, nil
}

func (r *repositorioFalso) ExisteEmail(email string) (bool, error) {
	_, existe := r.usuarios[email]
	return existe, nil
}

// TestServicioAuth_RegistrarYLogin prueba el flujo completo:
// registrar un usuario y después poder loguearse con esos datos.
func TestServicioAuth_RegistrarYLogin(t *testing.T) {
	repo := nuevoRepositorioFalso()
	auth := NuevoServicioAuth(repo)

	_, err := auth.Registrar("Ricardo", "ricardo@test.com", "12345678", "basico")
	if err != nil {
		t.Fatalf("no se esperaba error al registrar: %v", err)
	}

	usuario, err := auth.Login("ricardo@test.com", "12345678")
	if err != nil {
		t.Fatalf("no se esperaba error al hacer login: %v", err)
	}
	if usuario.Email() != "ricardo@test.com" {
		t.Errorf("email esperado 'ricardo@test.com', se obtuvo '%s'", usuario.Email())
	}
}

// TestServicioAuth_EmailDuplicado prueba que no se puede registrar
// dos veces el mismo correo.
func TestServicioAuth_EmailDuplicado(t *testing.T) {
	repo := nuevoRepositorioFalso()
	auth := NuevoServicioAuth(repo)

	auth.Registrar("Ricardo", "ricardo@test.com", "12345678", "basico")
	_, err := auth.Registrar("Otro Nombre", "ricardo@test.com", "87654321", "premium")

	if err != ErrEmailYaRegistrado {
		t.Errorf("se esperaba ErrEmailYaRegistrado, se obtuvo: %v", err)
	}
}

// TestServicioAuth_LoginPasswordIncorrecta prueba que el login
// falla si la contraseña no coincide.
func TestServicioAuth_LoginPasswordIncorrecta(t *testing.T) {
	repo := nuevoRepositorioFalso()
	auth := NuevoServicioAuth(repo)

	auth.Registrar("Ricardo", "ricardo@test.com", "12345678", "basico")
	_, err := auth.Login("ricardo@test.com", "claveincorrecta")

	if err != ErrCredenciales {
		t.Errorf("se esperaba ErrCredenciales, se obtuvo: %v", err)
	}
}
