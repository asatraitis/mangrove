import { Card, Button, Loader, TextInput } from '@mantine/core';
import { TbCircleKeyFilled } from "react-icons/tb";
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
                <TbCircleKeyFilled size={150} />
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
