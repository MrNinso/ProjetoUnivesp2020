import React, {FormEvent} from "react";
import { History } from 'history'
import Cookies from 'js-cookie';
import { sha256 } from 'js-sha256'

function sendEmail(event: FormEvent<HTMLFormElement>) {
        if (SenhaInput !== null &&  EmailInput !== null) {
            SenhaInput.value = require('md5')(SenhaInput.value)
            SenhaInput.value = sha256(EmailInput.value+"_&_"+SenhaInput.value)

            event.persist()
        }
}

var EmailInput: HTMLInputElement | null
var SenhaInput: HTMLInputElement | null

export const LoginPage = (http: XMLHttpRequest, history: History<any>) => {
    let email = Cookies.get('3ic7k5irhh2az9hkig1oy3')
    let token = Cookies.get('97b31ae2cd1a382f19a7b95f5ef98016')


    var invalidCookie = false //TODO MOSTRAR ERRO

    if (email !== undefined && token !== undefined) {
        http.open("POST", "/api/Login/")

        http.onload = () => {
            if (http.status === 200) {
                history.push("/app/rooms")
            } else {
                invalidCookie = true
            }
        }

        http.send()

    }

    return (
        <form method="POST" action="/api/Login/test" onSubmit={(e) => sendEmail(e)}>
            <label>
                Email:
                <input ref={i => EmailInput = i} type="email" name="Email" required/>
            </label>
            <label>
                Senha:
                <input ref={i => SenhaInput = i} type="password" name="Token" required/>
            </label>

            <input type="submit" value="enviar"/>
        </form>
    );
}
