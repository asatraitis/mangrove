import { useContext } from "react";
import { AuthCtx } from "./auth";

export const useAuthCtx = () => useContext(AuthCtx)
