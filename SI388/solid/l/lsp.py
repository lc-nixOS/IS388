class Usuario:
    def __init__(self, codigo, nombre, email):
        self.codigo = codigo
        self.nombre = nombre
        self.email = email

    def acceder_sistema(self):
        return f"{self.nombre} ha accedido al sistema UNSCH"

    def modificar_notas(self, curso, nueva_nota):
        # Todos los usuarios pueden modificar notas (problemático)
        return f"Nota modificada en {curso}: {nueva_nota}"


class EstudianteUsuario(Usuario):
    def modificar_notas(self, curso, nueva_nota):
        raise PermissionError("Los estudiantes no pueden modificar notas")


class ProfesorUsuario(Usuario):
    def modificar_notas(self, curso, nueva_nota):
        return f"Profesor modificó nota en {curso}: {nueva_nota}"
