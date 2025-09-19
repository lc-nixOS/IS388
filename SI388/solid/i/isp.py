from abc import ABC, abstractmethod


class ServicioAcademicoCompleto(ABC):
    @abstractmethod
    def matricular_curso(self, estudiante, curso):
        pass

    @abstractmethod
    def calificar_estudiante(self, estudiante, curso, nota):
        pass

    @abstractmethod
    def generar_horarios(self, semestre):
        pass

    @abstractmethod
    def gestionar_aulas(self, aula, operacion):
        pass

    @abstractmethod
    def procesar_pagos(self, estudiante, monto):
        pass


# los clientes se ven obligados a implementar métodos que no necesitan
class ServicoMatriculas(ServicioAcademicoCompleto):
    def gestionar_aulas(self, aula, operacion):
        pass

    def matricular_curso(self, estudiante, curso):
        print(f"Matriculando {estudiante} en el curso {curso}")

    # metodos que no deberia implementar pero esta obligado a hacerlo
    def calificar_estudiante(self, estudiante, curso, nota):
        raise NotImplementedError("No es responsabilidad de matriculas")

    def generar_horarios(self, semestre):
        raise NotImplementedError("No es responsabilidad de matriculas")

    def procesar_pagos(self, estudiante, monto):
        raise NotImplementedError("No es responsabilidad de matriculas")


if __name__ == "__main__":
    servicio = ServicoMatriculas()
    servicio.matricular_curso("Juan", "Matemáticas")

    # Si intentamos usar algo que no le corresponde:
    try:
        servicio.procesar_pagos("Juan", 200)
    except NotImplementedError as e:
        print(f"Error: {e}")
