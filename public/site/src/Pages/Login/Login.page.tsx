import React, {FormEvent} from "react";
import { History } from 'history'
import Cookies from 'js-cookie';
import { sha256 } from 'js-sha256'

var invalidLogin = false

function sendLogin(event: FormEvent<HTMLButtonElement>, http: XMLHttpRequest, history: History<any>) {
        if (SenhaInput !== null &&  EmailInput !== null) {
            var t : string
            t = require('md5')(SenhaInput.value)
            t = sha256(EmailInput.value+"_&_"+t)

            http.open("POST", "/api/Login/")

            http.setRequestHeader("email", EmailInput.value)
            http.setRequestHeader("token", t)

            http.onload = () => {
                if (http.status === 200) {
                    history.push("/app/home")
                } else {
                    invalidLogin = true
                }
            }
            http.send()
        }
}

var EmailInput: HTMLInputElement | null
var SenhaInput: HTMLInputElement | null

export const LoginPage = (http: XMLHttpRequest, history: History<any>) => {
    let email = Cookies.get('3ic7k5irhh2az9hkig1oy3')
    let token = Cookies.get('97b31ae2cd1a382f19a7b95f5ef98016')


    var invalidCookie = false

    if (email !== undefined && token !== undefined) {
        http.open("POST", "/api/Login/", false)

        http.send()

        if (http.status === 200) {
            history.push("/app/home")
        } else {
            invalidCookie = true
            Cookies.remove('3ic7k5irhh2az9hkig1oy3')
            Cookies.remove('97b31ae2cd1a382f19a7b95f5ef98016')
        }

    }

    return (
        <div>
            <label>
                Email:<br/>
                <input ref={i => EmailInput = i} type="email" name="Email" required/>
            </label>
            <br/>
            <label>
                Senha:
                <br/>
                <input ref={i => SenhaInput = i} type="password" name="Token" required/>
            </label>
            <br/>
            <button onClick={(e) => sendLogin(e, http, history)}>Logar</button>
            <br/>
            <h2 hidden={!invalidLogin}>usuario ou senha invalido</h2>
            <br/>
            <h2 hidden={!invalidCookie}>Sessão espirada</h2>
        </div>
    );
}
