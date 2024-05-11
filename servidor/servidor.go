package servidor

import (
	"crud/banco"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type usuario struct {
	ID    string `json:"id"`
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

	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar com o banco de dados"))
		return
	}

	statement, erro := db.Prepare("insert into usuario(nome, email) values (?, ?)")
	if erro != nil {
		w.Write([]byte("Erro ao criar statement!"))
		return
	}

	_, erro = statement.Exec(usuario.Nome, usuario.Email)
	if erro != nil {
		w.Write([]byte(erro.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Usuario criado com sucesso"))
}

func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar com o banco de dados"))
		return
	}
	defer db.Close()

	linhas, erro := db.Query("select * from usuario")
	if erro != nil {
		w.Write([]byte("Erro ao buscar usuários"))
		return
	}
	defer linhas.Close()

	var usuarios []usuario
	for linhas.Next() {
		var usuario usuario

		if erro := linhas.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); erro != nil {
			w.Write([]byte("Erro ao escanear usuário"))
			return
		}

		usuarios = append(usuarios, usuario)
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if erro := json.NewEncoder(w).Encode(usuarios); erro != nil {
		w.Write([]byte("Erro ao converter usuarios para JSON"))
		return
	}

}

func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	ID := parametros["id"]

	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar com o banco de dados"))
		return
	}

	linha, erro := db.Query("select * from usuario where id = ?", ID)
	if erro != nil {
		w.Write([]byte("Erro ao buscar o usuario"))
		return
	}

	var usuario usuario
	if linha.Next() {
		if erro := linha.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); erro != nil {
			w.Write([]byte("Erro ao escanear usuário"))
			return
		}
	}

	if usuario.ID == "" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if erro := json.NewEncoder(w).Encode(usuario); erro != nil {
		w.Write([]byte("Erro ao converter usuário"))
		return
	}

}

func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	ID := parametros["id"]

	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Erro ao ler corpo da requisição"))
		return
	}

	var usuario usuario
	if erro := json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erro ao converter usuário para struct"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erro ao conectar com o banco de dados"))
		return
	}
	defer db.Close()

	statement, erro := db.Prepare("update usuario set nome = ?, email = ? where id = ?")
	if erro != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erro ao criar statement"))
		fmt.Println(erro)
		return
	}
	defer statement.Close()

	if _, erro := statement.Exec(usuario.Nome, usuario.Email, ID); erro != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erro ao executar statement"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	ID := parametros["id"]

	db, erro := banco.Conectar()
	if erro != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erro ao conectar com o banco de dados"))
		return
	}
	defer db.Close()

	statement, erro := db.Prepare("delete from usuario where id = ?")
	if erro != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erro ao criar statement"))
		return
	}
	defer statement.Close()

	if _, erro := statement.Exec(ID); erro != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erro ao executar statement"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
