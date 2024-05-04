package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/guilhermeonrails/api-go-gin/controllers"
	"github.com/guilhermeonrails/api-go-gin/database"
	"github.com/guilhermeonrails/api-go-gin/models"
	"github.com/stretchr/testify/assert"
)

func SetupRotasDeTeste() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	return r
}

func Test_Saudacao_VerificaStatusCode(t *testing.T) {
	r := SetupRotasDeTeste()
	r.GET("/:nome", controllers.Saudacao)
	req, _ := http.NewRequest("GET", "/gui", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code, "Deveria ser Status OK (200)")

	mockResp := `{"API diz:":"E ai gui, tudo beleza?"}`
	respBody, _ := io.ReadAll(resp.Body)
	assert.Equal(t, mockResp, string(respBody), "Conteúdo do Body inválido")
}

func Test_Alunos_ListarTodos(t *testing.T) {
	database.ConectaComBancoDeDados()
	aluno := CriaAlunoMock()
	defer ExcluirAlunoMock(&aluno)

	r := SetupRotasDeTeste()
	r.GET("/alunos", controllers.ExibeTodosAlunos)
	req, _ := http.NewRequest("GET", "/alunos", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func CriaAlunoMock() models.Aluno {
	aluno := models.Aluno{
		Nome: "Aluno Mock",
		CPF:  "12345678901",
		RG:   "123456789",
	}
	database.DB.Create(&aluno)
	fmt.Println("ALUNO MOCK: ", aluno.ID)
	return aluno
}

func ExcluirAlunoMock(aluno *models.Aluno) {
	database.DB.Delete(&aluno)
}

func Test_Alunos_BuscarPorCpf(t *testing.T) {
	database.ConectaComBancoDeDados()
	aluno := CriaAlunoMock()
	defer ExcluirAlunoMock(&aluno)

	r := SetupRotasDeTeste()
	r.GET("/alunos/cpf/:cpf", controllers.BuscaAlunoPorCPF)
	req, _ := http.NewRequest("GET", "/alunos/cpf/"+aluno.CPF, nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func Test_Alunos_BuscarPorId(t *testing.T) {
	database.ConectaComBancoDeDados()
	aluno := CriaAlunoMock()
	defer ExcluirAlunoMock(&aluno)

	r := SetupRotasDeTeste()
	r.GET("/alunos/:id", controllers.BuscaAlunoPorID)
	req, _ := http.NewRequest("GET", "/alunos/"+fmt.Sprint(aluno.ID), nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	var alunoMock models.Aluno
	json.Unmarshal(resp.Body.Bytes(), &alunoMock)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, aluno.ID, alunoMock.ID)
	assert.Equal(t, aluno.Nome, alunoMock.Nome)
	assert.Equal(t, aluno.CPF, alunoMock.CPF)
	assert.Equal(t, aluno.RG, alunoMock.RG)
}

func Test_Alunos_Exclusao(t *testing.T) {
	database.ConectaComBancoDeDados()
	aluno := CriaAlunoMock()

	r := SetupRotasDeTeste()
	r.DELETE("/alunos/:id", controllers.DeletaAluno)
	req, _ := http.NewRequest("DELETE", "/alunos/"+fmt.Sprint(aluno.ID), nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func Test_Alunos_Atualizasao(t *testing.T) {
	database.ConectaComBancoDeDados()
	aluno := CriaAlunoMock()
	defer ExcluirAlunoMock(&aluno)

	aluno = models.Aluno{
		Nome: "Aluno Mock2",
		CPF:  "12345678902",
		RG:   "123456782",
	}

	alunoJson, _ := json.Marshal(&aluno)

	r := SetupRotasDeTeste()
	r.PATCH("/alunos/:id", controllers.EditaAluno)
	req, _ := http.NewRequest("PATCH", "/alunos/"+fmt.Sprint(aluno.ID), bytes.NewBuffer(alunoJson))
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	var alunoResp models.Aluno
	json.Unmarshal(resp.Body.Bytes(), &alunoResp)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, aluno.ID, alunoResp.ID)
	assert.Equal(t, aluno.Nome, alunoResp.Nome)
	assert.Equal(t, aluno.CPF, alunoResp.CPF)
	assert.Equal(t, aluno.RG, alunoResp.RG)
}
