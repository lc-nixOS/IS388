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

// ------------------------ Componentes con SRP -------------------------------

// 1) Validador: SOLO valida datos de Usuario.
type UserValidator struct{}

func (UserValidator) Validate(_ context.Context, u Usuario) error {
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

// 2) Calculador de costo: SOLO calcula el costo según el tipo.
type CostCalculator struct{}

func (CostCalculator) Calculate(_ context.Context, u Usuario) (float64, error) {
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

// 3) Generador de carné: SOLO genera el carné (contenido simulado).
type CardGenerator struct{}

func (CardGenerator) Generate(_ context.Context, u Usuario, costo float64) (Carne, error) {
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

// 4) Repositorio PostgreSQL: SOLO persiste datos (simulado).
type PgRepository struct {
	pgDSN string
}

func NewPgRepository(dsn string) PgRepository { return PgRepository{pgDSN: dsn} }

func (r PgRepository) SaveUser(_ context.Context, u Usuario) error {
	fmt.Printf("[PG] dsn=%s\nBEGIN;\n", r.pgDSN)
	fmt.Printf("INSERT INTO public.usuarios (id, email, tipo, vence_en) VALUES ('%s', '%s', '%s', '%s');\n",
		u.ID, u.Email, u.Tipo, u.VenceEn.Format(time.RFC3339))
	return nil
}

func (r PgRepository) SaveCard(_ context.Context, c Carne) error {
	fmt.Printf("INSERT INTO public.carnes (numero, usuario_id, costo, emitido_en) VALUES ('%s', '%s', %.2f, '%s');\n",
		c.Numero, c.Propietario.ID, c.Costo, c.EmitidoEn.Format(time.RFC3339))
	fmt.Println("COMMIT;")
	return nil
}

// 5) Notificador por email: SOLO envía correo (simulado).
type EmailNotifier struct {
	smtpHost string
	smtpUser string
	smtpPass string
}

func (n EmailNotifier) Send(_ context.Context, to string, subject string, body string) error {
	fmt.Printf("[EMAIL] smtp=%s user=%s to=%s\nSUBJECT: %s\nBODY:\n%s\n\n",
		n.smtpHost, n.smtpUser, to, subject, body)
	return nil
}

// 6) Impresora: SOLO imprime el carné (simulado).
type Printer struct {
	printerName string
}

func (p Printer) Print(_ context.Context, c Carne) error {
	fmt.Printf("[PRINT] printer=%s\nImprimiendo carné %s para %s %s (S/ %.2f)\n",
		p.printerName, c.Numero, c.Propietario.Nombres, c.Propietario.Apellidos, c.Costo)
	return nil
}

// ------------------------- Orquestador (Caso de uso) ------------------------

// GestorCarneService: SOLO orquesta el flujo completo,
// reutilizando componentes con una única responsabilidad cada uno.
type GestorCarneService struct {
	validator UserValidator
	costCalc  CostCalculator
	generator CardGenerator
	repo      PgRepository
	notifier  EmailNotifier
	printer   Printer
}

func NewGestorCarneService(
	validator UserValidator,
	costCalc CostCalculator,
	generator CardGenerator,
	repo PgRepository,
	notifier EmailNotifier,
	printer Printer,
) GestorCarneService {
	return GestorCarneService{
		validator: validator,
		costCalc:  costCalc,
		generator: generator,
		repo:      repo,
		notifier:  notifier,
		printer:   printer,
	}
}

func (s GestorCarneService) ProcesarSolicitud(ctx context.Context, u Usuario) (Carne, error) {
	// 1) Validar datos
	if err := s.validator.Validate(ctx, u); err != nil {
		return Carne{}, fmt.Errorf("validación fallida: %w", err)
	}
	// 2) Calcular costo
	costo, err := s.costCalc.Calculate(ctx, u)
	if err != nil {
		return Carne{}, fmt.Errorf("cálculo de costo falló: %w", err)
	}
	// 3) Generar carné
	carne, err := s.generator.Generate(ctx, u, costo)
	if err != nil {
		return Carne{}, fmt.Errorf("generación de carné falló: %w", err)
	}
	// 4) Persistir (usuario y carné)
	if err := s.repo.SaveUser(ctx, u); err != nil {
		return Carne{}, fmt.Errorf("persistencia de usuario falló: %w", err)
	}
	if err := s.repo.SaveCard(ctx, carne); err != nil {
		return Carne{}, fmt.Errorf("persistencia de carné falló: %w", err)
	}
	// 5) Notificar (no crítico)
	subject := "Tu carné de la Biblioteca UNSCH está listo"
	body := fmt.Sprintf("Hola %s, se emitió tu carné %s. Costo: S/ %.2f. Vence: %s",
		u.Nombres, carne.Numero, carne.Costo, u.VenceEn.Format("2006-01-02"))
	_ = s.notifier.Send(ctx, u.Email, subject, body)
	// 6) Imprimir
	if err := s.printer.Print(ctx, carne); err != nil {
		return Carne{}, fmt.Errorf("impresión falló: %w", err)
	}
	return carne, nil
}

// --------------------------------- main -------------------------------------

func main() {
	ctx := context.Background()

	// Componentes (cada uno con UNA responsabilidad)
	validator := UserValidator{}
	costCalc := CostCalculator{}
	generator := CardGenerator{}
	repo := NewPgRepository("postgres://user:pass@localhost:5432/biblioteca?sslmode=disable")
	notifier := EmailNotifier{
		smtpHost: "smtp.unsch.edu.pe:587",
		smtpUser: "no-reply@unsch.edu.pe",
		smtpPass: "secreto-super-seguro",
	}
	printer := Printer{printerName: "Printer-Thermal-001"}

	// Orquestador SRP
	gestor := NewGestorCarneService(validator, costCalc, generator, repo, notifier, printer)

	// Ejemplo
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
	fmt.Printf("OK: Carné emitido %s (S/ %.2f) para %s %s\n",
		carne.Numero, carne.Costo, usuario.Nombres, usuario.Apellidos)
}
