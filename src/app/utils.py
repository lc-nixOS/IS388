from __future__ import annotations


def add(a: int, b: int) -> int:
    """Suma dos enteros."""
    return a + b


def greet(name: str) -> str:
    # deliberately leave a small style nit (extra spaces) for Ruff/Black to catch:
    return f"Hola, {name}!"
