import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent.parent))
from calculator import Calculator
import pytest


class TestCalculator:
    def setup_method(self):
        self.calc = Calculator()

    def test_sumar(self):
        assert self.calc.sumar(2, 3) == 5
        assert self.calc.sumar(-1, 1) == 0

    def test_restar(self):
        assert self.calc.restar(5, 3) == 2
        assert self.calc.restar(0, 5) == -5

    def test_dividir(self):
        assert self.calc.dividir(10, 2) == 5
        assert self.calc.dividir(9, 3) == 3

    def test_dividir_por_cero(self):
        with pytest.raises(ValueError, match="No se puede dividir por cero"):
            self.calc.dividir(10, 0)

    def test_multiplicar(self):
        assert self.calc.multiplicar(2, 3) == 6
        assert self.calc.multiplicar(-1, 1) == -1

    def test_potencia_positiva(self):
        assert self.calc.potencia(2, 3) == 8
        assert self.calc.potencia(5, 2) == 25

    def test_potencia_negativa(self):
        assert pytest.approx(self.calc.potencia(2, -3)) == 0.125


# pytest --cov=calculator --cov-report=term-missing
