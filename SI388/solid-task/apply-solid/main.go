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
	VenceEn   time.Time
}

type Carne struct {
	Numero      string
	Propietario Usuario
	Costo       float64
	EmitidoEn   time.Time
	Contenido   []byte
}

// ------------------------ Componentes con SRP -------------------------------

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

// PriceStrategy define una estrategia de precio para uno o más tipos de usuario
type PriceStrategy interface {
	Supports(t TipoUsuario) bool
	Calculate(ctx context.Context, u Usuario) (float64, error)
}

// FlatPriceStrategy: precio fijo para un tipo concreto
type FlatPriceStrategy struct {
	Supported TipoUsuario
	Price     float64
}

func (s FlatPriceStrategy) Supports(t TipoUsuario) bool { return s.Supported == t }
func (s FlatPriceStrategy) Calculate(_ context.Context, _ Usuario) (float64, error) {
	return s.Price, nil
}

// CompositePriceCalculator: selecciona la primera estrategia que soporte el tipo
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

// (Opcional) Ejemplo de extensión sin tocar el servicio ni otras piezas:
// HappyHourStudent agrega un descuento temporal a PREGRADO entre horas dadas.
type HappyHourStudent struct {
	HourStart int // 0..23
	HourEnd   int // 0..23 (no inclusivo)
	Base      float64
	Discount  float64 // monto fijo a descontar (no negativo)
}

func (s HappyHourStudent) Supports(t TipoUsuario) bool { return t == TipoPregrado }
func (s HappyHourStudent) Calculate(_ context.Context, _ Usuario) (float64, error) {
	now := time.Now()
	h := now.Hour()
	price := s.Base
	if h >= s.HourStart && h < s.HourEnd {
		if s.Discount > price {
			price = 0
		} else {
			price -= s.Discount
		}
	}
	return price, nil
}

// 3) Generador de carné (SRP): SOLO genera el carné
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

// 4) Repositorio PostgreSQL (SRP): SOLO persiste (simulado)
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

// 5) Notificador (SRP): SOLO envía email (simulado)
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

// 6) Impresora (SRP): SOLO imprime (simulado)
type Printer struct {
	printerName string
}

func (p Printer) Print(_ context.Context, c Carne) error {
	fmt.Printf("[PRINT] printer=%s\nImprimiendo carné %s para %s %s (S/ %.2f)\n",
		p.printerName, c.Numero, c.Propietario.Nombres, c.Propietario.Apellidos, c.Costo)
	return nil
}

// ------------------------- Orquestador (Caso de uso) ------------------------

type priceCalculator interface {
	Calculate(ctx context.Context, u Usuario) (float64, error)
}

type GestorCarneService struct {
	validator UserValidator
	pricing   priceCalculator
	generator CardGenerator
	repo      PgRepository
	notifier  EmailNotifier
	printer   Printer
}

func NewGestorCarneService(
	validator UserValidator,
	pricing priceCalculator,
	generator CardGenerator,
	repo PgRepository,
	notifier EmailNotifier,
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
	// 2) Calcular costo (OCP vía estrategias)
	costo, err := s.pricing.Calculate(ctx, u)
	if err != nil {
		return Carne{}, fmt.Errorf("cálculo de costo falló: %w", err)
	}
	// 3) Generar carné
	carne, err := s.generator.Generate(ctx, u, costo)
	if err != nil {
		return Carne{}, fmt.Errorf("generación de carné falló: %w", err)
	}
	// 4) Persistir
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

	// Estrategias base (fácil de extender sin modificar nada del servicio)
	pricing := NewCompositePriceCalculator(
		FlatPriceStrategy{Supported: TipoPregrado, Price: 10.0},
		FlatPriceStrategy{Supported: TipoPosgrado, Price: 12.0},
		FlatPriceStrategy{Supported: TipoDocente, Price: 8.0},
		FlatPriceStrategy{Supported: TipoAdmin, Price: 7.0},
		FlatPriceStrategy{Supported: TipoExterno, Price: 20.0},

		HappyHourStudent{HourStart: 9, HourEnd: 11, Base: 10.0, Discount: 2.0},
	)

	validator := UserValidator{}
	generator := CardGenerator{}
	repo := NewPgRepository("postgres://user:pass@localhost:5432/biblioteca?sslmode=disable")
	notifier := EmailNotifier{
		smtpHost: "smtp.unsch.edu.pe:587",
		smtpUser: "no-reply@unsch.edu.pe",
		smtpPass: "secreto-super-seguro",
	}
	printer := Printer{printerName: "Printer-Thermal-001"}

	gestor := NewGestorCarneService(validator, pricing, generator, repo, notifier, printer)

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
