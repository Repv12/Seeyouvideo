package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"fmt"
	"seeyouvideo/src/usuarios"
	"seeyouvideo/src/videos"
)

// Servidor agrupa todo lo que los servicios web necesitan: el
// servicio de login/registro y el repositorio de videos.
type Servidor struct {
	auth         *usuarios.ServicioAuth
	repoUsuarios usuarios.RepositorioUsuarios
	repoVideos   *videos.RepositorioVideosMySQL
}

func NuevoServidor(auth *usuarios.ServicioAuth, repoUsuarios usuarios.RepositorioUsuarios, repoVideos *videos.RepositorioVideosMySQL) *Servidor {
	return &Servidor{auth: auth, repoUsuarios: repoUsuarios, repoVideos: repoVideos}
}

// responderJSON es un helper: convierte cualquier dato a JSON y lo
// manda como respuesta, así no repito esto en cada servicio.
func responderJSON(w http.ResponseWriter, status int, datos interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(datos)
}

func responderError(w http.ResponseWriter, status int, mensaje string) {
	responderJSON(w, status, map[string]string{"error": mensaje})
}

// Rutas registra los 8 servicios web en el enrutador de Go.
func (s *Servidor) Rutas() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", s.paginaPrincipal)
	mux.HandleFunc("POST /registro", s.registrar)                             // 1
	mux.HandleFunc("POST /login", s.login)                                    // 2
	mux.HandleFunc("GET /usuarios/{id}", s.obtenerUsuario)                    // 3
	mux.HandleFunc("PUT /usuarios/{id}/plan", s.cambiarPlan)                  // 4
	mux.HandleFunc("GET /videos", s.listarVideos)                             // 5
	mux.HandleFunc("GET /videos/{id}", s.obtenerVideo)                        // 6
	mux.HandleFunc("POST /videos", s.agregarVideo)                            // 7
	mux.HandleFunc("GET /videos/categoria/{categoria}", s.videosPorCategoria) // 8

	return mux
}

// --- 1. Registro ---
func (s *Servidor) registrar(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Nombre   string `json:"nombre"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Plan     string `json:"plan"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		responderError(w, http.StatusBadRequest, "cuerpo inválido")
		return
	}

	usuario, err := s.auth.Registrar(body.Nombre, body.Email, body.Password, body.Plan)
	if err != nil {
		responderError(w, http.StatusBadRequest, err.Error())
		return
	}
	responderJSON(w, http.StatusCreated, map[string]interface{}{
		"id":     usuario.ID(),
		"nombre": usuario.Nombre(),
		"email":  usuario.Email(),
		"plan":   usuario.Plan(),
	})
}

// --- 2. Login ---
func (s *Servidor) login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		responderError(w, http.StatusBadRequest, "cuerpo inválido")
		return
	}

	usuario, err := s.auth.Login(body.Email, body.Password)
	if err != nil {
		responderError(w, http.StatusUnauthorized, err.Error())
		return
	}
	responderJSON(w, http.StatusOK, map[string]interface{}{
		"id":     usuario.ID(),
		"nombre": usuario.Nombre(),
		"plan":   usuario.Plan(),
	})
}

// --- 3. Ver perfil de usuario ---
func (s *Servidor) obtenerUsuario(w http.ResponseWriter, r *http.Request) {
	// como todavía no tenemos búsqueda por ID, identificamos al
	// usuario por email vía query param: /usuarios/1?email=correo
	email := r.URL.Query().Get("email")
	if email == "" {
		responderError(w, http.StatusBadRequest, "usa ?email=correo para identificar al usuario")
		return
	}
	usuario, err := s.repoUsuarios.BuscarPorEmail(email)
	if err != nil {
		responderError(w, http.StatusNotFound, err.Error())
		return
	}
	responderJSON(w, http.StatusOK, map[string]interface{}{
		"id":     usuario.ID(),
		"nombre": usuario.Nombre(),
		"email":  usuario.Email(),
		"plan":   usuario.Plan(),
	})
}

