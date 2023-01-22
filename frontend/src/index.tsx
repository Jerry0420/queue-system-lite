import React from 'react'
import ReactDom from 'react-dom'
import { HashRouter, Route, Routes } from 'react-router-dom'
import { NoMatch } from './components/Default'
import { SignUp } from './components/Signup'
import { useApiRequest } from "./apis/reducer"
import { refreshToken } from './apis/StoreAPIs'
import { RefreshTokenContext } from './components/contexts'
import { StoreInfo } from './components/Store'
import { checkAuthFlow, validateResponseSuccess } from './apis/helper'
import { SignIn } from './components/Signin'
import { CreateCustomers } from './components/Customers'
import CssBaseline from '@mui/material/CssBaseline'
import { ThemeProvider } from '@emotion/react'
import { theme } from './styles/theme'
import { UpdatePasswordComponent } from './components/UpdatePassword'
// import './tailwind.scss'
// import './styles/style.scss'

function App() {
  
  const [refreshTokenAction, makeRefreshTokenRequest] = useApiRequest(...refreshToken())

  const wrapCheckAuthFlow = (nextStuff: () => void, redirectToMainPage: () => void) => {
    checkAuthFlow(refreshTokenAction.response, makeRefreshTokenRequest, 
      // nextStuff
      () => {
        if (validateResponseSuccess(refreshTokenAction.response) === true) {
          nextStuff()
        }
      }, 
      // redirectToMainPage
      () => {
        redirectToMainPage()
      }
    )
  }
  
  return (
    <ThemeProvider theme={theme}>
      <RefreshTokenContext.Provider value={{refreshTokenAction, makeRefreshTokenRequest, wrapCheckAuthFlow}}>
        <CssBaseline />
        <HashRouter>
          <Routes>
            <Route path="/" element={<SignUp />} />
            <Route path="/stores/:storeId/sessions/:sessionId" element={(<CreateCustomers />)} />
            <Route path="/stores/:storeId/password/update" element={(<UpdatePasswordComponent />)} />
            <Route path="/stores/:storeId" element={(<StoreInfo />)} />
            <Route path="/signin" element={(<SignIn />)} />
            <Route path="*" element={<NoMatch />} />
          </Routes>
        </HashRouter>
      </RefreshTokenContext.Provider>
    </ThemeProvider>
  )
}

ReactDom.render(<App />, document.getElementById('root'))
