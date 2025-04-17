import ctypes

# Initialiser la graine
libc = ctypes.CDLL("libc.so.6")
libc.srand(0xCFD2BC5B)

# Générer les nombres aléatoires
random_values = [libc.rand() for _ in range(0x9e)]
code_machine = bytes([r & 0xFF for r in random_values])

# Sauvegarder le code généré
with open("generated_code2.bin", "wb") as f:
    f.write(code_machine)
