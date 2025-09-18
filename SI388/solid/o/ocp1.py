# Aplicar el OCP: Usaremos herencia y polimorfismo. Creamos una clase base y luego subclases para cada tipo de estudiante.
class TipoEstudiante:
    def calcular_costo(self):
        pass


class Pregrado(TipoEstudiante):
    def calcular_costo(self):
        return 1500


class Maestria(TipoEstudiante):
    def calcular_costo(self):
        return 2500


class Doctorado(TipoEstudiante):
    def calcular_costo(self):
        return 3500


class CalculadorMatricula:
    @staticmethod
    def calcular_costo(tipo_estudiante: TipoEstudiante):
        return tipo_estudiante.calcular_costo()


ca = CalculadorMatricula()
print(ca.calcular_costo(Doctorado()))
