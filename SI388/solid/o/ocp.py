# Violar el OCP: Open/Closed Principle
class CalcularPagoMatricula:
    @staticmethod
    def calcular_costo(tipo_estudiante):
        if tipo_estudiante == "pregrado":
            return 1500
        elif tipo_estudiante == "maestria":
            return 2500
        elif tipo_estudiante == "doctorado":
            return 3500
        return None


if __name__ == "__main__":
    print("Pago pregrado:", CalcularPagoMatricula.calcular_costo("pregrado"))
    print("Pago maestría:", CalcularPagoMatricula.calcular_costo("maestria"))
    print("Pago doctorado:", CalcularPagoMatricula.calcular_costo("doctorado"))
