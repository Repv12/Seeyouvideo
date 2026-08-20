package usuarios

import "testing"

// TestNuevoUsuario_Valido prueba que un usuario con datos correctos
// se crea sin error.
func TestNuevoUsuario_Valido(t *testing.T) {
	usuario, err := NuevoUsuario("Ricardo", "ricardo@test.com", "12345678", "basico")
	if err != nil {
		t.Fatalf("no se esperaba error, se obtuvo: %v", err)
	}
	if usuario.Nombre() != "Ricardo" {
		t.Errorf("nombre esperado 'Ricardo', se obtuvo '%s'", usuario.Nombre())
	}
}

// TestNuevoUsuario_PasswordCorta prueba que se rechaza una
// contraseña de menos de 8 caracteres.
func TestNuevoUsuario_PasswordCorta(t *testing.T) {
	_, err := NuevoUsuario("Ricardo", "ricardo@test.com", "123", "basico")
	if err != ErrPasswordCorta {
		t.Errorf("se esperaba ErrPasswordCorta, se obtuvo: %v", err)
	}
}

// TestNuevoUsuario_PlanInvalido prueba que un plan que no sea
// basico o premium es rechazado.
func TestNuevoUsuario_PlanInvalido(t *testing.T) {
	_, err := NuevoUsuario("Ricardo", "ricardo@test.com", "12345678", "vip")
	if err != ErrPlanInvalido {
		t.Errorf("se esperaba ErrPlanInvalido, se obtuvo: %v", err)
	}
}

// TestVerificarPassword prueba que la contraseña correcta pasa la
// verificación y una incorrecta no.
func TestVerificarPassword(t *testing.T) {
	usuario, _ := NuevoUsuario("Ricardo", "ricardo@test.com", "12345678", "basico")

	if !usuario.VerificarPassword("12345678") {
		t.Error("la contraseña correcta debería pasar la verificación")
	}
	if usuario.VerificarPassword("otraclave") {
		t.Error("una contraseña incorrecta no debería pasar la verificación")
	}
}
