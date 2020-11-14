import React, { useEffect } from "react";

function updateImagesList(http: XMLHttpRequest, setListImagens: React.Dispatch<React.SetStateAction<any>>) {
    http.open("POST", "/api/ListImages/", false)

    http.send()

    if (http.status === 200) {
        console.log("Lista de Imagens recebida", http.response)
        setListImagens(JSON.parse(http.response).map((imagem: any) => {
            return { UId: imagem.UId, Name: imagem.Name, Created: imagem.Created }
        }))
    } else {
        console.log("erro ao carregar as imagens: ["+http.status+"] "+ http.response)
    }
}

function updateRoomsList(http: XMLHttpRequest, setListSalas: React.Dispatch<React.SetStateAction<any>>) {
    http.open("POST", "/api/ListRooms/", false)

    http.send()

    if (http.status === 200) {
        console.log("Lista de Sala recebida: "+ http.response)
        setListSalas(JSON.parse(http.response).map((sala: any) => {
            return {UId: sala.UId, title: sala.title, contentMd: sala.contentMd, imageUId: sala.imageUId}
        }))
    } else {
        console.error("erro ao carregar as salas: ["+http.status+"] "+ http.response)
    }
}

function updateUsersList(http: XMLHttpRequest) {
    console.log("TODO updateUsersList")
}

function createImage(http: XMLHttpRequest, NomeImagem: string, Dockerfile: string, setListImagens: React.Dispatch<React.SetStateAction<any>>) {
    http.open("POST", "/api/CreateImage/", false)

    http.setRequestHeader("IMAGE", JSON.stringify({
        Name: NomeImagem,
        Dockerfile: Dockerfile,
    }))

    http.send()

    if (http.status === 200) {
        console.log("Imagem Criada")
        updateImagesList(http, setListImagens)
    } else {
        console.error("erro ao criar: ["+http.status+"] "+ http.response)
    }
}

function createRoom(http: XMLHttpRequest, TituloSala: string, ConteudoMarkdown: string, ImagemSelecionada: string, setListSalas: React.Dispatch<React.SetStateAction<any>>) {
    console.log(":"+ImagemSelecionada) //TODO REMOVER
    http.open("POST", "/api/CreateRoom/", false)

    http.setRequestHeader("ROOM", JSON.stringify({
        UId: "",
        title: TituloSala,
        contentMd: ConteudoMarkdown,
        imageUId: ImagemSelecionada,
    }))

    http.send()

    if (http.status === 200) {
        console.log("Sala Criada")
        updateRoomsList(http, setListSalas)
    } else {
        console.error("erro ao criar: ["+http.status+"] "+ http.response)
    }
}

function listImages(http: XMLHttpRequest, ListaImagens: any, setListImagens: any) {
    updateImagesList(http, setListImagens)
    console.log(ListaImagens)
}

function listRooms(http: XMLHttpRequest, ListaSalas: any, setListSalas: any) {
    updateRoomsList(http, setListSalas)
    console.log(ListaSalas)
}

export const RoomsPage = (http: XMLHttpRequest) => {
    const [ NomeImagem, setNomeImagem ] = React.useState("")
    const [ Dockerfile, setDockerfile ] = React.useState("")
    const [ ListaImagens, setListImagens ] = React.useState<{UId: string, Name: string, Created: bigint} | any>([])

    const [ TituloSala, setTituloSala ] = React.useState("")
    const [ ImagemSelecionada, setImagemSelecionada ] = React.useState("")
    const [ ConteudoMarkdown, setConteudoMarkdown ] = React.useState("")
    const [ ListaSalas, setListaSalas ] = React.useState<{UId: string, title: string, contentMd: string, imageUId: string} | any>([])

    useEffect(() => {
        updateImagesList(http, setListImagens)
        updateRoomsList(http, setListaSalas)
    }, [http, setListImagens, setListaSalas])
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
                    <button onClick={() => createImage(http, NomeImagem, Dockerfile, setListImagens)}>Criar Imagem</button>
                    <button onClick={() => listImages(http, ListaImagens ,setListImagens)} >Listar Imagens</button>
                </div>
            </fieldset>
            <fieldset>
                <legend>Criar Sala</legend>
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
                        onClick={() => createRoom(http, TituloSala, ConteudoMarkdown, ImagemSelecionada, setListaSalas)}>
                            Criar Sala
                    </button>
                    <button onClick={() => listRooms(http, ListaSalas, setListaSalas)}>Listar Salas</button>
                </div>
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

//TODO Criar Sala
//TODO Listar Sala
//TODO Atulizar Sala
//TODO Deletar Sala
//TODO Renderizar Sala
//TODO Renderizar TODAS AS SALAS
//TODO LiveRoomRender