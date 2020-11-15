export class Room {
    UId: string
    title: string
    contentMd: string
    imageUId: string


    constructor(UId: string, title: string, contentMd: string, imageUId: string) {
        this.UId = UId;
        this.title = title;
        this.contentMd = contentMd;
        this.imageUId = imageUId;
    }
}