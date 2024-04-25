package servidor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type usuario struct {
	ID    uint32 `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		w.Write([]byte("Falha ao ler corpo da requisição"))
	}

	var usuario usuario

	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		w.Write([]byte("Erro ao conevrter para struct de usuario"))
	}

	fmt.Println(usuario)
}
