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


if __name__ == "__main__":
    pregrado = Pregrado()
    maestria = Maestria()
    doctorado = Doctorado()

    print("Pago pregrado:", CalculadorMatricula.calcular_costo(pregrado))
    print("Pago maestría:", CalculadorMatricula.calcular_costo(maestria))
    print("Pago doctorado:", CalculadorMatricula.calcular_costo(doctorado))

    # Agregar un nuevo tipo sin modificar código existente
    class Especializacion(TipoEstudiante):
        def calcular_costo(self):
            return 1000

    especializacion = Especializacion()
    print("Pago especialización:", CalculadorMatricula.calcular_costo(especializacion))
