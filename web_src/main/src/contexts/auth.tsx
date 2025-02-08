import { createContext, useContext, useState } from "react";
import {MeResponse} from "@dto/types"

export const defaultAuthUser: MeResponse = {displayName: "", id: "", role: "", status: ""}
interface IAuthCtx {
    user: MeResponse;
    setUser: React.Dispatch<React.SetStateAction<MeResponse>>
}
interface AuthProviderProps {
    children?: React.ReactNode
}

export const AuthCtx = createContext<IAuthCtx>({user: defaultAuthUser, setUser: () => {}})
export const useAuthCtx = () => useContext(AuthCtx)

export const AuthProvider = ({children}:AuthProviderProps) => {
    const [user, setUser] = useState(defaultAuthUser)   
    return (
        <AuthCtx.Provider value={{user, setUser}}>
            {children}
        </AuthCtx.Provider>
    )
}
