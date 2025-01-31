import { Link, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/router-devtools";
import { LuFileKey } from "react-icons/lu";
import { RiDoorLockLine } from "react-icons/ri";

import { AppShell, Burger, Flex, Container } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";

import { useAuthCtx } from "../contexts/auth/useAuthCtx";
import classes from "./index.module.css"

export default function Index() {
    const {user} = useAuthCtx()
    const [opened, {toggle}] = useDisclosure()



    if (user.status !== 'active') {
        return (
            <div>
                {'root?'}
                <Outlet />
            </div>
    )
    }
    
    return (
    <AppShell
        header={{ height: 60 }}
        navbar={{
            width: 250,
            breakpoint: 'sm',
            collapsed: { mobile: !opened },
        }}
        padding="md"
    >
        <AppShell.Header >
            <Flex className={classes.header} align="center">
                <Burger
                    opened={opened}
                    onClick={toggle}
                    hiddenFrom="sm"
                    size="sm"
                />
                <RiDoorLockLine className={classes.logoIcon} />
            </Flex>
        </AppShell.Header>
        <AppShell.Navbar p="md">
            <Link to="/clients" className={classes.link}>
                <LuFileKey className={classes.linkIcon} />
                <span>Clients</span>
            </Link>
        </AppShell.Navbar>
        <AppShell.Main>
            <Container size="xl">
                <Outlet />
            </Container>
        </AppShell.Main>
        <TanStackRouterDevtools />
    </AppShell>

    )
}