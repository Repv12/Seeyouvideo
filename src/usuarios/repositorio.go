package usuarios

// RepositorioUsuarios define QUÉ operaciones se pueden hacer con usuarios,
// sin decir CÓMO se guardan. Esto es una interfaz: cualquier struct que
// tenga  3 métodos "cumple" la interfaz automáticamente en Go
// (no hace falta declararlo explícitamente como en otros lenguajes).
// La ventaja: el resto del código (como ServicioAuth) solo depende de
// esta interfaz, no de MySQL directamente. Si mañana se cambia de base
// de datos, solo se toca la implementación, no la lógica general.
type RepositorioUsuarios interface {
	Guardar(u *Usuario) error
	BuscarPorEmail(email string) (*Usuario, error)
	ExisteEmail(email string) (bool, error)
}
