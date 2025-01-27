import { createLazyFileRoute } from '@tanstack/react-router'

export const Route = createLazyFileRoute('/clients')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/clients"!</div>
}
