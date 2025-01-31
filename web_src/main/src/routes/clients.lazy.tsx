import { createLazyFileRoute } from '@tanstack/react-router'

import Clients from "../pages/client"

export const Route = createLazyFileRoute('/clients')({
  component: Clients,
})
