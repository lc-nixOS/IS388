from abc import ABC, abstractmethod


class ServicioMatriculas(ABC):
    @abstractmethod
    def matricular_curso(self, estudiante, curso):
        pass


class ServicioCalificaciones(ABC):
    @abstractmethod
    def calificar_estudiante(self, estudiante, curso, nota):
        pass


class ServicioHorarios(ABC):
    @abstractmethod
    def generar_horarios(self, semestre):
        pass


class ServicioPagos(ABC):
    @abstractmethod
    def procesar_pagos(self, estudiante, monto):
        pass


# Ahora cada servicio implementa solo lo que necesita
class GestorMatriculas(ServicioMatriculas):
    def matricular_curso(self, estudiante, curso):
        return f"Matriculando {estudiante} en {curso} - UNSCH"


class GestorCalificaciones(ServicioCalificaciones):
    def calificar_estudiante(self, estudiante, curso, nota):
        return f"Calificando {estudiante} en {curso}: {nota}"


class GestorHorarios(ServicioHorarios):
    def generar_horarios(self, semestre):
        return f"Generando horarios para el semestre {semestre}"


class GestorPagos(ServicioPagos):
    def procesar_pagos(self, estudiante, monto):
        return f"Procesando pago de {monto} para {estudiante}"


if __name__ == "__main__":
    m = GestorMatriculas()
    c = GestorCalificaciones()
    h = GestorHorarios()
    p = GestorPagos()

    print(m.matricular_curso("Juan", "Matemáticas"))
    print(c.calificar_estudiante("Ana", "Física", 18))
    print(h.generar_horarios("2025-I"))
    print(p.procesar_pagos("Luis", 350.50))
