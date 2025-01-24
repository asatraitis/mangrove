import { createFileRoute, redirect } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
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
