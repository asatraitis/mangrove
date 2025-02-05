import { getRouteApi, linkOptions, useNavigate } from '@tanstack/react-router'
import { Flex, Card, Button } from '@mantine/core';

import ClientsTable from "../components/ClientsTable/ClientsTable"
import EmptyState from '../components/shared/EmptyState';
import ClientsHeader, {ActionType} from '../components/ClientsHeader/ClientsHeader';

const createClientLinkOptions = linkOptions({
    to: "/clients/create",
  })
const clientApi = getRouteApi("/clients")
export default function Clients() {
    const navigate = useNavigate()
    const {response: clients, error} = clientApi.useLoaderData()

    if (error) {
        console.error(error.message)
    }

    const handleAction = (type: ActionType) => {
        switch(type) {
            case "create":
                navigate(createClientLinkOptions)
                break
            default:
                console.warn("action not defined for ", type)
            
        }
    }

    return (
            <Card radius="md" p='xl' withBorder>
                <Flex direction="column">
                {
                    clients?.length ? (
                        <>
                            <ClientsHeader onAction={handleAction}/>
                            <ClientsTable data={clients!} />
                        </>
                    ) : (
                        <EmptyState title="No created items yet" description="Created clients will appear below.">
                                <Button mt="md" onClick={() => {navigate(createClientLinkOptions)}}>Create Client</Button>
                        </EmptyState>
                    )
                }
                </Flex>
            </Card>
    )
}
