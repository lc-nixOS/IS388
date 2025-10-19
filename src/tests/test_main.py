import pytest

from app.main import Calculator


class TestCalculator:
    def setup_method(self):
        self.calc = Calculator()

    def test_add(self):
        assert self.calc.add(2, 3) == 5
        assert self.calc.add(-2, 3) == 1
        assert self.calc.add(-2, -3) == -5

    def test_subtract(self):
        assert self.calc.subtract(2, 3) == -1
        assert self.calc.subtract(-2, 3) == -5
        assert self.calc.subtract(-2, -3) == 1

    def test_multiply(self):
        assert self.calc.multiply(2, 3) == 6
        assert self.calc.multiply(-2, 3) == -6
        assert self.calc.multiply(-2, -3) == 6

    def test_divide(self):
        assert self.calc.divide(2, 3) == 2 / 3
        assert self.calc.divide(-2, 3) == -2 / 3
        assert self.calc.divide(-2, -3) == 2 / 3
        with pytest.raises(ValueError):
            self.calc.divide(2, 0)
        with pytest.raises(ValueError):
            self.calc.divide(-2, 0)

    def test_power(self):
        assert self.calc.power(2, 3) == 8
        assert self.calc.power(-2, 3) == -8
        assert self.calc.power(-2, -3) == -0.125
