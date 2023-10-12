import datetime
import sys

def convert_windows_time(windows_time):
    # Define a época do Windows (1 de janeiro de 1601)
    epoch_start = datetime.datetime(year=1601, month=1, day=1)

    # Converte o tempo de arquivo para segundos
    seconds_since_epoch = windows_time / (10**7)

    # Adiciona os segundos à época do Windows
    timestamp = epoch_start + datetime.timedelta(seconds=seconds_since_epoch)

    return timestamp

# Verifique se um argumento foi passado
if len(sys.argv) < 2:
    print("Por favor, forneça o tempo do Windows como argumento.")
    sys.exit(1)

# Converta o primeiro argumento para um número inteiro e passe-o para a função
windows_time = int(sys.argv[1])

print(convert_windows_time(windows_time))
