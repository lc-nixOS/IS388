package main

import (
	"context"
	"encoding/json"
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
	VenceEn   time.Time
}

type Carne struct {
	Numero      string
	Propietario Usuario
	Costo       float64
	EmitidoEn   time.Time
	Contenido   []byte
}

// ------------------------- Abstracciones (LSP/DIP) --------------------------

type Validator interface {
	Validate(ctx context.Context, u Usuario) error
}

type PriceCalculator interface {
	Calculate(ctx context.Context, u Usuario) (float64, error)
}

type CardGenerator interface {
	Generate(ctx context.Context, u Usuario, costo float64) (Carne, error)
}

type Repository interface {
	SaveUser(ctx context.Context, u Usuario) error
	SaveCard(ctx context.Context, c Carne) error
}

type Notifier interface {
	Send(ctx context.Context, to string, subject string, body string) error
}

type Printer interface {
	Print(ctx context.Context, c Carne) error
}

// ------------------------ Implementaciones (SRP) ----------------------------

// 1) Validador
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
	default:
		return fmt.Errorf("tipo de usuario no soportado: %s", u.Tipo)
	}
	return nil
}

// 2) OCP: estrategias de precio
type PriceStrategy interface {
	Supports(t TipoUsuario) bool
	Calculate(ctx context.Context, u Usuario) (float64, error)
}

type FlatPriceStrategy struct {
	Supported TipoUsuario
	Price     float64
}

func (s FlatPriceStrategy) Supports(t TipoUsuario) bool { return s.Supported == t }
func (s FlatPriceStrategy) Calculate(_ context.Context, _ Usuario) (float64, error) {
	return s.Price, nil
}

type HappyHourStudent struct {
	HourStart int
	HourEnd   int // no inclusivo
	Base      float64
	Discount  float64
}

func (s HappyHourStudent) Supports(t TipoUsuario) bool { return t == TipoPregrado }
func (s HappyHourStudent) Calculate(_ context.Context, _ Usuario) (float64, error) {
	now := time.Now().Hour()
	price := s.Base
	if now >= s.HourStart && now < s.HourEnd {
		if s.Discount > price {
			price = 0
		} else {
			price -= s.Discount
		}
	}
	return price, nil
}

type CompositePriceCalculator struct {
	strategies []PriceStrategy
}

func NewCompositePriceCalculator(strategies ...PriceStrategy) *CompositePriceCalculator {
	return &CompositePriceCalculator{strategies: strategies}
}

func (c *CompositePriceCalculator) Calculate(ctx context.Context, u Usuario) (float64, error) {
	for _, s := range c.strategies {
		if s.Supports(u.Tipo) {
			return s.Calculate(ctx, u)
		}
	}
	return 0, fmt.Errorf("no hay estrategia de precio para el tipo: %s", u.Tipo)
}

// 3) Generadores de carné (LSP: PDF y JSON son sustituibles)
type PDFCardGenerator struct{}

func (PDFCardGenerator) Generate(_ context.Context, u Usuario, costo float64) (Carne, error) {
	card := Carne{
		Numero:      fmt.Sprintf("C-%s-%d", u.ID, time.Now().Unix()),
		Propietario: u,
		Costo:       costo,
		EmitidoEn:   time.Now(),
		Contenido:   []byte("%PDF-1.7\n... (PDF del carné simulado) ...\n%%EOF\n"),
	}
	fmt.Printf("[GEN/PDF] Carné %s para %s %s (S/ %.2f)\n", card.Numero, u.Nombres, u.Apellidos, costo)
	return card, nil
}

type JSONCardGenerator struct{}

func (JSONCardGenerator) Generate(_ context.Context, u Usuario, costo float64) (Carne, error) {
	card := Carne{
		Numero:      fmt.Sprintf("C-%s-%d", u.ID, time.Now().Unix()),
		Propietario: u,
		Costo:       costo,
		EmitidoEn:   time.Now(),
	}
	payload := map[string]any{
		"numero":      card.Numero,
		"propietario": u,
		"costo":       card.Costo,
		"emitido_en":  card.EmitidoEn.Format(time.RFC3339),
	}
	b, _ := json.Marshal(payload) // contrato: siempre retorna contenido válido
	card.Contenido = b
	fmt.Printf("[GEN/JSON] Carné %s para %s %s (S/ %.2f)\n", card.Numero, u.Nombres, u.Apellidos, costo)
	return card, nil
}

