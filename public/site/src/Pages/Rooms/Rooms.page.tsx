import React from "react";

var DockerFile : HTMLTextAreaElement | null
var NomeImagem : HTMLInputElement | null

function criarSala(http: XMLHttpRequest) {
    if (DockerFile != null && NomeImagem != null) {
        http.open("POST", "/api/CreateImage/", false)

        http.setRequestHeader("IMAGE", JSON.stringify({
            Name: NomeImagem.value,
            DockerFile: DockerFile.value,
        }))

        http.send()

        if (http.status === 200) {
            alert("Imagen Criada")
        } else {
            alert("erro ao criar: ["+http.status+"] "+ http.response)
        }


    }
}

export const RoomsPage = (http: XMLHttpRequest) => {
    return (
        <div>
            <h1>Demo da API</h1>
            <fieldset>
                <legend>Criar Imagem</legend>
                <label>
                    Nome da Imagem:
                    <br/>
                    <input type="text" ref={instance => NomeImagem = instance} />
                </label>
                <br/>
                <br/>
                <label>
                    Tutorial em Docker file:
                    <br/>
                    <textarea ref={instance => DockerFile = instance} />
                </label>
                <br/>
                <div>
                    <button onClick={(e) => criarSala(http)}>Criar Imagem</button>
                    <button>Listar Imagens</button>
                </div>
            </fieldset>
        </div>
    );
}

//TODO LISTAR Imagens
//TODO Atualizar Imagens
//TODO Deletar Imagem

//TODO Criar USUARIO
//TODO Atualizar Usuario
//TODO Deletar Usuario
//TODO Listar Usuarios

//TODO Criar Sala
//TODO Listar Sala
//TODO Atulizar Sala
//TODO Deletar Sala
//TODO Renderizar Sala
//TODO Renderizar TODAS AS SALAS
//TODO LiveRoomRender