import React from "react";

export function UpdateDockerImages() {
    return fetch("/api/ListImages/", {
        method: "POST",
    }).then((response) => {
        if (response.status === 200) {
            return response.json()
        }
        throw response.status
    })
}

export function CreateDockerImage(name: string, dockerfile: string) {
    let header = new Headers()
    header.append("image", JSON.stringify({
        Name: name,
        Dockerfile: dockerfile,
    }))
    return fetch("/api/CreateImage/", {
        method: "POST",
        headers: header,
        cache: "no-cache"
    }).then(response => {
        if (response.status === 200) {
            return UpdateDockerImages()
        }
        throw response.status
    })
}

export function UpdateRooms() {
    return fetch("/api/ListRooms/", {
        method: "POST",
    }).then((response) => {
        if (response.status === 200) {
            return response.json()
        }
        throw response.status
    })
}

export function CreateRoom(title: string, markdown: string, imageUId: string) {
    let header = new Headers()
    header.append("ROOM", JSON.stringify({
        title: title,
        contentMd: markdown,
        imageUId: imageUId,
    }))
    return fetch("/api/CreateRoom/", {
        method: "POST",
        headers: header,
        cache: "no-cache"
    }).then(response => {
        if (response.status === 200) {
            return UpdateRooms()
        }
        throw response.status
    })
}

export function GetRoomByUID(RoomId: string) {
    return fetch("/api/GetRoomByID/"+RoomId, {
        method: "POST",
    }).then(response => {
        if (response.status === 200) {
            return response.json()
        }
        throw response.status
    })
}
