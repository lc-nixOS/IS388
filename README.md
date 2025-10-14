# GitHub Actions CI

## 1 `requirements.txt` (con pre-commit)

Crea/actualiza este archivo en la raíz:

```txt
ruff==0.5.5
black==24.3.0
pytest==8.3.2
pre-commit==3.7.1
```

> Si ya usas `requirements-dev.txt`, puedes mantenerlo y duplicar ahí también; si sólo quieres uno, con `requirements.txt` basta.

---

## 2 `.pre-commit-config.yaml`

Ponlo en la raíz del repo:

```yaml
repos:
  - repo: https://github.com/astral-sh/ruff-pre-commit
    rev: v0.5.5
    hooks:
      - id: ruff
        args: [--fix] # auto-arregla lo posible
      - id: ruff-format # formateo estilo black (si lo usas, black puede sobrar)
  - repo: https://github.com/psf/black
    rev: 24.3.0
    hooks:
      - id: black
```

> Si prefieres que el formateo lo haga sólo Black, puedes quitar `ruff-format`.

---

## 3 Workflow de GitHub Actions

Archivo: `.github/workflows/ci.yml`

```yaml
name: CI (Lint, Types, Tests)

on:
  pull_request:
  push:
    branches: [main, master]

jobs:
  ci:
    runs-on: ubuntu-latest
    strategy:
      fail-fast: false
      matrix:
        python-version: ["3.11"]

    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Set up Python
        uses: actions/setup-python@v5
        with:
          python-version: ${{ matrix.python-version }}
          cache: "pip"

      - name: Install dev deps
        run: |
          python -m pip install --upgrade pip
          pip install -r requirements.txt

      - name: Black (check only)
        run: black --check src tests

      - name: Ruff (lint)
        run: ruff check --output-format=github src tests

      - name: Pytest
        run: pytest -q
```

---

## 4 Comandos para probar **en local**

### A Preparar entorno

```bash
# Crear y activar venv (Unix/macOS)
python -m venv .venv
source .venv/bin/activate

# (Windows PowerShell)
# python -m venv .venv
# .\.venv\Scripts\Activate.ps1

# Instalar dependencias
pip install --upgrade pip
pip install -r requirements.txt
```

### B Instalar y correr pre-commit

```bash
# Instala los hooks en .git/hooks
pre-commit install

# Probar todos los hooks sobre todos los archivos (una sola vez al inicio)
pre-commit run --all-files
```

### C Ejecutar herramientas manualmente (útil antes de commitear)

```bash
# Opción 1: rutas explícitas
black src src/tests
ruff format src src/tests
ruff check src src/tests
pytest -q src/tests

# Opción 2: todo el repo (si no tienes basura que formatear)
black .
ruff format .
ruff check .
pytest -q
```

> Repite estos pasos hasta que no haya errores. Luego haz commit/push.

**Qué debería pasar:**

- `black` no debe quejarse (si lo hace, corrige con `black ...`).
- `ruff format` formatea; `ruff check` marca problemas de lint (puedes auto-fijar muchos con `ruff check --fix src src/tests`).
- `pytest` corre tests (ya sin `ModuleNotFoundError` gracias a `pytest.ini` o `PYTHONPATH`).

---

### 5 ¿Qué debería suceder al subir?

1. **Push a una rama** o abre un **Pull Request** → se dispara el workflow `CI (Lint, Types, Tests)`.
2. En la pestaña **Actions** verás los pasos ejecutándose:

   - **Black (check only)**: ✔️ si todo está formateado; ❌ si falta formatear (solucionalo con `black src tests`).
   - **Ruff (lint)**: crea **anotaciones** en el PR mostrando los problemas exactos en cada línea.
   - **Pytest**: corre tests (si fallan, verás el log con el assert roto).

3. En el **PR**, los checks deben aparecer en **verde** para poder hacer merge (según tu configuración de branch protection).
4. Si algo falla:
   - Corre localmente los comandos del punto 4C.
   - Corrige y vuelve a hacer **commit + push**; el CI se re-ejecuta y debería pasar.

---

## 6 Mini ejemplo de carpetas

```txt
tu-repo/
├─ src/
│  └─ app/
│     ├─ __init__.py
│     ├─ main.py
│     └─ utils.py
├─ tests/
│  └─ test_utils.py
├─ requirements.txt
├─ .pre-commit-config.yaml
└─ .github/
   └─ workflows/
      └─ ci.yml
```

## 7 Comprobación rápida

1. Ejecuta local:

   ```powershell
   black src src/tests
   ruff check --fix src src/tests
   $env:PYTHONPATH = "$PWD\src"
   pytest -q src/tests
   ```

2. Si todo está OK, `git add . && git commit -m "fix(ci): rutas y config ruff/pytest" && git push`.
3. En el PR, todo verde ✅.

## 8 Actualiza pre-commit

```txt
git add .pre-commit-config.yaml
pre-commit clean
pre-commit install
pre-commit autoupdate
```
