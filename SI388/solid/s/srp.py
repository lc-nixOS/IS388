class Estudiante:
   def __init__(self, codigo, nombre, email, carrera):
       self.codigo = codigo
       self.nombre = nombre
       self.email = email
       self.carrera = carrera
       self.notas = []
   def agregar_nota(self, curso, nota, creditos):
       self.notas.append({
           'curso': curso,
           'nota': nota,
           'creditos': creditos
       })
   def calcular_promedio_ponderado(self):
       if not self.notas:
           return 0
       suma_ponderada = sum(nota['nota'] * nota['creditos']
                            for nota in self.notas)
       total_creditos = sum(nota['creditos'] for nota in self.notas)

       return suma_ponderada / total_creditos if total_creditos > 0 else 0

   def generar_reporte_academico(self):
       reporte = f"=== REPORTE ACADÉMICO UNSCH ===\n"
       reporte += f"Estudiante: {self.nombre} ({self.codigo})\n"
       reporte += f"Carrera: {self.carrera}\n"
       reporte += f"Email: {self.email}\n\n"
       reporte += "HISTORIAL ACADÉMICO:\n"
       for nota_info in self.notas:
           reporte += f"- {nota_info['curso']}: {nota_info['nota']} "
           reporte += f"({nota_info['creditos']} créditos)\n"
       reporte += f"\nPromedio Ponderado: {self.calcular_promedio_ponderado():.2f}"
       return reporte

   def enviar_email_notificacion(self, mensaje):
       print(f"Enviando email a {self.email}: {mensaje}")
       return True
   def guardar_en_base_datos(self):
       # Simulación de guardado en BD
       print(f"Guardando estudiante {self.codigo} en base de datos...")
       return True

Estudinate1=Estudiante("27150415","Pelayo Quispe Bautista","pelayo.quispeunsch.edu.pe","Ing. Sistemas")

Estudinate1.agregar_nota("Base de Datos",20,4)

Estudinate1.agregar_nota("algoritmos 1",15,3)

Estudinate1.guardar_en_base_datos()
Estudinate1.enviar_email_notificacion("Hola mundo")
print(Estudinate1.generar_reporte_academico())