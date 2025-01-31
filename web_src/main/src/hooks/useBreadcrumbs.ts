import { useLocation } from "@tanstack/react-router";

export function useBreadcrumbs() {
    const current = useLocation()
    const routeHistory = current.pathname.split("/").filter(x => x !== "")

    return routeHistory.map((route, idx) => ({name: route, to: `${idx === 0 ? '/'+route : '/'+routeHistory[idx-1]+'/'+route}`}))
}