# Violar el OCP: Open/Closed Principle
class CalcularPagoMatricula:
    @staticmethod
    def calcular_costo(tipo_estudiante):
        if tipo_estudiante == "pregrado":
            return 1500
        elif tipo_estudiante == "maestria":
            return 2500
        return None
