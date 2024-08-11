# goexpert-challenge-1-api

Esse repositório corresponde à solução do desafio 1 do curso de Go Expert, de 2024.
Arquivos inclusos são o `go.mod` e `go.sum`. Para instalar os pacotes, basta clonar
o repositório e executar `go mod tidy`. Depois, expor o servidor com `go run server.go`.
Para testar rapidamente que o servidor está online, basta usar `curl localhost:8080/cotacao`.
Após verificação, executar `go run client.go` para rodar o cliente e obter a cotação do dólar.

Para interagir diretamente com o Banco de Dados SQLite, usar `sqlite3 ./cotacoes.db` no Terminal.
O comando `.tables` mostrará as tables existentes, e o comando `SELECT * FROM cotacoes;` mostrará
todas as cotações registradas do dólar. Cada cotação corresponde a uma chamada ao endpoint em
`localhost:8080/cotacao`.