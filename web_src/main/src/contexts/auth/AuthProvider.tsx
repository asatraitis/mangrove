import { useState } from "react"
import { AuthCtx, defaultAuthUser } from "./auth"

interface AuthProviderProps {
    children?: React.ReactNode
}
export const AuthProvider = ({children}:AuthProviderProps) => {
    const [user, setUser] = useState(defaultAuthUser)   
    return (
        <AuthCtx.Provider value={{user, setUser}}>
            {children}
        </AuthCtx.Provider>
    )
}