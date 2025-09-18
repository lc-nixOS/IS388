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
