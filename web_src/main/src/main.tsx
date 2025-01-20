import { StrictMode } from 'react'
import ReactDOM from 'react-dom/client'
import { MantineProvider } from '@mantine/core';
import "@mantine/core/styles.css"

import { AuthProvider } from './contexts/auth/AuthProvider'
import App from './app'

const rootElement = document.getElementById('root')!
if (!rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement)
  root.render(
    <StrictMode>
      <AuthProvider>
        <MantineProvider defaultColorScheme="auto">
          <App />
        </MantineProvider>
      </AuthProvider>
    </StrictMode>,
  )
}