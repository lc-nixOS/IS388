class Estudiante:
    """Responsabilidad única: Representar la entidad Estudiante"""

    def __init__(self, codigo, nombre, email, carrera):
        self.codigo = codigo
        self.nombre = nombre
        self.email = email
        self.carrera = carrera


class RegistroAcademico:
    """Responsabilidad única: Gestionar el historial académico"""

    def __init__(self):
        self.notas = []

    def agregar_nota(self, curso, nota, creditos):
        self.notas.append({"curso": curso, "nota": nota, "creditos": creditos})

    def obtener_notas(self):
        return self.notas.copy()


class CalculadoraPromedio:
    """Responsabilidad única: Realizar cálculos académicos"""

    @staticmethod
    def calcular_promedio_ponderado(notas):
        if not notas:
            return 0
        suma_ponderada = sum(nota["nota"] * nota["creditos"] for nota in notas)
        total_creditos = sum(nota["creditos"] for nota in notas)
        return suma_ponderada / total_creditos if total_creditos > 0 else 0


class GeneradorReportes:
    """Responsabilidad única: Generar reportes académicos"""

    def generar_reporte_estudiante(self, estudiante, registro_academico):
        promedio = CalculadoraPromedio.calcular_promedio_ponderado(
            registro_academico.obtener_notas()
        )
        reporte = f"=== REPORTE ACADÉMICO UNSCH ===\n"
        reporte += f"Estudiante: {estudiante.nombre} ({estudiante.codigo})\n"
        reporte += f"Carrera: {estudiante.carrera}\n"
        reporte += f"Email: {estudiante.email}\n\n"
        reporte += "HISTORIAL ACADÉMICO:\n"
        for nota_info in registro_academico.obtener_notas():
            reporte += f"- {nota_info['curso']}: {nota_info['nota']} "
            reporte += f"({nota_info['creditos']} créditos)\n"
        reporte += f"\nPromedio Ponderado: {promedio:.2f}"
        return reporte


estudiante = Estudiante("27150415", "PELAYO", "pelayo.quispe@unsch.edu.pe", "Sistemas")
registro_academico = RegistroAcademico()
registro_academico.agregar_nota("Base de Datos", 15, 4)
reporte = GeneradorReportes().generar_reporte_estudiante(estudiante, registro_academico)
print(reporte)
