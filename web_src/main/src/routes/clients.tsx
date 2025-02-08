import { createFileRoute, redirect } from '@tanstack/react-router'
import { Skeleton } from '@mantine/core'

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
  loader: ({context}) => {
    return context.api.userClients()
  },
  pendingComponent: Pending,
})

function Pending() {
  return (
    <>
      <Skeleton height={8} radius="xl" />
      <Skeleton height={8} mt={6} radius="xl" />
      <Skeleton height={8} mt={6} radius="xl" />
      <Skeleton height={8} mt={6} radius="xl" />
      <Skeleton height={8} mt={6} radius="xl" />
    </>
  )
}

