class Calculator:
    def sumar(self, a, b):
        return a + b

    def restar(self, a, b):
        return a - b

    def dividir(self, a, b):
        if b == 0:
            raise ValueError("No se puede dividir por cero")
        return a / b

    def multiplicar(self, a, b):
        return a * b

    def potencia(self, base, exponente):
        if exponente < 0:
            return 1 / (base ** abs(exponente))
        return base**exponente
