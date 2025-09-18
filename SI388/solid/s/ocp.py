class CalcularPagoMatricula:
  def calcular_costo(self, tipo_estudiante):
    if tipo_estudiante == "pregrado":
      return 1500
    elif tipo_estudiante == "maestria":
      return 2500