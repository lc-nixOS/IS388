from __future__ import annotations

from src.app.utils import add, greet


def run() -> None:
    total = add(2, 3)
    mensaje = greet("Mundo")
    print(mensaje, "-", "Total:", total)


if __name__ == "__main__":
    run()
