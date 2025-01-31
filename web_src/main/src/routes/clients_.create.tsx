import { createFileRoute, redirect } from '@tanstack/react-router'

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
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/clients/new"!</div>
}
