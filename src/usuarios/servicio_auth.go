package usuarios

import "errors"

var ErrEmailYaRegistrado = errors.New("ya existe una cuenta con ese email")

// ServicioAuth agrupa la LÓGICA de negocio del login/registro. Nota que
// depende de la interfaz RepositorioUsuarios, no de RepositorioMySQL
// directamente — así este código no sabe (ni le importa) si los datos
// están en MySQL, en memoria, o en otro lado.
type ServicioAuth struct {
	repo RepositorioUsuarios
}

func NuevoServicioAuth(repo RepositorioUsuarios) *ServicioAuth {
	return &ServicioAuth{repo: repo}
}

// Registrar crea una cuenta nueva. Antes de crear el Usuario, valida
// que el email no esté ya registrado.
func (s *ServicioAuth) Registrar(nombre, email, password, plan string) (*Usuario, error) {
	existe, err := s.repo.ExisteEmail(email)
	if err != nil {
		return nil, err
	}
	if existe {
		return nil, ErrEmailYaRegistrado
	}

	// NuevoUsuario ya valida nombre/email/password/plan internamente.
	usuario, err := NuevoUsuario(nombre, email, password, plan)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Guardar(usuario); err != nil {
		return nil, err
	}
	return usuario, nil
}

// Login busca al usuario por email y verifica que la contraseña coincida.
func (s *ServicioAuth) Login(email, password string) (*Usuario, error) {
	usuario, err := s.repo.BuscarPorEmail(email)
	if err != nil {
		return nil, err
	}
	if !usuario.VerificarPassword(password) {
		return nil, ErrCredenciales
	}
	return usuario, nil
}
