class Usuario:
    def __init__(self, codigo, nombre, email):
        self.codigo = codigo
        self.nombre = nombre
        self.email = email

    def acceder_sistema(self):
        return f"{self.nombre} ha accedido al sistema UNSCH"


class UsuarioConPrivilegios(Usuario):
    @staticmethod
    def modificar_notas(curso, nueva_nota):
        return f"Nota modificada en {curso}: {nueva_nota}"


class EstudianteUsuario(Usuario):
    @staticmethod
    def consultar_notas():
        return "Consultando notas del estudiante"


class ProfesorUsuario(UsuarioConPrivilegios):
    @staticmethod
    def asignar_nota(_, estudiante, curso, nota):
        return f"Nota asignada a {estudiante} en {curso}: {nota}"


class AdministradorUsuario(UsuarioConPrivilegios):
    @staticmethod
    def generar_reporte_general(_):
        return "Generando reporte administrativo"


if __name__ == "__main__":
    estudiante = EstudianteUsuario("2023001", "Ana", "ana@unsch.edu.pe")
    profesor = ProfesorUsuario("P001", "Dr. Pérez", "perez@unsch.edu.pe")
    admin = AdministradorUsuario("A001", "Lic. Torres", "torres@unsch.edu.pe")

    # Acceder al sistema
    print(estudiante.acceder_sistema())
    print(profesor.acceder_sistema())
    print(admin.acceder_sistema())

    # Estudiante: solo consultar
    print(estudiante.consultar_notas())

    # Profesor: asignar y modificar
    print(profesor.asignar_nota("Dr. Pérez", "Ana", "Matemáticas", 18))
    print(profesor.modificar_notas("Física", 17))

    # Administrador: generar reportes y modificar
    print(admin.generar_reporte_general("Reporte 2025"))
    print(admin.modificar_notas("Química", 20))
