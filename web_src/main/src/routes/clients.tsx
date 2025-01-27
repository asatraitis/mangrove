import { createFileRoute, redirect } from '@tanstack/react-router'

export const Route = createFileRoute('/clients')({
  beforeLoad: ({context, location}) => {
      if (context?.user.status !== "active") {
        throw redirect({
          to: '/login',
          search: {
            redirect: location.href
          }
        })
      }
    },
})

