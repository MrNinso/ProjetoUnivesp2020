export class DockerImage {
    UId: string
    Name: string
    Created: string

    constructor(UId: string, Name: string, Created: string) {
        this.UId = UId;
        this.Name = Name;
        this.Created = Created;
    }
}