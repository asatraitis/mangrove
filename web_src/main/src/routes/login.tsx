import { createFileRoute, redirect } from '@tanstack/react-router'
import { USER_STATUS_ACTIVE, USER_STATUS_INACTIVE, USER_STATUS_PENDING, USER_STATUS_SUSPENDED } from '@dto/types'

type LoginSearch = {
  redirect: string
}
export const Route = createFileRoute('/login')({
  validateSearch: (search: Record<string, unknown>): LoginSearch => {
    return {
      redirect: (search.redirect as string) || "/"
    }
  },
  beforeLoad: async ({context: {setUser, api}, search}) => {
    const {response, error} = await api.me()
    if (error) {
      // handle error
      return
    }
    if (!response) {
      // TODO: handle no response
      return
    }

    // TODO: replace with switch case
    if (response.status === USER_STATUS_INACTIVE) {
      // TODO: handle inactive user status
      return
    }
    if (response.status === USER_STATUS_SUSPENDED) {
      // TODO: handle suspended user status
      return
    }
    if (response.status === USER_STATUS_PENDING) {
      // TODO: handler pending user status
      return
    }
    if (response.status === USER_STATUS_ACTIVE) {
      setUser(response)
      throw redirect({
        to: search.redirect
      })
    }
  }
})