// --- 4. Cambiar plan ---
func (s *Servidor) cambiarPlan(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	var body struct {
		Plan string `json:"plan"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		responderError(w, http.StatusBadRequest, "cuerpo inválido")
		return
	}

	usuario, err := s.repoUsuarios.BuscarPorEmail(email)
	if err != nil {
		responderError(w, http.StatusNotFound, err.Error())
		return
	}
	if err := usuario.SetPlan(body.Plan); err != nil {
		responderError(w, http.StatusBadRequest, err.Error())
		return
	}
	responderJSON(w, http.StatusOK, map[string]string{"mensaje": "plan actualizado a " + usuario.Plan()})
}

// --- 5. Listar todos los videos ---
func (s *Servidor) listarVideos(w http.ResponseWriter, r *http.Request) {
	lista, err := s.repoVideos.ListarTodos()
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responderJSON(w, http.StatusOK, videosAJSON(lista))
}

// --- 6. Ver un video específico ---
func (s *Servidor) obtenerVideo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responderError(w, http.StatusBadRequest, "id inválido")
		return
	}
	lista, err := s.repoVideos.ListarTodos()
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	for _, v := range lista {
		if v.ID() == id {
			responderJSON(w, http.StatusOK, videoAJSON(v))
			return
		}
	}
	responderError(w, http.StatusNotFound, "video no encontrado")
}

// --- 7. Agregar video ---
func (s *Servidor) agregarVideo(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Titulo    string `json:"titulo"`
		Categoria string `json:"categoria"`
		URL       string `json:"url"`
		IDUsuario int    `json:"id_usuario"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		responderError(w, http.StatusBadRequest, "cuerpo inválido")
		return
	}

	video, err := videos.NuevoVideo(body.Titulo, body.Categoria, body.URL, body.IDUsuario)
	if err != nil {
		responderError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.repoVideos.Guardar(video); err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responderJSON(w, http.StatusCreated, videoAJSON(video))
}

// --- 8. Videos por categoría ---
func (s *Servidor) videosPorCategoria(w http.ResponseWriter, r *http.Request) {
	categoria := r.PathValue("categoria")
	lista, err := s.repoVideos.ListarTodos()
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var filtrados []*videos.Video
	for _, v := range lista {
		if strings.EqualFold(v.Categoria(), categoria) {
			filtrados = append(filtrados, v)
		}
	}
	responderJSON(w, http.StatusOK, videosAJSON(filtrados))
}

// --- helpers de conversión a JSON ---
func videoAJSON(v *videos.Video) map[string]interface{} {
	return map[string]interface{}{
		"id":        v.ID(),
		"titulo":    v.Titulo(),
		"categoria": v.Categoria(),
		"url":       v.URL(),
	}
}

func videosAJSON(lista []*videos.Video) []map[string]interface{} {
	resultado := make([]map[string]interface{}, 0, len(lista))
	for _, v := range lista {
		resultado = append(resultado, videoAJSON(v))
	}
	return resultado
}

// extraerIDYoutube saca el ID del video de una URL de YouTube, para
// poder armar el link de la miniatura. Soporta los dos formatos más
// comunes: youtube.com/watch?v=XXXX y youtu.be/XXXX
func extraerIDYoutube(urlVideo string) string {
	if strings.Contains(urlVideo, "v=") {
		partes := strings.Split(urlVideo, "v=")
		id := partes[1]
		if amp := strings.Index(id, "&"); amp != -1 {
			id = id[:amp]
		}
		return id
	}
	if strings.Contains(urlVideo, "youtu.be/") {
		partes := strings.Split(urlVideo, "youtu.be/")
		id := partes[1]
		if q := strings.Index(id, "?"); q != -1 {
			id = id[:q]
		}
		return id
	}
	return ""
}

// paginaPrincipal sirve una página web sencilla que muestra el
// catálogo, consumiendo los mismos datos que el servicio /videos.
func (s *Servidor) paginaPrincipal(w http.ResponseWriter, r *http.Request) {
	lista, err := s.repoVideos.ListarTodos()
	if err != nil {
		http.Error(w, "error cargando videos", http.StatusInternalServerError)
		return
	}

	html := `<!DOCTYPE html>
<html lang="es">
<head>
<meta charset="UTF-8">
<title>SeeUvideo</title>
<style>
	body {
		background: #0a0c1e;
		color: white;
		font-family: sans-serif;
		padding: 40px;
		background-image: radial-gradient(white 1px, transparent 1px);
		background-size: 40px 40px;
	}
	h1 { font-size: 32px; }
	.video {
		background: #1e2346;
		padding: 12px 20px;
		border-radius: 10px;
		margin-bottom: 12px;
		display: flex;
		align-items: center;
		gap: 15px;
	}
	.miniatura {
		width: 120px;
		height: 68px;
		object-fit: cover;
		border-radius: 6px;
	}
	.info { flex: 1; }
	.video a {
		background: #7888ff;
		color: white;
		padding: 8px 16px;
		border-radius: 6px;
		text-decoration: none;
		font-weight: bold;
	}
	.categoria { color: #aab; font-size: 14px; }
</style>
</head>
<body>
	<h1>SeeUvideo</h1>`

	for _, v := range lista {
		idYoutube := extraerIDYoutube(v.URL())
		miniatura := "https://img.youtube.com/vi/" + idYoutube + "/hqdefault.jpg"

		html += fmt.Sprintf(`
	<div class="video">
		<img src="%s" class="miniatura" alt="miniatura">
		<div class="info">
			🎬 <strong>%s</strong>
			<div class="categoria">%s</div>
		</div>
		<a href="%s" target="_blank">Ver</a>
	</div>`, miniatura, v.Titulo(), v.Categoria(), v.URL())
	}

	html += `
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}
