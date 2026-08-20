package usuarios

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Usuario encapsula los datos de un usuario del sistema.
// Los campos empiezan con minúscula → son "privados": solo se pueden
// leer o modificar desde este mismo paquete (usuarios), nunca directamente
// desde afuera. Por eso más abajo hay funciones Getter/Setter.
type Usuario struct {
	id           int
	nombre       string
	email        string
	passwordHash string
	plan         string
}

// Errores propios del dominio de usuarios. Definirlos como variables
// permite compararlos después con errors.Is() y dar mensajes claros
// en vez de errores genéricos.
var (
	ErrNombreVacio   = errors.New("el nombre no puede estar vacío")
	ErrEmailInvalido = errors.New("el email no tiene un formato válido")
	ErrPasswordCorta = errors.New("la contraseña debe tener al menos 8 caracteres")
	ErrPlanInvalido  = errors.New("el plan debe ser: basico, estandar o premium")
	ErrCredenciales  = errors.New("email o contraseña incorrectos")
)

// Lista blanca de planes válidos, para no aceptar cualquier texto.
var planesValidos = map[string]bool{
	"basico":  true,
	"premium": true,
}

// NuevoUsuario es el CONSTRUCTOR: la única forma "correcta" de crear un
// Usuario. Valida cada dato antes de construir el objeto, y hashea la
// contraseña (nunca se guarda en texto plano). Si algo falla, retorna
// nil + el error correspondiente en vez de un objeto a medio armar.
func NuevoUsuario(nombre, email, password, plan string) (*Usuario, error) {
	if strings.TrimSpace(nombre) == "" {
		return nil, ErrNombreVacio
	}
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return nil, ErrEmailInvalido
	}
	if len(password) < 8 {
		return nil, ErrPasswordCorta
	}
	if !planesValidos[plan] {
		return nil, ErrPlanInvalido
	}

	// bcrypt genera un hash seguro de la contraseña (no reversible).
	// Así, ni siquiera si alguien ve la base de datos puede leer la
	// contraseña real del usuario.
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Usuario{
		nombre:       nombre,
		email:        email,
		passwordHash: string(hash),
		plan:         plan,
	}, nil
}

// --- Getters: acceso de SOLO LECTURA al estado interno. ---
// Permiten leer los datos desde otros paquetes sin exponer los campos
// directamente (que podrían modificarse sin control).

func (u *Usuario) ID() int        { return u.id }
func (u *Usuario) Nombre() string { return u.nombre }
func (u *Usuario) Email() string  { return u.email }
func (u *Usuario) Plan() string   { return u.plan }

// --- Setters: modifican el estado, pero VALIDANDO antes de aceptar el cambio. ---

// SetPlan cambia el plan del usuario, validando que sea uno de los permitidos.
func (u *Usuario) SetPlan(nuevoPlan string) error {
	if !planesValidos[nuevoPlan] {
		return ErrPlanInvalido
	}
	u.plan = nuevoPlan
	return nil
}

// SetID lo usa solo el repositorio, al leer un usuario ya guardado en
// la base de datos (donde ya se conoce su ID autogenerado por MySQL).
func (u *Usuario) SetID(id int) {
	u.id = id
}

// VerificarPassword compara una contraseña en texto plano contra el hash
// guardado sin nunca desencriptar el hash (bcrypt no funciona así,
// solo compara).
func (u *Usuario) VerificarPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.passwordHash), []byte(password))
	return err == nil
}

// PasswordHash expone el hash SOLO para que el repositorio lo guarde en
// la base de datos. Nunca se expone la contraseña real en ningún punto.
func (u *Usuario) PasswordHash() string {
	return u.passwordHash
}
