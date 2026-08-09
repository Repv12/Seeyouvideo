package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"seeyouvideo/src/db"
	"seeyouvideo/src/usuarios"
	"seeyouvideo/src/videos"
)

// usuarioActual guarda la sesión activa, pa que las demás pantallas
// sepan quién está logueado ahorita.
var usuarioActual *usuarios.Usuario

func main() {
	// conececto el mysql
	cfg := db.Config{
		Usuario:   "root",
		Password:  "0904",
		Host:      "127.0.0.1",
		Puerto:    "3306",
		BaseDatos: "streaming_db",
	}

	conexion, err := db.Conectar(cfg)
	if err != nil {
		fmt.Println("Error conectando a MySQL:", err)
		return
	}
	defer conexion.Close()

	repoUsuarios := usuarios.NuevoRepositorioMySQL(conexion)
	authServicio := usuarios.NuevoServicioAuth(repoUsuarios)

	repoVideos := videos.NuevoRepositorioVideosMySQL(conexion)

	a := app.New()
	a.Settings().SetTheme(&temaOscuro{}) // nuestro tema custom cielo nocturno

	w := a.NewWindow("SeeUvideo")
	w.Resize(fyne.NewSize(1000, 700))
	w.CenterOnScreen()

	mostrarPantallaLogin(w, a, authServicio, repoVideos)

	w.ShowAndRun()
}

// --- TEMA CUSTOM: cielo nocturno ---
// Fyne permite meter un tema propio implementando la interfaz fyne.Theme.
// Aquí solo cambiamos los colores clave, el resto lo dejamos por default.
type temaOscuro struct{}

func (temaOscuro) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 10, G: 12, B: 30, A: 255} // azul casi negro, tipo espacio
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 120, G: 140, B: 255, A: 255} // azul estrella
	case theme.ColorNameButton:
		return color.NRGBA{R: 30, G: 35, B: 70, A: 255}
	case theme.ColorNameForeground:
		return color.White
	}
	return theme.DefaultTheme().Color(name, variant)
}
func (temaOscuro) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}
func (temaOscuro) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}
func (temaOscuro) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

// fondoEstrellado genera puntitos random en el fondo, pa dar el efecto
// de universo. Es un truco simple: un montón de círculos chiquitos.
func fondoEstrellado(ancho, alto float32) *fyne.Container {
	fondo := container.NewWithoutLayout()
	for i := 0; i < 80; i++ {
		estrella := canvas.NewCircle(color.White)
		estrella.Resize(fyne.NewSize(2, 2))
		x := rand.Float32() * ancho
		y := rand.Float32() * alto
		estrella.Move(fyne.NewPos(x, y))
		fondo.Add(estrella)
	}
	return fondo
}

// --- PANTALLA DE LOGIN ---
func mostrarPantallaLogin(w fyne.Window, a fyne.App, auth *usuarios.ServicioAuth, repoVideos *videos.RepositorioVideosMySQL) {
	titulo := canvas.NewText("SeeUvideo", color.White)
	titulo.TextSize = 32
	titulo.Alignment = fyne.TextAlignCenter
	titulo.TextStyle = fyne.TextStyle{Bold: true}

	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("Email")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Contraseña")

	mensaje := widget.NewLabel("")
	mensaje.Wrapping = fyne.TextWrapWord

	botonLogin := widget.NewButton("Iniciar sesión", func() {
		usuario, err := auth.Login(emailEntry.Text, passwordEntry.Text)
		if err != nil {
			mensaje.SetText(err.Error())
			return
		}
		usuarioActual = usuario
		mostrarPantallaCatalogo(w, a, repoVideos)
	})

	botonIrRegistro := widget.NewButton("¿No tienes cuenta? Regístrate", func() {
		mostrarPantallaRegistro(w, a, auth, repoVideos)
	})

	formulario := container.NewVBox(
		titulo,
		widget.NewLabel(""), // espaciador
		emailEntry,
		passwordEntry,
		botonLogin,
		botonIrRegistro,
		mensaje,
	)

	// apilamos el fondo estrellado detrás del formulario
	fondo := fondoEstrellado(1000, 700)
	contenido := container.NewStack(fondo, container.NewPadded(formulario))

	w.SetContent(contenido)
}

