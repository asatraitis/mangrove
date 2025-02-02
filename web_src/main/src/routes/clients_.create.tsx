import { createFileRoute, redirect } from '@tanstack/react-router'
import CreateClient from '../pages/createClient'

export const Route = createFileRoute('/clients_/create')({
  beforeLoad: ({ context, location }) => {
    if (context?.user.status !== 'active') {
      throw redirect({
        to: '/login',
        search: {
          redirect: location.href,
        },
      })
    }
  },
  component: CreateClient,
})
