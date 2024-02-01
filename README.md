# Discord Token Brute Forcer

This project is a Go-based application designed to brute force Discord tokens. It uses a multi-threaded approach to maximize efficiency and speed.

## Description

The application reads a list of proxies from a file and uses them to distribute the brute force attempts across different IP addresses. This is done to avoid rate limiting and to increase the chances of finding a valid token.

The user is prompted to enter a Discord user ID and the number of threads per proxy. The user ID is then base64 encoded and used in the brute force attempts.

The application uses Go's native concurrency model, goroutines, to perform multiple brute force attempts simultaneously. Each goroutine represents a single brute force attempt.

## How it Works

1. The application starts by clearing the console for a clean start.
2. It then reads a list of proxies from a file named "proxies.txt".
3. The application displays a startup screen with the number of proxies and the proxies themselves.
4. The user is prompted to enter a Discord user ID, which is then base64 encoded.
5. The user is also prompted to enter the number of threads per proxy.
6. The application then starts the brute force attempts. For each proxy, it creates a number of goroutines equal to the number of threads specified by the user.
7. Each goroutine calls the `generateRandomToken` function with the base64 encoded user ID, a WaitGroup (for synchronization), and the proxy URL.
8. The application waits for all goroutines to finish before exiting.

## Disclaimer

This tool is intended for educational purposes only. Unauthorized attempts to brute force Discord tokens is against Discord's Terms of Service and can lead to your account being banned. Use this tool responsibly and ethically.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
