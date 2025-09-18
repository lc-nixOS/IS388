# Violar el OCP: Open/Closed Principle
class CalcularPagoMatricula:
    def calcular_costo(self, tipo_estudiante):
        if tipo_estudiante == "pregrado":
            return 1500
        elif tipo_estudiante == "maestria":
            return 2500


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
    def calcular_costo(self, tipo_estudiante: TipoEstudiante):
        return tipo_estudiante.calcular_costo()
