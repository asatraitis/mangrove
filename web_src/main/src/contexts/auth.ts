import { createContext } from "react";

export type AuthUser = {
    status: "pending" | "valid" | "invalid"
    displayName: string
}
export const defaultAuthUser: AuthUser = {status: 'pending', displayName: ""}
export const AuthCtx = createContext<AuthUser>(defaultAuthUser)
