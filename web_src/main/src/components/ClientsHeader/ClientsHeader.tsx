import { Button, Flex, Title } from "@mantine/core";

export type ActionType = "create"
export type Props = {
    onAction?: (type: ActionType) => void
}

export default function ClientsHeader({onAction = () => {}}: Props) {
    return (
        <Flex justify="space-between">
            <Title order={1}>Clients</Title>
            <Button onClick={() => {onAction("create")}}>Create Client</Button>
        </Flex>
    )
}