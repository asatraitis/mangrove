import { Link, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/router-devtools";
import { useAuthCtx } from "../contexts/auth/useAuthCtx";

export default function Index() {
    const {user} = useAuthCtx()
    if (user.status !== 'active') {
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