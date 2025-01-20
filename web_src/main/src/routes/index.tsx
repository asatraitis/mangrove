import { createFileRoute, redirect } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  beforeLoad: ({context, location}) => {
    if (context.auth.status !== "active") {
      throw redirect({
        to: '/login',
        search: {
          redirect: location.href
        }
      })
    }
  },
})
