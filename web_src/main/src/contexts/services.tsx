import { createContext, ReactNode, useContext } from "react";
import { IApiClient } from "@services/apiClient/apiClient";
import ApiClient from "@services/apiClient/apiClient"

export interface IServicesCtx {
    api: IApiClient
}

const ServicesCtx = createContext<IServicesCtx | null>(null)
export const useServices = () => useContext(ServicesCtx)

export const ServicesProvider = ({children}: {children: ReactNode}) => {
    return (
        <ServicesCtx.Provider value={{
            api: new ApiClient("")
        }}>{children}</ServicesCtx.Provider>
    )
}
