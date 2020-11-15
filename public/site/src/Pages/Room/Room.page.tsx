import React from "react";
import { Helmet } from "react-helmet";
import {TerminalCompoment} from "../../components/TerminalCompoment";
import {GetRoomByUID} from "../../services/API";

import "./Room.page.css"

export const RoomPage = () => {
    const [ RoomId ] = React.useState(window.location.pathname.replace("/app/room/", ""))
    const [ Room, setRoom ] = React.useState<{title: string, imageUId: string } | null>(null)

    if (Room === null) {
        GetRoomByUID(RoomId).then(s => {
            setRoom(s)
        })
    }

    return <div>
        <Helmet>
            <title>{Room?.title}</title>
        </Helmet>
        <div id="tutorial">
            <iframe title={Room?.title} id="tutorial" height="100%" width="47%" src={"/md/"+RoomId} />
        </div>
        <div id="terminal">
            {(() => {
                if (Room !== null) {
                    return <TerminalCompoment Uid={Room.imageUId} />
                }
            })()}
        </div>
    </div>;
}