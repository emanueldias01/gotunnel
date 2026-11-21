import socket

HOST = "127.0.0.1"
PORT = 3000

def main():
    try:
        with socket.create_connection((HOST, PORT), timeout=5) as client:
            print(f"Conectado ao servidor TCP em {HOST}:{PORT}")
            print("Digite mensagens para enviar. Use 'sair' para encerrar.")

            while True:
                message = input("> ")

                if message.lower() == "sair":
                    break

                client.sendall(message.encode("utf-8"))
                print(f"Enviado: {message}")

    except ConnectionRefusedError:
        print(f"Não foi possível conectar a {HOST}:{PORT}. Verifique se o servidor está rodando.")
    except OSError as error:
        print(f"Erro de conexão: {error}")

if __name__ == "__main__":
    main()
