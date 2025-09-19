# Módulo de bajo nivel (detalles)
class GeneradorPDF:
    @staticmethod
    # def generar_reporte(self, datos):
    def generar_reporte(_datos):
        print("Generando reporte en formato PDF...")


# Módulo de alto nivel
class GeneradorReporteUNSCH:
    def __init__(self, _pdf):
        # Dependencia directa de un detalle concreto
        self.generador = GeneradorPDF()

    def generar_reporte_anual(self, datos_estudiantes):
        self.generador.generar_reporte(datos_estudiantes)


if __name__ == "__main__":
    estudiantes = ["Juan", "Ana", "Luis"]

    # Usando PDF
    reporte_pdf = GeneradorReporteUNSCH(GeneradorPDF())
    reporte_pdf.generar_reporte_anual(estudiantes)
