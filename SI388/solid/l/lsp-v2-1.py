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
