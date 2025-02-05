import { Link } from "@tanstack/react-router";
import { useBreadcrumbs } from "../../hooks/useBreadcrumbs";
import { Breadcrumbs as MBreadcrumbs } from "@mantine/core";

export default function Breadcrumbs() {
    const bc = useBreadcrumbs()
    
    return <MBreadcrumbs mb="md" separator="/" separatorMargin="xs">
            {
                bc.map(item => (
                    <Link key={item.name} to={item.to}>{item.name}</Link>
                ))
            }
            </MBreadcrumbs>
}