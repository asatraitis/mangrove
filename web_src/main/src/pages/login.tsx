import { Card, Button, Loader, TextInput } from '@mantine/core';
import classes from './login.module.css';
import { FormEvent, useState } from 'react';
import { useRouter, getRouteApi } from '@tanstack/react-router';

import { startAuth } from '@services/auth/auth';

import { apiClient as api} from '@services/apiClient/apiClient'
import { useAuthCtx } from '../contexts/auth/useAuthCtx';

const routeApi = getRouteApi('/login')

export default function Login() {
    const search = routeApi.useSearch()

    const router = useRouter()
    const {setUser} = useAuthCtx()
    const [loading, setLoading] = useState(false)

    const [username, setUsername] = useState("")

    const handleAuth = async (e: FormEvent) => {
        e.preventDefault()
        setLoading(true)
        console.log(username)
        // hit the api and get options to be used with authenticator
        const {response, error: initLoginErr} = await api.initLogin(username)
        if (initLoginErr) {
            // TODO: handle error
            console.error("initLogin API error", initLoginErr)
            return
        }
        if (!response) {
            // TODO: handle error
            console.error("initLogin API returned null response")
            return
        }

        const [credential, startAuthErr] = await startAuth(response.publicKey)
        if (startAuthErr) {
            // TODO: handler err
            console.error("authenticator error", startAuthErr)
            return
        }
        const {response: finishLoginRes, error: finishLoginErr} = await api.finishLogin({credential, sessionKey: response.sessionKey})
        if (finishLoginErr) {
            // TODO: handle error
            console.error("finishLogin API error", finishLoginErr)
            return
        }
        if (!finishLoginRes) {
            // TODO: handle error
            console.error("initLogin API returned null response")
            return
        }
        setUser(finishLoginRes)
        setLoading(false)
        router.history.push(search.redirect)
    }
    
    return (
    <div className={classes.container}>
        <Card withBorder p="xl" radius="lg" style={{display: "flex", alignItems: "center"}}>
                <div style={{ width: "150px" }}>
                    <svg stroke="currentColor" fill="currentColor" stroke-width="0" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                        <path d="M12 2c5.523 0 10 4.477 10 10a10 10 0 0 1 -20 0c0 -5.523 4.477 -10 10 -10zm2 5a3 3 0 0 0 -2.98 2.65l-.015 .174l-.005 .176l.005 .176c.019 .319 .087 .624 .197 .908l.09 .209l-3.5 3.5l-.082 .094a1 1 0 0 0 0 1.226l.083 .094l1.5 1.5l.094 .083a1 1 0 0 0 1.226 0l.094 -.083l.083 -.094a1 1 0 0 0 0 -1.226l-.083 -.094l-.792 -.793l.585 -.585l.793 .792l.094 .083a1 1 0 0 0 1.403 -1.403l-.083 -.094l-.792 -.793l.792 -.792a3 3 0 1 0 1.293 -5.708zm0 2a1 1 0 1 1 0 2a1 1 0 0 1 0 -2z"></path>
                    </svg>
                </div>
                <form onSubmit={handleAuth} style={{display: "flex", flexDirection: "column"}}>
                    <TextInput value={username} onChange={(e) => {setUsername(e.target.value)}} label="Username" mt={30} />
                    <Button type='submit' style={{flexGrow: "1"}} variant='gradient' mt={5}>
                        { loading ? <Loader color="white" type="dots" /> : "Authenticate"}        
                    </Button>
                </form>     
        </Card>
    </div>
        
    )
}
