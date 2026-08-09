package videos

// RepositorioVideos dice qué se puede hacer con los videos,
// sin amarrarse a MySQL directo (mismo patrón que usamos en usuarios).
type RepositorioVideos interface {
	Guardar(v *Video) error
	ListarTodos() ([]*Video, error)
}
