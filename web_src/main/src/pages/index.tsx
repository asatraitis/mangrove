import { useContext } from "react";
import { AuthCtx } from "../contexts/auth";
import { Link, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/router-devtools";

export default function Index() {
    const auth = useContext(AuthCtx);
    if (auth.status !== 'valid') {
        return <Outlet />
    }
    
    return (
    <>
        <div className="p-2 flex gap-2">
        <Link to="/" className="[&.active]:font-bold">
            Home
        </Link>{' '}
        </div>
        <hr />
        <Outlet />
        <TanStackRouterDevtools />
    </>

    )
}