// 4) Repositorios (LSP: Postgres y Memory son sustituibles)
type PgRepository struct{ pgDSN string }

func NewPgRepository(dsn string) PgRepository { return PgRepository{pgDSN: dsn} }

func (r PgRepository) SaveUser(_ context.Context, u Usuario) error {
	fmt.Printf("[PG] dsn=%s\nBEGIN;\n", r.pgDSN)
	fmt.Printf("INSERT INTO public.usuarios (id, email, tipo, vence_en) VALUES ('%s','%s','%s','%s');\n",
		u.ID, u.Email, u.Tipo, u.VenceEn.Format(time.RFC3339))
	return nil
}
func (r PgRepository) SaveCard(_ context.Context, c Carne) error {
	fmt.Printf("INSERT INTO public.carnes (numero, usuario_id, costo, emitido_en) VALUES ('%s','%s',%.2f,'%s');\nCOMMIT;\n",
		c.Numero, c.Propietario.ID, c.Costo, c.EmitidoEn.Format(time.RFC3339))
	return nil
}

type MemoryRepository struct {
	users []Usuario
	cards []Carne
}

func (m *MemoryRepository) SaveUser(_ context.Context, u Usuario) error {
	m.users = append(m.users, u)
	fmt.Printf("[MEM] guardado usuario id=%s email=%s\n", u.ID, u.Email)
	return nil
}
func (m *MemoryRepository) SaveCard(_ context.Context, c Carne) error {
	m.cards = append(m.cards, c)
	fmt.Printf("[MEM] guardado carné numero=%s costo=%.2f\n", c.Numero, c.Costo)
	return nil
}

// 5) Notificadores (LSP: Email y Null son sustituibles)
type EmailNotifier struct {
	smtpHost string
	smtpUser string
	smtpPass string
}

func (n EmailNotifier) Send(_ context.Context, to, subject, body string) error {
	fmt.Printf("[EMAIL] smtp=%s user=%s to=%s\nSUBJECT: %s\nBODY:\n%s\n\n",
		n.smtpHost, n.smtpUser, to, subject, body)
	return nil
}

type NullNotifier struct{}

func (NullNotifier) Send(_ context.Context, _ string, _ string, _ string) error {
	// Implementación nula: cumple el contrato sin efectos secundarios.
	return nil
}

// 6) Impresoras (LSP: Térmica y DryRun son sustituibles)
type ThermalPrinter struct{ name string }

func (p ThermalPrinter) Print(_ context.Context, c Carne) error {
	fmt.Printf("[PRINT] printer=%s\nImprimiendo carné %s para %s %s (S/ %.2f)\n",
		p.name, c.Numero, c.Propietario.Nombres, c.Propietario.Apellidos, c.Costo)
	return nil
}

type DryRunPrinter struct{}

func (DryRunPrinter) Print(_ context.Context, c Carne) error {
	fmt.Printf("[PRINT/DRYRUN] (sin impresión) carné %s listo\n", c.Numero)
	return nil
}

// ------------------------- Orquestador (Caso de uso) ------------------------

type GestorCarneService struct {
	validator Validator
	pricing   PriceCalculator
	generator CardGenerator
	repo      Repository
	notifier  Notifier
	printer   Printer
}

func NewGestorCarneService(
	validator Validator,
	pricing PriceCalculator,
	generator CardGenerator,
	repo Repository,
	notifier Notifier,
	printer Printer,
) GestorCarneService {
	return GestorCarneService{
		validator: validator,
		pricing:   pricing,
		generator: generator,
		repo:      repo,
		notifier:  notifier,
		printer:   printer,
	}
}

