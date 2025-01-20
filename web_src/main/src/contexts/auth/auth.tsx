import { createContext } from "react";
import {MeResponse} from "@dto/types"

export const defaultAuthUser: MeResponse = {displayName: "", id: "", role: "", status: ""}
interface IAuthCtx {
    user: MeResponse;
    setUser: React.Dispatch<React.SetStateAction<MeResponse>>
}
export const AuthCtx = createContext<IAuthCtx>({user: defaultAuthUser, setUser: () => {}})
