//  Biblioteca Central UNSCH

package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ------------------------------- Dominio ------------------------------------

type TipoUsuario string

const (
	TipoPregrado TipoUsuario = "pregrado"
	TipoPosgrado TipoUsuario = "posgrado"
	TipoDocente  TipoUsuario = "docente"
	TipoAdmin    TipoUsuario = "administrativo"
	TipoExterno  TipoUsuario = "externo"
)

type Usuario struct {
	ID        string
	Nombres   string
	Apellidos string
	Email     string
	Tipo      TipoUsuario
	VenceEn   time.Time // fecha de expiración del carné
}

type Carne struct {
	Numero      string
	Propietario Usuario
	Costo       float64
	EmitidoEn   time.Time
	Contenido   []byte // bytes (por ejemplo, PDF simulado)
}

// --------------------------- Gestor monolítico -------------------------------

type GestorCarne struct {
	// Recursos/estado compartido simulados
	mysqlDSN    string
	printerName string
	smtpHost    string
	smtpUser    string
	smtpPass    string
}

func NuevoGestorCarne() *GestorCarne {
	return &GestorCarne{
		mysqlDSN:    "user:pass@tcp(localhost:3306)/biblioteca",
		printerName: "Printer-Thermal-001",
		smtpHost:    "smtp.unsch.edu.pe:587",
		smtpUser:    "no-reply@unsch.edu.pe",
		smtpPass:    "secreto-super-seguro",
	}
}

// ProcesarSolicitud: flujo secuencial end-to-end
func (g *GestorCarne) ProcesarSolicitud(ctx context.Context, u Usuario) (Carne, error) {
	// 1) Validar datos
	if err := g.ValidarDatos(ctx, u); err != nil {
		return Carne{}, fmt.Errorf("validación fallida: %w", err)
	}
	// 2) Calcular costo según tipo de usuario
	costo, err := g.CalcularCosto(ctx, u)
	if err != nil {
		return Carne{}, fmt.Errorf("cálculo de costo falló: %w", err)
	}
	// 3) Generar carné físico (PDF simulado)
	carne, err := g.GenerarCarne(ctx, u, costo)
	if err != nil {
		return Carne{}, fmt.Errorf("generación de carné falló: %w", err)
	}
	// 4) Guardar datos en MySQL
	if err := g.GuardarEnPostgresql(ctx, u, carne); err != nil {
		return Carne{}, fmt.Errorf("persistencia falló: %w", err)
	}
	// 5) Enviar notificación por email
	if err := g.EnviarNotificacionEmail(ctx, u, carne); err != nil {
		// No crítico, se registra pero no se detiene el flujo
		fmt.Println("[WARN] no se pudo enviar email:", err)
	}
	// 6) Imprimir carné
	if err := g.ImprimirCarne(ctx, carne); err != nil {
		return Carne{}, fmt.Errorf("impresión falló: %w", err)
	}
	return carne, nil
}

// ---------------------------- Métodos acoplados -----------------------------

func (g *GestorCarne) ValidarDatos(ctx context.Context, u Usuario) error {
	if strings.TrimSpace(u.Nombres) == "" || strings.TrimSpace(u.Apellidos) == "" {
		return errors.New("nombres y apellidos son obligatorios")
	}
	if !strings.Contains(u.Email, "@") || len(u.Email) < 5 {
		return errors.New("email inválido")
	}
	switch u.Tipo {
	case TipoPregrado, TipoPosgrado, TipoDocente, TipoAdmin, TipoExterno:
		// OK
	default:
		return fmt.Errorf("tipo de usuario no soportado: %s", u.Tipo)
	}
	return nil
}

func (g *GestorCarne) CalcularCosto(ctx context.Context, u Usuario) (float64, error) {
	switch u.Tipo {
	case TipoPregrado:
		return 10.0, nil
	case TipoPosgrado:
		return 12.0, nil
	case TipoDocente:
		return 8.0, nil
	case TipoAdmin:
		return 7.0, nil
	case TipoExterno:
		return 20.0, nil
	default:
		return 0, fmt.Errorf("no hay tarifa para el tipo: %s", u.Tipo)
	}
}

func (g *GestorCarne) GenerarCarne(ctx context.Context, u Usuario, costo float64) (Carne, error) {
	numero := fmt.Sprintf("C-%s-%d", u.ID, time.Now().Unix())
	contenido := []byte("%PDF-1.7\n... (PDF del carné simulado) ...\n%%EOF\n")
	c := Carne{
		Numero:      numero,
		Propietario: u,
		Costo:       costo,
		EmitidoEn:   time.Now(),
		Contenido:   contenido,
	}
	fmt.Printf("[GEN] Carné generado: %s para %s %s (S/ %.2f)\n",
		c.Numero, u.Nombres, u.Apellidos, costo)
	return c, nil
}

func (g *GestorCarne) EnviarNotificacionEmail(ctx context.Context, u Usuario, c Carne) error {
	// Simulación de envío de email (sin librerías SMTP)
	subject := "Tu carné de la Biblioteca UNSCH está listo"
	body := fmt.Sprintf("Hola %s, se emitió tu carné %s. Costo: S/ %.2f. Vence: %s",
		u.Nombres, c.Numero, c.Costo, u.VenceEn.Format("2006-01-02"))
	fmt.Printf("[EMAIL] smtp=%s user=%s to=%s\nSUBJECT: %s\nBODY:\n%s\n\n",
		g.smtpHost, g.smtpUser, u.Email, subject, body)
	return nil
}

func (g *GestorCarne) GuardarEnPostgresql(ctx context.Context, u Usuario, c Carne) error {
	// Simulación de INSERTs en PostgreSQL (sin driver)
	fmt.Printf("[PG] dsn=%s\nBEGIN;\n", g.mysqlDSN)
	fmt.Printf("INSERT INTO public.usuarios (id, email, tipo, vence_en) VALUES ('%s', '%s', '%s', '%s');\n",
		u.ID, u.Email, u.Tipo, u.VenceEn.Format(time.RFC3339))
	fmt.Printf("INSERT INTO public.carnes (numero, usuario_id, costo, emitido_en) VALUES ('%s', '%s', %.2f, '%s');\n",
		c.Numero, u.ID, c.Costo, c.EmitidoEn.Format(time.RFC3339))
	fmt.Println("COMMIT;")
	return nil
}

func (g *GestorCarne) ImprimirCarne(ctx context.Context, c Carne) error {
	fmt.Printf("[PRINT] printer=%s\nImprimiendo carné %s para %s %s (S/ %.2f)\n",
		g.printerName, c.Numero, c.Propietario.Nombres, c.Propietario.Apellidos, c.Costo)
	return nil
}

// ------------------------------- main ---------------------------------------

func main() {
	ctx := context.Background()
	gestor := NuevoGestorCarne()

	usuario := Usuario{
		ID:        "U-001",
		Nombres:   "Ana",
		Apellidos: "Quispe",
		Email:     "ana.quispe@unsch.edu.pe",
		Tipo:      TipoPregrado,
		VenceEn:   time.Now().AddDate(1, 0, 0),
	}

	carne, err := gestor.ProcesarSolicitud(ctx, usuario)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	fmt.Printf("OK: Carné emitido %s (S/ %.2f) para %s %s\n", carne.Numero, carne.Costo, usuario.Nombres, usuario.Apellidos)
}