func (s GestorCarneService) ProcesarSolicitud(ctx context.Context, u Usuario) (Carne, error) {
	// 1) Validar
	if err := s.validator.Validate(ctx, u); err != nil {
		return Carne{}, fmt.Errorf("validación fallida: %w", err)
	}
	// 2) Calcular costo
	costo, err := s.pricing.Calculate(ctx, u)
	if err != nil {
		return Carne{}, fmt.Errorf("cálculo de costo falló: %w", err)
	}
	// 3) Generar carné
	card, err := s.generator.Generate(ctx, u, costo)
	if err != nil {
		return Carne{}, fmt.Errorf("generación de carné falló: %w", err)
	}
	// 4) Persistir
	if err := s.repo.SaveUser(ctx, u); err != nil {
		return Carne{}, fmt.Errorf("persistencia de usuario falló: %w", err)
	}
	if err := s.repo.SaveCard(ctx, card); err != nil {
		return Carne{}, fmt.Errorf("persistencia de carné falló: %w", err)
	}
	// 5) Notificar (no crítico)
	subject := "Tu carné de la Biblioteca UNSCH está listo"
	body := fmt.Sprintf("Hola %s, se emitió tu carné %s. Costo: S/ %.2f. Vence: %s",
		u.Nombres, card.Numero, card.Costo, u.VenceEn.Format("2006-01-02"))
	_ = s.notifier.Send(ctx, u.Email, subject, body)
	// 6) Imprimir
	if err := s.printer.Print(ctx, card); err != nil {
		return Carne{}, fmt.Errorf("impresión falló: %w", err)
	}
	return card, nil
}

// --------------------------------- main -------------------------------------

func main() {
	ctx := context.Background()

	// OCP: fácilmente extensible con nuevas estrategias
	pricing := NewCompositePriceCalculator(
		FlatPriceStrategy{Supported: TipoPregrado, Price: 10.0},
		FlatPriceStrategy{Supported: TipoPosgrado, Price: 12.0},
		FlatPriceStrategy{Supported: TipoDocente, Price: 8.0},
		FlatPriceStrategy{Supported: TipoAdmin, Price: 7.0},
		FlatPriceStrategy{Supported: TipoExterno, Price: 20.0},
		// Extensión opcional
		HappyHourStudent{HourStart: 9, HourEnd: 11, Base: 10.0, Discount: 2.0},
	)

	// LSP en acción: prueba con diferentes implementaciones sin cambiar el servicio
	validator := UserValidator{}
	// Puedes alternar entre PDFCardGenerator{} y JSONCardGenerator{} sin romper nada.
	generator := PDFCardGenerator{} // o JSONCardGenerator{}
	// Puedes alternar entre repositorio real simulado o en memoria.
	repo := NewPgRepository("postgres://user:pass@localhost:5432/biblioteca?sslmode=disable")
	// repoMem := &MemoryRepository{} // <- sustituto válido
	// Puedes alternar entre notificador real o nulo.
	notifier := EmailNotifier{smtpHost: "smtp.unsch.edu.pe:587", smtpUser: "no-reply@unsch.edu.pe", smtpPass: "secreto-super-seguro"}
	// notifier := NullNotifier{} // <- sustituto válido
	// Puedes alternar entre impresora real o dry-run.
	printer := ThermalPrinter{name: "Printer-Thermal-001"}
	// printer := DryRunPrinter{} // <- sustituto válido

	gestor := NewGestorCarneService(validator, pricing, generator, repo, notifier, printer)

	usuario := Usuario{
		ID:        "U-001",
		Nombres:   "Ana",
		Apellidos: "Quispe",
		Email:     "ana.quispe@unsch.edu.pe",
		Tipo:      TipoPregrado,
		VenceEn:   time.Now().AddDate(1, 0, 0),
	}

	card, err := gestor.ProcesarSolicitud(ctx, usuario)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	fmt.Printf("OK: Carné emitido %s (S/ %.2f) para %s %s\n",
		card.Numero, card.Costo, usuario.Nombres, usuario.Apellidos)
}
