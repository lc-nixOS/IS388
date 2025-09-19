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

    @staticmethod
    def generar_reporte_estudiante(_self, estudiante, registro_academico):
        promedio = CalculadoraPromedio.calcular_promedio_ponderado(
            registro_academico.obtener_notas()
        )

        reporte = f"=== REPORTE ACADÉMICO UNSCH ===\n"
        reporte += f"Estudiante: {estudiante.nombre} ({estudiante.codigo})\n"
        reporte += f"Carrera: {estudiante.carrera}\n"
        reporte += f"Email: {estudiante.email}\n\n"
        reporte += "HISTORIAL ACADÉMICO:\n"
        for nota_info in registro_academico.obtener_notas():
            reporte += f"{nota_info['curso']}: {nota_info['nota']} "
            reporte += f"({nota_info['creditos']} créditos)\n"
        reporte += f"\nPromedio Ponderado: {promedio:.2f}"
        return reporte


class EnviarEmail:
    def __init__(self, email):
        self.email = email

    def enviar_email_notificacion(self, mensaje):
        # Simulación de envio de email
        print(f"Enviando email a {self.email}: {mensaje}")
        # Aquí iría la lógica real de envío
        return True


if __name__ == "__main__":
    Estudiante1 = Estudiante(
        "27202506",
        "Isaias Ramos Lopez",
        "isaias.ramos.27@unsch.edu.pe",
        "ing de sistemas",
    )

    registro = RegistroAcademico()
    registro.agregar_nota("Base de Datos", 20, 4)
    registro.agregar_nota("Algoritmos", 15, 3)

    print(GeneradorReportes.generar_reporte_estudiante(None, Estudiante1, registro))
    print(
        f"Promedio Ponderado: {CalculadoraPromedio.calcular_promedio_ponderado(registro.obtener_notas()):.2f}"
    )

    notificador = EnviarEmail(Estudiante1.email)
    notificador.enviar_email_notificacion("Hola mundo")
