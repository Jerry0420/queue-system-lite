import {createContext} from 'react'
import { initialState, Action } from '../apis/reducer'
import { JSONResponse } from '../apis/reducer'

const initialRefreshTokenContext: {
    refreshTokenAction: Action, 
    makeRefreshTokenRequest: (() => Promise<JSONResponse | null | undefined>),
    wrapCheckAuthFlow: ((nextStuff: () => void, redirectToMainPage: () => void) => void)
} = {
    refreshTokenAction: initialState,
    makeRefreshTokenRequest: (() => {return new Promise((resolve, reject) => {})}),
    wrapCheckAuthFlow: ((nextStuff: () => void, redirectToMainPage: () => void) => {})
}

const RefreshTokenContext = createContext(initialRefreshTokenContext)

export {
    RefreshTokenContext
}