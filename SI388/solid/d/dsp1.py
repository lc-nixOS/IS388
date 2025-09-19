from abc import ABC, abstractmethod


class GeneradorReporte(ABC):
    @abstractmethod
    def generar_reporte(self, datos):
        pass


# Módulos de bajo nivel (detalles) que dependen de la abstracción
class GeneradorPDF(GeneradorReporte):
    def generar_reporte(self, datos):
        print("Generando reporte en formato PDF...")


class GeneradorCSV(GeneradorReporte):
    def generar_reporte(self, datos):
        print("Generando reporte en formato CSV...")


# Módulo de alto nivel que depende de la abstracción
class GeneradorReporteUNSCH:
    def __init__(self, generador: GeneradorReporte):
        # Ahora el módulo de alto nivel depende de la abstracción.
        self.generador = generador

    def generar_reporte_anual(self, datos_estudiantes):
        self.generador.generar_reporte(datos_estudiantes)


# Uso del código
datos_ejemplo = {"estudiante": "Maria", "curso": "Fisica"}

# Inyectamos la dependencia en el constructor
generador_pdf = GeneradorReporteUNSCH(GeneradorPDF())
generador_csv = GeneradorReporteUNSCH(GeneradorCSV())
generador_pdf.generar_reporte_anual(datos_ejemplo)
generador_csv.generar_reporte_anual(datos_ejemplo)