// --- PANTALLA DE REGISTRO ---
func mostrarPantallaRegistro(w fyne.Window, a fyne.App, auth *usuarios.ServicioAuth, repoVideos *videos.RepositorioVideosMySQL) {
	titulo := canvas.NewText("Crear cuenta", color.White)
	titulo.TextSize = 28
	titulo.Alignment = fyne.TextAlignCenter
	titulo.TextStyle = fyne.TextStyle{Bold: true}

	nombreEntry := widget.NewEntry()
	nombreEntry.SetPlaceHolder("Nombre")

	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("Email")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Contraseña (mínimo 8 caracteres)")

	// dropdown pa elegir el plan, así no dejamos que escriban cualquier cosa
	planSelect := widget.NewSelect([]string{"basico", "premium"}, func(s string) {})
	planSelect.SetSelected("basico")

	mensaje := widget.NewLabel("")
	mensaje.Wrapping = fyne.TextWrapWord

	botonCrear := widget.NewButton("Crear cuenta", func() {
		usuario, err := auth.Registrar(nombreEntry.Text, emailEntry.Text, passwordEntry.Text, planSelect.Selected)
		if err != nil {
			mensaje.SetText(err.Error())
			return
		}
		usuarioActual = usuario
		mostrarPantallaCatalogo(w, a, repoVideos)
	})

	botonVolver := widget.NewButton("Volver al login", func() {
		mostrarPantallaLogin(w, a, auth, repoVideos)
	})

	formulario := container.NewVBox(
		titulo,
		nombreEntry,
		emailEntry,
		passwordEntry,
		planSelect,
		botonCrear,
		botonVolver,
		mensaje,
	)

	fondo := fondoEstrellado(1000, 700)
	contenido := container.NewStack(fondo, container.NewPadded(formulario))

	w.SetContent(contenido)
}

// --- PANTALLA DE CATÁLOGO (después de loguearse) ---
func mostrarPantallaCatalogo(w fyne.Window, a fyne.App, repoVideos *videos.RepositorioVideosMySQL) {
	titulo := canvas.NewText(fmt.Sprintf("Hola, %s", usuarioActual.Nombre()), color.White)
	titulo.TextSize = 24
	titulo.TextStyle = fyne.TextStyle{Bold: true}

	listaVideos := container.NewVBox() // en vez de un Label, usamos un contenedor de botones

	actualizarLista := func() {
		listaVideos.Objects = nil // limpiamos lo que había antes

		lista, err := repoVideos.ListarTodos()
		if err != nil {
			listaVideos.Add(widget.NewLabel("Error cargando videos: " + err.Error()))
			listaVideos.Refresh()
			return
		}
		if len(lista) == 0 {
			listaVideos.Add(widget.NewLabel("Todavía no hay videos. ¡Agrega el primero!"))
			listaVideos.Refresh()
			return
		}

		for _, v := range lista {
			v := v // copia local, pa que cada botón sepa cuál video es el suyo
			fila := container.NewHBox(
				widget.NewLabel(fmt.Sprintf("🎬 %s — %s", v.Titulo(), v.Categoria())),
				widget.NewButton("Ver", func() {
					u, err := url.Parse(v.URL())
					if err == nil {
						a.OpenURL(u)
					}
				}),
			)
			listaVideos.Add(fila)
		}
		listaVideos.Refresh()
	}
	actualizarLista()

	// según el plan, mostramos el botón de agregar o un mensaje invitando a mejorar el plan
	var botonAgregar fyne.CanvasObject
	if usuarioActual.Plan() == "premium" {
		botonAgregar = widget.NewButton("+ Agregar video", func() {
			mostrarDialogoAgregarVideo(w, a, repoVideos, actualizarLista)
		})
	} else {
		botonAgregar = widget.NewLabel("Consigue Premium para subir tus propios videos")
	}
	encabezado := container.NewVBox(
		titulo,
		widget.NewSeparator(),
	)

	scroll := container.NewVScroll(listaVideos)
	scroll.SetMinSize(fyne.NewSize(900, 500))

	contenido := container.NewBorder(
		encabezado,
		botonAgregar,
		nil, nil,
		scroll,
	)

	fondo := fondoEstrellado(1000, 700)
	final := container.NewStack(fondo, container.NewPadded(contenido))

	w.SetContent(final)
}

// --- FORMULARIO PARA AGREGAR VIDEO ---
func mostrarDialogoAgregarVideo(w fyne.Window, a fyne.App, repoVideos *videos.RepositorioVideosMySQL, alGuardar func()) {
	tituloEntry := widget.NewEntry()
	tituloEntry.SetPlaceHolder("Título del video")

	categoriaEntry := widget.NewEntry()
	categoriaEntry.SetPlaceHolder("Categoría")

	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Enlace o URL del video")

	mensaje := widget.NewLabel("")

	botonGuardar := widget.NewButton("Guardar video", func() {
		video, err := videos.NuevoVideo(tituloEntry.Text, categoriaEntry.Text, urlEntry.Text, usuarioActual.ID())
		if err != nil {
			mensaje.SetText(err.Error())
			return
		}
		if err := repoVideos.Guardar(video); err != nil {
			mensaje.SetText(err.Error())
			return
		}
		alGuardar() // refresca la lista del catálogo
		w.SetContent(container.NewVBox(widget.NewLabel("¡Video agregado!")))
		mostrarPantallaCatalogo(w, a, repoVideos)
	})

	botonCancelar := widget.NewButton("Cancelar", func() {
		mostrarPantallaCatalogo(w, a, repoVideos)
	})

	formulario := container.NewVBox(
		widget.NewLabel("Agregar nuevo video"),
		tituloEntry,
		categoriaEntry,
		urlEntry,
		botonGuardar,
		botonCancelar,
		mensaje,
	)

	fondo := fondoEstrellado(1000, 700)
	contenido := container.NewStack(fondo, container.NewPadded(formulario))

	w.SetContent(contenido)
}
