# Módulo de bajo nivel (detalles)
class GeneradorPDF:
    @staticmethod
    # def generar_reporte(self, datos):
    def generar_reporte(_datos):
        print("Generando reporte en formato PDF...")


# Módulo de alto nivel
class GeneradorReporteUNSCH:
    def __init__(self):
        # Dependencia directa de un detalle concreto
        self.generador = GeneradorPDF()

    def generar_reporte_anual(self, datos_estudiantes):
        self.generador.generar_reporte(datos_estudiantes)
