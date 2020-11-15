import React, { useEffect } from "react";
import { DockerImage } from "../../interfaces/DockerImage";
import {Room} from "../../interfaces/Room";
import { UpdateRooms, CreateDockerImage, CreateRoom, UpdateDockerImages } from "../../services/API";

import "./Demo.page.css"

export const DemoPage = () => {
    const [ NomeImagem, setNomeImagem ] = React.useState("")
    const [ Dockerfile, setDockerfile ] = React.useState("")
    const [ ListaImagens, setListImagens ] = React.useState<DockerImage | any>([])

    const [ TituloSala, setTituloSala ] = React.useState("")
    const [ ImagemSelecionada, setImagemSelecionada ] = React.useState("")
    const [ ConteudoMarkdown, setConteudoMarkdown ] = React.useState("")
    const [ ListaSalas, setListaSalas ] = React.useState<Room | any>([])

    const [ FirstUpdate, setFirstUpdate ] = React.useState<boolean>(false)

    if (!FirstUpdate) {
        UpdateDockerImages().catch((reason: any) => {
            console.error("UpdateDockerImages: ",reason)
        }).then(s => {
            console.log("UpdateDockerImages: ",s)
            setListImagens(s)
            return Promise.resolve()
        }).finally(() => {
            console.log("UpdateDockerImages Finish")
        })
        UpdateRooms().catch((reason: any) => {
            console.error("UpdateRooms: ",reason)
        }).then(s => {
            console.log("UpdateRooms: ",s)
            setListaSalas(s)
            return Promise.resolve()
        }).finally(() => {
            console.log("UpdateRooms Finish")

        })
        setFirstUpdate(true)
    }

    return (
        <div>
            <h1>Demo da API</h1>
            <fieldset>
                <legend>Criar Imagem</legend>
                <label>
                    Nome da Imagem:
                    <br/>
                    <input type="text"
                           value={ NomeImagem }
                           onChange={event => setNomeImagem(event.target.value)} />
                </label>
                <br/><br/>
                <label>
                    Dockerfile:
                    <br/>
                    <textarea
                        value={ Dockerfile }
                        onChange={event => setDockerfile(event.target.value)}/>
                </label>
                <br/>
                <div>
                    <button onClick={() => CreateDockerImage(NomeImagem, Dockerfile).then(s => {
                        console.log("Imagem criada")
                        setListImagens(s)
                    })}>Criar Imagem</button>
                    <button onClick={() => UpdateDockerImages().then(s => {setListImagens(s)
                        console.log(s)})} >Listar Imagens</button>
                </div>
            </fieldset>
            <fieldset>
                <legend>Salas</legend>
                <fieldset id="div1">
                    <legend>Criar</legend>
                    <label>
                        Titulo da sala:
                        <br/>
                        <input type="text"
                               value={TituloSala}
                               onChange={event => setTituloSala(event.target.value)}/>
                    </label>
                    <br/><br/>
                    <label>
                        Imagem:
                        <br/>
                        <select onChange={event => setImagemSelecionada(event.target.value)}>
                            <option value="" />
                            {ListaImagens.map((imagem: any) =><option value={imagem.UId}>{imagem.Name}</option>)}
                        </select>
                    </label>
                    <br/><br/>
                    <label>
                        Conteudo em markdown:
                        <br/>
                        <textarea
                            value={ ConteudoMarkdown }
                            onChange={event => setConteudoMarkdown(event.target.value)}/>
                    </label>
                    <div>
                        <button
                            onClick={() => CreateRoom(TituloSala, ConteudoMarkdown, ImagemSelecionada).then(s => {
                                console.log("Sala Criada")
                                setListaSalas(s)
                            })}>
                                Criar Sala
                        </button>
                        <button onClick={() => UpdateRooms().then(console.log)}>Listar Salas</button>
                    </div>
                </fieldset>
                <fieldset>
                    <legend>Salas</legend>
                    <div id="grid-container">
                        {ListaSalas.map((sala: Room) => {
                            return <a id="grid-item" href={"/app/room/"+sala.UId} >{sala.title}</a>
                        })}
                    </div>
                </fieldset>
            </fieldset>
        </div>
    );
}

//TODO Atualizar Imagens
//TODO Deletar Imagem

//TODO Criar USUARIO
//TODO Atualizar Usuario
//TODO Deletar Usuario
//TODO Listar Usuarios

//TODO Atulizar Sala
//TODO Deletar Sala
//TODO Renderizar Sala
//TODO Renderizar TODAS AS SALAS
//TODO LiveRoomRender