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


if __name__ == "__main__":
    estudiante = EstudianteUsuario("2023001", "Ana", "ana@unsch.edu.pe")
    profesor = ProfesorUsuario("P001", "Dr. Pérez", "perez@unsch.edu.pe")

    # Acceso al sistema
    print(estudiante.acceder_sistema())
    print(profesor.acceder_sistema())

    # Intentar modificar notas
    try:
        print(estudiante.modificar_notas("Matemáticas", 15))
    except PermissionError as e:
        print(f"Error: {e}")

    print(profesor.modificar_notas("Matemáticas", 18))
